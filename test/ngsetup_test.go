package test

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/calee0219/fatal"
	"github.com/free5gc/test/app"
	"github.com/free5gc/test/consumerTestdata/UDM/TestGenAuthData"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/nas/security"
	"github.com/free5gc/ngap"
	"github.com/free5gc/openapi/models"
)

var NFs = []app.NetworkFunction{}

const (
	noInit = iota
	initNF
	multiAMF
)

var initFlag int = initNF

const (
	ranNgapPort    int    = 9487
	amfNgapPort    int    = 38412
	ngapPPID       uint32 = 0x3c000000
	ranN2Ipv4Addr  string = "127.0.0.1"
	amfN2Ipv4Addr  string = "127.0.0.1"
	amfN2Ipv4Addr2 string = "127.0.0.50"
	ranN3Ipv4Addr  string = "10.200.200.1"
	upfN3Ipv4Addr  string = "10.200.200.102"
)

func NfTerminate() {
	if initFlag != noInit {
		nfNums := len(NFs)
		for i := nfNums - 1; i >= 0; i-- {
			NFs[i].Terminate()
		}
	}
}

func getNgapIp(amfIP, ranIP string, amfPort, ranPort int) (amfAddr, ranAddr *sctp.SCTPAddr, err error) {
	ips := []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", amfIP); err1 != nil {
		err = fmt.Errorf("Error resolving address '%s': %v", amfIP, err1)
		return nil, nil, err
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
		return nil, nil, err
	} else {
		ips = append(ips, *ip)
	}
	ranAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    ranPort,
	}
	return amfAddr, ranAddr, nil
}

func ConnectToAmf(amfIP, ranIP string, amfPort, ranPort int) (*sctp.SCTPConn, error) {
	amfAddr, ranAddr, err := getNgapIp(amfN2Ipv4Addr, ranN2Ipv4Addr, amfPort, ranPort)
	if err != nil {
		return nil, err
	}
	conn, err := sctp.DialSCTP("sctp", ranAddr, amfAddr)
	if err != nil {
		return nil, err
	}
	info, err := conn.GetDefaultSentParam()
	if err != nil {
		fatal.Fatalf("conn GetDefaultSentParam error in ConnectToAmf: %+v", err)
	}
	info.PPID = ngapPPID
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
	conn, err := ConnectToAmf(amfN2Ipv4Addr, ranN2Ipv4Addr, amfNgapPort, ranNgapPort)
	assert.Nil(t, err)

	// send NGSetupRequest Msg
	sendMsg, err = GetNGSetupRequest([]byte("\x00\x01\x02"), 24, "free5gc")
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

func TestCN(t *testing.T) {
	// New UE
	ue := NewRanUeContext("imsi-2089300007487", 1, security.AlgCiphering128NEA2, security.AlgIntegrity128NIA2,
		models.AccessType__3_GPP_ACCESS)
	// ue := NewRanUeContext("imsi-2089300007487", 1, security.AlgCiphering128NEA0, security.AlgIntegrity128NIA0, models.AccessType__3_GPP_ACCESS)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = GetAuthSubscription(TestGenAuthData.MilenageTestSet19.K,
		TestGenAuthData.MilenageTestSet19.OPC,
		TestGenAuthData.MilenageTestSet19.OP)
	// insert UE data to MongoDB

	servingPlmnId := "20893"
	InsertAuthSubscriptionToMongoDB(ue.Supi, ue.AuthenticationSubs)
	getData := GetAuthSubscriptionFromMongoDB(ue.Supi)
	assert.NotNil(t, getData)
	{
		amData := GetAccessAndMobilitySubscriptionData()
		InsertAccessAndMobilitySubscriptionDataToMongoDB(ue.Supi, amData, servingPlmnId)
		getData := GetAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		smfSelData := GetSmfSelectionSubscriptionData()
		InsertSmfSelectionSubscriptionDataToMongoDB(ue.Supi, smfSelData, servingPlmnId)
		getData := GetSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		smSelData := GetSessionManagementSubscriptionData()
		InsertSessionManagementSubscriptionDataToMongoDB(ue.Supi, servingPlmnId, smSelData)
		getData := GetSessionManagementDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		amPolicyData := GetAmPolicyData()
		InsertAmPolicyDataToMongoDB(ue.Supi, amPolicyData)
		getData := GetAmPolicyDataFromMongoDB(ue.Supi)
		assert.NotNil(t, getData)
	}
	{
		smPolicyData := GetSmPolicyData()
		InsertSmPolicyDataToMongoDB(ue.Supi, smPolicyData)
		getData := GetSmPolicyDataFromMongoDB(ue.Supi)
		assert.NotNil(t, getData)
	}

	defer beforeClose(ue)

	// subscribe os signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Signal(syscall.SIGUSR1))
	<-c
}

func beforeClose(ue *RanUeContext) {
	// delete test data
	DelAuthSubscriptionToMongoDB(ue.Supi)
	DelAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, "20893")
	DelSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, "20893")
}
