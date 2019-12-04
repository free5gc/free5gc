package ngap_message_test

import (
	"encoding/hex"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/AMF/TestComm"
	"free5gc/lib/aper"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/ngap"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_ngap/ngap_message"
	"free5gc/src/test/ngapTestpacket"
	"reflect"
	"testing"
)

func TestBuildPDUSessionResourceReleaseCommand(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	transfer := TestComm.GetPDUSessionResourceReleaseCommandTransfer()
	pduSessionResourceToReleaseListRelCmd := ngapType.PDUSessionResourceToReleaseListRelCmd{}

	item := ngapType.PDUSessionResourceToReleaseItemRelCmd{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 10,
		},
		PDUSessionResourceReleaseCommandTransfer: transfer,
		// TODO: use real PDUSessionResourceSetupRequestTransfer
		// PDUSessionResourceSetupRequestTransfer: aper.OctetString("\x01\x02\x03"),
	}
	pduSessionResourceToReleaseListRelCmd.List = append(pduSessionResourceToReleaseListRelCmd.List, item)
	pkg, err := ngap_message.BuildPDUSessionResourceReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS], []byte{12}, pduSessionResourceToReleaseListRelCmd)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))
	}
}

func TestBuildNGSetupResponse(t *testing.T) {

	pdu, err := ngap_message.BuildNGSetupResponse()
	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildNGSetupFailure(t *testing.T) {

	if len(TestAmf.TestAmf.UePool) == 0 {
		TestAmf.AmfInit()
		TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	}

	cause := ngapType.Cause{
		Present: ngapType.CausePresentMisc,
		Misc: &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		},
	}
	pdu, err := ngap_message.BuildNGSetupFailure(cause)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("  Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildNGReset(t *testing.T) {

	cause := ngapType.Cause{
		Present: ngapType.CausePresentMisc,
		Misc: &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		},
	}

	pdu, err := ngap_message.BuildNGReset(cause, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildNGResetAcknowledge(t *testing.T) {

	AmfUeNgapID := ngapType.AMFUENGAPID{Value: 123}
	RanUeNgapID := ngapType.RANUENGAPID{Value: 456}

	item := ngapType.UEAssociatedLogicalNGConnectionItem{
		AMFUENGAPID: &AmfUeNgapID,
		RANUENGAPID: &RanUeNgapID,
	}

	partOfNGInterface := ngapType.UEAssociatedLogicalNGConnectionList{
		List: []ngapType.UEAssociatedLogicalNGConnectionItem{
			item,
		},
	}

	pdu, err := ngap_message.BuildNGResetAcknowledge(&partOfNGInterface, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildDownlinkNasTransport(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	nasPdu := ngapType.NASPDU{
		Value: aper.OctetString("\x01\x02"),
	}

	pdu, err := ngap_message.BuildDownlinkNasTransport(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasPdu.Value)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildUEContextReleaseCommand(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	pdu, err := ngap_message.BuildUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.CausePresentMisc, ngapType.CauseMiscPresentUnspecified)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildErrorIndication(t *testing.T) {

	amfUeNgapID := int64(123)
	ranUeNgapID := int64(456)

	cause := ngapType.Cause{
		Present: ngapType.CausePresentMisc,
		Misc: &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		},
	}

	pdu, err := ngap_message.BuildErrorIndication(&amfUeNgapID, &ranUeNgapID, &cause, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}
func TestBuildUERadioCapabilityCheckRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	pdu, err := ngap_message.BuildUERadioCapabilityCheckRequest(ue.RanUe[models.AccessType__3_GPP_ACCESS])

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildHandoverCancelAcknowledge(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	pdu, err := ngap_message.BuildHandoverCancelAcknowledge(ue.RanUe[models.AccessType__3_GPP_ACCESS], nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildPDUSessionResourceSetupRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	nasPdu := []byte{0x01, 0x02}
	pduSessionResourceSetupListSUReq := ngapType.PDUSessionResourceSetupListSUReq{}

	item := ngapType.PDUSessionResourceSetupItemSUReq{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 1,
		},
		SNSSAI: ngapType.SNSSAI{
			SST: ngapType.SST{
				Value: aper.OctetString("\x01"),
			},
		},
		PDUSessionResourceSetupRequestTransfer: aper.OctetString("\x01\x02\x03"),
	}

	pduSessionResourceSetupListSUReq.List = append(pduSessionResourceSetupListSUReq.List, item)

	pkg, err := ngap_message.BuildPDUSessionResourceSetupRequest(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasPdu, pduSessionResourceSetupListSUReq)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))
	}
}

func TestBuildPDUSessionResourceModifyConfirm(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 1
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 2
	pduSessionResourceModifyConfirmList := ngapType.PDUSessionResourceModifyListModCfm{}
	item := ngapType.PDUSessionResourceModifyItemModCfm{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 10,
		},
		PDUSessionResourceModifyConfirmTransfer: aper.OctetString("\x01\x02"),
	}
	pduSessionResourceModifyConfirmList.List = append(pduSessionResourceModifyConfirmList.List, item)

	pduSessionResourceFailedToModifyListModCfm := ngapType.PDUSessionResourceFailedToModifyListModCfm{}
	item2 := ngapType.PDUSessionResourceFailedToModifyItemModCfm{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 5,
		},
		PDUSessionResourceModifyIndicationUnsuccessfulTransfer: aper.OctetString("\x01\x02"),
	}
	pduSessionResourceFailedToModifyListModCfm.List = append(pduSessionResourceFailedToModifyListModCfm.List, item2)

	pkg, err := ngap_message.BuildPDUSessionResourceModifyConfirm(ue.RanUe[models.AccessType__3_GPP_ACCESS], pduSessionResourceModifyConfirmList, pduSessionResourceFailedToModifyListModCfm, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))
	}
}

func TestBuildPDUSessionResourceModifyRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	pduModifyRequestList := ngapType.PDUSessionResourceModifyListModReq{}
	item := ngapType.PDUSessionResourceModifyItemModReq{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 1,
		},
		NASPDU: &ngapType.NASPDU{
			Value: aper.OctetString("\x01\x02"),
		},
		PDUSessionResourceModifyRequestTransfer: aper.OctetString("\x01\x02\x03"),
	}

	pduModifyRequestList.List = append(pduModifyRequestList.List, item)

	pkg, err := ngap_message.BuildPDUSessionResourceModifyRequest(ue.RanUe[models.AccessType__3_GPP_ACCESS], pduModifyRequestList)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))
	}
}

func TestBuildInitialContextSetupRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.PlmnId = models.PlmnId{
		Mcc: "208",
		Mnc: "93",
	}
	ue.AccessAndMobilitySubscriptionData = &models.AccessAndMobilitySubscriptionData{}
	ue.AccessAndMobilitySubscriptionData.RatRestrictions = append(ue.AccessAndMobilitySubscriptionData.RatRestrictions, models.RatType_EUTRA)

	ue.AmPolicyAssociation = &models.PolicyAssociation{}
	ue.AmPolicyAssociation.ServAreaRes = &models.ServiceAreaRestriction{}
	ue.AmPolicyAssociation.ServAreaRes.Areas = append(ue.AmPolicyAssociation.ServAreaRes.Areas, models.Area{
		Tacs: []string{
			"000102",
		},
	})

	pkg, err := ngap_message.BuildInitialContextSetupRequest(ue, models.AccessType__3_GPP_ACCESS, []byte{0x01, 0x02}, nil, nil, nil, nil, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))
	}
}

func TestBuildUEContextModificationRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	oldAmfUeNgapID := int64(1234)

	emergencyFallbackIndicator := ngapType.EmergencyFallbackIndicator{}
	emergencyFallbackIndicator.EmergencyFallbackRequestIndicator.Value = ngapType.EmergencyFallbackRequestIndicatorPresentEmergencyFallbackRequested
	emergencyFallbackIndicator.EmergencyServiceTargetCN = new(ngapType.EmergencyServiceTargetCN)
	emergencyFallbackIndicator.EmergencyServiceTargetCN.Value = ngapType.EmergencyServiceTargetCNPresentFiveGC

	pkg, err := ngap_message.BuildUEContextModificationRequest(ue, models.AccessType__3_GPP_ACCESS, &oldAmfUeNgapID, nil, nil, nil, &emergencyFallbackIndicator)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))
	}
}

