package pcf_consumer_test

import (
	"flag"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/pcf/factory"
	"free5gc/src/pcf/pcf_consumer"
	"free5gc/src/pcf/pcf_service"
	"testing"
	"time"

	"github.com/antihax/optional"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/bson"
)

var flags flag.FlagSet
var c = cli.NewContext(nil, &flags, nil)

func nrfInit() {
	nrf := &nrf_service.NRF{}
	nrf.Initialize(c)
	go nrf.Start()
	time.Sleep(100 * time.Millisecond)
}
func TestRegisterNFInstance(t *testing.T) {
	// init NRF
	nrfInit()
	// Clear DB
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})
	time.Sleep(50 * time.Millisecond)
	// Init PCF and register to NRF
	pcf := pcf_service.PCF{}
	pcf.Initialize(c)
	go pcf.Start()
	time.Sleep(50 * time.Millisecond)
	// Send NF Discovery to discover PCF
	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NPCF_AM_POLICY_CONTROL, models.ServiceName_NPCF_BDTPOLICYCONTROL, models.ServiceName_NPCF_POLICYAUTHORIZATION, models.ServiceName_NPCF_SMPOLICYCONTROL, models.ServiceName_NPCF_UE_POLICY_CONTROL}),
	}
	result, err := pcf_consumer.SendSearchNFInstances(factory.PcfConfig.Configuration.NrfUri, models.NfType_PCF, models.NfType_UDR, param)
	if err != nil {
		t.Error(err.Error())
	} else if result.NfInstances == nil {
		t.Error("NF Instances is nil")
	}
}
