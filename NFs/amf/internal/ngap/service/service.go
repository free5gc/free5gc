package service

import (
	"encoding/hex"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/ngap"
	"github.com/free5gc/sctp"
)

type NGAPHandler struct {
	HandleMessage         func(conn net.Conn, msg []byte)
	HandleNotification    func(conn net.Conn, notification sctp.Notification)
	HandleConnectionError func(conn net.Conn)
}

const (
	notimeout   int    = -1
	readBufSize uint32 = 262144
)

// set default read timeout to 2 seconds
var readTimeout syscall.Timeval = syscall.Timeval{Sec: 2, Usec: 0}

var (
	sctpListener *sctp.SCTPListener
	connections  sync.Map
)

func NewSctpConfig(cfg *factory.Sctp) *sctp.SocketConfig {
	sctpConfig := &sctp.SocketConfig{
		InitMsg: sctp.InitMsg{
			NumOstreams:    uint16(cfg.NumOstreams),
			MaxInstreams:   uint16(cfg.MaxInstreams),
			MaxAttempts:    uint16(cfg.MaxAttempts),
			MaxInitTimeout: uint16(cfg.MaxInitTimeout),
		},
		RtoInfo:   &sctp.RtoInfo{SrtoAssocID: 0, SrtoInitial: 500, SrtoMax: 1500, StroMin: 100},
		AssocInfo: &sctp.AssocInfo{AsocMaxRxt: 4},
	}
	return sctpConfig
}

func Run(addresses []string, port int, handler NGAPHandler, sctpConfig *sctp.SocketConfig) {
	ips := []net.IPAddr{}

	for _, addr := range addresses {
		if netAddr, err := net.ResolveIPAddr("ip", addr); err != nil {
			logger.NgapLog.Errorf("Error resolving address '%s': %v\n", addr, err)
		} else {
			logger.NgapLog.Debugf("Resolved address '%s' to %s\n", addr, netAddr)
			ips = append(ips, *netAddr)
		}
	}

	addr := &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    port,
	}

	go listenAndServe(addr, handler, sctpConfig)
}

func listenAndServe(addr *sctp.SCTPAddr, handler NGAPHandler, sctpConfig *sctp.SocketConfig) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.NgapLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if sctpConfig == nil {
		logger.NgapLog.Errorf("Error sctp SocketConfig is nil")
		return
	}

	if listener, err := sctpConfig.Listen("sctp", addr); err != nil {
		logger.NgapLog.Errorf("Failed to listen: %+v", err)
		return
	} else {
		sctpListener = listener
	}

	logger.NgapLog.Infof("Listen on %s", sctpListener.Addr())

	for {
		newConn, err := sctpListener.AcceptSCTP(notimeout)
		if err != nil {
			switch err {
			case syscall.EINTR, syscall.EAGAIN:
				logger.NgapLog.Debugf("AcceptSCTP: %+v", err)
			default:
				logger.NgapLog.Errorf("Failed to accept: %+v", err)
			}
			continue
		}

		var info *sctp.SndRcvInfo
		if infoTmp, errGetDefaultSentParam := newConn.GetDefaultSentParam(); errGetDefaultSentParam != nil {
			logger.NgapLog.Errorf("Get default sent param error: %+v, accept failed", errGetDefaultSentParam)
			if errGetDefaultSentParam = newConn.Close(); errGetDefaultSentParam != nil {
				logger.NgapLog.Errorf("Close error: %+v", errGetDefaultSentParam)
			}
			continue
		} else {
			info = infoTmp
			logger.NgapLog.Debugf("Get default sent param[value: %+v]", info)
		}

		info.PPID = ngap.PPID
		if errSetDefaultSentParam := newConn.SetDefaultSentParam(info); errSetDefaultSentParam != nil {
			logger.NgapLog.Errorf("Set default sent param error: %+v, accept failed", errSetDefaultSentParam)
			if errSetDefaultSentParam = newConn.Close(); errSetDefaultSentParam != nil {
				logger.NgapLog.Errorf("Close error: %+v", errSetDefaultSentParam)
			}
			continue
		} else {
			logger.NgapLog.Debugf("Set default sent param[value: %+v]", info)
		}

		events := sctp.SCTP_EVENT_DATA_IO | sctp.SCTP_EVENT_SHUTDOWN | sctp.SCTP_EVENT_ASSOCIATION
		if errSubscribeEvents := newConn.SubscribeEvents(events); errSubscribeEvents != nil {
			logger.NgapLog.Errorf("Failed to accept: %+v", errSubscribeEvents)
			if errSubscribeEvents = newConn.Close(); errSubscribeEvents != nil {
				logger.NgapLog.Errorf("Close error: %+v", errSubscribeEvents)
			}
			continue
		} else {
			logger.NgapLog.Debugln("Subscribe SCTP event[DATA_IO, SHUTDOWN_EVENT, ASSOCIATION_CHANGE]")
		}

		if errSetReadBuffer := newConn.SetReadBuffer(int(readBufSize)); errSetReadBuffer != nil {
			logger.NgapLog.Errorf("Set read buffer error: %+v, accept failed", errSetReadBuffer)
			if errSetReadBuffer = newConn.Close(); errSetReadBuffer != nil {
				logger.NgapLog.Errorf("Close error: %+v", errSetReadBuffer)
			}
			continue
		} else {
			logger.NgapLog.Debugf("Set read buffer to %d bytes", readBufSize)
		}

		if errSetReadTimeout := newConn.SetReadTimeout(readTimeout); errSetReadTimeout != nil {
			logger.NgapLog.Errorf("Set read timeout error: %+v, accept failed", errSetReadTimeout)
			if errSetReadTimeout = newConn.Close(); errSetReadTimeout != nil {
				logger.NgapLog.Errorf("Close error: %+v", errSetReadTimeout)
			}
			continue
		} else {
			logger.NgapLog.Debugf("Set read timeout: %+v", readTimeout)
		}
		logger.NgapLog.Infof("[AMF] SCTP Accept from: %+v", newConn.RemoteAddr())
		connections.Store(newConn, newConn)

		go handleConnection(newConn, readBufSize, handler)
	}
}

