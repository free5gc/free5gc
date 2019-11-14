package amf_util_test

import (
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_util"
	"free5gc/src/amf/factory"
	"reflect"
	"testing"
)

var (
	ConfigFileFromArgs string
	DefaultConfigFile  string = path_util.Gofree5gcPath("free5gc/src/amf/amf_util/test/testAmfcfg.conf")
)

func compareContext(t *testing.T, context, testContext *amf_context.AMFContext) {
	if testContext.Name != context.Name {
		t.Errorf("Expect : %s\nOuput: %s", context.Name, testContext.Name)
	}
	if testContext.UriScheme != context.UriScheme {
		t.Errorf("Expect : %s\nOuput: %s", context.UriScheme, testContext.UriScheme)
	}
	if testContext.HttpIPv4Address != context.HttpIPv4Address {
		t.Errorf("Expect : %s\nOuput: %s", context.HttpIPv4Address, testContext.HttpIPv4Address)
	}
	if testContext.HttpIpv4Port != context.HttpIpv4Port {
		t.Errorf("Expect : %d\nOuput: %d", context.HttpIpv4Port, testContext.HttpIpv4Port)
	}
	if testContext.NrfUri != context.NrfUri {
		t.Errorf("Expect : %s\nOuput: %s", context.NrfUri, testContext.NrfUri)
	}
	if !reflect.DeepEqual(testContext.ServedGuamiList, context.ServedGuamiList) {
		t.Errorf("Expect : %s\nOuput: %s", TestAmf.Config.Sdump(context.ServedGuamiList), TestAmf.Config.Sdump(testContext.ServedGuamiList))
	}
	if !reflect.DeepEqual(testContext.NfService, context.NfService) {
		t.Errorf("Expect : %s\nOuput: %s", TestAmf.Config.Sdump(context.NfService), TestAmf.Config.Sdump(testContext.NfService))
	}
	if !reflect.DeepEqual(testContext.PlmnSupportList, context.PlmnSupportList) {
		t.Errorf("Expect : %s\nOuput: %s", TestAmf.Config.Sdump(context.PlmnSupportList), TestAmf.Config.Sdump(testContext.PlmnSupportList))
	}
	if !reflect.DeepEqual(testContext.SupportDnnLists, context.SupportDnnLists) {
		t.Errorf("Expect : %s\nOuput: %s", TestAmf.Config.Sdump(context.SupportDnnLists), TestAmf.Config.Sdump(testContext.SupportDnnLists))
	}
	if !reflect.DeepEqual(testContext.SupportTaiLists, context.SupportTaiLists) {
		t.Errorf("Expect : %s\nOuput: %s", TestAmf.Config.Sdump(context.SupportTaiLists), TestAmf.Config.Sdump(testContext.SupportTaiLists))
	}
	if !reflect.DeepEqual(testContext.SecurityAlgorithm, context.SecurityAlgorithm) {
		t.Errorf("Expect : %s\nOuput: %s", TestAmf.Config.Sdump(context.SecurityAlgorithm), TestAmf.Config.Sdump(testContext.SecurityAlgorithm))
	}
	if !reflect.DeepEqual(testContext.NetworkName, context.NetworkName) {
		t.Errorf("Expect : %s\nOuput: %s", TestAmf.Config.Sdump(context.NetworkName), TestAmf.Config.Sdump(testContext.NetworkName))
	}
}

func TestInitAmfContext1(t *testing.T) {

	configFile := path_util.Gofree5gcPath("free5gc/src/amf/amf_util/test/testAmfcfg.conf")
	context := getExpAmf1()
	testContext := amf_context.AMF_Self()
	factory.InitConfigFactory(configFile)
	amf_util.InitAmfContext(testContext)
	compareContext(t, context, testContext)
}

func TestInitAmfContext2(t *testing.T) {
	amf_context.AMF_Self().Reset()
	configFile := path_util.Gofree5gcPath("free5gc/src/amf/amf_util/test/testAmfcfg2.conf")
	context := getExpAmf2()
	testContext := amf_context.AMF_Self()
	factory.InitConfigFactory(configFile)
	amf_util.InitAmfContext(testContext)
	compareContext(t, context, testContext)
}

var dnn = []string{
	"internet",
	"wire.cs.nctu.edu.tw",
}
var serivceName = []string{
	"namf-comm",
	"namf-evts",
	"namf-mt",
	"namf-loc",
}
var plmnId = []models.PlmnId{
	{
		Mcc: "208",
		Mnc: "93",
	},
	{
		Mcc: "466",
		Mnc: "46",
	},
}

