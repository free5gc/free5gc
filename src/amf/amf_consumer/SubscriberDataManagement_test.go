package amf_consumer_test

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
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
	ueId := "imsi-2089300007487"
	servingPlmnId := "20893"
	testData := models.AccessAndMobilitySubscriptionData{
		UeUsageType: 1,
	}
	tmp, _ := json.Marshal(testData)
	var insertTestData = bson.M{}
	if err := json.Unmarshal(tmp, &insertTestData); err != nil {
		t.Error(err)
	}
	insertTestData["ueId"] = ueId
	insertTestData["servingPlmnId"] = servingPlmnId
	if result, err := collection.InsertOne(context.TODO(), insertTestData); err != nil {
		t.Errorf("insert test data error: %+v", err)
	} else {
		t.Logf("insert test data result: %+v", result)
	}

	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	// udmUri := "https://localhost:29503"
	problemDetails, err := amf_consumer.SDMGetAmData(ue)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	}

	if _, err := collection.DeleteOne(context.TODO(), bson.M{"ueId": "imsi-2089300007487"}); err != nil {
		t.Error(err.Error())
	}
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
