package Communication_test

import (
	"context"
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/AMF/TestComm"
	Namf_Communication_Client "free5gc/lib/Namf_Communication"
	"free5gc/lib/http2_util"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	Namf_Communication_Server "free5gc/src/amf/Communication"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/amf_producer/amf_producer_callback"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func sendRequestAndPrintResult(client *Namf_Communication_Client.APIClient, supi string, request models.UeN1N2InfoSubscriptionCreateData) {
	ueContextInfo, httpResponse, err := client.N1N2SubscriptionsCollectionForIndividualUEContextsDocumentApi.N1N2MessageSubscribe(context.Background(), supi, request)
	if err != nil {
		if httpResponse == nil {
			log.Panic(err)
		} else if err.Error() != httpResponse.Status {
			log.Panic(err)
		} else {
			var probelmDetail models.ProblemDetails
			probelmDetail = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			TestAmf.Config.Dump(probelmDetail)
		}
	} else {
		TestAmf.Config.Dump(ueContextInfo)
	}
}

func TestN1N2MessageSubscribe(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		go func() {
			router := Namf_Communication_Server.NewRouter()
			server, err := http2_util.NewServer(":29518", TestAmf.AmfLogPath, router)
			if err == nil && server != nil {
				err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
			}
			assert.True(t, err == nil)
		}()

		go amf_handler.Handle()
		TestAmf.AmfInit()
		TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
		time.Sleep(100 * time.Millisecond)
	}
	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	/* init ue info*/
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ueN1N2InfoSubscriptionCreateData := TestComm.ConsumerAMFN1N2MessageSubscribeRequsetTable[TestComm.UeN1N2InfoSubsriptionCreateData]
	sendRequestAndPrintResult(client, ue.Supi, *ueN1N2InfoSubscriptionCreateData)
}

func TestN1MessageNotify(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestN1N2MessageSubscribe(t)
	}
	go func() {
		keylogFile, err := os.OpenFile(TestAmf.AmfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		assert.True(t, err == nil)
		server := http.Server{
			Addr: ":29507",
			TLSConfig: &tls.Config{
				KeyLogWriter: keylogFile,
			},
		}
		http2.ConfigureServer(&server, nil)
		http.HandleFunc("/n1MessageNotify", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
		assert.True(t, err == nil)
	}()
	time.Sleep(200 * time.Millisecond)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	n1msg := []byte{0x12}
	amf_producer_callback.SendN1MessageNotify(ue, models.N1MessageClass__5_GMM, n1msg, nil)
}

func TestN2InfoNotifyUri(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestN1N2MessageSubscribe(t)
	}
	go func() {
		keylogFile, err := os.OpenFile(TestAmf.AmfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		assert.True(t, err == nil)
		server := http.Server{
			Addr: ":29503",
			TLSConfig: &tls.Config{
				KeyLogWriter: keylogFile,
			},
		}
		http2.ConfigureServer(&server, nil)
		http.HandleFunc("/n2InfoNotify", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
		assert.True(t, err == nil)
	}()
	time.Sleep(200 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	tmp := []byte{0x12}
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", nil)
	amf_producer_callback.SendN2InfoNotify(ue, models.N2InformationClass_NRP_PA, nasPdu, tmp)
}