var tai = []models.Tai{
	{
		PlmnId: &plmnId[0],
		Tac:    "000001",
	},
	{
		PlmnId: &plmnId[0],
		Tac:    "000102",
	},
	{
		PlmnId: &plmnId[1],
		Tac:    "000201",
	},
}
var guami = []models.Guami{
	{
		PlmnId: &plmnId[0],
		AmfId:  "cafe00",
	},
	{
		PlmnId: &plmnId[1],
		AmfId:  "123456",
	},
}

var sNssaiList = []models.Snssai{
	{
		Sst: 1,
		Sd:  "010203",
	},
	{
		Sst: 2,
		Sd:  "112233",
	},
	{
		Sst: 2,
		Sd:  "445566",
	},
}

func initTestAmfContext() (context *amf_context.AMFContext) {
	context = new(amf_context.AMFContext)
	context.UriScheme = models.UriScheme_HTTPS
	context.ServedGuamiList = make([]models.Guami, 0, amf_context.MaxNumOfServedGuamiList)
	context.PlmnSupportList = make([]amf_context.PlmnSupportItem, 0, amf_context.MaxNumOfPLMNs)
	context.NfService = make(map[models.ServiceName]models.NfService)
	return
}

func getExpAmf1() (context *amf_context.AMFContext) {
	context = initTestAmfContext()
	context.Name = "AMF"
	context.NgapIpList = []string{"127.0.0.1"}
	context.UriScheme = models.UriScheme_HTTPS
	context.HttpIPv4Address = "127.0.0.1"
	context.HttpIPv6Address = "2001:0db8:85a3:08d3:1319:8a2e:0370:7344"
	context.HttpIpv4Port = 29518
	context.NrfUri = "https://localhost:29510"
	context.InitNFService(serivceName, "1.0.0")
	context.SecurityAlgorithm.IntegrityOrder = []uint8{0x40}
	context.SecurityAlgorithm.CipheringOrder = []uint8{0x40}
	context.NetworkName.Full = "free5GC"
	context.NetworkName.Short = "free"
	context.SupportDnnLists = append(context.SupportDnnLists, dnn[0])
	context.SupportTaiLists = append(context.SupportTaiLists, tai[0])
	context.ServedGuamiList = append(context.ServedGuamiList, guami[0])
	plmnSupportItem := amf_context.NewPlmnSupportItem()
	plmnSupportItem.PlmnId = plmnId[0]
	plmnSupportItem.SNssaiList = append(plmnSupportItem.SNssaiList, sNssaiList[:1]...)
	context.PlmnSupportList = append(context.PlmnSupportList, plmnSupportItem)
	return
}

func getExpAmf2() (context *amf_context.AMFContext) {
	context = initTestAmfContext()
	context.Name = "Wirelab"
	context.NgapIpList = []string{"127.0.0.1", "192.188.2.2"}
	context.UriScheme = models.UriScheme_HTTP
	context.HttpIPv4Address = "192.168.0.1"
	context.HttpIpv4Port = 8888
	context.NrfUri = "https://192.168.0.2:29510"
	context.InitNFService(serivceName[:2], "2.0.0")
	context.SecurityAlgorithm.IntegrityOrder = []uint8{0x40, 0x80}
	context.SecurityAlgorithm.CipheringOrder = []uint8{0x40, 0x20, 0x08}
	context.NetworkName.Full = "HAHAHAHA"
	context.SupportDnnLists = append(context.SupportDnnLists, dnn...)
	context.SupportTaiLists = append(context.SupportTaiLists, tai...)
	context.ServedGuamiList = append(context.ServedGuamiList, guami...)
	plmnSupportItem := amf_context.NewPlmnSupportItem()
	plmnSupportItem.PlmnId = plmnId[0]
	plmnSupportItem.SNssaiList = append(plmnSupportItem.SNssaiList, sNssaiList[:2]...)
	context.PlmnSupportList = append(context.PlmnSupportList, plmnSupportItem)
	plmnSupportItem = amf_context.NewPlmnSupportItem()
	plmnSupportItem.PlmnId = plmnId[1]
	plmnSupportItem.SNssaiList = append(plmnSupportItem.SNssaiList, sNssaiList[2])
	context.PlmnSupportList = append(context.PlmnSupportList, plmnSupportItem)
	return
}
