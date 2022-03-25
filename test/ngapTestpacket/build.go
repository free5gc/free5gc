package ngapTestpacket

import (
	"encoding/hex"

	"github.com/calee0219/fatal"

	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
)

// TODO: check test data
var TestPlmn ngapType.PLMNIdentity

func init() {
	TestPlmn.Value = aper.OctetString("\x02\xf8\x39")
}

func BuildNGSetupRequest() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentNGSetupRequest
	initiatingMessage.Value.NGSetupRequest = new(ngapType.NGSetupRequest)

	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	nGSetupRequestIEs := &nGSetupRequest.ProtocolIEs

	// GlobalRANNodeID
	ie := ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDGlobalRANNodeID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentGlobalRANNodeID
	ie.Value.GlobalRANNodeID = new(ngapType.GlobalRANNodeID)

	globalRANNodeID := ie.Value.GlobalRANNodeID
	globalRANNodeID.Present = ngapType.GlobalRANNodeIDPresentGlobalGNBID
	globalRANNodeID.GlobalGNBID = new(ngapType.GlobalGNBID)

	globalGNBID := globalRANNodeID.GlobalGNBID
	globalGNBID.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	globalGNBID.GNBID.Present = ngapType.GNBIDPresentGNBID
	globalGNBID.GNBID.GNBID = new(aper.BitString)

	gNBID := globalGNBID.GNBID.GNBID

	*gNBID = aper.BitString{
		Bytes:     []byte{0x45, 0x46, 0x47},
		BitLength: 24,
	}
	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	// RANNodeName
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANNodeName
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentRANNodeName
	ie.Value.RANNodeName = new(ngapType.RANNodeName)

	rANNodeName := ie.Value.RANNodeName
	rANNodeName.Value = "free5GC"
	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)
	// SupportedTAList
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSupportedTAList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentSupportedTAList
	ie.Value.SupportedTAList = new(ngapType.SupportedTAList)

	supportedTAList := ie.Value.SupportedTAList

	// SupportedTAItem in SupportedTAList
	supportedTAItem := ngapType.SupportedTAItem{}
	supportedTAItem.TAC.Value = aper.OctetString("\x00\x00\x01")

	broadcastPLMNList := &supportedTAItem.BroadcastPLMNList
	// BroadcastPLMNItem in BroadcastPLMNList
	broadcastPLMNItem := ngapType.BroadcastPLMNItem{}
	broadcastPLMNItem.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")

	sliceSupportList := &broadcastPLMNItem.TAISliceSupportList
	// SliceSupportItem in SliceSupportList
	sliceSupportItem := ngapType.SliceSupportItem{}
	sliceSupportItem.SNSSAI.SST.Value = aper.OctetString("\x01")
	// optional
	sliceSupportItem.SNSSAI.SD = new(ngapType.SD)
	sliceSupportItem.SNSSAI.SD.Value = aper.OctetString("\x01\x02\x03")

	sliceSupportList.List = append(sliceSupportList.List, sliceSupportItem)

	broadcastPLMNList.List = append(broadcastPLMNList.List, broadcastPLMNItem)

	supportedTAList.List = append(supportedTAList.List, supportedTAItem)

	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	// PagingDRX
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDDefaultPagingDRX
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentDefaultPagingDRX
	ie.Value.DefaultPagingDRX = new(ngapType.PagingDRX)

	pagingDRX := ie.Value.DefaultPagingDRX
	pagingDRX.Value = ngapType.PagingDRXPresentV128
	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	return pdu
}

func BuildNGReset(partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeNGReset
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentNGReset
	initiatingMessage.Value.NGReset = new(ngapType.NGReset)

	nGReset := initiatingMessage.Value.NGReset
	nGResetIEs := &nGReset.ProtocolIEs

	// Cause
	ie := ngapType.NGResetIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGResetIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentNas
	cause.Nas = new(ngapType.CauseNas)
	cause.Nas.Value = ngapType.CauseNasPresentNormalRelease

	nGResetIEs.List = append(nGResetIEs.List, ie)

	// CHOICE ResetType (NG interface; Part of NG interface)
	ie = ngapType.NGResetIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDResetType
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGResetIEsPresentResetType
	ie.Value.ResetType = new(ngapType.ResetType)

	resetType := ie.Value.ResetType
	if partOfNGInterface == nil {
		resetType.Present = ngapType.ResetTypePresentNGInterface
		resetType.NGInterface = new(ngapType.ResetAll)
		resetType.NGInterface.Value = ngapType.ResetAllPresentResetAll
	} else {
		resetType.Present = ngapType.ResetTypePresentPartOfNGInterface
		resetType.PartOfNGInterface = new(ngapType.UEAssociatedLogicalNGConnectionList)
		resetType.PartOfNGInterface = partOfNGInterface
	}

	nGResetIEs.List = append(nGResetIEs.List, ie)

	return pdu
}

func BuildNGResetAcknowledge() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeNGReset
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentNGResetAcknowledge
	successfulOutcome.Value.NGResetAcknowledge = new(ngapType.NGResetAcknowledge)

	nGResetAcknowledge := successfulOutcome.Value.NGResetAcknowledge
	nGResetAcknowledgeIEs := &nGResetAcknowledge.ProtocolIEs

	// UE-associated Logical NG-connection List (optional)
	ie := ngapType.NGResetAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGResetAcknowledgeIEsPresentUEAssociatedLogicalNGConnectionList
	ie.Value.UEAssociatedLogicalNGConnectionList = new(ngapType.UEAssociatedLogicalNGConnectionList)

	uEAssociatedLogicalNGConnectionList := ie.Value.UEAssociatedLogicalNGConnectionList

	// UE-associated Logical NG-connection Item in UE-associated Logical NG-connection List
	uEAssociatedLogicalNGConnectionItem := ngapType.UEAssociatedLogicalNGConnectionItem{}
	// AMF UE NGAP ID in UE-associated Logical NG-connection Item (optional)
	uEAssociatedLogicalNGConnectionItem.AMFUENGAPID = new(ngapType.AMFUENGAPID)
	uEAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value = 123
	// RAN UE NGAP ID in UE-associated Logical NG-connection Item (optional)
	uEAssociatedLogicalNGConnectionItem.RANUENGAPID = new(ngapType.RANUENGAPID)
	uEAssociatedLogicalNGConnectionItem.RANUENGAPID.Value = 456

	uEAssociatedLogicalNGConnectionList.List =
		append(uEAssociatedLogicalNGConnectionList.List, uEAssociatedLogicalNGConnectionItem)

	nGResetAcknowledgeIEs.List = append(nGResetAcknowledgeIEs.List, ie)

	// Criticality Diagnostics (optional)
	return pdu
}

func BuildInitialUEMessage(ranUeNgapID int64, nasPdu []byte, fiveGSTmsi string) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeInitialUEMessage
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentInitialUEMessage
	initiatingMessage.Value.InitialUEMessage = new(ngapType.InitialUEMessage)

	initialUEMessage := initiatingMessage.Value.InitialUEMessage
	initialUEMessageIEs := &initialUEMessage.ProtocolIEs

	// RAN UE NGAP ID
	ie := ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// NAS-PDU
	ie = ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNASPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentNASPDU
	ie.Value.NASPDU = new(ngapType.NASPDU)

	// TODO: complete NAS-PDU
	nASPDU := ie.Value.NASPDU
	nASPDU.Value = nasPdu

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// User Location Information
	ie = ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity.Value = TestPlmn.Value
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}

	userLocationInformationNR.TAI.PLMNIdentity.Value = TestPlmn.Value
	userLocationInformationNR.TAI.TAC.Value = aper.OctetString("\x00\x00\x01")

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// RRC Establishment Cause
	ie = ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRRCEstablishmentCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentRRCEstablishmentCause
	ie.Value.RRCEstablishmentCause = new(ngapType.RRCEstablishmentCause)

	rRCEstablishmentCause := ie.Value.RRCEstablishmentCause
	rRCEstablishmentCause.Value = ngapType.RRCEstablishmentCausePresentMtAccess

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// 5G-S-TSMI (optional)
	if fiveGSTmsi != "" {
		ie = ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDFiveGSTMSI
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentFiveGSTMSI
		ie.Value.FiveGSTMSI = new(ngapType.FiveGSTMSI)

		fiveGSTMSI := ie.Value.FiveGSTMSI
		amfSetID, err := hex.DecodeString(fiveGSTmsi[:4])
		if err != nil {
			fatal.Fatalf("DecodeString error in BuildInitialUEMessage: %+v", err)
		}
		fiveGSTMSI.AMFSetID.Value = aper.BitString{
			Bytes:     amfSetID,
			BitLength: 10,
		}
		amfPointer, err := hex.DecodeString(fiveGSTmsi[2:4])
		if err != nil {
			fatal.Fatalf("DecodeString error in BuildInitialUEMessage: %+v", err)
		}
		fiveGSTMSI.AMFPointer.Value = aper.BitString{
			Bytes:     amfPointer,
			BitLength: 6,
		}
		tmsi, err := hex.DecodeString(fiveGSTmsi[4:])
		if err != nil {
			fatal.Fatalf("DecodeString error in BuildInitialUEMessage: %+v", err)
		}
		fiveGSTMSI.FiveGTMSI.Value = aper.OctetString(tmsi)

		initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)
	}
	// AMF Set ID (optional)

	// UE Context Request (optional)
	ie = ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUEContextRequest
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentUEContextRequest
	ie.Value.UEContextRequest = new(ngapType.UEContextRequest)
	ie.Value.UEContextRequest.Value = ngapType.UEContextRequestPresentRequested
	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// Allowed NSSAI (optional)
	return pdu
}

func BuildErrorIndication() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeErrorIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentErrorIndication
	initiatingMessage.Value.ErrorIndication = new(ngapType.ErrorIndication)

	errorIndication := initiatingMessage.Value.ErrorIndication
	errorIndicationIEs := &errorIndication.ProtocolIEs

	// AMF UE NGAP ID (optional)
	ie := ngapType.ErrorIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.ErrorIndicationIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 12345

	errorIndicationIEs.List = append(errorIndicationIEs.List, ie)

	// RAN UE NGAP ID (optional)
	ie = ngapType.ErrorIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.ErrorIndicationIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 1234

	errorIndicationIEs.List = append(errorIndicationIEs.List, ie)

	// Cause (optional)
	ie = ngapType.ErrorIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.ErrorIndicationIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentNas
	cause.Nas = new(ngapType.CauseNas)
	cause.Nas.Value = ngapType.CauseNasPresentNormalRelease

	errorIndicationIEs.List = append(errorIndicationIEs.List, ie)

	// Criticality Diagnostics (optional)
	ie = ngapType.ErrorIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.ErrorIndicationIEsPresentCriticalityDiagnostics
	ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

	criticalityDiagnostics := ie.Value.CriticalityDiagnostics

	criticalityDiagnostics.IEsCriticalityDiagnostics = new(ngapType.CriticalityDiagnosticsIEList)
	criticalityDiagnosticsIEItem := ngapType.CriticalityDiagnosticsIEItem{}

	criticalityDiagnosticsIEItem.IECriticality.Value = ngapType.CriticalityPresentReject
	criticalityDiagnosticsIEItem.IEID.Value = ngapType.ProtocolIEIDAMFUENGAPID
	criticalityDiagnosticsIEItem.TypeOfError.Value = ngapType.TypeOfErrorPresentMissing

	criticalityDiagnostics.IEsCriticalityDiagnostics.List =
		append(criticalityDiagnostics.IEsCriticalityDiagnostics.List, criticalityDiagnosticsIEItem)
	errorIndicationIEs.List = append(errorIndicationIEs.List, ie)

	return pdu
}

func BuildUEContextReleaseRequest(amfUeNgapID, ranUeNgapID int64, pduSessionIDList []int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUEContextReleaseRequest
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUEContextReleaseRequest
	initiatingMessage.Value.UEContextReleaseRequest = new(ngapType.UEContextReleaseRequest)

	uEContextReleaseRequest := initiatingMessage.Value.UEContextReleaseRequest
	uEContextReleaseRequestIEs := &uEContextReleaseRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	// PDU Session Resource List
	if pduSessionIDList != nil {
		ie = ngapType.UEContextReleaseRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceListCxtRelReq
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentPDUSessionResourceListCxtRelReq
		ie.Value.PDUSessionResourceListCxtRelReq = new(ngapType.PDUSessionResourceListCxtRelReq)

		pDUSessionResourceListCxtRelReq := ie.Value.PDUSessionResourceListCxtRelReq

		// PDU Session Resource Item in PDU session Resource List
		for _, pduSessionID := range pduSessionIDList {
			pDUSessionResourceItem := ngapType.PDUSessionResourceItemCxtRelReq{}
			pDUSessionResourceItem.PDUSessionID.Value = pduSessionID
			pDUSessionResourceListCxtRelReq.List = append(pDUSessionResourceListCxtRelReq.List, pDUSessionResourceItem)
		}
		uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)
	}

	// Cause
	ie = ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentRadioNetwork
	cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
	cause.RadioNetwork.Value = ngapType.CauseRadioNetworkPresentTxnrelocoverallExpiry

	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	return pdu
}

func BuildUEContextReleaseComplete(amfUeNgapID, ranUeNgapID int64, pduSessionIDList []int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextRelease
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUEContextReleaseComplete
	successfulOutcome.Value.UEContextReleaseComplete = new(ngapType.UEContextReleaseComplete)

	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	uEContextReleaseCompleteIEs := &uEContextReleaseComplete.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// User Location Information (optional)
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}

	userLocationInformationNR.TAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.TAI.TAC.Value = aper.OctetString("\x00\x00\x11")

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)
	// Information on Recommended Cells and RAN Nodes for Paging (optional)

	// PDU Session Resource List
	if pduSessionIDList != nil {
		ie = ngapType.UEContextReleaseCompleteIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceListCxtRelCpl
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentPDUSessionResourceListCxtRelCpl
		ie.Value.PDUSessionResourceListCxtRelCpl = new(ngapType.PDUSessionResourceListCxtRelCpl)

		pDUSessionResourceListCxtRelCpl := ie.Value.PDUSessionResourceListCxtRelCpl

		// PDU Session Resource Item (in PDU Session Resource List)
		for _, pduSessionID := range pduSessionIDList {
			pDUSessionResourceItemCxtRelCpl := ngapType.PDUSessionResourceItemCxtRelCpl{}
			pDUSessionResourceItemCxtRelCpl.PDUSessionID.Value = pduSessionID
			pDUSessionResourceListCxtRelCpl.List = append(pDUSessionResourceListCxtRelCpl.List, pDUSessionResourceItemCxtRelCpl)
		}

		uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)
	}

	// Criticality Diagnostics (optional)
	return pdu
}

