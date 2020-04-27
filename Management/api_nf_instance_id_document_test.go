package Management_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/Management"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/nrf_handler"
	"free5gc/src/nrf/nrf_util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	. "free5gc/lib/openapi/models"
	//"free5gc/src/nrf/nrf_handler/nrf_message"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

//var collName = "NfProfileManagementTest"
var dbName = "free5gc"
var collName = "NfProfile"
var dbUrl = "mongodb://140.113.214.205:30030"

func TestGetNFInstance(t *testing.T) {
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
	MongoDBLibrary.SetMongoDB(dbName, dbUrl)

	nfInstanceId := "97"
	filter := bson.M{"nfInstanceId": nfInstanceId}

	//Create test data
	var testData = models.NfProfile{
		NfInstanceId: nfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "NRF",
	}

	//Convert into map[string] interface
	tmp, _ := json.Marshal(testData)
	var putData = bson.M{}
	json.Unmarshal(tmp, &putData)

	//Put one into mongoDB
	if MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData) {
		t.Logf("Put UpdateOne")
	} else {
		t.Logf("Put InsertOne")
	}

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Set test data (with expected data)

	// Check test data (Use RESTful GET)
	rep, res, err := client.NFInstanceIDDocumentApi.GetNFInstance(context.TODO(), nfInstanceId)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		} else if reflect.DeepEqual(testData, rep) == false {
			t.Errorf("reflect DeepEqual wrong : got\n %v \nwant \n%v",
				rep, testData)
		}
	}
	t.Logf("%+v", rep)
}

func TestDeregisterNFInstance(t *testing.T) {
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

	go func() {
		router := gin.Default()
		router.POST("", func(c *gin.Context) {
			var ND NotificationData

			if err := c.ShouldBindJSON(&ND); err != nil {
				log.Panic(err.Error())
			}
			c.JSON(http.StatusNoContent, gin.H{})
		})

		srv, err := http2_util.NewServer(":30678", nrf_util.NrfLogPath, router)
		if err != nil {
			log.Panic(err.Error())
		}

		err2 := srv.ListenAndServeTLS(nrf_util.NrfPemPath, nrf_util.NrfKeyPath)
		if err2 != nil && err2 != http.ErrServerClosed {
			log.Panic(err2.Error())
		}
	}()

	time.Sleep(time.Duration(2) * time.Second)
	go nrf_handler.Handle()

	//Connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbUrl)
	//collName := "NfProfileManagementTest"
	nfInstanceId := "97"
	filter := bson.M{"nfInstanceId": nfInstanceId}

	//Create test data
	var testData = models.NfProfile{
		NfInstanceId: nfInstanceId,
		NfType:       "NRF",
	}
	//Convert into map[string] interface
	tmp, _ := json.Marshal(testData)
	var putData = bson.M{}
	json.Unmarshal(tmp, &putData)
	//Put one into mongoDB
	if MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData) {
		t.Logf("Put Success")
	} else {
		t.Logf("Put Fail")
	}

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Set test data (with expected data)

	// Check test data (Use RESTful GET)
	res, err := client.NFInstanceIDDocumentApi.DeregisterNFInstance(context.TODO(), nfInstanceId)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNoContent)
		}
	}
	//t.Logf("%+v", rep)
}

