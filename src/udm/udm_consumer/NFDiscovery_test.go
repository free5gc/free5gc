package udm_consumer_test

import (
	"flag"
	"fmt"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/udm/udm_consumer"
	"free5gc/src/udm/udm_context"
	"reflect"
	"testing"
	"time"

	"github.com/antihax/optional"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
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
	time.Sleep(200 * time.Millisecond)

	udm_context.TestInit()
	self := udm_context.UDM_Self()
	nfprofile, err := udm_consumer.BuildNFInstance(self)
	if err != nil {
		t.Error(err.Error())
	}

	uri, _, err1 := udm_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, nfprofile)
	if err1 != nil {
		t.Error(err1.Error())
	} else {
		fmt.Println("uri: ", uri)
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NUDM_SDM}),
	}

	result, err2 := udm_consumer.SendNFIntances(self.NrfUri, models.NfType_UDM, models.NfType_UDM, param)
	if err2 != nil {
		t.Error(err2.Error())
	} else if len(result.NfInstances) > 0 && !reflect.DeepEqual(nfprofile, result.NfInstances[0]) {
		t.Error("failed for expecte value mismatch")
	} else if len(result.NfInstances) == 0 {
		t.Error("len(result.NfInstances) is 0\n")
	}

}

func TestSendSearchNFInstancesUDR(t *testing.T) {

	nrfInit()

	time.Sleep(200 * time.Millisecond)
	MongoDBLibrary.RestfulAPIDeleteMany("Nfprofile", bson.M{})
	time.Sleep(200 * time.Millisecond)

	udm_context.TestInit()
	self := udm_context.UDM_Self()
	nfprofile, err := udm_consumer.BuildNFInstance(self)
	if err != nil {
		t.Error(err.Error())
	}

	uri, _, err1 := udm_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, nfprofile)
	if err1 != nil {
		t.Error(err1.Error())
	} else {
		fmt.Println("uri: ", uri)
	}

	udmUeContext := udm_context.UdmUe_self()
	result := udm_consumer.SendNFIntancesUDR(udmUeContext.Supi, udm_consumer.NFDiscoveryToUDRParamSupi)
	fmt.Println("result: ", result)
}
