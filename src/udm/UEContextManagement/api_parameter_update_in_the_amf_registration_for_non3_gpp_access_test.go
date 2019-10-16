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

func TestUpdateAmfNon3gppAccess(t *testing.T) {
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

		router.PATCH("/nudr-dr/v1/subscription-data/:ueId/context-data/amf-non-3gpp-access", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("==========Parameter update in the AMF registration for non-3GPP access==========")
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

		router.PUT("/nudr-dr/v1/subscription-data/:ueId/context-data/amf-non-3gpp-access", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("==========AMF registration for non-3GPP access==========")
			fmt.Println("ueId: ", ueId)

			var amfNon3GppAccessRegistration models.AmfNon3GppAccessRegistration
			if err := c.ShouldBindJSON(&amfNon3GppAccessRegistration); err != nil {
				fmt.Println("fake udr server error")
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			fmt.Println("amfNon3GppAccessRegistration - ", amfNon3GppAccessRegistration.AmfInstanceId)
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
	putGuami.AmfId = "NON_3GPP_TEST_GUAMI_001"
	putGuami.PlmnId = new(models.PlmnId)
	putGuami.PlmnId.Mcc = "208"
	putGuami.PlmnId.Mnc = "93"
	var amfNon3GppAccessRegistration models.AmfNon3GppAccessRegistration
	amfNon3GppAccessRegistration.AmfInstanceId = "NON_3GPP_PUT_TEST_001"
	amfNon3GppAccessRegistration.Guami = &putGuami
	_, putresp, puterr := clientAPI.AMFRegistrationForNon3GPPAccessApi.Register(context.Background(), ueId, amfNon3GppAccessRegistration)
	if puterr != nil {
		fmt.Println(puterr.Error())
	} else {
		fmt.Println("PUT resp: ", putresp)
	}

	var patchGuami models.Guami
	patchGuami.AmfId = "NON_3GPP_TEST_GUAMI_001"
	patchGuami.PlmnId = new(models.PlmnId)
	patchGuami.PlmnId.Mcc = "208"
	patchGuami.PlmnId.Mnc = "93"
	var amfNon3GppAccessRegistrationModification models.AmfNon3GppAccessRegistrationModification
	amfNon3GppAccessRegistrationModification.Pei = "NON_3GPP_testPEI"
	amfNon3GppAccessRegistrationModification.Guami = &patchGuami
	patchresp, patcherr := clientAPI.ParameterUpdateInTheAMFRegistrationForNon3GPPAccessApi.UpdateAmfNon3gppAccess(context.Background(), ueId, amfNon3GppAccessRegistrationModification)
	if patcherr != nil {
		fmt.Println(patcherr.Error())
	} else {
		fmt.Println("PATCH resp: ", patchresp)
	}
}
