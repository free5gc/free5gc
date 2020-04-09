package BDTPolicy_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	TestBDTPolicy "free5gc/lib/CommonConsumerTestData/PCF/TestBDTPolicy"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Npcf_BDTPolicyControl"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/pcf/pcf_service"
	"free5gc/src/udr/DataRepository"
	"free5gc/src/udr/factory"
	"free5gc/src/udr/logger"
	"free5gc/src/udr/udr_consumer"
	"free5gc/src/udr/udr_service"
	"free5gc/src/udr/udr_util"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

var NFs = []app.NetworkFunction{
	&nrf_service.NRF{},
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
	newNrfUri, _, err := udr_consumer.SendRegisterNFInstance(nrfUri, profile.NfInstanceId, profile)
	if err == nil {
		config.Configuration.NrfUri = newNrfUri
	} else {
		fmt.Errorf("Send Register NFInstance Error[%s]", err.Error())
	}
	go func() { // fake udr server
		router := gin.Default()

		router.GET("/nudr-dr/v1/policy-data/bdt-data", func(c *gin.Context) {
			fmt.Println("==========GET BDT Policy Data==========")
			rsp := []models.BdtData{}
			c.JSON(http.StatusOK, rsp)
		})
		router.PUT("/nudr-dr/v1/policy-data/bdt-data/:bdtReferenceId", func(c *gin.Context) {
			bdtReferenceId := c.Param("bdtReferenceId")
			fmt.Println("==========PUT BDT Policy Data==========")
			fmt.Println("bdtReferenceId: ", bdtReferenceId)
			var bdtData models.BdtData

			if err := c.ShouldBindJSON(&bdtData); err != nil {
				log.Panic(err.Error())
			}
			spew.Dump(bdtData)
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
		if i == 1 {
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
func TestCreateBDTPolicy(t *testing.T) {

	configuration := Npcf_BDTPolicyControl.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29507")
	client := Npcf_BDTPolicyControl.NewAPIClient(configuration)

	// get test data
	bdtReqData := TestBDTPolicy.GetCreateTestData()

	// test create service
	rsp, httpRsp, err := client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqData)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	if assert.Equal(t, http.StatusCreated, httpRsp.StatusCode) {
		if assert.NotNil(t, rsp.BdtReqData) {
			assert.Equal(t, *rsp.BdtReqData, bdtReqData)
		}
		if assert.NotNil(t, rsp.BdtPolData) {
			assert.True(t, rsp.BdtPolData.SelTransPolicyId == 1)
			assert.Equal(t, rsp.BdtPolData.SuppFeat, "")
			if assert.True(t, len(rsp.BdtPolData.TransfPolicies) == 1) {
				assert.Equal(t, rsp.BdtPolData.TransfPolicies[0], models.TransferPolicy{
					RatingGroup:   1,
					RecTimeInt:    bdtReqData.DesTimeInt,
					TransPolicyId: 1,
				})
			}
		}
	}
	time.Sleep(30 * time.Millisecond)
}
