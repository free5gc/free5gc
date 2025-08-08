package test_test

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"test"
	"test/app"
	"test/consumerTestdata/UDM/TestGenAuthData"

	"github.com/free5gc/nas/security"
	"github.com/free5gc/ngap"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/mongoapi"
	"github.com/stretchr/testify/assert"
)

var NFstructs = []app.NFstruct{}

const (
	noInit = iota
	initNF
	multiAMF
)

var initFlag int = initNF

func init() {
	initNfCfg := test.StartNFsConfig{
		OAuth:  false,
		TestId: "",
	}

	for _, arg := range os.Args {
		values := strings.Split(arg, "=")
		for k, v := range values {
			if v == "-test.run" {
				initNfCfg.TestId = test.TestId(values[k+1])
				fmt.Printf("Run %s\n", initNfCfg.TestId)
			}
		}

		if arg == "noinit" {
			initFlag = noInit
		}

		if arg == "multiAmf" {
			initFlag = multiAMF
		}

		if arg == "oauth" {
			initNfCfg.OAuth = true
		}
	}
	fmt.Println("os.Args:", os.Args)

	switch initFlag {
	case multiAMF:
		initLogPath()

		startNfCfg := test.StartNFsConfig{
			OAuth:  initNfCfg.OAuth,
			TestId: initNfCfg.TestId,

			Nrf:  true,
			Smf:  true,
			Udr:  true,
			Pcf:  true,
			Udm:  true,
			Nssf: true,
			Ausf: true,
			Chf:  true,

			Amf: false,
		}

		NFstructs = test.CreateNFs(startNfCfg)

		NfStart()
		setMongoDB()

	case initNF:
		initLogPath()

		startNfCfg := test.StartNFsConfig{
			OAuth:  initNfCfg.OAuth,
			TestId: initNfCfg.TestId,

			Nrf:  true,
			Amf:  true,
			Smf:  true,
			Udr:  true,
			Pcf:  true,
			Udm:  true,
			Nssf: true,
			Ausf: true,
			Chf:  true,
		}
		NFstructs = test.CreateNFs(startNfCfg)
		NfStart()

	default:
		setMongoDB()
	}

}

func initLogPath() {
	// default key log path
	if err := os.MkdirAll("./log/", 0775); err != nil {
		fmt.Printf("Make directory %s failed: %+v", "./log/", err)
	}
}

func setMongoDB() {
	if err := mongoapi.SetMongoDB("free5gc", "mongodb://127.0.0.1:27017"); err != nil {
		fmt.Printf("SetMongoDB failed: %v\n", err)
		return
	}
	fmt.Println("MongoDB Set")
}

var nfWaitingGroup sync.WaitGroup

func NfStart() {
	fmt.Println("NfStart", len(NFstructs))

	for _, app := range NFstructs {
		go func() {
			nfWaitingGroup.Add(1)
			defer nfWaitingGroup.Done()

			app.Nf.Start()
		}()
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)
}

func NfTerminate() {
	if initFlag != noInit {
		test.NfCancel()
		nfWaitingGroup.Wait()
	}
}

func TestNGSetup(t *testing.T) {
	var n int
	var sendMsg []byte
	var recvMsg = make([]byte, 2048)

	// RAN connect to AMF
	conn, err := test.ConnectToAmf("127.0.0.1", "127.0.0.1", 38412, 9487)
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

func TestCN(t *testing.T) {
	// New UE
	ue := test.NewRanUeContext("imsi-208930000007487", 1, security.AlgCiphering128NEA2, security.AlgIntegrity128NIA2,
		models.AccessType__3_GPP_ACCESS)
	// ue := test.NewRanUeContext("imsi-208930000007487", 1, security.AlgCiphering128NEA0, security.AlgIntegrity128NIA0, models.AccessType__3_GPP_ACCESS)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = test.GetAuthSubscription(TestGenAuthData.MilenageTestSet19.K,
		TestGenAuthData.MilenageTestSet19.OPC,
		TestGenAuthData.MilenageTestSet19.OP)
	// insert UE data to MongoDB

	servingPlmnId := "20893"
	// test.InsertAuthSubscriptionToMongoDB(ue.Supi, ue.AuthenticationSubs)

	test.InsertUeToMongoDB(t, ue, servingPlmnId)

	getData := test.GetAuthSubscriptionFromMongoDB(ue.Supi)
	assert.NotNil(t, getData)
	{
		amData := test.GetAccessAndMobilitySubscriptionData()
		test.InsertAccessAndMobilitySubscriptionDataToMongoDB(ue.Supi, amData, servingPlmnId)
		getData := test.GetAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		smfSelData := test.GetSmfSelectionSubscriptionData()
		test.InsertSmfSelectionSubscriptionDataToMongoDB(ue.Supi, smfSelData, servingPlmnId)
		getData := test.GetSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		smSelData := test.GetSessionManagementSubscriptionData()
		test.InsertSessionManagementSubscriptionDataToMongoDB(ue.Supi, servingPlmnId, smSelData)
		getData := test.GetSessionManagementDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		amPolicyData := test.GetAmPolicyData()
		test.InsertAmPolicyDataToMongoDB(ue.Supi, amPolicyData)
		getData := test.GetAmPolicyDataFromMongoDB(ue.Supi)
		assert.NotNil(t, getData)
	}
	{
		smPolicyData := test.GetSmPolicyData()
		test.InsertSmPolicyDataToMongoDB(ue.Supi, smPolicyData)
		getData := test.GetSmPolicyDataFromMongoDB(ue.Supi)
		assert.NotNil(t, getData)
	}

	defer beforeClose(t, ue, servingPlmnId)

	// subscribe os signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Signal(syscall.SIGUSR1))
	<-c
}

func beforeClose(t *testing.T, ue *test.RanUeContext, servingPlmnId string) {
	// delete test data
	// test.DelAuthSubscriptionToMongoDB(ue.Supi)
	// test.DelAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, "20893")
	// test.DelSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, "20893")

	test.DelUeFromMongoDB(t, ue, servingPlmnId)
}
