package ike

import (
	"net"
	"runtime/debug"

	ike_message "github.com/free5gc/ike/message"
	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/logger"
)

func (s *Server) Dispatch(
	udpConn *net.UDPConn,
	localAddr, remoteAddr *net.UDPAddr,
	ikeMessage *ike_message.IKEMessage, msg []byte,
	ikeSA *n3iwf_context.IKESecurityAssociation,
) {
	ikeLog := logger.IKELog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			ikeLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	switch ikeMessage.ExchangeType {
	case ike_message.IKE_SA_INIT:
		s.HandleIKESAINIT(udpConn, localAddr, remoteAddr, ikeMessage, msg)
	case ike_message.IKE_AUTH:
		s.HandleIKEAUTH(udpConn, localAddr, remoteAddr, ikeMessage, ikeSA)
	case ike_message.CREATE_CHILD_SA:
		s.HandleCREATECHILDSA(udpConn, localAddr, remoteAddr, ikeMessage, ikeSA)
	case ike_message.INFORMATIONAL:
		s.HandleInformational(udpConn, localAddr, remoteAddr, ikeMessage, ikeSA)
	default:
		ikeLog.Warnf("Unimplemented IKE message type, exchange type: %d", ikeMessage.ExchangeType)
	}
}
