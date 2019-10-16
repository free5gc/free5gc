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

func TestCreateBDTPolicy(t *testing.T) {
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
	bdtReqData1, _, _, _ := TestBDTPolicy.GetCreateTestData()
	pcf_context.NewPCFUe(bdtReqData1.AspId)
	pcf_context.AddAspIdToUe(bdtReqData1.AspId, bdtReqData1.AspId)

	// test create service
	_, httpRsp, err := client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqData1)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)

}
