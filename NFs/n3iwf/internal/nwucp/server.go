package nwucp

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"net"
	"runtime/debug"
	"strings"
	"sync"

	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/n3iwf/internal/ngap/message"
	"github.com/free5gc/n3iwf/pkg/factory"
)

type n3iwf interface {
	Config() *factory.Config
	Context() *n3iwf_context.N3IWFContext
	CancelContext() context.Context

	SendNgapEvt(n3iwf_context.NgapEvt)
}

type Server struct {
	n3iwf

	tcpListener net.Listener
}

func NewServer(n3iwf n3iwf) (*Server, error) {
	s := &Server{
		n3iwf: n3iwf,
	}
	return s, nil
}

// Run setup N3IWF NAS for UE to forward NAS message
// to AMF
func (s *Server) Run(wg *sync.WaitGroup) error {
	cfg := s.Config()
	listener, err := net.Listen("tcp", cfg.GetNasTcpAddr())
	if err != nil {
		return err
	}
	s.tcpListener = listener

	wg.Add(1)
	go s.listenAndServe(wg)

	return nil
}

// listenAndServe handle TCP listener and accept incoming
// requests. It also stores accepted connection into UE
// context, and finally, call serveConn() to serve the messages
// received from the connection.
func (s *Server) listenAndServe(wg *sync.WaitGroup) {
	nwucpLog := logger.NWuCPLog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			nwucpLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}

		err := s.tcpListener.Close()
		if err != nil {
			nwucpLog.Errorf("Error closing tcpListener: %+v", err)
		}
		wg.Done()
	}()

	n3iwfCtx := s.Context()
	for {
		connection, err := s.tcpListener.Accept()
		if err != nil {
			nwucpLog.Errorf("TCP server accept failed : %+v. Close the listener...", err)
			return
		}

		nwucpLog.Tracef("Accepted one UE from %+v", connection.RemoteAddr())

		// Find UE context and store this connection in to it, then check if
		// there is any cached NAS message for this UE. If yes, send to it.

		ueIP := strings.Split(connection.RemoteAddr().String(), ":")[0]
		ikeUe, ok := n3iwfCtx.AllocatedUEIPAddressLoad(ueIP)
		if !ok {
			nwucpLog.Errorf("UE context not found for peer %+v", ueIP)
			continue
		}

		ranUe, err := n3iwfCtx.RanUeLoadFromIkeSPI(ikeUe.N3IWFIKESecurityAssociation.LocalSPI)
		if err != nil {
			nwucpLog.Errorf("RanUe context not found : %v", err)
			continue
		}

		n3iwfUe, ok := ranUe.(*n3iwf_context.N3IWFRanUe)
		if !ok {
			nwucpLog.Errorf("listenAndServe(): [Type Assertion] RanUe -> N3iwfUe failed")
			continue
		}

		// Store connection
		n3iwfUe.TCPConnection = connection

		s.SendNgapEvt(n3iwf_context.NewNASTCPConnEstablishedCompleteEvt(n3iwfUe.RanUeNgapId))

		wg.Add(1)
		go serveConn(n3iwfUe, connection, wg)
	}
}

func (s *Server) Stop() {
	nwucpLog := logger.NWuCPLog
	nwucpLog.Infof("Close Nwucp server...")

	if err := s.tcpListener.Close(); err != nil {
		nwucpLog.Errorf("Stop nwucp server error : %+v", err)
	}

	// TODO: [Bug] TCPConnection may close twice, need to check
	s.Context().RANUePool.Range(
		func(key, value interface{}) bool {
			ranUe, ok := value.(*n3iwf_context.N3IWFRanUe)
			if ok && ranUe.TCPConnection != nil {
				if err := ranUe.TCPConnection.Close(); err != nil {
					logger.InitLog.Errorf("Stop nwucp server error : %+v", err)
				}
			}
			return true
		},
	)
}

// serveConn handle accepted TCP connection. It reads NAS packets
// from the connection and call forward() to forward NAS messages
// to AMF
func serveConn(ranUe *n3iwf_context.N3IWFRanUe, connection net.Conn, wg *sync.WaitGroup) {
	nwucpLog := logger.NWuCPLog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			nwucpLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}

		err := connection.Close()
		if err != nil {
			nwucpLog.Errorf("Error closing connection: %+v", err)
		}
		wg.Done()
	}()

	connReader := bufio.NewReader(connection)
	buf := make([]byte, factory.MAX_BUF_MSG_LEN)
	for {
		// Read the length of NAS message
		n, err := io.ReadFull(connReader, buf[:2])
		if err != nil {
			nwucpLog.Errorf("Read the length of NAS message failed: %+v", err)
			ranUe.TCPConnection = nil
			return
		}
		nasLen := binary.BigEndian.Uint16(buf[:n])
		if uint64(nasLen) > uint64(cap(buf)) {
			buf = make([]byte, 0, nasLen)
		}

		// Read the NAS message
		n, err = io.ReadFull(connReader, buf[:nasLen])
		if err != nil {
			nwucpLog.Errorf("Read the NAS message failed: %+v", err)
			ranUe.TCPConnection = nil
			return
		}
		fwdNas := make([]byte, n)
		copy(fwdNas, buf[:n])

		wg.Add(1)
		go forward(ranUe, fwdNas, wg)
	}
}

// forward forwards NAS messages sent from UE to the
// associated AMF
func forward(ranUe *n3iwf_context.N3IWFRanUe, packet []byte, wg *sync.WaitGroup) {
	nwucpLog := logger.NWuCPLog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			nwucpLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		wg.Done()
	}()

	nwucpLog.Trace("Forward NWu -> N2")
	message.SendUplinkNASTransport(ranUe, packet)
}
