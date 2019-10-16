package nrf_consumer_test

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/Management"
	"free5gc/src/nrf/nrf_context"
	"free5gc/src/nrf/nrf_service"
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

func TestNRFRegisterSearchNFInstances(t *testing.T) {

	nrfInit()

	time.Sleep(200 * time.Millisecond)
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})

	time.Sleep(100 * time.Millisecond)

	nfInstanceId := nrf_context.Nrf_NfInstanceID
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	//Create nfProfile
	var nfProfile = models.NfProfile{
		NfInstanceId: nfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "NRF",
		RecoveryTime: &dateFormat,
	}

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	// Check test data (Use RESTful GET)
	_, _, err := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, nfProfile)
	if err != nil {
		t.Error(err.Error())
	}

	// Set client and set url
	configuration1 := Nnrf_NFDiscovery.NewConfiguration()
	configuration1.SetBasePath("https://127.0.0.1:29510")
	client1 := Nnrf_NFDiscovery.NewAPIClient(configuration1)

	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		TargetNfInstanceId: optional.NewInterface(nfInstanceId),
	}
	_, _, err1 := client1.NFInstancesStoreApi.SearchNFInstances(context.TODO(), "NRF", "NRF", &localVarOptionals)
	if err != nil {
		t.Error(err1.Error())
	}
}

func TestNRFRegisterSearchNFInstancesExtend(t *testing.T) {

	nrfInit()

	time.Sleep(200 * time.Millisecond)
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})

	time.Sleep(100 * time.Millisecond)

	//Create nfProfile
	var nfProfile models.NfProfile
	//nfProfile.RecoveryTime = &dateFormat

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	PutTestNfProfile(client, t)

	nfProfile = nrf_context.NrfNfProfile
	nfProfile.NrfInfo = Management.GetNrfInfo()

	// Check test data (Use RESTful PUT)
	_, _, err := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfProfile.NfInstanceId, nfProfile)
	if err != nil {
		t.Error(err.Error())
	}
	// Set client and set url
	configuration1 := Nnrf_NFDiscovery.NewConfiguration()
	configuration1.SetBasePath("https://127.0.0.1:29510")
	client1 := Nnrf_NFDiscovery.NewAPIClient(configuration1)

	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		TargetNfInstanceId: optional.NewInterface(nrf_context.NrfNfProfile.NfInstanceId),
	}
	_, _, err1 := client1.NFInstancesStoreApi.SearchNFInstances(context.TODO(), "NRF", "NRF", &localVarOptionals)
	if err != nil {
		t.Error(err1.Error())
	}
}

