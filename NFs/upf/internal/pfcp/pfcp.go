package pfcp

import (
	"encoding/hex"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wmnsk/go-pfcp/message"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/internal/report"
	"github.com/free5gc/go-upf/pkg/factory"
	logger_util "github.com/free5gc/util/logger"
)

const (
	RECEIVE_CHANNEL_LEN       = 512
	REPORT_CHANNEL_LEN        = 128
	TRANS_TIMEOUT_CHANNEL_LEN = 64
	MAX_PFCP_MSG_LEN          = 65536
)

type ReceivePacket struct {
	RemoteAddr net.Addr
	Buf        []byte
}

type TransType int

const (
	TX TransType = iota
	RX
)

type TransactionTimeout struct {
	TrType TransType
	TrID   string
}

type PfcpServer struct {
	cfg          *factory.Config
	listen       string
	nodeID       string
	rcvCh        chan ReceivePacket
	srCh         chan report.SessReport
	trToCh       chan TransactionTimeout
	conn         *net.UDPConn
	recoveryTime time.Time
	driver       forwarder.Driver
	lnode        LocalNode
	rnodes       map[string]*RemoteNode
	txTrans      map[string]*TxTransaction // key: RemoteAddr-Sequence
	rxTrans      map[string]*RxTransaction // key: RemoteAddr-Sequence
	txSeq        uint32
	log          *logrus.Entry
}

func NewPfcpServer(cfg *factory.Config, driver forwarder.Driver) *PfcpServer {
	listen := fmt.Sprintf("%s:%d", cfg.Pfcp.Addr, factory.UpfPfcpDefaultPort)
	return &PfcpServer{
		cfg:          cfg,
		listen:       listen,
		nodeID:       cfg.Pfcp.NodeID,
		rcvCh:        make(chan ReceivePacket, RECEIVE_CHANNEL_LEN),
		srCh:         make(chan report.SessReport, REPORT_CHANNEL_LEN),
		trToCh:       make(chan TransactionTimeout, TRANS_TIMEOUT_CHANNEL_LEN),
		recoveryTime: time.Now(),
		driver:       driver,
		rnodes:       make(map[string]*RemoteNode),
		txTrans:      make(map[string]*TxTransaction),
		rxTrans:      make(map[string]*RxTransaction),
		log:          logger.PfcpLog.WithField(logger_util.FieldListenAddr, listen),
	}
}

func (s *PfcpServer) main(wg *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			s.log.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}

		s.log.Infoln("pfcp server stopped")
		s.stopTrTimers()
		close(s.rcvCh)
		close(s.srCh)
		close(s.trToCh)
		wg.Done()
	}()

	var err error
	laddr, err := net.ResolveUDPAddr("udp4", s.listen)
	if err != nil {
		s.log.Errorf("Resolve err: %+v", err)
		return
	}

	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		s.log.Errorf("Listen err: %+v", err)
		return
	}
	s.conn = conn

	wg.Add(1)
	go s.receiver(wg)

	for {
		select {
		case sr := <-s.srCh:
			s.log.Tracef("receive SessReport from srCh")
			s.ServeReport(&sr)
		case rcvPkt := <-s.rcvCh:
			s.log.Tracef("receive buf(len=%d) from rcvCh", len(rcvPkt.Buf))
			if len(rcvPkt.Buf) == 0 {
				// receiver closed
				return
			}
			msg, err := message.Parse(rcvPkt.Buf)
			if err != nil {
				s.log.Errorln(err)
				s.log.Tracef("ignored undecodable message:\n%+v", hex.Dump(rcvPkt.Buf))
				continue
			}

			trID := fmt.Sprintf("%s-%d", rcvPkt.RemoteAddr, msg.Sequence())
			if isRequest(msg) {
				s.log.Tracef("receive req pkt from %s", trID)
				rx, ok := s.rxTrans[trID]
				if !ok {
					rx = NewRxTransaction(s, rcvPkt.RemoteAddr, msg.Sequence())
					s.rxTrans[trID] = rx
				}
				needDispatch, err1 := rx.recv(msg, ok)
				if err1 != nil {
					s.log.Warnf("rcvCh: %v", err1)
					continue
				} else if !needDispatch {
					s.log.Debugf("rcvCh: rxtr[%s] req no need to dispatch", trID)
					continue
				}
				err = s.reqDispacher(msg, rcvPkt.RemoteAddr)
				if err != nil {
					s.log.Errorln(err)
					s.log.Tracef("ignored undecodable message:\n%+v", hex.Dump(rcvPkt.Buf))
				}
			} else if isResponse(msg) {
				s.log.Tracef("receive rsp pkt from %s", trID)
				tx, ok := s.txTrans[trID]
				if !ok {
					s.log.Debugf("rcvCh: No txtr[%s] found for rsp", trID)
					continue
				}
				req := tx.recv(msg)
				err = s.rspDispacher(msg, rcvPkt.RemoteAddr, req)
				if err != nil {
					s.log.Errorln(err)
					s.log.Tracef("ignored undecodable message:\n%+v", hex.Dump(rcvPkt.Buf))
				}
			}
		case trTo := <-s.trToCh:
			s.log.Tracef("receive tr timeout (%v) from trToCh", trTo)
			if trTo.TrType == TX {
				tx, ok := s.txTrans[trTo.TrID]
				if !ok {
					s.log.Warnf("trToCh: txtr[%s] not found", trTo.TrID)
					continue
				}
				tx.handleTimeout()
			} else { // RX
				rx, ok := s.rxTrans[trTo.TrID]
				if !ok {
					s.log.Warnf("trToCh: rxtr[%s] not found", trTo.TrID)
					continue
				}
				rx.handleTimeout()
			}
		}
	}
}

