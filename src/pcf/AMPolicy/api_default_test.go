package AMPolicy_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"free5gc/lib/CommonConsumerTestData/PCF/TestAMPolicy"
	"free5gc/lib/Npcf_AMPolicy"
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
	"free5gc/src/udr/udr_service"
	"net/http"
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

func init() {
	app.AppInitializeWillInitialize("")
	flag := flag.FlagSet{}
	cli := cli.NewContext(nil, &flag, nil)
	for _, service := range NFs {
		service.Initialize(cli)
		go service.Start()
		time.Sleep(300 * time.Millisecond)
	}
}
func TestCreateAMPolicy(t *testing.T) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	client := Npcf_AMPolicy.NewAPIClient(configuration)

	//Test PostPolicies
	{
		amCreateReqData := TestAMPolicy.GetAMreqdata()
		_, httpRsp, err := client.DefaultApi.PoliciesPost(context.Background(), amCreateReqData)
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
		amCreateReqData := TestAMPolicy.GetamCreatefailnotifyURIData()
		_, httpRsp, err := client.DefaultApi.PoliciesPost(context.Background(), amCreateReqData)
		assert.True(t, err != nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusBadRequest, httpRsp.StatusCode)
			problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			assert.Equal(t, "ERROR_REQUEST_PARAMETERS", problem.Cause)
			assert.Equal(t, "Miss Mandotory IE", problem.Detail)
		}
	}
	{
		amCreateReqData := TestAMPolicy.GetamCreatefailsupiData()
		_, httpRsp, err := client.DefaultApi.PoliciesPost(context.Background(), amCreateReqData)
		assert.True(t, err != nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusBadRequest, httpRsp.StatusCode)
			problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			assert.Equal(t, "ERROR_REQUEST_PARAMETERS", problem.Cause)
			assert.Equal(t, "Supi Format Error", problem.Detail)
		}
	}

}

func TestGetAMPolicy(t *testing.T) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	client := Npcf_AMPolicy.NewAPIClient(configuration)

	amCreateReqData := TestAMPolicy.GetAMreqdata()
	polAssoId := "imsi-2089300007487-1"
	//Test PostPolicies
	{
		_, httpRsp, err := client.DefaultApi.PoliciesPost(context.Background(), amCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			index := strings.LastIndex(locationHeader, "/")
			assert.True(t, index != -1)
			polAssoId = locationHeader[index+1:]
			assert.True(t, strings.HasPrefix(polAssoId, "imsi-2089300007487"))
		}
	}
	{
		//Test GetPoliciesPolAssoId
		rsp, httpRsp, err := client.DefaultApi.PoliciesPolAssoIdGet(context.Background(), polAssoId)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
			assert.Equal(t, amCreateReqData.SuppFeat, rsp.SuppFeat)
			assert.Equal(t, amCreateReqData.ServAreaRes, rsp.ServAreaRes)
			assert.Equal(t, amCreateReqData.Rfsp, rsp.Rfsp)
			assert.True(t, rsp.Triggers == nil)
			assert.True(t, rsp.Pras == nil)
		}
	}

}

func TestDelAMPolicy(t *testing.T) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	client := Npcf_AMPolicy.NewAPIClient(configuration)

	amCreateReqData := TestAMPolicy.GetAMreqdata()
	polAssoId := "imsi-2089300007487-1"
	//Test PostPolicies
	{
		_, httpRsp, err := client.DefaultApi.PoliciesPost(context.Background(), amCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			index := strings.LastIndex(locationHeader, "/")
			assert.True(t, index != -1)
			polAssoId = locationHeader[index+1:]
			assert.True(t, strings.HasPrefix(polAssoId, "imsi-2089300007487"))
		}
	}
	{
		//Test DelPoliciesPolAssoId
		httpRsp, err := client.DefaultApi.PoliciesPolAssoIdDelete(context.Background(), polAssoId)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusNoContent, httpRsp.StatusCode)
		}
	}
	{
		//Test GetPoliciesPolAssoId
		_, httpRsp, err := client.DefaultApi.PoliciesPolAssoIdGet(context.Background(), polAssoId)
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
	client := Npcf_AMPolicy.NewAPIClient(configuration)

	amCreateReqData := TestAMPolicy.GetAMreqdata()
	polAssoId := "imsi-2089300007487-1"
	//Test PostPolicies
	{
		_, httpRsp, err := client.DefaultApi.PoliciesPost(context.Background(), amCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			index := strings.LastIndex(locationHeader, "/")
			assert.True(t, index != -1)
			polAssoId = locationHeader[index+1:]
			assert.True(t, strings.HasPrefix(polAssoId, "imsi-2089300007487"))
		}
	}
	updateReq := TestAMPolicy.GetAMUpdateReqData()
	{
		//Test UpdatePoliciesPolAssoId
		rsp, httpRsp, err := client.DefaultApi.PoliciesPolAssoIdUpdatePost(context.Background(), polAssoId, updateReq)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusOK, httpRsp.StatusCode)
			assert.Equal(t, updateReq.ServAreaRes, rsp.ServAreaRes)
			assert.Equal(t, updateReq.Rfsp, rsp.Rfsp)
			assert.True(t, rsp.Triggers == nil)
			assert.True(t, rsp.Pras == nil)
		}
	}
}

