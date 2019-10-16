package ausf_consumer_test

import (
	"flag"
	"fmt"
	"github.com/antihax/optional"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/openapi/models"
	"free5gc/src/ausf/ausf_consumer"
	"free5gc/src/ausf/ausf_context"
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

	time.Sleep(100 * time.Millisecond)

	ausf_context.TestInit()
	self := ausf_context.GetSelf()
	nfprofile, err := ausf_consumer.BuildNFInstance(self)
	if err != nil {
		t.Error(err.Error())
	}

	uri, err1 := ausf_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, nfprofile)
	if err1 != nil {
		t.Error(err1.Error())
	} else {
		fmt.Println("uri: ", uri)
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NAUSF_AUTH}),
	}
	result, err2 := ausf_consumer.SendSearchNFInstances(self.NrfUri, models.NfType_AUSF, models.NfType_AUSF, param)
	/*v := reflect.ValueOf(result.NfInstances[0])
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	fmt.Println("result ", result.NfInstances[0])
	fmt.Println("=====================================")
	v2 := reflect.ValueOf(nfprofile)
	typeOfS2 := v2.Type()
	for i := 0; i < v2.NumField(); i++ {
		fmt.Printf("Field: %s\tValue: %v\n", typeOfS2.Field(i).Name, v2.Field(i).Interface())
	}
	fmt.Println("profile ", nfprofile)*/
	if err2 != nil {
		t.Error(err1.Error())
	} else if len(result.NfInstances) > 0 && !reflect.DeepEqual(nfprofile, result.NfInstances[0]) {
		t.Error("failed for expected value mismatch")
	} else if len(result.NfInstances) == 0 {
		t.Error("len(result.NfInstances) is 0\n")
	}
}
