package SMPolicy_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/CommonConsumerTestData/PCF/TestAMPolicy"
	"free5gc/lib/CommonConsumerTestData/PCF/TestSMPolicy"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Npcf_AMPolicy"
	"free5gc/lib/Npcf_SMPolicyControl"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/amf/amf_service"
	"free5gc/src/app"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_producer"
	"free5gc/src/pcf/pcf_service"
	"free5gc/src/udr/DataRepository"
	"free5gc/src/udr/factory"
	"free5gc/src/udr/udr_consumer"
	"free5gc/src/udr/udr_service"
	"free5gc/src/udr/udr_util"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

var NFs = []app.NetworkFunction{
	&nrf_service.NRF{},
	&amf_service.AMF{},
	&udr_service.UDR{},
	&pcf_service.PCF{},
}

func fakeudrInit() {
	config := factory.UdrConfig
	sbi := config.Configuration.Sbi
	mongodb := config.Configuration.Mongodb
	nrfUri := config.Configuration.NrfUri

	// Connect to MongoDB
	DataRepository.SetMongoDB(mongodb.Name, mongodb.Url)

	udrLogPath := udr_util.UdrLogPath
	udrPemPath := udr_util.UdrPemPath
	udrKeyPath := udr_util.UdrKeyPath
	if sbi.Tls != nil {
		udrLogPath = path_util.Gofree5gcPath(sbi.Tls.Log)
		udrPemPath = path_util.Gofree5gcPath(sbi.Tls.Pem)
		udrKeyPath = path_util.Gofree5gcPath(sbi.Tls.Key)
	}

	profile := udr_consumer.BuildNFInstance()
	newNrfUri, err := udr_consumer.SendRegisterNFInstance(nrfUri, profile.NfInstanceId, profile)
	if err == nil {
		config.Configuration.NrfUri = newNrfUri
	} else {
		fmt.Errorf("Send Register NFInstance Error[%s]", err.Error())
	}
	go func() { // fake udr server
		router := gin.Default()

		router.GET("/nudr-dr/v1/policy-data/ues/:ueId/sm-data", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("==========GET SM Policy Data==========")
			fmt.Println("ueId: ", ueId)
			queryParameters := c.Request.URL.Query()
			snssai := &models.Snssai{}
			dnn := ""
			if queryParameters["dnn"] != nil {
				dnn = queryParameters["dnn"][0]
				fmt.Println("dnn : ", dnn)
			}
			if queryParameters["snssai"] != nil {
				tmp := queryParameters["snssai"][0]
				err := json.Unmarshal([]byte(tmp), snssai)
				if err != nil {
					fmt.Errorf(err.Error())
				} else {
					fmt.Printf("snssai : %+v\n", snssai)
				}
			}
			key := fmt.Sprintf("%02x%s", snssai.Sst, snssai.Sd)
			if len(key) != 8 {
				c.JSON(http.StatusBadRequest, gin.H{})
				return
			}
			rsp := models.SmPolicyData{
				SmPolicySnssaiData: make(map[string]models.SmPolicySnssaiData),
			}

			rsp.SmPolicySnssaiData[key] = models.SmPolicySnssaiData{
				Snssai:          snssai,
				SmPolicyDnnData: make(map[string]models.SmPolicyDnnData),
			}
			rsp.SmPolicySnssaiData[key].SmPolicyDnnData[dnn] = models.SmPolicyDnnData{
				Dnn: dnn,
				// AllowedServices
				// SubscCats
				GbrUl:      "500Mbps",
				GbrDl:      "500Mbps",
				AdcSupport: false,
				// SubscSpendingLimit
				Ipv4Index: 6,
				Ipv6Index: 6,
				Offline:   false,
				Online:    false,
				// ChfInfo
				// RefUmDataLimitIds
				// MpsPriority
				// ImsSignallingPrio
				// MpsPriorityLevel
			}
			c.JSON(http.StatusOK, rsp)
		})
		server, err := http2_util.NewServer(":29504", udrLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udrPemPath, udrKeyPath))
		}
	}()
}

