package n3iwf_sctp

import (
	"errors"
	"io"
	"net"
	"sync"

	"github.com/ishidawataru/sctp"
	"github.com/sirupsen/logrus"

	"free5gc/src/n3iwf/factory"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_handler/n3iwf_message"
)

type SCTPSession struct {
	Address   string
	Port      int
	SessionID string
	Conn      *sctp.SCTPConn
}

var ngapLog *logrus.Entry

var peerAMFs map[string]*SCTPSession

const NGAP_PPID_BigEndian = 0x3c000000

func init() {
	peerAMFs = make(map[string]*SCTPSession)
	ngapLog = logger.NgapLog
}

func (session *SCTPSession) Connect() (sessionID string, err error) {
	// Check Address defined.
	var ipAddr *net.IPAddr

	ipAddr, err = net.ResolveIPAddr("ip", session.Address)
	if err != nil {
		ngapLog.Errorf("[SCTP] ResolveIPAddr(): %s", err.Error())
		err = errors.New("Failed to connect given AMF.")
		return
	}

	ip := []net.IPAddr{*ipAddr}

	// TODO: Bind local address according to configuration
	localAddr := new(sctp.SCTPAddr)

	remoteAddr := &sctp.SCTPAddr{
		IPAddrs: ip,
		Port:    session.Port,
	}

	session.Conn, err = sctp.DialSCTP("sctp", localAddr, remoteAddr)
	if err != nil {
		ngapLog.Errorf("[SCTP] DialSCTP(): %s", err.Error())
		err = errors.New("Failed to connect given AMF.")
		return
	}

	// Set default sender SCTP infomation sinfo_ppid = NGAP_PPID = 60
	info, err := session.Conn.GetDefaultSentParam()
	if err != nil {
		ngapLog.Errorf("[SCTP] GetDefaultSentParam(): %s", err.Error())
		err = errors.New("Failed to get socket infomation of given AMF.")
		session.Conn.Close()
		return
	}
	info.PPID = NGAP_PPID_BigEndian
	err = session.Conn.SetDefaultSentParam(info)
	if err != nil {
		ngapLog.Errorf("[SCTP] SetDefaultSentParam(): %s", err.Error())
		err = errors.New("Failed to set socket infomation of given AMF.")
		session.Conn.Close()
		return
	}

	// Subscribe receiver SCTP information
	err = session.Conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		ngapLog.Errorf("[SCTP] SubscribeEvents(): %s", err.Error())
		err = errors.New("Failed to subscribe SCTP event of given AMF socket.")
		session.Conn.Close()
		return
	}

	session.SessionID = session.Conn.RemoteAddr().String()
	sessionID = session.SessionID

	// Send EventSCTPConnectMessage to trigger NGSetup procedure
	handlerMessage := n3iwf_message.HandlerMessage{
		Event:         n3iwf_message.EventSCTPConnectMessage,
		SCTPSessionID: sessionID,
	}
	n3iwf_message.SendMessage(handlerMessage)

	ngapLog.Info("[SCTP] Successfully send event to N3IWF event queue.")

	return
}

func (session *SCTPSession) ClientListen(wg *sync.WaitGroup) {
	// Create a go routine to keep reading the connection.
	readData := make([]byte, 8192)

	go func(wg *sync.WaitGroup) {
		for {
			n, info, err := session.Conn.SCTPRead(readData)

			if err != nil {
				ngapLog.Errorf("[SCTP] SCTPRead(): %s", err.Error())
				ngapLog.Error("[SCTP] Failed to read from SCTP connection.")
				ngapLog.Debugf("[SCTP] AMF Address: %s\n[SCTP] Port: %d\n[SCTP] Session ID: %s", session.Address, session.Port, session.SessionID)

				if err == io.EOF || err == io.ErrUnexpectedEOF {
					ngapLog.Warn("[SCTP] Close connection.")
					ReleaseSession(session.SessionID)
					wg.Done()
					return
				}

			} else {
				ngapLog.Infof("[SCTP] Successfully read %d bytes.", n)

				if info == nil || info.PPID != NGAP_PPID_BigEndian {
					ngapLog.Warn("Recv SCTP PPID != 60")
					continue
				}

				handlerMessage := n3iwf_message.HandlerMessage{
					Event:         n3iwf_message.EventNGAPMessage,
					SCTPSessionID: session.SessionID,
					Value:         readData,
				}
				n3iwf_message.SendMessage(handlerMessage)

				ngapLog.Info("[SCTP] Successfully send data to N3IWF event queue.")
			}
		}
	}(wg)
}

