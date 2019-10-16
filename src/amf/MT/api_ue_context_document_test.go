package Namf_MT_test

import (
	"context"
	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	Namf_MT_Clinet "free5gc/lib/Namf_MT"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	Namf_MT_Server "free5gc/src/amf/MT"
	"free5gc/src/amf/amf_handler"
	"log"
	"testing"
	"time"
)

func sendRequestAndPrintResult(client *Namf_MT_Clinet.APIClient, supi string, request *Namf_MT_Clinet.ProvideDomainSelectionInfoParamOpts) {
	ueContextInfo, httpResponse, err := client.UeContextDocumentApi.ProvideDomainSelectionInfo(context.Background(), supi, request)
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

func TestProvideDomainSelectionInfo(t *testing.T) {
	go func() {
		router := Namf_MT_Server.NewRouter()
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
	configuration := Namf_MT_Clinet.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_MT_Clinet.NewAPIClient(configuration)

	/* init ue info*/
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	anType := ue.GetAnType()
	ue.RanUe[anType].SupportVoPSn3gpp = false
	ue.RanUe[anType].SupportVoPS = false
	ue.RanUe[anType].SupportedFeatures = "nothing support"
	time := time.Now()
	ue.RanUe[anType].LastActTime = &time
	ue.RatType = "RatType_NR"
	ue.RanUe[anType].Ran.AnType = "AccessType__3_GPP_ACCESS"
	ue.Supi = "imsi-2089300007487"
	var ProvideDomainSelectionInfoParamOpts Namf_MT_Clinet.ProvideDomainSelectionInfoParamOpts

	//without info-class
	sendRequestAndPrintResult(client, ue.Supi, &ProvideDomainSelectionInfoParamOpts)
	ProvideDomainSelectionInfoParamOpts.InfoClass = optional.NewInterface("TADS")

	sendRequestAndPrintResult(client, ue.Supi, &ProvideDomainSelectionInfoParamOpts)
	// 404 CONTEXT_NOT_FOUND
	sendRequestAndPrintResult(client, "imsi-0010202", &ProvideDomainSelectionInfoParamOpts)
}