func BuildUEContextModificationResponse(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextModification
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUEContextModificationResponse
	successfulOutcome.Value.UEContextModificationResponse = new(ngapType.UEContextModificationResponse)

	uEContextModificationResponse := successfulOutcome.Value.UEContextModificationResponse
	uEContextModificationResponseIEs := &uEContextModificationResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	// RRC State (optional)
	ie = ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRRCState
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationResponseIEsPresentRRCState
	ie.Value.RRCState = new(ngapType.RRCState)

	ie.Value.RRCState.Value = ngapType.RRCStatePresentConnected

	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	// User Location Information (optional)
	ie = ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationResponseIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInfo := ie.Value.UserLocationInformation
	userLocationInfo.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInfo.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	locationNR := userLocationInfo.UserLocationInformationNR
	locationNR.NRCGI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	locationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}

	locationNR.TAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	locationNR.TAI.TAC.Value = aper.OctetString("\x00\x00\x11")

	locationNR.TimeStamp = new(ngapType.TimeStamp)
	locationNR.TimeStamp.Value = aper.OctetString("\x00\x00\x11\x11")

	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	// Criticality Diagnostics (optional)
	return pdu
}

func BuildUplinkNasTransport(amfUeNgapID, ranUeNgapID int64, nasPdu []byte) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUplinkNASTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUplinkNASTransport
	initiatingMessage.Value.UplinkNASTransport = new(ngapType.UplinkNASTransport)

	uplinkNasTransport := initiatingMessage.Value.UplinkNASTransport
	uplinkNasTransportIEs := &uplinkNasTransport.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// NAS-PDU
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNASPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentNASPDU
	ie.Value.NASPDU = new(ngapType.NASPDU)

	// TODO: complete NAS-PDU
	nASPDU := ie.Value.NASPDU
	nASPDU.Value = nasPdu

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// User Location Information
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}

	userLocationInformationNR.TAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.TAI.TAC.Value = aper.OctetString("\x00\x00\x01")

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	return pdu
}

func BuildInitialContextSetupResponse(amfUeNgapID, ranUeNgapID int64, ipv4 string,
	pduSessionFailedList *ngapType.PDUSessionResourceFailedToSetupListCxtRes) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentInitialContextSetupResponse
	successfulOutcome.Value.InitialContextSetupResponse = new(ngapType.InitialContextSetupResponse)

	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	initialContextSetupResponseIEs := &initialContextSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// PDU Session Resource Setup Response List
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtRes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentPDUSessionResourceSetupListCxtRes
	ie.Value.PDUSessionResourceSetupListCxtRes = new(ngapType.PDUSessionResourceSetupListCxtRes)

	pDUSessionResourceSetupListCxtRes := ie.Value.PDUSessionResourceSetupListCxtRes

	// PDU Session Resource Setup Response Item in PDU Session Resource Setup Response List
	pDUSessionResourceSetupItemCxtRes := ngapType.PDUSessionResourceSetupItemCxtRes{}
	pDUSessionResourceSetupItemCxtRes.PDUSessionID.Value = 10
	pDUSessionResourceSetupItemCxtRes.PDUSessionResourceSetupResponseTransfer =
		GetPDUSessionResourceSetupResponseTransfer(ipv4)

	pDUSessionResourceSetupListCxtRes.List =
		append(pDUSessionResourceSetupListCxtRes.List, pDUSessionResourceSetupItemCxtRes)

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// PDU Session Resource Failed to Setup List
	if pduSessionFailedList != nil {
		ie = ngapType.InitialContextSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtRes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentPDUSessionResourceFailedToSetupListCxtRes
		ie.Value.PDUSessionResourceFailedToSetupListCxtRes = pduSessionFailedList
		initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)
	}
	// Criticality Diagnostics (optional)
	return pdu
}

func BuildInitialContextSetupFailure(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentInitialContextSetupFailure
	unsuccessfulOutcome.Value.InitialContextSetupFailure = new(ngapType.InitialContextSetupFailure)

	initialContextSetupFailure := unsuccessfulOutcome.Value.InitialContextSetupFailure
	initialContextSetupFailureIEs := &initialContextSetupFailure.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// PDU Session Resource Failed to Setup List
	ie = ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtFail
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentPDUSessionResourceFailedToSetupListCxtFail
	ie.Value.PDUSessionResourceFailedToSetupListCxtFail = new(ngapType.PDUSessionResourceFailedToSetupListCxtFail)

	pDUSessionResourceFailedToSetupListCxtFail := ie.Value.PDUSessionResourceFailedToSetupListCxtFail

	// PDU Session Resource Failed to Setup Item in PDU Sessuin Resource Failed to Setup List
	pDUSessionResourceFailedToSetupItemCxtFail := ngapType.PDUSessionResourceFailedToSetupItemCxtFail{}
	pDUSessionResourceFailedToSetupItemCxtFail.PDUSessionID.Value = 10
	pDUSessionResourceFailedToSetupItemCxtFail.PDUSessionResourceSetupUnsuccessfulTransfer = aper.OctetString("\x11\x22")

	pDUSessionResourceFailedToSetupListCxtFail.List =
		append(pDUSessionResourceFailedToSetupListCxtFail.List, pDUSessionResourceFailedToSetupItemCxtFail)

	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// Cause
	ie = ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentNas
	cause.Nas = new(ngapType.CauseNas)
	cause.Nas.Value = ngapType.CauseNasPresentNormalRelease

	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// Criticality Diagnostics (optional)
	return pdu
}

func BuildPathSwitchRequest(sourceAmfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePathSwitchRequest
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPathSwitchRequest
	initiatingMessage.Value.PathSwitchRequest = new(ngapType.PathSwitchRequest)

	pathSwitchRequest := initiatingMessage.Value.PathSwitchRequest
	pathSwitchRequestIEs := &pathSwitchRequest.ProtocolIEs

	// RAN UE NGAP ID
	ie := ngapType.PathSwitchRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PathSwitchRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pathSwitchRequestIEs.List = append(pathSwitchRequestIEs.List, ie)

	// Source AMF UE NGAP ID (equal to AMF UE NGAP ID)
	ie = ngapType.PathSwitchRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSourceAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PathSwitchRequestIEsPresentSourceAMFUENGAPID
	ie.Value.SourceAMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.SourceAMFUENGAPID
	aMFUENGAPID.Value = sourceAmfUeNgapID

	pathSwitchRequestIEs.List = append(pathSwitchRequestIEs.List, ie)

	// User Location Information
	ie = ngapType.PathSwitchRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PathSwitchRequestIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x20},
		BitLength: 36,
	}

	userLocationInformationNR.TAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.TAI.TAC.Value = aper.OctetString("\x00\x00\x11")

	pathSwitchRequestIEs.List = append(pathSwitchRequestIEs.List, ie)

	// UE Security Capabilities
	ie = ngapType.PathSwitchRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUESecurityCapabilities
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PathSwitchRequestIEsPresentUESecurityCapabilities
	ie.Value.UESecurityCapabilities = new(ngapType.UESecurityCapabilities)

	uESecurityCapabilities := ie.Value.UESecurityCapabilities
	uESecurityCapabilities.NRencryptionAlgorithms.Value = aper.BitString{
		Bytes:     []byte{0xff, 0xff},
		BitLength: 16,
	}
	uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value = aper.BitString{
		Bytes:     []byte{0xff, 0xff},
		BitLength: 16,
	}
	uESecurityCapabilities.EUTRAencryptionAlgorithms.Value = aper.BitString{
		Bytes:     []byte{0xff, 0xff},
		BitLength: 16,
	}
	uESecurityCapabilities.EUTRAintegrityProtectionAlgorithms.Value = aper.BitString{
		Bytes:     []byte{0xff, 0xff},
		BitLength: 16,
	}

	pathSwitchRequestIEs.List = append(pathSwitchRequestIEs.List, ie)

	// PDU Session Resource to be Switched in Downlink List
	ie = ngapType.PathSwitchRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceToBeSwitchedDLList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PathSwitchRequestIEsPresentPDUSessionResourceToBeSwitchedDLList
	ie.Value.PDUSessionResourceToBeSwitchedDLList = new(ngapType.PDUSessionResourceToBeSwitchedDLList)

	pDUSessionResourceToBeSwitchedDLList := ie.Value.PDUSessionResourceToBeSwitchedDLList

	// PDU Session Resource to be Switched in Downlink Item (in PDU Session Resource to be Switched in Downlink List)
	pDUSessionResourceToBeSwitchedDLItem := ngapType.PDUSessionResourceToBeSwitchedDLItem{}
	pDUSessionResourceToBeSwitchedDLItem.PDUSessionID.Value = 10
	pDUSessionResourceToBeSwitchedDLItem.PathSwitchRequestTransfer = GetPathSwitchRequestTransfer()

	pDUSessionResourceToBeSwitchedDLList.List =
		append(pDUSessionResourceToBeSwitchedDLList.List, pDUSessionResourceToBeSwitchedDLItem)

	pathSwitchRequestIEs.List = append(pathSwitchRequestIEs.List, ie)

	// PDU Session Resource Failed to Setup List
	ie = ngapType.PathSwitchRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListPSReq
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PathSwitchRequestIEsPresentPDUSessionResourceFailedToSetupListPSReq
	ie.Value.PDUSessionResourceFailedToSetupListPSReq = new(ngapType.PDUSessionResourceFailedToSetupListPSReq)

	pDUSessionResourceFailedToSetupListPSReq := ie.Value.PDUSessionResourceFailedToSetupListPSReq

	// PDU Session Resource Failed to Setup Item (in PDU Session Resource Failed to Setup List)
	pDUSessionResourceFailedToSetupItemPSReq := ngapType.PDUSessionResourceFailedToSetupItemPSReq{}
	pDUSessionResourceFailedToSetupItemPSReq.PDUSessionID.Value = 11
	pDUSessionResourceFailedToSetupItemPSReq.PathSwitchRequestSetupFailedTransfer =
		GetPathSwitchRequestSetupFailedTransfer()

	pDUSessionResourceFailedToSetupListPSReq.List =
		append(pDUSessionResourceFailedToSetupListPSReq.List, pDUSessionResourceFailedToSetupItemPSReq)

	pathSwitchRequestIEs.List = append(pathSwitchRequestIEs.List, ie)

	return pdu
}

func BuildHandoverRequestAcknowledge(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeHandoverResourceAllocation
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentHandoverRequestAcknowledge
	successfulOutcome.Value.HandoverRequestAcknowledge = new(ngapType.HandoverRequestAcknowledge)

	handoverRequestAcknowledge := successfulOutcome.Value.HandoverRequestAcknowledge
	handoverRequestAcknowledgeIEs := &handoverRequestAcknowledge.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.HandoverRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverRequestAcknowledgeIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	handoverRequestAcknowledgeIEs.List = append(handoverRequestAcknowledgeIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.HandoverRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverRequestAcknowledgeIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	handoverRequestAcknowledgeIEs.List = append(handoverRequestAcknowledgeIEs.List, ie)

	//PDU Session Resource Admitted List
	ie = ngapType.HandoverRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceAdmittedList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverRequestAcknowledgeIEsPresentPDUSessionResourceAdmittedList
	ie.Value.PDUSessionResourceAdmittedList = new(ngapType.PDUSessionResourceAdmittedList)

	pDUSessionResourceAdmittedList := ie.Value.PDUSessionResourceAdmittedList

	//PDU SessionResource Admittedy Item
	pDUSessionResourceAdmittedItem := ngapType.PDUSessionResourceAdmittedItem{}
	pDUSessionResourceAdmittedItem.PDUSessionID.Value = 10
	pDUSessionResourceAdmittedItem.HandoverRequestAcknowledgeTransfer = GetHandoverRequestAcknowledgeTransfer()

	pDUSessionResourceAdmittedList.List = append(pDUSessionResourceAdmittedList.List, pDUSessionResourceAdmittedItem)

	handoverRequestAcknowledgeIEs.List = append(handoverRequestAcknowledgeIEs.List, ie)

	//PDU Session Resource Failed to setup List
	ie = ngapType.HandoverRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListHOAck
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverRequestAcknowledgeIEsPresentPDUSessionResourceFailedToSetupListHOAck
	ie.Value.PDUSessionResourceFailedToSetupListHOAck = new(ngapType.PDUSessionResourceFailedToSetupListHOAck)

	pDUSessionResourceFailedToSetupListHOAck := ie.Value.PDUSessionResourceFailedToSetupListHOAck

	//PDU Session Resource Failed to setup Item
	pDUSessionResourceFailedToSetupItemHOAck := ngapType.PDUSessionResourceFailedToSetupItemHOAck{}
	pDUSessionResourceFailedToSetupItemHOAck.PDUSessionID.Value = 11
	pDUSessionResourceFailedToSetupItemHOAck.HandoverResourceAllocationUnsuccessfulTransfer =
		GetHandoverResourceAllocationUnsuccessfulTransfer()

	pDUSessionResourceFailedToSetupListHOAck.List =
		append(pDUSessionResourceFailedToSetupListHOAck.List, pDUSessionResourceFailedToSetupItemHOAck)

	handoverRequestAcknowledgeIEs.List = append(handoverRequestAcknowledgeIEs.List, ie)

	//Target To Source TransparentContainer
	ie = ngapType.HandoverRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDTargetToSourceTransparentContainer
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestAcknowledgeIEsPresentTargetToSourceTransparentContainer
	ie.Value.TargetToSourceTransparentContainer = new(ngapType.TargetToSourceTransparentContainer)

	targetToSourceTransparentContainer := ie.Value.TargetToSourceTransparentContainer
	targetToSourceTransparentContainer.Value = aper.OctetString("\x00\x01\x00\x00")

	handoverRequestAcknowledgeIEs.List = append(handoverRequestAcknowledgeIEs.List, ie)

	// Criticality Diagnostics (optional)
	return pdu
}

func BuildHandoverFailure(amfUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	UnsuccessfulOutcome := pdu.UnsuccessfulOutcome
	UnsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeHandoverResourceAllocation
	UnsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentIgnore

	UnsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentHandoverFailure
	UnsuccessfulOutcome.Value.HandoverFailure = new(ngapType.HandoverFailure)

	handoverFailure := UnsuccessfulOutcome.Value.HandoverFailure
	handoverFailureIEs := &handoverFailure.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.HandoverFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverFailureIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	handoverFailureIEs.List = append(handoverFailureIEs.List, ie)
	// Cause
	ie = ngapType.HandoverFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverFailureIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentNas
	cause.Nas = new(ngapType.CauseNas)
	cause.Nas.Value = ngapType.CauseNasPresentNormalRelease

	handoverFailureIEs.List = append(handoverFailureIEs.List, ie)

	//Criticality Diagnostics (optional)

	return pdu
}

func BuildPDUSessionResourceReleaseResponse() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceRelease
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceReleaseResponse
	successfulOutcome.Value.PDUSessionResourceReleaseResponse = new(ngapType.PDUSessionResourceReleaseResponse)

	pDUSessionResourceReleaseResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	pDUSessionResourceReleaseResponseIEs := &pDUSessionResourceReleaseResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 4294967295

	pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)

	// PDU Session Resource Released List
	ie = ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceReleasedListRelRes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentPDUSessionResourceReleasedListRelRes
	ie.Value.PDUSessionResourceReleasedListRelRes = new(ngapType.PDUSessionResourceReleasedListRelRes)

	pDUSessionResourceReleasedListRelRes := ie.Value.PDUSessionResourceReleasedListRelRes

	// PDU Session Resource Released Item
	pDUSessionResourceReleasedItemRelRes := ngapType.PDUSessionResourceReleasedItemRelRes{}
	pDUSessionResourceReleasedItemRelRes.PDUSessionID.Value = 10
	pDUSessionResourceReleasedItemRelRes.PDUSessionResourceReleaseResponseTransfer = aper.OctetString("\x01\x02\x03")

	pDUSessionResourceReleasedListRelRes.List =
		append(pDUSessionResourceReleasedListRelRes.List, pDUSessionResourceReleasedItemRelRes)

	pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)

	// User Location Information (optional)
	ie = ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}

	userLocationInformationNR.TAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.TAI.TAC.Value = aper.OctetString("\x00\x00\x11")

	pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)

	// Criticality Diagnostics (optional)
	return pdu
}

