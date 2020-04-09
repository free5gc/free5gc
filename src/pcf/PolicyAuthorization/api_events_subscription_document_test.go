package PolicyAuthorization_test

import (
	"context"
	"fmt"
	"free5gc/lib/CommonConsumerTestData/PCF/TestAMPolicy"
	"free5gc/lib/CommonConsumerTestData/PCF/TestPolicyAuthorization"
	"free5gc/lib/CommonConsumerTestData/PCF/TestSMPolicy"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Npcf_AMPolicy"
	"free5gc/lib/Npcf_PolicyAuthorization"
	"free5gc/lib/Npcf_SMPolicyControl"
	"free5gc/src/pcf/pcf_context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventSubscription(t *testing.T) {
	defer MongoDBLibrary.RestfulAPIDeleteMany(amPolicyDataColl, filterUeIdOnly)
	defer MongoDBLibrary.RestfulAPIDeleteMany(smPolicyDataColl, filterUeIdOnly)

	fakeSmfServer(t, ":29502")

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	amclient := Npcf_AMPolicy.NewAPIClient(configuration)
	configuration1 := Npcf_SMPolicyControl.NewConfiguration()
	configuration1.SetBasePath("https://127.0.0.1:29507")
	smclient := Npcf_SMPolicyControl.NewAPIClient(configuration1)
	configuration2 := Npcf_PolicyAuthorization.NewConfiguration()
	configuration2.SetBasePath("https://127.0.0.1:29507")
	afclient := Npcf_PolicyAuthorization.NewAPIClient(configuration2)
	smPolicyId := ""
	//Test PostPolicies
	{
		amCreateReqData := TestAMPolicy.GetAMreqdata()
		_, httpRsp, err := amclient.DefaultApi.PoliciesPost(context.Background(), amCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			index := strings.LastIndex(locationHeader, "/")
			assert.True(t, index != -1)
			polAssoId := locationHeader[index+1:]
			assert.True(t, strings.HasPrefix(polAssoId, "imsi-2089300007487"))
		}
	}
	ue := pcf_context.PCF_Self().UePool["imsi-2089300007487"]
	pcfSelf := pcf_context.PCF_Self()
	{
		smCreateReqData := TestSMPolicy.CreateTestData()
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesPost(context.Background(), smCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			index := strings.LastIndex(locationHeader, "/")
			assert.True(t, index != -1)
			smPolicyId = locationHeader[index+1:]
			assert.True(t, locationHeader == "https://127.0.0.1:29507/npcf-smpolicycontrol/v1/sm-policies/imsi-2089300007487-1")
		}
	}
	appSessionId := ""
	{
		// Create App Session
		afReqData := TestPolicyAuthorization.GetPostAppSessionsData_NoEvent()
		_, httpRsp, err := afclient.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), afReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			appSessionId = fmt.Sprintf("%s-%d", ue.Supi, ue.AppSessionIdGenerator-1)
			header := fmt.Sprintf("https://127.0.0.1:29507/npcf-policyauthorization/v1/app-sessions/%s", appSessionId)
			assert.True(t, locationHeader == header)
			appSession := pcfSelf.AppSessionPool[appSessionId]
			if assert.NotNil(t, appSession) {
				smPolicy := ue.SmPolicyData[smPolicyId]
				assert.Equal(t, smPolicy, appSession.SmPolicyData)
				pccRuleId := fmt.Sprintf("PccRuleId-%d", smPolicy.PccRuleIdGenarator-1)
				assert.Equal(t, 1, len(appSession.RelatedPccRuleIds))
				pccRule := smPolicy.PolicyDecision.PccRules[pccRuleId]
				if assert.NotNil(t, pccRule) {
					assert.Equal(t, afReqData.AscReqData.MedComponents["1"].MedSubComps["1"].FDescs[0], pccRule.FlowInfos[0].FlowDescription)
				}
			}
		}
	}
	appSession := pcfSelf.AppSessionPool[appSessionId]
	{
		// Create App Session Subscription (201)
		reqData := TestPolicyAuthorization.GetUpdateEventsSubsc201Data()
		resp, httpRsp, err := afclient.EventsSubscriptionDocumentApi.UpdateEventsSubsc(context.Background(), appSessionId, reqData)
		assert.True(t, err == nil)
		if assert.True(t, httpRsp != nil) {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			if assert.NotNil(t, resp) {
				assert.Equal(t, reqData, resp.EvSubsc)
				assert.NotNil(t, resp.EvsNotif)
				assert.Equal(t, 2, len(appSession.Events))
				assert.Equal(t, reqData.NotifUri, appSession.EventUri)
			}
		}
	}
	{
		// Modify App Session Subscription (200)
		reqData := TestPolicyAuthorization.GetUpdateEventsSubsc200Data()
		resp, httpRsp, err := afclient.EventsSubscriptionDocumentApi.UpdateEventsSubsc(context.Background(), appSessionId, reqData)
		assert.True(t, err == nil)
		if assert.True(t, httpRsp != nil) {
			assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
			if assert.NotNil(t, resp) {
				assert.Equal(t, reqData, resp.EvSubsc)
				assert.NotNil(t, resp.EvsNotif)
				assert.Equal(t, 1, len(appSession.Events))
				assert.Equal(t, reqData.NotifUri, appSession.EventUri)
			}
		}
	}
	{
		// Modify App Session Subscription (204)
		reqData := TestPolicyAuthorization.GetUpdateEventsSubsc204Data()
		_, httpRsp, err := afclient.EventsSubscriptionDocumentApi.UpdateEventsSubsc(context.Background(), appSessionId, reqData)
		assert.True(t, err == nil)
		if assert.True(t, httpRsp != nil) {
			assert.Equal(t, http.StatusNoContent, httpRsp.StatusCode)
			assert.Equal(t, 1, len(appSession.Events))
			assert.Equal(t, reqData.NotifUri, appSession.EventUri)
		}
	}
	{
		// Create App Session Subscription (400)
		reqData := TestPolicyAuthorization.GetUpdateEventsSubsc400Data()
		_, httpRsp, err := afclient.EventsSubscriptionDocumentApi.UpdateEventsSubsc(context.Background(), appSessionId, reqData)
		assert.True(t, err != nil)
		if assert.True(t, httpRsp != nil) {
			assert.Equal(t, http.StatusBadRequest, httpRsp.StatusCode)
		}
	}
	{
		// Delete App Session Subscription (204)
		httpRsp, err := afclient.EventsSubscriptionDocumentApi.DeleteEventsSubsc(context.Background(), appSessionId)
		assert.True(t, err == nil)
		if assert.True(t, httpRsp != nil) {
			assert.Equal(t, http.StatusNoContent, httpRsp.StatusCode)
			assert.Nil(t, appSession.Events)
		}
	}

	time.Sleep(100 * time.Millisecond)
}
