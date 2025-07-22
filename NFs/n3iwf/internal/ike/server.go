package ike

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pkg/errors"

	"github.com/free5gc/ike"
	ike_message "github.com/free5gc/ike/message"
	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/n3iwf/pkg/factory"
	"github.com/free5gc/util/safe_channel"
)

const (
	RECEIVE_IKEPACKET_CHANNEL_LEN = 512
	RECEIVE_IKEEVENT_CHANNEL_LEN  = 512

	DEFAULT_IKE_PORT  = 500
	DEFAULT_NATT_PORT = 4500
)

type n3iwf interface {
	Config() *factory.Config
	Context() *n3iwf_context.N3IWFContext
	CancelContext() context.Context

	SendNgapEvt(n3iwf_context.NgapEvt)
}

type EspHandler func(srcIP, dstIP *net.UDPAddr, espPkt []byte) error

type Server struct {
	n3iwf

	Listener   map[int]*net.UDPConn
	StopServer chan struct{}
	rcvPktCh   *safe_channel.SafeCh[IkeReceivePacket]
	rcvEvtCh   *safe_channel.SafeCh[n3iwf_context.IkeEvt]
}

type IkeReceivePacket struct {
	Listener   net.UDPConn
	LocalAddr  net.UDPAddr
	RemoteAddr net.UDPAddr
	Msg        []byte
}

func NewServer(n3iwf n3iwf) (*Server, error) {
	s := &Server{
		n3iwf:      n3iwf,
		Listener:   make(map[int]*net.UDPConn),
		StopServer: make(chan struct{}),
	}
	s.rcvPktCh = safe_channel.NewSafeCh[IkeReceivePacket](RECEIVE_IKEPACKET_CHANNEL_LEN)
	s.rcvEvtCh = safe_channel.NewSafeCh[n3iwf_context.IkeEvt](RECEIVE_IKEEVENT_CHANNEL_LEN)
	return s, nil
}

func (s *Server) Run(wg *sync.WaitGroup) error {
	cfg := s.Config()

	// Resolve UDP addresses
	ip := cfg.GetIKEBindAddr()
	ikeAddrPort, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, DEFAULT_IKE_PORT))
	if err != nil {
		return err
	}
	nattAddrPort, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, DEFAULT_NATT_PORT))
	if err != nil {
		return err
	}

	// Listen and serve
	var errChan chan error

	wg.Add(1)
	errChan = make(chan error)
	go s.receiver(ikeAddrPort, errChan, wg)
	if err, ok := <-errChan; ok {
		return errors.Wrapf(err, "ikeAddrPort")
	}

	wg.Add(1)
	errChan = make(chan error)
	go s.receiver(nattAddrPort, errChan, wg)
	if err, ok := <-errChan; ok {
		return errors.Wrapf(err, "nattAddrPort")
	}

	wg.Add(1)
	go s.server(wg)

	return nil
}

func (s *Server) server(wg *sync.WaitGroup) {
	ikeLog := logger.IKELog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			ikeLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		ikeLog.Infof("Ike server stopped")
		s.rcvPktCh.Close()
		s.rcvEvtCh.Close()
		close(s.StopServer)
		wg.Done()
	}()

	rcvEvtCh := s.rcvEvtCh.GetRcvChan()
	rcvPktCh := s.rcvPktCh.GetRcvChan()

	for {
		select {
		case rcvPkt := <-rcvPktCh:
			ikeMsg, ikeSA, err := s.checkIKEMessage(
				rcvPkt.Msg, &rcvPkt.Listener, &rcvPkt.LocalAddr, &rcvPkt.RemoteAddr)
			if err != nil {
				ikeLog.Warnln(err)
				continue
			}
			s.Dispatch(&rcvPkt.Listener, &rcvPkt.LocalAddr, &rcvPkt.RemoteAddr,
				ikeMsg, rcvPkt.Msg, ikeSA)
		case rcvIkeEvent := <-rcvEvtCh:
			s.HandleEvent(rcvIkeEvent)
		case <-s.StopServer:
			return
		}
	}
}