func TestBuildHandoverCommand(t *testing.T) {

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	pduSessionResourceHandoverList := ngapType.PDUSessionResourceHandoverList{}
	pduSessionResourceToReleaseList := ngapType.PDUSessionResourceToReleaseListHOCmd{}

	item1 := ngapType.PDUSessionResourceHandoverItem{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 10,
		},
		HandoverCommandTransfer: aper.OctetString("\x11\x22\x33"),
	}
	pduSessionResourceHandoverList.List = append(pduSessionResourceHandoverList.List, item1)

	item2 := ngapType.PDUSessionResourceToReleaseItemHOCmd{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 5,
		},
		HandoverPreparationUnsuccessfulTransfer: aper.OctetString("\x01\x02"),
	}
	pduSessionResourceToReleaseList.List = append(pduSessionResourceToReleaseList.List, item2)
	container := ngapType.TargetToSourceTransparentContainer{
		Value: []byte{0x11, 0x22},
	}
	pdu, err := ngap_message.BuildHandoverCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS], pduSessionResourceHandoverList, pduSessionResourceToReleaseList, container, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildHandoverPreparationFailure(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	cause := ngapType.Cause{
		Present: ngapType.CausePresentMisc,
		Misc: &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		},
	}

	pdu, err := ngap_message.BuildHandoverPreparationFailure(ue.RanUe[models.AccessType__3_GPP_ACCESS], cause, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildHandoverRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.NCC = 5
	ue.NH, _ = hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	ue.SecurityCapabilities.NREncryptionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.NRIntegrityProtectionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.EUTRAEncryptionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.EUTRAIntegrityProtectionAlgorithms = [2]byte{0xe0, 0x00}

	cause := ngapType.Cause{
		Present: ngapType.CausePresentMisc,
		Misc: &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		},
	}

	pduSessionResourceSetupListHOReq := ngapType.PDUSessionResourceSetupListHOReq{}
	setupItem := ngapType.PDUSessionResourceSetupItemHOReq{

		PDUSessionID: ngapType.PDUSessionID{
			Value: 10,
		},
		SNSSAI: ngapType.SNSSAI{
			SST: ngapType.SST{
				Value: aper.OctetString("\x11"),
			},
			// SD(optional)
		},
		HandoverRequestTransfer: aper.OctetString("\x11\x22\x33"),
	}

	pduSessionResourceSetupListHOReq.List = append(pduSessionResourceSetupListHOReq.List, setupItem)

	sourceToTargetTransparentContainer := ngapType.SourceToTargetTransparentContainer{
		Value: aper.OctetString("\x30\x33\x99"),
	}

	pdu, err := ngap_message.BuildHandoverRequest(ue.RanUe[models.AccessType__3_GPP_ACCESS], cause, pduSessionResourceSetupListHOReq, sourceToTargetTransparentContainer, false)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildPathSwitchRequestAcknowledge(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	pduSessionResourceSwitchedList := ngapType.PDUSessionResourceSwitchedList{}
	pduSessionResourceReleasedList := ngapType.PDUSessionResourceReleasedListPSAck{}

	switchedItem := ngapType.PDUSessionResourceSwitchedItem{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 10,
		},
		PathSwitchRequestAcknowledgeTransfer: aper.OctetString("\x11\x22\x33"),
	}
	pduSessionResourceSwitchedList.List = append(pduSessionResourceSwitchedList.List, switchedItem)

	releasedItem := ngapType.PDUSessionResourceReleasedItemPSAck{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 5,
		},
		PathSwitchRequestUnsuccessfulTransfer: aper.OctetString("\x22\x33\x44"),
	}
	pduSessionResourceReleasedList.List = append(pduSessionResourceReleasedList.List, releasedItem)

	ue.NCC = 5
	ue.NH, _ = hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	ue.SecurityCapabilities.NREncryptionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.NRIntegrityProtectionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.EUTRAEncryptionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.EUTRAIntegrityProtectionAlgorithms = [2]byte{0xe0, 0x00}

	pdu, err := ngap_message.BuildPathSwitchRequestAcknowledge(ue.RanUe[models.AccessType__3_GPP_ACCESS], pduSessionResourceSwitchedList, pduSessionResourceReleasedList, false, nil, nil, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildPathSwitchRequestFailure(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	pduSessionResourceReleasedList := ngapType.PDUSessionResourceReleasedListPSFail{}

	releasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{
		PDUSessionID: ngapType.PDUSessionID{
			Value: 5,
		},
		PathSwitchRequestUnsuccessfulTransfer: aper.OctetString("\x22\x33\x44"),
	}
	pduSessionResourceReleasedList.List = append(pduSessionResourceReleasedList.List, releasedItem)
	ranUe := ue.RanUe[models.AccessType__3_GPP_ACCESS]
	pdu, err := ngap_message.BuildPathSwitchRequestFailure(ranUe.AmfUeNgapId, ranUe.RanUeNgapId, &pduSessionResourceReleasedList, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildDownlinkRanStatusTransfer(t *testing.T) {
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	item := ngapType.DRBsSubjectToStatusTransferItem{
		DRBID: ngapType.DRBID{Value: int64(8)},
		DRBStatusUL: ngapType.DRBStatusUL{
			Present: ngapType.DRBStatusULPresentDRBStatusUL18,
			DRBStatusUL12: &ngapType.DRBStatusUL12{
				ULCOUNTValue: ngapType.COUNTValueForPDCPSN12{
					PDCPSN12:    int64(2),
					HFNPDCPSN12: int64(1),
				},
			},
			DRBStatusUL18: &ngapType.DRBStatusUL18{
				ULCOUNTValue: ngapType.COUNTValueForPDCPSN18{
					PDCPSN18:    int64(4),
					HFNPDCPSN18: int64(3),
				},
			},
		},
		DRBStatusDL: ngapType.DRBStatusDL{
			Present: ngapType.DRBStatusDLPresentDRBStatusDL12,
			DRBStatusDL12: &ngapType.DRBStatusDL12{
				DLCOUNTValue: ngapType.COUNTValueForPDCPSN12{
					PDCPSN12:    int64(2),
					HFNPDCPSN12: int64(1),
				},
			},
			DRBStatusDL18: &ngapType.DRBStatusDL18{
				DLCOUNTValue: ngapType.COUNTValueForPDCPSN18{
					PDCPSN18:    int64(4),
					HFNPDCPSN18: int64(3),
				},
			},
		},
	}
	ranStatusTransferTransparentContainer := ngapType.RANStatusTransferTransparentContainer{
		DRBsSubjectToStatusTransferList: ngapType.DRBsSubjectToStatusTransferList{
			List: []ngapType.DRBsSubjectToStatusTransferItem{
				item,
			},
		},
	}

	pdu, err := ngap_message.BuildDownlinkRanStatusTransfer(ue.RanUe[models.AccessType__3_GPP_ACCESS], ranStatusTransferTransparentContainer)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildPaging(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.Tai.PlmnId = &models.PlmnId{
		Mcc: "208",
		Mnc: "f93",
	}
	ue.Tai.Tac = "000001"
	ue.RegistrationArea[models.AccessType__3_GPP_ACCESS] = append(ue.RegistrationArea[models.AccessType__3_GPP_ACCESS], ue.Tai)

	tai2 := models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "000002",
	}
	ue.RegistrationArea[models.AccessType__3_GPP_ACCESS] = append(ue.RegistrationArea[models.AccessType__3_GPP_ACCESS], tai2)

	ue.Guti = "20893cafe0000000012"

	pagingPriority := ngapType.PagingPriority{
		Value: ngapType.PagingPriorityPresentPriolevel1,
	}

	recommendedCell := amf_context.RecommendedCell{
		NgRanCGI: amf_context.NGRANCGI{
			Present: amf_context.NgRanCgiPresentNRCGI,
			NRCGI: &models.Ncgi{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				NrCellId: "000000001",
			},
		},
	}
	ue.InfoOnRecommendedCellsAndRanNodesForPaging = new(amf_context.InfoOnRecommendedCellsAndRanNodesForPaging)
	ue.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedCells = append(ue.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedCells, recommendedCell)

	pkg, err := ngap_message.BuildPaging(ue, &pagingPriority, false)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))
	}
}

func TestBuildRerouteNasRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    12, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}
	nasPdu := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration, mobileIdentity5GS, nil, nil)
	initialUeMessage := ngapTestpacket.BuildInitialUEMessage(1, nasPdu, "")
	initialUeMessagePkg, _ := ngap.Encoder(initialUeMessage)

	ue.Guti = "20893cafe0000000001"
	pdu, err := ngap_message.BuildRerouteNasRequest(ue, models.AccessType__3_GPP_ACCESS, nil, initialUeMessagePkg, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildRanConfigurationUpdateAcknowledge(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	//ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	pdu, err := ngap_message.BuildRanConfigurationUpdateAcknowledge(nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildRanConfigurationUpdateFailure(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	cause := ngapType.Cause{
		Present: ngapType.CausePresentMisc,
		Misc: &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		},
	}

	pdu, err := ngap_message.BuildRanConfigurationUpdateFailure(cause, nil)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}
func TestBuildAMFStatusIndication(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	item := ngapType.UnavailableGUAMIItem{
		GUAMI: ngapType.GUAMI{
			PLMNIdentity: ngapType.PLMNIdentity{
				Value: aper.OctetString("\x02\xf8\x39"),
			},
			AMFRegionID: ngapType.AMFRegionID{
				Value: aper.BitString{
					Bytes:     []byte{0x45, 0x46},
					BitLength: 8,
				},
			},
			AMFSetID: ngapType.AMFSetID{
				Value: aper.BitString{
					Bytes:     []byte{0x45, 0x46},
					BitLength: 10,
				},
			},
			AMFPointer: ngapType.AMFPointer{
				Value: aper.BitString{
					Bytes:     []byte{0x45},
					BitLength: 6,
				},
			},
		},
	}
	unavailableGUAMIList := ngapType.UnavailableGUAMIList{
		List: []ngapType.UnavailableGUAMIItem{
			item,
		},
	}

	pdu, err := ngap_message.BuildAMFStatusIndication(unavailableGUAMIList)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildOverloadStart(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	overloadResponse := ngapType.OverloadResponse{
		Present: ngapType.OverloadResponsePresentOverloadAction,
		OverloadAction: &ngapType.OverloadAction{
			Value: ngapType.OverloadActionPresentRejectNonEmergencyMoDt,
		},
	}

	snssai := ngapConvert.SNssaiToNgap(TestAmf.TestAmf.PlmnSupportList[0].SNssaiList[0])
	sliceOverloadItem := ngapType.SliceOverloadItem{
		SNSSAI: snssai,
	}
	overloadStartNSSAIItem := ngapType.OverloadStartNSSAIItem{}
	overloadStartNSSAIItem.SliceOverloadList.List = append(overloadStartNSSAIItem.SliceOverloadList.List, sliceOverloadItem)

	overloadStartNSSAIList := ngapType.OverloadStartNSSAIList{}
	overloadStartNSSAIList.List = append(overloadStartNSSAIList.List, overloadStartNSSAIItem)

	pdu, err := ngap_message.BuildOverloadStart(&overloadResponse, 80, &overloadStartNSSAIList)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildOverloadStop(t *testing.T) {

	pdu, err := ngap_message.BuildOverloadStop()

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildDownlinkRanConfigurationTransfer(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	TLA := ngapType.TransportLayerAddress{
		Value: aper.BitString{
			Bytes:     []byte{0x12, 0x34, 0x50},
			BitLength: 22,
		},
	}

	var xnTLAs ngapType.XnTLAs

	xnTLAs.List = append(xnTLAs.List, TLA)

	sONConfigurationTransfer := ngapType.SONConfigurationTransfer{
		TargetRANNodeID: ngapType.TargetRANNodeID{
			GlobalRANNodeID: ngapType.GlobalRANNodeID{
				Present: ngapType.GlobalRANNodeIDPresentGlobalGNBID,
				GlobalGNBID: &ngapType.GlobalGNBID{
					PLMNIdentity: ngapType.PLMNIdentity{
						Value: aper.OctetString("\x02\xf8\x39"),
					},
					GNBID: ngapType.GNBID{
						Present: ngapType.GNBIDPresentGNBID,
						GNBID: &aper.BitString{
							Bytes:     []byte{0x45, 0x46, 0x47},
							BitLength: 24,
						},
					},
				},
			},
			SelectedTAI: ngapType.TAI{
				PLMNIdentity: ngapType.PLMNIdentity{
					Value: aper.OctetString("\x02\xf8\x39"),
				},
				TAC: ngapType.TAC{
					Value: aper.OctetString("\x00\x00\x01"),
				},
			},
		},
		SourceRANNodeID: ngapType.SourceRANNodeID{
			GlobalRANNodeID: ngapType.GlobalRANNodeID{
				Present: ngapType.GlobalRANNodeIDPresentGlobalGNBID,
				GlobalGNBID: &ngapType.GlobalGNBID{
					PLMNIdentity: ngapType.PLMNIdentity{
						Value: aper.OctetString("\x02\xf8\x39"),
					},
					GNBID: ngapType.GNBID{
						Present: ngapType.GNBIDPresentGNBID,
						GNBID: &aper.BitString{
							Bytes:     []byte{0x41, 0x42, 0x43},
							BitLength: 24,
						},
					},
				},
			},
			SelectedTAI: ngapType.TAI{
				PLMNIdentity: ngapType.PLMNIdentity{
					Value: aper.OctetString("\x02\xf8\x39"),
				},
				TAC: ngapType.TAC{
					Value: aper.OctetString("\x00\x00\x01"),
				},
			},
		},
		SONInformation: ngapType.SONInformation{
			Present: ngapType.SONInformationPresentSONInformationRequest,
			SONInformationRequest: &ngapType.SONInformationRequest{
				Value: ngapType.SONInformationRequestPresentXnTNLConfigurationInfo,
			},
		},
		XnTNLConfigurationInfo: ngapType.XnTNLConfigurationInfo{
			XnTransportLayerAddresses: ngapType.XnTLAs{
				List: xnTLAs.List,
			},
		},
	}

	pkg, err := ngap_message.BuildDownlinkRanConfigurationTransfer(&sONConfigurationTransfer)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))

	}
}

