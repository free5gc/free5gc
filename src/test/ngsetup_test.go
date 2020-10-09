package test_test

import (
	"flag"
	"fmt"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/nas/security"
	"free5gc/lib/ngap"
	"free5gc/lib/path_util"
	amf_service "free5gc/src/amf/service"
	"free5gc/src/app"
	ausf_service "free5gc/src/ausf/service"
	nrf_service "free5gc/src/nrf/service"
	nssf_service "free5gc/src/nssf/service"
	pcf_service "free5gc/src/pcf/service"
	smf_service "free5gc/src/smf/service"
	"free5gc/src/test"
	udm_service "free5gc/src/udm/service"
	udr_service "free5gc/src/udr/service"
	"log"
	"os"
	"sync"
	"testing"
	"time"

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
	//&n3iwf_service.N3IWF{},
}

func init() {
	var init bool = true

	for _, arg := range os.Args {
		if arg == "noinit" {
			init = false
		}
	}

	if init {
		app.AppInitializeWillInitialize("")
		flagSet := flag.NewFlagSet("free5gc", 0)
		flagSet.String("smfcfg", "", "SMF Config Path")
		cli := cli.NewContext(nil, flagSet, nil)
		err := cli.Set("smfcfg", path_util.Gofree5gcPath("free5gc/config/test/smfcfg.test.conf"))
		if err != nil {
			log.Fatal("SMF test config error")
			return
		}

		for _, service := range NFs {
			service.Initialize(cli)
			go service.Start()
			time.Sleep(200 * time.Millisecond)
		}
	} else {
		MongoDBLibrary.SetMongoDB("free5gc", "mongodb://127.0.0.1:27017")
		fmt.Println("MongoDB Set")
	}

}

func TestNGSetup(t *testing.T) {
	var n int
	var sendMsg []byte
	var recvMsg = make([]byte, 2048)

	// RAN connect to AMF
	conn, err := test.ConntectToAmf("127.0.0.1", "127.0.0.1", 38412, 9487)
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
	ue := test.NewRanUeContext("imsi-2089300007487", 1, security.AlgCiphering128NEA2, security.AlgIntegrity128NIA2)
	// ue := test.NewRanUeContext("imsi-2089300007487", 1, security.AlgCiphering128NEA0, security.AlgIntegrity128NIA0)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = test.GetAuthSubscription(TestGenAuthData.MilenageTestSet19.K,
		TestGenAuthData.MilenageTestSet19.OPC,
		TestGenAuthData.MilenageTestSet19.OP)
	// insert UE data to MongoDB

	servingPlmnId := "20893"
	test.InsertAuthSubscriptionToMongoDB(ue.Supi, ue.AuthenticationSubs)
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

	defer beforeClose(ue)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func beforeClose(ue *test.RanUeContext) {
	// delete test data
	test.DelAuthSubscriptionToMongoDB(ue.Supi)
	test.DelAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, "20893")
	test.DelSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, "20893")
}
