package UEAuthentication

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/Nudm_UEAuthentication"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_handler"
	"free5gc/src/udm/udm_util"
	"net/http"
	"testing"
)

func TestUeAuthenticationsPost(t *testing.T) {
	go func() { // udm server
		router := gin.Default()
		AddService(router)

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

		router.GET("/nudr-dr/v1/subscription-data/:ueId/authentication-data/authentication-subscription", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("ueId: ", ueId)
			var authSubs models.AuthenticationSubscription
			var pk models.PermanentKey
			var opc models.Opc
			var var_milenage models.Milenage
			var op models.Op

			pk.PermanentKeyValue = TestGenAuthData.MilenageTestSet19.K
			opc.OpcValue = TestGenAuthData.MilenageTestSet19.OPC
			op.OpValue = TestGenAuthData.MilenageTestSet19.OP
			var_milenage.Op = &op

			authSubs.PermanentKey = &pk
			authSubs.Opc = &opc
			authSubs.Milenage = &var_milenage
			authSubs.SequenceNumber = TestGenAuthData.MilenageTestSet19.SQN
			authSubs.AuthenticationMethod = models.AuthMethod__5_G_AKA
			// authSubs.AuthenticationMethod = models.AuthMethod_EAP_AKA_PRIME

			c.JSON(http.StatusOK, authSubs)
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

	var authInfoReq models.AuthenticationInfoRequest
	authInfoReq.ServingNetworkName = TestGenAuthData.TestGenAuthDataTable[TestGenAuthData.SUCCESS_CASE].ServingNetworkName
	supiOrSuci := TestGenAuthData.SUPI

	cfg := Nudm_UEAuthentication.NewConfiguration()
	cfg.SetBasePath("https://localhost:29503")
	client := Nudm_UEAuthentication.NewAPIClient(cfg)

	authInfoRes, resp, err := client.GenerateAuthDataApi.GenerateAuthData(context.TODO(), supiOrSuci, authInfoReq)
	fmt.Println("=====")
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println("resp: ", resp)
		fmt.Println("authInfoRes: ", authInfoRes)
	}

	switch authInfoRes.AuthType {
	case models.AuthType__5_G_AKA:
		fmt.Printf("auth type: 5G AKA\n")
		// rand, xresStar, autn, kausf
		av := authInfoRes.AuthenticationVector
		if av != nil {
			fmt.Printf("rand: %s\nxresStar: %s\nautn: %s\nkausf: %s\n", av.Rand, av.XresStar, av.Autn, av.Kausf)
		} else {
			fmt.Printf("nil av\n")
		}

	case models.AuthType_EAP_AKA_PRIME:
		fmt.Printf("auth type: EAP-AKA'\n")
		// rand, xres, autn, ckPrime, ikPrime
		av := authInfoRes.AuthenticationVector
		if av != nil {
			fmt.Printf("rand: %s\nxres: %s\nautn: %s\nCK': %s\nIK': %s", av.Rand, av.Xres, av.Autn, av.CkPrime, av.IkPrime)
		} else {
			fmt.Printf("nil av\n")
		}

	default:
		fmt.Println("authInfoRes authType error")
	}
}
