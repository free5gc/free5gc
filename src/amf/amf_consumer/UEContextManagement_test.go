package amf_consumer_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/logger"
	// "free5gc/src/udm/udm_context"
	"net/http"
	"testing"
	"time"
)

func TestUeCmRegistration(t *testing.T) {

	go func() { // fake udr server
		router := gin.Default()

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
	udminit()

	// udm_context.Init()

	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	// udmUri := "https://localhost:29503"
	problemDetails, err := amf_consumer.UeCmRegistration(ue, models.AccessType__3_GPP_ACCESS, true)
	if err != nil {
		fmt.Println(err.Error())
	} else if problemDetails != nil {
		fmt.Println("problemDetails: ", problemDetails)
	}
}
