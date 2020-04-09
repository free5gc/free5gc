package PolicyAuthorization_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/CommonConsumerTestData/PCF/TestAMPolicy"
	"free5gc/lib/CommonConsumerTestData/PCF/TestPolicyAuthorization"
	"free5gc/lib/CommonConsumerTestData/PCF/TestSMPolicy"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Npcf_AMPolicy"
	"free5gc/lib/Npcf_PolicyAuthorization"
	"free5gc/lib/Npcf_SMPolicyControl"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/amf/amf_service"
	"free5gc/src/app"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_service"
	"free5gc/src/udr/udr_service"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

const amPolicyDataColl = "policyData.ues.amData"
const smPolicyDataColl = "policyData.ues.smData"

var NFs = []app.NetworkFunction{
	&nrf_service.NRF{},
	&amf_service.AMF{},
	&udr_service.UDR{},
	&pcf_service.PCF{},
}

var filterUeIdOnly bson.M

func toBsonM(data interface{}) bson.M {
	tmp, _ := json.Marshal(data)
	var putData = bson.M{}
	_ = json.Unmarshal(tmp, &putData)
	return putData
}
func insertDefaultPoliciesToDb(ueId string) {
	amPolicyData := models.AmPolicyData{
		SubscCats: []string{
			"free5gc",
		},
	}

	smPolicyData := models.SmPolicyData{
		SmPolicySnssaiData: map[string]models.SmPolicySnssaiData{
			"01010203": {
				Snssai: &models.Snssai{
					Sd:  "010203",
					Sst: 1,
				},
				SmPolicyDnnData: map[string]models.SmPolicyDnnData{
					"internet": {
						Dnn:        "internet",
						GbrUl:      "500 Mbps",
						GbrDl:      "500 Mbps",
						AdcSupport: false,
						Ipv4Index:  6,
						Ipv6Index:  6,
						Offline:    true,
						Online:     false,
						// ChfInfo
						// RefUmDataLimitIds
						// MpsPriority
						// ImsSignallingPrio
						// MpsPriorityLevel
						// AllowedServices
						// SubscCats
						// SubscSpendingLimit

					},
				},
			},
			"01112233": {
				Snssai: &models.Snssai{
					Sd:  "112233",
					Sst: 1,
				},
				SmPolicyDnnData: map[string]models.SmPolicyDnnData{
					"internet": {
						Dnn: "internet",
					},
				},
			},
		},
	}

	filterUeIdOnly = bson.M{"ueId": ueId}
	amPolicyDataBsonM := toBsonM(amPolicyData)
	amPolicyDataBsonM["ueId"] = ueId
	MongoDBLibrary.RestfulAPIPutOne(amPolicyDataColl, filterUeIdOnly, amPolicyDataBsonM)
	smPolicyDataBsonM := toBsonM(smPolicyData)
	smPolicyDataBsonM["ueId"] = ueId
	MongoDBLibrary.RestfulAPIPost(smPolicyDataColl, filterUeIdOnly, smPolicyDataBsonM)
}