func init() {
	app.AppInitializeWillInitialize("")
	flag := flag.FlagSet{}
	cli := cli.NewContext(nil, &flag, nil)
	for i, service := range NFs {
		if i == 2 {
			service.Initialize(cli)
			fakeudrInit()
		} else {
			service.Initialize(cli)
			go service.Start()
		}
		time.Sleep(300 * time.Millisecond)
		if i == 0 {
			MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})
			time.Sleep(300 * time.Millisecond)
		}
	}
}
func TestCreateSMPolicy(t *testing.T) {
	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	amclient := Npcf_AMPolicy.NewAPIClient(configuration)
	configuration1 := Npcf_SMPolicyControl.NewConfiguration()
	configuration1.SetBasePath("https://127.0.0.1:29507")
	smclient := Npcf_SMPolicyControl.NewAPIClient(configuration1)
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
	{
		smCreateReqData := TestSMPolicy.CreateTestData()
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesPost(context.Background(), smCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			// index := strings.LastIndex(locationHeader, "/")
			assert.True(t, locationHeader == "https://127.0.0.1:29507/npcf-smpolicycontrol/v1/sm-policies/imsi-2089300007487-1")
		}
	}
	{
		smCreateReqData := TestSMPolicy.CreateTestData()
		smCreateReqData.Supi = ""
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesPost(context.Background(), smCreateReqData)
		assert.True(t, err != nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusBadRequest, httpRsp.StatusCode)
			if httpRsp != nil {
				assert.Equal(t, http.StatusBadRequest, httpRsp.StatusCode)
				problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
				assert.Equal(t, "ERROR_INITIAL_PARAMETERS", problem.Cause)
				// assert.Equal(t, "Supi Format Error", problem.Detail)
			}
		}
	}
	{
		smCreateReqData := TestSMPolicy.CreateTestData()
		smCreateReqData.SliceInfo.Sd = "123"
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesPost(context.Background(), smCreateReqData)
		assert.True(t, err != nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusBadRequest, httpRsp.StatusCode)
			if httpRsp != nil {
				assert.Equal(t, http.StatusBadRequest, httpRsp.StatusCode)
				problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
				assert.Equal(t, "ERROR_INITIAL_PARAMETERS", problem.Cause)
				// assert.Equal(t, "Supi Format Error", problem.Detail)
			}
		}
	}
}

func TestGetSMPolicy(t *testing.T) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	amclient := Npcf_AMPolicy.NewAPIClient(configuration)
	configuration1 := Npcf_SMPolicyControl.NewConfiguration()
	configuration1.SetBasePath("https://127.0.0.1:29507")
	smclient := Npcf_SMPolicyControl.NewAPIClient(configuration1)
	smPolicyId := "imsi-2089300007487-1"
	smCreateReqData := TestSMPolicy.CreateTestData()
	var decision models.SmPolicyDecision
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
	{
		tmp, httpRsp, err := smclient.DefaultApi.SmPoliciesPost(context.Background(), smCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			// index := strings.LastIndex(locationHeader, "/")
			assert.True(t, locationHeader == "https://127.0.0.1:29507/npcf-smpolicycontrol/v1/sm-policies/imsi-2089300007487-1")
		}
		decision = tmp
	}
	{
		//Test GetPoliciesPolAssoId
		rsp, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
			assert.True(t, reflect.DeepEqual(smCreateReqData, *rsp.Context))
			assert.True(t, reflect.DeepEqual(decision, *rsp.Policy))
		}
	}

}