func BuildAMFConfigurationUpdateFailure() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeAMFConfigurationUpdate
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentAMFConfigurationUpdateFailure
	unsuccessfulOutcome.Value.AMFConfigurationUpdateFailure = new(ngapType.AMFConfigurationUpdateFailure)

	AMFConfigurationUpdateFailure := unsuccessfulOutcome.Value.AMFConfigurationUpdateFailure
	AMFConfigurationUpdateFailureIEs := &AMFConfigurationUpdateFailure.ProtocolIEs

	//	Cause
	ie := ngapType.AMFConfigurationUpdateFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.AMFConfigurationUpdateFailureIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentRadioNetwork
	cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
	cause.RadioNetwork.Value = ngapType.CauseRadioNetworkPresentTxnrelocoverallExpiry

	AMFConfigurationUpdateFailureIEs.List = append(AMFConfigurationUpdateFailureIEs.List, ie)

	//	TODO: Time to wait (optional)

	//	TODO: Criticality Diagnostics (optional)

	return pdu

}

func BuildUERadioCapabilityCheckRequest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUERadioCapabilityCheck
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUERadioCapabilityCheckRequest
	initiatingMessage.Value.UERadioCapabilityCheckRequest = new(ngapType.UERadioCapabilityCheckRequest)

	uERadioCapabilityCheckRequest := initiatingMessage.Value.UERadioCapabilityCheckRequest
	uERadioCapabilityCheckRequestIEs := &uERadioCapabilityCheckRequest.ProtocolIEs
	// AMFUENGAPID
	{
		ie := ngapType.UERadioCapabilityCheckRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.UERadioCapabilityCheckRequestIEsPresentAMFUENGAPID
		ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

		aMFUENGAPID := ie.Value.AMFUENGAPID
		aMFUENGAPID.Value = amfUeNgapID
		uERadioCapabilityCheckRequestIEs.List = append(uERadioCapabilityCheckRequestIEs.List, ie)
	}
	// RANUENGAPID
	{
		ie := ngapType.UERadioCapabilityCheckRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.UERadioCapabilityCheckRequestIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = ranUeNgapID
		uERadioCapabilityCheckRequestIEs.List = append(uERadioCapabilityCheckRequestIEs.List, ie)
	}
	// UERadioCapability
	{
		ie := ngapType.UERadioCapabilityCheckRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUERadioCapability
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UERadioCapabilityCheckRequestIEsPresentUERadioCapability
		ie.Value.UERadioCapability = new(ngapType.UERadioCapability)

		uERadioCapability := ie.Value.UERadioCapability
		uERadioCapability.Value = aper.OctetString("\x00\x00\x01")

		uERadioCapabilityCheckRequestIEs.List = append(uERadioCapabilityCheckRequestIEs.List, ie)
	}

	return pdu

}

func BuildUERadioCapabilityCheckResponse() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUERadioCapabilityCheck
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUERadioCapabilityCheckResponse
	successfulOutcome.Value.UERadioCapabilityCheckResponse = new(ngapType.UERadioCapabilityCheckResponse)

	uERadioCapabilityCheckResponse := successfulOutcome.Value.UERadioCapabilityCheckResponse
	uERadioCapabilityCheckResponseIEs := &uERadioCapabilityCheckResponse.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.UERadioCapabilityCheckResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UERadioCapabilityCheckResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	uERadioCapabilityCheckResponseIEs.List = append(uERadioCapabilityCheckResponseIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.UERadioCapabilityCheckResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UERadioCapabilityCheckResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 0xFFFFFFFF

	uERadioCapabilityCheckResponseIEs.List = append(uERadioCapabilityCheckResponseIEs.List, ie)

	//IMS Voice Support Indicator
	ie = ngapType.UERadioCapabilityCheckResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDIMSVoiceSupportIndicator
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UERadioCapabilityCheckResponseIEsPresentIMSVoiceSupportIndicator
	ie.Value.IMSVoiceSupportIndicator = new(ngapType.IMSVoiceSupportIndicator)

	iMSVoiceSupportIndicator := ie.Value.IMSVoiceSupportIndicator
	iMSVoiceSupportIndicator.Value = ngapType.IMSVoiceSupportIndicatorPresentNotSupported

	uERadioCapabilityCheckResponseIEs.List = append(uERadioCapabilityCheckResponseIEs.List, ie)

	//TODO:Criticality Diagnostics (optional)

	return pdu
}

func BuildHandoverCancel() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeHandoverCancel
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentHandoverCancel
	initiatingMessage.Value.HandoverCancel = new(ngapType.HandoverCancel)

	handoverCancel := initiatingMessage.Value.HandoverCancel
	handoverCancelIEs := &handoverCancel.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.HandoverCancelIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverCancelIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	handoverCancelIEs.List = append(handoverCancelIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.HandoverCancelIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverCancelIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 0xffffffff

	handoverCancelIEs.List = append(handoverCancelIEs.List, ie)

	//Cause
	ie = ngapType.HandoverCancelIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverCancelIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentRadioNetwork
	cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
	cause.RadioNetwork.Value = ngapType.CauseRadioNetworkPresentHandoverCancelled

	handoverCancelIEs.List = append(handoverCancelIEs.List, ie)

	return pdu
}
func BuildLocationReportingFailureIndication() (pdu ngapType.NGAPPDU) {
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeLocationReportingFailureIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentLocationReportingFailureIndication
	initiatingMessage.Value.LocationReportingFailureIndication = new(ngapType.LocationReportingFailureIndication)

	locationReportingFailureIndication := initiatingMessage.Value.LocationReportingFailureIndication
	locationReportingFailureIndicationIEs := &locationReportingFailureIndication.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.LocationReportingFailureIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.LocationReportingFailureIndicationIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	locationReportingFailureIndicationIEs.List = append(locationReportingFailureIndicationIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.LocationReportingFailureIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.LocationReportingFailureIndicationIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 0xffffffff

	locationReportingFailureIndicationIEs.List = append(locationReportingFailureIndicationIEs.List, ie)

	//Cause
	ie = ngapType.LocationReportingFailureIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.LocationReportingFailureIndicationIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentRadioNetwork
	cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
	cause.RadioNetwork.Value = ngapType.CauseRadioNetworkPresentFailureInRadioInterfaceProcedure

	locationReportingFailureIndicationIEs.List = append(locationReportingFailureIndicationIEs.List, ie)

	return pdu
}

func BuildPDUSessionResourceSetupResponse(amfUeNgapID, ranUeNgapID int64, ipv4 string) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceSetupResponse
	successfulOutcome.Value.PDUSessionResourceSetupResponse = new(ngapType.PDUSessionResourceSetupResponse)

	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	pDUSessionResourceSetupResponseIEs := &pDUSessionResourceSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// PDU Session Resource Setup Response List
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceSetupListSURes
	ie.Value.PDUSessionResourceSetupListSURes = new(ngapType.PDUSessionResourceSetupListSURes)

	pDUSessionResourceSetupListSURes := ie.Value.PDUSessionResourceSetupListSURes

	// PDU Session Resource Setup Response Item in PDU Session Resource Setup Response List
	pDUSessionResourceSetupItemSURes := ngapType.PDUSessionResourceSetupItemSURes{}
	pDUSessionResourceSetupItemSURes.PDUSessionID.Value = 10

	pDUSessionResourceSetupItemSURes.PDUSessionResourceSetupResponseTransfer =
		GetPDUSessionResourceSetupResponseTransfer(ipv4)

	pDUSessionResourceSetupListSURes.List = append(pDUSessionResourceSetupListSURes.List, pDUSessionResourceSetupItemSURes)

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// PDU Sessuin Resource Failed to Setup List
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceFailedToSetupListSURes
	ie.Value.PDUSessionResourceFailedToSetupListSURes = new(ngapType.PDUSessionResourceFailedToSetupListSURes)

	pDUSessionResourceFailedToSetupListSURes := ie.Value.PDUSessionResourceFailedToSetupListSURes

	// PDU Session Resource Failed to Setup Item in PDU Sessuin Resource Failed to Setup List
	pDUSessionResourceFailedToSetupItemSURes := ngapType.PDUSessionResourceFailedToSetupItemSURes{}
	pDUSessionResourceFailedToSetupItemSURes.PDUSessionID.Value = 10
	pDUSessionResourceFailedToSetupItemSURes.PDUSessionResourceSetupUnsuccessfulTransfer =
		GetPDUSessionResourceSetupUnsucessfulTransfer()

	pDUSessionResourceFailedToSetupListSURes.List =
		append(pDUSessionResourceFailedToSetupListSURes.List, pDUSessionResourceFailedToSetupItemSURes)

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)
	// Criticality Diagnostics (optional)
	return pdu
}

func BuildPDUSessionResourceSetupResponseForPaging(amfUeNgapID, ranUeNgapID int64, ipv4 string) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceSetupResponse
	successfulOutcome.Value.PDUSessionResourceSetupResponse = new(ngapType.PDUSessionResourceSetupResponse)

	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	pDUSessionResourceSetupResponseIEs := &pDUSessionResourceSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// PDU Session Resource Setup Response List
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceSetupListSURes
	ie.Value.PDUSessionResourceSetupListSURes = new(ngapType.PDUSessionResourceSetupListSURes)

	pDUSessionResourceSetupListSURes := ie.Value.PDUSessionResourceSetupListSURes

	// PDU Session Resource Setup Response Item in PDU Session Resource Setup Response List
	pDUSessionResourceSetupItemSURes := ngapType.PDUSessionResourceSetupItemSURes{}
	pDUSessionResourceSetupItemSURes.PDUSessionID.Value = 10

	pDUSessionResourceSetupItemSURes.PDUSessionResourceSetupResponseTransfer =
		GetPDUSessionResourceSetupResponseTransfer(ipv4)

	pDUSessionResourceSetupListSURes.List = append(pDUSessionResourceSetupListSURes.List, pDUSessionResourceSetupItemSURes)

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// PDU Sessuin Resource Failed to Setup List
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceFailedToSetupListSURes
	ie.Value.PDUSessionResourceFailedToSetupListSURes = new(ngapType.PDUSessionResourceFailedToSetupListSURes)

	// Criticality Diagnostics (optional)
	return pdu
}

func BuildPDUSessionResourceModifyResponse(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceModify
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceModifyResponse
	successfulOutcome.Value.PDUSessionResourceModifyResponse = new(ngapType.PDUSessionResourceModifyResponse)

	pDUSessionResourceModifyResponse := successfulOutcome.Value.PDUSessionResourceModifyResponse
	pDUSessionResourceModifyResponseIEs := &pDUSessionResourceModifyResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceModifyResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	pDUSessionResourceModifyResponseIEs.List = append(pDUSessionResourceModifyResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceModifyResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pDUSessionResourceModifyResponseIEs.List = append(pDUSessionResourceModifyResponseIEs.List, ie)

	// PDU Session Resource Modify Response List
	ie = ngapType.PDUSessionResourceModifyResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceModifyListModRes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentPDUSessionResourceModifyListModRes
	ie.Value.PDUSessionResourceModifyListModRes = new(ngapType.PDUSessionResourceModifyListModRes)

	pDUSessionResourceModifyListModRes := ie.Value.PDUSessionResourceModifyListModRes

	// PDU Session Resource Modify Response Item in PDU Session Resource Modify Response List
	pDUSessionResourceModifyResponseItem := ngapType.PDUSessionResourceModifyItemModRes{}
	pDUSessionResourceModifyResponseItem.PDUSessionID.Value = 10
	// transfer := GetPDUSessionResourceModifyResponseTransfer()

	pDUSessionResourceModifyResponseItem.PDUSessionResourceModifyResponseTransfer =
		GetPDUSessionResourceModifyResponseTransfer()

	pDUSessionResourceModifyListModRes.List =
		append(pDUSessionResourceModifyListModRes.List, pDUSessionResourceModifyResponseItem)

	pDUSessionResourceModifyResponseIEs.List = append(pDUSessionResourceModifyResponseIEs.List, ie)

	// PDU Session Resource Failed to Modify List
	ie = ngapType.PDUSessionResourceModifyResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToModifyListModRes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentPDUSessionResourceFailedToModifyListModRes
	ie.Value.PDUSessionResourceFailedToModifyListModRes = new(ngapType.PDUSessionResourceFailedToModifyListModRes)

	pDUSessionResourceFailedToModifyListModRes := ie.Value.PDUSessionResourceFailedToModifyListModRes

	// PDU Session Resource Failed to Modify Item in PDU Session Resource Failed to Modify List
	pDUSessionResourceFailedToModifyItem := ngapType.PDUSessionResourceFailedToModifyItemModRes{}
	pDUSessionResourceFailedToModifyItem.PDUSessionID.Value = 10
	pDUSessionResourceFailedToModifyItem.PDUSessionResourceModifyUnsuccessfulTransfer =
		GetPDUSessionResourceModifyUnsuccessfulTransfer()

	pDUSessionResourceFailedToModifyListModRes.List =
		append(pDUSessionResourceFailedToModifyListModRes.List, pDUSessionResourceFailedToModifyItem)

	pDUSessionResourceModifyResponseIEs.List = append(pDUSessionResourceModifyResponseIEs.List, ie)

	// User Location Information (optional)
	ie = ngapType.PDUSessionResourceModifyResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x20},
		BitLength: 36,
	}

	userLocationInformationNR.TAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.TAI.TAC.Value = aper.OctetString("\x00\x00\x11")

	pDUSessionResourceModifyResponseIEs.List = append(pDUSessionResourceModifyResponseIEs.List, ie)

	// Criticality Diagnostics (optional)
	return pdu
}