func (s *PfcpServer) receiver(wg *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			s.log.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}

		s.log.Infoln("pfcp reciver stopped")
		wg.Done()
	}()

	buf := make([]byte, MAX_PFCP_MSG_LEN)
	for {
		s.log.Tracef("receiver starts to read...")
		n, addr, err := s.conn.ReadFrom(buf)
		if err != nil {
			s.log.Errorf("%+v", err)
			s.rcvCh <- ReceivePacket{}
			break
		}

		s.log.Tracef("receiver reads message(len=%d)", n)
		msgBuf := make([]byte, n)
		copy(msgBuf, buf)
		s.rcvCh <- ReceivePacket{
			RemoteAddr: addr,
			Buf:        msgBuf,
		}
	}
}

func (s *PfcpServer) Start(wg *sync.WaitGroup) {
	s.log.Infoln("starting pfcp server")
	wg.Add(1)
	go s.main(wg)
	s.log.Infoln("pfcp server started")
}

func (s *PfcpServer) Stop() {
	s.log.Infoln("Stopping pfcp server")
	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			s.log.Errorf("Stop pfcp server err: %+v", err)
		}
	}
}

func (s *PfcpServer) NewNode(id string, addr net.Addr, driver forwarder.Driver) *RemoteNode {
	n := NewRemoteNode(
		id,
		addr,
		&s.lnode,
		driver,
		s.log.WithField(logger_util.FieldControlPlaneNodeID, id),
	)
	n.log.Infoln("New node")
	return n
}

func (s *PfcpServer) UpdateNodeID(n *RemoteNode, newId string) {
	s.log.Infof("Update nodeId %q to %q", n.ID, newId)
	delete(s.rnodes, n.ID)
	n.ID = newId
	n.log = s.log.WithField(logger_util.FieldControlPlaneNodeID, newId)
	s.rnodes[newId] = n
}

func (s *PfcpServer) NotifySessReport(sr report.SessReport) {
	s.srCh <- sr
}

func (s *PfcpServer) NotifyTransTimeout(trType TransType, trID string) {
	s.trToCh <- TransactionTimeout{TrType: trType, TrID: trID}
}

func (s *PfcpServer) PopBufPkt(seid uint64, pdrid uint16) ([]byte, bool) {
	sess, err := s.lnode.Sess(seid)
	if err != nil {
		s.log.Errorln(err)
		return nil, false
	}
	return sess.Pop(pdrid)
}