func TestDelSMPolicy(t *testing.T) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	amclient := Npcf_AMPolicy.NewAPIClient(configuration)
	configuration1 := Npcf_SMPolicyControl.NewConfiguration()
	configuration1.SetBasePath("https://127.0.0.1:29507")
	smclient := Npcf_SMPolicyControl.NewAPIClient(configuration1)
	smPolicyId := "imsi-2089300007487-1"
	smCreateReqData := TestSMPolicy.CreateTestData()
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
	{
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesPost(context.Background(), smCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			// index := strings.LastIndex(locationHeader, "/")
			assert.True(t, locationHeader == "https://127.0.0.1:29507/npcf-smpolicycontrol/v1/sm-policies/imsi-2089300007487-1")
		}
	}
	{
		//Test DelPoliciesPolAssoId
		httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdDeletePost(context.Background(), smPolicyId, models.SmPolicyDeleteData{})
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusNoContent, httpRsp.StatusCode)
		}
	}
	{
		//Test GetPoliciesPolAssoId
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
		assert.True(t, err != nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusNotFound, httpRsp.StatusCode)
			problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			assert.Equal(t, "CONTEXT_NOT_FOUND", problem.Cause)
		}
	}
}

func TestUpdateAMPolicy(t *testing.T) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	amclient := Npcf_AMPolicy.NewAPIClient(configuration)
	configuration1 := Npcf_SMPolicyControl.NewConfiguration()
	configuration1.SetBasePath("https://127.0.0.1:29507")
	smclient := Npcf_SMPolicyControl.NewAPIClient(configuration1)
	smPolicyId := "imsi-2089300007487-1"
	smCreateReqData := TestSMPolicy.CreateTestData()
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
	{
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesPost(context.Background(), smCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			// index := strings.LastIndex(locationHeader, "/")
			assert.True(t, locationHeader == "https://127.0.0.1:29507/npcf-smpolicycontrol/v1/sm-policies/imsi-2089300007487-1")
		}
	}
	{
		trigger := []models.PolicyControlRequestTrigger{
			models.PolicyControlRequestTrigger_PLMN_CH,
			models.PolicyControlRequestTrigger_AC_TY_CH,
			models.PolicyControlRequestTrigger_UE_IP_CH,
			models.PolicyControlRequestTrigger_PS_DA_OFF,
			models.PolicyControlRequestTrigger_DEF_QOS_CH,
			models.PolicyControlRequestTrigger_SE_AMBR_CH,
			models.PolicyControlRequestTrigger_SAREA_CH,
			models.PolicyControlRequestTrigger_SCNN_CH,
			models.PolicyControlRequestTrigger_RAT_TY_CH,
			models.PolicyControlRequestTrigger_UE_TZ_CH,
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
			// 	{
			// 		//Test GetPoliciesPolAssoId
			// 		rsp, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
			// 		assert.True(t, err == nil)
			// 		assert.True(t, httpRsp != nil)
			// 		if httpRsp != nil {
			// 			assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
			// 		}
			// 	}
			// }
		}
	}
	{
		trigger := []models.PolicyControlRequestTrigger{
			models.PolicyControlRequestTrigger_AC_TY_CH,
			models.PolicyControlRequestTrigger_RAT_TY_CH,
			models.PolicyControlRequestTrigger_RES_MO_RE,
		}
		op := models.RuleOperation_CREATE_PCC_RULE
		updateReq := TestSMPolicy.UpdateTestData(trigger, &op)
		//Test UpdatePoliciesPolAssoId
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdUpdatePost(context.Background(), smPolicyId, updateReq)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			if assert.Equal(t, http.StatusOK, httpRsp.StatusCode) {
				{
					//Test GetPoliciesPolAssoId
					rsp, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
					assert.True(t, err == nil)
					assert.True(t, httpRsp != nil)
					if httpRsp != nil {
						assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
						rule, exist := rsp.Policy.PccRules[updateReq.UeInitResReq.PccRuleId]
						assert.True(t, exist)
						req := updateReq.UeInitResReq
						assert.Equal(t, rule.FlowInfos[0].FlowDescription, req.PackFiltInfo[0].PackFiltCont)
						assert.True(t, models.FlowDirection(rule.FlowInfos[0].FlowDirection) == req.PackFiltInfo[0].FlowDirection)
						assert.Equal(t, req.ReqQos.Var5qi, rsp.Policy.QosDecs[rule.RefQosData[0]].Var5qi)
						assert.Equal(t, req.ReqQos.GbrUl, rsp.Policy.QosDecs[rule.RefQosData[0]].GbrUl)
						assert.Equal(t, req.ReqQos.GbrDl, rsp.Policy.QosDecs[rule.RefQosData[0]].GbrDl)
					}
				}
			}
		}
	}
	{
		trigger := []models.PolicyControlRequestTrigger{
			models.PolicyControlRequestTrigger_RES_MO_RE,
		}
		op := models.RuleOperation_MODIFY_PCC_RULE_AND_ADD_PACKET_FILTERS
		updateReq := TestSMPolicy.UpdateTestData(trigger, &op)
		//Test UpdatePoliciesPolAssoId
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdUpdatePost(context.Background(), smPolicyId, updateReq)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			if assert.Equal(t, http.StatusOK, httpRsp.StatusCode) {
				{
					//Test GetPoliciesPolAssoId
					rsp, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
					assert.True(t, err == nil)
					assert.True(t, httpRsp != nil)
					if httpRsp != nil {
						assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
						rule, exist := rsp.Policy.PccRules[updateReq.UeInitResReq.PccRuleId]
						assert.True(t, exist)
						req := updateReq.UeInitResReq
						assert.Equal(t, rule.FlowInfos[1].FlowDescription, req.PackFiltInfo[0].PackFiltCont)
						assert.True(t, models.FlowDirection(rule.FlowInfos[1].FlowDirection) == req.PackFiltInfo[0].FlowDirection)
						assert.Equal(t, req.ReqQos.Var5qi, rsp.Policy.QosDecs[rule.RefQosData[0]].Var5qi)
						assert.Equal(t, req.ReqQos.GbrUl, rsp.Policy.QosDecs[rule.RefQosData[0]].GbrUl)
					}
				}
			}
		}
	}
	{
		trigger := []models.PolicyControlRequestTrigger{
			models.PolicyControlRequestTrigger_RES_MO_RE,
		}
		op := models.RuleOperation_MODIFY_PCC_RULE_AND_DELETE_PACKET_FILTERS
		updateReq := TestSMPolicy.UpdateTestData(trigger, &op)
		//Test UpdatePoliciesPolAssoId
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdUpdatePost(context.Background(), smPolicyId, updateReq)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			if assert.Equal(t, http.StatusOK, httpRsp.StatusCode) {
				{
					//Test GetPoliciesPolAssoId
					rsp, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
					assert.True(t, err == nil)
					assert.True(t, httpRsp != nil)
					if httpRsp != nil {
						assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
						rule, exist := rsp.Policy.PccRules[updateReq.UeInitResReq.PccRuleId]
						assert.True(t, exist)
						assert.True(t, len(rule.FlowInfos) == 1)
					}
				}
			}
		}
	}
	{
		trigger := []models.PolicyControlRequestTrigger{
			models.PolicyControlRequestTrigger_RES_MO_RE,
		}
		op := models.RuleOperation_MODIFY_PCC_RULE_AND_REPLACE_PACKET_FILTERS
		updateReq := TestSMPolicy.UpdateTestData(trigger, &op)
		//Test UpdatePoliciesPolAssoId
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdUpdatePost(context.Background(), smPolicyId, updateReq)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			if assert.Equal(t, http.StatusOK, httpRsp.StatusCode) {
				{
					//Test GetPoliciesPolAssoId
					rsp, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
					assert.True(t, err == nil)
					assert.True(t, httpRsp != nil)
					if httpRsp != nil {
						assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
						rule, exist := rsp.Policy.PccRules[updateReq.UeInitResReq.PccRuleId]
						assert.True(t, exist)
						assert.True(t, len(rule.FlowInfos) == 1)
						req := updateReq.UeInitResReq
						assert.Equal(t, rule.FlowInfos[0].FlowDescription, req.PackFiltInfo[0].PackFiltCont)
						assert.True(t, models.FlowDirection(rule.FlowInfos[0].FlowDirection) == req.PackFiltInfo[0].FlowDirection)
					}
				}
			}
		}
	}
	{
		trigger := []models.PolicyControlRequestTrigger{
			models.PolicyControlRequestTrigger_RES_MO_RE,
		}
		op := models.RuleOperation_MODIFY_PCC_RULE_WITHOUT_MODIFY_PACKET_FILTERS
		updateReq := TestSMPolicy.UpdateTestData(trigger, &op)
		//Test UpdatePoliciesPolAssoId
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdUpdatePost(context.Background(), smPolicyId, updateReq)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			if assert.Equal(t, http.StatusOK, httpRsp.StatusCode) {
				{
					//Test GetPoliciesPolAssoId
					rsp, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
					assert.True(t, err == nil)
					assert.True(t, httpRsp != nil)
					if httpRsp != nil {
						assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
						rule, exist := rsp.Policy.PccRules[updateReq.UeInitResReq.PccRuleId]
						assert.True(t, exist)
						req := updateReq.UeInitResReq
						assert.Equal(t, req.ReqQos.Var5qi, rsp.Policy.QosDecs[rule.RefQosData[0]].Var5qi)
						assert.Equal(t, req.ReqQos.GbrUl, rsp.Policy.QosDecs[rule.RefQosData[0]].GbrUl)
						assert.Equal(t, req.ReqQos.GbrDl, rsp.Policy.QosDecs[rule.RefQosData[0]].GbrDl)
					}
				}
			}
		}
	}
	{
		trigger := []models.PolicyControlRequestTrigger{
			models.PolicyControlRequestTrigger_RES_MO_RE,
		}
		op := models.RuleOperation_DELETE_PCC_RULE
		updateReq := TestSMPolicy.UpdateTestData(trigger, &op)
		//Test UpdatePoliciesPolAssoId
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdUpdatePost(context.Background(), smPolicyId, updateReq)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			if assert.Equal(t, http.StatusOK, httpRsp.StatusCode) {
				{
					//Test GetPoliciesPolAssoId
					rsp, httpRsp, err := smclient.DefaultApi.SmPoliciesSmPolicyIdGet(context.Background(), smPolicyId)
					assert.True(t, err == nil)
					assert.True(t, httpRsp != nil)
					if httpRsp != nil {
						assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
						assert.True(t, len(rsp.Policy.PccRules) == 0)
						assert.True(t, len(rsp.Policy.QosDecs) == 0)
					}
				}
			}
		}
	}
}

