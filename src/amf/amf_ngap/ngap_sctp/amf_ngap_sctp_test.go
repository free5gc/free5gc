package amf_ngap_sctp_test

import (
	"encoding/binary"
	"encoding/hex"
	"free5gc/src/amf/amf_handler"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/ishidawataru/sctp"

	amf_ngap_sctp "free5gc/src/amf/amf_ngap/ngap_sctp"
	"free5gc/src/amf/logger"
)

var testClientNum = 2

func TestSctpServer(t *testing.T) {
	runtime.GOMAXPROCS(20)

	go amf_handler.Handle()
	time.Sleep(200 * time.Microsecond)

	amf_ngap_sctp.Server("127.0.0.1")
	logger.NgapLog.Print("Start Client")
	for i := 0; i < testClientNum; i++ {
		time.Sleep(100 * time.Millisecond)
		go func(clientOrder int) {
			testClient(clientOrder)
		}(i)
	}
	time.Sleep(800 * time.Millisecond)
}

func testClient(clientOrder int) {
	logger.NgapLog.Printf("Inside client %d", clientOrder)
	ipStr := "127.0.0.1"
	ips := []net.IPAddr{}
	if ip, err := net.ResolveIPAddr("ip", ipStr); err != nil {
		logger.NgapLog.Errorf("Error resolving address '%s': %v", ipStr, err)
	} else {
		ips = append(ips, *ip)
	}
	addr := &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    38412,
	}
	logger.NgapLog.Printf("raw addr: %+v\n", addr.ToRawSockAddrBuf())

	var laddr *sctp.SCTPAddr
	conn, err := sctp.DialSCTP("sctp", laddr, addr)

	if err != nil {
		logger.NgapLog.Errorf("failed to dial: %v\n", err)
	}
	logger.NgapLog.Printf("Dail LocalAddr: %s; RemoteAddr: %s", conn.LocalAddr(), conn.RemoteAddr())
	time.Sleep(time.Millisecond)
	for {
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, 60)
		ppid := binary.LittleEndian.Uint32(bs)
		info := &sctp.SndRcvInfo{
			Stream: uint16(ppid),
			PPID:   uint32(ppid),
		}
		err := conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
		if err != nil {
			logger.NgapLog.Fatalf("Connection Error %v", err)
		}
		msg, err := hex.DecodeString("00150035000004001B00080002F839104546470052400903006672656535474300660010000000112200208F93000010080102030015400140")
		if err != nil {
			logger.NgapLog.Fatalf("failed to deocde hex string: %v", err)
		}
		n, err := conn.SCTPWrite(msg, info)
		if err != nil {
			logger.NgapLog.Fatalf("failed to write: %v", err)
		}
		logger.NgapLog.Printf("write: %d", n)
		time.Sleep(time.Second)
	}
}
