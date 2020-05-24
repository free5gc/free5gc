package management_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/openapi/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/management"
	"free5gc/src/nrf/urilist"
	"free5gc/src/nrf/util"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//var collName = "NfProfileManagementTest"
var GetNFInstancesdbName = "free5gc"
var GetNFInstancescollName = "UriListTest"
var GetNFInstancesdbUrl = "mongodb://140.113.214.205:30030"

func TestGetNFInstances(t *testing.T) {
	// run GetNFInstances Server Routine
	go func() {
		kl, _ := os.OpenFile(util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		management.AddService(router)

		server := http.Server{
			Addr: "127.0.0.1:29510",
			TLSConfig: &tls.Config{
				KeyLogWriter: kl,
			},

			Handler: router,
		}
		_ = server.ListenAndServeTLS(util.NrfPemPath, util.NrfKeyPath)

	}()
	time.Sleep(time.Duration(2) * time.Second)

	//Connect to mongoDB
	MongoDBLibrary.SetMongoDB(GetNFInstancesdbName, GetNFInstancesdbUrl)

	//set filter
	nftype := "SMF"
	filter := bson.M{"nfType": nftype}

	//set test data
	var testData = urilist.urilist{
		NfType: models.NfType_SMF,
		Link: urilist.Links{
			Item: []urilist.Item{
				{
					Href: "wwww.12",
				},
				{
					Href: "wwww.22",
				},
				{
					Href: "wwww.33",
				},
			},
		},
	}

	//Convert into map[string] interface
	tmp, _ := json.Marshal(testData)
	var putData map[string]interface{}
	json.Unmarshal(tmp, &putData)

	//Put one into mongoDB
	if MongoDBLibrary.RestfulAPIPutOne(GetNFInstancescollName, filter, putData) {
		t.Logf("Put UpdateOne")
	} else {
		t.Logf("Put InsertOne")
	}

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Set query data (with expected data)
	var localVarOptionals Nnrf_NFManagement.GetNFInstancesParamOpts

	localVarOptionals.NfType = optional.NewInterface("SMF")
	localVarOptionals.Limit = optional.NewInt32(3)

	// Check test data (Use RESTful GET)
	rep, res, err := client.NFInstancesStoreApi.GetNFInstances(context.TODO(), &localVarOptionals)
	delete(rep, "_id")
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
		if reflect.DeepEqual(putData, rep) == false {
			t.Errorf("handler returned wrong DeepEqual: gott \n%v \nwant \n%v",
				rep, putData)
		}
	}

	//t.Logf("\n%+v\n%+v", rep, putData)
}
