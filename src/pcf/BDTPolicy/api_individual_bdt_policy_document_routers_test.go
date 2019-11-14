package BDTPolicy_test

import (
	"context"
	TestBDTPolicy "free5gc/lib/CommonConsumerTestData/PCF/TestBDTPolicy"
	Npcf_BDTPolicy "free5gc/lib/Npcf_BDTPolicy"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUpdateBDTPolicy(t *testing.T) {
	pcfInit()

	configuration := Npcf_BDTPolicy.NewConfiguration()
	configuration.SetBasePath(pcf_util.PCF_BASIC_PATH + pcf_context.BdtUri)
	client := Npcf_BDTPolicy.NewAPIClient(configuration)

	// get test data
	bdtReqData1, bdtReqData2, bdtReqData3, bdtReqDataNil := TestBDTPolicy.GetCreateTestData()
	bdtPolicyDataPatch := TestBDTPolicy.GetUpdateTestData()

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
