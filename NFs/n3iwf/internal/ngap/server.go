package ngap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"

	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/n3iwf/internal/ngap/message"
	"github.com/free5gc/n3iwf/pkg/factory"
	lib_ngap "github.com/free5gc/ngap"
	"github.com/free5gc/sctp"
	"github.com/free5gc/util/safe_channel"
)

const (
	RECEIVE_NGAPPACKET_CHANNEL_LEN = 512
	RECEIVE_NGAPEVENT_CHANNEL_LEN  = 512
)

type n3iwf interface {
	Config() *factory.Config
	Context() *n3iwf_context.N3IWFContext
	CancelContext() context.Context

	SendIkeEvt(n3iwf_context.IkeEvt)
}

type Server struct {
	n3iwf

	conn     []*sctp.SCTPConn
	rcvPktCh *safe_channel.SafeCh[ReceiveNGAPPacket]
	rcvEvtCh *safe_channel.SafeCh[n3iwf_context.NgapEvt]
}

type ReceiveNGAPPacket struct {
	Conn *sctp.SCTPConn
	Buf  []byte
}

func NewServer(n3iwf n3iwf) (*Server, error) {
	s := &Server{
		n3iwf: n3iwf,
	}
	s.rcvPktCh = safe_channel.NewSafeCh[ReceiveNGAPPacket](RECEIVE_NGAPPACKET_CHANNEL_LEN)
	s.rcvEvtCh = safe_channel.NewSafeCh[n3iwf_context.NgapEvt](RECEIVE_NGAPEVENT_CHANNEL_LEN)
	return s, nil
}

// Run start the N3IWF SCTP process.
func (s *Server) Run(wg *sync.WaitGroup) error {
	// n3iwf context
	cfg := s.Config()

	localAddr := cfg.GetLocalSctpAddr()

	for _, remoteAddr := range cfg.GetAmfSctpAddrs() {
		errChan := make(chan error)
		wg.Add(1)
		go s.listenAndServe(localAddr, remoteAddr, errChan, wg)
		if err, ok := <-errChan; ok {
			return err
		}
	}

	wg.Add(1)
	go s.runNgapEventHandler(wg)

	return nil
}

func (s *Server) runNgapEventHandler(wg *sync.WaitGroup) {
	ngapLog := logger.NgapLog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			ngapLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		ngapLog.Infof("NGAP server stopped")
		s.rcvEvtCh.Close()
		s.rcvPktCh.Close()
		wg.Done()
	}()

	rcvEvtCh := s.rcvEvtCh.GetRcvChan()
	rcvPktCh := s.rcvPktCh.GetRcvChan()

	for {
		select {
		case rcvPkt := <-rcvPktCh:
			if len(rcvPkt.Buf) == 0 { // receiver closed
				return
			}
			s.NGAPDispatch(rcvPkt.Conn, rcvPkt.Buf)
		case rcvEvt := <-rcvEvtCh:
			s.HandleEvent(rcvEvt)
		}
	}
}

func (s *Server) listenAndServe(
	localAddr, remoteAddr *sctp.SCTPAddr,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {
	ngapLog := logger.NgapLog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			ngapLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		ngapLog.Infof("NGAP receiver stopped")
		wg.Done()
	}()

	var conn *sctp.SCTPConn
	var err error
	// Connect the session
	for i := 0; i < 3; i++ {
		conn, err = sctp.DialSCTP("sctp", localAddr, remoteAddr)
		if err != nil {
			ngapLog.Errorf("[SCTP] DialSCTP(): %+v", err)
		} else {
			break
		}

		if i != 2 {
			ngapLog.Info("Retry to connect AMF after 1 second...")
			time.Sleep(1 * time.Second)
		} else {
			ngapLog.Debugf("[SCTP] AMF SCTP address: %s", remoteAddr)
			errChan <- fmt.Errorf("failed to connect to AMF")
			return
		}
	}

	// Set default sender SCTP information sinfo_ppid = NGAP_PPID = 60
	info, err := conn.GetDefaultSentParam()
	if err != nil {
		ngapLog.Errorf("[SCTP] GetDefaultSentParam(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			ngapLog.Errorf("conn close error in GetDefaultSentParam(): %+v", errConn)
		}
		errChan <- fmt.Errorf("get socket information failed")
		return
	}
	info.PPID = lib_ngap.PPID
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		ngapLog.Errorf("[SCTP] SetDefaultSentParam(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			ngapLog.Errorf("conn close error in SetDefaultSentParam(): %+v", errConn)
		}
		errChan <- errors.New("set socket parameter failed")
		return
	}

	// Subscribe receiver SCTP information
	err = conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		ngapLog.Errorf("[SCTP] SubscribeEvents(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			ngapLog.Errorf("conn close error in SubscribeEvents(): %+v", errConn)
		}
		errChan <- errors.New("subscribe SCTP event failed")
		return
	}

	// Send NG setup request
	n3iwfCtx := s.Context()
	message.SendNGSetupRequest(conn, n3iwfCtx)

	close(errChan)

	s.conn = append(s.conn, conn)

	buf := make([]byte, factory.MAX_BUF_MSG_LEN)
	for {
		n, info, _, err := conn.SCTPRead(buf)
		if err != nil {
			ngapLog.Debugf("[SCTP] AMF SCTP address: %s", remoteAddr)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				ngapLog.Warn("[SCTP] Close connection.")
				errConn := conn.Close()
				if errConn != nil {
					ngapLog.Errorf("conn close error: %+v", errConn)
				}
				s.rcvPktCh.Send(ReceiveNGAPPacket{})
				return
			}
			ngapLog.Errorf("[SCTP] Read from SCTP connection failed: %+v", err)
			return
		}

		ngapLog.Tracef("[SCTP] Successfully read %d bytes.", n)

		if info == nil || info.PPID != lib_ngap.PPID {
			ngapLog.Warn("Received SCTP PPID != 60")
			continue
		}

		forwardData := make([]byte, n)
		copy(forwardData, buf[:n])

		ngapPkt := ReceiveNGAPPacket{
			Conn: conn,
			Buf:  forwardData[:n],
		}
		s.rcvPktCh.Send(ngapPkt)
	}
}

func (s *Server) SendNgapEvt(evt n3iwf_context.NgapEvt) {
	s.rcvEvtCh.Send(evt)
}

func (s *Server) Stop() {
	ngapLog := logger.NgapLog
	ngapLog.Infof("Close NGAP server....")

	for _, ngapServerConn := range s.conn {
		if err := ngapServerConn.Close(); err != nil {
			ngapLog.Errorf("Stop ngap server error : %+v", err)
		}
	}
}