func fakeSmfServer(t *testing.T, port string) {
	go func() { // fake smf server
		router := gin.Default()

		router.POST("nsmf-callback/v1/sm-policies/:smPolicyId/update", func(c *gin.Context) {
			smPolicyId := c.Param("smPolicyId")
			fmt.Println("==========SM Policy Update Notification Callback=============")
			fmt.Println("smPolicyId: ", smPolicyId)

			var notification models.SmPolicyNotification
			if err := c.ShouldBindJSON(&notification); err != nil {
				fmt.Println("fake smf server error")
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			c.JSON(http.StatusNoContent, gin.H{})
		})

		router.POST("nsmf-callback/v1/sm-policies/:smPolicyId/terminate", func(c *gin.Context) {
			smPolicyId := c.Param("smPolicyId")
			fmt.Println("==========SM Policy Terminate Callback=============")
			fmt.Println("smPolicyId: ", smPolicyId)

			var terminationNotification models.TerminationNotification
			if err := c.ShouldBindJSON(&terminationNotification); err != nil {
				fmt.Println("fake smf server error")
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			c.JSON(http.StatusNoContent, gin.H{})
			configuration := Npcf_SMPolicyControl.NewConfiguration()
			configuration.SetBasePath("https://127.0.0.1:29507")
			smclient := Npcf_SMPolicyControl.NewAPIClient(configuration)
			httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdDeletePost(context.Background(), smPolicyId, models.SmPolicyDeleteData{})
			assert.True(t, err == nil)
			assert.True(t, httpRsp != nil)
			if httpRsp != nil {
				assert.Equal(t, http.StatusNoContent, httpRsp.StatusCode)
			}
		})

		smfLogPath := path_util.Gofree5gcPath("free5gc/smfsslkey.log")
		smfPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/smf.pem")
		smfKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/smf.key")

		server, err := http2_util.NewServer(port, smfLogPath, router)
		if err == nil && server != nil {
			fmt.Println(server.ListenAndServeTLS(smfPemPath, smfKeyPath))
		}
		assert.True(t, err == nil)
	}()
}

func init() {
	app.AppInitializeWillInitialize("")
	flag := flag.FlagSet{}
	cli := cli.NewContext(nil, &flag, nil)
	for i, service := range NFs {
		service.Initialize(cli)
		go service.Start()
		time.Sleep(300 * time.Millisecond)
		if i == 0 {
			MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})
			time.Sleep(300 * time.Millisecond)
		}
	}
	insertDefaultPoliciesToDb("imsi-2089300007487")

	time.Sleep(100 * time.Millisecond)

}

func TestApplicationSessionsCollection(t *testing.T) {
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
	{
		afReqData := TestPolicyAuthorization.GetPostAppSessionsData_Normal()
		_, httpRsp, err := afclient.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), afReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			appSessionId := fmt.Sprintf("%s-%d", ue.Supi, ue.AppSessionIdGenerator-1)
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
		afReqData := TestPolicyAuthorization.GetPostAppSessionsData_Flow3()
		_, httpRsp, err := afclient.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), afReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			appSessionId := fmt.Sprintf("%s-%d", ue.Supi, ue.AppSessionIdGenerator-1)
			header := fmt.Sprintf("https://127.0.0.1:29507/npcf-policyauthorization/v1/app-sessions/%s", appSessionId)
			assert.True(t, locationHeader == header)
			appSession := pcfSelf.AppSessionPool[appSessionId]
			if assert.NotNil(t, appSession) {
				assert.True(t, appSession.EventUri == afReqData.AscReqData.EvSubsc.NotifUri)
				assert.True(t, len(appSession.Events) == len(afReqData.AscReqData.EvSubsc.Events))
				smPolicy := ue.SmPolicyData[smPolicyId]
				assert.Equal(t, smPolicy, appSession.SmPolicyData)
				assert.Equal(t, 3, len(appSession.RelatedPccRuleIds))
				assert.Equal(t, 3, len(smPolicy.PolicyDecision.PccRules))
			}
		}
	}
	{
		afReqData := TestPolicyAuthorization.GetPostAppSessionsData_403Forbidden()
		_, httpRsp, err := afclient.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), afReqData)
		assert.True(t, err != nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusForbidden, httpRsp.StatusCode)
			problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			assert.Equal(t, "REQUESTED_SERVICE_NOT_AUTHORIZED", problem.Cause)
		}
	}
	{
		afReqData := TestPolicyAuthorization.GetPostAppSessionsData_400()
		_, httpRsp, err := afclient.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), afReqData)
		assert.True(t, err != nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusBadRequest, httpRsp.StatusCode)
			problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			assert.Equal(t, "ERROR_REQUEST_PARAMETERS", problem.Cause)
		}
	}
}