func (s *PfcpServer) sendReqTo(msg message.Message, addr net.Addr) error {
	if !isRequest(msg) {
		return errors.Errorf("sendReqTo: invalid req type(%d)", msg.MessageType())
	}

	txtr := NewTxTransaction(s, addr, s.txSeq)
	s.txSeq++
	s.txTrans[txtr.id] = txtr

	return txtr.send(msg)
}

func (s *PfcpServer) sendRspTo(msg message.Message, addr net.Addr) error {
	if !isResponse(msg) {
		return errors.Errorf("sendRspTo: invalid rsp type(%d)", msg.MessageType())
	}

	// find transaction
	trID := fmt.Sprintf("%s-%d", addr, msg.Sequence())
	rxtr, ok := s.rxTrans[trID]
	if !ok {
		return errors.Errorf("sendRspTo: rxtr(%s) not found", trID)
	}

	return rxtr.send(msg)
}

func (s *PfcpServer) stopTrTimers() {
	for _, tx := range s.txTrans {
		if tx.timer == nil {
			continue
		}
		tx.timer.Stop()
		tx.timer = nil
	}
	for _, rx := range s.rxTrans {
		if rx.timer == nil {
			continue
		}
		rx.timer.Stop()
		rx.timer = nil
	}
}

func isRequest(msg message.Message) bool {
	switch msg.MessageType() {
	case message.MsgTypeHeartbeatRequest:
		return true
	case message.MsgTypePFDManagementRequest:
		return true
	case message.MsgTypeAssociationSetupRequest:
		return true
	case message.MsgTypeAssociationUpdateRequest:
		return true
	case message.MsgTypeAssociationReleaseRequest:
		return true
	case message.MsgTypeNodeReportRequest:
		return true
	case message.MsgTypeSessionSetDeletionRequest:
		return true
	case message.MsgTypeSessionEstablishmentRequest:
		return true
	case message.MsgTypeSessionModificationRequest:
		return true
	case message.MsgTypeSessionDeletionRequest:
		return true
	case message.MsgTypeSessionReportRequest:
		return true
	default:
	}
	return false
}

func isResponse(msg message.Message) bool {
	switch msg.MessageType() {
	case message.MsgTypeHeartbeatResponse:
		return true
	case message.MsgTypePFDManagementResponse:
		return true
	case message.MsgTypeAssociationSetupResponse:
		return true
	case message.MsgTypeAssociationUpdateResponse:
		return true
	case message.MsgTypeAssociationReleaseResponse:
		return true
	case message.MsgTypeNodeReportResponse:
		return true
	case message.MsgTypeSessionSetDeletionResponse:
		return true
	case message.MsgTypeSessionEstablishmentResponse:
		return true
	case message.MsgTypeSessionModificationResponse:
		return true
	case message.MsgTypeSessionDeletionResponse:
		return true
	case message.MsgTypeSessionReportResponse:
		return true
	default:
	}
	return false
}

func setReqSeq(msgtmp message.Message, seq uint32) {
	switch msg := msgtmp.(type) {
	case *message.HeartbeatRequest:
		msg.SetSequenceNumber(seq)
	case *message.PFDManagementRequest:
		msg.SetSequenceNumber(seq)
	case *message.AssociationSetupRequest:
		msg.SetSequenceNumber(seq)
	case *message.AssociationUpdateRequest:
		msg.SetSequenceNumber(seq)
	case *message.AssociationReleaseRequest:
		msg.SetSequenceNumber(seq)
	case *message.NodeReportRequest:
		msg.SetSequenceNumber(seq)
	case *message.SessionSetDeletionRequest:
		msg.SetSequenceNumber(seq)
	case *message.SessionEstablishmentRequest:
		msg.SetSequenceNumber(seq)
	case *message.SessionModificationRequest:
		msg.SetSequenceNumber(seq)
	case *message.SessionDeletionRequest:
		msg.SetSequenceNumber(seq)
	case *message.SessionReportRequest:
		msg.SetSequenceNumber(seq)
	default:
	}
}
