package amf_consumer_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bronze1man/radius"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/AUSF/TestUEAuth"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/ausf/ausf_context"
	"free5gc/src/ausf/ausf_producer"
	"free5gc/src/ausf/ausf_service"
	"testing"
	"time"
)

var collName = "subscriptionData.authenticationData.authenticationSubscription"
var filter = bson.M{"ueId": "imsi-2089300007487"}

var testServingNetworkName = "5G:mnc093.mcc208.3gppnetwork.org"

func getAuthSubscription() (authSubs models.AuthenticationSubscription) {
	authSubs.PermanentKey = &models.PermanentKey{
		PermanentKeyValue: TestGenAuthData.MilenageTestSet19.K,
	}
	authSubs.Opc = &models.Opc{
		OpcValue: TestGenAuthData.MilenageTestSet19.OPC,
	}
	authSubs.Milenage = &models.Milenage{
		Op: &models.Op{
			OpValue: TestGenAuthData.MilenageTestSet19.OP,
		},
	}
	authSubs.AuthenticationManagementField = "8000"

	authSubs.SequenceNumber = TestGenAuthData.MilenageTestSet19.SQN
	authSubs.AuthenticationMethod = models.AuthMethod__5_G_AKA
	return
}
func toBsonM(data interface{}) bson.M {
	tmp, _ := json.Marshal(data)
	var putData = bson.M{}
	_ = json.Unmarshal(tmp, &putData)
	return putData
}
func insertAuthSubscriptionToMongoDB(ueId string, authSubs models.AuthenticationSubscription) {
	putData := toBsonM(authSubs)
	putData["ueId"] = ueId
	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func ausfInit() {
	ausf := &ausf_service.AUSF{}
	ausf.Initialize(testC)
	go ausf.Start()
	time.Sleep(100 * time.Millisecond)
}

func ausfUeAdd(ausfUe *ausf_context.AusfUeContext) {
	ausf_context.AddAusfUeContextToPool(ausfUe)
}

func Init() {
	if len(TestAmf.TestAmf.UePool) > 0 {
		insertAuthSubscriptionToMongoDB("imsi-2089300007487", getAuthSubscription())
		return
	}
	udminit()
	udrinit()
	ausfInit()

	time.Sleep(100 * time.Millisecond)

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
}

func TestUeAuthenticationAuthenticateRequest(t *testing.T) {
	Init()
	defer MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	// ausfUri := "https://localhost:29509"

	response, _, err := amf_consumer.SendUEAuthenticationAuthenticateRequest(ue, nil)
	if response != nil {
		fmt.Printf("response: %+v\n", response)
	} else if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		// TODO: error handling
	}
}

func TestEapConfirm(t *testing.T) {
	Init()
	defer MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
	ausfUe := ausf_context.NewAusfUeContext("imsi-2089300007487")
	ausfUe.ServingNetworkName = testServingNetworkName
	ausfUe.XRES = TestUEAuth.TestUeEapAuthTable[TestUEAuth.SUCCESS_CASE].Xres
	ausfUe.K_aut = TestUEAuth.TestUeEapAuthTable[TestUEAuth.SUCCESS_CASE].K_aut
	ausfUe.AuthStatus = models.AuthResult_ONGOING
	ausfUeAdd(ausfUe)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	// ausfUri := "https://localhost:29509"

	var eapPkt radius.EapPacket

	eapPkt.Code = radius.EapCode(radius.EapCodeResponse)
	eapPkt.Type = radius.EapType(50) // accroding to RFC5448 6.1
	eapPkt.Identifier = 0x01
	atRes, _ := ausf_producer.EapEncodeAttribute("AT_RES", TestUEAuth.TestUeEapAuthTable[TestUEAuth.SUCCESS_CASE].Res)
	atMAC, _ := ausf_producer.EapEncodeAttribute("AT_MAC", "")

	dataArrayBeforeMAC := atRes + atMAC
	eapPkt.Data = []byte(dataArrayBeforeMAC)
	encodedPktBeforeMAC := eapPkt.Encode()

	MACvalue := ausf_producer.CalculateAtMAC([]byte(TestUEAuth.TestUeEapAuthTable[TestUEAuth.SUCCESS_CASE].K_aut), encodedPktBeforeMAC)

	atMacNum := fmt.Sprintf("%02x", ausf_context.AT_MAC_ATTRIBUTE)
	atMACfirstRow, _ := hex.DecodeString(atMacNum + "05" + "0000")
	wholeAtMAC := append(atMACfirstRow, MACvalue...)

	atMAC = string(wholeAtMAC)
	dataArrayAfterMAC := atRes + atMAC

	eapPkt.Data = []byte(dataArrayAfterMAC)
	encodedPktAfterMAC := eapPkt.Encode()

	eapMsg := nasType.EAPMessage{
		Len:    uint16(len(encodedPktAfterMAC)),
		Buffer: encodedPktAfterMAC,
	}

	response, problemDetails, err := amf_consumer.SendEapAuthConfirmRequest(ue, eapMsg)
	if err != nil {
		t.Error(err)
	} else if problemDetails != nil {
		t.Logf("problemDetail: %+v", problemDetails)
	} else {
		t.Logf("response: %+v", response)
	}
}

func Test5gAkaConfirm(t *testing.T) {
	Init()
	ausfUe := ausf_context.NewAusfUeContext("imsi-2089300007487")
	ausfUe.ServingNetworkName = testServingNetworkName
	ausfUe.XresStar = TestUEAuth.TestUe5gAuthTable[TestUEAuth.SUCCESS_CASE].XresStar
	ausfUe.AuthStatus = models.AuthResult_ONGOING
	defer MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
	ausfUeAdd(ausfUe)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	// ausfUri := "https://localhost:29509"

	resStar := TestUEAuth.TestUe5gAuthTable[TestUEAuth.SUCCESS_CASE].ResStar
	response, problemDetails, err := amf_consumer.SendAuth5gAkaConfirmRequest(ue, resStar)
	if err != nil {
		t.Error(err)
	} else if problemDetails != nil {
		t.Logf("problemDetail: %+v", problemDetails)
	} else {
		t.Logf("response: %+v", response)
	}
}
