package Nudm_Callback_test

import (
	"context"
	"flag"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_consumer"
	"free5gc/src/udm/udm_context"
	"free5gc/src/udr/udr_service"
	"net/http"
	"strings"
	"testing"
	"time"
)

func nrfInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	nrf := &nrf_service.NRF{}
	nrf.Initialize(c)
	go nrf.Start()
	time.Sleep(100 * time.Millisecond)
}

func udrInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	udr := &udr_service.UDR{}
	udr.Initialize(c)
	go udr.Start()
	time.Sleep(100 * time.Millisecond)
}

func TestOnDataChangeNotifification(t *testing.T) {
	nrfInit()
	udrInit()

	go func() { // Simulate NF consumer server
		udmLogPath := path_util.Gofree5gcPath("free5gc/udrmslkey.log")
		udmPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.pem")
		udmKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.key")

		router := gin.Default()

		router.POST("/DataChangeNotify", func(c *gin.Context) {
			logger.Handlelog.Infoln("== Comsumer received notifiction from UDR ==")
			var dataChangeNotify models.DataChangeNotify
			if err := c.ShouldBindJSON(&dataChangeNotify); err != nil {
				logger.Handlelog.Panic(err.Error())
			}
			spew.Dump(dataChangeNotify)
			c.JSON(http.StatusNoContent, nil)
		})

		server, err := http2_util.NewServer(":9999", udmLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udmPemPath, udmKeyPath))
			assert.True(t, err == nil)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	udm_context.TestInit()

	// Set client and set url
	ueId := "imsi-0123456789"
	date := time.Now()
	duration, _ := time.ParseDuration("30m")
	date = date.Add(duration)
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	clientAPI := udm_consumer.UDMClientToUDR(ueId, false)

	var subscriptionDataSubscriptions = models.SubscriptionDataSubscriptions{
		CallbackReference: "https://127.0.0.1:9999/DataChangeNotify",
		Expiry:            &dateFormat,
		UeId:              ueId,
	}

	var subsUri string
	var subsId string

	// Create SubscriptionDataSubscriptions
	{
		subscriptionDataSubscriptions, res, err := clientAPI.SubsToNofifyCollectionApi.PostSubscriptionDataSubscriptions(context.TODO(), subscriptionDataSubscriptions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if status := res.StatusCode; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		subsUri = res.Header.Get("Location")
		spew.Printf("[subsUri_Header_Location] %s\n", subsUri)
		subsId = subsUri[strings.LastIndex(subsUri, "/")+1:]
		spew.Printf("[subsId] %s\n", subsId)
		spew.Dump(subscriptionDataSubscriptions)
	}
}
