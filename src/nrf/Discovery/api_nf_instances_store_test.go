package Discovery_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/AccessToken"
	"free5gc/src/nrf/Discovery"
	"free5gc/src/nrf/Management"
	"free5gc/src/nrf/dataconv"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/nrf_handler"
	"free5gc/src/nrf/nrf_util"
	"log"
	"net/http"
	"os"
	"reflect"
	//"strconv"
	"testing"
	"time"
	//"github.com/stretchr/testify/assert"
)

var dbName = "free5gc"
var dbAddr = "mongodb://140.113.214.205:30030"

func TestQueryParamRequestNfType(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo: &models.UpfInfo{
				SNssaiUpfInfoList: &[]models.SnssaiUpfInfoItem{
					{
						SNssai: &models.Snssai{
							Sst: 0,
						},
						DnnUpfInfoList: &[]models.DnnUpfInfoItem{
							{
								Dnn: "upf_dnn",
							},
						},
					},
				},
			},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
					{
						SNssai: &models.Snssai{},
						DnnSmfInfoList: &[]models.DnnSmfInfoItem{
							{
								Dnn: "smf_dnn",
							},
						},
					},
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo: &models.BsfInfo{
				DnnList: []string{
					"bsf_dnn",
					"nrf_dnn",
				},
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second)

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0, 2},
		{1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{},
		{},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_NRF,
		models.NfType_SMF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_NRF,
		models.NfType_UDM,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamRequesterNfInstanceFqdn(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo: &models.UpfInfo{
				SNssaiUpfInfoList: &[]models.SnssaiUpfInfoItem{
					{
						SNssai: &models.Snssai{
							Sst: 0,
						},
						DnnUpfInfoList: &[]models.DnnUpfInfoItem{
							{
								Dnn: "upf_dnn",
							},
						},
					},
				},
			},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					/*AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},*/
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
					{
						SNssai: &models.Snssai{},
						DnnSmfInfoList: &[]models.DnnSmfInfoItem{
							{
								Dnn: "smf_dnn",
							},
						},
					},
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo: &models.BsfInfo{
				DnnList: []string{
					"bsf_dnn",
					"nrf_dnn",
				},
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second)

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0, 1, 2},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			RequesterNfInstanceFqdn: optional.NewString("nfdomain4"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_NRF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_NRF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestSearchNFInstances(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority:             1,
			Capacity:             1,
			Load:                 1,
			Locality:             "NCTU",
			UdrInfo:              &models.UdrInfo{},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "112",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 221,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi1",
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority:             1,
			Capacity:             1,
			Load:                 1,
			Locality:             "NCTU",
			UdrInfo:              &models.UdrInfo{},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NSMF_PDUSESSION,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	err1 := MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	if err1 == true {
		t.Errorf("1 error")
	}
	err2 := MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	if err2 == true {
		t.Errorf("2 error")
	}
	time.Sleep(time.Duration(2) * time.Second)
	//log.Println(testDataSliceMapInterface[0])
	//log.Println(testDataSliceMapInterface[1])

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
		{0, 1},
		{0, 1},
		{0},
		{1},
		{0, 1},
		{0, 1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NNRF_DISC}),
		},
		{
			RequesterNfInstanceFqdn: optional.NewString("nfdomain4"),
		},
		{
			TargetPlmnList: optional.NewInterface("{\"mcc\": \"111\",\"mnc\": \"111\"},{\"mcc\": \"112\",\"mnc\": \"111\"}"),
		},
		{
			Snssais: optional.NewInterface(
				/*[]models.Snssai{
					{
						Sst: 1,
						Sd:  "2",
					},
				},*/
				[]string{
					"{\"sst\":222,\"sd\":\"SNssais\"}",
					//"{sst: 222, sd: \"SNssais\"}",
				},
			),
		},
		{
			TargetNfInstanceId: optional.NewInterface("1"),
		},
		{
			TargetNfFqdn: optional.NewString("fqdn"),
		},
		{
			NsiList: optional.NewInterface("nsi0"),
		},
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := models.NfType_NRF
		requesterNfType := models.NfType_NRF
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamServiceNames(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority:             1,
			Capacity:             1,
			Load:                 1,
			Locality:             "NCTU",
			UdrInfo:              &models.UdrInfo{},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "112",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 221,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi1",
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority:             1,
			Capacity:             1,
			Load:                 1,
			Locality:             "NCTU",
			UdrInfo:              &models.UdrInfo{},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NSMF_PDUSESSION,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "112",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 221,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi1",
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority:             1,
			Capacity:             1,
			Load:                 1,
			Locality:             "NCTU",
			UdrInfo:              &models.UdrInfo{},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NSMF_PDUSESSION,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NNRF_DISC}),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_NRF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_NRF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamSNSSAI(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority:             1,
			Capacity:             1,
			Load:                 1,
			Locality:             "NCTU",
			UdrInfo:              &models.UdrInfo{},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "112",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 221,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi1",
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority:             1,
			Capacity:             1,
			Load:                 1,
			Locality:             "NCTU",
			UdrInfo:              &models.UdrInfo{},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NSMF_PDUSESSION,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority:             1,
			Capacity:             1,
			Load:                 1,
			Locality:             "NCTU",
			UdrInfo:              &models.UdrInfo{},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0, 2},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			Snssais: optional.NewInterface(
				[]string{
					"{\"sst\":222,\"sd\":\"SNssais\"}",
				},
			),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_NRF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_NRF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamDnn(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_UPF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo: &models.UpfInfo{
				SNssaiUpfInfoList: &[]models.SnssaiUpfInfoItem{
					{
						SNssai: &models.Snssai{
							Sst: 0,
						},
						DnnUpfInfoList: &[]models.DnnUpfInfoItem{
							{
								Dnn: "upf_dnn",
							},
						},
					},
				},
			},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
					{
						SNssai: &models.Snssai{},
						DnnSmfInfoList: &[]models.DnnSmfInfoItem{
							{
								Dnn: "smf_dnn",
							},
						},
					},
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{
				/*DnnList: []string{
					"bsf_dnn",
					"nrf_dnn",
				},*/
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "3",
			NfType:         models.NfType_PCF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo: &models.PcfInfo{
				DnnList: []string{
					"pcf_dnn",
					"nrf_dnn",
				},
			},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "3"}, testDataSliceMapInterface[3])
	time.Sleep(time.Duration(1) * time.Second)

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
		{1},
		{2},
		{3},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			Dnn: optional.NewString("upf_dnn"),
		},
		{
			Dnn: optional.NewString("smf_dnn"),
		},
		{
			Dnn: optional.NewString("bsf_dnn"),
		},
		{
			Dnn: optional.NewString("pcf_dnn"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_UPF,
		models.NfType_SMF,
		models.NfType_BSF,
		models.NfType_PCF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_UPF,
		models.NfType_SMF,
		models.NfType_BSF,
		models.NfType_PCF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamSmfServingArea(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_UPF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo: &models.UpfInfo{
				SNssaiUpfInfoList: &[]models.SnssaiUpfInfoItem{
					{
						SNssai: &models.Snssai{
							Sst: 0,
						},
						DnnUpfInfoList: &[]models.DnnUpfInfoItem{
							{
								Dnn: "upf_dnn",
							},
						},
					},
				},
				/*SmfServingArea: []string{
					"nctu",
				},*/
			},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
					{
						SNssai: &models.Snssai{},
						DnnSmfInfoList: &[]models.DnnSmfInfoItem{
							{
								Dnn: "smf_dnn",
							},
						},
					},
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo: &models.BsfInfo{
				DnnList: []string{
					"bsf_dnn",
					"nrf_dnn",
				},
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			SmfServingArea: optional.NewString("nctu"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_UPF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_UPF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamTai(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_AMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo: &models.AmfInfo{
				TaiList: &[]models.Tai{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						Tac: "ABC",
					},
				},
			},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
					{
						SNssai: &models.Snssai{},
						DnnSmfInfoList: &[]models.DnnSmfInfoItem{
							{
								Dnn: "smf_dnn",
							},
						},
					},
				},
				TaiList: &[]models.Tai{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						Tac: "ABC",
					},
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo: &models.BsfInfo{
				DnnList: []string{
					"bsf_dnn",
					"nrf_dnn",
				},
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
		{1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			Tai: optional.NewInterface("{\"plmnId\":{\"mcc\": \"111\", \"mnc\":\"111\"}, \"tac\": \"ABC\"}"),
		},
		{
			Tai: optional.NewInterface("{\"plmnId\":{\"mcc\": \"111\", \"mnc\":\"111\"}, \"tac\": \"ABC\"}"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_AMF,
		models.NfType_SMF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_AMF,
		models.NfType_SMF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamAmfRegionId(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_AMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo: &models.AmfInfo{
				TaiList: &[]models.Tai{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						Tac: "ABC",
					},
				},
				AmfRegionId: "nctu",
			},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
					{
						SNssai: &models.Snssai{},
						DnnSmfInfoList: &[]models.DnnSmfInfoItem{
							{
								Dnn: "smf_dnn",
							},
						},
					},
				},
				TaiList: &[]models.Tai{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						Tac: "ABC",
					},
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo: &models.BsfInfo{
				DnnList: []string{
					"bsf_dnn",
					"nrf_dnn",
				},
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			AmfRegionId: optional.NewString("nctu"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_AMF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_AMF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamAmfSetId(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_AMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo: &models.AmfInfo{
				TaiList: &[]models.Tai{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						Tac: "ABC",
					},
				},
				AmfRegionId: "nctu",
				AmfSetId:    "123",
			},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
					{
						SNssai: &models.Snssai{},
						DnnSmfInfoList: &[]models.DnnSmfInfoItem{
							{
								Dnn: "smf_dnn",
							},
						},
					},
				},
				TaiList: &[]models.Tai{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						Tac: "ABC",
					},
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo: &models.BsfInfo{
				DnnList: []string{
					"bsf_dnn",
					"nrf_dnn",
				},
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			AmfSetId: optional.NewString("123"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_AMF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_AMF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamGuami(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_AMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo: &models.AmfInfo{
				TaiList: &[]models.Tai{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						Tac: "ABC",
					},
				},
				AmfRegionId: "nctu",
				AmfSetId:    "123",
				GuamiList: &[]models.Guami{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						AmfId: "123",
					},
				},
			},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
					{
						SNssai: &models.Snssai{},
						DnnSmfInfoList: &[]models.DnnSmfInfoItem{
							{
								Dnn: "smf_dnn",
							},
						},
					},
				},
				TaiList: &[]models.Tai{
					{
						PlmnId: &models.PlmnId{
							Mcc: "111",
							Mnc: "111",
						},
						Tac: "ABC",
					},
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo: &models.BsfInfo{
				DnnList: []string{
					"bsf_dnn",
					"nrf_dnn",
				},
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			Guami: optional.NewInterface("{\"plmnId\":{\"mcc\":\"111\", \"mnc\":\"111\"},\"amfId\":\"123\"}"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_AMF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_AMF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamSupi(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_PCF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{
				/*SupiRanges: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},*/
			},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo:  &models.ChfInfo{
				/*SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},*/
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_AUSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{
				/*SupiRanges: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},*/
			},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "3",
			NfType:         models.NfType_UDM,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo: &models.UdmInfo{
				SupiRanges: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "4",
			NfType:         models.NfType_UDR,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{
				/*SupiRanges: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},*/
			},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "3"}, testDataSliceMapInterface[3])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "4"}, testDataSliceMapInterface[4])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
		{1},
		{2},
		{3},
		{4},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			Supi: optional.NewString("imsi-10001"),
		},
		{
			Supi: optional.NewString("imsi-10001"),
		},
		{
			Supi: optional.NewString("imsi-10001"),
		},
		{
			Supi: optional.NewString("imsi-10001"),
		},
		{
			Supi: optional.NewString("imsi-10001"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_PCF,
		models.NfType_CHF,
		models.NfType_AUSF,
		models.NfType_UDM,
		models.NfType_UDR,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_PCF,
		models.NfType_CHF,
		models.NfType_AUSF,
		models.NfType_UDM,
		models.NfType_UDR,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamUeIpv4Address(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{
				/*Ipv4AddressRanges: &[]models.Ipv4AddressRange{
					{
						Start: strconv.Itoa(int(dataconv.Ipv4ToInt("140.113.1.1"))),
						End:   strconv.Itoa(int(dataconv.Ipv4ToInt("140.113.1.5"))),
					},
				},*/
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "2",
			NfType:         models.NfType_AUSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{
				SupiRanges: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var expectedDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{
				/*Ipv4AddressRanges: &[]models.Ipv4AddressRange{
					{
						Start: "140.113.1.1",
						End:   "140.113.1.5",
					},
				},*/
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			UeIpv4Address: optional.NewString("140.113.1.2"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_BSF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_BSF,
	}

	for i := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		/*var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}*/

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedDatas)

		if reflect.DeepEqual(rep.NfInstances, expectedDatas) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamIpDomain(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{
				/*IpDomainList: []string{
					"1234",
				},*/
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			IpDomain: optional.NewString("1234"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_BSF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_BSF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamUeIpv6Prefix(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{
				/*Ipv6PrefixRanges: &[]models.Ipv6PrefixRange{
					{
						Start: dataconv.Ipv6ToInt("2001:db6::").String(),
						End:   dataconv.Ipv6ToInt("2001:db7::").String(),
					},
				},*/
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	var expectedDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_BSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{
				/*Ipv6PrefixRanges: &[]models.Ipv6PrefixRange{
					{
						Start: dataconv.Ipv6IntToIpv6String(dataconv.Ipv6ToInt("2001:db6::")),
						End:   dataconv.Ipv6IntToIpv6String(dataconv.Ipv6ToInt("2001:db9::")),
					},
				},*/
			},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			UeIpv6Prefix: optional.NewInterface("2001:db8::"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_BSF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_BSF,
	}

	for i := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		/*var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}*/

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedDatas)

		if reflect.DeepEqual(rep.NfInstances, expectedDatas) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamPgwInd(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				PgwFqdn: "123",
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	err := MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	if err {
		t.Logf("2222")
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			PgwInd: optional.NewBool(true),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_SMF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_SMF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamPgw(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				PgwFqdn: "123",
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			Pgw: optional.NewString("123"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_SMF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_SMF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamGpsi(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo:  &models.ChfInfo{
				/*GpsiRangeList: &[]models.IdentityRange{
					{
						Start: "00000",
						End:   "10000",
					},
				},*/
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_UDR,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{
				/*GpsiRanges: &[]models.IdentityRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},*/
			},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_UDM,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo: &models.UdmInfo{
				GpsiRanges: &[]models.IdentityRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
		{1},
		{2},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			Gpsi: optional.NewString("msisdn-10000"),
		},
		{
			Gpsi: optional.NewString("msisdn-10000"),
		},
		{
			Gpsi: optional.NewString("msisdn-10000"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_CHF,
		models.NfType_UDR,
		models.NfType_UDM,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_CHF,
		models.NfType_UDR,
		models.NfType_UDM,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamDataSet(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_UDR,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo: &models.UdrInfo{
				SupportedDataSets: []models.DataSetId{
					models.DataSetId_SUBSCRIPTION,
				},
			},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_UDR,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0, 1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			DataSet: optional.NewInterface(models.DataSetId_SUBSCRIPTION),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_UDR,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_UDR,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamRoutingIndicator(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_AUSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{
				/*RoutingIndicators: []string{
					"routingIndicators",
				},*/
			},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_UDM,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo: &models.UdmInfo{
				RoutingIndicators: []string{
					"routingIndicators",
				},
			},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
		{1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			RoutingIndicator: optional.NewString("routingIndicators"),
		},
		{
			RoutingIndicator: optional.NewString("routingIndicators"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_AUSF,
		models.NfType_UDM,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_AUSF,
		models.NfType_UDM,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamGroupIdList(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_AUSF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{
				GroupId: "12345",
			},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_UDM,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo: &models.UdmInfo{
				GroupId: "12345",
			},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_UDR,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo: &models.UdrInfo{
				GroupId: "12345",
			},
			UdmInfo:              &models.UdmInfo{},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "2"}, testDataSliceMapInterface[2])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
		{1},
		{2},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			GroupIdList: optional.NewInterface("12345,123,123"),
		},
		{
			GroupIdList: optional.NewInterface("12345,123,123"),
		},
		{
			GroupIdList: optional.NewInterface("12345,123,123"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_AUSF,
		models.NfType_UDM,
		models.NfType_UDR,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_AUSF,
		models.NfType_UDM,
		models.NfType_UDR,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamDnaiList(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_UPF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo: &models.UpfInfo{
				SNssaiUpfInfoList: &[]models.SnssaiUpfInfoItem{
					{
						DnnUpfInfoList: &[]models.DnnUpfInfoItem{
							{
								DnaiList: []string{
									"111",
									"222",
								},
							},
						},
					},
				},
			},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			DnaiList: optional.NewInterface("111"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_UPF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_UPF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamUpfIwkEpsInd(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_UPF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo: &models.UpfInfo{
				IwkEpsInd: true,
			},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			UpfIwkEpsInd: optional.NewBool(true),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_UPF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_UPF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamChfSupportedPlmn(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo:  &models.ChfInfo{
				/*PlmnRangeList: &[]models.PlmnRange{
					{
						Start: "100100",
						End:   "100200",
					},
				},*/
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				PlmnRangeList: &[]models.PlmnRange{
					{
						Start: "20000",
						End:   "20000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			ChfSupportedPlmn: optional.NewInterface("{\"mcc\":\"100\",\"mnc\":\"101\"}"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_CHF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_CHF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamPreferredLocality(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_AMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				PlmnRangeList: &[]models.PlmnRange{
					{
						Start: "100100",
						End:   "100200",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_CHF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			PreferredLocality: optional.NewString("NCTU"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_CHF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_CHF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamAccessType(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				AccessType: []models.AccessType{
					models.AccessType__3_GPP_ACCESS,
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			AccessType: optional.NewInterface(models.AccessType_NON_3_GPP_ACCESS),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_SMF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_SMF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamSupportedFeatures(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_SMF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo: &models.SmfInfo{
				AccessType: []models.AccessType{
					models.AccessType__3_GPP_ACCESS,
				},
			},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_NRF,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfTypes: []models.NfType{
				models.NfType_NRF,
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
					SupportedFeatures: "discovery",
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			SupportedFeatures: optional.NewString("discovery"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_NRF,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_NRF,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}

func TestQueryParamExternalGroupIdentity(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := gin.Default()

		AccessToken.AddService(router)
		Discovery.AddService(router)
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
	time.Sleep(time.Duration(1) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB(dbName, dbAddr)
	collectionName := "NfProfile"

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	// clear mongoDB
	MongoDBLibrary.RestfulAPIDeleteMany(collectionName, bson.M{})

	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	// set data into DB
	var testDatas = []models.NfProfile{
		{
			NfInstanceId:   "0",
			NfType:         models.NfType_UDM,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo:  &models.UdrInfo{},
			UdmInfo: &models.UdmInfo{
				ExternalGroupIdentifiersRanges: &[]models.IdentityRange{
					{
						Start: dataconv.EncodeGroupId("abcdef00-100-100-1000000001"),
						End:   dataconv.EncodeGroupId("abcdef00-100-100-1000000005"),
					},
				},
			},
			AusfInfo:             &models.AusfInfo{},
			AmfInfo:              &models.AmfInfo{},
			SmfInfo:              &models.SmfInfo{},
			UpfInfo:              &models.UpfInfo{},
			PcfInfo:              &models.PcfInfo{},
			BsfInfo:              &models.BsfInfo{},
			ChfInfo:              &models.ChfInfo{},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
				},
			},
		},
		{
			NfInstanceId:   "1",
			NfType:         models.NfType_UDR,
			NfStatus:       models.NfStatus_REGISTERED,
			HeartBeatTimer: 10,
			PlmnList: &[]models.PlmnId{ // Pattern: '^[0-9]{3}[0-9]{2,3}$'
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			SNssais: &[]models.Snssai{ // range 0-255
				{
					Sst: 222,
					Sd:  "SNssais",
				},
			},
			NsiList: []string{
				"nsi0",
			},
			Fqdn:          "fqdn",
			InterPlmnFqdn: "InterPlmnFqdn",
			Ipv4Addresses: []string{
				"140.113.1.1",
			},
			Ipv6Addresses: []string{
				"fc00::",
			},
			AllowedPlmns: &[]models.PlmnId{
				{
					Mcc: "111",
					Mnc: "111",
				},
			},
			AllowedNfDomains: []string{
				"nfdomain1",
			},
			AllowedNssais: &[]models.Snssai{
				{
					Sst: 333,
					Sd:  "AllowedNssais",
				},
			},
			Priority: 1,
			Capacity: 1,
			Load:     1,
			Locality: "NCTU",
			UdrInfo: &models.UdrInfo{
				ExternalGroupIdentifiersRanges: &[]models.IdentityRange{
					{
						Start: dataconv.EncodeGroupId("abcdef00-100-100-1000000001"),
						End:   dataconv.EncodeGroupId("abcdef00-100-100-1000000005"),
					},
				},
			},
			UdmInfo:  &models.UdmInfo{},
			AusfInfo: &models.AusfInfo{},
			AmfInfo:  &models.AmfInfo{},
			SmfInfo:  &models.SmfInfo{},
			UpfInfo:  &models.UpfInfo{},
			PcfInfo:  &models.PcfInfo{},
			BsfInfo:  &models.BsfInfo{},
			ChfInfo: &models.ChfInfo{
				SupiRangeList: &[]models.SupiRange{
					{
						Start: "10000",
						End:   "11000",
					},
				},
			},
			NrfInfo:              &models.NrfInfo{},
			CustomInfo:           &map[string]interface{}{},
			RecoveryTime:         &dateFormat,
			NfServicePersistence: true,
			NfServices: &[]models.NfService{
				{
					ServiceName:     models.ServiceName_NNRF_DISC,
					NfServiceStatus: models.NfServiceStatus_REGISTERED,
					AllowedNfDomains: []string{
						"nfdomain3",
						"nfdomain4",
					},
					SupportedFeatures: "discovery",
				},
			},
		},
	}

	var testDataSliceMapInterface []map[string]interface{}

	for _, testData := range testDatas {
		var testDataMapInterface map[string]interface{}

		testDataByteArray, _ := json.Marshal(testData)
		_ = json.Unmarshal(testDataByteArray, &testDataMapInterface)

		testDataSliceMapInterface = append(testDataSliceMapInterface, testDataMapInterface)
	}
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "0"}, testDataSliceMapInterface[0])
	MongoDBLibrary.RestfulAPIPutOne(collectionName, bson.M{"nfInstanceId": "1"}, testDataSliceMapInterface[1])
	time.Sleep(time.Duration(2) * time.Second) // must set sleep time to wait for storing data!!!

	// Set Expected Data Index
	var expectedDataIndexTable = [][]int{
		{0},
		{1},
	}

	// Set Test Query Param
	localVarOptionalsTable := []Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		{
			ExternalGroupIdentity: optional.NewString("abcdef00-100-100-1000000001"),
		},
		{
			ExternalGroupIdentity: optional.NewString("abcdef00-100-100-1000000003"),
		},
	}

	// Set targetNfType
	targetNfTypeTable := []models.NfType{
		models.NfType_UDM,
		models.NfType_UDR,
	}

	// Set requesterNfType
	requesterNfTypeTable := []models.NfType{
		models.NfType_UDM,
		models.NfType_UDR,
	}

	for i, expectedDataIndex := range expectedDataIndexTable {
		// Set test data (with expected data)
		targetNfType := targetNfTypeTable[i]
		requesterNfType := requesterNfTypeTable[i]
		localVarOptionals := localVarOptionalsTable[i]

		// Check test data (Use RESTful GET)
		rep, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
		if err != nil {
			logger.AppLog.Panic(err)
		}
		if res != nil {
			if status := res.StatusCode; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}

		// build Expected Data
		var expectedData []models.NfProfile
		for _, index := range expectedDataIndex {
			expectedData = append(expectedData, testDatas[index])
		}

		log.Println("Test: ", i)
		log.Printf("Output:\t%+v\n", rep.NfInstances)
		log.Printf("Expect:\t%+v\n\n", expectedData)

		if reflect.DeepEqual(rep.NfInstances, expectedData) != true {
			t.Errorf("TEST %d Not correct", i)
		}
	}
}
