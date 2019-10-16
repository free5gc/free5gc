package amf_consumer_test

import (
	"flag"
	"github.com/antihax/optional"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/nrf/nrf_service"
	"reflect"
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

func TestSendSearchNFInstances(t *testing.T) {

	nrfInit()

	time.Sleep(200 * time.Millisecond)
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})

	// Init AMF
	TestAmf.AmfInit()

	time.Sleep(100 * time.Millisecond)

	nfprofile, err := amf_consumer.BuildNFInstance(TestAmf.TestAmf)
	if err != nil {
		t.Error(err.Error())
	}

	uri, err1 := amf_consumer.SendRegisterNFInstance(TestAmf.TestAmf.NrfUri, TestAmf.TestAmf.NfId, nfprofile)
	if err1 != nil {
		t.Error(err1.Error())
	} else {
		TestAmf.Config.Dump(uri)
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NAMF_COMM}),
	}
	result, err2 := amf_consumer.SendSearchNFInstances(TestAmf.TestAmf.NrfUri, models.NfType_AMF, models.NfType_AMF, param)
	if err2 != nil {
		t.Error(err1.Error())
	} else if !reflect.DeepEqual(nfprofile, result.NfInstances[0]) {
		t.Error("failed for expected value mismatch")
	}
}
