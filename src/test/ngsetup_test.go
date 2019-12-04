package test_test

import (
	"flag"
	"fmt"
	"free5gc/lib/ngap"
	"free5gc/lib/ngap/ngapSctp"
	"free5gc/lib/path_util"
	"free5gc/src/amf/amf_service"
	"free5gc/src/app"
	"free5gc/src/ausf/ausf_service"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/nssf/nssf_service"
	"free5gc/src/pcf/pcf_service"
	"free5gc/src/smf/smf_service"
	"free5gc/src/test"
	"free5gc/src/udm/udm_service"
	"free5gc/src/udr/udr_service"
	"log"
	"net"
	"testing"
	"time"

	"github.com/ishidawataru/sctp"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

var NFs = []app.NetworkFunction{
	&nrf_service.NRF{},
	&amf_service.AMF{},
	&smf_service.SMF{},
	&udr_service.UDR{},
	&pcf_service.PCF{},
	&udm_service.UDM{},
	&nssf_service.NSSF{},
	&ausf_service.AUSF{},
}

func init() {
	app.AppInitializeWillInitialize("")
	flagSet := flag.NewFlagSet("free5gc", 0)
	flagSet.String("smfcfg", "", "SMF Config Path")
	cli := cli.NewContext(nil, flagSet, nil)
	err := cli.Set("smfcfg", path_util.Gofree5gcPath("free5gc/config/smfcfg.test.conf"))
	if err != nil {
		log.Fatal("SMF test config error")
		return
	}
	for _, service := range NFs {
		service.Initialize(cli)
		go service.Start()
		time.Sleep(200 * time.Millisecond)
	}
}

func getNgapIp(amfIP, ranIP string, amfPort, ranPort int) (amfAddr, ranAddr *sctp.SCTPAddr, err error) {
	ips := []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", amfIP); err1 != nil {
		err = fmt.Errorf("Error resolving address '%s': %v", amfIP, err1)
		return
	} else {
		ips = append(ips, *ip)
	}
	amfAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    amfPort,
	}
	ips = []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", ranIP); err1 != nil {
		err = fmt.Errorf("Error resolving address '%s': %v", ranIP, err1)
		return
	} else {
		ips = append(ips, *ip)
	}
	ranAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    ranPort,
	}
	return
}

func conntectToAmf(amfIP, ranIP string, amfPort, ranPort int) (*sctp.SCTPConn, error) {
	amfAddr, ranAddr, err := getNgapIp(amfIP, ranIP, amfPort, ranPort)
	if err != nil {
		return nil, err
	}
	conn, err := sctp.DialSCTP("sctp", ranAddr, amfAddr)
	if err != nil {
		return nil, err
	}
	info, _ := conn.GetDefaultSentParam()
	info.PPID = ngapSctp.NGAP_PPID
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func TestNGSetup(t *testing.T) {
	var n int
	var sendMsg []byte
	var recvMsg = make([]byte, 2048)

	// RAN connect to AMF
	conn, err := conntectToAmf("127.0.0.1", "127.0.0.1", 38412, 9487)
	assert.Nil(t, err)

	// send NGSetupRequest Msg
	sendMsg, err = test.GetNGSetupRequest([]byte("\x00\x01\x02"), 24, "free5gc")
	assert.Nil(t, err)
	_, err = conn.Write(sendMsg)
	assert.Nil(t, err)

	// receive NGSetupResponse Msg
	n, err = conn.Read(recvMsg)
	assert.Nil(t, err)
	_, err = ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)

	// close Connection
	conn.Close()
}
