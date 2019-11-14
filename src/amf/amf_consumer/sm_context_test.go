package amf_consumer_test

import (
	"encoding/hex"
	"flag"
	"fmt"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/SMF/TestPDUSession"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_context"
	"free5gc/src/smf/smf_service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
	"github.com/urfave/cli"
)

func smfInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	smf := &smf_service.SMF{}
	smf.Initialize(c)
	go smf.Start()
	time.Sleep(100 * time.Millisecond)
}

func TestSMContextCreate(t *testing.T) {
	smfInit()
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	payload := TestPDUSession.GetEstablishmentRequestData(TestPDUSession.SERVICE_REQUEST)
	pduSession := models.PduSessionContext{
		PduSessionId: 10,
		Dnn:          "nctu.edu.tw",
		SNssai: &models.Snssai{
			Sst: 1,
			Sd:  "020304",
		},
	}
	requestType := models.RequestType_INITIAL_REQUEST
	if anType := ue.GetAnType(); anType == "" {
		pduSession.AccessType = models.AccessType__3_GPP_ACCESS
	} else {
		pduSession.AccessType = anType
	}
	smContextCreateData := amf_consumer.BuildCreateSmContextRequest(ue, pduSession, requestType)
	// TODO: http://localhost:29502/ -> smfD smfUri which required from NRF
	smfUri := "https://localhost:29502"

	response, smContextRef, _, _, err := amf_consumer.SendCreateSmContextRequest(ue, smfUri, payload, smContextCreateData)
	if response != nil {
		var smContext amf_context.SmContext
		pduSession.SmContextRef = smContextRef
		smContext.PduSessionContext = &pduSession
		smContext.UserLocation = deepcopy.Copy(ue.Location).(models.UserLocation)
		smContext.SmfUri = smfUri
		ue.SmContextList[pduSession.PduSessionId] = &smContext
		// TODO: handle response(response N2SmInfo to RAN if exists)
		TestAmf.Config.Dump(*ue.SmContextList[10])

	} else if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		// TODO: error handling
	}

}

