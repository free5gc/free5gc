package PolicyAuthorization_test

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"free5gc/lib/CommonConsumerTestData/PCF/TestAMPolicy"
	"free5gc/lib/CommonConsumerTestData/PCF/TestPolicyAuthorization"
	"free5gc/lib/CommonConsumerTestData/PCF/TestSMPolicy"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Npcf_AMPolicy"
	"free5gc/lib/Npcf_PolicyAuthorization"
	"free5gc/lib/Npcf_SMPolicyControl"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/pcf/pcf_context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func fakeAfServer(t *testing.T, port string) {
	go func() { // fake af server
		router := gin.Default()

		router.POST("/notify", func(c *gin.Context) {
			fmt.Println("==========App Session Event Notification Callback=============")

			var notification models.EventsNotification
			if err := c.ShouldBindJSON(&notification); err != nil {
				fmt.Println("fake AF server error")
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			spew.Dump(notification.EvNotifs)
			c.JSON(http.StatusNoContent, gin.H{})
		})

		router.POST("/terminate", func(c *gin.Context) {
			fmt.Println("==========App Session Teimination Callback=============")

			var terminationInfo models.TerminationInfo
			if err := c.ShouldBindJSON(&terminationInfo); err != nil {
				fmt.Println("fake AF server error")
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			spew.Dump(terminationInfo)
			c.JSON(http.StatusNoContent, gin.H{})
		})

		pcfPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/pcf.pem")
		pcfKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/pcf.key")

		server := &http.Server{
			Addr:    port,
			Handler: router,
		}

		fmt.Println(server.ListenAndServeTLS(pcfPemPath, pcfKeyPath))
	}()
}

func TestIndividualAppSessionContext(t *testing.T) {
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
		afReqData := TestPolicyAuthorization.GetPostAppSessionsData_Normal()
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
				assert.True(t, appSession.EventUri == afReqData.AscReqData.EvSubsc.NotifUri)
				assert.True(t, len(appSession.Events) == len(afReqData.AscReqData.EvSubsc.Events))
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
	{
		// Get App Session
		resp, httpRsp, err := afclient.IndividualApplicationSessionContextDocumentApi.GetAppSession(context.Background(), appSessionId)
		assert.True(t, err == nil)
		if assert.True(t, httpRsp != nil) {
			assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
			if assert.NotNil(t, resp) {
				assert.NotNil(t, resp.AscRespData)
				assert.NotNil(t, resp.AscReqData)
			}
		}
	}
	appSession := pcfSelf.AppSessionPool[appSessionId]
	smPolicy := ue.SmPolicyData[smPolicyId]
	{
		// Update App Session
		modData := TestPolicyAuthorization.GetModAppSession200Data()
		_, httpRsp, err := afclient.IndividualApplicationSessionContextDocumentApi.ModAppSession(context.Background(), appSessionId, modData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
			if assert.NotNil(t, appSession) {
				assert.True(t, appSession.EventUri == appSession.AppSessionContext.AscReqData.EvSubsc.NotifUri)
				assert.True(t, len(appSession.Events) == len(modData.EvSubsc.Events))
				assert.Equal(t, smPolicy, appSession.SmPolicyData)
				pccRuleId := fmt.Sprintf("PccRuleId-%d", smPolicy.PccRuleIdGenarator-1)
				assert.Equal(t, 1, len(appSession.RelatedPccRuleIds))
				assert.Equal(t, 1, len(smPolicy.PolicyDecision.PccRules))
				pccRule := smPolicy.PolicyDecision.PccRules[pccRuleId]
				if assert.NotNil(t, pccRule) {
					assert.Equal(t, modData.MedComponents["1"].MedSubComps["1"].FDescs[0], pccRule.FlowInfos[0].FlowDescription)
				}
			}
		}
	}
	{
		// Delete App Session
		_, httpRsp, err := afclient.IndividualApplicationSessionContextDocumentApi.DeleteAppSession(context.Background(), appSessionId, nil)
		assert.True(t, err == nil)
		if assert.True(t, httpRsp != nil) {
			assert.Equal(t, http.StatusNoContent, httpRsp.StatusCode)
			assert.Equal(t, 0, len(pcfSelf.AppSessionPool))
			assert.Equal(t, 0, len(smPolicy.AppSessions))
			assert.Equal(t, 0, len(smPolicy.PolicyDecision.PccRules))
		}
	}
}

func TestAppSessionNotification(t *testing.T) {

	defer MongoDBLibrary.RestfulAPIDeleteMany(amPolicyDataColl, filterUeIdOnly)
	defer MongoDBLibrary.RestfulAPIDeleteMany(smPolicyDataColl, filterUeIdOnly)

	fakeSmfServer(t, ":29502")
	fakeAfServer(t, ":12345")

	time.Sleep(100 * time.Millisecond)

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
		afReqData := TestPolicyAuthorization.GetPostAppSessionsData_Normal()
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
				assert.True(t, appSession.EventUri == afReqData.AscReqData.EvSubsc.NotifUri)
				assert.True(t, len(appSession.Events) == len(afReqData.AscReqData.EvSubsc.Events))
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
	time.Sleep(100 * time.Millisecond)

	// Test AppSessionEventNotify (Send SM policy Update which is related to App Session)
	{
		trigger := []models.PolicyControlRequestTrigger{
			models.PolicyControlRequestTrigger_PLMN_CH,
			models.PolicyControlRequestTrigger_AC_TY_CH,
			models.PolicyControlRequestTrigger_RAT_TY_CH,
		}
		updateReq := TestSMPolicy.UpdateTestData(trigger, nil)
		updateReq.AccessType = models.AccessType_NON_3_GPP_ACCESS
		updateReq.RatType = models.RatType_WLAN
		//Test UpdatePoliciesPolAssoId
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdUpdatePost(context.Background(), smPolicyId, updateReq)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
		}
	}
	time.Sleep(100 * time.Millisecond)
	// Test AppSession Termination (Send SM policy Update which is related to App Session)
	{
		//Test DelPoliciesPolAssoId
		httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdDeletePost(context.Background(), smPolicyId, models.SmPolicyDeleteData{})
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusNoContent, httpRsp.StatusCode)
		}
	}
	time.Sleep(100 * time.Millisecond)
	{
		// Get App Session
		_, httpRsp, err := afclient.IndividualApplicationSessionContextDocumentApi.GetAppSession(context.Background(), appSessionId)
		assert.True(t, err != nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusNotFound, httpRsp.StatusCode)
			problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			assert.Equal(t, "APPLICATION_SESSION_CONTEXT_NOT_FOUND", problem.Cause)
		}
	}
	assert.True(t, len(pcf_context.PCF_Self().AppSessionPool) == 0)
	assert.True(t, len(ue.SmPolicyData) == 0)

}