func (s *Server) receiver(
	localAddr *net.UDPAddr,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {
	ikeLog := logger.IKELog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			ikeLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		ikeLog.Infof("Ike receiver stopped")
		wg.Done()
	}()

	listener, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		ikeLog.Errorf("Listen UDP failed: %+v", err)
		errChan <- errors.New("listenAndServe failed")
		return
	}

	close(errChan)

	s.Listener[localAddr.Port] = listener

	buf := make([]byte, factory.MAX_BUF_MSG_LEN)

	for {
		n, remoteAddr, err := listener.ReadFromUDP(buf)
		if err != nil {
			ikeLog.Errorf("ReadFromUDP failed: %+v", err)
			return
		}

		msgBuf := make([]byte, n)
		copy(msgBuf, buf)
		ikeLog.Tracef("recv from port(%d):\n%s", localAddr.Port, hex.Dump(msgBuf))

		// As specified in RFC 7296 section 3.1, the IKE message send from/to UDP port 4500
		// should prepend a 4 bytes zero
		if localAddr.Port == DEFAULT_NATT_PORT {
			msgBuf, err = handleNattMsg(msgBuf, remoteAddr, localAddr, handleESPPacket)
			if err != nil {
				ikeLog.Errorf("Handle NATT msg: %v", err)
				continue
			}
			if msgBuf == nil {
				continue
			}
		}

		if len(msgBuf) < ike_message.IKE_HEADER_LEN {
			ikeLog.Warnf("Received IKE msg is too short from %s", remoteAddr)
			continue
		}

		ikePkt := IkeReceivePacket{
			RemoteAddr: *remoteAddr,
			Listener:   *listener,
			LocalAddr:  *localAddr,
			Msg:        msgBuf,
		}
		s.rcvPktCh.Send(ikePkt)
	}
}

func handleNattMsg(
	msgBuf []byte,
	rAddr, lAddr *net.UDPAddr,
	espHandler EspHandler,
) ([]byte, error) {
	if len(msgBuf) == 1 && msgBuf[0] == 0xff {
		// skip NAT-T Keepalive
		return nil, nil
	}

	nonEspMarker := []byte{0, 0, 0, 0} // Non-ESP Marker
	nonEspMarkerLen := len(nonEspMarker)
	if len(msgBuf) < nonEspMarkerLen {
		return nil, errors.Errorf("Received msg is too short")
	}
	if !bytes.Equal(msgBuf[:nonEspMarkerLen], nonEspMarker) {
		// ESP packet
		if espHandler != nil {
			err := espHandler(rAddr, lAddr, msgBuf)
			if err != nil {
				return nil, errors.Wrapf(err, "Handle ESP")
			}
		}
		return nil, nil
	}

	// IKE message: skip Non-ESP Marker
	msgBuf = msgBuf[nonEspMarkerLen:]
	return msgBuf, nil
}

func (s *Server) SendIkeEvt(evt n3iwf_context.IkeEvt) {
	s.rcvEvtCh.Send(evt)
}

func (s *Server) Stop() {
	ikeLog := logger.IKELog
	ikeLog.Infof("Close Ike server...")

	for _, ikeServerListener := range s.Listener {
		if err := ikeServerListener.Close(); err != nil {
			ikeLog.Errorf("Stop ike server : %s error : %+v", err, ikeServerListener.LocalAddr().String())
		}
	}

	s.StopServer <- struct{}{}
}