func BuildPDUSessionResourceNotify() (pdu ngapType.NGAPPDU) {
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceNotify
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPDUSessionResourceNotify
	initiatingMessage.Value.PDUSessionResourceNotify = new(ngapType.PDUSessionResourceNotify)

	pDUSessionResourceNotify := initiatingMessage.Value.PDUSessionResourceNotify
	pDUSessionResourceNotifyIEs := &pDUSessionResourceNotify.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceNotifyIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceNotifyIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 0xffffffff

	pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)

	//PDU Session Resource Notify List
	ie = ngapType.PDUSessionResourceNotifyIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceNotifyList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentPDUSessionResourceNotifyList
	ie.Value.PDUSessionResourceNotifyList = new(ngapType.PDUSessionResourceNotifyList)

	pDUSessionResourceNotifyList := ie.Value.PDUSessionResourceNotifyList

	//PDU Session Resource Setup Request Item in (PDU Session Resource Setup Request List)
	pDUSessionResourceNotifyItem := ngapType.PDUSessionResourceNotifyItem{}
	pDUSessionResourceNotifyItem.PDUSessionID.Value = 10
	pDUSessionResourceNotifyItem.PDUSessionResourceNotifyTransfer = aper.OctetString("\x12\x34\x56")

	pDUSessionResourceNotifyList.List = append(pDUSessionResourceNotifyList.List, pDUSessionResourceNotifyItem)

	pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)

	//PDU Session Resource Released List
	ie = ngapType.PDUSessionResourceNotifyIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceReleasedListNot
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentPDUSessionResourceReleasedListNot
	ie.Value.PDUSessionResourceReleasedListNot = new(ngapType.PDUSessionResourceReleasedListNot)

	pDUSessionResourceReleasedListNot := ie.Value.PDUSessionResourceReleasedListNot

	//PDU Session Resource Released Item in (PDU Session Resource Released List)
	pDUSessionResourceReleasedItemNot := ngapType.PDUSessionResourceReleasedItemNot{}
	pDUSessionResourceReleasedItemNot.PDUSessionID.Value = 11
	pDUSessionResourceReleasedItemNot.PDUSessionResourceNotifyReleasedTransfer = aper.OctetString("\x65\x43\x21")

	pDUSessionResourceReleasedListNot.List =
		append(pDUSessionResourceReleasedListNot.List, pDUSessionResourceReleasedItemNot)

	pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)

	// User Location Information [optional]
	ie = ngapType.PDUSessionResourceNotifyIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}

	userLocationInformationNR.TAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	userLocationInformationNR.TAI.TAC.Value = aper.OctetString("\x00\x00\x11")

	pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)

	return pdu
}

func BuildPDUSessionResourceModifyIndication(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceModifyIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPDUSessionResourceModifyIndication
	initiatingMessage.Value.PDUSessionResourceModifyIndication = new(ngapType.PDUSessionResourceModifyIndication)

	pDUSessionResourceModifyIndication := initiatingMessage.Value.PDUSessionResourceModifyIndication
	pDUSessionResourceModifyIndicationIEs := &pDUSessionResourceModifyIndication.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceModifyIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	pDUSessionResourceModifyIndicationIEs.List = append(pDUSessionResourceModifyIndicationIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceModifyIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pDUSessionResourceModifyIndicationIEs.List = append(pDUSessionResourceModifyIndicationIEs.List, ie)

	// PDU Session Resource Modify Indication List
	ie = ngapType.PDUSessionResourceModifyIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentPDUSessionResourceModifyListModInd
	ie.Value.PDUSessionResourceModifyListModInd = new(ngapType.PDUSessionResourceModifyListModInd)

	pDUSessionResourceModifyListModInd := ie.Value.PDUSessionResourceModifyListModInd

	// PDU Session Resource Modify Indication Item (in PDU Session Resource Modify Indication List)
	pDUSessionResourceModifyItemModInd := ngapType.PDUSessionResourceModifyItemModInd{}
	pDUSessionResourceModifyItemModInd.PDUSessionID.Value = 10
	pDUSessionResourceModifyItemModInd.PDUSessionResourceModifyIndicationTransfer =
		GetPDUSessionResourceModifyIndicationTransfer()

	pDUSessionResourceModifyListModInd.List =
		append(pDUSessionResourceModifyListModInd.List, pDUSessionResourceModifyItemModInd)

	pDUSessionResourceModifyIndicationIEs.List = append(pDUSessionResourceModifyIndicationIEs.List, ie)

	return pdu
}

func BuildUEContextModificationFailure(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextModification
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentUEContextModificationFailure
	unsuccessfulOutcome.Value.UEContextModificationFailure = new(ngapType.UEContextModificationFailure)

	uEContextModificationFailure := unsuccessfulOutcome.Value.UEContextModificationFailure
	uEContextModificationFailureIEs := &uEContextModificationFailure.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	ie = ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentTransport
	cause.Transport = new(ngapType.CauseTransport)
	cause.Transport.Value = ngapType.CauseTransportPresentTransportResourceUnavailable

	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	// Criticality Diagnostics (optional)

	return pdu
}

func BuildRRCInactiveTransitionReport() (pdu ngapType.NGAPPDU) {
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeRRCInactiveTransitionReport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentRRCInactiveTransitionReport
	initiatingMessage.Value.RRCInactiveTransitionReport = new(ngapType.RRCInactiveTransitionReport)

	rRCInactiveTransitionReport := initiatingMessage.Value.RRCInactiveTransitionReport
	rRCInactiveTransitionReportIEs := &rRCInactiveTransitionReport.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.RRCInactiveTransitionReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.RRCInactiveTransitionReportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	rRCInactiveTransitionReportIEs.List = append(rRCInactiveTransitionReportIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.RRCInactiveTransitionReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.RRCInactiveTransitionReportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 0xFFFFFFFF

	rRCInactiveTransitionReportIEs.List = append(rRCInactiveTransitionReportIEs.List, ie)

	//RRC State
	ie = ngapType.RRCInactiveTransitionReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRRCState
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.RRCInactiveTransitionReportIEsPresentRRCState
	ie.Value.RRCState = new(ngapType.RRCState)

	rRCState := ie.Value.RRCState
	rRCState.Value = ngapType.RRCStatePresentConnected
	rRCInactiveTransitionReportIEs.List = append(rRCInactiveTransitionReportIEs.List, ie)

	//User Location Information
	ie = ngapType.RRCInactiveTransitionReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.RRCInactiveTransitionReportIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	//NR user location information

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity.Value = aper.OctetString("\x0f\x01\x22")
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}
	userLocationInformationNR.TAI.PLMNIdentity.Value = aper.OctetString("\x0f\x01\x22")
	userLocationInformationNR.TAI.TAC.Value = aper.OctetString("\x0f\x01\x22")
	//optional
	userLocationInformationNR.TimeStamp = new(ngapType.TimeStamp)
	userLocationInformationNR.TimeStamp.Value = aper.OctetString("\x0f\x01\x22\x21")

	//E-UTRA user location information
	/*
		userLocationInformation := ie.Value.UserLocationInformation
		userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationEUTRA
		userLocationInformation.UserLocationInformationEUTRA = new(ngapType.UserLocationInformationEUTRA)

		userLocationInformationEUTRA := userLocationInformation.UserLocationInformationEUTRA
		userLocationInformationEUTRA.EUTRACGI.EUTRACellIdentity.Value = aper.BitString{
			Bytes:     []byte{0x02, 0x42, 0x07, 0x30},
			BitLength: 28,
		}
		userLocationInformationEUTRA.EUTRACGI.PLMNIdentity.Value = aper.OctetString("\x0f\x01\x22")
		userLocationInformationEUTRA.TAI.PLMNIdentity.Value = aper.OctetString("\x0f\x01\x22")
		userLocationInformationEUTRA.TAI.TAC.Value = aper.OctetString("\x0f\x01\x22")

		//optional
		userLocationInformationEUTRA.TimeStamp = new(ngapType.TimeStamp)
		userLocationInformationEUTRA.TimeStamp.Value = aper.OctetString("\x0f\x01\x22\x21")
	*/

	rRCInactiveTransitionReportIEs.List = append(rRCInactiveTransitionReportIEs.List, ie)

	return pdu
}

func BuildHandoverNotify(amfUeNgapID int64, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeHandoverNotification
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentHandoverNotify
	initiatingMessage.Value.HandoverNotify = new(ngapType.HandoverNotify)

	handoverNotify := initiatingMessage.Value.HandoverNotify
	handoverNotifyIEs := &handoverNotify.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.HandoverNotifyIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverNotifyIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	handoverNotifyIEs.List = append(handoverNotifyIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.HandoverNotifyIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverNotifyIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	handoverNotifyIEs.List = append(handoverNotifyIEs.List, ie)

	//User Location Information
	ie = ngapType.HandoverNotifyIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverNotifyIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationEUTRA
	userLocationInformation.UserLocationInformationEUTRA = new(ngapType.UserLocationInformationEUTRA)

	userLocationInformationEUTRA := userLocationInformation.UserLocationInformationEUTRA
	userLocationInformationEUTRA.TAI.PLMNIdentity.Value = aper.OctetString("\x30\x33\x99")
	userLocationInformationEUTRA.TAI.TAC.Value = aper.OctetString("\x30\x33\x99")

	userLocationInformationEUTRA.EUTRACGI.PLMNIdentity.Value = aper.OctetString("\x30\x33\x99")
	userLocationInformationEUTRA.EUTRACGI.EUTRACellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x24, 0x16, 0x08, 0xFF},
		BitLength: 28,
	}

	handoverNotifyIEs.List = append(handoverNotifyIEs.List, ie)

	return pdu
}

func BuildUplinkRanStatusTransfer(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUplinkRANStatusTransfer
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUplinkRANStatusTransfer
	initiatingMessage.Value.UplinkRANStatusTransfer = new(ngapType.UplinkRANStatusTransfer)

	uplinkRANStatusTransfer := initiatingMessage.Value.UplinkRANStatusTransfer
	uplinkRANStatusTransferIEs := &uplinkRANStatusTransfer.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UplinkRANStatusTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkRANStatusTransferIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	uplinkRANStatusTransferIEs.List = append(uplinkRANStatusTransferIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UplinkRANStatusTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkRANStatusTransferIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	uplinkRANStatusTransferIEs.List = append(uplinkRANStatusTransferIEs.List, ie)

	// RAN Status Transfer Transparent Container
	ie = ngapType.UplinkRANStatusTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANStatusTransferTransparentContainer
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkRANStatusTransferIEsPresentRANStatusTransferTransparentContainer
	ie.Value.RANStatusTransferTransparentContainer = new(ngapType.RANStatusTransferTransparentContainer)

	rANStatusTransferTransparentContainer := ie.Value.RANStatusTransferTransparentContainer
	dRBsSubjectToStatusTransferList := &rANStatusTransferTransparentContainer.DRBsSubjectToStatusTransferList
	dRBsSubjectToStatusTransferItem := ngapType.DRBsSubjectToStatusTransferItem{}
	dRBsSubjectToStatusTransferItem.DRBID.Value = 123

	dRBStatusUL := &dRBsSubjectToStatusTransferItem.DRBStatusUL
	dRBStatusUL.Present = ngapType.DRBStatusULPresentDRBStatusUL12
	dRBStatusUL.DRBStatusUL12 = new(ngapType.DRBStatusUL12)

	dRBStatusUL12 := dRBStatusUL.DRBStatusUL12
	dRBStatusUL12.ULCOUNTValue.HFNPDCPSN12 = 345
	dRBStatusUL12.ULCOUNTValue.PDCPSN12 = 898

	dRBStatusDL := &dRBsSubjectToStatusTransferItem.DRBStatusDL
	dRBStatusDL.Present = ngapType.DRBStatusDLPresentDRBStatusDL12
	dRBStatusDL.DRBStatusDL12 = new(ngapType.DRBStatusDL12)

	dRBStatusDL12 := dRBStatusDL.DRBStatusDL12
	dRBStatusDL12.DLCOUNTValue.HFNPDCPSN12 = 987
	dRBStatusDL12.DLCOUNTValue.PDCPSN12 = 907

	dRBsSubjectToStatusTransferList.List = append(dRBsSubjectToStatusTransferList.List, dRBsSubjectToStatusTransferItem)
	uplinkRANStatusTransferIEs.List = append(uplinkRANStatusTransferIEs.List, ie)

	return pdu
}

func BuildNasNonDeliveryIndication(amfUeNgapID, ranUeNgapID int64, naspdu aper.OctetString) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeNASNonDeliveryIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentNASNonDeliveryIndication
	initiatingMessage.Value.NASNonDeliveryIndication = new(ngapType.NASNonDeliveryIndication)

	nasNonDeliveryIndication := initiatingMessage.Value.NASNonDeliveryIndication
	nasNonDeliveryIndicationIEs := &nasNonDeliveryIndication.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.NASNonDeliveryIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NASNonDeliveryIndicationIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	nasNonDeliveryIndicationIEs.List = append(nasNonDeliveryIndicationIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.NASNonDeliveryIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NASNonDeliveryIndicationIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	nasNonDeliveryIndicationIEs.List = append(nasNonDeliveryIndicationIEs.List, ie)

	// NAS-PDU
	ie = ngapType.NASNonDeliveryIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNASPDU
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NASNonDeliveryIndicationIEsPresentNASPDU
	ie.Value.NASPDU = new(ngapType.NASPDU)

	ie.Value.NASPDU.Value = naspdu

	nasNonDeliveryIndicationIEs.List = append(nasNonDeliveryIndicationIEs.List, ie)

	// Cause
	ie = ngapType.NASNonDeliveryIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NASNonDeliveryIndicationIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentRadioNetwork
	cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
	cause.RadioNetwork.Value = ngapType.CauseRadioNetworkPresentNgIntraSystemHandoverTriggered

	nasNonDeliveryIndicationIEs.List = append(nasNonDeliveryIndicationIEs.List, ie)

	return pdu
}

func BuildRanConfigurationUpdate() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeRANConfigurationUpdate
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentRANConfigurationUpdate
	initiatingMessage.Value.RANConfigurationUpdate = new(ngapType.RANConfigurationUpdate)

	rANConfigurationUpdate := initiatingMessage.Value.RANConfigurationUpdate
	rANConfigurationUpdateIEs := &rANConfigurationUpdate.ProtocolIEs

	// RanNodeName(optional)
	ie := ngapType.RANConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANNodeName
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.RANConfigurationUpdateIEsPresentRANNodeName
	ie.Value.RANNodeName = new(ngapType.RANNodeName)

	rANNodeName := ie.Value.RANNodeName
	rANNodeName.Value = "Free5GC"
	rANConfigurationUpdateIEs.List = append(rANConfigurationUpdateIEs.List, ie)

	// SupportTAList
	ie = ngapType.RANConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSupportedTAList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.RANConfigurationUpdateIEsPresentSupportedTAList
	ie.Value.SupportedTAList = new(ngapType.SupportedTAList)

	supportedTAList := ie.Value.SupportedTAList
	// SupportTAItem in SupportTAlist
	supportedTAItem := ngapType.SupportedTAItem{}
	supportedTAItem.TAC.Value = aper.OctetString("\x00\x00\x01")

	broadcastPLMNList := &supportedTAItem.BroadcastPLMNList
	// BroadcastPLMNItem in BroadcastPLMNList
	broadcastPLMNLItem := ngapType.BroadcastPLMNItem{}
	broadcastPLMNLItem.PLMNIdentity.Value = aper.OctetString("\x00\x1D\x5C")

	sliceSupportList := &broadcastPLMNLItem.TAISliceSupportList
	// SlicesupportItem in SliceSupportList
	sliceSupportItem := ngapType.SliceSupportItem{}
	sliceSupportItem.SNSSAI.SST.Value = aper.OctetString("\x57")
	// Optional
	sliceSupportItem.SNSSAI.SD = new(ngapType.SD)
	sliceSupportItem.SNSSAI.SD.Value = aper.OctetString("\x00\x01\x02")

	sliceSupportList.List = append(sliceSupportList.List, sliceSupportItem)
	broadcastPLMNList.List = append(broadcastPLMNList.List, broadcastPLMNLItem)
	supportedTAList.List = append(supportedTAList.List, supportedTAItem)
	rANConfigurationUpdateIEs.List = append(rANConfigurationUpdateIEs.List, ie)

	// DefaultPagingDRX
	ie = ngapType.RANConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDDefaultPagingDRX
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.RANConfigurationUpdateIEsPresentDefaultPagingDRX
	ie.Value.DefaultPagingDRX = new(ngapType.PagingDRX)

	pagingDRX := ie.Value.DefaultPagingDRX
	pagingDRX.Value = ngapType.PagingDRXPresentV128
	rANConfigurationUpdateIEs.List = append(rANConfigurationUpdateIEs.List, ie)

	return pdu
}