func TestAMPolicyNotification(t *testing.T) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	client := Npcf_AMPolicy.NewAPIClient(configuration)
	go func() { // fake udr server
		router := gin.Default()

		router.POST("/namf-callback/v1/am-policy/:polAssoId/update", func(c *gin.Context) {
			polAssoId := c.Param("polAssoId")
			fmt.Println("==========AMF Policy Association Update Callback=============")
			fmt.Println("polAssoId: ", polAssoId)

			var policyUpdate models.PolicyUpdate
			if err := c.ShouldBindJSON(&policyUpdate); err != nil {
				fmt.Println("fake amf server error")
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			c.JSON(http.StatusNoContent, gin.H{})
		})

		router.POST("/namf-callback/v1/am-policy/:polAssoId/terminate", func(c *gin.Context) {
			polAssoId := c.Param("polAssoId")
			fmt.Println("==========AMF Policy Association Terminate Callback=============")
			fmt.Println("polAssoId: ", polAssoId)

			var terminationNotification models.TerminationNotification
			if err := c.ShouldBindJSON(&terminationNotification); err != nil {
				fmt.Println("fake amf server error")
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			c.JSON(http.StatusNoContent, gin.H{})
			httpRsp, err := client.DefaultApi.PoliciesPolAssoIdDelete(context.Background(), polAssoId)
			assert.True(t, err == nil)
			assert.True(t, httpRsp != nil)
			if httpRsp != nil {
				assert.Equal(t, http.StatusNoContent, httpRsp.StatusCode)
			}
		})

		amfLogPath := path_util.Gofree5gcPath("free5gc/amfsslkey.log")
		amfPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/amf.pem")
		amfKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/amf.key")

		server, err := http2_util.NewServer(":8888", amfLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(amfPemPath, amfKeyPath))
		}
		assert.True(t, err == nil)
	}()

	time.Sleep(100 * time.Millisecond)

	var polAssoId string
	amCreateReqData := TestAMPolicy.GetAMreqdata()
	amCreateReqData.NotificationUri = "https://127.0.0.1:8888/namf-callback/v1/am-policy/imsi-2089300007487-1"
	//Test PostPolicies
	{
		_, httpRsp, err := client.DefaultApi.PoliciesPost(context.Background(), amCreateReqData)
		assert.True(t, err == nil)
		assert.True(t, httpRsp != nil)
		if httpRsp != nil {
			assert.Equal(t, http.StatusCreated, httpRsp.StatusCode)
			locationHeader := httpRsp.Header.Get("Location")
			index := strings.LastIndex(locationHeader, "/")
			assert.True(t, index != -1)
			polAssoId = locationHeader[index+1:]
			assert.True(t, strings.HasPrefix(polAssoId, "imsi-2089300007487"))
		}
	}
	ue := pcf_context.PCF_Self().UePool["imsi-2089300007487"]
	//Test Policies Update Notify
	policyUpdate := models.PolicyUpdate{
		ResourceUri: amCreateReqData.NotificationUri,
	}
	pcf_producer.SendAMPolicyUpdateNotification(ue, polAssoId, policyUpdate)

	//Test Policies Termination Notify
	notification := models.TerminationNotification{
		ResourceUri: amCreateReqData.NotificationUri,
		Cause:       models.PolicyAssociationReleaseCause_UNSPECIFIED,
	}
	pcf_producer.SendAMPolicyTerminationRequestNotification(ue, polAssoId, notification)

	time.Sleep(200 * time.Millisecond)
}
