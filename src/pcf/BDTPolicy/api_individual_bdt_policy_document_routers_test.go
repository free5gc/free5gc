package BDTPolicy_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	TestBDTPolicy "free5gc/lib/CommonConsumerTestData/PCF/TestBDTPolicy"
	"free5gc/lib/Npcf_BDTPolicyControl"
	"free5gc/lib/openapi/models"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGetUpdateBDTPolicy(t *testing.T) {

	configuration := Npcf_BDTPolicyControl.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	client := Npcf_BDTPolicyControl.NewAPIClient(configuration)

	var bdtPolicyId string
	// get test data
	bdtReqData := TestBDTPolicy.GetCreateTestData()

	// test create service
	{
		rsp, httpRsp, err := client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if assert.Equal(t, http.StatusCreated, httpRsp.StatusCode) {
			if assert.NotNil(t, rsp.BdtReqData) {
				assert.Equal(t, *rsp.BdtReqData, bdtReqData)
			}
			if assert.NotNil(t, rsp.BdtPolData) {
				assert.True(t, rsp.BdtPolData.SelTransPolicyId == 1)
				assert.Equal(t, rsp.BdtPolData.SuppFeat, "")
				if assert.True(t, len(rsp.BdtPolData.TransfPolicies) == 1) {
					assert.Equal(t, rsp.BdtPolData.TransfPolicies[0], models.TransferPolicy{
						RatingGroup:   1,
						RecTimeInt:    bdtReqData.DesTimeInt,
						TransPolicyId: 1,
					})
				}
			}
			locationHeader := httpRsp.Header.Get("Location")
			index := strings.LastIndex(locationHeader, "/")
			assert.True(t, index != -1)
			bdtPolicyId = locationHeader[index+1:]
			assert.True(t, strings.HasPrefix(bdtPolicyId, "BdtPolicyId-1"))

		}
	}

	time.Sleep(30 * time.Millisecond)
	// test update service
	{
		bdtPolicyDataPatch := TestBDTPolicy.GetUpdateTestData()
		_, httpRsp, err := client.IndividualBDTPolicyDocumentApi.UpdateBDTPolicy(context.Background(), bdtPolicyId, bdtPolicyDataPatch)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
	}
	time.Sleep(30 * time.Millisecond)
	// test get service
	{
		rsp, httpRsp, err := client.IndividualBDTPolicyDocumentApi.GetBDTPolicy(context.Background(), bdtPolicyId)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if assert.Equal(t, http.StatusOK, httpRsp.StatusCode) {
			if assert.NotNil(t, rsp.BdtReqData) {
				assert.Equal(t, *rsp.BdtReqData, bdtReqData)
			}
			if assert.NotNil(t, rsp.BdtPolData) {
				assert.True(t, rsp.BdtPolData.SelTransPolicyId == 1)
				assert.Equal(t, rsp.BdtPolData.SuppFeat, "")
				if assert.True(t, len(rsp.BdtPolData.TransfPolicies) == 1) {
					assert.Equal(t, rsp.BdtPolData.TransfPolicies[0], models.TransferPolicy{
						RatingGroup:   1,
						RecTimeInt:    bdtReqData.DesTimeInt,
						TransPolicyId: 1,
					})
				}
			}
		}
	}

}
