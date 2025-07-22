package message

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/util"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/aper"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

func BuildPDUSessionResourceReleaseCommand(ue *context.RanUe, nasPdu []byte,
	pduSessionResourceReleasedList ngapType.PDUSessionResourceToReleaseListRelCmd,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU
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
	aMFUENGAPID.Value = ue.AmfUeNgapId

	PDUSessionResourceReleaseCommandIEs.List = append(PDUSessionResourceReleaseCommandIEs.List, ie)

	// RANUENGAPID
	ie = ngapType.PDUSessionResourceReleaseCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceReleaseCommandIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	PDUSessionResourceReleaseCommandIEs.List = append(PDUSessionResourceReleaseCommandIEs.List, ie)

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

	return ngap.Encoder(pdu)
}

func BuildNGSetupResponse() ([]byte, error) {
	amfSelf := context.GetSelf()
	var pdu ngapType.NGAPPDU
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
	aMFName.Value = amfSelf.Name

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// ServedGUAMIList
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDServedGUAMIList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentServedGUAMIList
	ie.Value.ServedGUAMIList = new(ngapType.ServedGUAMIList)

	servedGUAMIList := ie.Value.ServedGUAMIList
	for _, guami := range amfSelf.ServedGuamiList {
		servedGUAMIItem := ngapType.ServedGUAMIItem{}
		servedGUAMIItem.GUAMI.PLMNIdentity = ngapConvert.PlmnIdToNgap(util.PlmnIdNidToModelsPlmnId(*guami.PlmnId))
		regionId, setId, prtId := ngapConvert.AmfIdToNgap(guami.AmfId)
		servedGUAMIItem.GUAMI.AMFRegionID.Value = regionId
		servedGUAMIItem.GUAMI.AMFSetID.Value = setId
		servedGUAMIItem.GUAMI.AMFPointer.Value = prtId
		servedGUAMIList.List = append(servedGUAMIList.List, servedGUAMIItem)
	}

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// relativeAMFCapacity
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRelativeAMFCapacity
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentRelativeAMFCapacity
	ie.Value.RelativeAMFCapacity = new(ngapType.RelativeAMFCapacity)
	relativeAMFCapacity := ie.Value.RelativeAMFCapacity
	relativeAMFCapacity.Value = amfSelf.RelativeCapacity

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// ServedGUAMIList
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPLMNSupportList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentPLMNSupportList
	ie.Value.PLMNSupportList = new(ngapType.PLMNSupportList)

	pLMNSupportList := ie.Value.PLMNSupportList
	for _, plmnItem := range amfSelf.PlmnSupportList {
		pLMNSupportItem := ngapType.PLMNSupportItem{}
		pLMNSupportItem.PLMNIdentity = ngapConvert.PlmnIdToNgap(*plmnItem.PlmnId)
		for _, snssai := range plmnItem.SNssaiList {
			sliceSupportItem := ngapType.SliceSupportItem{}
			sliceSupportItem.SNSSAI = ngapConvert.SNssaiToNgap(snssai)
			pLMNSupportItem.SliceSupportList.List = append(pLMNSupportItem.SliceSupportList.List, sliceSupportItem)
		}
		pLMNSupportList.List = append(pLMNSupportList.List, pLMNSupportItem)
	}

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildNGSetupFailure(cause ngapType.Cause) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentNGSetupFailure
	unsuccessfulOutcome.Value.NGSetupFailure = new(ngapType.NGSetupFailure)

	nGSetupFailure := unsuccessfulOutcome.Value.NGSetupFailure
	nGSetupFailureIEs := &nGSetupFailure.ProtocolIEs

	// Cause
	ie := ngapType.NGSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupFailureIEsPresentCause
	ie.Value.Cause = &cause

	nGSetupFailureIEs.List = append(nGSetupFailureIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildNGReset(
	cause ngapType.Cause, partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	logger.NgapLog.Trace("Build NG Reset message")

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
	ie.Value.Cause = &cause

	nGResetIEs.List = append(nGResetIEs.List, ie)

	// Reset Type
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

	return ngap.Encoder(pdu)
}

func BuildNGResetAcknowledge(partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU

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
	if partOfNGInterface != nil && len(partOfNGInterface.List) > 0 {
		ie := ngapType.NGResetAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.NGResetAcknowledgeIEsPresentUEAssociatedLogicalNGConnectionList
		ie.Value.UEAssociatedLogicalNGConnectionList = new(ngapType.UEAssociatedLogicalNGConnectionList)

		uEAssociatedLogicalNGConnectionList := ie.Value.UEAssociatedLogicalNGConnectionList

		for i, item := range partOfNGInterface.List {
			if item.AMFUENGAPID == nil && item.RANUENGAPID == nil {
				logger.NgapLog.Warn("[Build NG Reset Ack] No AmfUeNgapID & RanUeNgapID")
				continue
			}

			uEAssociatedLogicalNGConnectionItem := ngapType.UEAssociatedLogicalNGConnectionItem{}

			if item.AMFUENGAPID != nil {
				uEAssociatedLogicalNGConnectionItem.AMFUENGAPID = new(ngapType.AMFUENGAPID)
				uEAssociatedLogicalNGConnectionItem.AMFUENGAPID = item.AMFUENGAPID
				logger.NgapLog.Tracef(
					"[Build NG Reset Ack] (pair %d) AmfUeNgapID[%d]", i, uEAssociatedLogicalNGConnectionItem.AMFUENGAPID)
			}
			if item.RANUENGAPID != nil {
				uEAssociatedLogicalNGConnectionItem.RANUENGAPID = new(ngapType.RANUENGAPID)
				uEAssociatedLogicalNGConnectionItem.RANUENGAPID = item.RANUENGAPID
				logger.NgapLog.Tracef(
					"[Build NG Reset Ack] (pair %d) RanUeNgapID[%d]", i, uEAssociatedLogicalNGConnectionItem.RANUENGAPID)
			}

			uEAssociatedLogicalNGConnectionList.List = append(uEAssociatedLogicalNGConnectionList.List,
				uEAssociatedLogicalNGConnectionItem)
		}

		nGResetAcknowledgeIEs.List = append(nGResetAcknowledgeIEs.List, ie)
	}

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie := ngapType.NGResetAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.NGResetAcknowledgeIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		ie.Value.CriticalityDiagnostics = criticalityDiagnostics

		nGResetAcknowledgeIEs.List = append(nGResetAcknowledgeIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildDownlinkNasTransport(ue *context.RanUe, nasPdu []byte,
	mobilityRestrictionList *ngapType.MobilityRestrictionList,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeDownlinkNASTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentDownlinkNASTransport
	initiatingMessage.Value.DownlinkNASTransport = new(ngapType.DownlinkNASTransport)

	downlinkNasTransport := initiatingMessage.Value.DownlinkNASTransport
	downlinkNasTransportIEs := &downlinkNasTransport.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.DownlinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkNASTransportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	downlinkNasTransportIEs.List = append(downlinkNasTransportIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.DownlinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkNASTransportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	downlinkNasTransportIEs.List = append(downlinkNasTransportIEs.List, ie)

	// NAS PDU
	ie = ngapType.DownlinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNASPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkNASTransportIEsPresentNASPDU
	ie.Value.NASPDU = new(ngapType.NASPDU)

	ie.Value.NASPDU.Value = nasPdu

	downlinkNasTransportIEs.List = append(downlinkNasTransportIEs.List, ie)

	// Old AMF (optional)
	if ue.OldAmfName != "" {
		ie = ngapType.DownlinkNASTransportIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDOldAMF
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.DownlinkNASTransportIEsPresentOldAMF
		ie.Value.OldAMF = new(ngapType.AMFName)

		ie.Value.OldAMF.Value = ue.OldAmfName

		downlinkNasTransportIEs.List = append(downlinkNasTransportIEs.List, ie)
		ue.OldAmfName = "" // clear data
	}

	// RAN Paging Priority (optional)
	// Mobility Restriction List (optional)
	if c := factory.AmfConfig.GetNgapIEMobilityRestrictionList(); c != nil && c.Enable &&
		ue.Ran.AnType == models.AccessType__3_GPP_ACCESS && mobilityRestrictionList != nil {
		amfUe := ue.AmfUe
		if amfUe == nil {
			return nil, fmt.Errorf("amfUe is nil")
		}

		ie = ngapType.DownlinkNASTransportIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDMobilityRestrictionList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.DownlinkNASTransportIEsPresentMobilityRestrictionList
		ie.Value.MobilityRestrictionList = mobilityRestrictionList
		downlinkNasTransportIEs.List = append(downlinkNasTransportIEs.List, ie)
	}
	// Index to RAT/Frequency Selection Priority (optional)
	// UE Aggregate Maximum Bit Rate (optional)
	// Allowed NSSAI (optional)

	return ngap.Encoder(pdu)
}

func BuildUEContextReleaseCommand(
	ue *context.RanUe, causePresent int, cause aper.Enumerated,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUEContextRelease
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUEContextReleaseCommand
	initiatingMessage.Value.UEContextReleaseCommand = new(ngapType.UEContextReleaseCommand)

	ueContextReleaseCommand := initiatingMessage.Value.UEContextReleaseCommand
	ueContextReleaseCommandIEs := &ueContextReleaseCommand.ProtocolIEs

	// UE NGAP IDs
	ie := ngapType.UEContextReleaseCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUENGAPIDs
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseCommandIEsPresentUENGAPIDs
	ie.Value.UENGAPIDs = new(ngapType.UENGAPIDs)

	ueNGAPIDs := ie.Value.UENGAPIDs

	if ue.RanUeNgapId == context.RanUeNgapIdUnspecified {
		ueNGAPIDs.Present = ngapType.UENGAPIDsPresentAMFUENGAPID
		ueNGAPIDs.AMFUENGAPID = new(ngapType.AMFUENGAPID)

		ueNGAPIDs.AMFUENGAPID.Value = ue.AmfUeNgapId
	} else {
		ueNGAPIDs.Present = ngapType.UENGAPIDsPresentUENGAPIDPair
		ueNGAPIDs.UENGAPIDPair = new(ngapType.UENGAPIDPair)

		ueNGAPIDs.UENGAPIDPair.AMFUENGAPID.Value = ue.AmfUeNgapId
		ueNGAPIDs.UENGAPIDPair.RANUENGAPID.Value = ue.RanUeNgapId
	}

	ueContextReleaseCommandIEs.List = append(ueContextReleaseCommandIEs.List, ie)

	// Cause
	ie = ngapType.UEContextReleaseCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCommandIEsPresentCause
	ngapCause := ngapType.Cause{
		Present: causePresent,
	}
	switch causePresent {
	case ngapType.CausePresentNothing:
		return nil, fmt.Errorf("cause present is nothing")
	case ngapType.CausePresentRadioNetwork:
		ngapCause.RadioNetwork = new(ngapType.CauseRadioNetwork)
		ngapCause.RadioNetwork.Value = cause
	case ngapType.CausePresentTransport:
		ngapCause.Transport = new(ngapType.CauseTransport)
		ngapCause.Transport.Value = cause
	case ngapType.CausePresentNas:
		ngapCause.Nas = new(ngapType.CauseNas)
		ngapCause.Nas.Value = cause
	case ngapType.CausePresentProtocol:
		ngapCause.Protocol = new(ngapType.CauseProtocol)
		ngapCause.Protocol.Value = cause
	case ngapType.CausePresentMisc:
		ngapCause.Misc = new(ngapType.CauseMisc)
		ngapCause.Misc.Value = cause
	default:
		return nil, fmt.Errorf("cause present is unknown")
	}
	ie.Value.Cause = &ngapCause

	ueContextReleaseCommandIEs.List = append(ueContextReleaseCommandIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildErrorIndication(amfUeNgapId, ranUeNgapId *int64, cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeErrorIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentErrorIndication
	initiatingMessage.Value.ErrorIndication = new(ngapType.ErrorIndication)

	errorIndication := initiatingMessage.Value.ErrorIndication
	errorIndicationIEs := &errorIndication.ProtocolIEs

	if cause == nil && criticalityDiagnostics == nil {
		logger.NgapLog.Error(
			"[Build Error Indication] shall contain at least either the Cause or the Criticality Diagnostics")
	}

	if amfUeNgapId != nil {
		ie := ngapType.ErrorIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.ErrorIndicationIEsPresentAMFUENGAPID
		ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

		aMFUENGAPID := ie.Value.AMFUENGAPID
		aMFUENGAPID.Value = *amfUeNgapId

		errorIndicationIEs.List = append(errorIndicationIEs.List, ie)
	}

	if ranUeNgapId != nil {
		ie := ngapType.ErrorIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.ErrorIndicationIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = *ranUeNgapId

		errorIndicationIEs.List = append(errorIndicationIEs.List, ie)
	}

	if cause != nil {
		ie := ngapType.ErrorIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCause
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.ErrorIndicationIEsPresentCause
		ie.Value.Cause = new(ngapType.Cause)

		ie.Value.Cause = cause

		errorIndicationIEs.List = append(errorIndicationIEs.List, ie)
	}

	if criticalityDiagnostics != nil {
		ie := ngapType.ErrorIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.ErrorIndicationIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		ie.Value.CriticalityDiagnostics = criticalityDiagnostics

		errorIndicationIEs.List = append(errorIndicationIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildUERadioCapabilityCheckRequest(ue *context.RanUe) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUERadioCapabilityCheck
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUERadioCapabilityCheckRequest
	initiatingMessage.Value.UERadioCapabilityCheckRequest = new(ngapType.UERadioCapabilityCheckRequest)

	uERadioCapabilityCheckRequest := initiatingMessage.Value.UERadioCapabilityCheckRequest
	uERadioCapabilityCheckRequestIEs := &uERadioCapabilityCheckRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UERadioCapabilityCheckRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UERadioCapabilityCheckRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	uERadioCapabilityCheckRequestIEs.List = append(uERadioCapabilityCheckRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UERadioCapabilityCheckRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UERadioCapabilityCheckRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uERadioCapabilityCheckRequestIEs.List = append(uERadioCapabilityCheckRequestIEs.List, ie)

	// TODO:UE Radio Capability(optional)
	return ngap.Encoder(pdu)
}

func BuildHandoverCancelAcknowledge(
	ue *context.RanUe, criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeHandoverCancel
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentHandoverCancelAcknowledge
	successfulOutcome.Value.HandoverCancelAcknowledge = new(ngapType.HandoverCancelAcknowledge)

	handoverCancelAcknowledge := successfulOutcome.Value.HandoverCancelAcknowledge
	handoverCancelAcknowledgeIEs := &handoverCancelAcknowledge.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.HandoverCancelAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	handoverCancelAcknowledgeIEs.List = append(handoverCancelAcknowledgeIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.HandoverCancelAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	handoverCancelAcknowledgeIEs.List = append(handoverCancelAcknowledgeIEs.List, ie)

	// Criticality Diagnostics [optional]
	if criticalityDiagnostics != nil {
		handoverCancelAcknowledgeIEsie := ngapType.HandoverCancelAcknowledgeIEs{}
		handoverCancelAcknowledgeIEsie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		handoverCancelAcknowledgeIEsie.Criticality.Value = ngapType.CriticalityPresentIgnore
		handoverCancelAcknowledgeIEsie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentCriticalityDiagnostics
		handoverCancelAcknowledgeIEsie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		handoverCancelAcknowledgeIEsie.Value.CriticalityDiagnostics = criticalityDiagnostics

		handoverCancelAcknowledgeIEs.List = append(handoverCancelAcknowledgeIEs.List, handoverCancelAcknowledgeIEsie)
	}

	return ngap.Encoder(pdu)
}

// nasPDU: from nas layer
// pduSessionResourceSetupRequestList: provided by AMF, and transfer data is from SMF
func BuildPDUSessionResourceSetupRequest(ue *context.RanUe, nasPdu []byte,
	pduSessionResourceSetupRequestList *ngapType.PDUSessionResourceSetupListSUReq,
) ([]byte, error) {
	// TODO: Ran Paging Priority (optional)

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPDUSessionResourceSetupRequest
	initiatingMessage.Value.PDUSessionResourceSetupRequest = new(ngapType.PDUSessionResourceSetupRequest)

	pDUSessionResourceSetupRequest := initiatingMessage.Value.PDUSessionResourceSetupRequest
	pDUSessionResourceSetupRequestIEs := &pDUSessionResourceSetupRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceSetupRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	pDUSessionResourceSetupRequestIEs.List = append(pDUSessionResourceSetupRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceSetupRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	pDUSessionResourceSetupRequestIEs.List = append(pDUSessionResourceSetupRequestIEs.List, ie)

	// Ran Paging Priority (optional)

	// NAS-PDU (optional)
	if nasPdu != nil {
		ie = ngapType.PDUSessionResourceSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNASPDU
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PDUSessionResourceSetupRequestIEsPresentNASPDU
		ie.Value.NASPDU = new(ngapType.NASPDU)

		ie.Value.NASPDU.Value = nasPdu

		pDUSessionResourceSetupRequestIEs.List = append(pDUSessionResourceSetupRequestIEs.List, ie)
	}

	// PDU Session Resource Setup Request list
	ie = ngapType.PDUSessionResourceSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListSUReq
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceSetupRequestIEsPresentPDUSessionResourceSetupListSUReq
	ie.Value.PDUSessionResourceSetupListSUReq = pduSessionResourceSetupRequestList
	pDUSessionResourceSetupRequestIEs.List = append(pDUSessionResourceSetupRequestIEs.List, ie)

	// UE AggreateMaximum Bit Rate
	ie = ngapType.PDUSessionResourceSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUEAggregateMaximumBitRate
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupRequestIEsPresentUEAggregateMaximumBitRate
	ie.Value.UEAggregateMaximumBitRate = new(ngapType.UEAggregateMaximumBitRate)
	ueAmbrUL := ngapConvert.UEAmbrToInt64(ue.AmfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Uplink)
	ueAmbrDL := ngapConvert.UEAmbrToInt64(ue.AmfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Downlink)
	ie.Value.UEAggregateMaximumBitRate.UEAggregateMaximumBitRateUL.Value = ueAmbrUL
	ie.Value.UEAggregateMaximumBitRate.UEAggregateMaximumBitRateDL.Value = ueAmbrDL
	pDUSessionResourceSetupRequestIEs.List = append(pDUSessionResourceSetupRequestIEs.List, ie)

	return ngap.Encoder(pdu)
}

// pduSessionResourceModifyConfirmList: provided by AMF, and transfer data is return from SMF
// pduSessionResourceFailedToModifyList: provided by AMF, and transfer data is return from SMF
func BuildPDUSessionResourceModifyConfirm(
	ue *context.RanUe,
	pduSessionResourceModifyConfirmList ngapType.PDUSessionResourceModifyListModCfm,
	pduSessionResourceFailedToModifyList ngapType.PDUSessionResourceFailedToModifyListModCfm,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU
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
	aMFUENGAPID.Value = ue.AmfUeNgapId

	pDUSessionResourceModifyConfirmIEs.List = append(pDUSessionResourceModifyConfirmIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceModifyConfirmIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyConfirmIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

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

	return ngap.Encoder(pdu)
}

// pduSessionResourceModifyRequestList: from SMF
func BuildPDUSessionResourceModifyRequest(ue *context.RanUe,
	pduSessionResourceModifyRequestList ngapType.PDUSessionResourceModifyListModReq,
) ([]byte, error) {
	// TODO: Ran Paging Priority (optional)

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceModify
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPDUSessionResourceModifyRequest
	initiatingMessage.Value.PDUSessionResourceModifyRequest = new(ngapType.PDUSessionResourceModifyRequest)

	pDUSessionResourceModifyRequest := initiatingMessage.Value.PDUSessionResourceModifyRequest
	pDUSessionResourceModifyRequestIEs := &pDUSessionResourceModifyRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceModifyRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceModifyRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	pDUSessionResourceModifyRequestIEs.List = append(pDUSessionResourceModifyRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceModifyRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceModifyRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	pDUSessionResourceModifyRequestIEs.List = append(pDUSessionResourceModifyRequestIEs.List, ie)

	// Ran Paging Priority (optional)

	// PDU Session Resource Modify Request List
	ie = ngapType.PDUSessionResourceModifyRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceModifyListModReq
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceModifyRequestIEsPresentPDUSessionResourceModifyListModReq
	ie.Value.PDUSessionResourceModifyListModReq = &pduSessionResourceModifyRequestList
	pDUSessionResourceModifyRequestIEs.List = append(pDUSessionResourceModifyRequestIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildInitialContextSetupRequest(
	amfUe *context.AmfUe,
	anType models.AccessType,
	nasPdu []byte,
	pduSessionResourceSetupRequestList *ngapType.PDUSessionResourceSetupListCxtReq,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	coreNetworkAssistanceInfo *ngapType.CoreNetworkAssistanceInformation,
	emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator,
) ([]byte, error) {
	// Old AMF: new amf should get old amf's amf name

	// rrcInactiveTransitionReportRequest: configured by amf
	// This IE is used to request the NG-RAN node to report or stop reporting to the 5GC
	// when the UE enters or leaves RRC_INACTIVE state. (TS 38.413 9.3.1.91)

	// accessType indicate amfUe send this msg for which accessType
	// emergencyFallbackIndicator: configured by amf (TS 23.501 5.16.4.11)
	// coreNetworkAssistanceInfo TS 23.501 5.4.6, 5.4.6.2

	// Mobility Restriction List TS 23.501 5.3.4
	// TS 23.501 5.3.4.1.1: For a given UE, the core network determines the Mobility restrictions
	// based on UE subscription information.
	// TS 38.413 9.3.1.85: This IE defines roaming or access restrictions for subsequent mobility action for
	// which the NR-RAN provides information about the target of the mobility action towards
	// the UE, e.g., handover, or for SCG selection during dual connectivity operation or for
	// assigning proper RNAs. If the NG-RAN receives the Mobility Restriction List IE, it shall
	// overwrite previously received mobility restriction information.

	if amfUe == nil {
		return nil, fmt.Errorf("amfUe is nil")
	}

	var pdu ngapType.NGAPPDU
	ranUe, ok := amfUe.RanUe[anType]
	if !ok {
		return nil, fmt.Errorf("ranUe for %s is nil", anType)
	}
	amfSelf := context.GetSelf()

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentInitialContextSetupRequest
	initiatingMessage.Value.InitialContextSetupRequest = new(ngapType.InitialContextSetupRequest)

	initialContextSetupRequest := initiatingMessage.Value.InitialContextSetupRequest
	initialContextSetupRequestIEs := &initialContextSetupRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ranUe.AmfUeNgapId

	initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUe.RanUeNgapId

	initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)

	// Old AMF (optional)
	if ranUe.OldAmfName != "" {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDOldAMF
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentOldAMF
		ie.Value.OldAMF = new(ngapType.AMFName)
		ie.Value.OldAMF.Value = ranUe.OldAmfName
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
		ranUe.OldAmfName = "" // clear data
	}

	// UE Aggregate Maximum Bit Rate (conditional: if pdu session resource setup)
	// The subscribed UE-AMBR is a subscription parameter which is
	// retrieved from UDM and provided to the (R)AN by the AMF
	if pduSessionResourceSetupRequestList != nil && len(pduSessionResourceSetupRequestList.List) > 0 {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUEAggregateMaximumBitRate
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentUEAggregateMaximumBitRate
		ie.Value.UEAggregateMaximumBitRate = new(ngapType.UEAggregateMaximumBitRate)

		ueAmbrUL := ngapConvert.UEAmbrToInt64(amfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Uplink)
		ueAmbrDL := ngapConvert.UEAmbrToInt64(amfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Downlink)
		ie.Value.UEAggregateMaximumBitRate.UEAggregateMaximumBitRateUL.Value = ueAmbrUL
		ie.Value.UEAggregateMaximumBitRate.UEAggregateMaximumBitRateDL.Value = ueAmbrDL

		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// Core Network Assistance Information (optional)
	if coreNetworkAssistanceInfo != nil {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCoreNetworkAssistanceInformation
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentCoreNetworkAssistanceInformation
		ie.Value.CoreNetworkAssistanceInformation = coreNetworkAssistanceInfo
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// GUAMI
	ie = ngapType.InitialContextSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDGUAMI
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentGUAMI
	ie.Value.GUAMI = new(ngapType.GUAMI)

	guami := ie.Value.GUAMI
	plmnID := &guami.PLMNIdentity
	amfRegionID := &guami.AMFRegionID
	amfSetID := &guami.AMFSetID
	amfPtrID := &guami.AMFPointer

	servedGuami := amfSelf.ServedGuamiList[0]

	*plmnID = ngapConvert.PlmnIdToNgap(util.PlmnIdNidToModelsPlmnId(*servedGuami.PlmnId))
	amfRegionID.Value, amfSetID.Value, amfPtrID.Value = ngapConvert.AmfIdToNgap(servedGuami.AmfId)

	initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)

	// PDU Session Resource Setup Request List
	if pduSessionResourceSetupRequestList != nil && len(pduSessionResourceSetupRequestList.List) > 0 {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtReq
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentPDUSessionResourceSetupListCxtReq
		ie.Value.PDUSessionResourceSetupListCxtReq = pduSessionResourceSetupRequestList
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// Allowed NSSAI
	ie = ngapType.InitialContextSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAllowedNSSAI
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentAllowedNSSAI
	ie.Value.AllowedNSSAI = new(ngapType.AllowedNSSAI)

	allowedNSSAI := ie.Value.AllowedNSSAI

	for _, allowedSnssai := range amfUe.AllowedNssai[anType] {
		allowedNSSAIItem := ngapType.AllowedNSSAIItem{}
		ngapSnssai := ngapConvert.SNssaiToNgap(*allowedSnssai.AllowedSnssai)
		allowedNSSAIItem.SNSSAI = ngapSnssai
		allowedNSSAI.List = append(allowedNSSAI.List, allowedNSSAIItem)
	}

	initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)

	// UE Security Capabilities
	ie = ngapType.InitialContextSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUESecurityCapabilities
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentUESecurityCapabilities
	ie.Value.UESecurityCapabilities = new(ngapType.UESecurityCapabilities)

	ueSecurityCapabilities := ie.Value.UESecurityCapabilities
	nrEncryptionAlgorighm := []byte{0x00, 0x00}

	nrEncryptionAlgorighm[0] |= amfUe.UESecurityCapability.GetEA1_128_5G() << 7
	nrEncryptionAlgorighm[0] |= amfUe.UESecurityCapability.GetEA2_128_5G() << 6
	nrEncryptionAlgorighm[0] |= amfUe.UESecurityCapability.GetEA3_128_5G() << 5
	ueSecurityCapabilities.NRencryptionAlgorithms.Value = ngapConvert.ByteToBitString(nrEncryptionAlgorighm, 16)

	nrIntegrityAlgorithm := []byte{0x00, 0x00}

	nrIntegrityAlgorithm[0] |= amfUe.UESecurityCapability.GetIA1_128_5G() << 7
	nrIntegrityAlgorithm[0] |= amfUe.UESecurityCapability.GetIA2_128_5G() << 6
	nrIntegrityAlgorithm[0] |= amfUe.UESecurityCapability.GetIA3_128_5G() << 5

	ueSecurityCapabilities.NRintegrityProtectionAlgorithms.Value = ngapConvert.ByteToBitString(nrIntegrityAlgorithm, 16)

	// only support NR algorithms
	eutraEncryptionAlgorithm := []byte{0x00, 0x00}
	ueSecurityCapabilities.EUTRAencryptionAlgorithms.Value = ngapConvert.ByteToBitString(eutraEncryptionAlgorithm, 16)

	eutraIntegrityAlgorithm := []byte{0x00, 0x00}
	ueSecurityCapabilities.EUTRAintegrityProtectionAlgorithms.Value = ngapConvert.
		ByteToBitString(eutraIntegrityAlgorithm, 16)

	initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)

	// Security Key
	ie = ngapType.InitialContextSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSecurityKey
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentSecurityKey
	ie.Value.SecurityKey = new(ngapType.SecurityKey)

	securityKey := ie.Value.SecurityKey
	switch ranUe.Ran.AnType {
	case models.AccessType__3_GPP_ACCESS:
		securityKey.Value = ngapConvert.ByteToBitString(amfUe.Kgnb, 256)
	case models.AccessType_NON_3_GPP_ACCESS:
		securityKey.Value = ngapConvert.ByteToBitString(amfUe.Kn3iwf, 256)
	}

	initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)

	// Trace Activation (optional)
	if amfUe.TraceData != nil {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDTraceActivation
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentTraceActivation
		ie.Value.TraceActivation = new(ngapType.TraceActivation)
		// TS 32.422 4.2.2.9
		// TODO: AMF allocate Trace Recording Session Reference
		traceActivation := ngapConvert.TraceDataToNgap(*amfUe.TraceData, ranUe.Trsr)
		ie.Value.TraceActivation = &traceActivation
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// Mobility Restriction List (optional)
	if c := factory.AmfConfig.GetNgapIEMobilityRestrictionList(); c != nil && c.Enable &&
		anType == models.AccessType__3_GPP_ACCESS {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDMobilityRestrictionList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentMobilityRestrictionList
		ie.Value.MobilityRestrictionList = new(ngapType.MobilityRestrictionList)

		mobilityRestrictionList := BuildIEMobilityRestrictionList(amfUe)
		ie.Value.MobilityRestrictionList = &mobilityRestrictionList
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// UE Radio Capability (optional)
	if amfUe.UeRadioCapability != "" {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUERadioCapability
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentUERadioCapability
		ie.Value.UERadioCapability = new(ngapType.UERadioCapability)
		uecapa, err := hex.DecodeString(amfUe.UeRadioCapability)
		if err != nil {
			return nil, err
		}
		ie.Value.UERadioCapability.Value = uecapa
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// Index to RAT/Frequency Selection Priority (optional)
	if amfUe.AmPolicyAssociation != nil && amfUe.AmPolicyAssociation.Rfsp != 0 {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDIndexToRFSP
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentIndexToRFSP
		ie.Value.IndexToRFSP = new(ngapType.IndexToRFSP)

		ie.Value.IndexToRFSP.Value = int64(amfUe.AmPolicyAssociation.Rfsp)

		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// Masked IMEISV (optional)
	// TS 38.413 9.3.1.54; TS 23.003 6.2; TS 23.501 5.9.3
	// last 4 digits of the SNR masked by setting the corresponding bits to 1.
	// The first to fourth bits correspond to the first digit of the IMEISV,
	// the fifth to eighth bits correspond to the second digit of the IMEISV, and so on
	if c := factory.AmfConfig.GetNgapIEMaskedIMEISV(); c != nil && c.Enable &&
		amfUe.Pei != "" && strings.HasPrefix(amfUe.Pei, "imeisv") {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDMaskedIMEISV
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentMaskedIMEISV
		ie.Value.MaskedIMEISV = new(ngapType.MaskedIMEISV)

		imeisv := strings.TrimPrefix(amfUe.Pei, "imeisv-")
		imeisvBytes, err := hex.DecodeString(imeisv)
		if err != nil {
			logger.NgapLog.Errorf("[Build Error] DecodeString imeisv error: %+v", err)
		}

		var maskedImeisv []byte
		maskedImeisv = append(maskedImeisv, imeisvBytes[:5]...)
		maskedImeisv = append(maskedImeisv, []byte{0xff, 0xff}...)
		maskedImeisv = append(maskedImeisv, imeisvBytes[7])
		ie.Value.MaskedIMEISV.Value = aper.BitString{
			BitLength: 64,
			Bytes:     maskedImeisv,
		}
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// NAS-PDU (optional)
	if nasPdu != nil {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNASPDU
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentNASPDU
		ie.Value.NASPDU = new(ngapType.NASPDU)

		ie.Value.NASPDU.Value = nasPdu

		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// Emergency Fallback indicator (optional)
	if emergencyFallbackIndicator != nil {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDEmergencyFallbackIndicator
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentEmergencyFallbackIndicator
		ie.Value.EmergencyFallbackIndicator = emergencyFallbackIndicator
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// RRC Inactive Transition Report Request (optional)
	if rrcInactiveTransitionReportRequest != nil {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRRCInactiveTransitionReportRequest
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentRRCInactiveTransitionReportRequest
		ie.Value.RRCInactiveTransitionReportRequest = rrcInactiveTransitionReportRequest
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// UE Radio Capability for Paging (optional)
	if amfUe.UeRadioCapabilityForPaging != nil {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUERadioCapabilityForPaging
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentUERadioCapabilityForPaging
		ie.Value.UERadioCapabilityForPaging = new(ngapType.UERadioCapabilityForPaging)
		uERadioCapabilityForPaging := ie.Value.UERadioCapabilityForPaging
		var err error
		if amfUe.UeRadioCapabilityForPaging.NR != "" {
			uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR.Value, err = hex.
				DecodeString(amfUe.UeRadioCapabilityForPaging.NR)
			if err != nil {
				logger.NgapLog.Errorf("[Build Error] DecodeString amfUe.UeRadioCapabilityForPaging.NR error: %+v", err)
			}
		}
		if amfUe.UeRadioCapabilityForPaging.EUTRA != "" {
			uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA.Value, err = hex.
				DecodeString(amfUe.UeRadioCapabilityForPaging.EUTRA)
			if err != nil {
				logger.NgapLog.Errorf("[Build Error] DecodeString amfUe.UeRadioCapabilityForPaging.NR error: %+v", err)
			}
		}
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	// Redirection for Voice EPS Fallback (optional)
	if c := factory.AmfConfig.GetNgapIERedirectionVoiceFallback(); c != nil && c.Enable {
		ie = ngapType.InitialContextSetupRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRedirectionVoiceFallback
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupRequestIEsPresentRedirectionVoiceFallback
		ie.Value.RedirectionVoiceFallback = new(ngapType.RedirectionVoiceFallback)
		ie.Value.RedirectionVoiceFallback.Value = ngapType.RedirectionVoiceFallbackPresentNotPossible
		initialContextSetupRequestIEs.List = append(initialContextSetupRequestIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildUEContextModificationRequest(
	amfUe *context.AmfUe,
	anType models.AccessType,
	oldAmfUeNgapID *int64,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	coreNetworkAssistanceInfo *ngapType.CoreNetworkAssistanceInformation,
	mobilityRestrictionList *ngapType.MobilityRestrictionList,
	emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator,
) ([]byte, error) {
	// accessType indicate amfUe send this msg for which accessType
	// oldAmfUeNgapID: if amf allocate a new amf ue ngap id to amfUe, the caller should
	// update the context by itself, and pass the old AmfUeNgapID to this function
	// for other parameters, please reference the comments in BuildInitialContextSetupRequest

	// TODO: Ran Paging Priority (optional) [int: 1~256] TS 38.413 9.3.3.15, TS 23.501
	// TODO: fill IE securityKey & ueSecurityCapabilities to code

	if amfUe == nil {
		return nil, fmt.Errorf("amfUe is nil")
	}

	ranUe, ok := amfUe.RanUe[anType]
	if !ok {
		return nil, fmt.Errorf("ranUe for %s is nil", anType)
	}

	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUEContextModification
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUEContextModificationRequest
	initiatingMessage.Value.UEContextModificationRequest = new(ngapType.UEContextModificationRequest)

	uEContextModificationRequest := initiatingMessage.Value.UEContextModificationRequest
	uEContextModificationRequestIEs := &uEContextModificationRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextModificationRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextModificationRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	if oldAmfUeNgapID != nil {
		aMFUENGAPID.Value = *oldAmfUeNgapID
	} else {
		aMFUENGAPID.Value = ranUe.AmfUeNgapId
	}

	uEContextModificationRequestIEs.List = append(uEContextModificationRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextModificationRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextModificationRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUe.RanUeNgapId

	uEContextModificationRequestIEs.List = append(uEContextModificationRequestIEs.List, ie)

	// Ran Paging Priority (optional)

	// Security Key (optional)

	// Index to RAT/Frequency Selection Priority (optional)
	if amfUe.AmPolicyAssociation != nil && amfUe.AmPolicyAssociation.Rfsp != 0 {
		ie = ngapType.UEContextModificationRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDIndexToRFSP
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UEContextModificationRequestIEsPresentIndexToRFSP
		ie.Value.IndexToRFSP = new(ngapType.IndexToRFSP)

		ie.Value.IndexToRFSP.Value = int64(amfUe.AmPolicyAssociation.Rfsp)

		uEContextModificationRequestIEs.List = append(uEContextModificationRequestIEs.List, ie)
	}

	// UE Aggregate Maximum Bit Rate (optional)
	if amfUe.AccessAndMobilitySubscriptionData != nil &&
		amfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr != nil {
		ie = ngapType.UEContextModificationRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUEAggregateMaximumBitRate
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UEContextModificationRequestIEsPresentUEAggregateMaximumBitRate
		ie.Value.UEAggregateMaximumBitRate = new(ngapType.UEAggregateMaximumBitRate)

		ueAmbrUL := ngapConvert.UEAmbrToInt64(amfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Uplink)
		ueAmbrDL := ngapConvert.UEAmbrToInt64(amfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Downlink)
		ie.Value.UEAggregateMaximumBitRate.UEAggregateMaximumBitRateUL.Value = ueAmbrUL
		ie.Value.UEAggregateMaximumBitRate.UEAggregateMaximumBitRateDL.Value = ueAmbrDL

		uEContextModificationRequestIEs.List = append(uEContextModificationRequestIEs.List, ie)
	}

	// UE Security Capabilities (optional)

	// Core Network Assistance Information (optional)
	if coreNetworkAssistanceInfo != nil {
		ie = ngapType.UEContextModificationRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCoreNetworkAssistanceInformation
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UEContextModificationRequestIEsPresentCoreNetworkAssistanceInformation
		ie.Value.CoreNetworkAssistanceInformation = coreNetworkAssistanceInfo
		uEContextModificationRequestIEs.List = append(uEContextModificationRequestIEs.List, ie)
	}

	// Emergency Fallback Indicator (optional)
	if emergencyFallbackIndicator != nil {
		ie = ngapType.UEContextModificationRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDEmergencyFallbackIndicator
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.UEContextModificationRequestIEsPresentEmergencyFallbackIndicator
		ie.Value.EmergencyFallbackIndicator = emergencyFallbackIndicator
		uEContextModificationRequestIEs.List = append(uEContextModificationRequestIEs.List, ie)
	}

	// New AMF UE NGAP ID (optional)
	if oldAmfUeNgapID != nil {
		ie = ngapType.UEContextModificationRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNewAMFUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UEContextModificationRequestIEsPresentNewAMFUENGAPID
		ie.Value.NewAMFUENGAPID = new(ngapType.AMFUENGAPID)

		ie.Value.NewAMFUENGAPID.Value = ranUe.AmfUeNgapId

		uEContextModificationRequestIEs.List = append(uEContextModificationRequestIEs.List, ie)
	}

	// RRC Inactive Transition Report Request (optional)
	if rrcInactiveTransitionReportRequest != nil {
		ie = ngapType.UEContextModificationRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRRCInactiveTransitionReportRequest
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UEContextModificationRequestIEsPresentRRCInactiveTransitionReportRequest
		ie.Value.RRCInactiveTransitionReportRequest = rrcInactiveTransitionReportRequest
		uEContextModificationRequestIEs.List = append(uEContextModificationRequestIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

// pduSessionResourceHandoverList: provided by amf and transfer is return from smf
// pduSessionResourceToReleaseList: provided by amf and transfer is return from smf
// criticalityDiagnostics = criticalityDiagonstics IE in receiver node's error indication
// when received node can't comprehend the IE or missing IE
func BuildHandoverCommand(
	sourceUe *context.RanUe,
	pduSessionResourceHandoverList ngapType.PDUSessionResourceHandoverList,
	pduSessionResourceToReleaseList ngapType.PDUSessionResourceToReleaseListHOCmd,
	container ngapType.TargetToSourceTransparentContainer,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeHandoverPreparation
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentHandoverCommand
	successfulOutcome.Value.HandoverCommand = new(ngapType.HandoverCommand)

	handoverCommand := successfulOutcome.Value.HandoverCommand
	handoverCommandIEs := &handoverCommand.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.HandoverCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverCommandIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = sourceUe.AmfUeNgapId

	handoverCommandIEs.List = append(handoverCommandIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.HandoverCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = sourceUe.RanUeNgapId

	handoverCommandIEs.List = append(handoverCommandIEs.List, ie)

	// Handover Type
	ie = ngapType.HandoverCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDHandoverType
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverCommandIEsPresentHandoverType
	ie.Value.HandoverType = new(ngapType.HandoverType)

	handoverType := ie.Value.HandoverType
	handoverType.Value = sourceUe.HandOverType.Value

	handoverCommandIEs.List = append(handoverCommandIEs.List, ie)

	// NAS Security Parameters from NG-RAN [C-iftoEPS]
	if handoverType.Value == ngapType.HandoverTypePresentFivegsToEps {
		ie = ngapType.HandoverCommandIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNASSecurityParametersFromNGRAN
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.HandoverCommandIEsPresentNASSecurityParametersFromNGRAN
		ie.Value.NASSecurityParametersFromNGRAN = new(ngapType.NASSecurityParametersFromNGRAN)

		handoverCommandIEs.List = append(handoverCommandIEs.List, ie)
	}

	// PDU Session Resource Handover List
	ie = ngapType.HandoverCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceHandoverList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverCommandIEsPresentPDUSessionResourceHandoverList
	ie.Value.PDUSessionResourceHandoverList = &pduSessionResourceHandoverList
	handoverCommandIEs.List = append(handoverCommandIEs.List, ie)

	// PDU Session Resource to Release List
	if len(pduSessionResourceToReleaseList.List) > 0 {
		ie = ngapType.HandoverCommandIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceToReleaseListHOCmd
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.HandoverCommandIEsPresentPDUSessionResourceToReleaseListHOCmd
		ie.Value.PDUSessionResourceToReleaseListHOCmd = &pduSessionResourceToReleaseList
		handoverCommandIEs.List = append(handoverCommandIEs.List, ie)
	}

	// Target to Source Transparent Container
	ie = ngapType.HandoverCommandIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDTargetToSourceTransparentContainer
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverCommandIEsPresentTargetToSourceTransparentContainer
	ie.Value.TargetToSourceTransparentContainer = &container

	handoverCommandIEs.List = append(handoverCommandIEs.List, ie)

	// Criticality Diagnostics [optional]
	if criticalityDiagnostics != nil {
		handoverCommandIEsie := ngapType.HandoverCommandIEs{}
		handoverCommandIEsie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		handoverCommandIEsie.Criticality.Value = ngapType.CriticalityPresentIgnore
		handoverCommandIEsie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentCriticalityDiagnostics
		handoverCommandIEsie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		handoverCommandIEsie.Value.CriticalityDiagnostics = criticalityDiagnostics

		handoverCommandIEs.List = append(handoverCommandIEs.List, handoverCommandIEsie)
	}

	return ngap.Encoder(pdu)
}

func BuildHandoverPreparationFailure(sourceUe *context.RanUe, cause ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	// cause = initiate the Handover Cancel procedure with the appropriate value for the Cause IE.

	// criticalityDiagnostics = criticalityDiagonstics IE in receiver node's error indication
	// when received node can't comprehend the IE or missing IE

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeHandoverPreparation
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentHandoverPreparationFailure
	unsuccessfulOutcome.Value.HandoverPreparationFailure = new(ngapType.HandoverPreparationFailure)

	handoverPreparationFailure := unsuccessfulOutcome.Value.HandoverPreparationFailure
	handoverPreparationFailureIEs := &handoverPreparationFailure.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.HandoverPreparationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = sourceUe.AmfUeNgapId

	handoverPreparationFailureIEs.List = append(handoverPreparationFailureIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.HandoverPreparationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = sourceUe.RanUeNgapId

	handoverPreparationFailureIEs.List = append(handoverPreparationFailureIEs.List, ie)

	// Cause
	ie = ngapType.HandoverPreparationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentCriticalityDiagnostics
	ie.Value.Cause = new(ngapType.Cause)

	ie.Value.Cause = &cause

	handoverPreparationFailureIEs.List = append(handoverPreparationFailureIEs.List, ie)

	// Criticality Diagnostics [optional]
	if criticalityDiagnostics != nil {
		HandoverPreparationFailureIEsie := ngapType.HandoverPreparationFailureIEs{}
		HandoverPreparationFailureIEsie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		HandoverPreparationFailureIEsie.Criticality.Value = ngapType.CriticalityPresentIgnore
		HandoverPreparationFailureIEsie.Value.Present = ngapType.HandoverCancelAcknowledgeIEsPresentCriticalityDiagnostics
		HandoverPreparationFailureIEsie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		HandoverPreparationFailureIEsie.Value.CriticalityDiagnostics = criticalityDiagnostics

		handoverPreparationFailureIEs.List = append(handoverPreparationFailureIEs.List, HandoverPreparationFailureIEsie)
	}

	return ngap.Encoder(pdu)
}

/*The PGW-C+SMF (V-SMF in the case of home-routed roaming scenario only) sends
a Nsmf_PDUSession_CreateSMContext Response(N2 SM Information (PDU Session ID, cause code)) to the AMF.*/
// Cause is from SMF
// pduSessionResourceSetupList provided by AMF, and the transfer data is from SMF
// sourceToTargetTransparentContainer is received from S-RAN
// nsci: new security context indicator, if amfUe has updated security context,
// set nsci to true, otherwise set to false
func BuildHandoverRequest(ue *context.RanUe, cause ngapType.Cause,
	pduSessionResourceSetupListHOReq ngapType.PDUSessionResourceSetupListHOReq,
	sourceToTargetTransparentContainer ngapType.SourceToTargetTransparentContainer, nsci bool,
) ([]byte, error) {
	amfSelf := context.GetSelf()
	amfUe := ue.AmfUe
	if amfUe == nil {
		return nil, fmt.Errorf("AmfUe is nil")
	}

	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeHandoverResourceAllocation
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentHandoverRequest
	initiatingMessage.Value.HandoverRequest = new(ngapType.HandoverRequest)

	handoverRequest := initiatingMessage.Value.HandoverRequest
	handoverRequestIEs := &handoverRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// Handover Type
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDHandoverType
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentHandoverType
	ie.Value.HandoverType = new(ngapType.HandoverType)

	handoverType := ie.Value.HandoverType
	handoverType.Value = ue.HandOverType.Value

	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// Cause
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.HandoverRequestIEsPresentCause
	ie.Value.Cause = &cause

	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// UE Aggregate Maximum Bit Rate
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUEAggregateMaximumBitRate
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentUEAggregateMaximumBitRate
	ie.Value.UEAggregateMaximumBitRate = new(ngapType.UEAggregateMaximumBitRate)

	ueAmbrUL := ngapConvert.UEAmbrToInt64(amfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Uplink)
	ueAmbrDL := ngapConvert.UEAmbrToInt64(amfUe.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Downlink)
	ie.Value.UEAggregateMaximumBitRate.UEAggregateMaximumBitRateUL.Value = ueAmbrUL
	ie.Value.UEAggregateMaximumBitRate.UEAggregateMaximumBitRateDL.Value = ueAmbrDL

	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// UE Security Capabilities
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUESecurityCapabilities
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentUESecurityCapabilities
	ie.Value.UESecurityCapabilities = new(ngapType.UESecurityCapabilities)

	ueSecurityCapabilities := ie.Value.UESecurityCapabilities

	nrEncryptionAlgorighm := []byte{0x00, 0x00}
	nrEncryptionAlgorighm[0] |= amfUe.UESecurityCapability.GetEA1_128_5G() << 7
	nrEncryptionAlgorighm[0] |= amfUe.UESecurityCapability.GetEA2_128_5G() << 6
	nrEncryptionAlgorighm[0] |= amfUe.UESecurityCapability.GetEA3_128_5G() << 5
	ueSecurityCapabilities.NRencryptionAlgorithms.Value = ngapConvert.ByteToBitString(nrEncryptionAlgorighm, 16)

	nrIntegrityAlgorithm := []byte{0x00, 0x00}
	nrIntegrityAlgorithm[0] |= amfUe.UESecurityCapability.GetIA1_128_5G() << 7
	nrIntegrityAlgorithm[0] |= amfUe.UESecurityCapability.GetIA2_128_5G() << 6
	nrIntegrityAlgorithm[0] |= amfUe.UESecurityCapability.GetIA3_128_5G() << 5
	ueSecurityCapabilities.NRintegrityProtectionAlgorithms.Value = ngapConvert.ByteToBitString(nrIntegrityAlgorithm, 16)

	// only support NR algorithms
	eutraEncryptionAlgorithm := []byte{0x00, 0x00}
	ueSecurityCapabilities.EUTRAencryptionAlgorithms.Value = ngapConvert.ByteToBitString(eutraEncryptionAlgorithm, 16)

	eutraIntegrityAlgorithm := []byte{0x00, 0x00}
	ueSecurityCapabilities.EUTRAintegrityProtectionAlgorithms.Value = ngapConvert.
		ByteToBitString(eutraIntegrityAlgorithm, 16)

	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// Security Context
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSecurityContext
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentSecurityContext
	ie.Value.SecurityContext = new(ngapType.SecurityContext)

	securityContext := ie.Value.SecurityContext
	securityContext.NextHopChainingCount.Value = int64(ue.AmfUe.NCC)
	securityContext.NextHopNH.Value = ngapConvert.HexToBitString(hex.EncodeToString(ue.AmfUe.NH), 256)

	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// PDU Session Resource Setup List
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListHOReq
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentPDUSessionResourceSetupListHOReq
	ie.Value.PDUSessionResourceSetupListHOReq = &pduSessionResourceSetupListHOReq
	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// Allowed NSSAI
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAllowedNSSAI
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentAllowedNSSAI
	ie.Value.AllowedNSSAI = new(ngapType.AllowedNSSAI)

	allowedNSSAI := ie.Value.AllowedNSSAI
	for _, snssaiItem := range amfSelf.PlmnSupportList[0].SNssaiList {
		allowedNSSAIItem := ngapType.AllowedNSSAIItem{}

		ngapSnssai := ngapConvert.SNssaiToNgap(snssaiItem)
		allowedNSSAIItem.SNSSAI = ngapSnssai
		allowedNSSAI.List = append(allowedNSSAI.List, allowedNSSAIItem)
	}
	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// Masked IMEISV (optional)
	// TS 38.413 9.3.1.54; TS 23.003 6.2; TS 23.501 5.9.3
	// last 4 digits of the SNR masked by setting the corresponding bits to 1.
	// The first to fourth bits correspond to the first digit of the IMEISV,
	// the fifth to eighth bits correspond to the second digit of the IMEISV, and so on
	if c := factory.AmfConfig.GetNgapIEMaskedIMEISV(); c != nil && c.Enable &&
		amfUe.Pei != "" && strings.HasPrefix(amfUe.Pei, "imeisv") {
		ie = ngapType.HandoverRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDMaskedIMEISV
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.HandoverRequestIEsPresentMaskedIMEISV
		ie.Value.MaskedIMEISV = new(ngapType.MaskedIMEISV)

		imeisv := strings.TrimPrefix(amfUe.Pei, "imeisv-")
		imeisvBytes, err := hex.DecodeString(imeisv)
		if err != nil {
			logger.NgapLog.Errorf("[Build Error] DecodeString imeisv error: %+v", err)
		}

		var maskedImeisv []byte
		maskedImeisv = append(maskedImeisv, imeisvBytes[:5]...)
		maskedImeisv = append(maskedImeisv, []byte{0xff, 0xff}...)
		maskedImeisv = append(maskedImeisv, imeisvBytes[7])
		ie.Value.MaskedIMEISV.Value = aper.BitString{
			BitLength: 64,
			Bytes:     maskedImeisv,
		}
		handoverRequestIEs.List = append(handoverRequestIEs.List, ie)
	}

	// Source To Target Transparent Container
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSourceToTargetTransparentContainer
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentSourceToTargetTransparentContainer
	ie.Value.SourceToTargetTransparentContainer = new(ngapType.SourceToTargetTransparentContainer)

	sourceToTargetTransparentContaine := ie.Value.SourceToTargetTransparentContainer
	sourceToTargetTransparentContaine.Value = sourceToTargetTransparentContainer.Value

	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// Mobility Restriction List (optional)
	if c := factory.AmfConfig.GetNgapIEMobilityRestrictionList(); c != nil && c.Enable {
		ie = ngapType.HandoverRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDMobilityRestrictionList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.HandoverRequestIEsPresentMobilityRestrictionList
		ie.Value.MobilityRestrictionList = new(ngapType.MobilityRestrictionList)

		mobilityRestrictionList := BuildIEMobilityRestrictionList(amfUe)
		ie.Value.MobilityRestrictionList = &mobilityRestrictionList
		handoverRequestIEs.List = append(handoverRequestIEs.List, ie)
	}

	// GUAMI
	ie = ngapType.HandoverRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDGUAMI
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.HandoverRequestIEsPresentGUAMI
	ie.Value.GUAMI = new(ngapType.GUAMI)

	guami := ie.Value.GUAMI
	plmnID := &guami.PLMNIdentity
	amfRegionID := &guami.AMFRegionID
	amfSetID := &guami.AMFSetID
	amfPtrID := &guami.AMFPointer

	servedGuami := amfSelf.ServedGuamiList[0]

	*plmnID = ngapConvert.PlmnIdToNgap(util.PlmnIdNidToModelsPlmnId(*servedGuami.PlmnId))
	amfRegionID.Value, amfSetID.Value, amfPtrID.Value = ngapConvert.AmfIdToNgap(servedGuami.AmfId)

	handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// //Core Network Assistance Information(optional)
	// ie = ngapType.HandoverRequestIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDCoreNetworkAssistanceInformation
	// ie.Criticality.Value = ngapType.CriticalityPresentReject
	// ie.Value.Present = ngapType.HandoverRequestIEsPresentCoreNetworkAssistanceInformation
	// ie.Value.CoreNetworkAssistanceInformation = new(ngapType.CoreNetworkAssistanceInformation)
	// handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// New Security ContextInd(optional)
	if nsci {
		ie = ngapType.HandoverRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNewSecurityContextInd
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.HandoverRequestIEsPresentNewSecurityContextInd
		ie.Value.NewSecurityContextInd = new(ngapType.NewSecurityContextInd)
		ie.Value.NewSecurityContextInd.Value = ngapType.NewSecurityContextIndPresentTrue
		handoverRequestIEs.List = append(handoverRequestIEs.List, ie)
	}

	// NASC(optional)
	// ie.Criticality.Value = ngapType.CriticalityPresentReject
	// ie.Value.Present = ngapType.HandoverRequestIEsPresentNASC
	// ie.Id.Value = ngapType.ProtocolIEIDNASC
	// ie.Criticality.Value = ngapType.CriticalityPresentReject
	// ie.Value.Present = ngapType.HandoverRequestIEsPresentNASC
	// ie.Value.NASC = new(ngapType.)
	// handoverRequestIEs.List = append(handoverRequestIEs.List, ie)

	// Trace Activation(optional)
	// Masked IMEISV(optional)
	// Mobility Restriction List(optional)
	// Location Reporting Request Type(optional)
	// RRC Inactive Transition Report Reques(optional)

	// Redirection for Voice EPS Fallback (optional)
	if c := factory.AmfConfig.GetNgapIERedirectionVoiceFallback(); c != nil && c.Enable {
		ie = ngapType.HandoverRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRedirectionVoiceFallback
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.HandoverRequestIEsPresentRedirectionVoiceFallback
		ie.Value.RedirectionVoiceFallback = new(ngapType.RedirectionVoiceFallback)
		ie.Value.RedirectionVoiceFallback.Value = ngapType.RedirectionVoiceFallbackPresentNotPossible
		handoverRequestIEs.List = append(handoverRequestIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

// pduSessionResourceSwitchedList: provided by AMF, and the transfer data is from SMF
// pduSessionResourceReleasedList: provided by AMF, and the transfer data is from SMF
// newSecurityContextIndicator: if AMF has activated a new 5G NAS security context,
// set it to true, otherwise set to false
// coreNetworkAssistanceInformation: provided by AMF,
// based on collection of UE behavior statistics and/or other available
// information about the expected UE behavior. TS 23.501 5.4.6, 5.4.6.2
// rrcInactiveTransitionReportRequest: configured by amf
// criticalityDiagnostics: from received node when received not comprehended IE or missing IE
func BuildPathSwitchRequestAcknowledge(
	ue *context.RanUe,
	pduSessionResourceSwitchedList ngapType.PDUSessionResourceSwitchedList,
	pduSessionResourceReleasedList ngapType.PDUSessionResourceReleasedListPSAck,
	newSecurityContextIndicator bool,
	coreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	amfSelf := context.GetSelf()

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePathSwitchRequest
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPathSwitchRequestAcknowledge
	successfulOutcome.Value.PathSwitchRequestAcknowledge = new(ngapType.PathSwitchRequestAcknowledge)

	pathSwitchRequestAck := successfulOutcome.Value.PathSwitchRequestAcknowledge
	pathSwitchRequestAckIEs := &pathSwitchRequestAck.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PathSwitchRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)

	// UE Security Capabilities (optional)
	ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUESecurityCapabilities
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentUESecurityCapabilities
	ie.Value.UESecurityCapabilities = new(ngapType.UESecurityCapabilities)

	ueSecurityCapabilities := ie.Value.UESecurityCapabilities
	nrEncryptionAlgorighm := []byte{0x00, 0x00}
	nrEncryptionAlgorighm[0] |= ue.AmfUe.UESecurityCapability.GetEA1_128_5G() << 7
	nrEncryptionAlgorighm[0] |= ue.AmfUe.UESecurityCapability.GetEA2_128_5G() << 6
	nrEncryptionAlgorighm[0] |= ue.AmfUe.UESecurityCapability.GetEA3_128_5G() << 5
	ueSecurityCapabilities.NRencryptionAlgorithms.Value = ngapConvert.ByteToBitString(nrEncryptionAlgorighm, 16)

	nrIntegrityAlgorithm := []byte{0x00, 0x00}
	nrIntegrityAlgorithm[0] |= ue.AmfUe.UESecurityCapability.GetIA1_128_5G() << 7
	nrIntegrityAlgorithm[0] |= ue.AmfUe.UESecurityCapability.GetIA2_128_5G() << 6
	nrIntegrityAlgorithm[0] |= ue.AmfUe.UESecurityCapability.GetIA3_128_5G() << 5
	ueSecurityCapabilities.NRintegrityProtectionAlgorithms.Value = ngapConvert.ByteToBitString(nrIntegrityAlgorithm, 16)

	// only support NR algorithms
	eutraEncryptionAlgorithm := []byte{0x00, 0x00}
	ueSecurityCapabilities.EUTRAencryptionAlgorithms.Value = ngapConvert.
		ByteToBitString(eutraEncryptionAlgorithm, 16)

	eutraIntegrityAlgorithm := []byte{0x00, 0x00}
	ueSecurityCapabilities.EUTRAintegrityProtectionAlgorithms.Value = ngapConvert.
		ByteToBitString(eutraIntegrityAlgorithm, 16)

	pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)

	// Security Context
	ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSecurityContext
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentSecurityContext
	ie.Value.SecurityContext = new(ngapType.SecurityContext)

	securityContext := ie.Value.SecurityContext
	securityContext.NextHopChainingCount.Value = int64(ue.AmfUe.NCC)
	securityContext.NextHopNH.Value = ngapConvert.HexToBitString(hex.EncodeToString(ue.AmfUe.NH), 256)

	pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)

	// New Security Context Indicator (optional)
	if newSecurityContextIndicator {
		ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNewSecurityContextInd
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentNewSecurityContextInd
		ie.Value.NewSecurityContextInd = new(ngapType.NewSecurityContextInd)
		ie.Value.NewSecurityContextInd.Value = ngapType.NewSecurityContextIndPresentTrue
		pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)
	}

	// PDU Session Resource Switched List
	ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSwitchedList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentPDUSessionResourceSwitchedList
	ie.Value.PDUSessionResourceSwitchedList = &pduSessionResourceSwitchedList
	pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)

	// PDU Session Resource Released List
	if len(pduSessionResourceReleasedList.List) > 0 {
		ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceReleasedListPSAck
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentPDUSessionResourceReleasedListPSAck
		ie.Value.PDUSessionResourceReleasedListPSAck = &pduSessionResourceReleasedList
		pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)
	}

	// Allowed NSSAI
	ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAllowedNSSAI
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentAllowedNSSAI
	ie.Value.AllowedNSSAI = new(ngapType.AllowedNSSAI)

	allowedNSSAI := ie.Value.AllowedNSSAI
	// plmnSupportList[0] is serving plmn
	for _, modelSnssai := range amfSelf.PlmnSupportList[0].SNssaiList {
		allowedNSSAIItem := ngapType.AllowedNSSAIItem{}

		ngapSnssai := ngapConvert.SNssaiToNgap(modelSnssai)
		allowedNSSAIItem.SNSSAI = ngapSnssai
		allowedNSSAI.List = append(allowedNSSAI.List, allowedNSSAIItem)
	}
	pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)

	// Core Network Assistance Information (optional)
	if coreNetworkAssistanceInformation != nil {
		ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCoreNetworkAssistanceInformation
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentCoreNetworkAssistanceInformation
		ie.Value.CoreNetworkAssistanceInformation = coreNetworkAssistanceInformation
		pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)
	}

	// RRC Inactive Transition Report Request (optional)
	if rrcInactiveTransitionReportRequest != nil {
		ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRRCInactiveTransitionReportRequest
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentRRCInactiveTransitionReportRequest
		ie.Value.RRCInactiveTransitionReportRequest = rrcInactiveTransitionReportRequest
		pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)
	}

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)
	}

	// Redirection for Voice EPS Fallback (optional)
	if c := factory.AmfConfig.GetNgapIERedirectionVoiceFallback(); c != nil && c.Enable {
		ie = ngapType.PathSwitchRequestAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRedirectionVoiceFallback
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PathSwitchRequestAcknowledgeIEsPresentRedirectionVoiceFallback
		ie.Value.RedirectionVoiceFallback = new(ngapType.RedirectionVoiceFallback)
		ie.Value.RedirectionVoiceFallback.Value = ngapType.RedirectionVoiceFallbackPresentNotPossible
		pathSwitchRequestAckIEs.List = append(pathSwitchRequestAckIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

// pduSessionResourceReleasedList: provided by AMF, and the transfer data is from SMF
// criticalityDiagnostics: from received node when received not comprehended IE or missing IE
func BuildPathSwitchRequestFailure(
	amfUeNgapId,
	ranUeNgapId int64,
	pduSessionResourceReleasedList *ngapType.PDUSessionResourceReleasedListPSFail,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePathSwitchRequest
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentPathSwitchRequestFailure
	unsuccessfulOutcome.Value.PathSwitchRequestFailure = new(ngapType.PathSwitchRequestFailure)

	pathSwitchRequestFailure := unsuccessfulOutcome.Value.PathSwitchRequestFailure
	pathSwitchRequestFailureIEs := &pathSwitchRequestFailure.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PathSwitchRequestFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PathSwitchRequestFailureIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapId

	pathSwitchRequestFailureIEs.List = append(pathSwitchRequestFailureIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PathSwitchRequestFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PathSwitchRequestFailureIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapId

	pathSwitchRequestFailureIEs.List = append(pathSwitchRequestFailureIEs.List, ie)

	// PDU Session Resource Released List
	if pduSessionResourceReleasedList != nil {
		ie = ngapType.PathSwitchRequestFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceReleasedListPSFail
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PathSwitchRequestFailureIEsPresentPDUSessionResourceReleasedListPSFail
		ie.Value.PDUSessionResourceReleasedListPSFail = pduSessionResourceReleasedList
		pathSwitchRequestFailureIEs.List = append(pathSwitchRequestFailureIEs.List, ie)
	}

	if criticalityDiagnostics != nil {
		ie = ngapType.PathSwitchRequestFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PathSwitchRequestFailureIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		pathSwitchRequestFailureIEs.List = append(pathSwitchRequestFailureIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildDownlinkRanStatusTransfer(ue *context.RanUe,
	ranStatusTransferTransparentContainer ngapType.RANStatusTransferTransparentContainer,
) ([]byte, error) {
	// ranStatusTransferTransparentContainer from Uplink Ran Configuration Transfer
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeDownlinkRANStatusTransfer
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore
	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentDownlinkRANStatusTransfer
	initiatingMessage.Value.DownlinkRANStatusTransfer = new(ngapType.DownlinkRANStatusTransfer)

	downlinkRanStatusTransfer := initiatingMessage.Value.DownlinkRANStatusTransfer
	downlinkRanStatusTransferIEs := &downlinkRanStatusTransfer.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.DownlinkRANStatusTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkRANStatusTransferIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	downlinkRanStatusTransferIEs.List = append(downlinkRanStatusTransferIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.DownlinkRANStatusTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkRANStatusTransferIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	downlinkRanStatusTransferIEs.List = append(downlinkRanStatusTransferIEs.List, ie)

	// RAN Status Transfer Transparent Container
	ie = ngapType.DownlinkRANStatusTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANStatusTransferTransparentContainer
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkRANStatusTransferIEsPresentRANStatusTransferTransparentContainer

	ie.Value.RANStatusTransferTransparentContainer = &ranStatusTransferTransparentContainer

	downlinkRanStatusTransferIEs.List = append(downlinkRanStatusTransferIEs.List, ie)

	return ngap.Encoder(pdu)
}

// anType indicate amfUe send this msg for which accessType
// Paging Priority: is included only if the AMF receives an Namf_Communication_N1N2MessageTransfer message
// with an ARP value associated with
// priority services (e.g., MPS, MCS), as configured by the operator. (TS 23.502 4.2.3.3, TS 23.501 5.22.3)
// pagingOriginNon3GPP: TS 23.502 4.2.3.3 step 4b: If the UE is simultaneously registered over
// 3GPP and non-3GPP accesses in the same PLMN,
// the UE is in CM-IDLE state in both 3GPP access and non-3GPP access, and the PDU Session ID in step 3a
// is associated with non-3GPP access, the AMF sends a Paging message with associated access "non-3GPP" to
// NG-RAN node(s) via 3GPP access.
// more paging policy with 3gpp/non-3gpp access is described in TS 23.501 5.6.8
func BuildPaging(
	ue *context.AmfUe, pagingPriority *ngapType.PagingPriority, pagingOriginNon3GPP bool,
) ([]byte, error) {
	// TODO: Paging DRX (optional)

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePaging
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPaging
	initiatingMessage.Value.Paging = new(ngapType.Paging)

	paging := initiatingMessage.Value.Paging
	pagingIEs := &paging.ProtocolIEs

	// UE Paging Identity
	ie := ngapType.PagingIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUEPagingIdentity
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PagingIEsPresentUEPagingIdentity
	ie.Value.UEPagingIdentity = new(ngapType.UEPagingIdentity)

	uePagingIdentity := ie.Value.UEPagingIdentity
	uePagingIdentity.Present = ngapType.UEPagingIdentityPresentFiveGSTMSI
	uePagingIdentity.FiveGSTMSI = new(ngapType.FiveGSTMSI)

	var amfID string
	var tmsi string
	if len(ue.Guti) == 19 {
		amfID = ue.Guti[5:11]
		tmsi = ue.Guti[11:]
	} else {
		amfID = ue.Guti[6:12]
		tmsi = ue.Guti[12:]
	}
	_, amfSetID, amfPointer := ngapConvert.AmfIdToNgap(amfID)

	var err error
	uePagingIdentity.FiveGSTMSI.AMFSetID.Value = amfSetID
	uePagingIdentity.FiveGSTMSI.AMFPointer.Value = amfPointer
	uePagingIdentity.FiveGSTMSI.FiveGTMSI.Value, err = hex.DecodeString(tmsi)
	if err != nil {
		logger.NgapLog.Errorf("[Build Error] DecodeString tmsi error: %+v", err)
	}

	pagingIEs.List = append(pagingIEs.List, ie)

	// Paging DRX (optional)

	// TAI List for Paging
	ie = ngapType.PagingIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDTAIListForPaging
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PagingIEsPresentTAIListForPaging
	ie.Value.TAIListForPaging = new(ngapType.TAIListForPaging)

	taiListForPaging := ie.Value.TAIListForPaging
	if ue.RegistrationArea[models.AccessType__3_GPP_ACCESS] == nil {
		err = fmt.Errorf("registration area of Ue[%s] is empty", ue.Supi)
		return nil, err
	} else {
		for _, tai := range ue.RegistrationArea[models.AccessType__3_GPP_ACCESS] {
			var tac []byte
			taiListforPagingItem := ngapType.TAIListForPagingItem{}
			taiListforPagingItem.TAI.PLMNIdentity = ngapConvert.PlmnIdToNgap(*tai.PlmnId)
			tac, err = hex.DecodeString(tai.Tac)
			if err != nil {
				logger.NgapLog.Errorf("[Build Error] DecodeString tai.Tac error: %+v", err)
			}
			taiListforPagingItem.TAI.TAC.Value = tac
			taiListForPaging.List = append(taiListForPaging.List, taiListforPagingItem)
		}
	}

	pagingIEs.List = append(pagingIEs.List, ie)

	// Paging Priority (optional)
	if pagingPriority != nil {
		ie = ngapType.PagingIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPagingPriority
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PagingIEsPresentPagingPriority
		ie.Value.PagingPriority = pagingPriority
		pagingIEs.List = append(pagingIEs.List, ie)
	}

	// UE Radio Capability for Paging (optional)
	if ue.UeRadioCapabilityForPaging != nil {
		ie = ngapType.PagingIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUERadioCapabilityForPaging
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PagingIEsPresentUERadioCapabilityForPaging
		ie.Value.UERadioCapabilityForPaging = new(ngapType.UERadioCapabilityForPaging)
		uERadioCapabilityForPaging := ie.Value.UERadioCapabilityForPaging
		if ue.UeRadioCapabilityForPaging.NR != "" {
			uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR.Value, err = hex.
				DecodeString(ue.UeRadioCapabilityForPaging.NR)
			if err != nil {
				logger.NgapLog.Errorf(
					"[Build Error] DecodeString ue.UeRadioCapabilityForPaging.NR error: %+v", err)
			}
		}
		if ue.UeRadioCapabilityForPaging.EUTRA != "" {
			uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA.Value, err = hex.
				DecodeString(ue.UeRadioCapabilityForPaging.EUTRA)
			if err != nil {
				logger.NgapLog.Errorf("[Build Error] DecodeString ue.UeRadioCapabilityForPaging.EUTRA error: %+v", err)
			}
		}
		pagingIEs.List = append(pagingIEs.List, ie)
	}

	// Assistance Data for Paing (optional)
	if ue.InfoOnRecommendedCellsAndRanNodesForPaging != nil {
		ie = ngapType.PagingIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAssistanceDataForPaging
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PagingIEsPresentAssistanceDataForPaging
		ie.Value.AssistanceDataForPaging = new(ngapType.AssistanceDataForPaging)

		assistanceDataForPaging := ie.Value.AssistanceDataForPaging
		assistanceDataForPaging.AssistanceDataForRecommendedCells = new(ngapType.AssistanceDataForRecommendedCells)
		recommendedCellList := &assistanceDataForPaging.
			AssistanceDataForRecommendedCells.RecommendedCellsForPaging.RecommendedCellList

		for _, recommendedCell := range ue.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedCells {
			recommendedCellItem := ngapType.RecommendedCellItem{}
			switch recommendedCell.NgRanCGI.Present {
			case context.NgRanCgiPresentNRCGI:
				recommendedCellItem.NGRANCGI.Present = ngapType.NGRANCGIPresentNRCGI
				recommendedCellItem.NGRANCGI.NRCGI = new(ngapType.NRCGI)
				nrCGI := recommendedCellItem.NGRANCGI.NRCGI
				nrCGI.PLMNIdentity = ngapConvert.PlmnIdToNgap(*recommendedCell.NgRanCGI.NRCGI.PlmnId)
				nrCGI.NRCellIdentity.Value = ngapConvert.HexToBitString(recommendedCell.NgRanCGI.NRCGI.NrCellId, 36)
			case context.NgRanCgiPresentEUTRACGI:
				recommendedCellItem.NGRANCGI.Present = ngapType.NGRANCGIPresentEUTRACGI
				recommendedCellItem.NGRANCGI.EUTRACGI = new(ngapType.EUTRACGI)
				eutraCGI := recommendedCellItem.NGRANCGI.EUTRACGI
				eutraCGI.PLMNIdentity = ngapConvert.PlmnIdToNgap(*recommendedCell.NgRanCGI.EUTRACGI.PlmnId)
				eutraCGI.EUTRACellIdentity.Value = ngapConvert.HexToBitString(recommendedCell.NgRanCGI.EUTRACGI.EutraCellId, 28)
			}

			if recommendedCell.TimeStayedInCell != nil {
				recommendedCellItem.TimeStayedInCell = recommendedCell.TimeStayedInCell
			}
			recommendedCellList.List = append(recommendedCellList.List, recommendedCellItem)
		}

		// TODO: Paging Attempt Information (optional): provided by AMF (TS 23.502 4.2.3.3, TS 38.300 9.2.5)
		pagingIEs.List = append(pagingIEs.List, ie)
	}

	// Paging Origin (optional)
	if pagingOriginNon3GPP {
		ie = ngapType.PagingIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPagingOrigin
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PagingIEsPresentPagingOrigin
		ie.Value.PagingOrigin = new(ngapType.PagingOrigin)
		ie.Value.PagingOrigin.Value = ngapType.PagingOriginPresentNon3gpp
		pagingIEs.List = append(pagingIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

// TS 23.502 4.2.2.2.3
// anType: indicate amfUe send this msg for which accessType
// amfUeNgapID: initial AMF get it from target AMF
// ngapMessage: initial UE Message to reroute
// allowedNSSAI: provided by AMF, and AMF get it from NSSF (4.2.2.2.3 step 4b)
func BuildRerouteNasRequest(ue *context.AmfUe, anType models.AccessType, amfUeNgapID *int64,
	ngapMessage []byte, allowedNSSAI *ngapType.AllowedNSSAI,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeRerouteNASRequest
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentRerouteNASRequest
	initiatingMessage.Value.RerouteNASRequest = new(ngapType.RerouteNASRequest)

	rerouteNasRequest := initiatingMessage.Value.RerouteNASRequest
	rerouteNasRequestIEs := &rerouteNasRequest.ProtocolIEs

	// RAN UE NGAP ID
	ie := ngapType.RerouteNASRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.RerouteNASRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUe[anType].RanUeNgapId

	rerouteNasRequestIEs.List = append(rerouteNasRequestIEs.List, ie)

	// AMF UE NGAP ID (optional)
	if amfUeNgapID != nil {
		ie = ngapType.RerouteNASRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.RerouteNASRequestIEsPresentAMFUENGAPID
		ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

		aMFUENGAPID := ie.Value.AMFUENGAPID
		aMFUENGAPID.Value = *amfUeNgapID

		rerouteNasRequestIEs.List = append(rerouteNasRequestIEs.List, ie)
	}

	// NGAP Message (Contains the initial ue message)
	ie = ngapType.RerouteNASRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNGAPMessage
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.RerouteNASRequestIEsPresentNGAPMessage

	msg := aper.OctetString(ngapMessage)
	ie.Value.NGAPMessage = &msg

	rerouteNasRequestIEs.List = append(rerouteNasRequestIEs.List, ie)

	// AMF Set ID
	ie = ngapType.RerouteNASRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFSetID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.RerouteNASRequestIEsPresentAMFSetID

	// <MCC><MNC><AMF Region ID><AMF Set ID><AMF Pointer><5G-TMSI>
	// <MCC><MNC> is 3 bytes, <AMF Region ID><AMF Set ID><AMF Pointer> is 3 bytes
	// 1 byte is 2 characters
	var amfID string
	if len(ue.Guti) == 19 { // MNC is 2 char
		amfID = ue.Guti[5:11]
	} else {
		amfID = ue.Guti[6:12]
	}
	_, amfSetID, _ := ngapConvert.AmfIdToNgap(amfID)

	ie.Value.AMFSetID = new(ngapType.AMFSetID)
	ie.Value.AMFSetID.Value = amfSetID

	rerouteNasRequestIEs.List = append(rerouteNasRequestIEs.List, ie)

	// Allowed NSSAI
	if allowedNSSAI != nil {
		ie = ngapType.RerouteNASRequestIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAllowedNSSAI
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.RerouteNASRequestIEsPresentAllowedNSSAI

		ie.Value.AllowedNSSAI = allowedNSSAI

		rerouteNasRequestIEs.List = append(rerouteNasRequestIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildRanConfigurationUpdateAcknowledge(
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	// criticality ->from received node when received node can't comprehend the IE or missing IE

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeRANConfigurationUpdate
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentRANConfigurationUpdateAcknowledge
	successfulOutcome.Value.RANConfigurationUpdateAcknowledge = new(ngapType.RANConfigurationUpdateAcknowledge)

	rANConfigurationUpdateAcknowledge := successfulOutcome.Value.RANConfigurationUpdateAcknowledge
	rANConfigurationUpdateAcknowledgeIEs := &rANConfigurationUpdateAcknowledge.ProtocolIEs

	// Criticality Doagnostics(Optional)
	if criticalityDiagnostics != nil {
		ie := ngapType.RANConfigurationUpdateAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.RANConfigurationUpdateAcknowledgeIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		rANConfigurationUpdateAcknowledgeIEs.List = append(rANConfigurationUpdateAcknowledgeIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildRanConfigurationUpdateFailure(
	cause ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) ([]byte, error) {
	// criticality ->from received node when received node can't comprehend the IE or missing IE
	// If the AMF cannot accept the update,
	// it shall respond with a RAN CONFIGURATION UPDATE FAILURE message and appropriate cause value.

	var pdu ngapType.NGAPPDU
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
	ie := ngapType.RANConfigurationUpdateFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.RANConfigurationUpdateFailureIEsPresentCause
	ie.Value.Cause = &cause

	rANConfigurationUpdateFailureIEs.List = append(rANConfigurationUpdateFailureIEs.List, ie)

	// Time To Wait(Optional)
	ie = ngapType.RANConfigurationUpdateFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDTimeToWait
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.RANConfigurationUpdateFailureIEsPresentTimeToWait
	ie.Value.TimeToWait = new(ngapType.TimeToWait)

	timeToWait := ie.Value.TimeToWait
	timeToWait.Value = ngapType.TimeToWaitPresentV1s

	rANConfigurationUpdateFailureIEs.List = append(rANConfigurationUpdateFailureIEs.List, ie)

	// Criticality Doagnostics(Optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.RANConfigurationUpdateFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.RANConfigurationUpdateFailureIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		rANConfigurationUpdateFailureIEs.List = append(rANConfigurationUpdateFailureIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

// An AMF shall be able to instruct other peer CP NFs, subscribed to receive such a notification,
// that it will be unavailable on this AMF and its corresponding target AMF(s).
// If CP NF does not subscribe to receive AMF unavailable notification, the CP NF may attempt
// forwarding the transaction towards the old AMF and detect that the AMF is unavailable. When
// it detects unavailable, it marks the AMF and its associated GUAMI(s) as unavailable.
// Defined in 23.501 5.21.2.2.2
func BuildAMFStatusIndication(unavailableGUAMIList ngapType.UnavailableGUAMIList) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	logger.NgapLog.Trace("Build AMF Status Indication message")

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeAMFStatusIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentAMFStatusIndication
	initiatingMessage.Value.AMFStatusIndication = new(ngapType.AMFStatusIndication)

	aMFStatusIndication := initiatingMessage.Value.AMFStatusIndication
	aMFStatusIndicationIEs := &aMFStatusIndication.ProtocolIEs

	//	Unavailable GUAMI List
	ie := ngapType.AMFStatusIndicationIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUnavailableGUAMIList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.AMFStatusIndicationIEsPresentUnavailableGUAMIList
	ie.Value.UnavailableGUAMIList = new(ngapType.UnavailableGUAMIList)

	ie.Value.UnavailableGUAMIList = &unavailableGUAMIList

	aMFStatusIndicationIEs.List = append(aMFStatusIndicationIEs.List, ie)

	return ngap.Encoder(pdu)
}

// TS 23.501 5.19.5.2
// amfOverloadResponse: the required behavior of NG-RAN, provided by AMF
// amfTrafficLoadReductionIndication(int 1~99): indicates the percentage of the type
// of traffic relative to the instantaneous incoming rate at the NG-RAN node, provided by AMF
// overloadStartNSSAIList: overload slices, provide by AMF
func BuildOverloadStart(
	amfOverloadResponse *ngapType.OverloadResponse,
	amfTrafficLoadReductionIndication int64,
	overloadStartNSSAIList *ngapType.OverloadStartNSSAIList,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeOverloadStart
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentOverloadStart
	initiatingMessage.Value.OverloadStart = new(ngapType.OverloadStart)

	overloadStart := initiatingMessage.Value.OverloadStart
	overloadStartIEs := &overloadStart.ProtocolIEs

	// AMF Overload Response (optional)
	if amfOverloadResponse != nil {
		ie := ngapType.OverloadStartIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFOverloadResponse
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.OverloadStartIEsPresentAMFOverloadResponse
		ie.Value.AMFOverloadResponse = amfOverloadResponse
		overloadStartIEs.List = append(overloadStartIEs.List, ie)
	}

	// AMF Traffic Load Reduction Indication (optional)
	if amfTrafficLoadReductionIndication != 0 {
		ie := ngapType.OverloadStartIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFTrafficLoadReductionIndication
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.OverloadStartIEsPresentAMFTrafficLoadReductionIndication
		ie.Value.AMFTrafficLoadReductionIndication = &ngapType.TrafficLoadReductionIndication{
			Value: amfTrafficLoadReductionIndication,
		}
		overloadStartIEs.List = append(overloadStartIEs.List, ie)
	}

	// Overload Start NSSAI List (optional)
	if overloadStartNSSAIList != nil {
		ie := ngapType.OverloadStartIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDOverloadStartNSSAIList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.OverloadStartIEsPresentOverloadStartNSSAIList
		ie.Value.OverloadStartNSSAIList = overloadStartNSSAIList
		overloadStartIEs.List = append(overloadStartIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildOverloadStop() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeOverloadStop
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentOverloadStop
	initiatingMessage.Value.OverloadStop = new(ngapType.OverloadStop)

	return ngap.Encoder(pdu)
}

func BuildDownlinkRanConfigurationTransfer(
	sONConfigurationTransfer *ngapType.SONConfigurationTransfer,
) ([]byte, error) {
	// sONConfigurationTransfer = sONConfigurationTransfer from uplink Ran Configuration Transfer

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeDownlinkRANConfigurationTransfer
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore
	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentDownlinkRANConfigurationTransfer
	initiatingMessage.Value.DownlinkRANConfigurationTransfer = new(ngapType.DownlinkRANConfigurationTransfer)

	downlinkRANConfigurationTransfer := initiatingMessage.Value.DownlinkRANConfigurationTransfer
	downlinkRANConfigurationTransferIEs := &downlinkRANConfigurationTransfer.ProtocolIEs

	// SON Configuration Transfer [optional]
	if sONConfigurationTransfer != nil {
		ie := ngapType.DownlinkRANConfigurationTransferIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDSONConfigurationTransferDL
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.DownlinkRANConfigurationTransferIEsPresentSONConfigurationTransferDL
		ie.Value.SONConfigurationTransferDL = new(ngapType.SONConfigurationTransfer)

		ie.Value.SONConfigurationTransferDL = sONConfigurationTransfer

		downlinkRANConfigurationTransferIEs.List = append(downlinkRANConfigurationTransferIEs.List, ie)
	}
	return ngap.Encoder(pdu)
}

func BuildDownlinkNonUEAssociatedNRPPATransport(
	ue *context.RanUe, nRPPaPDU ngapType.NRPPaPDU,
) ([]byte, error) {
	// NRPPa PDU is by pass
	// NRPPa PDU is from LMF define in 4.13.5.6

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeDownlinkNonUEAssociatedNRPPaTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentDownlinkNonUEAssociatedNRPPaTransport
	initiatingMessage.Value.DownlinkNonUEAssociatedNRPPaTransport = new(ngapType.DownlinkNonUEAssociatedNRPPaTransport)

	downlinkNonUEAssociatedNRPPaTransport := initiatingMessage.Value.DownlinkNonUEAssociatedNRPPaTransport
	downlinkNonUEAssociatedNRPPaTransportIEs := &downlinkNonUEAssociatedNRPPaTransport.ProtocolIEs

	// Routing ID
	// Routing id in the ran context
	ie := ngapType.DownlinkNonUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRoutingID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkNonUEAssociatedNRPPaTransportIEsPresentRoutingID
	ie.Value.RoutingID = new(ngapType.RoutingID)

	var err error
	routingID := ie.Value.RoutingID
	routingID.Value, err = hex.DecodeString(ue.RoutingID)
	if err != nil {
		logger.NgapLog.Errorf("[Build Error] DecodeString ue.RoutingID error: %+v", err)
	}

	downlinkNonUEAssociatedNRPPaTransportIEs.List = append(downlinkNonUEAssociatedNRPPaTransportIEs.List, ie)

	// NRPPa-PDU
	ie = ngapType.DownlinkNonUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNRPPaPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkNonUEAssociatedNRPPaTransportIEsPresentNRPPaPDU
	ie.Value.NRPPaPDU = new(ngapType.NRPPaPDU)

	ie.Value.NRPPaPDU = &nRPPaPDU

	downlinkNonUEAssociatedNRPPaTransportIEs.List = append(downlinkNonUEAssociatedNRPPaTransportIEs.List, ie)
	return ngap.Encoder(pdu)
}

func BuildTraceStart() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildDeactivateTrace(amfUe *context.AmfUe, anType models.AccessType) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	ranUe, ok := amfUe.RanUe[anType]
	if !ok {
		return nil, fmt.Errorf("ranUe for %s is nil", anType)
	}

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeDeactivateTrace
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentDeactivateTrace
	initiatingMessage.Value.DeactivateTrace = new(ngapType.DeactivateTrace)

	deactivateTrace := initiatingMessage.Value.DeactivateTrace
	deactivateTraceIEs := &deactivateTrace.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.DeactivateTraceIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DeactivateTraceIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ranUe.AmfUeNgapId

	deactivateTraceIEs.List = append(deactivateTraceIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.DeactivateTraceIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DeactivateTraceIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUe.RanUeNgapId

	deactivateTraceIEs.List = append(deactivateTraceIEs.List, ie)
	if amfUe.TraceData != nil {
		// NG-RAN TraceID
		ie = ngapType.DeactivateTraceIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNGRANTraceID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.DeactivateTraceIEsPresentNGRANTraceID
		ie.Value.NGRANTraceID = new(ngapType.NGRANTraceID)

		// TODO:composed of the following TS:32.422
		traceData := *amfUe.TraceData
		subStringSlice := strings.Split(traceData.TraceRef, "-")

		if len(subStringSlice) != 2 {
			logger.NgapLog.Warningln("TraceRef format is not correct")
		}

		plmnID := models.PlmnId{}
		plmnID.Mcc = subStringSlice[0][:3]
		plmnID.Mnc = subStringSlice[0][3:]
		traceID, err := hex.DecodeString(subStringSlice[1])
		if err != nil {
			logger.NgapLog.Errorf("[Build Error] DecodeString traceID error: %+v", err)
		}

		tmp := ngapConvert.PlmnIdToNgap(plmnID)
		traceReference := []byte{}
		traceReference = append(traceReference, tmp.Value...)
		traceReference = append(traceReference, traceID...)
		trsr := ranUe.Trsr
		trsrNgap, err := hex.DecodeString(trsr)
		if err != nil {
			logger.NgapLog.Errorf(
				"[Build Error] DecodeString trsr error: %+v", err)
		}
		ie.Value.NGRANTraceID.Value = append(ie.Value.NGRANTraceID.Value, traceReference...)
		ie.Value.NGRANTraceID.Value = append(ie.Value.NGRANTraceID.Value, trsrNgap...)
		deactivateTraceIEs.List = append(deactivateTraceIEs.List, ie)
	}
	return ngap.Encoder(pdu)
}

// AOI List is from SMF
// The SMF may subscribe to the UE mobility event notification from the AMF
// (e.g. location reporting, UE moving into or out of Area Of Interest) TS 23.502 4.3.2.2.1 Step.17
// The Location Reporting Control message shall identify the UE for which reports are requested and
// may include Reporting Type, Location Reporting Level, Area Of Interest and
// Request Reference ID TS 23.502 4.10 LocationReportingProcedure
// The AMF may request the NG-RAN location reporting with event reporting type
// (e.g. UE location or UE presence in Area of Interest),
// reporting mode and its related parameters (e.g. number of reporting) TS 23.501 5.4.7
// Location Reference ID To Be Canceled IE shall be present if
// the Event Type IE is set to "Stop UE presence in the area of interest".
func BuildLocationReportingControl(
	ue *context.RanUe,
	aoiList *ngapType.AreaOfInterestList,
	locationReportingReferenceIDToBeCancelled int64,
	eventType ngapType.EventType,
) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeLocationReportingControl
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentLocationReportingControl
	initiatingMessage.Value.LocationReportingControl = new(ngapType.LocationReportingControl)

	locationReportingControl := initiatingMessage.Value.LocationReportingControl
	locationReportingControlIEs := &locationReportingControl.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.LocationReportingControlIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.LocationReportingControlIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	locationReportingControlIEs.List = append(locationReportingControlIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.LocationReportingControlIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.LocationReportingControlIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	locationReportingControlIEs.List = append(locationReportingControlIEs.List, ie)

	// Location Reporting Request Type
	ie = ngapType.LocationReportingControlIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDLocationReportingRequestType
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.LocationReportingControlIEsPresentLocationReportingRequestType
	ie.Value.LocationReportingRequestType = new(ngapType.LocationReportingRequestType)

	locationReportingRequestType := ie.Value.LocationReportingRequestType

	// Event Type
	locationReportingRequestType.EventType = eventType

	// Report Area in Location Reporting Request Type
	locationReportingRequestType.ReportArea.Value = ngapType.ReportAreaPresentCell // only this enum

	// AOI List in Location Reporting Request Type
	if aoiList != nil {
		locationReportingRequestType.AreaOfInterestList = new(ngapType.AreaOfInterestList)
		areaOfInterestList := locationReportingRequestType.AreaOfInterestList
		areaOfInterestList.List = aoiList.List
	}

	// location reference ID to be Canceled [Conditional]
	if locationReportingRequestType.EventType.Value ==
		ngapType.EventTypePresentStopUePresenceInAreaOfInterest {
		locationReportingRequestType.LocationReportingReferenceIDToBeCancelled = new(ngapType.LocationReportingReferenceID)
		locationReportingRequestType.
			LocationReportingReferenceIDToBeCancelled.Value = locationReportingReferenceIDToBeCancelled
	}

	locationReportingControlIEs.List = append(locationReportingControlIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildUETNLABindingReleaseRequest(ue *context.RanUe) ([]byte, error) {
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUETNLABindingRelease
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUETNLABindingReleaseRequest
	initiatingMessage.Value.UETNLABindingReleaseRequest = new(ngapType.UETNLABindingReleaseRequest)

	uETNLABindingReleaseRequest := initiatingMessage.Value.UETNLABindingReleaseRequest
	uETNLABindingReleaseRequestIEs := &uETNLABindingReleaseRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UETNLABindingReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UETNLABindingReleaseRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	uETNLABindingReleaseRequestIEs.List = append(uETNLABindingReleaseRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UETNLABindingReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UETNLABindingReleaseRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uETNLABindingReleaseRequestIEs.List = append(uETNLABindingReleaseRequestIEs.List, ie)

	return ngap.Encoder(pdu)
}

// Weight Factor associated with each of the TNL association within the AMF
func BuildAMFConfigurationUpdate(tNLassociationUsage ngapType.TNLAssociationUsage,
	tNLAddressWeightFactor ngapType.TNLAddressWeightFactor,
) ([]byte, error) {
	amfSelf := context.GetSelf()
	var pdu ngapType.NGAPPDU

	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeAMFConfigurationUpdate
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject
	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentAMFConfigurationUpdate
	initiatingMessage.Value.AMFConfigurationUpdate = new(ngapType.AMFConfigurationUpdate)

	aMFConfigurationUpdate := initiatingMessage.Value.AMFConfigurationUpdate
	aMFConfigurationUpdateIEs := &aMFConfigurationUpdate.ProtocolIEs

	//	AMF Name(optional)
	ie := ngapType.AMFConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFName
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentAMFName
	ie.Value.AMFName = new(ngapType.AMFName)

	aMFName := ie.Value.AMFName
	aMFName.Value = amfSelf.Name

	aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)

	//	Served GUAMI List
	ie = ngapType.AMFConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDServedGUAMIList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentServedGUAMIList
	ie.Value.ServedGUAMIList = new(ngapType.ServedGUAMIList)

	servedGUAMIList := ie.Value.ServedGUAMIList
	for _, guami := range amfSelf.ServedGuamiList {
		servedGUAMIItem := ngapType.ServedGUAMIItem{}
		servedGUAMIItem.GUAMI.PLMNIdentity = ngapConvert.PlmnIdToNgap(util.PlmnIdNidToModelsPlmnId(*guami.PlmnId))
		regionId, setId, prtId := ngapConvert.AmfIdToNgap(guami.AmfId)
		servedGUAMIItem.GUAMI.AMFRegionID.Value = regionId
		servedGUAMIItem.GUAMI.AMFSetID.Value = setId
		servedGUAMIItem.GUAMI.AMFPointer.Value = prtId
		servedGUAMIList.List = append(servedGUAMIList.List, servedGUAMIItem)
	}

	aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)

	//	relative AMF Capability
	ie = ngapType.AMFConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRelativeAMFCapacity
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentRelativeAMFCapacity
	ie.Value.RelativeAMFCapacity = new(ngapType.RelativeAMFCapacity)
	relativeAMFCapacity := ie.Value.RelativeAMFCapacity
	relativeAMFCapacity.Value = amfSelf.RelativeCapacity

	aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)

	//	PLMN Support List
	ie = ngapType.AMFConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPLMNSupportList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentPLMNSupportList
	ie.Value.PLMNSupportList = new(ngapType.PLMNSupportList)

	pLMNSupportList := ie.Value.PLMNSupportList
	for _, plmnItem := range amfSelf.PlmnSupportList {
		pLMNSupportItem := ngapType.PLMNSupportItem{}
		pLMNSupportItem.PLMNIdentity = ngapConvert.PlmnIdToNgap(*plmnItem.PlmnId)
		for _, snssai := range plmnItem.SNssaiList {
			sliceSupportItem := ngapType.SliceSupportItem{}
			sliceSupportItem.SNSSAI = ngapConvert.SNssaiToNgap(snssai)
			pLMNSupportItem.SliceSupportList.List = append(pLMNSupportItem.SliceSupportList.List, sliceSupportItem)
		}
		pLMNSupportList.List = append(pLMNSupportList.List, pLMNSupportItem)
	}

	aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)

	//	AMF TNL Association to Add List
	ie = ngapType.AMFConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationToAddList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentAMFTNLAssociationToAddList
	ie.Value.AMFTNLAssociationToAddList = new(ngapType.AMFTNLAssociationToAddList)

	aMFTNLAssociationToAddList := ie.Value.AMFTNLAssociationToAddList

	//	AMFTNLAssociationToAddItem in AMFTNLAssociationToAddList
	aMFTNLAssociationToAddItem := ngapType.AMFTNLAssociationToAddItem{}
	aMFTNLAssociationToAddItem.AMFTNLAssociationAddress.Present = ngapType.
		CPTransportLayerInformationPresentEndpointIPAddress
	aMFTNLAssociationToAddItem.AMFTNLAssociationAddress.EndpointIPAddress = new(ngapType.TransportLayerAddress)
	*aMFTNLAssociationToAddItem.AMFTNLAssociationAddress.EndpointIPAddress = ngapConvert.
		IPAddressToNgap(amfSelf.RegisterIPv4, amfSelf.HttpIPv6Address)

	//	AMF TNL Association Usage[optional]
	if aMFTNLAssociationToAddItem.TNLAssociationUsage != nil {
		aMFTNLAssociationToAddItem.TNLAssociationUsage = new(ngapType.TNLAssociationUsage)
		aMFTNLAssociationToAddItem.TNLAssociationUsage = &tNLassociationUsage
	}

	//	AMF TNL Address Weight Factor
	aMFTNLAssociationToAddItem.TNLAddressWeightFactor = tNLAddressWeightFactor

	aMFTNLAssociationToAddList.List = append(aMFTNLAssociationToAddList.List, aMFTNLAssociationToAddItem)
	aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)

	//	AMF TNL Association to Remove List
	ie = ngapType.AMFConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationToRemoveList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentAMFTNLAssociationToRemoveList
	ie.Value.AMFTNLAssociationToRemoveList = new(ngapType.AMFTNLAssociationToRemoveList)

	aMFTNLAssociationToRemoveList := ie.Value.AMFTNLAssociationToRemoveList

	//	AMFTNLAssociationToRemoveItem
	aMFTNLAssociationToRemoveItem := ngapType.AMFTNLAssociationToRemoveItem{}
	aMFTNLAssociationToRemoveItem.AMFTNLAssociationAddress.Present = ngapType.
		CPTransportLayerInformationPresentEndpointIPAddress
	aMFTNLAssociationToRemoveItem.AMFTNLAssociationAddress.EndpointIPAddress = new(ngapType.TransportLayerAddress)
	*aMFTNLAssociationToRemoveItem.AMFTNLAssociationAddress.EndpointIPAddress = ngapConvert.
		IPAddressToNgap(amfSelf.RegisterIPv4, amfSelf.HttpIPv6Address)

	aMFTNLAssociationToRemoveList.List = append(aMFTNLAssociationToRemoveList.List, aMFTNLAssociationToRemoveItem)
	aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)

	//	AMFTNLAssociationToUpdateList
	ie = ngapType.AMFConfigurationUpdateIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationToUpdateList
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.AMFConfigurationUpdateIEsPresentAMFTNLAssociationToUpdateList
	ie.Value.AMFTNLAssociationToUpdateList = new(ngapType.AMFTNLAssociationToUpdateList)

	aMFTNLAssociationToUpdateList := ie.Value.AMFTNLAssociationToUpdateList

	//	AMFTNLAssociationAddress in AMFTNLAssociationtoUpdateItem
	aMFTNLAssociationToUpdateItem := ngapType.AMFTNLAssociationToUpdateItem{}
	aMFTNLAssociationToUpdateItem.AMFTNLAssociationAddress.Present = ngapType.
		CPTransportLayerInformationPresentEndpointIPAddress
	aMFTNLAssociationToUpdateItem.AMFTNLAssociationAddress.EndpointIPAddress = new(ngapType.TransportLayerAddress)
	*aMFTNLAssociationToUpdateItem.AMFTNLAssociationAddress.EndpointIPAddress = ngapConvert.
		IPAddressToNgap(amfSelf.RegisterIPv4, amfSelf.HttpIPv6Address)

	//	TNLAssociationUsage in AMFTNLAssociationtoUpdateItem [optional]
	if aMFTNLAssociationToUpdateItem.TNLAssociationUsage != nil {
		aMFTNLAssociationToUpdateItem.TNLAssociationUsage = new(ngapType.TNLAssociationUsage)
		aMFTNLAssociationToUpdateItem.TNLAssociationUsage = &tNLassociationUsage
	}
	//	TNLAddressWeightFactor in AMFTNLAssociationtoUpdateItem [optional]
	if aMFTNLAssociationToUpdateItem.TNLAddressWeightFactor != nil {
		aMFTNLAssociationToUpdateItem.TNLAddressWeightFactor = new(ngapType.TNLAddressWeightFactor)
		aMFTNLAssociationToUpdateItem.TNLAddressWeightFactor = &tNLAddressWeightFactor
	}
	aMFTNLAssociationToUpdateList.List = append(aMFTNLAssociationToUpdateList.List, aMFTNLAssociationToUpdateItem)
	aMFConfigurationUpdateIEs.List = append(aMFConfigurationUpdateIEs.List, ie)

	return ngap.Encoder(pdu)
}

// NRPPa PDU is a pdu from LMF to RAN defined in TS 23.502 4.13.5.5 step 3
// NRPPa PDU is by pass
func BuildDownlinkUEAssociatedNRPPaTransport(ue *context.RanUe, nRPPaPDU ngapType.NRPPaPDU) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeDownlinkUEAssociatedNRPPaTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentDownlinkUEAssociatedNRPPaTransport
	initiatingMessage.Value.DownlinkUEAssociatedNRPPaTransport = new(ngapType.DownlinkUEAssociatedNRPPaTransport)

	downlinkUEAssociatedNRPPaTransport := initiatingMessage.Value.DownlinkUEAssociatedNRPPaTransport
	downlinkUEAssociatedNRPPaTransportIEs := &downlinkUEAssociatedNRPPaTransport.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.DownlinkUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkUEAssociatedNRPPaTransportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	downlinkUEAssociatedNRPPaTransportIEs.List = append(downlinkUEAssociatedNRPPaTransportIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.DownlinkUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkUEAssociatedNRPPaTransportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	downlinkUEAssociatedNRPPaTransportIEs.List = append(downlinkUEAssociatedNRPPaTransportIEs.List, ie)

	// Routing ID
	ie = ngapType.DownlinkUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRoutingID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkUEAssociatedNRPPaTransportIEsPresentRoutingID
	ie.Value.RoutingID = new(ngapType.RoutingID)

	var err error
	routingID := ie.Value.RoutingID
	routingID.Value, err = hex.DecodeString(ue.RoutingID)
	if err != nil {
		logger.NgapLog.Errorf("[Build Error] DecodeString ue.RoutingID error: %+v", err)
	}

	downlinkUEAssociatedNRPPaTransportIEs.List = append(downlinkUEAssociatedNRPPaTransportIEs.List, ie)

	// NRPPa-PDU
	ie = ngapType.DownlinkUEAssociatedNRPPaTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNRPPaPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.DownlinkUEAssociatedNRPPaTransportIEsPresentNRPPaPDU
	ie.Value.NRPPaPDU = new(ngapType.NRPPaPDU)

	ie.Value.NRPPaPDU = &nRPPaPDU

	downlinkUEAssociatedNRPPaTransportIEs.List = append(downlinkUEAssociatedNRPPaTransportIEs.List, ie)

	return ngap.Encoder(pdu)
}
