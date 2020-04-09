package amf_consumer_test

import (
	"context"
	"flag"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/UDR/TestRegistrationProcedure"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/udm/udm_service"
	"free5gc/src/udr/udr_service"
	"strings"
	"testing"
	"time"
)

var testflags flag.FlagSet
var testC = cli.NewContext(nil, &testflags, nil)
var testAmData = TestRegistrationProcedure.TestAmDataTable[TestRegistrationProcedure.FREE5GC_CASE]
var servingPlmnId = "20893"

func udminit() {
	udm := &udm_service.UDM{}
	udm.Initialize(testC)
	go udm.Start()
	time.Sleep(100 * time.Millisecond)
}

func udrinit() {
	udr := &udr_service.UDR{}
	udr.Initialize(testC)
	go udr.Start()
	time.Sleep(100 * time.Millisecond)
}

func insertAccessAndMobilitySubscriptionDataToMongoDB(ueId string, amData models.AccessAndMobilitySubscriptionData, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	putData := toBsonM(amData)
	putData["ueId"] = ueId
	putData["servingPlmnId"] = servingPlmnId
	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func delAccessAndMobilitySubscriptionDataFromMongoDB(ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}

func TestPutUpuAck(t *testing.T) {
	udminit()
	udrinit()

	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	// udmUri := "https://localhost:29503"

	upuMacIue := strings.Repeat("1", 32)
	err := amf_consumer.PutUpuAck(ue, upuMacIue)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	}

}

func TestSDMGetAmData(t *testing.T) {
	nrfInit()
	if len(TestAmf.TestAmf.AmfRanPool) == 0 {
		udminit()
		udrinit()
	}

	Client := MongoDBLibrary.Client

	// Drop old data
	collection := Client.Database("free5gc").Collection("subscriptionData.provisionedData.amData")
	if _, err := collection.DeleteOne(context.TODO(), bson.M{"ueId": "imsi-2089300007487"}); err != nil {
		t.Errorf("delete old test data error: %+v", err)
	}

	// Set test data
	insertAccessAndMobilitySubscriptionDataToMongoDB("imsi-2089300007487", testAmData, servingPlmnId)

	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ue.NudmSDMUri = "https://localhost:29503"
	problemDetails, err := amf_consumer.SDMGetAmData(ue)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	} else {
		t.Logf("Get AM Data: %+v", ue.AccessAndMobilitySubscriptionData)
	}

	delAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
}

func TestSDMGetSmfSelectData(t *testing.T) {
	// TODO: finish test
	if len(TestAmf.TestAmf.AmfRanPool) == 0 {
		udminit()
		udrinit()
	}
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	// udmUri := "https://localhost:29503"
	problemDetails, err := amf_consumer.SDMGetSmfSelectData(ue)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	}
}

func TestSDMGetUeContextInSmfData(t *testing.T) {
	// TODO: finish test
	if len(TestAmf.TestAmf.AmfRanPool) == 0 {
		udminit()
		udrinit()
	}
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	// udmUri := "https://localhost:29503"
	problemDetails, err := amf_consumer.SDMGetUeContextInSmfData(ue)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	}
}

func TestSDMSubscribe(t *testing.T) {
	// TODO: finish test
	if len(TestAmf.TestAmf.AmfRanPool) == 0 {
		udminit()
		udrinit()
	}
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	// udmUri := "https://localhost:29503"
	problemDetails, err := amf_consumer.SDMSubscribe(ue)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	}
}

func TestSDMGetSliceSelectionSubscriptionData(t *testing.T) {
	nrfInit()
	if len(TestAmf.TestAmf.AmfRanPool) == 0 {
		udminit()
		udrinit()
	}

	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	insertAccessAndMobilitySubscriptionDataToMongoDB("imsi-2089300007487", testAmData, servingPlmnId)

	ue.NudmSDMUri = "https://localhost:29503"
	problemDetails, err := amf_consumer.SDMGetSliceSelectionSubscriptionData(ue)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	} else {
		t.Logf("Get Nssai: %+v", ue.SubscribedNssai)
	}

	delAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
}