func TestSMContextUpdate(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestSMContextCreate(t)
	}

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	/***ActivateUpCnxState***/
	ue.Location.NrLocation.Tai.Tac = "0001"
	ue.Location.NrLocation.Ncgi.NrCellId = "000000002"
	anType := ue.GetAnType()
	if anType == "" {
		anType = models.AccessType__3_GPP_ACCESS
	}
	response, _, _, err := amf_consumer.SendUpdateSmContextActivateUpCnxState(ue, 10, anType)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS ActivateUpCnxState--")
	}
	if response != nil {
		//TODO: send n2InfoToRAN
		ue.SmContextList[10].UserLocation = deepcopy.Copy(ue.Location).(models.UserLocation)
	} else {
		// TODO: error handling
	}

	/***DeactivateUpCnxState***/
	cause := models.Cause__5_G_AN_NOT_RESPONDING
	causes := amf_context.CauseAll{
		Cause: &cause,
	}
	response, _, _, err = amf_consumer.SendUpdateSmContextDeactivateUpCnxState(ue, 10, causes)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS DeactivateUpCnxState--")
	}
	if response == nil {
		// TODO: error handling
	}

	/***ChangeAccessType***/
	response, _, _, err = amf_consumer.SendUpdateSmContextChangeAccessType(ue, 10, false)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS ChangeAccessType--")
	}
	if response != nil && response.BinaryDataN2SmInformation != nil {
		// TODO: n2SmInfo send to RAN
	} else if response == nil {
		// TODO: error handling
	}

	/***N2Info***/
	hexString := "00150035000004001B000800535455104546470052400903006672656535474300660010000053545500535455000012985354550015400140"
	n2SmInfo, err := hex.DecodeString(hexString)
	response, _, _, err = amf_consumer.SendUpdateSmContextN2Info(ue, 10, models.N2SmInfoType_PDU_RES_SETUP_RSP, n2SmInfo)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS N2Info--")
	}
	// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
	if response != nil && response.BinaryDataN2SmInformation != nil {
		// TODO: n2SmInfo send to RAN
	} else if response == nil {
		// TODO: error handling
	}

	/***XnHandover***/
	response, _, _, err = amf_consumer.SendUpdateSmContextXnHandover(ue, 10, models.N2SmInfoType_PATH_SWITCH_REQ, n2SmInfo)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS XnHandover--")
	}
	if response == nil {
		// TODO: error handling
	}

	/***XnHandoverFailed***/
	response, _, _, err = amf_consumer.SendUpdateSmContextXnHandoverFailed(ue, 10, models.N2SmInfoType_PATH_SWITCH_REQ_FAIL, n2SmInfo)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS XnHandoverFailed--")
	}
	if response == nil {
		// TODO: error handling
	}

	/***N2HandoverPreparing***/
	handoverAmfId := uuid.New().String()
	if handoverAmfId == TestAmf.TestAmf.NfId {
		handoverAmfId = ""
	}
	targetId := &models.NgRanTargetId{
		RanNodeId: &models.GlobalRanNodeId{
			PlmnId: &models.PlmnId{
				Mcc: "208",
				Mnc: "93",
			},
			GNbId: &models.GNbId{
				BitLength: 24,
				GNBValue:  "454647",
			},
		},
		Tai: &models.Tai{
			PlmnId: ue.Location.NrLocation.Tai.PlmnId,
			Tac:    "000001",
		},
	}
	response, _, _, err = amf_consumer.SendUpdateSmContextN2HandoverPreparing(ue, 10, models.N2SmInfoType_HANDOVER_REQUIRED, n2SmInfo, handoverAmfId, targetId)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS N2HandoverPreparing--")
	}
	if response != nil && response.BinaryDataN2SmInformation != nil {
		// TODO: n2SmInfo send to RAN
	} else if response == nil {
		// TODO: error handling
	}

	/***N2HandoverPrepared***/
	response, _, _, err = amf_consumer.SendUpdateSmContextN2HandoverPrepared(ue, 10, models.N2SmInfoType_HANDOVER_REQ_ACK, n2SmInfo)
	// failure type: models.N2SmInfoType_HANDOVER_RES_ALLOC_FAIL
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS N2HandoverPrepared--")
	}
	if response != nil && response.BinaryDataN2SmInformation != nil {
		// TODO: n2SmInfo send to RAN
	} else if response == nil {
		// TODO: error handling
	}

	/***N2HandoverComplete***/
	handoverAmfGuami := &models.Guami{
		PlmnId: TestAmf.TestAmf.ServedGuamiList[0].PlmnId,
		AmfId:  "CAFE01",
	}
	response, _, _, err = amf_consumer.SendUpdateSmContextN2HandoverComplete(ue, 10, handoverAmfId, handoverAmfGuami)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS N2HandoverComplete--")
	}
	if response == nil {
		// TODO: error handling
	}

	/***N2HandoverCanceled***/
	cause = models.Cause_HO_FAILURE
	causes = amf_context.CauseAll{
		Cause: &cause,
	}
	response, _, _, err = amf_consumer.SendUpdateSmContextN2HandoverCanceled(ue, 10, causes)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS N2HandoverCanceled--")
	}
	if response == nil {
		// TODO: error handling
	}

	/***HandoverBetweenAccessType***/
	accessType := models.AccessType_NON_3_GPP_ACCESS
	response, _, _, err = amf_consumer.SendUpdateSmContextHandoverBetweenAccessType(ue, 10, accessType, nil)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS HandoverBetweenAccessType--")
	}
	if response != nil {
		ue.SmContextList[10].PduSessionContext.AccessType = accessType
	} else {
		// TODO: error handling
	}

	/***HandoverBetweenAMF***/
	handoverBetweenAmfId := uuid.New().String()
	handoverBetweenAmfGuami := &models.Guami{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "94",
		},
		AmfId: "010203",
	}
	activation := false
	response, _, _, err = amf_consumer.SendUpdateSmContextHandoverBetweenAMF(ue, 10, handoverBetweenAmfId, handoverBetweenAmfGuami, activation)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		fmt.Println("--PASS HandoverBetweenAMF--")
	}
	if response != nil && response.BinaryDataN2SmInformation != nil {
		// TODO: n2SmInfo send to RAN
	} else if response == nil {
		// TODO: error handling
	}

}

func TestSMContextRelease(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestSMContextCreate(t)
	}
	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	releaseDate := amf_consumer.BuildReleaseSmContextRequest(ue, nil, "", nil)
	_, err := amf_consumer.SendReleaseSmContextRequest(ue, 10, releaseDate)
	if err != nil {
		t.Errorf("[ERROR] " + err.Error())
	} else {
		TestAmf.Config.Dump(ue.SmContextList)
	}

}
