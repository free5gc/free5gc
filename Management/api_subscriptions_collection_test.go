package Management_test

import (
	"context"
	"crypto/tls"
	"encoding/json"

	//"encoding/json"

	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/Management"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/nrf_handler"
	"free5gc/src/nrf/nrf_util"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var SubscriptionsDbName = "free5gc"
var SubscriptionsCollName = "Subscriptions"
var SubscriptionsDbUrl = "mongodb://140.113.214.205:30030"

func TestCreateSubscription(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		Management.AddService(router)

		server := http.Server{
			Addr: "127.0.0.1:29510",
			TLSConfig: &tls.Config{
				KeyLogWriter: kl,
			},

			Handler: router,
		}
		_ = server.ListenAndServeTLS(nrf_util.NrfPemPath, nrf_util.NrfKeyPath)

	}()
	time.Sleep(time.Duration(2) * time.Second)

	go nrf_handler.Handle()

	// Connect to mongoDB
	MongoDBLibrary.SetMongoDB(SubscriptionsDbName, SubscriptionsDbUrl)

	// Clear Subscription DB
	MongoDBLibrary.RestfulAPIDeleteMany("Subscriptions", bson.M{})

	var subscrCond = models.NfTypeCond{
		NfType: "NRF",
	}

	var subscrCondInterface interface{}
	tmp, _ := json.Marshal(subscrCond)
	json.Unmarshal(tmp, &subscrCondInterface)

	//Create test data
	var testData = models.NrfSubscriptionData{
		NfStatusNotificationUri: "nctu.edu.tw",
		SubscriptionId:          "",
		SubscrCond:              &subscrCondInterface,
		PlmnId: &models.PlmnId{
			Mcc: "100",
			Mnc: "10",
		},
	}

	//Create Correct format
	var correctData = models.NrfSubscriptionData{
		NfStatusNotificationUri: "nctu.edu.tw",
		SubscriptionId:          "1",
		SubscrCond:              &subscrCondInterface,
		PlmnId: &models.PlmnId{
			Mcc: "100",
			Mnc: "10",
		},
	}

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Set test data (with expected data)

	// Check test data (Use RESTful GET)
	rep, res, err := client.SubscriptionsCollectionApi.CreateSubscription(context.TODO(), testData)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

	if status := res.StatusCode; status == http.StatusCreated {
		if reflect.DeepEqual(rep, correctData) == false {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rep, correctData)
		}
	}
	//t.Logf("%+v", rep)

}