func TestBuildDownlinkNonUEAssociatedNRPPATransport(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	nRPPaPDU := ngapType.NRPPaPDU{
		Value: aper.OctetString("\x02\x03"),
	}

	pdu, err := ngap_message.BuildDownlinkNonUEAssociatedNRPPATransport(ue.RanUe[models.AccessType__3_GPP_ACCESS], nRPPaPDU)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildDeactivateTrace(t *testing.T) {
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 123
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 456

	pdu, err := ngap_message.BuildDeactivateTrace(ue, models.AccessType__3_GPP_ACCESS)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildLocationReportingControl(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 123
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 456

	eventType := ngapType.EventType{
		Value: ngapType.EventTypePresentStopChangeOfServeCell,
	}

	AOITAIItem := ngapType.AreaOfInterestTAIItem{
		TAI: ngapType.TAI{
			PLMNIdentity: ngapType.PLMNIdentity{
				Value: aper.OctetString("\x02\xf8\x39"),
			},
			TAC: ngapType.TAC{
				Value: aper.OctetString("\x00\x00\x01"),
			},
		},
	}

	AOICellItem := ngapType.AreaOfInterestCellItem{
		NGRANCGI: ngapType.NGRANCGI{
			Present: ngapType.NGRANCGIPresentNRCGI,
			NRCGI: &ngapType.NRCGI{
				PLMNIdentity: ngapType.PLMNIdentity{
					Value: aper.OctetString("\x02\xf8\x39"),
				},
				NRCellIdentity: ngapType.NRCellIdentity{
					Value: aper.BitString{
						Bytes:     []byte{0x11, 0x12, 0x21, 0x22, 0x20},
						BitLength: 36,
					},
				},
			},
		},
	}

	AOIRANNodeItem := ngapType.AreaOfInterestRANNodeItem{
		GlobalRANNodeID: ngapType.GlobalRANNodeID{
			Present: ngapType.GlobalRANNodeIDPresentGlobalGNBID,
			GlobalGNBID: &ngapType.GlobalGNBID{
				PLMNIdentity: ngapType.PLMNIdentity{
					Value: aper.OctetString("\x02\xf8\x39"),
				},
				GNBID: ngapType.GNBID{
					Present: ngapType.GNBIDPresentGNBID,
					GNBID: &aper.BitString{
						Bytes:     []byte{0x41, 0x42, 0x43},
						BitLength: 24,
					},
				},
			},
		},
	}

	AOIItem := ngapType.AreaOfInterestItem{
		AreaOfInterest: ngapType.AreaOfInterest{
			AreaOfInterestTAIList: &ngapType.AreaOfInterestTAIList{
				List: []ngapType.AreaOfInterestTAIItem{
					AOITAIItem,
				},
			},
			AreaOfInterestCellList: &ngapType.AreaOfInterestCellList{
				List: []ngapType.AreaOfInterestCellItem{
					AOICellItem,
				},
			},
			AreaOfInterestRANNodeList: &ngapType.AreaOfInterestRANNodeList{
				List: []ngapType.AreaOfInterestRANNodeItem{
					AOIRANNodeItem,
				},
			},
		},
		LocationReportingReferenceID: ngapType.LocationReportingReferenceID{
			Value: 12,
		},
	}

	var AOIList ngapType.AreaOfInterestList
	AOIList.List = append(AOIList.List, AOIItem)

	pkg, err := ngap_message.BuildLocationReportingControl(ue.RanUe[models.AccessType__3_GPP_ACCESS], &AOIList, 0, eventType)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pkg)
	}
	pdu, err1 := ngap.Decoder(pkg)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pdu).Elem(), 0))

	}
}

func TestBuildUETNLABindingReleaseRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 123
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 456

	pdu, err := ngap_message.BuildUETNLABindingReleaseRequest(ue.RanUe[models.AccessType__3_GPP_ACCESS])

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}

func TestBuildAMFConfigurationUpdate(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	amfSelf := amf_context.AMF_Self()
	amfSelf.HttpIPv4Address = "127.0.0.1"
	amfSelf.HttpIPv6Address = "2001:0db8:85a3:08d3:1319:8a2e:0370:7344"
	amfSelf.TNLWeightFactor = 123

	tNLassociationUsage := ngapType.TNLAssociationUsage{
		Value: ngapType.TNLAssociationUsagePresentBoth,
	}
	tNLAddressWeightFactor := ngapType.TNLAddressWeightFactor{
		Value: amfSelf.TNLWeightFactor,
	}

	pdu, err := ngap_message.BuildAMFConfigurationUpdate(tNLassociationUsage, tNLAddressWeightFactor)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}
func TestBuildDownlinkUEAssociatedNRPPaTransport(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	nRPPaPDU := ngapType.NRPPaPDU{
		Value: aper.OctetString("\x03\x02"),
	}

	pdu, err := ngap_message.BuildDownlinkUEAssociatedNRPPaTransport(ue.RanUe[models.AccessType__3_GPP_ACCESS], nRPPaPDU)

	if err != nil {
		t.Errorf("Encode is FAILED: %+v", err)
	} else {
		t.Logf("Encode : %0x", pdu)
	}
	pkg, err1 := ngap.Decoder(pdu)
	if err1 != nil {
		t.Errorf("Decode is FAILED: %+v", err1)
	} else {
		t.Logf("Decode : \n%s", ngap.PrintResult(reflect.ValueOf(pkg).Elem(), 0))
	}
}
