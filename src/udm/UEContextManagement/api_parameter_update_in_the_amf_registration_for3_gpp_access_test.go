package UEContextManagement_test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	Nudm_UECM_Client "free5gc/lib/Nudm_UEContextManagement"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	Nudm_UECM_Server "free5gc/src/udm/UEContextManagement"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_context"
	"free5gc/src/udm/udm_handler"
	"free5gc/src/udm/udm_util"
	"net/http"
	"testing"
)

func TestUpdateAmf3gppAccess(t *testing.T) {
	go func() { // udm server
		router := gin.Default()
		Nudm_UECM_Server.AddService(router)

		udmLogPath := path_util.Gofree5gcPath("free5gc/udmsslkey.log")
		udmPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.pem")
		udmKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.key")

		server, err := http2_util.NewServer(":29503", udmLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udmPemPath, udmKeyPath))
			assert.True(t, err == nil)
		}
	}()
	udm_util.testInitUdmConfig()
	go udm_handler.Handle()

	go func() { // fake udr server
		router := gin.Default()

		router.PATCH("/nudr-dr/v1/subscription-data/:ueId/context-data/amf-3gpp-access", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("==========Parameter update in the AMF registration for 3GPP access==========")
			fmt.Println("ueId: ", ueId)
			var patchItems []models.PatchItem
			if err := c.ShouldBindJSON(&patchItems); err != nil {
				fmt.Println("fake udm server error: ", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			fmt.Println("patchItems: ", patchItems)
			c.JSON(http.StatusNoContent, nil)
		})

		router.PUT("/nudr-dr/v1/subscription-data/:ueId/context-data/amf-3gpp-access", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("==========AMF registration for 3GPP access==========")
			fmt.Println("ueId: ", ueId)

			var amf3GppAccessRegistration models.Amf3GppAccessRegistration
			if err := c.ShouldBindJSON(&amf3GppAccessRegistration); err != nil {
				fmt.Println("fake udr server error")
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			fmt.Println("amf3GppAccessRegistration - ", amf3GppAccessRegistration.AmfInstanceId)
			c.JSON(http.StatusNoContent, gin.H{})
		})

		udrLogPath := path_util.Gofree5gcPath("free5gc/udrsslkey.log")
		udrPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udr.pem")
		udrKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udr.key")

		server, err := http2_util.NewServer(":29504", udrLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udrPemPath, udrKeyPath))
			assert.True(t, err == nil)
		}
	}()

	udm_context.Init()
	cfg := Nudm_UECM_Client.NewConfiguration()
	cfg.SetBasePath("https://localhost:29503")
	clientAPI := Nudm_UECM_Client.NewAPIClient(cfg)

	ueId := "UECM1234"

	var putGuami models.Guami
	putGuami.AmfId = "TEST_GUAMI_001"
	putGuami.PlmnId = new(models.PlmnId)
	putGuami.PlmnId.Mcc = "208"
	putGuami.PlmnId.Mnc = "93"
	var amf3GppAccessRegistration models.Amf3GppAccessRegistration
	amf3GppAccessRegistration.AmfInstanceId = "PUT_TEST_001"
	amf3GppAccessRegistration.Guami = &putGuami
	_, putresp, puterr := clientAPI.AMFRegistrationFor3GPPAccessApi.Registration(context.Background(), ueId, amf3GppAccessRegistration)
	if puterr != nil {
		fmt.Println(puterr.Error())
	} else {
		fmt.Println("PUT resp: ", putresp)
	}

	var patchGuami models.Guami
	patchGuami.AmfId = "TEST_GUAMI_001"
	patchGuami.PlmnId = new(models.PlmnId)
	patchGuami.PlmnId.Mcc = "208"
	patchGuami.PlmnId.Mnc = "93"
	var amf3GppAccessRegistrationModification models.Amf3GppAccessRegistrationModification
	amf3GppAccessRegistrationModification.Pei = "testPEI"
	amf3GppAccessRegistrationModification.Guami = &patchGuami
	patchresp, patcherr := clientAPI.ParameterUpdateInTheAMFRegistrationFor3GPPAccessApi.Update(context.Background(), ueId, amf3GppAccessRegistrationModification)
	if patcherr != nil {
		fmt.Println(patcherr.Error())
	} else {
		fmt.Println("PATCH resp: ", patchresp)
	}
}