func BuildRanConfigurationUpdateAck(diagnostics *ngapType.CriticalityDiagnostics) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeRANConfigurationUpdate
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentRANConfigurationUpdateAcknowledge
	successfulOutcome.Value.RANConfigurationUpdateAcknowledge = new(ngapType.RANConfigurationUpdateAcknowledge)

	rANConfigurationUpdateAcknowledge := successfulOutcome.Value.RANConfigurationUpdateAcknowledge
	rANConfigurationUpdateAcknowledgeIEs := &rANConfigurationUpdateAcknowledge.ProtocolIEs
	// CriticalityDiagnostics
	if diagnostics != nil {
		ie := ngapType.RANConfigurationUpdateAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.RANConfigurationUpdateAcknowledgeIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		criticalityDiagnostics := ie.Value.CriticalityDiagnostics
		*criticalityDiagnostics = *diagnostics

		rANConfigurationUpdateAcknowledgeIEs.List = append(rANConfigurationUpdateAcknowledgeIEs.List, ie)
	}

	return pdu
}

func BuildRanConfigurationUpdateFailure(
	time *ngapType.TimeToWait,
	diagnostics *ngapType.CriticalityDiagnostics) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeRANConfigurationUpdate
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentRANConfigurationUpdateFailure
	unsuccessfulOutcome.Value.RANConfigurationUpdateFailure = new(ngapType.RANConfigurationUpdateFailure)

	rANConfigurationUpdateFailure := unsuccessfulOutcome.Value.RANConfigurationUpdateFailure
	rANConfigurationUpdateFailureIEs := &rANConfigurationUpdateFailure.ProtocolIEs
	// Cause
	{
		ie := ngapType.RANConfigurationUpdateFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCause
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.RANConfigurationUpdateFailureIEsPresentCause
		ie.Value.Cause = new(ngapType.Cause)

		cause := ie.Value.Cause
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentControlProcessingOverload,
		}

		rANConfigurationUpdateFailureIEs.List = append(rANConfigurationUpdateFailureIEs.List, ie)
	}
	// TimeToWait
	if time != nil {
		ie := ngapType.RANConfigurationUpdateFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDTimeToWait
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.RANConfigurationUpdateFailureIEsPresentTimeToWait
		ie.Value.TimeToWait = new(ngapType.TimeToWait)

		timeToWait := ie.Value.TimeToWait
		*timeToWait = *time

		rANConfigurationUpdateFailureIEs.List = append(rANConfigurationUpdateFailureIEs.List, ie)
	}
	// CriticalityDiagnostics
	if diagnostics != nil {
		ie := ngapType.RANConfigurationUpdateFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.RANConfigurationUpdateFailureIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		criticalityDiagnostics := ie.Value.CriticalityDiagnostics
		*criticalityDiagnostics = *diagnostics

		rANConfigurationUpdateFailureIEs.List = append(rANConfigurationUpdateFailureIEs.List, ie)
	}

	return pdu
}

func BuildAMFStatusIndication() (pdu ngapType.NGAPPDU) {
	return pdu
}

func BuildUplinkRanConfigurationTransfer() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUplinkRANConfigurationTransfer
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUplinkRANConfigurationTransfer
	initiatingMessage.Value.UplinkRANConfigurationTransfer = new(ngapType.UplinkRANConfigurationTransfer)

	uplinkRANConfigurationTransfer := initiatingMessage.Value.UplinkRANConfigurationTransfer
	uplinkRANConfigurationTransferIEs := &uplinkRANConfigurationTransfer.ProtocolIEs

	// SON Configuration Transfer [optional]
	ie := ngapType.UplinkRANConfigurationTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSONConfigurationTransferUL
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UplinkRANConfigurationTransferIEsPresentSONConfigurationTransferUL
	ie.Value.SONConfigurationTransferUL = new(ngapType.SONConfigurationTransfer)

	sONConfigurationTransferUL := ie.Value.SONConfigurationTransferUL

	// Target Ran Node ID in (SON Configuration Transfer)
	targetRANNodeID := &sONConfigurationTransferUL.TargetRANNodeID
	targetRANNodeID.GlobalRANNodeID.Present = ngapType.GlobalRANNodeIDPresentGlobalGNBID
	targetRANNodeID.GlobalRANNodeID.GlobalGNBID = new(ngapType.GlobalGNBID)
	targetRANNodeID.GlobalRANNodeID.GlobalGNBID.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	targetRANNodeID.GlobalRANNodeID.GlobalGNBID.GNBID.Present = ngapType.GNBIDPresentGNBID
	targetRANNodeID.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID = new(aper.BitString)

	gNBID := targetRANNodeID.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID
	*gNBID = aper.BitString{
		Bytes:     []byte{0x41, 0x42, 0x40},
		BitLength: 22,
	}
	targetRANNodeID.SelectedTAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	targetRANNodeID.SelectedTAI.TAC.Value = aper.OctetString("\x00\x00\x01")

	// Source Ran Node ID in (SON Configuration Transfer)
	sourceRANNodeID := &sONConfigurationTransferUL.SourceRANNodeID
	sourceRANNodeID.GlobalRANNodeID.Present = ngapType.GlobalRANNodeIDPresentGlobalGNBID
	sourceRANNodeID.GlobalRANNodeID.GlobalGNBID = new(ngapType.GlobalGNBID)
	sourceRANNodeID.GlobalRANNodeID.GlobalGNBID.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	sourceRANNodeID.GlobalRANNodeID.GlobalGNBID.GNBID.Present = ngapType.GNBIDPresentGNBID
	sourceRANNodeID.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID = new(aper.BitString)

	gNBID = sourceRANNodeID.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID
	*gNBID = aper.BitString{
		Bytes:     []byte{0x45, 0x46, 0x47},
		BitLength: 24,
	}
	sourceRANNodeID.SelectedTAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	sourceRANNodeID.SelectedTAI.TAC.Value = aper.OctetString("\x00\x00\x01")

	// SON Information in (SON Configuration Transfer)
	sONInformation := &sONConfigurationTransferUL.SONInformation
	sONInformation.Present = ngapType.SONInformationPresentSONInformationRequest
	sONInformation.SONInformationRequest = new(ngapType.SONInformationRequest)
	sONInformation.SONInformationRequest.Value = ngapType.SONInformationRequestPresentXnTNLConfigurationInfo
	// sONInformation.SONInformationReply.XnTNLConfigurationInfo = new(ngapType.XnTNLConfigurationInfo)

	// xnTNLConfigurationInfo := sONInformation.SONInformationReply.XnTNLConfigurationInfo

	// Xn TNL Configuration Info [C-ifSONInformationRequest]
	xnTNLConfigurationInfo := sONConfigurationTransferUL.XnTNLConfigurationInfo
	xnTransportLayerAddresses := xnTNLConfigurationInfo.XnTransportLayerAddresses

	TLA := ngapType.TransportLayerAddress{}
	TLA.Value = aper.BitString{
		Bytes:     []byte{0x12, 0x34, 0x56, 0x78},
		BitLength: 32,
	}
	xnTransportLayerAddresses.List = append(xnTransportLayerAddresses.List, TLA)

	uplinkRANConfigurationTransferIEs.List = append(uplinkRANConfigurationTransferIEs.List, ie)

	return pdu
}

func BuildUplinkUEAssociatedNRPPATransport() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUplinkUEAssociatedNRPPaTransport
	initiatingMessage.Value.UplinkUEAssociatedNRPPaTransport = new(ngapType.UplinkUEAssociatedNRPPaTransport)

	uplinkUEAssociatedNRPPaTransport := initiatingMessage.Value.UplinkUEAssociatedNRPPaTransport
	uplinkUEAssociatedNRPPaTransportIEs := &uplinkUEAssociatedNRPPaTransport.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.UplinkUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkUEAssociatedNRPPaTransportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	uplinkUEAssociatedNRPPaTransportIEs.List = append(uplinkUEAssociatedNRPPaTransportIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.UplinkUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkUEAssociatedNRPPaTransportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 4294967295

	uplinkUEAssociatedNRPPaTransportIEs.List = append(uplinkUEAssociatedNRPPaTransportIEs.List, ie)

	//Routing ID
	ie = ngapType.UplinkUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRoutingID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkUEAssociatedNRPPaTransportIEsPresentRoutingID
	ie.Value.RoutingID = new(ngapType.RoutingID)

	routingID := ie.Value.RoutingID
	routingID.Value = aper.OctetString("\x03\x02")

	uplinkUEAssociatedNRPPaTransportIEs.List = append(uplinkUEAssociatedNRPPaTransportIEs.List, ie)

	//NRPPa-PDU
	ie = ngapType.UplinkUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNRPPaPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkUEAssociatedNRPPaTransportIEsPresentNRPPaPDU
	ie.Value.NRPPaPDU = new(ngapType.NRPPaPDU)

	nRPPaPDU := ie.Value.NRPPaPDU
	nRPPaPDU.Value = aper.OctetString("\x03\x02")

	uplinkUEAssociatedNRPPaTransportIEs.List = append(uplinkUEAssociatedNRPPaTransportIEs.List, ie)

	return pdu
}

func BuildUplinkNonUEAssociatedNRPPATransport() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUplinkNonUEAssociatedNRPPaTransport
	initiatingMessage.Value.UplinkNonUEAssociatedNRPPaTransport = new(ngapType.UplinkNonUEAssociatedNRPPaTransport)

	uplinkNonUEAssociatedNRPPaTransport := initiatingMessage.Value.UplinkNonUEAssociatedNRPPaTransport
	uplinkNonUEAssociatedNRPPaTransportIEs := &uplinkNonUEAssociatedNRPPaTransport.ProtocolIEs

	//Routing ID
	ie := ngapType.UplinkNonUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRoutingID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNonUEAssociatedNRPPaTransportIEsPresentRoutingID
	ie.Value.RoutingID = new(ngapType.RoutingID)
	routingID := ie.Value.RoutingID
	routingID.Value = aper.OctetString("\x03\x02")

	uplinkNonUEAssociatedNRPPaTransportIEs.List = append(uplinkNonUEAssociatedNRPPaTransportIEs.List, ie)

	//NRPPa-PDU
	ie = ngapType.UplinkNonUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNRPPaPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNonUEAssociatedNRPPaTransportIEsPresentNRPPaPDU
	ie.Value.NRPPaPDU = new(ngapType.NRPPaPDU)
	nRPPaPDU := ie.Value.NRPPaPDU
	nRPPaPDU.Value = aper.OctetString("\x03\x02")

	uplinkNonUEAssociatedNRPPaTransportIEs.List = append(uplinkNonUEAssociatedNRPPaTransportIEs.List, ie)
	return pdu
}

func BuildLocationReport() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeLocationReport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentLocationReport
	initiatingMessage.Value.LocationReport = new(ngapType.LocationReport)

	locationReport := initiatingMessage.Value.LocationReport
	locationReportIEs := &locationReport.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.LocationReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.LocationReportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	locationReportIEs.List = append(locationReportIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.LocationReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.LocationReportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 0xffffffff

	locationReportIEs.List = append(locationReportIEs.List, ie)

	//User Location Information
	ie = ngapType.LocationReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.LocationReportIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationEUTRA
	userLocationInformation.UserLocationInformationEUTRA = new(ngapType.UserLocationInformationEUTRA)

	userLocationInformationEUTRA := userLocationInformation.UserLocationInformationEUTRA
	userLocationInformationEUTRA.EUTRACGI.PLMNIdentity.Value = aper.OctetString("\x53\x54\x55")
	userLocationInformationEUTRA.EUTRACGI.EUTRACellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x10, 0x11, 0x12, 0x13},
		BitLength: 28,
	}

	userLocationInformationEUTRA.TAI.PLMNIdentity.Value = aper.OctetString("\x53\x54\x55")
	userLocationInformationEUTRA.TAI.TAC.Value = aper.OctetString("\x53\x54\x55")

	locationReportIEs.List = append(locationReportIEs.List, ie)

	//UE Presence in Area of Interest List [optional]	// if EventType = ngapType.EventTypePresentUePresenceInAreaOfInterest
	ie = ngapType.LocationReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUEPresenceInAreaOfInterestList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.LocationReportIEsPresentUEPresenceInAreaOfInterestList
	ie.Value.UEPresenceInAreaOfInterestList = new(ngapType.UEPresenceInAreaOfInterestList)
	uEPresenceInAreaOfInterestList := ie.Value.UEPresenceInAreaOfInterestList

	uEPresenceInAreaOfInterestItem := ngapType.UEPresenceInAreaOfInterestItem{}
	uEPresenceInAreaOfInterestItem.LocationReportingReferenceID.Value = 23
	uEPresenceInAreaOfInterestItem.UEPresence.Value = ngapType.UEPresencePresentIn

	uEPresenceInAreaOfInterestList.List = append(uEPresenceInAreaOfInterestList.List, uEPresenceInAreaOfInterestItem)
	locationReportIEs.List = append(locationReportIEs.List, ie)

	//Location Reporting Request Type
	ie = ngapType.LocationReportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDLocationReportingRequestType
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.LocationReportIEsPresentLocationReportingRequestType
	ie.Value.LocationReportingRequestType = new(ngapType.LocationReportingRequestType)

	locationReportingRequestType := ie.Value.LocationReportingRequestType
	locationReportingRequestType.EventType.Value = ngapType.EventTypePresentStopUePresenceInAreaOfInterest
	locationReportingRequestType.ReportArea.Value = ngapType.ReportAreaPresentCell

	locationReportingRequestType.AreaOfInterestList = new(ngapType.AreaOfInterestList)
	AOIList := locationReportingRequestType.AreaOfInterestList

	areaOfInterestItem := ngapType.AreaOfInterestItem{}

	areaOfInterestItem.AreaOfInterest.AreaOfInterestTAIList = new(ngapType.AreaOfInterestTAIList)
	AOITAIList := areaOfInterestItem.AreaOfInterest.AreaOfInterestTAIList

	areaOfInterestTAIItem := ngapType.AreaOfInterestTAIItem{}
	areaOfInterestTAIItem.TAI.PLMNIdentity.Value = aper.OctetString("\x53\x54\x55")
	areaOfInterestTAIItem.TAI.TAC.Value = aper.OctetString("\x53\x54\x55")
	AOITAIList.List = append(AOITAIList.List, areaOfInterestTAIItem)

	areaOfInterestItem.AreaOfInterest.AreaOfInterestCellList = new(ngapType.AreaOfInterestCellList)
	AOICellList := areaOfInterestItem.AreaOfInterest.AreaOfInterestCellList

	areaOfInterestCellItem := ngapType.AreaOfInterestCellItem{}
	areaOfInterestCellItem.NGRANCGI.Present = ngapType.NGRANCGIPresentEUTRACGI
	areaOfInterestCellItem.NGRANCGI.EUTRACGI = new(ngapType.EUTRACGI)
	areaOfInterestCellItem.NGRANCGI.EUTRACGI.PLMNIdentity.Value = aper.OctetString("\x53\x54\x55")
	areaOfInterestCellItem.NGRANCGI.EUTRACGI.EUTRACellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x10, 0x11, 0x12, 0x13},
		BitLength: 28,
	}
	AOICellList.List = append(AOICellList.List, areaOfInterestCellItem)

	areaOfInterestItem.AreaOfInterest.AreaOfInterestRANNodeList = new(ngapType.AreaOfInterestRANNodeList)
	AOIRanNodeList := areaOfInterestItem.AreaOfInterest.AreaOfInterestRANNodeList

	areaOfInterestRANNodeItem := ngapType.AreaOfInterestRANNodeItem{}
	areaOfInterestRANNodeItem.GlobalRANNodeID.Present = ngapType.GlobalRANNodeIDPresentGlobalGNBID
	areaOfInterestRANNodeItem.GlobalRANNodeID.GlobalGNBID = new(ngapType.GlobalGNBID)
	areaOfInterestRANNodeItem.GlobalRANNodeID.GlobalGNBID.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	areaOfInterestRANNodeItem.GlobalRANNodeID.GlobalGNBID.GNBID.Present = ngapType.GNBIDPresentGNBID
	areaOfInterestRANNodeItem.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID = new(aper.BitString)

	gNBID := areaOfInterestRANNodeItem.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID
	*gNBID = aper.BitString{
		Bytes:     []byte{0x45, 0x46, 0x47},
		BitLength: 24,
	}
	AOIRanNodeList.List = append(AOIRanNodeList.List, areaOfInterestRANNodeItem)

	areaOfInterestItem.LocationReportingReferenceID.Value = 23
	AOIList.List = append(AOIList.List, areaOfInterestItem)

	// Location Reporting Reference ID to be Cancelled in (Location Reporting Request Type)
	// [C- ifEventTypeisStopUEPresinAoI]
	locationReportingRequestType.LocationReportingReferenceIDToBeCancelled = new(ngapType.LocationReportingReferenceID)
	locationReportingRequestType.LocationReportingReferenceIDToBeCancelled.Value = 56

	locationReportIEs.List = append(locationReportIEs.List, ie)

	return pdu
}