func TestSMPolicyNotification(t *testing.T) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	amclient := Npcf_AMPolicy.NewAPIClient(configuration)
	configuration1 := Npcf_SMPolicyControl.NewConfiguration()
	configuration1.SetBasePath("https://127.0.0.1:29507")
	smclient := Npcf_SMPolicyControl.NewAPIClient(configuration1)
	go func() { // fake udr server
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

		server, err := http2_util.NewServer(":8888", smfLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(smfPemPath, smfKeyPath))
		}
		assert.True(t, err == nil)
	}()

	time.Sleep(100 * time.Millisecond)

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
	smPolicyId := "imsi-2089300007487-1"
	smCreateReqData := TestSMPolicy.CreateTestData()
	smCreateReqData.NotificationUri = "https://127.0.0.1:8888/nsmf-callback/v1/sm-policies/imsi-2089300007487-1"
	{
		_, httpRsp, err := smclient.DefaultApi.SmPoliciesPost(context.Background(), smCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			// index := strings.LastIndex(locationHeader, "/")
			assert.True(t, locationHeader == "https://127.0.0.1:29507/npcf-smpolicycontrol/v1/sm-policies/imsi-2089300007487-1")
		}
	}

	ue := pcf_context.PCF_Self().UePool["imsi-2089300007487"]
	//Test Policies Update Notify
	notification := models.SmPolicyNotification{
		ResourceUri: smCreateReqData.NotificationUri,
	}
	pcf_producer.SendSMPolicyUpdateNotification(ue, smPolicyId, notification)

	//Test Policies Termination Notify
	termination := models.TerminationNotification{
		ResourceUri: smCreateReqData.NotificationUri,
		Cause:       models.PolicyAssociationReleaseCause_UNSPECIFIED,
	}
	pcf_producer.SendSMPolicyTerminationRequestNotification(ue, smPolicyId, termination)

	time.Sleep(200 * time.Millisecond)
}
