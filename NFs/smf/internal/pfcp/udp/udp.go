package udp

import (
	"errors"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/free5gc/pfcp"
	"github.com/free5gc/pfcp/pfcpUdp"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/internal/logger"
)

const MaxPfcpUdpDataSize = 1024

var Server *pfcpUdp.PfcpServer

var ServerStartTime time.Time

func Run(dispatch func(*pfcpUdp.Message)) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.PfcpLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	smfContext := smf_context.GetSelf()

	serverIP := smfContext.ListenIP().To4()
	Server = pfcpUdp.NewPfcpServer(serverIP.String())

	err := Server.Listen()
	if err != nil {
		logger.PfcpLog.Errorf("Failed to listen: %v", err)
	}

	logger.PfcpLog.Infof("Listen on %s", Server.Conn.LocalAddr().String())

	go func(p *pfcpUdp.PfcpServer) {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.PfcpLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		for {
			msg, errReadFrom := p.ReadFrom()
			if errReadFrom != nil {
				if errReadFrom == pfcpUdp.ErrReceivedResentRequest {
					logger.PfcpLog.Infoln(errReadFrom)
				} else if strings.Contains(errReadFrom.Error(), "use of closed network connection") {
					continue
				} else {
					logger.PfcpLog.Warnf("Read PFCP error: %v, msg: [%v]", errReadFrom, msg)
					select {
					case <-smfContext.PfcpContext.Done():
						// PFCP is closing
						return
					default:
						continue
					}
				}
				continue
			}

			if msg.PfcpMessage.IsRequest() {
				go dispatch(msg)
			}
		}
	}(Server)

	ServerStartTime = time.Now()

	logger.PfcpLog.Infof("Pfcp running... [%v]", ServerStartTime)
}

func SendPfcpResponse(sndMsg *pfcp.Message, addr *net.UDPAddr) {
	Server.WriteResponseTo(sndMsg, addr)
}

func SendPfcpRequest(sndMsg *pfcp.Message, addr *net.UDPAddr) (rsvMsg *pfcpUdp.Message, err error) {
	if addr.IP.Equal(net.IPv4zero) {
		return nil, errors.New("no destination IP address is specified")
	}
	return Server.WriteRequestTo(sndMsg, addr)
}

func ClosePfcp() error {
	smf_context.GetSelf().PfcpCancelFunc()

	closeErr := Server.Close()
	if closeErr != nil {
		logger.PfcpLog.Errorf("Pfcp close err: %+v", closeErr)
	} else {
		logger.PfcpLog.Infof("Pfcp server closed")
	}
	return closeErr
}