func BuildUETNLABindingReleaseRequest() (pdu ngapType.NGAPPDU) {
	return pdu
}

func BuildUERadioCapabilityInfoIndication() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUERadioCapabilityInfoIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUERadioCapabilityInfoIndication
	initiatingMessage.Value.UERadioCapabilityInfoIndication = new(ngapType.UERadioCapabilityInfoIndication)

	uERadioCapabilityInfoIndication := initiatingMessage.Value.UERadioCapabilityInfoIndication
	uERadioCapabilityInfoIndicationIEs := &uERadioCapabilityInfoIndication.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.UERadioCapabilityInfoIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UERadioCapabilityInfoIndicationIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = 1

	uERadioCapabilityInfoIndicationIEs.List = append(uERadioCapabilityInfoIndicationIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.UERadioCapabilityInfoIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UERadioCapabilityInfoIndicationIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = 4294967295

	uERadioCapabilityInfoIndicationIEs.List = append(uERadioCapabilityInfoIndicationIEs.List, ie)

	//UE Radio Capability
	ie = ngapType.UERadioCapabilityInfoIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUERadioCapability
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UERadioCapabilityInfoIndicationIEsPresentUERadioCapability
	ie.Value.UERadioCapability = new(ngapType.UERadioCapability)

	uERadioCapability := ngapType.UERadioCapability{}
	uERadioCapability.Value = aper.OctetString("\x00\x00\x01")

	uERadioCapabilityInfoIndicationIEs.List = append(uERadioCapabilityInfoIndicationIEs.List, ie)

	//	TODO: UE Radio Capability for Paging (optional)
	return pdu
}

func BuildAMFConfigurationUpdateAcknowledge() (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeAMFConfigurationUpdate
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentAMFConfigurationUpdateAcknowledge
	successfulOutcome.Value.AMFConfigurationUpdateAcknowledge = new(ngapType.AMFConfigurationUpdateAcknowledge)

	AMFConfigurationUpdateAcknowledge := successfulOutcome.Value.AMFConfigurationUpdateAcknowledge
	AMFConfigurationUpdateAcknowledgeIEs := &AMFConfigurationUpdateAcknowledge.ProtocolIEs

	//	AMF TNL Association Setup List
	ie := ngapType.AMFConfigurationUpdateAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationSetupList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.AMFConfigurationUpdateAcknowledgeIEsPresentAMFTNLAssociationSetupList
	ie.Value.AMFTNLAssociationSetupList = new(ngapType.AMFTNLAssociationSetupList)

	aMFTNLAssociationSetupList := ie.Value.AMFTNLAssociationSetupList

	//	AMF TNL Association Setup Item
	aMFTNLAssociationSetupItem := ngapType.AMFTNLAssociationSetupItem{}
	aMFTNLAssociationSetupItem.AMFTNLAssociationAddress.Present =
		ngapType.CPTransportLayerInformationPresentEndpointIPAddress
	aMFTNLAssociationSetupItem.AMFTNLAssociationAddress.EndpointIPAddress = new(ngapType.TransportLayerAddress)
	aMFTNLAssociationSetupItem.AMFTNLAssociationAddress.EndpointIPAddress.Value = aper.BitString{
		Bytes:     []byte{0x12, 0x34, 0x56, 0x78},
		BitLength: 32,
	}

	aMFTNLAssociationSetupList.List = append(aMFTNLAssociationSetupList.List, aMFTNLAssociationSetupItem)
	AMFConfigurationUpdateAcknowledgeIEs.List = append(AMFConfigurationUpdateAcknowledgeIEs.List, ie)

	//	AMF TNL Association Failed to Setup List (optional)
	ie = ngapType.AMFConfigurationUpdateAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationFailedToSetupList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.AMFConfigurationUpdateAcknowledgeIEsPresentAMFTNLAssociationFailedToSetupList
	ie.Value.Present = ngapType.AMFConfigurationUpdateAcknowledgeIEsPresentAMFTNLAssociationFailedToSetupList
	ie.Value.AMFTNLAssociationFailedToSetupList = new(ngapType.TNLAssociationList)

	aMFTNLAssociationFailedToSetupList := ie.Value.AMFTNLAssociationFailedToSetupList

	//	TNLAssociationItem
	tNLAssociationItem := ngapType.TNLAssociationItem{}
	tNLAssociationItem.Cause.Present = ngapType.CausePresentMisc
	tNLAssociationItem.Cause.Misc = new(ngapType.CauseMisc)
	tNLAssociationItem.Cause.Misc.Value = ngapType.CauseMiscPresentUnspecified
	tNLAssociationItem.TNLAssociationAddress.Present = ngapType.CPTransportLayerInformationPresentEndpointIPAddress
	tNLAssociationItem.TNLAssociationAddress.EndpointIPAddress = new(ngapType.TransportLayerAddress)
	tNLAssociationItem.TNLAssociationAddress.EndpointIPAddress.Value = aper.BitString{
		Bytes:     []byte{0x12, 0x34, 0x56, 0x78},
		BitLength: 32,
	}

	aMFTNLAssociationFailedToSetupList.List = append(aMFTNLAssociationFailedToSetupList.List, tNLAssociationItem)
	AMFConfigurationUpdateAcknowledgeIEs.List = append(AMFConfigurationUpdateAcknowledgeIEs.List, ie)

	//	Criticality Diagnostics (optional)

	return pdu
}

func BuildAMFConfigurationUpdate(amfName string, guamiList []ngapType.ServedGUAMIItem,
	plmnList []ngapType.PLMNSupportItem, amfRelativeCapacity int64,
	addList *ngapType.AMFTNLAssociationToAddList, removeList *ngapType.AMFTNLAssociationToRemoveList,
	updateList *ngapType.AMFTNLAssociationToUpdateList) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeAMFConfigurationUpdate
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentAMFConfigurationUpdate
	initiatingMessage.Value.AMFConfigurationUpdate = new(ngapType.AMFConfigurationUpdate)

	aMFConfigurationUpdate := initiatingMessage.Value.AMFConfigurationUpdate
	aMFConfigurationUpdateIEs := &aMFConfigurationUpdate.ProtocolIEs
	// AMFName
	{
		ie := ngapType.AMFConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFName
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentAMFName
		ie.Value.AMFName = new(ngapType.AMFName)

		aMFName := ie.Value.AMFName
		aMFName.Value = amfName

		aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)
	}
	// ServedGUAMIList
	{
		ie := ngapType.AMFConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDServedGUAMIList
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentServedGUAMIList
		ie.Value.ServedGUAMIList = new(ngapType.ServedGUAMIList)

		servedGUAMIList := ie.Value.ServedGUAMIList
		servedGUAMIList.List = guamiList

		aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)
	}
	// RelativeAMFCapacity
	{
		ie := ngapType.AMFConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRelativeAMFCapacity
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentRelativeAMFCapacity
		ie.Value.RelativeAMFCapacity = new(ngapType.RelativeAMFCapacity)

		relativeAMFCapacity := ie.Value.RelativeAMFCapacity
		relativeAMFCapacity.Value = amfRelativeCapacity

		aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)
	}
	// PLMNSupportList
	{
		ie := ngapType.AMFConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPLMNSupportList
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentPLMNSupportList
		ie.Value.PLMNSupportList = new(ngapType.PLMNSupportList)

		pLMNSupportList := ie.Value.PLMNSupportList
		pLMNSupportList.List = plmnList

		aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)
	}
	// AMFTNLAssociationToAddList
	if addList != nil {
		ie := ngapType.AMFConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationToAddList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentAMFTNLAssociationToAddList
		ie.Value.AMFTNLAssociationToAddList = new(ngapType.AMFTNLAssociationToAddList)

		aMFTNLAssociationToAddList := ie.Value.AMFTNLAssociationToAddList
		*aMFTNLAssociationToAddList = *addList

		aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)
	}
	// AMFTNLAssociationToRemoveList
	if removeList != nil {
		ie := ngapType.AMFConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationToRemoveList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentAMFTNLAssociationToRemoveList
		ie.Value.AMFTNLAssociationToRemoveList = new(ngapType.AMFTNLAssociationToRemoveList)

		aMFTNLAssociationToRemoveList := ie.Value.AMFTNLAssociationToRemoveList
		*aMFTNLAssociationToRemoveList = *removeList

		aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)
	}
	// AMFTNLAssociationToUpdateList
	if updateList != nil {
		ie := ngapType.AMFConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationToUpdateList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentAMFTNLAssociationToUpdateList
		ie.Value.AMFTNLAssociationToUpdateList = new(ngapType.AMFTNLAssociationToUpdateList)

		aMFTNLAssociationToUpdateList := ie.Value.AMFTNLAssociationToUpdateList
		*aMFTNLAssociationToUpdateList = *updateList

		aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)
	}

	return pdu
}

func BuildHandoverRequired(
	amfUeNgapID, ranUeNgapID int64, targetGNBID []byte, targetCellID []byte) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeHandoverPreparation
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentHandoverRequired
	initiatingMessage.Value.HandoverRequired = new(ngapType.HandoverRequired)

	handoverRequired := initiatingMessage.Value.HandoverRequired
	handoverRequiredIEs := &handoverRequired.ProtocolIEs

	//AMF UE NGAP ID
	ie := ngapType.HandoverRequiredIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequiredIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	handoverRequiredIEs.List = append(handoverRequiredIEs.List, ie)

	//RAN UE NGAP ID
	ie = ngapType.HandoverRequiredIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequiredIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	handoverRequiredIEs.List = append(handoverRequiredIEs.List, ie)

	// Handover Type
	ie = ngapType.HandoverRequiredIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDHandoverType
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequiredIEsPresentHandoverType
	ie.Value.HandoverType = new(ngapType.HandoverType)

	handoverType := ie.Value.HandoverType
	handoverType.Value = ngapType.HandoverTypePresentIntra5gs

	handoverRequiredIEs.List = append(handoverRequiredIEs.List, ie)

	//Cause
	ie = ngapType.HandoverRequiredIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverRequiredIEsPresentCause
	ie.Value.Cause = new(ngapType.Cause)

	cause := ie.Value.Cause
	cause.Present = ngapType.CausePresentRadioNetwork
	cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
	cause.RadioNetwork.Value = ngapType.CauseRadioNetworkPresentHandoverDesirableForRadioReason

	handoverRequiredIEs.List = append(handoverRequiredIEs.List, ie)

	//Target ID
	ie = ngapType.HandoverRequiredIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDTargetID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequiredIEsPresentTargetID
	ie.Value.TargetID = new(ngapType.TargetID)

	targetID := ie.Value.TargetID
	targetID.Present = ngapType.TargetIDPresentTargetRANNodeID
	targetID.TargetRANNodeID = new(ngapType.TargetRANNodeID)

	targetRANNodeID := targetID.TargetRANNodeID
	targetRANNodeID.GlobalRANNodeID.Present = ngapType.GlobalRANNodeIDPresentGlobalGNBID
	targetRANNodeID.GlobalRANNodeID.GlobalGNBID = new(ngapType.GlobalGNBID)

	globalRANNodeID := targetRANNodeID.GlobalRANNodeID
	globalRANNodeID.GlobalGNBID.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	globalRANNodeID.GlobalGNBID.GNBID.Present = ngapType.GNBIDPresentGNBID

	globalRANNodeID.GlobalGNBID.GNBID.GNBID = new(aper.BitString)

	gNBID := globalRANNodeID.GlobalGNBID.GNBID.GNBID

	*gNBID = aper.BitString{
		Bytes:     targetGNBID,
		BitLength: uint64(len(targetGNBID) * 8),
	}
	globalRANNodeID.GlobalGNBID.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")

	targetRANNodeID.SelectedTAI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	targetRANNodeID.SelectedTAI.TAC.Value = aper.OctetString("\x30\x33\x99")

	handoverRequiredIEs.List = append(handoverRequiredIEs.List, ie)

	// Direct Forwarding Path Availability [optional]

	// PDU Session Resource List
	ie = ngapType.HandoverRequiredIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceListHORqd
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequiredIEsPresentPDUSessionResourceListHORqd
	ie.Value.PDUSessionResourceListHORqd = new(ngapType.PDUSessionResourceListHORqd)

	pDUSessionResourceListHORqd := ie.Value.PDUSessionResourceListHORqd

	//PDU Session Resource Item (in PDU Session Resource List)
	pDUSessionResourceItem := ngapType.PDUSessionResourceItemHORqd{}
	pDUSessionResourceItem.PDUSessionID.Value = 10
	pDUSessionResourceItem.HandoverRequiredTransfer = GetHandoverRequiredTransfer()

	pDUSessionResourceListHORqd.List = append(pDUSessionResourceListHORqd.List, pDUSessionResourceItem)

	handoverRequiredIEs.List = append(handoverRequiredIEs.List, ie)

	//Source to Target Transparent Container
	ie = ngapType.HandoverRequiredIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSourceToTargetTransparentContainer
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequiredIEsPresentSourceToTargetTransparentContainer
	ie.Value.SourceToTargetTransparentContainer = new(ngapType.SourceToTargetTransparentContainer)

	ie.Value.SourceToTargetTransparentContainer.Value = GetSourceToTargetTransparentTransfer(targetGNBID, targetCellID)

	handoverRequiredIEs.List = append(handoverRequiredIEs.List, ie)

	return pdu
}

