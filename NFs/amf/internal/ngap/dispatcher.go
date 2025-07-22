package ngap

import (
	"net"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/ngap"
	"github.com/free5gc/sctp"
)

func Dispatch(conn net.Conn, msg []byte) {
	var ran *context.AmfRan
	amfSelf := context.GetSelf()

	ran, ok := amfSelf.AmfRanFindByConn(conn)
	if !ok {
		addr := conn.RemoteAddr()
		if addr == nil {
			logger.NgapLog.Warn("Addr of new NG connection is nii")
			return
		}
		logger.NgapLog.Infof("Create a new NG connection for: %s", addr.String())
		ran = amfSelf.NewAmfRan(conn)
	}

	if len(msg) == 0 {
		ran.Log.Infof("RAN close the connection.")
		ran.Remove()
		return
	}

	pdu, err := ngap.Decoder(msg)
	if err != nil {
		ran.Log.Errorf("NGAP decode error : %+v", err)
		return
	}

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}

	if pdu == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	dispatchMain(ran, pdu)
}

func HandleSCTPNotification(conn net.Conn, notification sctp.Notification) {
	amfSelf := context.GetSelf()

	logger.NgapLog.Infof("Handle SCTP Notification[addr: %+v]", conn.RemoteAddr())

	ran, ok := amfSelf.AmfRanFindByConn(conn)
	if !ok {
		logger.NgapLog.Warnf("RAN context has been removed[addr: %+v]", conn.RemoteAddr())
		return
	}

	switch notification.Type() {
	case sctp.SCTP_ASSOC_CHANGE:
		ran.Log.Infof("SCTP_ASSOC_CHANGE notification")
		event := notification.(*sctp.SCTPAssocChangeEvent)
		switch event.State() {
		case sctp.SCTP_COMM_LOST:
			ran.Log.Infof("SCTP state is SCTP_COMM_LOST, close the connection")
			ran.Remove()
		case sctp.SCTP_SHUTDOWN_COMP:
			ran.Log.Infof("SCTP state is SCTP_SHUTDOWN_COMP, close the connection")
			ran.Remove()
		default:
			ran.Log.Warnf("SCTP state[%+v] is not handled", event.State())
		}
	case sctp.SCTP_SHUTDOWN_EVENT:
		ran.Log.Infof("SCTP_SHUTDOWN_EVENT notification, close the connection")
		ran.Remove()
	default:
		ran.Log.Warnf("Non handled notification type: 0x%x", notification.Type())
	}
}

func HandleSCTPConnError(conn net.Conn) {
	amfSelf := context.GetSelf()

	logger.NgapLog.Infof("Handle SCTP Connection Error[addr: %+v] - remove RAN", conn.RemoteAddr())

	ran, ok := amfSelf.AmfRanFindByConn(conn)
	if !ok {
		logger.NgapLog.Warnf("RAN context has been removed[addr: %+v]", conn.RemoteAddr())
		return
	}
	ran.Remove()
}