func Stop() {
	logger.NgapLog.Infof("Close SCTP server...")
	if err := sctpListener.Close(); err != nil {
		logger.NgapLog.Error(err)
		logger.NgapLog.Infof("SCTP server may not close normally.")
	}

	connections.Range(func(key, value interface{}) bool {
		conn := value.(net.Conn)
		if err := conn.Close(); err != nil {
			logger.NgapLog.Error(err)
		}
		return true
	})

	logger.NgapLog.Infof("SCTP server closed")
}

func handleConnection(conn *sctp.SCTPConn, bufsize uint32, handler NGAPHandler) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.NgapLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}

		// if AMF call Stop(), then conn.Close() will return EBADF because conn has been closed inside Stop()
		if err := conn.Close(); err != nil && err != syscall.EBADF {
			logger.NgapLog.Errorf("close connection error: %+v", err)
		}
		connections.Delete(conn)
	}()

	for {
		buf := make([]byte, bufsize)

		n, info, notification, err := conn.SCTPRead(buf)
		if err != nil {
			switch err {
			case io.EOF, io.ErrUnexpectedEOF:
				logger.NgapLog.Debugln("Read EOF from client")
				handler.HandleConnectionError(conn)
				return
			case syscall.EAGAIN:
				logger.NgapLog.Debugln("SCTP read timeout")
				continue
			case syscall.EINTR:
				logger.NgapLog.Debugf("SCTPRead: %+v", err)
				continue
			default:
				logger.NgapLog.Errorf(
					"Handle connection[addr: %+v] error: %+v",
					conn.RemoteAddr(),
					err,
				)
				handler.HandleConnectionError(conn)
				return
			}
		}

		if notification != nil {
			if handler.HandleNotification != nil {
				handler.HandleNotification(conn, notification)
			} else {
				logger.NgapLog.Warnf("Received sctp notification[type 0x%x] but not handled", notification.Type())
			}
		} else {
			if info == nil || info.PPID != ngap.PPID {
				logger.NgapLog.Warnln("Received SCTP PPID != 60, discard this packet")
				continue
			}

			logger.NgapLog.Tracef("Read %d bytes", n)
			logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(buf[:n]))

			// TODO: concurrent on per-UE message
			handler.HandleMessage(conn, buf[:n])
		}
	}
}