func PutTestNfProfile(client *Nnrf_NFManagement.APIClient, t *testing.T) {
	// Input Other Services (Use RESTful PUT)

	//set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	//UdrInfo
	UdrInfo := models.UdrInfo{
		GroupId: "1",
	}
	UdrInfoNfInstanceId := uuid.New().String()

	var UdrInfonfProfile = models.NfProfile{
		NfInstanceId: UdrInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "UDR",
		RecoveryTime: &dateFormat,
		UdrInfo:      &UdrInfo,
	}
	_, _, Udrerr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), UdrInfoNfInstanceId, UdrInfonfProfile)
	if Udrerr != nil {
		t.Error(Udrerr.Error())
	}
	//UdmInfo
	UdmInfo := models.UdmInfo{
		GroupId: "2",
	}
	UdmInfoNfInstanceId := uuid.New().String()

	var UdmInfonfProfile = models.NfProfile{
		NfInstanceId: UdmInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "UDM",
		RecoveryTime: &dateFormat,
		UdmInfo:      &UdmInfo,
	}
	_, _, Udmerr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), UdmInfoNfInstanceId, UdmInfonfProfile)
	if Udmerr != nil {
		t.Error(Udmerr.Error())
	}
	//AusfInfo
	AusfInfo := models.AusfInfo{
		GroupId: "3",
	}
	AusfInfoNfInstanceId := uuid.New().String()

	var AusfInfonfProfile = models.NfProfile{
		NfInstanceId: AusfInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "AUSF",
		RecoveryTime: &dateFormat,
		AusfInfo:     &AusfInfo,
	}
	_, _, Ausferr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), AusfInfoNfInstanceId, AusfInfonfProfile)
	if Ausferr != nil {
		t.Error(Ausferr.Error())
	}
	//AmfInfo
	AmfInfo := models.AmfInfo{
		AmfSetId: "44",
	}
	AmfInfoNfInstanceId := uuid.New().String()

	var AmfInfonfProfile = models.NfProfile{
		NfInstanceId: AmfInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "AMF",
		RecoveryTime: &dateFormat,
		AmfInfo:      &AmfInfo,
	}
	_, _, Amferr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), AmfInfoNfInstanceId, AmfInfonfProfile)
	if Amferr != nil {
		t.Error(Amferr.Error())
	}
	//AmfInfo2
	AmfInfo2 := models.AmfInfo{
		AmfSetId: "56",
	}
	AmfInfoNfInstanceId2 := uuid.New().String()

	var AmfInfonfProfile2 = models.NfProfile{
		NfInstanceId: AmfInfoNfInstanceId2,
		NfStatus:     "REGISTERED",
		NfType:       "AMF",
		//RecoveryTime: &dateFormat,
		AmfInfo: &AmfInfo2,
	}
	_, _, Amferr2 := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), AmfInfoNfInstanceId2, AmfInfonfProfile2)
	if Amferr2 != nil {
		t.Error(Amferr2.Error())
	}
	//AmfInfo3
	AmfInfo3 := models.AmfInfo{
		AmfSetId: "56",
	}
	AmfInfoNfInstanceId3 := uuid.New().String()

	var AmfInfonfProfile3 = models.NfProfile{
		NfInstanceId: AmfInfoNfInstanceId3,
		NfStatus:     "REGISTERED",
		NfType:       "AMF",
		//RecoveryTime: &dateFormat,
		AmfInfo: &AmfInfo3,
	}
	_, _, Amferr3 := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), AmfInfoNfInstanceId3, AmfInfonfProfile3)
	if Amferr3 != nil {
		t.Error(Amferr3.Error())
	}
	//SmfInfo
	SmfInfo := models.SmfInfo{
		AccessType: []models.AccessType{
			models.AccessType__3_GPP_ACCESS,
		},
	}
	SmfInfoNfInstanceId := uuid.New().String()

	var SmfInfonfProfile = models.NfProfile{
		NfInstanceId: SmfInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "SMF",
		RecoveryTime: &dateFormat,
		SmfInfo:      &SmfInfo,
	}
	_, _, Smferr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), SmfInfoNfInstanceId, SmfInfonfProfile)
	if Smferr != nil {
		t.Error(Smferr.Error())
	}
	//UpfInfo
	UpfInfo := models.UpfInfo{
		SmfServingArea: []string{
			"1", "2",
		},
	}
	UpfInfoNfInstanceId := uuid.New().String()

	var UpfInfonfProfile = models.NfProfile{
		NfInstanceId: UpfInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "UPF",
		RecoveryTime: &dateFormat,
		UpfInfo:      &UpfInfo,
	}
	_, _, Upferr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), UpfInfoNfInstanceId, UpfInfonfProfile)
	if Upferr != nil {
		t.Error(Upferr.Error())
	}
	//PcfInfo
	PcfInfo := models.PcfInfo{
		DnnList: []string{
			"1", "2",
		},
	}
	PcfInfoNfInstanceId := uuid.New().String()

	var PcfInfonfProfile = models.NfProfile{
		NfInstanceId: PcfInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "PCF",
		RecoveryTime: &dateFormat,
		PcfInfo:      &PcfInfo,
	}
	_, _, Pcferr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), PcfInfoNfInstanceId, PcfInfonfProfile)
	if Pcferr != nil {
		t.Error(Pcferr.Error())
	}
	//BsfInfo
	BsfInfo := models.BsfInfo{
		DnnList: []string{
			"1", "2",
		},
	}
	BsfInfoNfInstanceId := uuid.New().String()

	var BsfInfonfProfile = models.NfProfile{
		NfInstanceId: BsfInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "BSF",
		RecoveryTime: &dateFormat,
		BsfInfo:      &BsfInfo,
	}
	_, _, Bsferr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), BsfInfoNfInstanceId, BsfInfonfProfile)
	if Bsferr != nil {
		t.Error(Bsferr.Error())
	}
	//ChfInfo
	ChfInfo := models.ChfInfo{
		SupiRangeList: &[]models.SupiRange{
			{
				Start:   "1",
				End:     "2",
				Pattern: "3",
			},
		},
	}
	ChfInfoNfInstanceId := uuid.New().String()

	var ChfInfonfProfile = models.NfProfile{
		NfInstanceId: ChfInfoNfInstanceId,
		NfStatus:     "REGISTERED",
		NfType:       "CHF",
		RecoveryTime: &dateFormat,
		ChfInfo:      &ChfInfo,
	}
	_, _, Chferr := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), ChfInfoNfInstanceId, ChfInfonfProfile)
	if Chferr != nil {
		t.Error(Chferr.Error())
	}
}
