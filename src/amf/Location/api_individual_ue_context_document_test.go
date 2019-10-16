package Namf_Location_test

import (
	"context"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	Namf_Loc_Clinet "free5gc/lib/Namf_Location"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	Namf_Loc_Server "free5gc/src/amf/Location"
	"free5gc/src/amf/amf_handler"
	"log"

	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func sendRequestAndPrintResult(client *Namf_Loc_Clinet.APIClient, supi string, request models.RequestLocInfo) {
	ueContextInfo, httpResponse, err := client.IndividualUEContextDocumentApi.ProvideLocationInfo(context.Background(), supi, request)
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

func TestProvideLocationInfo(t *testing.T) {
	go func() {
		router := Namf_Loc_Server.NewRouter()
		server, err := http2_util.NewServer(":29518", TestAmf.AmfLogPath, router)
		if err == nil && server != nil {
			err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
		}
		assert.True(t, err == nil, err.Error())
	}()

	go amf_handler.Handle()
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)
	configuration := Namf_Loc_Clinet.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29518")
	client := Namf_Loc_Clinet.NewAPIClient(configuration)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.Supi = "imsi-2089300007487"
	var requestLocInfo models.RequestLocInfo

	sendRequestAndPrintResult(client, ue.Supi, requestLocInfo)

	requestLocInfo.Req5gsLoc = true
	requestLocInfo.ReqCurrentLoc = true
	requestLocInfo.ReqRatType = true
	requestLocInfo.ReqTimeZone = true
	sendRequestAndPrintResult(client, ue.Supi, requestLocInfo)

	// 404 CONTEXT_NOT_FOUND
	sendRequestAndPrintResult(client, "imsi-0010202", requestLocInfo)
}