func (session *SCTPSession) Send(data []byte) (err error) {
	if len(data) == 0 {
		ngapLog.Warn("[SCTP] Sending data is empty. Skipped.")
		return
	}

	var wroteBytes int

	wroteBytes, err = session.Conn.Write(data)
	if err != nil {
		ngapLog.Errorf("[SCTP] Write(): %s", err.Error())
		err = errors.New("Failed to send to AMF.")
	} else {
		ngapLog.Infof("[SCTP] Successfully sent %d bytes.", wroteBytes)
	}

	return
}

func (session *SCTPSession) Close() (err error) {
	err = session.Conn.Close()
	if err != nil {
		ngapLog.Errorf("[SCTP] %s", err.Error())
		err = errors.New("Failed to close session.")
	}
	return
}

// InitiateSCTP initiate the N3IWF SCTP process.
func InitiateSCTP(wg *sync.WaitGroup) {

	amfAddr := factory.N3iwfConfig.Configuration.AMFAddress

	for _, iterator := range amfAddr {
		// Create the session
		sctpSession := &SCTPSession{
			Address: iterator.NetworkAddress,
			Port:    assignPort(iterator.Port),
		}

		// Connect the session
		sessionID, err := sctpSession.Connect()
		if err != nil {
			ngapLog.Errorf("[SCTP] %s", err.Error())
			ngapLog.Debugf("[SCTP] AMF address: %s\n[SCTP] Remote port: %d", sctpSession.Address, sctpSession.Port)
			continue
		}

		// Add the session to map
		_, ok := peerAMFs[sessionID]
		if ok {
			ngapLog.Warn("[SCTP] InitiateSCTP(): SCTP session exists. The existing session will be released.")
			ngapLog.Debugf("[SCTP] Session ID: %s", sessionID)

			if ok := ReleaseSession(sessionID); !ok {
				// Improvement: retry mechanism
				continue
			} else {
				peerAMFs[sessionID] = sctpSession
			}

		} else {
			peerAMFs[sessionID] = sctpSession
		}

		// Add wait group number
		wg.Add(1)

		// Listen the session
		sctpSession.ClientListen(wg)
	}
}

func Send(sessionID string, data []byte) (ok bool) {
	if value, ok := peerAMFs[sessionID]; ok {
		if err := value.Send(data); err != nil {
			ngapLog.Errorf("[SCTP] %s", err.Error())
			return false
		}
		return true
	} else {
		ngapLog.Error("[SCTP] Send(): SCTP session not found.")
		ngapLog.Debugf("[SCTP] Session ID: %s", sessionID)
		return false
	}
}

func ReleaseSession(sessionID string) (ok bool) {
	if value, ok := peerAMFs[sessionID]; ok {
		if err := value.Close(); err != nil {
			ngapLog.Errorf("[SCTP] %s", err.Error())
			ngapLog.Debugf("[SCTP] Session ID: %s", sessionID)
			return false
		}

		delete(peerAMFs, sessionID)

		return true
	} else {
		ngapLog.Error("[SCTP] ReleaseSession(): SCTP session not found.")
		ngapLog.Debugf("[SCTP] Session ID: %s", sessionID)
		return false
	}
}

func assignPort(port int) int {
	if port == 0 {
		return 38412
	}
	return port
}
