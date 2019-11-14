package Namf_Callback_test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/http2_util"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/amf/HttpCallback"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/amf_nas"
	"free5gc/src/nrf/Discovery"
	"free5gc/src/nrf/Management"
	"free5gc/src/nrf/nrf_handler"
	"free5gc/src/smf/PDUSession"
	"free5gc/src/smf/smf_handler"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestSmContextStatusNotify(t *testing.T) {
	go func() {
		router := Namf_Callback.NewRouter()
		server, err := http2_util.NewServer(":29518", TestAmf.AmfLogPath, router)
		if err == nil && server != nil {
			err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
		}
		assert.True(t, err == nil, err.Error())
	}()
	go func() {
		router := PDUSession.NewRouter()
		server, err := http2_util.NewServer(":29502", TestAmf.AmfLogPath, router)
		if err == nil && server != nil {
			err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
			assert.True(t, err == nil, err.Error())
		}
	}()
	go func() {
		router := gin.Default()

		Management.AddService(router)
		Discovery.AddService(router)

		nrfLogPath := path_util.Gofree5gcPath("free5gc/nrfsslkey.log")
		nrfPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/nrf.pem")
		nrfKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/nrf.key")

		server, err := http2_util.NewServer(":29510", nrfLogPath, router)
		if err == nil && server != nil {
			err = server.ListenAndServeTLS(nrfPemPath, nrfKeyPath)
			assert.True(t, err == nil, err.Error())
		}

	}()
	TestAmf.SctpSever()
	go amf_handler.Handle()
	go smf_handler.Handle()
	go nrf_handler.Handle()
	time.Sleep(10 * time.Millisecond)
	//Connect to mongoDB
	MongoDBLibrary.SetMongoDB("free5gc", "mongodb://localhost:27017")
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})

	time.Sleep(10 * time.Millisecond)
	uuid, profile := TestAmf.BuildSmfNfProfile()
	uri, err := amf_consumer.SendRegisterNFInstance("https://localhost:29510", uuid, profile)
	if err != nil {
		t.Error(err.Error())
	} else {
		TestAmf.Config.Dump(uri)
	}

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)

	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	// InitialRequest with known dnn and Snssai (success)
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	m := nas.NewMessage()
	err = m.GmmMessageDecode(&nasPdu)
	err = amf_nas.Dispatch(ue, models.AccessType__3_GPP_ACCESS, ngapType.ProcedureCodeUplinkNASTransport, m)
	assert.True(t, err == nil)
	assert.True(t, err == nil)

	TestAmf.Config.Dump(ue.SmContextList)

	// Notification
	client := &http.Client{}
	client.Transport = &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	smContextStatusNotification := models.SmContextStatusNotification{
		StatusInfo: &models.StatusInfo{
			ResourceStatus: models.ResourceStatus_RELEASED,
		},
	}
	reqByte, err := json.Marshal(&smContextStatusNotification)
	assert.True(t, err == nil)

	reqBuffer := bytes.NewBuffer([]byte(reqByte))

	req, err := http.NewRequest("POST", "https://localhost:29518/namf-callback/v1/smContextStatus/"+ue.Guti+"/10", reqBuffer)
	assert.True(t, err == nil)

	resp, err := client.Do(req)
	assert.True(t, err == nil)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	defer resp.Body.Close()
	TestAmf.Config.Dump(string(body))

	TestAmf.Config.Dump(ue.SmContextList)

}