func BuildCellTrafficTrace(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeCellTrafficTrace
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentCellTrafficTrace
	initiatingMessage.Value.CellTrafficTrace = new(ngapType.CellTrafficTrace)

	cellTrafficTrace := initiatingMessage.Value.CellTrafficTrace
	cellTrafficTraceIEs := &cellTrafficTrace.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.CellTrafficTraceIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.CellTrafficTraceIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	cellTrafficTraceIEs.List = append(cellTrafficTraceIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.CellTrafficTraceIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.CellTrafficTraceIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	cellTrafficTraceIEs.List = append(cellTrafficTraceIEs.List, ie)

	// NG-RAN Trace ID
	ie = ngapType.CellTrafficTraceIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNGRANTraceID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.CellTrafficTraceIEsPresentNGRANTraceID
	ie.Value.NGRANTraceID = new(ngapType.NGRANTraceID)

	ie.Value.NGRANTraceID.Value = aper.OctetString("\x02\xf8\x39\x11\x22\x33\x00\x01")

	cellTrafficTraceIEs.List = append(cellTrafficTraceIEs.List, ie)

	// NG-RAN CGI
	ie = ngapType.CellTrafficTraceIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNGRANCGI
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.CellTrafficTraceIEsPresentNGRANCGI
	ie.Value.NGRANCGI = new(ngapType.NGRANCGI)

	nGRANCGI := ie.Value.NGRANCGI
	nGRANCGI.Present = ngapType.NGRANCGIPresentNRCGI
	nGRANCGI.NRCGI = new(ngapType.NRCGI)
	nGRANCGI.NRCGI.PLMNIdentity.Value = aper.OctetString("\x02\xf8\x39")
	nGRANCGI.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}

	cellTrafficTraceIEs.List = append(cellTrafficTraceIEs.List, ie)

	// Trace Collection Entity IP Address
	ie = ngapType.CellTrafficTraceIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDTraceCollectionEntityIPAddress
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.CellTrafficTraceIEsPresentTraceCollectionEntityIPAddress
	ie.Value.TraceCollectionEntityIPAddress = new(ngapType.TransportLayerAddress)

	ie.Value.TraceCollectionEntityIPAddress.Value = aper.BitString{
		Bytes:     []byte{0x7f, 0x00, 0x00, 0x01},
		BitLength: 32,
	}

	cellTrafficTraceIEs.List = append(cellTrafficTraceIEs.List, ie)

	return pdu
}

func buildPDUSessionResourceSetupResponseTransfer(ipv4 string) (data ngapType.PDUSessionResourceSetupResponseTransfer) {

	// QoS Flow per TNL Information
	qosFlowPerTNLInformation := &data.DLQosFlowPerTNLInformation
	qosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

	// UP Transport Layer Information in QoS Flow per TNL Information
	upTransportLayerInformation := &qosFlowPerTNLInformation.UPTransportLayerInformation
	upTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	upTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
	upTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString("\x00\x00\x00\x01")
	upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(ipv4, "")

	// Associated QoS Flow List in QoS Flow per TNL Information
	associatedQosFlowList := &qosFlowPerTNLInformation.AssociatedQosFlowList

	associatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
	associatedQosFlowItem.QosFlowIdentifier.Value = 1
	associatedQosFlowList.List = append(associatedQosFlowList.List, associatedQosFlowItem)

	return data
}

func buildPDUSessionResourceModifyResponseTransfer() (data ngapType.PDUSessionResourceModifyResponseTransfer) {

	// Qos Flow Add or Modify Response List
	data.QosFlowAddOrModifyResponseList = new(ngapType.QosFlowAddOrModifyResponseList)
	qosFlowAddOrModifyResponseList := data.QosFlowAddOrModifyResponseList

	qosFlowAddOrModifyResponseItem := ngapType.QosFlowAddOrModifyResponseItem{
		QosFlowIdentifier: ngapType.QosFlowIdentifier{
			Value: 5,
		},
	}

	qosFlowAddOrModifyResponseList.List = append(qosFlowAddOrModifyResponseList.List, qosFlowAddOrModifyResponseItem)

	return data
}

func buildPDUSessionResourceSetupUnsucessfulTransfer() (data ngapType.PDUSessionResourceSetupUnsuccessfulTransfer) {

	// Cause
	data.Cause.Present = ngapType.CausePresentRadioNetwork
	data.Cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
	data.Cause.RadioNetwork.Value = ngapType.CauseRadioNetworkPresentCellNotAvailable

	return data
}

func buildPDUSessionResourceModifyUnsuccessfulTransfer() (data ngapType.PDUSessionResourceModifyUnsuccessfulTransfer) {

	// Cause
	data.Cause = ngapType.Cause{
		Present: ngapType.CausePresentRadioNetwork,
		RadioNetwork: &ngapType.CauseRadioNetwork{
			Value: ngapType.CauseRadioNetworkPresentUnknownPDUSessionID,
		},
	}

	return data
}

func buildPDUSessionResourceReleaseResponseTransfer() (data ngapType.PDUSessionResourceReleaseResponseTransfer) {
	// PDU Session Resource Release Response Transfer

	return data
}

func buildPDUSessionResourceNotifyTransfer(
	qfis []int64, notiCause []uint64, relQfis []int64) (data ngapType.PDUSessionResourceNotifyTransfer) {

	if len(qfis) > 0 {
		data.QosFlowNotifyList = new(ngapType.QosFlowNotifyList)
	}
	if len(relQfis) > 0 {
		data.QosFlowReleasedList = new(ngapType.QosFlowListWithCause)
	}
	for i, qfi := range qfis {
		item := ngapType.QosFlowNotifyItem{
			QosFlowIdentifier: ngapType.QosFlowIdentifier{
				Value: qfi,
			},
			NotificationCause: ngapType.NotificationCause{
				Value: aper.Enumerated(notiCause[i]),
			},
		}
		data.QosFlowNotifyList.List = append(data.QosFlowNotifyList.List, item)
	}
	for _, qfi := range relQfis {
		item := ngapType.QosFlowWithCauseItem{
			QosFlowIdentifier: ngapType.QosFlowIdentifier{
				Value: qfi,
			},
			Cause: ngapType.Cause{
				Present: ngapType.CausePresentMisc,
				Misc: &ngapType.CauseMisc{
					Value: ngapType.CauseMiscPresentNotEnoughUserPlaneProcessingResources,
				},
			},
		}
		data.QosFlowReleasedList.List = append(data.QosFlowReleasedList.List, item)
	}
	return data
}

func buildPDUSessionResourceNotifyReleasedTransfer() (data ngapType.PDUSessionResourceNotifyReleasedTransfer) {
	// Cause
	data.Cause = ngapType.Cause{
		Present: ngapType.CausePresentRadioNetwork,
		RadioNetwork: &ngapType.CauseRadioNetwork{
			Value: ngapType.CauseRadioNetworkPresentUnknownPDUSessionID,
		},
	}
	return data
}

func buildPathSwitchRequestTransfer() (data ngapType.PathSwitchRequestTransfer) {

	// DL NG-U UP TNL information
	upTransportLayerInformation := &data.DLNGUUPTNLInformation
	upTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	upTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
	upTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString("\x00\x00\x00\x02")
	upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap("127.0.0.15", "")

	// Qos Flow Accepted List
	qosFlowAcceptedList := &data.QosFlowAcceptedList
	qosFlowAcceptedItem := ngapType.QosFlowAcceptedItem{
		QosFlowIdentifier: ngapType.QosFlowIdentifier{
			Value: 1,
		},
	}
	qosFlowAcceptedList.List = append(qosFlowAcceptedList.List, qosFlowAcceptedItem)

	return data
}

func buildPDUSessionResourceModifyIndicationTransfer() (
	data ngapType.PDUSessionResourceModifyIndicationTransfer) {

	// DL UP TNL Information
	data.DLQosFlowPerTNLInformation = ngapType.QosFlowPerTNLInformation{
		UPTransportLayerInformation: ngapType.UPTransportLayerInformation{
			Present: ngapType.UPTransportLayerInformationPresentGTPTunnel,
			GTPTunnel: &ngapType.GTPTunnel{
				TransportLayerAddress: ngapConvert.IPAddressToNgap("127.0.0.1", ""),
				GTPTEID: ngapType.GTPTEID{
					Value: aper.OctetString("\x00\x00\x00\x01"),
				},
			},
		},
	}

	return data
}

func buildPDUSessionResourceModifyConfirmTransfer(qfis []int64) (
	data ngapType.PDUSessionResourceModifyConfirmTransfer) {
	for _, qfi := range qfis {
		item := ngapType.QosFlowModifyConfirmItem{
			QosFlowIdentifier: ngapType.QosFlowIdentifier{
				Value: qfi,
			},
		}
		data.QosFlowModifyConfirmList.List = append(data.QosFlowModifyConfirmList.List, item)
	}
	return data
}

func buildPDUSessionResourceModifyIndicationUnsuccessfulTransfer() (
	data ngapType.PDUSessionResourceModifyIndicationUnsuccessfulTransfer) {
	data.Cause = ngapType.Cause{
		Present: ngapType.CausePresentTransport,
		Transport: &ngapType.CauseTransport{
			Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
		},
	}
	return data
}

func buildPDUSessionResourceReleaseCommandTransferr() (
	data ngapType.PDUSessionResourceReleaseCommandTransfer) {
	// Cause
	data.Cause = ngapType.Cause{
		Present: ngapType.CausePresentNas,
		Nas: &ngapType.CauseNas{
			Value: ngapType.CauseNasPresentNormalRelease,
		},
	}
	return data
}

func buildPathSwitchRequestSetupFailedTransfer() (data ngapType.PathSwitchRequestSetupFailedTransfer) {

	// Cause
	data.Cause = ngapType.Cause{
		Present: ngapType.CausePresentTransport,
		Transport: &ngapType.CauseTransport{
			Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
		},
	}

	return data
}

func buildHandoverRequestAcknowledgeTransfer() (data ngapType.HandoverRequestAcknowledgeTransfer) {

	// DL NG-U UP TNL information
	upTransportLayerInformation := &data.DLNGUUPTNLInformation
	upTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	upTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
	upTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString("\x00\x00\x00\x01")
	upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap("10.200.200.2", "")

	// Qos Flow Setup Response List
	qosFlowSetupResponseItem := ngapType.QosFlowItemWithDataForwarding{
		QosFlowIdentifier: ngapType.QosFlowIdentifier{
			Value: 9,
		},
	}

	data.DLForwardingUPTNLInformation = new(ngapType.UPTransportLayerInformation)
	dlForwardingUPTNLInfo := data.DLForwardingUPTNLInformation
	dlForwardingUPTNLInfo.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	dlForwardingUPTNLInfo.GTPTunnel = new(ngapType.GTPTunnel)
	dlForwardingUPTNLInfo.GTPTunnel.GTPTEID.Value = aper.OctetString("\x00\x00\x00\x02")
	dlForwardingUPTNLInfo.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap("10.200.200.2", "")

	data.QosFlowSetupResponseList.List = append(data.QosFlowSetupResponseList.List, qosFlowSetupResponseItem)

	return data
}

func buildHandoverResourceAllocationUnsuccessfulTransfer() (
	data ngapType.HandoverResourceAllocationUnsuccessfulTransfer) {

	data.Cause = ngapType.Cause{
		Present: ngapType.CausePresentRadioNetwork,
		RadioNetwork: &ngapType.CauseRadioNetwork{
			Value: ngapType.CauseRadioNetworkPresentHandoverCancelled,
		},
	}

	return data
}

func buildHandoverRequiredTransfer() (data ngapType.HandoverRequiredTransfer) {
	data.DirectForwardingPathAvailability = new(ngapType.DirectForwardingPathAvailability)
	data.DirectForwardingPathAvailability.Value = ngapType.DirectForwardingPathAvailabilityPresentDirectPathAvailable
	return data
}

func buildSourceToTargetTransparentTransfer(
	targetGNBID []byte, targetCellID []byte) (data ngapType.SourceNGRANNodeToTargetNGRANNodeTransparentContainer) {

	// RRC Container
	data.RRCContainer.Value = aper.OctetString("\x00\x00\x11")

	// PDU Session Resource Information List
	data.PDUSessionResourceInformationList = new(ngapType.PDUSessionResourceInformationList)
	infoItem := ngapType.PDUSessionResourceInformationItem{}
	infoItem.PDUSessionID.Value = 10
	qosItem := ngapType.QosFlowInformationItem{}
	qosItem.QosFlowIdentifier.Value = 1
	infoItem.QosFlowInformationList.List = append(infoItem.QosFlowInformationList.List, qosItem)
	data.PDUSessionResourceInformationList.List = append(data.PDUSessionResourceInformationList.List, infoItem)

	// Target Cell ID
	data.TargetCellID.Present = ngapType.TargetIDPresentTargetRANNodeID
	data.TargetCellID.NRCGI = new(ngapType.NRCGI)
	data.TargetCellID.NRCGI.PLMNIdentity = TestPlmn
	data.TargetCellID.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     append(targetGNBID, targetCellID...),
		BitLength: 36,
	}

	// UE History Information
	lastVisitedCellItem := ngapType.LastVisitedCellItem{}
	lastVisitedCellInfo := &lastVisitedCellItem.LastVisitedCellInformation
	lastVisitedCellInfo.Present = ngapType.LastVisitedCellInformationPresentNGRANCell
	lastVisitedCellInfo.NGRANCell = new(ngapType.LastVisitedNGRANCellInformation)
	ngRanCell := lastVisitedCellInfo.NGRANCell
	ngRanCell.GlobalCellID.Present = ngapType.NGRANCGIPresentNRCGI
	ngRanCell.GlobalCellID.NRCGI = new(ngapType.NRCGI)
	ngRanCell.GlobalCellID.NRCGI.PLMNIdentity = TestPlmn
	ngRanCell.GlobalCellID.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}
	ngRanCell.CellType.CellSize.Value = ngapType.CellSizePresentVerysmall
	ngRanCell.TimeUEStayedInCell.Value = 10

	data.UEHistoryInformation.List = append(data.UEHistoryInformation.List, lastVisitedCellItem)
	return data
}

