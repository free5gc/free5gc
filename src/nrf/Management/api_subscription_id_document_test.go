package Management_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/Management"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/nrf_util"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

var SubscriptionDbName = "free5gc"
var SubscriptionCollName = "Subscriptions"
var SubscriptionDbUrl = "mongodb://140.113.214.205:30030"

func TestRemoveSubscription(t *testing.T) {
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

	//Connect to mongoDB
	MongoDBLibrary.SetMongoDB(SubscriptionDbName, SubscriptionDbUrl)
	subscriptionID := "96"

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Set test data (with expected data)

	// Check test data (Use RESTful GET)
	res, err := client.SubscriptionIDDocumentApi.RemoveSubscription(context.TODO(), subscriptionID)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNoContent)
		}
	}
	t.Logf("%+v", res)
}

func TestUpdateSubscription(t *testing.T) {
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

	//Connect to mongoDB
	MongoDBLibrary.SetMongoDB(SubscriptionDbName, SubscriptionDbUrl)

	SubscriptionID := "90"
	filter := bson.M{"subscriptionID": SubscriptionID}

	//Create test data
	var testData = models.NrfSubscriptionData{
		NfStatusNotificationUri: "www.nthu.edu.tw",
		SubscriptionId:          SubscriptionID,
		PlmnId: &models.PlmnId{
			Mcc: "100",
			Mnc: "10",
		},
	}
	//Create Correct format
	var correctData = models.NrfSubscriptionData{
		NfStatusNotificationUri: "www.nctu.edu.tw",
		SubscriptionId:          SubscriptionID,
		PlmnId: &models.PlmnId{
			Mcc: "100",
			Mnc: "10",
		},
	}
	//Create PatchItem
	patchItemArray := []models.PatchItem{
		{
			Op:    models.PatchOperation_REPLACE,
			Path:  "/nfStatusNotificationUri",
			Value: "www.nctu.edu.tw",
		},
	}

	//Convert into map[string] interface
	tmp, _ := json.Marshal(testData)
	var putData = bson.M{}
	json.Unmarshal(tmp, &putData)

	//Put one into mongoDB
	if MongoDBLibrary.RestfulAPIPost(SubscriptionCollName, filter, putData) {
		t.Logf("Post UpdateOne")
	} else {
		t.Logf("Post InsertOne")
	}
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Set test data (with expected data)

	// Check test data (Use RESTful GET)
	rep, res, err := client.SubscriptionIDDocumentApi.UpdateSubscription(context.TODO(), SubscriptionID, patchItemArray)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

	if status := res.StatusCode; status == http.StatusOK {
		if reflect.DeepEqual(rep, correctData) == false {
			t.Errorf("handler returned wrong DeepEqual: got %v want %v",
				rep, correctData)
		}
	}
	//t.Logf("%+v", rep)

}
