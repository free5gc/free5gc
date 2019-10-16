package UEContextManagement_test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	Nudm_UECM_Client "free5gc/lib/Nudm_UEContextManagement"
	"free5gc/lib/http2_util"
	// "free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	Nudm_UECM_Server "free5gc/src/udm/UEContextManagement"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_context"
	"free5gc/src/udm/udm_handler"
	"free5gc/src/udm/udm_util"
	"net/http"
	"testing"
)

// GetAmf3gppAccess - retrieve the AMF registration for 3GPP access information
func TestDeregistrationSmfRegistrations(t *testing.T) {
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

		router.DELETE("/nudr-dr/v1/subscription-data/:ueId/context-data/smf-registrations/:pduSessionId", func(c *gin.Context) {
			ueId := c.Param("ueId")
			pduSessionId := c.Param("pduSessionId")
			fmt.Println("==========SMF Deregistration==========")
			fmt.Println("ueId: ", ueId)
			fmt.Println("pduSessionId: ", pduSessionId)

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
	pduSessionIdInt := 222
	pduSessionId := int32(pduSessionIdInt)
	resp, err := clientAPI.SMFDeregistrationApi.Deregistration(context.Background(), ueId, pduSessionId)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("resp: ", resp)
	}
}