func TestRegisterNFInstance(t *testing.T) {
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
	go func() {
		router := gin.Default()
		router.POST("", func(c *gin.Context) {
			var ND NotificationData

			if err := c.ShouldBindJSON(&ND); err != nil {
				log.Panic(err.Error())
			}
			c.JSON(http.StatusNoContent, gin.H{})
		})

		srv, err := http2_util.NewServer(":30678", nrf_util.NrfLogPath, router)
		if err != nil {
			log.Panic(err.Error())
		}

		err2 := srv.ListenAndServeTLS(nrf_util.NrfPemPath, nrf_util.NrfKeyPath)
		if err2 != nil && err2 != http.ErrServerClosed {
			log.Panic(err2.Error())
		}
	}()
	time.Sleep(time.Duration(2) * time.Second)

	//Connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbUrl)

	// Clear NfProfile DB
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})
	time.Sleep(time.Duration(2) * time.Second)

	nfInstanceId := "97"
	//filter := bson.M{"nfInstanceId": nfInstanceId}

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	//Create test data
	var testData = models.NfProfile{
		NfInstanceId: nfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "NRF",
		RecoveryTime: &dateFormat,
	}

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Set test data (with expected data)

	// Check test data (Use RESTful GET)
	rep, res, err := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, testData)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			if status != http.StatusCreated {
				t.Errorf("handler returned wrong status code: got %v want %v or %v",
					status, http.StatusOK, http.StatusCreated)
			}
		}
	}
	if status := res.StatusCode; status == http.StatusOK || status == http.StatusCreated {
		if reflect.DeepEqual(testData, rep) == false {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rep, testData)
		}
	}
	//t.Logf("%+v \n %+v", rep, testData)
}

func TestUpdateNFInstance(t *testing.T) {
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
	go func() {
		router := gin.Default()
		router.POST("", func(c *gin.Context) {
			var ND NotificationData

			if err := c.ShouldBindJSON(&ND); err != nil {
				log.Panic(err.Error())
			}
			c.JSON(http.StatusNoContent, gin.H{})
		})

		srv, err := http2_util.NewServer(":30676", nrf_util.NrfLogPath, router)
		if err != nil {
			log.Panic(err.Error())
		}

		err2 := srv.ListenAndServeTLS(nrf_util.NrfPemPath, nrf_util.NrfKeyPath)
		if err2 != nil && err2 != http.ErrServerClosed {
			log.Panic(err2.Error())
		}
	}()
	time.Sleep(time.Duration(2) * time.Second)

	// Connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbUrl)

	// Clear NfProfile DB
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})

	nfInstanceId := "95"
	filter := bson.M{"nfInstanceId": nfInstanceId}

	//Create test data
	var testData = models.NfProfile{
		NfInstanceId: nfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "NRF",
	}
	//Create Correct format
	var correctData = models.NfProfile{
		NfInstanceId: "95",
		NfStatus:     "SUSPENDED",
		NfType:       "NRF",
	}
	//Create PatchItem
	patchItemArray := []models.PatchItem{
		{
			Op:    models.PatchOperation_REPLACE,
			Path:  "/nfStatus",
			Value: "DEREGISTERED",
		},
	}

	//Convert into map[string] interface
	tmp, _ := json.Marshal(testData)
	var putData = bson.M{}
	json.Unmarshal(tmp, &putData)

	//Put one into mongoDB
	if MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData) {
		t.Logf("Put Update")
	} else {
		t.Logf("Put Insert")
	}
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Set test data (with expected data)

	// Check test data (Use RESTful GET)
	rep, res, err := client.NFInstanceIDDocumentApi.UpdateNFInstance(context.TODO(), nfInstanceId, patchItemArray)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v or %v",
				status, http.StatusOK, http.StatusCreated)
		}
	}

	if status := res.StatusCode; status == http.StatusOK {
		if reflect.DeepEqual(rep, correctData) == false {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rep, correctData)
		}
	}
	//t.Logf("%+v", rep)
}

func TestRegisterNotification(t *testing.T) {
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
	MongoDBLibrary.SetMongoDB(dbName, dbUrl)

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// clear NfProfile DB & Subscription DB
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})
	MongoDBLibrary.RestfulAPIDeleteMany("Subscriptions", bson.M{})

	time.Sleep(time.Duration(2) * time.Second)

	// Subscription data (empty condition)
	var testData = models.NrfSubscriptionData{
		NfStatusNotificationUri: "nctu.edu.tw",
		SubscriptionId:          "",
		PlmnId: &models.PlmnId{
			Mcc: "100",
			Mnc: "10",
		},
	}

	_, res, err := client.SubscriptionsCollectionApi.CreateSubscription(context.TODO(), testData)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

	// Register Data
	nfInstanceId := "1234"

	var registerData = models.NfProfile{
		NfInstanceId: nfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "NRF",
	}

	// Register
	_, res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, registerData)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

	// Check Notification
}