func (s *Server) checkIKEMessage(
	msg []byte, udpConn *net.UDPConn,
	localAddr, remoteAddr *net.UDPAddr,
) (*ike_message.IKEMessage,
	*n3iwf_context.IKESecurityAssociation, error,
) {
	var ikeHeader *ike_message.IKEHeader
	var ikeMessage *ike_message.IKEMessage
	var ikeSA *n3iwf_context.IKESecurityAssociation
	var err error

	// parse IKE header and setup IKE context
	ikeHeader, err = ike_message.ParseHeader(msg)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "IKE msg decode header")
	}

	// check major version
	if ikeHeader.MajorVersion > 2 {
		// send INFORMATIONAL type message with INVALID_MAJOR_VERSION Notify payload
		// For response or needed data
		payload := new(ike_message.IKEPayloadContainer)
		payload.BuildNotification(ike_message.TypeNone,
			ike_message.INVALID_MAJOR_VERSION, nil, nil)
		responseIKEMessage := ike_message.NewMessage(ikeHeader.InitiatorSPI, ikeHeader.ResponderSPI,
			ike_message.INFORMATIONAL, true, false, ikeHeader.MessageID, *payload)

		err = SendIKEMessageToUE(udpConn, localAddr, remoteAddr, responseIKEMessage, nil)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Received an IKE message with higher major version "+
				"(%d>2)", ikeHeader.MajorVersion)
		}
		return nil, nil, errors.Errorf("Received an IKE message with higher major version (%d>2)", ikeHeader.MajorVersion)
	}

	if ikeHeader.ExchangeType == ike_message.IKE_SA_INIT {
		ikeMessage, err = ike.DecodeDecrypt(msg, ikeHeader,
			nil, ike_message.Role_Responder)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Decrypt IkeMsg error")
		}
	} else if ikeHeader.ExchangeType != ike_message.IKE_SA_INIT {
		localSPI := ikeHeader.ResponderSPI
		var ok bool
		n3iwfCtx := s.Context()

		ikeSA, ok = n3iwfCtx.IKESALoad(localSPI)
		if !ok {
			payload := new(ike_message.IKEPayloadContainer)
			// send INFORMATIONAL type message with INVALID_IKE_SPI Notify payload ( OUTSIDE IKE SA )
			payload.BuildNotification(ike_message.TypeNone, ike_message.INVALID_IKE_SPI, nil, nil)
			responseIKEMessage := ike_message.NewMessage(ikeHeader.InitiatorSPI, ikeHeader.ResponderSPI,
				ike_message.INFORMATIONAL, true, false, ikeHeader.MessageID, *payload)

			err = SendIKEMessageToUE(udpConn, localAddr, remoteAddr, responseIKEMessage, nil)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "checkIKEMessage():")
			}
			return nil, nil, errors.Errorf("Received an unrecognized SPI message: %d", localSPI)
		}

		ikeMessage, err = ike.DecodeDecrypt(msg, ikeHeader,
			ikeSA.IKESAKey, ike_message.Role_Responder)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Decrypt IkeMsg error")
		}
	}

	return ikeMessage, ikeSA, nil
}

func constructPacketWithESP(srcIP, dstIP *net.UDPAddr, espPacket []byte) ([]byte, error) {
	ipLayer := &layers.IPv4{
		SrcIP:    srcIP.IP,
		DstIP:    dstIP.IP,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolESP,
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	err := gopacket.SerializeLayers(buffer,
		options,
		ipLayer,
		gopacket.Payload(espPacket),
	)
	if err != nil {
		return nil, errors.Errorf("Error serializing layers: %v", err)
	}

	packetData := buffer.Bytes()
	return packetData, nil
}

func handleESPPacket(srcIP, dstIP *net.UDPAddr, espPacket []byte) error {
	ikeLog := logger.IKELog
	ikeLog.Tracef("Handle ESPPacket")

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return errors.Errorf("socket error: %v", err)
	}

	defer func() {
		if err = syscall.Close(fd); err != nil {
			ikeLog.Errorf("Close fd error : %v", err)
		}
	}()

	ipPacket, err := constructPacketWithESP(srcIP, dstIP, espPacket)
	if err != nil {
		return err
	}

	addr := syscall.SockaddrInet4{
		Addr: [4]byte(dstIP.IP),
		Port: dstIP.Port,
	}

	err = syscall.Sendto(fd, ipPacket, 0, &addr)
	if err != nil {
		return errors.Errorf("sendto error: %v", err)
	}

	return nil
}
