package amf_ngap_sctp

import (
	"github.com/ishidawataru/sctp"
	"free5gc/src/amf/amf_handler/amf_message"
	"net"
	"sync"
	"time"

	"free5gc/lib/ngap/ngapSctp"
	"free5gc/src/amf/logger"
)

var readChan chan ngapSctp.ConnData

func init() {
	readChan = make(chan ngapSctp.ConnData, 1024)
}

type SCTPListener struct {
	ln   *sctp.SCTPListener
	mtx  sync.Mutex
	conn map[string]net.Conn
}

func Server(addrStr string) (listener *SCTPListener) {
	ln := ngapSctp.Server(addrStr)
	listener = &SCTPListener{
		ln:   ln,
		conn: make(map[string]net.Conn),
	}
	go listener.forwardData()

	// Wait for accept connection
	go func(l *SCTPListener) {
		for {
			conn, err := ngapSctp.Accept(ln)
			if err != nil {
				logger.NgapLog.Warn(err.Error())
				continue
			}
			logger.NgapLog.Infof("[AMF] NGAP SCTP Accept from: %s", conn.RemoteAddr().String())
			// send connection to amf handler
			msg := amf_message.HandlerMessage{}
			msg.Event = amf_message.EventNGAPAcceptConn
			msg.Value = conn
			amf_message.SendMessage(msg)
			l.mtx.Lock()
			l.conn[conn.RemoteAddr().String()] = conn
			l.mtx.Unlock()
			go ngapSctp.Start(conn, readChan)

			// put connection into global conn
			// SctpConnection = append(SctpConnection, conn)
		}

	}(listener)
	return
}
func (l *SCTPListener) Close() {
	logger.NgapLog.Infoln("Close listener")
	l.mtx.Lock()
	msg := amf_message.HandlerMessage{}
	for key, conn := range l.conn {
		msg.Event = amf_message.EventNGAPCloseConn
		msg.NgapAddr = conn.RemoteAddr().String()
		delete(l.conn, key)
		amf_message.SendMessage(msg)
	}
	logger.NgapLog.Errorln(l.ln.Addr())
	l.ln.Close()
	time.Sleep(10 * time.Millisecond)
	l.mtx.Unlock()

}

func (l *SCTPListener) forwardData() {
	// time.Sleep(3000 * time.Microsecond)
	defer close(readChan)
	for {
		// logger.NgapLog.Printf("Channel buffer size: %d", len(readChan))
		ngapChan := <-readChan
		msg := amf_message.HandlerMessage{}
		raddr := ngapChan.GetRAddr()
		if ngapChan.GetError() != nil {
			if _, ok := l.conn[raddr]; ok {
				msg.Event = amf_message.EventNGAPCloseConn
				msg.NgapAddr = raddr
				l.mtx.Lock()
				delete(l.conn, raddr)
				l.mtx.Unlock()
			} else {
				continue
			}

		} else {
			msg.Event = amf_message.EventNGAPMessage
			msg.NgapAddr = raddr
			msg.Value = ngapChan.GetData()
			logger.NgapLog.Debugf("Packet get: 0x%x", msg.Value)
		}
		amf_message.SendMessage(msg)
		// logger.NgapLog.Printf("Packet get: %s", packet)
	}
}