func GetPDUSessionResourceSetupResponseTransfer(ipv4 string) []byte {
	data := buildPDUSessionResourceSetupResponseTransfer(ipv4)
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceSetupResponseTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceModifyResponseTransfer() []byte {
	data := buildPDUSessionResourceModifyResponseTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceModifyResponseTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceSetupUnsucessfulTransfer() []byte {
	data := buildPDUSessionResourceSetupUnsucessfulTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceSetupUnsucessfulTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceModifyUnsuccessfulTransfer() []byte {
	data := buildPDUSessionResourceModifyUnsuccessfulTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceModifyUnsuccessfulTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceModifyConfirmTransfer(qfis []int64) []byte {
	data := buildPDUSessionResourceModifyConfirmTransfer(qfis)
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceModifyConfirmTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceModifyIndicationUnsuccessfulTransfer() []byte {
	data := buildPDUSessionResourceModifyIndicationUnsuccessfulTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf(
			"aper MarshalWithParams error in GetPDUSessionResourceModifyIndicationUnsuccessfulTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceReleaseCommandTransfer() []byte {
	data := buildPDUSessionResourceReleaseCommandTransferr()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceReleaseCommandTransfer: %+v", err)
	}
	return encodeData
}

func GetPathSwitchRequestTransfer() []byte {
	data := buildPathSwitchRequestTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPathSwitchRequestTransfer: %+v", err)
	}
	return encodeData
}

func GetPathSwitchRequestSetupFailedTransfer() []byte {
	data := buildPathSwitchRequestSetupFailedTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPathSwitchRequestSetupFailedTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceModifyIndicationTransfer() []byte {
	data := buildPDUSessionResourceModifyIndicationTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceModifyIndicationTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceReleaseResponseTransfer() []byte {
	data := buildPDUSessionResourceReleaseResponseTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceReleaseResponseTransfer: %+v", err)
	}
	return encodeData
}

func GetPDUSessionResourceNotifyTransfer(qfis []int64, notiCause []uint64, relQfis []int64) []byte {
	data := buildPDUSessionResourceNotifyTransfer(qfis, notiCause, relQfis)
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceNotifyTransfer: %+v", err)
	}
	return encodeData
}
func GetPDUSessionResourceNotifyReleasedTransfer() []byte {
	data := buildPDUSessionResourceNotifyReleasedTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceNotifyReleasedTransfer: %+v", err)
	}
	return encodeData
}

func GetHandoverRequestAcknowledgeTransfer() []byte {
	data := buildHandoverRequestAcknowledgeTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetHandoverRequestAcknowledgeTransfer: %+v", err)
	}
	return encodeData
}

func GetHandoverResourceAllocationUnsuccessfulTransfer() []byte {
	data := buildHandoverResourceAllocationUnsuccessfulTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetHandoverResourceAllocationUnsuccessfulTransfer: %+v", err)
	}
	return encodeData
}

func GetHandoverRequiredTransfer() []byte {
	data := buildHandoverRequiredTransfer()
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetHandoverRequiredTransfer: %+v", err)
	}
	return encodeData
}

func GetSourceToTargetTransparentTransfer(targetGNBID []byte, targetCellID []byte) []byte {
	data := buildSourceToTargetTransparentTransfer(targetGNBID, targetCellID)
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		fatal.Fatalf("aper MarshalWithParams error in GetSourceToTargetTransparentTransfer: %+v", err)
	}
	return encodeData
}

func BuildInitialContextSetupResponseForRegistraionTest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentInitialContextSetupResponse
	successfulOutcome.Value.InitialContextSetupResponse = new(ngapType.InitialContextSetupResponse)

	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	initialContextSetupResponseIEs := &initialContextSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	return pdu
}

func BuildPDUSessionResourceSetupResponseForRegistrationTest(
	pduSessionId int64, amfUeNgapID, ranUeNgapID int64, ipv4 string) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceSetupResponse
	successfulOutcome.Value.PDUSessionResourceSetupResponse = new(ngapType.PDUSessionResourceSetupResponse)

	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	pDUSessionResourceSetupResponseIEs := &pDUSessionResourceSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// PDU Session Resource Setup Response List
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceSetupListSURes
	ie.Value.PDUSessionResourceSetupListSURes = new(ngapType.PDUSessionResourceSetupListSURes)

	pDUSessionResourceSetupListSURes := ie.Value.PDUSessionResourceSetupListSURes

	// PDU Session Resource Setup Response Item in PDU Session Resource Setup Response List
	pDUSessionResourceSetupItemSURes := ngapType.PDUSessionResourceSetupItemSURes{}
	pDUSessionResourceSetupItemSURes.PDUSessionID.Value = pduSessionId

	pDUSessionResourceSetupItemSURes.PDUSessionResourceSetupResponseTransfer =
		GetPDUSessionResourceSetupResponseTransfer(ipv4)

	pDUSessionResourceSetupListSURes.List = append(pDUSessionResourceSetupListSURes.List, pDUSessionResourceSetupItemSURes)

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// PDU Sessuin Resource Failed to Setup List
	// ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes
	// ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	// ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceFailedToSetupListSURes
	// ie.Value.PDUSessionResourceFailedToSetupListSURes = new(ngapType.PDUSessionResourceFailedToSetupListSURes)

	// pDUSessionResourceFailedToSetupListSURes := ie.Value.PDUSessionResourceFailedToSetupListSURes

	// // PDU Session Resource Failed to Setup Item in PDU Sessuin Resource Failed to Setup List
	// pDUSessionResourceFailedToSetupItemSURes := ngapType.PDUSessionResourceFailedToSetupItemSURes{}
	// pDUSessionResourceFailedToSetupItemSURes.PDUSessionID.Value = 10
	// pDUSessionResourceFailedToSetupItemSURes.PDUSessionResourceSetupUnsuccessfulTransfer =
	// 	GetPDUSessionResourceSetupUnsucessfulTransfer()

	// pDUSessionResourceFailedToSetupListSURes.List =
	// 	append(pDUSessionResourceFailedToSetupListSURes.List, pDUSessionResourceFailedToSetupItemSURes)

	// pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)
	// Criticality Diagnostics (optional)
	return pdu
}

func BuildPDUSessionResourceReleaseResponseForReleaseTest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceRelease
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceReleaseResponse
	successfulOutcome.Value.PDUSessionResourceReleaseResponse = new(ngapType.PDUSessionResourceReleaseResponse)

	pDUSessionResourceReleaseResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	pDUSessionResourceReleaseResponseIEs := &pDUSessionResourceReleaseResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)

	// PDU Session Resource Released List
	ie = ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceReleasedListRelRes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentPDUSessionResourceReleasedListRelRes
	ie.Value.PDUSessionResourceReleasedListRelRes = new(ngapType.PDUSessionResourceReleasedListRelRes)

	pDUSessionResourceReleasedListRelRes := ie.Value.PDUSessionResourceReleasedListRelRes

	// PDU Session Resource Released Item
	pDUSessionResourceReleasedItemRelRes := ngapType.PDUSessionResourceReleasedItemRelRes{}
	pDUSessionResourceReleasedItemRelRes.PDUSessionID.Value = 10

	pDUSessionResourceReleasedItemRelRes.PDUSessionResourceReleaseResponseTransfer =
		GetPDUSessionResourceReleaseResponseTransfer()
	// pDUSessionResourceReleasedItemRelRes.PDUSessionResourceReleaseResponseTransfer =aper.OctetString("\x01\x02\x03")

	pDUSessionResourceReleasedListRelRes.List =
		append(pDUSessionResourceReleasedListRelRes.List, pDUSessionResourceReleasedItemRelRes)

	pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)

	return pdu
}

func BuildNGSetupResponse(amfName string, guamiList []ngapType.ServedGUAMIItem,
	plmnList []ngapType.PLMNSupportItem, amfRelativeCapacity int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentNGSetupResponse
	successfulOutcome.Value.NGSetupResponse = new(ngapType.NGSetupResponse)

	nGSetupResponse := successfulOutcome.Value.NGSetupResponse
	nGSetupResponseIEs := &nGSetupResponse.ProtocolIEs

	// AMFName
	ie := ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFName
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentAMFName
	ie.Value.AMFName = new(ngapType.AMFName)

	aMFName := ie.Value.AMFName
	aMFName.Value = amfName

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// ServedGUAMIList
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDServedGUAMIList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentServedGUAMIList
	ie.Value.ServedGUAMIList = new(ngapType.ServedGUAMIList)

	servedGUAMIList := ie.Value.ServedGUAMIList
	servedGUAMIList.List = guamiList

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// relativeAMFCapacity
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRelativeAMFCapacity
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentRelativeAMFCapacity
	ie.Value.RelativeAMFCapacity = new(ngapType.RelativeAMFCapacity)
	relativeAMFCapacity := ie.Value.RelativeAMFCapacity
	relativeAMFCapacity.Value = amfRelativeCapacity

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// PLMNSupportList
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPLMNSupportList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentPLMNSupportList
	ie.Value.PLMNSupportList = new(ngapType.PLMNSupportList)

	pLMNSupportList := ie.Value.PLMNSupportList
	pLMNSupportList.List = plmnList

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	return pdu
}

func BuildPDUSessionResourceModifyConfirm(
	amfUeNgapId int64,
	ranUeNgapId int64,
	pduSessionResourceModifyConfirmList ngapType.PDUSessionResourceModifyListModCfm,
	pduSessionResourceFailedToModifyList ngapType.PDUSessionResourceFailedToModifyListModCfm,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceModifyIndication
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceModifyConfirm
	successfulOutcome.Value.PDUSessionResourceModifyConfirm = new(ngapType.PDUSessionResourceModifyConfirm)

	pDUSessionResourceModifyConfirm := successfulOutcome.Value.PDUSessionResourceModifyConfirm
	pDUSessionResourceModifyConfirmIEs := &pDUSessionResourceModifyConfirm.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceModifyConfirmIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyConfirmIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapId

	pDUSessionResourceModifyConfirmIEs.List = append(pDUSessionResourceModifyConfirmIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceModifyConfirmIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyConfirmIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapId

	pDUSessionResourceModifyConfirmIEs.List = append(pDUSessionResourceModifyConfirmIEs.List, ie)

	// PDU Session Resource Modify Confirm List
	ie = ngapType.PDUSessionResourceModifyConfirmIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceModifyListModCfm
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyConfirmIEsPresentPDUSessionResourceModifyListModCfm
	ie.Value.PDUSessionResourceModifyListModCfm = &pduSessionResourceModifyConfirmList
	pDUSessionResourceModifyConfirmIEs.List = append(pDUSessionResourceModifyConfirmIEs.List, ie)

	// PDU Session Resource Failed to Modify List
	if len(pduSessionResourceFailedToModifyList.List) > 0 {
		ie = ngapType.PDUSessionResourceModifyConfirmIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToModifyListModCfm
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceModifyConfirmIEsPresentPDUSessionResourceFailedToModifyListModCfm
		ie.Value.PDUSessionResourceFailedToModifyListModCfm = &pduSessionResourceFailedToModifyList
		pDUSessionResourceModifyConfirmIEs.List = append(pDUSessionResourceModifyConfirmIEs.List, ie)
	}

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.PDUSessionResourceModifyConfirmIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceModifyConfirmIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		pDUSessionResourceModifyConfirmIEs.List = append(pDUSessionResourceModifyConfirmIEs.List, ie)
	}

	return pdu
}

func BuildPDUSessionResourceReleaseCommand(
	amfUeNgapId int64,
	ranUeNgapId int64,
	pagingPriority *ngapType.RANPagingPriority,
	nasPdu []byte,
	pduSessionResourceReleasedList ngapType.PDUSessionResourceToReleaseListRelCmd) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceRelease
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject
	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPDUSessionResourceReleaseCommand
	initiatingMessage.Value.PDUSessionResourceReleaseCommand = new(ngapType.PDUSessionResourceReleaseCommand)

	pDUSessionResourceReleaseCommand := initiatingMessage.Value.PDUSessionResourceReleaseCommand
	PDUSessionResourceReleaseCommandIEs := &pDUSessionResourceReleaseCommand.ProtocolIEs

	// AMFUENGAPID
	ie := ngapType.PDUSessionResourceReleaseCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceReleaseCommandIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapId

	PDUSessionResourceReleaseCommandIEs.List = append(PDUSessionResourceReleaseCommandIEs.List, ie)

	// RANUENGAPID
	ie = ngapType.PDUSessionResourceReleaseCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceReleaseCommandIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapId

	PDUSessionResourceReleaseCommandIEs.List = append(PDUSessionResourceReleaseCommandIEs.List, ie)

	// RAN Paging Priority (optional)

	if pagingPriority != nil {
		ie = ngapType.PDUSessionResourceReleaseCommandIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPagingPriority
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceReleaseCommandIEsPresentRANPagingPriority
		ie.Value.RANPagingPriority = pagingPriority

		PDUSessionResourceReleaseCommandIEs.List = append(PDUSessionResourceReleaseCommandIEs.List, ie)
	}

	// NAS-PDU (optional)
	if nasPdu != nil {
		ie = ngapType.PDUSessionResourceReleaseCommandIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNASPDU
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceReleaseCommandIEsPresentNASPDU
		ie.Value.NASPDU = new(ngapType.NASPDU)

		ie.Value.NASPDU.Value = nasPdu

		PDUSessionResourceReleaseCommandIEs.List = append(PDUSessionResourceReleaseCommandIEs.List, ie)
	}

	// PDUSessionResourceToReleaseListRelCmd
	ie = ngapType.PDUSessionResourceReleaseCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceToReleaseListRelCmd
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceReleaseCommandIEsPresentPDUSessionResourceToReleaseListRelCmd
	ie.Value.PDUSessionResourceToReleaseListRelCmd = &pduSessionResourceReleasedList
	PDUSessionResourceReleaseCommandIEs.List = append(PDUSessionResourceReleaseCommandIEs.List, ie)

	return pdu
}

func BuildOverloadStart(
	action *ngapType.OverloadAction,
	ind *int64,
	list []ngapType.OverloadStartNSSAIItem) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeOverloadStart
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentOverloadStart
	initiatingMessage.Value.OverloadStart = new(ngapType.OverloadStart)

	overloadStart := initiatingMessage.Value.OverloadStart
	overloadStartIEs := &overloadStart.ProtocolIEs
	// AMFOverloadResponse
	if action != nil {
		ie := ngapType.OverloadStartIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFOverloadResponse
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.OverloadStartIEsPresentAMFOverloadResponse
		ie.Value.AMFOverloadResponse = new(ngapType.OverloadResponse)

		aMFOverloadResponse := ie.Value.AMFOverloadResponse
		aMFOverloadResponse.Present = ngapType.OverloadResponsePresentOverloadAction
		aMFOverloadResponse.OverloadAction = action

		overloadStartIEs.List = append(overloadStartIEs.List, ie)
	}
	// AMFTrafficLoadReductionIndication
	if ind != nil {
		ie := ngapType.OverloadStartIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFTrafficLoadReductionIndication
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.OverloadStartIEsPresentAMFTrafficLoadReductionIndication
		ie.Value.AMFTrafficLoadReductionIndication = new(ngapType.TrafficLoadReductionIndication)

		aMFTrafficLoadReductionIndication := ie.Value.AMFTrafficLoadReductionIndication
		aMFTrafficLoadReductionIndication.Value = *ind

		overloadStartIEs.List = append(overloadStartIEs.List, ie)
	}
	// OverloadStartNSSAIList
	if len(list) > 0 {
		ie := ngapType.OverloadStartIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDOverloadStartNSSAIList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.OverloadStartIEsPresentOverloadStartNSSAIList
		ie.Value.OverloadStartNSSAIList = new(ngapType.OverloadStartNSSAIList)

		overloadStartNSSAIList := ie.Value.OverloadStartNSSAIList
		overloadStartNSSAIList.List = list

		overloadStartIEs.List = append(overloadStartIEs.List, ie)
	}
	return pdu
}

func BuildOverloadStop() (pdu ngapType.NGAPPDU) {
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeOverloadStop
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentOverloadStop
	initiatingMessage.Value.OverloadStop = new(ngapType.OverloadStop)

	return pdu
}
