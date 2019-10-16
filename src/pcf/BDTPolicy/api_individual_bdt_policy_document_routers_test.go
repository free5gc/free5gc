package BDTPolicy_test

import (
	"context"
	TestBDTPolicy "free5gc/lib/CommonConsumerTestData/PCF/TestBDTPolicy"
	Npcf_BDTPolicy "free5gc/lib/Npcf_BDTPolicy"
	"free5gc/lib/http2_util"
	BDTPolicy "free5gc/src/pcf/BDTPolicy"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler"
	"free5gc/src/pcf/pcf_util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUpdateBDTPolicy(t *testing.T) {
	go func() {
		pcfRouter := BDTPolicy.NewRouter()
		pcfServer, err := http2_util.NewServer(":29507", pcf_util.PCF_LOG_PATH, pcfRouter)
		if err == nil && pcfServer != nil {
			err := pcfServer.ListenAndServeTLS(pcf_util.PCF_PEM_PATH, pcf_util.PCF_KEY_PATH)
			assert.True(t, err == nil)
		}
	}()

	go pcf_handler.Handle()
	configuration := Npcf_BDTPolicy.NewConfiguration()
	configuration.SetBasePath(pcf_util.PCF_BASIC_PATH + pcf_context.BdtUri)
	client := Npcf_BDTPolicy.NewAPIClient(configuration)

	// get test data
	bdtReqData1, bdtReqData2, bdtReqData3, bdtReqDataNil := TestBDTPolicy.GetCreateTestData()
	bdtPolicyDataPatch := TestBDTPolicy.GetUpdateTestData()
	pcf_context.NewPCFUe(bdtReqData1.AspId)
	pcf_context.NewPCFUe(bdtReqData2.AspId)
	pcf_context.NewPCFUe(bdtReqData3.AspId)
	pcf_context.NewPCFUe(bdtReqDataNil.AspId)
	pcf_context.AddAspIdToUe(bdtReqData1.AspId, bdtReqData1.AspId)
	pcf_context.AddAspIdToUe(bdtReqData2.AspId, bdtReqData2.AspId)
	pcf_context.AddAspIdToUe(bdtReqData3.AspId, bdtReqData3.AspId)
	pcf_context.AddAspIdToUe(bdtReqDataNil.AspId, bdtReqDataNil.AspId)
	// test create service
	_, httpRsp, err := client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqData1)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)

	// test get service
	_, httpRsp, err = client.IndividualBDTPolicyDocumentApi.GetBDTPolicy(context.Background(), pcf_context.DefaultBdtRefId+bdtReqData1.AspId)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	// test update service
	_, httpRsp, err = client.IndividualBDTPolicyDocumentApi.UpdateBDTPolicy(context.Background(), pcf_context.DefaultBdtRefId+bdtReqData1.AspId, bdtPolicyDataPatch)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "204 No Content", httpRsp.Status)

	///////////// test error////////////////
	// test nil
	_, httpRsp, err = client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqDataNil)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "404 Not Found", httpRsp.Status)

	/////////////  Multiple data //////////////
	// test mutiple create
	_, httpRsp, err = client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqData2)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)

	_, httpRsp, err = client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqData3)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)
	// test mutiple get
	_, httpRsp, err = client.IndividualBDTPolicyDocumentApi.GetBDTPolicy(context.Background(), pcf_context.DefaultBdtRefId+bdtReqData2.AspId)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	_, httpRsp, err = client.IndividualBDTPolicyDocumentApi.GetBDTPolicy(context.Background(), pcf_context.DefaultBdtRefId+bdtReqData3.AspId)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	// recreate test
	_, httpRsp, err = client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqData1)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

}
