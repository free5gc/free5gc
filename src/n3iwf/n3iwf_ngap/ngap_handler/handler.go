package ngap_handler

import (
	"encoding/binary"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/timer"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_handler/n3iwf_message"
	"free5gc/src/n3iwf/n3iwf_ike/ike_handler"
	"free5gc/src/n3iwf/n3iwf_ike/ike_message"
	"free5gc/src/n3iwf/n3iwf_ike/udp_server"
	"free5gc/src/n3iwf/n3iwf_ngap/ngap_message"
	"math/rand"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

var ngapLog *logrus.Entry

func init() {
	ngapLog = logger.NgapLog
}

func HandleEventSCTPConnect(sctpAddr string) {
	ngapLog.Infoln("[N3IWF] Handle SCTP connect event")
	ngap_message.SendNGSetupRequest(sctpAddr)
}

func HandleNGSetupResponse(sctpAddr string, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle NG Setup Response")

	var amfName *ngapType.AMFName
	var servedGUAMIList *ngapType.ServedGUAMIList
	var relativeAMFCapacity *ngapType.RelativeAMFCapacity
	var plmnSupportList *ngapType.PLMNSupportList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	n3iwfSelf := n3iwf_context.N3IWFSelf()

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ngapLog.Error("Successful Outcome is nil")
		return
	}

	ngSetupResponse := successfulOutcome.Value.NGSetupResponse
	if ngSetupResponse == nil {
		ngapLog.Error("ngSetupResponse is nil")
		return
	}

	for _, ie := range ngSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFName:
			ngapLog.Traceln("[NGAP] Decode IE AMFName")
			amfName = ie.Value.AMFName
			if amfName == nil {
				ngapLog.Errorf("AMFName is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDServedGUAMIList:
			ngapLog.Traceln("[NGAP] Decode IE ServedGUAMIList")
			servedGUAMIList = ie.Value.ServedGUAMIList
			if servedGUAMIList == nil {
				ngapLog.Errorf("ServedGUAMIList is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRelativeAMFCapacity:
			ngapLog.Traceln("[NGAP] Decode IE RelativeAMFCapacity")
			relativeAMFCapacity = ie.Value.RelativeAMFCapacity
		case ngapType.ProtocolIEIDPLMNSupportList:
			ngapLog.Traceln("[NGAP] Decode IE PLMNSupportList")
			plmnSupportList = ie.Value.PLMNSupportList
			if plmnSupportList == nil {
				ngapLog.Errorf("PLMNSupportList is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			ngapLog.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if len(iesCriticalityDiagnostics.List) != 0 {
		ngapLog.Traceln("[NGAP] Sending error indication to AMF, because some mandatory IEs were not included")

		cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)

		procedureCode := ngapType.ProcedureCodeNGSetup
		triggeringMessage := ngapType.TriggeringMessagePresentSuccessfulOutcome
		procedureCriticality := ngapType.CriticalityPresentReject

		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)

		ngap_message.SendErrorIndicationWithSctpAddr(sctpAddr, nil, nil, cause, &criticalityDiagnostics)

		return
	}

	amfInfo := n3iwfSelf.NewN3iwfAmf(sctpAddr)

	if amfName != nil {
		amfInfo.AMFName = amfName
	}

	if servedGUAMIList != nil {
		amfInfo.ServedGUAMIList = servedGUAMIList
	}

	if relativeAMFCapacity != nil {
		amfInfo.RelativeAMFCapacity = relativeAMFCapacity
	}

	if plmnSupportList != nil {
		amfInfo.PLMNSupportList = plmnSupportList
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

}

func HandleNGSetupFailure(sctpAddr string, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle NG Setup Failure")

	var cause *ngapType.Cause
	var timeToWait *ngapType.TimeToWait
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		ngapLog.Error("Unseccessful Message is nil")
		return
	}

	ngSetupFailure := unsuccessfulOutcome.Value.NGSetupFailure
	if ngSetupFailure == nil {
		ngapLog.Error("NGSetupFailure is nil")
		return
	}

	for _, ie := range ngSetupFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			ngapLog.Traceln("[NGAP] Decode IE Cause")
			cause = ie.Value.Cause
			if cause == nil {
				ngapLog.Error("Cause is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDTimeToWait:
			ngapLog.Traceln("[NGAP] Decode IE TimeToWait")
			timeToWait = ie.Value.TimeToWait
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			ngapLog.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		// TODO: Send error indication
		ngapLog.Traceln("[NGAP] Sending error indication to AMF, because some mandatory IEs were not included")

		cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)

		procedureCode := ngapType.ProcedureCodeNGSetup
		triggeringMessage := ngapType.TriggeringMessagePresentUnsuccessfullOutcome
		procedureCriticality := ngapType.CriticalityPresentReject

		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)

		ngap_message.SendErrorIndicationWithSctpAddr(sctpAddr, nil, nil, cause, &criticalityDiagnostics)

		return
	}

	if cause != nil {
		printAndGetCause(cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

	var waitingTime int

	if timeToWait != nil {

		switch timeToWait.Value {
		case ngapType.TimeToWaitPresentV1s:
			waitingTime = 1
		case ngapType.TimeToWaitPresentV2s:
			waitingTime = 2
		case ngapType.TimeToWaitPresentV5s:
			waitingTime = 5
		case ngapType.TimeToWaitPresentV10s:
			waitingTime = 10
		case ngapType.TimeToWaitPresentV20s:
			waitingTime = 20
		case ngapType.TimeToWaitPresentV60s:
			waitingTime = 60
		}

	}

	if waitingTime != 0 {
		ngapLog.Infof("Wait at lease  %ds to reinitialize with same AMF[%s]", waitingTime, sctpAddr)
		n3iwf_context.N3IWFSelf().AMFReInitAvailableList[sctpAddr] = false
		msg := n3iwf_message.HandlerMessage{
			Event:    n3iwf_message.EventSCTPConnectMessage,
			SCTPAddr: sctpAddr,
		}
		timer.StartTimer(waitingTime, func(msg interface{}) {
			handlerMessage := msg.(n3iwf_message.HandlerMessage)
			n3iwf_context.N3IWFSelf().AMFReInitAvailableList[handlerMessage.SCTPAddr] = true
			n3iwf_message.SendMessage(handlerMessage)
		}, msg)
		return
	}

	// TODO: Limited retry mechanism
	ngap_message.SendNGSetupRequest(sctpAddr)
}

func HandleNGReset(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {

	ngapLog.Infoln("[N3IWF] Handle NG Reset")

	var cause *ngapType.Cause
	var resetType *ngapType.ResetType

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("InitiatingMessage is nil")
		return
	}

	nGReset := initiatingMessage.Value.NGReset
	if nGReset == nil {
		ngapLog.Error("nGReset is nil")
		return
	}

	for _, ie := range nGReset.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			ngapLog.Traceln("[NGAP] Decode IE Cause")
			cause = ie.Value.Cause
		case ngapType.ProtocolIEIDResetType:
			ngapLog.Traceln("[NGAP] Decode IE ResetType")
			resetType = ie.Value.ResetType
			if resetType == nil {
				ngapLog.Error("ResetType is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodeNGReset
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(amf, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	printAndGetCause(cause)

	switch resetType.Present {
	case ngapType.ResetTypePresentNGInterface:
		ngapLog.Trace("ResetType Present: NG Interface")
		// TODO: Release Uu Interface related to this amf(IPSec)
		// Remove all Ue
		amf.RemoveAllRelatedUe()
		ngap_message.SendNGResetAcknowledge(amf, nil, nil)
	case ngapType.ResetTypePresentPartOfNGInterface:
		ngapLog.Trace("ResetType Present: Part of NG Interface")

		partOfNGInterface := resetType.PartOfNGInterface
		if partOfNGInterface == nil {
			ngapLog.Error("PartOfNGInterface is nil")
			return
		}

		var ue *n3iwf_context.N3IWFUe

		for _, ueAssociatedLogicalNGConnectionItem := range partOfNGInterface.List {
			if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
				ngapLog.Tracef("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				ue = n3iwf_context.N3IWFSelf().FindUeByRanUeNgapID(ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
			} else if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
				ngapLog.Tracef("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				ue = amf.FindUeByAmfUeNgapID(ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
			}

			if ue == nil {
				ngapLog.Warn("Cannot not find UE Context")
				if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
					ngapLog.Warnf("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				}
				if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
					ngapLog.Warnf("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				}
				continue
			}
			// TODO: Release Uu Interface (IPSec)
			ue.Remove()
		}
		ngap_message.SendNGResetAcknowledge(amf, partOfNGInterface, nil)
	default:
		ngapLog.Warnf("Invalid ResetType[%d]", resetType.Present)
	}
}

func HandleNGResetAcknowledge(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {

	ngapLog.Infoln("[N3IWF] Handle NG Reset Acknowledge")

	var uEAssociatedLogicalNGConnectionList *ngapType.UEAssociatedLogicalNGConnectionList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ngapLog.Error("SuccessfulOutcome is nil")
		return
	}

	nGResetAcknowledge := successfulOutcome.Value.NGResetAcknowledge
	if nGResetAcknowledge == nil {
		ngapLog.Error("nGResetAcknowledge is nil")
		return
	}

	for _, ie := range nGResetAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList:
			ngapLog.Traceln("[NGAP] Decode IE UEAssociatedLogicalNGConnectionList")
			uEAssociatedLogicalNGConnectionList = ie.Value.UEAssociatedLogicalNGConnectionList
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			ngapLog.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if uEAssociatedLogicalNGConnectionList != nil {
		ngapLog.Tracef("%d UE association(s) has been reset", len(uEAssociatedLogicalNGConnectionList.List))
		for i, item := range uEAssociatedLogicalNGConnectionList.List {
			if item.AMFUENGAPID != nil && item.RANUENGAPID != nil {
				ngapLog.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[%d]", i+1, item.AMFUENGAPID.Value, item.RANUENGAPID.Value)
			} else if item.AMFUENGAPID != nil {
				ngapLog.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[unknown]", i+1, item.AMFUENGAPID.Value)
			} else if item.RANUENGAPID != nil {
				ngapLog.Tracef("%d: AmfUeNgapID[unknown] RanUeNgapID[%d]", i+1, item.RANUENGAPID.Value)
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandleInitialContextSetupRequest(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Initial Context Setup Request")

	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	var oldAMF *ngapType.AMFName
	var ueAggregateMaximumBitRate *ngapType.UEAggregateMaximumBitRate
	var coreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation
	var guami *ngapType.GUAMI
	var pduSessionResourceSetupListCxtReq *ngapType.PDUSessionResourceSetupListCxtReq
	var allowedNSSAI *ngapType.AllowedNSSAI
	var ueSecurityCapabilities *ngapType.UESecurityCapabilities
	var securityKey *ngapType.SecurityKey
	var traceActivation *ngapType.TraceActivation
	var ueRadioCapability *ngapType.UERadioCapability
	var indexToRFSP *ngapType.IndexToRFSP
	var maskedIMEISV *ngapType.MaskedIMEISV
	var nasPDU *ngapType.NASPDU
	var emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var n3iwfUe *n3iwf_context.N3IWFUe
	var n3iwfSelf = n3iwf_context.N3IWFSelf()

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("Initiating Message is nil")
		return
	}

	initialContextSetupRequest := initiatingMessage.Value.InitialContextSetupRequest
	if initialContextSetupRequest == nil {
		ngapLog.Error("InitialContextSetupRequest is nil")
		return
	}

	for _, ie := range initialContextSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			amfUeNgapID = ie.Value.AMFUENGAPID
			if amfUeNgapID == nil {
				ngapLog.Errorf("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Errorf("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDOldAMF:
			ngapLog.Traceln("[NGAP] Decode IE OldAMF")
			oldAMF = ie.Value.OldAMF
		case ngapType.ProtocolIEIDUEAggregateMaximumBitRate:
			ngapLog.Traceln("[NGAP] Decode IE UEAggregateMaximumBitRate")
			ueAggregateMaximumBitRate = ie.Value.UEAggregateMaximumBitRate
		case ngapType.ProtocolIEIDCoreNetworkAssistanceInformation:
			ngapLog.Traceln("[NGAP] Decode IE CoreNetworkAssistanceInformation")
			coreNetworkAssistanceInformation = ie.Value.CoreNetworkAssistanceInformation
			if coreNetworkAssistanceInformation != nil {
				ngapLog.Warnln("Not Supported IE [CoreNetworkAssistanceInformation]")
			}
		case ngapType.ProtocolIEIDGUAMI:
			ngapLog.Traceln("[NGAP] Decode IE GUAMI")
			guami = ie.Value.GUAMI
			if guami == nil {
				ngapLog.Errorf("GUAMI is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtReq:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionResourceSetupListCxtReq")
			pduSessionResourceSetupListCxtReq = ie.Value.PDUSessionResourceSetupListCxtReq
		case ngapType.ProtocolIEIDAllowedNSSAI:
			ngapLog.Traceln("[NGAP] Decode IE AllowedNSSAI")
			allowedNSSAI = ie.Value.AllowedNSSAI
			if allowedNSSAI == nil {
				ngapLog.Errorf("AllowedNSSAI is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDUESecurityCapabilities:
			ngapLog.Traceln("[NGAP] Decode IE UESecurityCapabilities")
			ueSecurityCapabilities = ie.Value.UESecurityCapabilities
			if ueSecurityCapabilities == nil {
				ngapLog.Errorf("UESecurityCapabilities is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDSecurityKey:
			ngapLog.Traceln("[NGAP] Decode IE SecurityKey")
			securityKey = ie.Value.SecurityKey
			if securityKey == nil {
				ngapLog.Errorf("SecurityKey is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDTraceActivation:
			ngapLog.Traceln("[NGAP] Decode IE TraceActivation")
			traceActivation = ie.Value.TraceActivation
			if traceActivation != nil {
				ngapLog.Warnln("Not Supported IE [TraceActivation]")
			}
		case ngapType.ProtocolIEIDUERadioCapability:
			ngapLog.Traceln("[NGAP] Decode IE UERadioCapability")
			ueRadioCapability = ie.Value.UERadioCapability
		case ngapType.ProtocolIEIDIndexToRFSP:
			ngapLog.Traceln("[NGAP] Decode IE IndexToRFSP")
			indexToRFSP = ie.Value.IndexToRFSP
		case ngapType.ProtocolIEIDMaskedIMEISV:
			ngapLog.Traceln("[NGAP] Decode IE MaskedIMEISV")
			maskedIMEISV = ie.Value.MaskedIMEISV
		case ngapType.ProtocolIEIDNASPDU:
			ngapLog.Traceln("[NGAP] Decode IE NAS PDU")
			nasPDU = ie.Value.NASPDU
		case ngapType.ProtocolIEIDEmergencyFallbackIndicator:
			ngapLog.Traceln("[NGAP] Decode IE EmergencyFallbackIndicator")
			emergencyFallbackIndicator = ie.Value.EmergencyFallbackIndicator
			if emergencyFallbackIndicator != nil {
				ngapLog.Warnln("Not Supported IE [EmergencyFallbackIndicator]")
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		ngapLog.Traceln("[NGAP] Sending unsuccessful outcome to AMF, because some mandatory IEs were not included")
		cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage)

		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)

		failedListCxtFail := new(ngapType.PDUSessionResourceFailedToSetupListCxtFail)
		for _, item := range pduSessionResourceSetupListCxtReq.List {
			transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
			if err != nil {
				ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", err)
			}
			ngap_message.AppendPDUSessionResourceFailedToSetupListCxtfail(failedListCxtFail, item.PDUSessionID.Value, transfer)
		}

		ngap_message.SendInitialContextSetupFailure(amf, n3iwfUe, *cause, failedListCxtFail, &criticalityDiagnostics)
		return
	}

	if (amfUeNgapID != nil) && (ranUeNgapID != nil) {
		// Find UE context
		n3iwfUe = n3iwfSelf.FindUeByRanUeNgapID(ranUeNgapID.Value)
		if n3iwfUe == nil {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
			// TODO: build cause and handle error
			// Cause: Unknown local UE NGAP ID
			return
		} else {
			if n3iwfUe.AmfUeNgapId != amfUeNgapID.Value {
				// TODO: build cause and handle error
				// Cause: Inconsistent remote UE NGAP ID
				return
			}
		}
	}

	n3iwfUe.AmfUeNgapId = amfUeNgapID.Value
	n3iwfUe.RanUeNgapId = ranUeNgapID.Value

	if pduSessionResourceSetupListCxtReq != nil {
		if ueAggregateMaximumBitRate != nil {
			n3iwfUe.Ambr = ueAggregateMaximumBitRate
		} else {
			ngapLog.Errorln("IE[UEAggregateMaximumBitRate] is nil")
			cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage)

			criticalityDiagnosticsIEItem := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDUEAggregateMaximumBitRate, ngapType.TypeOfErrorPresentMissing)
			iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, criticalityDiagnosticsIEItem)
			criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)

			failedListCxtFail := new(ngapType.PDUSessionResourceFailedToSetupListCxtFail)
			for _, item := range pduSessionResourceSetupListCxtReq.List {
				transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
				if err != nil {
					ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", err)
				}
				ngap_message.AppendPDUSessionResourceFailedToSetupListCxtfail(failedListCxtFail, item.PDUSessionID.Value, transfer)
			}

			ngap_message.SendInitialContextSetupFailure(amf, n3iwfUe, *cause, failedListCxtFail, &criticalityDiagnostics)
			return
		}

		setupListCxtRes := new(ngapType.PDUSessionResourceSetupListCxtRes)
		failedListCxtRes := new(ngapType.PDUSessionResourceFailedToSetupListCxtRes)
		// UE temporary data for PDU session setup response
		n3iwfUe.TemporaryPDUSessionSetupData = &n3iwf_context.PDUSessionSetupTemporaryData{
			SetupListCxtRes:  setupListCxtRes,
			FailedListCxtRes: failedListCxtRes,
		}
		n3iwfUe.TemporaryPDUSessionSetupData.NGAPProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup

		for _, item := range pduSessionResourceSetupListCxtReq.List {
			pduSessionID := item.PDUSessionID.Value
			// TODO: send NAS to UE
			// pduSessionNasPdu := item.NASPDU
			snssai := item.SNSSAI

			transfer := ngapType.PDUSessionResourceSetupRequestTransfer{}
			err := aper.UnmarshalWithParams(item.PDUSessionResourceSetupRequestTransfer, &transfer, "valueExt")
			if err != nil {
				ngapLog.Errorf("[PDUSessionID: %d] PDUSessionResourceSetupRequestTransfer Decode Error: %+v\n", pduSessionID, err)
			}

			pduSession, err := n3iwfUe.CreatePDUSession(pduSessionID, snssai)
			if err != nil {
				ngapLog.Errorf("Create PDU Session Error: %+v\n", err)

				cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentMultiplePDUSessionIDInstances)
				unsuccessfulTransfer, buildErr := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
				if buildErr != nil {
					ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", buildErr)
				}
				ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(failedListCxtRes, pduSessionID, unsuccessfulTransfer)
				continue
			}

			success, resTransfer := handlePDUSessionResourceSetupRequestTransfer(n3iwfUe, pduSession, transfer)
			if success {
				// Append this PDU session to unactivated PDU session list
				n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession = append(n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession, pduSessionID)
			} else {
				// Delete the pdusession store in UE conext
				delete(n3iwfUe.PduSessionList, pduSessionID)
				ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(failedListCxtRes, pduSessionID, resTransfer)
			}
		}
	}

	if oldAMF != nil {
		ngapLog.Debugf("Old AMF: %s\n", oldAMF.Value)
	}

	if guami != nil {
		n3iwfUe.Guami = guami
	}

	if allowedNSSAI != nil {
		n3iwfUe.AllowedNssai = allowedNSSAI
	}

	if maskedIMEISV != nil {
		n3iwfUe.MaskedIMEISV = maskedIMEISV
	}

	if ueRadioCapability != nil {
		n3iwfUe.RadioCapability = ueRadioCapability
	}

	if coreNetworkAssistanceInformation != nil {
		n3iwfUe.CoreNetworkAssistanceInformation = coreNetworkAssistanceInformation
	}

	if indexToRFSP != nil {
		n3iwfUe.IndexToRfsp = indexToRFSP.Value
	}

	if ueSecurityCapabilities != nil {
		n3iwfUe.SecurityCapabilities = ueSecurityCapabilities
	}

	if securityKey != nil {
		n3iwfUe.Kn3iwf = securityKey.Value.Bytes
	}

	if nasPDU != nil {
		// TODO: Send NAS UE
	}

	// Send EAP Success to UE
	ikeSecurityAssociation := n3iwfUe.N3IWFIKESecurityAssociation

	// IKEHDR-SK-{response}
	var eap *ike_message.EAP

	// Build IKE message
	responseIKEMessage := ike_message.BuildIKEHeader(ikeSecurityAssociation.RemoteSPI, ikeSecurityAssociation.LocalSPI, ike_message.IKE_AUTH, ike_message.ResponseBitCheck, ikeSecurityAssociation.MessageID)

	// Build response
	var ikePayload []ike_message.IKEPayloadType

	// EAP Success
	var identifier uint8
	for {
		identifier = uint8(rand.Uint32())
		if identifier != ikeSecurityAssociation.LastEAPIdentifier {
			ikeSecurityAssociation.LastEAPIdentifier = identifier
			break
		}
	}
	eap = ike_message.BuildEAPSuccess(identifier)
	ikePayload = append(ikePayload, eap)

	if err := ike_handler.EncryptProcedure(ikeSecurityAssociation, ikePayload, responseIKEMessage); err != nil {
		ngapLog.Errorf("Encrypting IKE message failed: %+v", err)
		return
	}

	// Send IKE message to UE
	ike_handler.SendIKEMessageToUE(n3iwfUe.UDPSendInfoGroup, responseIKEMessage)

	n3iwfUe.UDPSendInfoGroup = nil
	n3iwfUe.N3IWFIKESecurityAssociation.State++
}

// handlePDUSessionResourceSetupRequestTransfer parse and store needed information
// for setup user plane connection for UE
// When success, it will return a status "success" set to "true" for indication and set
// unsuccessfulTransfer to nil, otherwise, return false and corresponding transfer
func handlePDUSessionResourceSetupRequestTransfer(ue *n3iwf_context.N3IWFUe, pduSession *n3iwf_context.PDUSession, transfer ngapType.PDUSessionResourceSetupRequestTransfer) (success bool, unsuccessfulTransfer []byte) {

	var pduSessionAMBR *ngapType.PDUSessionAggregateMaximumBitRate
	var ulNGUUPTNLInformation *ngapType.UPTransportLayerInformation
	var pduSessionType *ngapType.PDUSessionType
	var securityIndication *ngapType.SecurityIndication
	var networkInstance *ngapType.NetworkInstance
	var qosFlowSetupRequestList *ngapType.QosFlowSetupRequestList
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	for _, ie := range transfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDPDUSessionAggregateMaximumBitRate:
			pduSessionAMBR = ie.Value.PDUSessionAggregateMaximumBitRate
		case ngapType.ProtocolIEIDULNGUUPTNLInformation:
			ulNGUUPTNLInformation = ie.Value.ULNGUUPTNLInformation
			if ulNGUUPTNLInformation == nil {
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionType:
			pduSessionType = ie.Value.PDUSessionType
			if pduSessionType == nil {
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDSecurityIndication:
			securityIndication = ie.Value.SecurityIndication
		case ngapType.ProtocolIEIDNetworkInstance:
			networkInstance = ie.Value.NetworkInstance
		case ngapType.ProtocolIEIDQosFlowSetupRequestList:
			qosFlowSetupRequestList = ie.Value.QosFlowSetupRequestList
			if qosFlowSetupRequestList == nil {
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		success = false
		cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage)
		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)
		responseTransfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, &criticalityDiagnostics)
		if err != nil {
			ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", err)
		}
		unsuccessfulTransfer = responseTransfer
		return
	}

	pduSession.Ambr = pduSessionAMBR
	pduSession.Type = pduSessionType
	pduSession.NetworkInstance = networkInstance

	// Security Indication
	if securityIndication != nil {
		switch securityIndication.IntegrityProtectionIndication.Value {
		case ngapType.IntegrityProtectionIndicationPresentNotNeeded:
			pduSession.SecurityIntegrity = false
		case ngapType.IntegrityProtectionIndicationPresentPreferred:
			pduSession.SecurityIntegrity = true
		case ngapType.IntegrityProtectionIndicationPresentRequired:
			pduSession.SecurityIntegrity = true
		default:
			ngapLog.Error("Unknown security integrity indication")
			success = false
			cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentSemanticError)
			responseTransfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
			if err != nil {
				ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", err)
			}
			unsuccessfulTransfer = responseTransfer
			return
		}

		switch securityIndication.ConfidentialityProtectionIndication.Value {
		case ngapType.ConfidentialityProtectionIndicationPresentNotNeeded:
			pduSession.SecurityCipher = false
		case ngapType.ConfidentialityProtectionIndicationPresentPreferred:
			pduSession.SecurityCipher = true
		case ngapType.ConfidentialityProtectionIndicationPresentRequired:
			pduSession.SecurityCipher = true
		default:
			ngapLog.Error("Unknown security confidentiality indication")
			success = false
			cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentSemanticError)
			responseTransfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
			if err != nil {
				ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", err)
			}
			unsuccessfulTransfer = responseTransfer
			return
		}
	} else {
		pduSession.SecurityIntegrity = false
		pduSession.SecurityCipher = true
	}

	// TODO: apply qos rule
	for _, item := range qosFlowSetupRequestList.List {
		// QoS Flow
		qosFlow := new(n3iwf_context.QosFlow)
		qosFlow.Identifier = item.QosFlowIdentifier.Value
		qosFlow.Parameters = item.QosFlowLevelQosParameters
		pduSession.QosFlows[item.QosFlowIdentifier.Value] = qosFlow
		// QFI List
		pduSession.QFIList = append(pduSession.QFIList, uint8(item.QosFlowIdentifier.Value))
	}

	// Setup GTP tunnel with UPF
	// TODO: Support IPv6
	upfIPv4, _ := ngapConvert.IPAddressToString(ulNGUUPTNLInformation.GTPTunnel.TransportLayerAddress)
	if upfIPv4 != "" {
		gtpConnection := &n3iwf_context.GTPConnectionInfo{
			UPFIPAddr:    upfIPv4,
			OutgoingTEID: binary.BigEndian.Uint32(ulNGUUPTNLInformation.GTPTunnel.GTPTEID.Value),
		}
		pduSession.GTPConnection = gtpConnection
	}

	success = true
	unsuccessfulTransfer = nil
	return
}

func HandleUEContextModificationRequest(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle UE Context Modification Request")

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}

	var amfUeNgapID *ngapType.AMFUENGAPID
	var newAmfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	var ueAggregateMaximumBitRate *ngapType.UEAggregateMaximumBitRate
	var ueSecurityCapabilities *ngapType.UESecurityCapabilities
	var securityKey *ngapType.SecurityKey
	var indexToRFSP *ngapType.IndexToRFSP
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var n3iwfUe *n3iwf_context.N3IWFUe
	var n3iwfSelf = n3iwf_context.N3IWFSelf()

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("Initiating Message is nil")
		return
	}

	ueContextModificationRequest := initiatingMessage.Value.UEContextModificationRequest
	if ueContextModificationRequest == nil {
		ngapLog.Error("UEContextModificationRequest is nil")
		return
	}

	for _, ie := range ueContextModificationRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			amfUeNgapID = ie.Value.AMFUENGAPID
			if amfUeNgapID == nil {
				ngapLog.Errorf("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Errorf("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDSecurityKey:
			ngapLog.Traceln("[NGAP] Decode IE SecurityKey")
			securityKey = ie.Value.SecurityKey
		case ngapType.ProtocolIEIDIndexToRFSP:
			ngapLog.Traceln("[NGAP] Decode IE IndexToRFSP")
			indexToRFSP = ie.Value.IndexToRFSP
		case ngapType.ProtocolIEIDUEAggregateMaximumBitRate:
			ngapLog.Traceln("[NGAP] Decode IE UEAggregateMaximumBitRate")
			ueAggregateMaximumBitRate = ie.Value.UEAggregateMaximumBitRate
		case ngapType.ProtocolIEIDUESecurityCapabilities:
			ngapLog.Traceln("[NGAP] Decode IE UESecurityCapabilities")
			ueSecurityCapabilities = ie.Value.UESecurityCapabilities
		case ngapType.ProtocolIEIDCoreNetworkAssistanceInformation:
			ngapLog.Traceln("[NGAP] Decode IE CoreNetworkAssistanceInformation")
			ngapLog.Warnln("Not Supported IE [CoreNetworkAssistanceInformation]")
		case ngapType.ProtocolIEIDEmergencyFallbackIndicator:
			ngapLog.Traceln("[NGAP] Decode IE EmergencyFallbackIndicator")
			ngapLog.Warnln("Not Supported IE [EmergencyFallbackIndicator]")
		case ngapType.ProtocolIEIDNewAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE NewAMFUENGAPID")
			newAmfUeNgapID = ie.Value.NewAMFUENGAPID
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		// TODO: send unsuccessful outcome or error indication
		return
	}

	if (amfUeNgapID != nil) && (ranUeNgapID != nil) {
		// Find UE context
		n3iwfUe = n3iwfSelf.FindUeByRanUeNgapID(ranUeNgapID.Value)
		if n3iwfUe == nil {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
			// TODO: build cause and handle error
			// Cause: Unknown local UE NGAP ID
			return
		} else {
			if n3iwfUe.AmfUeNgapId != amfUeNgapID.Value {
				// TODO: build cause and handle error
				// Cause: Inconsistent remote UE NGAP ID
				return
			}
		}
	}

	if newAmfUeNgapID != nil {
		ngapLog.Debugf("New AmfUeNgapID[%d]\n", newAmfUeNgapID.Value)
		n3iwfUe.AmfUeNgapId = newAmfUeNgapID.Value
	}

	if ueAggregateMaximumBitRate != nil {
		n3iwfUe.Ambr = ueAggregateMaximumBitRate
		// TODO: use the received UE Aggregate Maximum Bit Rate for all non-GBR QoS flows
	}

	if ueSecurityCapabilities != nil {
		n3iwfUe.SecurityCapabilities = ueSecurityCapabilities
	}

	if securityKey != nil {
		n3iwfUe.Kn3iwf = securityKey.Value.Bytes
	}

	// TODO: use new security key to update security context

	if indexToRFSP != nil {
		n3iwfUe.IndexToRfsp = indexToRFSP.Value
	}

	ngap_message.SendUEContextModificationResponse(amf, n3iwfUe, nil)
}

func HandleUEContextReleaseCommand(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle UE Context Release Command")

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}

	var ueNgapIDs *ngapType.UENGAPIDs
	var cause *ngapType.Cause
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var n3iwfUe *n3iwf_context.N3IWFUe
	var n3iwfSelf = n3iwf_context.N3IWFSelf()

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("Initiating Message is nil")
		return
	}

	ueContextReleaseCommand := initiatingMessage.Value.UEContextReleaseCommand
	if ueContextReleaseCommand == nil {
		ngapLog.Error("UEContextReleaseCommand is nil")
		return
	}

	for _, ie := range ueContextReleaseCommand.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUENGAPIDs:
			ngapLog.Traceln("[NGAP] Decode IE UENGAPIDs")
			ueNgapIDs = ie.Value.UENGAPIDs
			if ueNgapIDs == nil {
				ngapLog.Errorf("UENGAPIDs is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDCause:
			ngapLog.Traceln("[NGAP] Decode IE Cause")
			cause = ie.Value.Cause
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		// TODO: send error indication
		return
	}

	switch ueNgapIDs.Present {
	case ngapType.UENGAPIDsPresentUENGAPIDPair:
		n3iwfUe = n3iwfSelf.FindUeByRanUeNgapID(ueNgapIDs.UENGAPIDPair.RANUENGAPID.Value)
		if n3iwfUe == nil {
			n3iwfUe = amf.FindUeByAmfUeNgapID(ueNgapIDs.UENGAPIDPair.AMFUENGAPID.Value)
		}
	case ngapType.UENGAPIDsPresentAMFUENGAPID:
		// TODO: find UE according to specific AMF
		// The implementation here may have error when N3IWF need to
		// connect multiple AMFs.
		// Use UEpool in AMF context can solve this problem
		n3iwfUe = amf.FindUeByAmfUeNgapID(ueNgapIDs.AMFUENGAPID.Value)
	}

	if n3iwfUe == nil {
		// TODO: send error indication(unknown local ngap ue id)
		return
	}

	if cause != nil {
		printAndGetCause(cause)
	}

	// TODO: release pdu session and gtp info for ue
	n3iwfUe.Remove()

	ngap_message.SendUEContextReleaseComplete(amf, n3iwfUe, nil)
}

func HandleDownlinkNASTransport(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Downlink NAS Transport")

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}

	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	var oldAMF *ngapType.AMFName
	var nasPDU *ngapType.NASPDU
	var indexToRFSP *ngapType.IndexToRFSP
	var ueAggregateMaximumBitRate *ngapType.UEAggregateMaximumBitRate
	var allowedNSSAI *ngapType.AllowedNSSAI
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var n3iwfUe *n3iwf_context.N3IWFUe
	var n3iwfSelf = n3iwf_context.N3IWFSelf()

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("Initiating Message is nil")
		return
	}

	downlinkNASTransport := initiatingMessage.Value.DownlinkNASTransport
	if downlinkNASTransport == nil {
		ngapLog.Error("DownlinkNASTransport is nil")
		return
	}

	for _, ie := range downlinkNASTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			amfUeNgapID = ie.Value.AMFUENGAPID
			if amfUeNgapID == nil {
				ngapLog.Errorf("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Errorf("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDOldAMF:
			ngapLog.Traceln("[NGAP] Decode IE OldAMF")
			oldAMF = ie.Value.OldAMF
		case ngapType.ProtocolIEIDNASPDU:
			ngapLog.Traceln("[NGAP] Decode IE NASPDU")
			nasPDU = ie.Value.NASPDU
			if nasPDU == nil {
				ngapLog.Errorf("NASPDU is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDIndexToRFSP:
			ngapLog.Traceln("[NGAP] Decode IE IndexToRFSP")
			indexToRFSP = ie.Value.IndexToRFSP
		case ngapType.ProtocolIEIDUEAggregateMaximumBitRate:
			ngapLog.Traceln("[NGAP] Decode IE UEAggregateMaximumBitRate")
			ueAggregateMaximumBitRate = ie.Value.UEAggregateMaximumBitRate
		case ngapType.ProtocolIEIDAllowedNSSAI:
			ngapLog.Traceln("[NGAP] Decode IE AllowedNSSAI")
			allowedNSSAI = ie.Value.AllowedNSSAI
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		// TODO: Send Error Indication
	}

	if ranUeNgapID != nil {
		n3iwfUe = n3iwfSelf.FindUeByRanUeNgapID(ranUeNgapID.Value)
		if n3iwfUe == nil {
			ngapLog.Warnf("No UE Context[RanUeNgapID:%d]\n", ranUeNgapID.Value)
			return
		}
	}

	if amfUeNgapID != nil {
		if n3iwfUe.AmfUeNgapId == n3iwf_context.AmfUeNgapIdUnspecified {
			ngapLog.Tracef("Create new logical UE-associated NG-connection")
			n3iwfUe.AmfUeNgapId = amfUeNgapID.Value
			// n3iwfUe.SCTPAddr = amf.SCTPAddr
		} else {
			if n3iwfUe.AmfUeNgapId != amfUeNgapID.Value {
				ngapLog.Warn("AMFUENGAPID unmatched")
				return
			}
		}
	}

	if oldAMF != nil {
		ngapLog.Debugf("Old AMF: %s\n", oldAMF.Value)
	}

	if indexToRFSP != nil {
		n3iwfUe.IndexToRfsp = indexToRFSP.Value
	}

	if ueAggregateMaximumBitRate != nil {
		n3iwfUe.Ambr = ueAggregateMaximumBitRate
	}

	if allowedNSSAI != nil {
		n3iwfUe.AllowedNssai = allowedNSSAI
	}

	if nasPDU != nil {
		// TODO: Send NAS PDU to UE
		if n3iwfUe.UDPSendInfoGroup != nil {
			var identifier uint8
			ikeSecurityAssociation := n3iwfUe.N3IWFIKESecurityAssociation

			for {
				identifier = uint8(rand.Uint32())
				if identifier != n3iwfUe.N3IWFIKESecurityAssociation.LastEAPIdentifier {
					n3iwfUe.N3IWFIKESecurityAssociation.LastEAPIdentifier = identifier
					break
				}
			}
			// Send NAS via IKE EAP
			// IKEHDR-SK-{response}
			var eap *ike_message.EAP

			// Build IKE message
			responseIKEMessage := ike_message.BuildIKEHeader(ikeSecurityAssociation.RemoteSPI, ikeSecurityAssociation.LocalSPI, ike_message.IKE_AUTH, ike_message.ResponseBitCheck, ikeSecurityAssociation.MessageID)

			// Build response
			var ikePayload []ike_message.IKEPayloadType

			// EAP-5G
			eap = ike_message.BuildEAP5GNAS(identifier, nasPDU.Value)
			ikePayload = append(ikePayload, eap)

			if err := ike_handler.EncryptProcedure(ikeSecurityAssociation, ikePayload, responseIKEMessage); err != nil {
				ngapLog.Errorf("[NGAP] Encrypting IKE message failed: %+v", err)
				return
			}

			// Send IKE message to UE
			ike_handler.SendIKEMessageToUE(n3iwfUe.UDPSendInfoGroup, responseIKEMessage)
		} else {
			for i := 0; i < 3; i++ {
				if n3iwfUe.TCPConnection == nil {
					ngapLog.Warnf("No NAS signalling session found, retry %d", i+1)
					time.Sleep(500 * time.Millisecond)
				} else {
					break
				}
			}
			if n3iwfUe.TCPConnection != nil {
				if n, err := n3iwfUe.TCPConnection.Write(nasPDU.Value); err != nil {
					ngapLog.Errorf("Writing via IPSec signalling SA failed: %+v", err)
				} else {
					ngapLog.Infof("Wrote %d bytes", n)
				}
			} else {
				ngapLog.Error("No connection found for UE to send NAS message")
			}
		}
	}
}

func HandlePDUSessionResourceSetupRequest(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle PDU Session Resource Setup Request")

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}

	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	var nasPDU *ngapType.NASPDU
	var pduSessionResourceSetupListSUReq *ngapType.PDUSessionResourceSetupListSUReq
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var n3iwfUe *n3iwf_context.N3IWFUe
	var n3iwfSelf = n3iwf_context.N3IWFSelf()

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("Initiating Message is nil")
		return
	}

	pduSessionResourceSetupRequest := initiatingMessage.Value.PDUSessionResourceSetupRequest
	if pduSessionResourceSetupRequest == nil {
		ngapLog.Error("PDUSessionResourceSetupRequest is nil")
		return
	}

	for _, ie := range pduSessionResourceSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			amfUeNgapID = ie.Value.AMFUENGAPID
			if amfUeNgapID == nil {
				ngapLog.Errorf("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Errorf("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDNASPDU:
			ngapLog.Traceln("[NGAP] Decode IE NASPDU")
			nasPDU = ie.Value.NASPDU
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListSUReq:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionResourceSetupRequestList")
			pduSessionResourceSetupListSUReq = ie.Value.PDUSessionResourceSetupListSUReq
			if pduSessionResourceSetupListSUReq == nil {
				ngapLog.Errorf("PDUSessionResourceSetupRequestList is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		// TODO: Send error indication to AMF
		ngapLog.Errorln("Sending error indication to AMF")
		return
	}

	if (amfUeNgapID != nil) && (ranUeNgapID != nil) {
		// Find UE context
		n3iwfUe = n3iwfSelf.FindUeByRanUeNgapID(ranUeNgapID.Value)
		if n3iwfUe == nil {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
			// TODO: build cause and handle error
			// Cause: Unknown local UE NGAP ID
			return
		} else {
			if n3iwfUe.AmfUeNgapId != amfUeNgapID.Value {
				// TODO: build cause and handle error
				// Cause: Inconsistent remote UE NGAP ID
				return
			}
		}
	}

	if nasPDU != nil {
		// TODO: Send NAS to UE
		if n3iwfUe.TCPConnection == nil {
			ngapLog.Error("No IPSec NAS signalling SA for this UE")
			return
		} else {
			if n, err := n3iwfUe.TCPConnection.Write(nasPDU.Value); err != nil {
				ngapLog.Errorf("Send NAS to UE failed: %+v", err)
				return
			} else {
				ngapLog.Tracef("Wrote %d bytes", n)
			}
		}
	}

	if pduSessionResourceSetupListSUReq != nil {
		setupListSURes := new(ngapType.PDUSessionResourceSetupListSURes)
		failedListSURes := new(ngapType.PDUSessionResourceFailedToSetupListSURes)
		// UE temporary data for PDU session setup response
		n3iwfUe.TemporaryPDUSessionSetupData = &n3iwf_context.PDUSessionSetupTemporaryData{
			SetupListSURes:  setupListSURes,
			FailedListSURes: failedListSURes,
		}
		n3iwfUe.TemporaryPDUSessionSetupData.NGAPProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup

		for _, item := range pduSessionResourceSetupListSUReq.List {
			pduSessionID := item.PDUSessionID.Value
			// TODO: send NAS to UE
			// pduSessionNasPdu := item.NASPDU
			snssai := item.SNSSAI

			transfer := ngapType.PDUSessionResourceSetupRequestTransfer{}
			err := aper.UnmarshalWithParams(item.PDUSessionResourceSetupRequestTransfer, &transfer, "valueExt")
			if err != nil {
				ngapLog.Errorf("[PDUSessionID: %d] PDUSessionResourceSetupRequestTransfer Decode Error: %+v\n", pduSessionID, err)
			}

			pduSession, err := n3iwfUe.CreatePDUSession(pduSessionID, snssai)
			if err != nil {
				ngapLog.Errorf("Create PDU Session Error: %+v\n", err)

				cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentMultiplePDUSessionIDInstances)
				unsuccessfulTransfer, buildErr := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
				if buildErr != nil {
					ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", buildErr)
				}
				ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(failedListSURes, pduSessionID, unsuccessfulTransfer)
				continue
			}

			success, resTransfer := handlePDUSessionResourceSetupRequestTransfer(n3iwfUe, pduSession, transfer)
			if success {
				// Append this PDU session to unactivated PDU session list
				n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession = append(n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession, pduSessionID)
			} else {
				// Delete the pdusession store in UE conext
				delete(n3iwfUe.PduSessionList, pduSessionID)
				ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(failedListSURes, pduSessionID, resTransfer)
			}
		}
	}

	if n3iwfUe.TemporaryPDUSessionSetupData != nil {
		for {
			if len(n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession) != 0 {
				pduSessionID := n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession[0]
				pduSession := n3iwfUe.PduSessionList[pduSessionID]

				ikeSecurityAssociation := n3iwfUe.N3IWFIKESecurityAssociation

				// Send CREATE_CHILD_SA to UE
				// Add MessageID for IKE security association
				ikeSecurityAssociation.MessageID++
				ikeMessage := ike_message.BuildIKEHeader(ikeSecurityAssociation.LocalSPI, ikeSecurityAssociation.RemoteSPI, ike_message.CREATE_CHILD_SA, ike_message.InitiatorBitCheck, ikeSecurityAssociation.MessageID)

				// IKE payload
				var ikePayload []ike_message.IKEPayloadType

				// Build SA
				// Proposals
				var proposals []*ike_message.Proposal

				// Allocate SPI
				var spi uint32
				spiByte := make([]byte, 4)
				for {
					randomUint64 := ike_handler.GenerateRandomNumber().Uint64()
					if _, ok := n3iwfSelf.ChildSA[uint32(randomUint64)]; !ok {
						spi = uint32(randomUint64)
						break
					}
				}
				binary.BigEndian.PutUint32(spiByte, spi)

				// First Proposal - Proposal No.1
				proposal := ike_message.BuildProposal(1, ike_message.TypeESP, spiByte)

				// Encryption transform
				var attributeType uint16 = ike_message.AttributeTypeKeyLength
				var attributeValue uint16 = 256
				encryptionTransform := ike_message.BuildTransform(ike_message.TypeEncryptionAlgorithm, ike_message.ENCR_AES_CBC, &attributeType, &attributeValue, nil)
				if ok := ike_message.AppendTransformToProposal(proposal, encryptionTransform); !ok {
					ngapLog.Error("Generate IKE message failed: Cannot append to proposal")
					n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession = n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
					cause := buildCause(ngapType.CausePresentTransport, ngapType.CauseTransportPresentTransportResourceUnavailable)
					transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
					if err != nil {
						ngapLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
						continue
					}
					ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(n3iwfUe.TemporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
					continue
				}
				// Integrity transform
				if pduSession.SecurityIntegrity {
					integrityTransform := ike_message.BuildTransform(ike_message.TypeIntegrityAlgorithm, ike_message.AUTH_HMAC_SHA1_96, nil, nil, nil)
					if ok := ike_message.AppendTransformToProposal(proposal, integrityTransform); !ok {
						ngapLog.Error("Generate IKE message failed: Cannot append to proposal")
						n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession = n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
						cause := buildCause(ngapType.CausePresentTransport, ngapType.CauseTransportPresentTransportResourceUnavailable)
						transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
						if err != nil {
							ngapLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
							continue
						}
						ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(n3iwfUe.TemporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
						continue
					}
				}
				// ESN transform
				esnTransform := ike_message.BuildTransform(ike_message.TypeExtendedSequenceNumbers, ike_message.ESN_NO, nil, nil, nil)
				if ok := ike_message.AppendTransformToProposal(proposal, esnTransform); !ok {
					ngapLog.Error("Generate IKE message failed: Cannot append to proposal")
					n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession = n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
					cause := buildCause(ngapType.CausePresentTransport, ngapType.CauseTransportPresentTransportResourceUnavailable)
					transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
					if err != nil {
						ngapLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
						continue
					}
					ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(n3iwfUe.TemporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
					continue
				}

				proposals = append(proposals, proposal)

				securityAssociation := ike_message.BuildSecurityAssociation(proposals)

				ikePayload = append(ikePayload, securityAssociation)

				// Build Nonce
				nonceData := ike_handler.GenerateRandomNumber().Bytes()
				nonce := ike_message.BuildNonce(nonceData)

				// Store nonce into context
				ikeSecurityAssociation.ConcatenatedNonce = nonceData

				ikePayload = append(ikePayload, nonce)

				// TSi
				n3iwfIPAddr := net.ParseIP(n3iwfSelf.IPSecGatewayAddress)
				individualTrafficSelector := ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
					0, 65535, n3iwfIPAddr.To4(), n3iwfIPAddr.To4())
				trafficSelectorInitiator := ike_message.BuildTrafficSelectorInitiator([]*ike_message.IndividualTrafficSelector{individualTrafficSelector})

				ikePayload = append(ikePayload, trafficSelectorInitiator)

				// TSr
				ueIPAddr := net.ParseIP(n3iwfUe.IPSecInnerIP)
				individualTrafficSelector = ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
					0, 65535, ueIPAddr.To4(), ueIPAddr.To4())
				trafficSelectorResponder := ike_message.BuildTrafficSelectorResponder([]*ike_message.IndividualTrafficSelector{individualTrafficSelector})

				ikePayload = append(ikePayload, trafficSelectorResponder)

				// Notify-Qos
				notifyQos := ike_message.BuildNotify5G_QOS_INFO(uint8(pduSessionID), pduSession.QFIList, true)

				ikePayload = append(ikePayload, notifyQos)

				// Notify-UP_IP_ADDRESS
				notifyUPIPAddr := ike_message.BuildNotifyUP_IP4_ADDRESS(n3iwfSelf.IPSecGatewayAddress)

				ikePayload = append(ikePayload, notifyUPIPAddr)

				if err := ike_handler.EncryptProcedure(n3iwfUe.N3IWFIKESecurityAssociation, ikePayload, ikeMessage); err != nil {
					ngapLog.Errorf("Encrypting IKE message failed: %+v", err)
					n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession = n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
					cause := buildCause(ngapType.CausePresentTransport, ngapType.CauseTransportPresentTransportResourceUnavailable)
					transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
					if err != nil {
						ngapLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
						continue
					}
					ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(n3iwfUe.TemporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
					continue
				}

				ueUDPAddr, err := net.ResolveUDPAddr("udp", n3iwfUe.IPAddrv4+":500")
				if err != nil {
					ngapLog.Errorf("Resolve UE UDP address failed: %+v", err)
					n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession = n3iwfUe.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
					cause := buildCause(ngapType.CausePresentTransport, ngapType.CauseTransportPresentTransportResourceUnavailable)
					transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
					if err != nil {
						ngapLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
						continue
					}
					ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(n3iwfUe.TemporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
					continue
				}

				ueSendInfo := &n3iwf_message.UDPSendInfoGroup{
					ChannelID: udp_server.ChannelIDForPort500,
					Addr:      ueUDPAddr,
				}

				ike_handler.SendIKEMessageToUE(ueSendInfo, ikeMessage)
				break
			} else {
				// Send PDU Session Resource Setup Response to AMF
				ngap_message.SendPDUSessionResourceSetupResponse(amf, n3iwfUe, n3iwfUe.TemporaryPDUSessionSetupData.SetupListSURes, n3iwfUe.TemporaryPDUSessionSetupData.FailedListSURes, nil)
				break
			}
		}
	}

}

func HandlePDUSessionResourceModifyRequest(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle PDU Session Resource Modify Request")

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}

	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	var pduSessionResourceModifyListModReq *ngapType.PDUSessionResourceModifyListModReq
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var n3iwfUe *n3iwf_context.N3IWFUe
	var n3iwfSelf = n3iwf_context.N3IWFSelf()

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("Initiating Message is nil")
		return
	}

	pduSessionResourceModifyRequest := initiatingMessage.Value.PDUSessionResourceModifyRequest
	if pduSessionResourceModifyRequest == nil {
		ngapLog.Error("PDUSessionResourceModifyRequest is nil")
		return
	}

	for _, ie := range pduSessionResourceModifyRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			amfUeNgapID = ie.Value.AMFUENGAPID
			if amfUeNgapID == nil {
				ngapLog.Error("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModReq:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionResourceModifyListModReq")
			pduSessionResourceModifyListModReq = ie.Value.PDUSessionResourceModifyListModReq
			if pduSessionResourceModifyListModReq == nil {
				ngapLog.Error("PDUSessionResourceModifyListModReq is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)
		ngap_message.SendPDUSessionResourceModifyResponse(amf, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	if (amfUeNgapID != nil) && (ranUeNgapID != nil) {
		// Find UE context
		n3iwfUe = n3iwfSelf.FindUeByRanUeNgapID(ranUeNgapID.Value)
		if n3iwfUe == nil {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
			// TODO: build cause and send error indication
			// Cause: Unknown local UE NGAP ID
			return
		} else {
			if n3iwfUe.AmfUeNgapId != amfUeNgapID.Value {
				// TODO: build cause and send error indication
				// Cause: Inconsistent remote UE NGAP ID
				return
			}
		}
	}

	responseList := new(ngapType.PDUSessionResourceModifyListModRes)
	failedListModRes := new(ngapType.PDUSessionResourceFailedToModifyListModRes)
	if pduSessionResourceModifyListModReq != nil {
		var pduSession *n3iwf_context.PDUSession
		for _, item := range pduSessionResourceModifyListModReq.List {
			pduSessionID := item.PDUSessionID.Value
			// TODO: send NAS to UE
			// pduSessionNasPdu := item.NASPDU
			transfer := ngapType.PDUSessionResourceModifyRequestTransfer{}
			err := aper.UnmarshalWithParams(item.PDUSessionResourceModifyRequestTransfer, transfer, "valueExt")
			if err != nil {
				ngapLog.Errorf("[PDUSessionID: %d] PDUSessionResourceModifyRequestTransfer Decode Error: %+v\n", pduSessionID, err)
			}

			if pduSession = n3iwfUe.FindPDUSession(pduSessionID); pduSession == nil {
				ngapLog.Errorf("[PDUSessionID: %d] Unknown PDU session ID", pduSessionID)

				cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentUnknownPDUSessionID)
				unsuccessfulTransfer, buildErr := ngap_message.BuildPDUSessionResourceModifyUnsuccessfulTransfer(*cause, nil)
				if buildErr != nil {
					ngapLog.Errorf("Build PDUSessionResourceModifyUnsuccessfulTransfer Error: %+v\n", buildErr)
				}
				ngap_message.AppendPDUSessionResourceFailedToModifyListModRes(failedListModRes, pduSessionID, unsuccessfulTransfer)
				continue
			}

			success, resTransfer := handlePDUSessionResourceModifyRequestTransfer(pduSession, transfer)
			if success {
				ngap_message.AppendPDUSessionResourceModifyListModRes(responseList, pduSessionID, resTransfer)
			} else {
				ngap_message.AppendPDUSessionResourceFailedToModifyListModRes(failedListModRes, pduSessionID, resTransfer)
			}
		}
	}

	ngap_message.SendPDUSessionResourceModifyResponse(amf, n3iwfUe, responseList, failedListModRes, nil)
}

func handlePDUSessionResourceModifyRequestTransfer(pduSession *n3iwf_context.PDUSession, transfer ngapType.PDUSessionResourceModifyRequestTransfer) (success bool, responseTransfer []byte) {
	ngapLog.Trace("[N3IWF] Handle PDU Session Resource Modify Request Transfer")

	var pduSessionAMBR *ngapType.PDUSessionAggregateMaximumBitRate
	var ulNGUUPTNLModifyList *ngapType.ULNGUUPTNLModifyList
	var networkInstance *ngapType.NetworkInstance
	var qosFlowAddOrModifyRequestList *ngapType.QosFlowAddOrModifyRequestList
	var qosFlowToReleaseList *ngapType.QosFlowList
	var additionalULNGUUPTNLInformation *ngapType.UPTransportLayerInformation

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	// used for building response transfer
	var resDLNGUUPTNLInfo *ngapType.UPTransportLayerInformation
	var resULNGUUPTNLInfo *ngapType.UPTransportLayerInformation
	var resQosFlowAddOrModifyRequestList ngapType.QosFlowAddOrModifyResponseList
	var resQosFlowFailedToAddOrModifyList ngapType.QosFlowList

	for _, ie := range transfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDPDUSessionAggregateMaximumBitRate:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionAggregateMaximumBitRate")
			pduSessionAMBR = ie.Value.PDUSessionAggregateMaximumBitRate
		case ngapType.ProtocolIEIDULNGUUPTNLModifyList:
			ngapLog.Traceln("[NGAP] Decode IE ULNGUUPTNLModifyList")
			ulNGUUPTNLModifyList = ie.Value.ULNGUUPTNLModifyList
			if ulNGUUPTNLModifyList != nil && len(ulNGUUPTNLModifyList.List) == 0 {
				ngapLog.Error("ULNGUUPTNLModifyList should have at least one element")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDNetworkInstance:
			ngapLog.Traceln("[NGAP] Decode IE NetworkInstance")
			networkInstance = ie.Value.NetworkInstance
		case ngapType.ProtocolIEIDQosFlowAddOrModifyRequestList:
			ngapLog.Traceln("[NGAP] Decode IE QosFLowAddOrModifyRequestList")
			qosFlowAddOrModifyRequestList = ie.Value.QosFlowAddOrModifyRequestList
			if qosFlowAddOrModifyRequestList != nil && len(qosFlowAddOrModifyRequestList.List) == 0 {
				ngapLog.Error("QosFlowAddOrModifyRequestList should have at least one element")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDQosFlowToReleaseList:
			ngapLog.Traceln("[NGAP] Decode IE QosFlowToReleaseList")
			qosFlowToReleaseList = ie.Value.QosFlowToReleaseList
			if qosFlowToReleaseList != nil && len(qosFlowToReleaseList.List) == 0 {
				ngapLog.Error("qosFlowToReleaseList should have at least one element")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDAdditionalULNGUUPTNLInformation:
			ngapLog.Traceln("[NGAP] Decode IE AdditionalULNGUUPTNLInformation")
			additionalULNGUUPTNLInformation = ie.Value.AdditionalULNGUUPTNLInformation
		}
	}

	if len(iesCriticalityDiagnostics.List) != 0 {
		// build unsuccessful transfer
		cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)
		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)
		unsuccessfulTransfer, err := ngap_message.BuildPDUSessionResourceModifyUnsuccessfulTransfer(*cause, &criticalityDiagnostics)
		if err != nil {
			ngapLog.Errorf("Build PDUSessionResourceModifyUnsuccessfulTransfer Error: %+v\n", err)
		}

		responseTransfer = unsuccessfulTransfer
		return
	}

	if ulNGUUPTNLModifyList != nil {
		updateItem := ulNGUUPTNLModifyList.List[0]

		// TODO: update GTP tunnel

		ngapLog.Info("Update uplink NG-U user plane tunnel information")

		resULNGUUPTNLInfo = &updateItem.ULNGUUPTNLInformation
		resDLNGUUPTNLInfo = &updateItem.DLNGUUPTNLInformation
	}

	if qosFlowAddOrModifyRequestList != nil {
		for _, updateItem := range qosFlowAddOrModifyRequestList.List {
			target, ok := pduSession.QosFlows[updateItem.QosFlowIdentifier.Value]
			if ok {
				ngapLog.Trace("Update qos flow level qos parameters")

				target.Parameters = *updateItem.QosFlowLevelQosParameters

				item := ngapType.QosFlowAddOrModifyResponseItem{
					QosFlowIdentifier: updateItem.QosFlowIdentifier,
				}

				resQosFlowAddOrModifyRequestList.List = append(resQosFlowAddOrModifyRequestList.List, item)
			} else {
				ngapLog.Errorf("Requested Qos flow not found, QosFlowID: %d", updateItem.QosFlowIdentifier)

				cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentUnkownQosFlowID)

				item := ngapType.QosFlowItem{
					QosFlowIdentifier: updateItem.QosFlowIdentifier,
					Cause:             *cause,
				}

				resQosFlowFailedToAddOrModifyList.List = append(resQosFlowFailedToAddOrModifyList.List, item)
			}
		}
	}

	if pduSessionAMBR != nil {
		ngapLog.Trace("Store PDU session AMBR")
		pduSession.Ambr = pduSessionAMBR
	}

	if networkInstance != nil {
		// Used to select transport layer resource
		ngapLog.Trace("Store network instance")
		pduSession.NetworkInstance = networkInstance
	}

	if qosFlowToReleaseList != nil {
		for _, releaseItem := range qosFlowToReleaseList.List {
			_, ok := pduSession.QosFlows[releaseItem.QosFlowIdentifier.Value]
			if ok {
				ngapLog.Tracef("Delete QosFlow. ID: %d", releaseItem.QosFlowIdentifier.Value)
				printAndGetCause(&releaseItem.Cause)
				delete(pduSession.QosFlows, releaseItem.QosFlowIdentifier.Value)
			}
		}
	}

	if additionalULNGUUPTNLInformation != nil {
		// TODO: forward AdditionalULNGUUPTNLInfomation to S-NG-RAN
	}

	encodeData, err := ngap_message.BuildPDUSessionResourceModifyResponseTransfer(resULNGUUPTNLInfo, resDLNGUUPTNLInfo, &resQosFlowAddOrModifyRequestList, &resQosFlowFailedToAddOrModifyList)
	if err != nil {
		ngapLog.Errorf("Build PDUSessionResourceModifyTransfer Error: %+v\n", err)
	}

	success = true
	responseTransfer = encodeData

	return
}

func HandlePDUSessionResourceModifyConfirm(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle PDU Session Resource Modify Confirm")

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceModifyListModCfm *ngapType.PDUSessionResourceModifyListModCfm
	var pDUSessionResourceFailedToModifyListModCfm *ngapType.PDUSessionResourceFailedToModifyListModCfm
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ngapLog.Error("Successful Outcome is nil")
		return
	}

	pDUSessionResourceModifyConfirm := successfulOutcome.Value.PDUSessionResourceModifyConfirm
	if pDUSessionResourceModifyConfirm == nil {
		ngapLog.Error("pDUSessionResourceModifyConfirm is nil")
		return
	}

	for _, ie := range pDUSessionResourceModifyConfirm.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModCfm:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionResourceModifyListModCfm")
			pDUSessionResourceModifyListModCfm = ie.Value.PDUSessionResourceModifyListModCfm
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToModifyListModCfm:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionResourceFailedToModifyListModCfm")
			pDUSessionResourceFailedToModifyListModCfm = ie.Value.PDUSessionResourceFailedToModifyListModCfm
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			ngapLog.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	var ue *n3iwf_context.N3IWFUe
	var n3iwfSelf = n3iwf_context.N3IWFSelf()

	if rANUENGAPID != nil {
		ue = n3iwfSelf.FindUeByRanUeNgapID(rANUENGAPID.Value)
		if ue == nil {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", rANUENGAPID.Value)
			return
		}
	}
	if aMFUENGAPID != nil {
		if ue != nil {
			if ue.AmfUeNgapId != aMFUENGAPID.Value {
				ngapLog.Errorf("Inconsistent remote UE NGAP ID, AMFUENGAPID: %d, ue.AmfUeNgapId: %d", aMFUENGAPID.Value, ue.AmfUeNgapId)
				return
			}
		} else {
			ue = amf.FindUeByAmfUeNgapID(aMFUENGAPID.Value)
			if ue == nil {
				ngapLog.Errorf("Inconsistent remote UE NGAP ID, AMFUENGAPID: %d, ue.AmfUeNgapId: %d", aMFUENGAPID.Value, ue.AmfUeNgapId)
				return
			}
		}
	}
	if ue == nil {
		ngapLog.Warn("RANUENGAPID and  AMFUENGAPID are both nil")
		return
	}
	if pDUSessionResourceModifyListModCfm != nil {
		for _, item := range pDUSessionResourceModifyListModCfm.List {
			pduSessionId := item.PDUSessionID.Value
			ngapLog.Tracef("PDU Session Id[%d] in Pdu Session Resouce Modification Confrim List", pduSessionId)
			sess, exist := ue.PduSessionList[pduSessionId]
			if !exist {
				ngapLog.Warnf("PDU Session Id[%d] is not exist in Ue[ranUeNgapId:%d]", pduSessionId, ue.RanUeNgapId)
			} else {
				transfer := ngapType.PDUSessionResourceModifyConfirmTransfer{}
				err := aper.UnmarshalWithParams(item.PDUSessionResourceModifyConfirmTransfer, &transfer, "valueExt")
				if err != nil {
					ngapLog.Warnf("[PDUSessionID: %d] PDUSessionResourceSetupRequestTransfer Decode Error: %+v\n", pduSessionId, err)
				} else if transfer.QosFlowFailedToModifyList != nil {
					for _, flow := range transfer.QosFlowFailedToModifyList.List {
						ngapLog.Warnf("Delete QFI[%d] due to Qos Flow Failure in Pdu Session Resouce Modification Confrim List", flow.QosFlowIdentifier.Value)
						delete(sess.QosFlows, flow.QosFlowIdentifier.Value)
					}
				}
			}
		}
	}
	if pDUSessionResourceFailedToModifyListModCfm != nil {
		for _, item := range pDUSessionResourceFailedToModifyListModCfm.List {
			pduSessionId := item.PDUSessionID.Value
			transfer := ngapType.PDUSessionResourceModifyIndicationUnsuccessfulTransfer{}
			err := aper.UnmarshalWithParams(item.PDUSessionResourceModifyIndicationUnsuccessfulTransfer, &transfer, "valueExt")
			if err != nil {
				ngapLog.Warnf("[PDUSessionID: %d] PDUSessionResourceModifyIndicationUnsuccessfulTransfer Decode Error: %+v\n", pduSessionId, err)
			} else {
				printAndGetCause(&transfer.Cause)
			}
			ngapLog.Tracef("Release PDU Session Id[%d] due to PDU Session Resource Modify Indication Unsuccessful", pduSessionId)
			delete(ue.PduSessionList, pduSessionId)
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

}

func HandlePDUSessionResourceReleaseCommand(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle PDU Session Resource Release Command")
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rANPagingPriority *ngapType.RANPagingPriority
	var nASPDU *ngapType.NASPDU
	var pDUSessionResourceToReleaseListRelCmd *ngapType.PDUSessionResourceToReleaseListRelCmd

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("Initiating Message is nil")
		return
	}

	pDUSessionResourceReleaseCommand := initiatingMessage.Value.PDUSessionResourceReleaseCommand
	if pDUSessionResourceReleaseCommand == nil {
		ngapLog.Error("pDUSessionResourceReleaseCommand is nil")
		return
	}

	for _, ie := range pDUSessionResourceReleaseCommand.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				ngapLog.Error("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				ngapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANPagingPriority:
			ngapLog.Traceln("[NGAP] Decode IE RANPagingPriority")
			rANPagingPriority = ie.Value.RANPagingPriority
		case ngapType.ProtocolIEIDNASPDU:
			ngapLog.Traceln("[NGAP] Decode IE NASPDU")
			nASPDU = ie.Value.NASPDU
		case ngapType.ProtocolIEIDPDUSessionResourceToReleaseListRelCmd:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionResourceToReleaseListRelCmd")
			pDUSessionResourceToReleaseListRelCmd = ie.Value.PDUSessionResourceToReleaseListRelCmd
			if pDUSessionResourceToReleaseListRelCmd == nil {
				ngapLog.Error("PDUSessionResourceToReleaseListRelCmd is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodePDUSessionResourceRelease
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(amf, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	n3iwfSelf := n3iwf_context.N3IWFSelf()
	ue := n3iwfSelf.FindUeByRanUeNgapID(rANUENGAPID.Value)
	if ue == nil {
		ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", rANUENGAPID.Value)
		cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID)
		ngap_message.SendErrorIndication(amf, nil, nil, cause, nil)
		return
	}

	if ue.AmfUeNgapId != aMFUENGAPID.Value {
		ngapLog.Errorf("Inconsistent remote UE NGAP ID, AMFUENGAPID: %d, ue.AmfUeNgapId: %d", aMFUENGAPID.Value, ue.AmfUeNgapId)
		cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentInconsistentRemoteUENGAPID)
		ngap_message.SendErrorIndication(amf, nil, &rANUENGAPID.Value, cause, nil)
		return
	}

	if rANPagingPriority != nil {
		//n3iwf does not support paging
	}

	releaseList := ngapType.PDUSessionResourceReleasedListRelRes{}
	for _, item := range pDUSessionResourceToReleaseListRelCmd.List {
		pduSessionId := item.PDUSessionID.Value
		transfer := ngapType.PDUSessionResourceReleaseCommandTransfer{}
		err := aper.UnmarshalWithParams(item.PDUSessionResourceReleaseCommandTransfer, &transfer, "valueExt")
		if err != nil {
			ngapLog.Warnf("[PDUSessionID: %d] PDUSessionResourceReleaseCommandTransfer Decode Error: %+v\n", pduSessionId, err)
		} else {
			printAndGetCause(&transfer.Cause)
		}
		ngapLog.Tracef("Release PDU Session Id[%d] due to PDU Session Resource Release Command", pduSessionId)
		delete(ue.PduSessionList, pduSessionId)

		// reponse list
		releaseItem := ngapType.PDUSessionResourceReleasedItemRelRes{
			PDUSessionID: item.PDUSessionID,
			PDUSessionResourceReleaseResponseTransfer: getPDUSessionResourceReleaseResponseTransfer(),
		}
		releaseList.List = append(releaseList.List, releaseItem)
	}

	if nASPDU != nil {
		// TODO: Send NAS to UE
	}
	ngap_message.SendPDUSessionResourceReleaseResponse(amf, ue, releaseList, nil)

}

func HandleErrorIndication(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Error Indication")

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}
	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("InitiatingMessage is nil")
		return
	}
	errorIndication := initiatingMessage.Value.ErrorIndication
	if errorIndication == nil {
		ngapLog.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range errorIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			ngapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ngapLog.Trace("[NGAP] Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			ngapLog.Trace("[NGAP] Decode IE Cause")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ngapLog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	if cause == nil && criticalityDiagnostics == nil {
		ngapLog.Error("Both Cause IE and CriticalityDiagnostics IE are nil, should have at least one")
		return
	}

	if (aMFUENGAPID == nil) != (rANUENGAPID == nil) {
		ngapLog.Error("One of UE NGAP ID is not included in this message")
		return
	}

	if (aMFUENGAPID != nil) && (rANUENGAPID != nil) {
		ngapLog.Trace("UE-associated procedure error")
		ngapLog.Warnf("AMF UE NGAP ID is defined, value = %d", aMFUENGAPID.Value)
		ngapLog.Warnf("RAN UE NGAP ID is defined, value = %d", rANUENGAPID.Value)
	}

	if cause != nil {
		printAndGetCause(cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

	// TODO: handle error based on cause/criticalityDiagnostics
}

func HandleUERadioCapabilityCheckRequest(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle UE Radio Capability Check Request")
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var uERadioCapability *ngapType.UERadioCapability

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("InitiatingMessage is nil")
		return
	}

	uERadioCapabilityCheckRequest := initiatingMessage.Value.UERadioCapabilityCheckRequest
	if uERadioCapabilityCheckRequest == nil {
		ngapLog.Error("uERadioCapabilityCheckRequest is nil")
		return
	}

	for _, ie := range uERadioCapabilityCheckRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				ngapLog.Error("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				ngapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDUERadioCapability:
			ngapLog.Traceln("[NGAP] Decode IE UERadioCapability")
			uERadioCapability = ie.Value.UERadioCapability
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodeUERadioCapabilityCheck
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(amf, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	n3iwfSelf := n3iwf_context.N3IWFSelf()
	ue := n3iwfSelf.FindUeByRanUeNgapID(rANUENGAPID.Value)
	if ue == nil {
		ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", rANUENGAPID.Value)
		cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID)
		ngap_message.SendErrorIndication(amf, nil, nil, cause, nil)
		return
	}

	ue.RadioCapability = uERadioCapability

}

func HandleAMFConfigurationUpdate(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {

	ngapLog.Infoln("[N3IWF] Handle AMF Configuration Updaet")

	var aMFName *ngapType.AMFName
	var servedGUAMIList *ngapType.ServedGUAMIList
	var relativeAMFCapacity *ngapType.RelativeAMFCapacity
	var pLMNSupportList *ngapType.PLMNSupportList
	var aMFTNLAssociationToAddList *ngapType.AMFTNLAssociationToAddList
	var aMFTNLAssociationToRemoveList *ngapType.AMFTNLAssociationToRemoveList
	var aMFTNLAssociationToUpdateList *ngapType.AMFTNLAssociationToUpdateList

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("InitiatingMessage is nil")
		return
	}

	aMFConfigurationUpdate := initiatingMessage.Value.AMFConfigurationUpdate
	if aMFConfigurationUpdate == nil {
		ngapLog.Error("aMFConfigurationUpdate is nil")
		return
	}

	for _, ie := range aMFConfigurationUpdate.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFName:
			ngapLog.Traceln("[NGAP] Decode IE AMFName")
			aMFName = ie.Value.AMFName
		case ngapType.ProtocolIEIDServedGUAMIList:
			ngapLog.Traceln("[NGAP] Decode IE ServedGUAMIList")
			servedGUAMIList = ie.Value.ServedGUAMIList
		case ngapType.ProtocolIEIDRelativeAMFCapacity:
			ngapLog.Traceln("[NGAP] Decode IE RelativeAMFCapacity")
			relativeAMFCapacity = ie.Value.RelativeAMFCapacity
		case ngapType.ProtocolIEIDPLMNSupportList:
			ngapLog.Traceln("[NGAP] Decode IE PLMNSupportList")
			pLMNSupportList = ie.Value.PLMNSupportList
		case ngapType.ProtocolIEIDAMFTNLAssociationToAddList:
			ngapLog.Traceln("[NGAP] Decode IE AMFTNLAssociationToAddList")
			aMFTNLAssociationToAddList = ie.Value.AMFTNLAssociationToAddList
		case ngapType.ProtocolIEIDAMFTNLAssociationToRemoveList:
			ngapLog.Traceln("[NGAP] Decode IE AMFTNLAssociationToRemoveList")
			aMFTNLAssociationToRemoveList = ie.Value.AMFTNLAssociationToRemoveList
		case ngapType.ProtocolIEIDAMFTNLAssociationToUpdateList:
			ngapLog.Traceln("[NGAP] Decode IE AMFTNLAssociationToUpdateList")
			aMFTNLAssociationToUpdateList = ie.Value.AMFTNLAssociationToUpdateList
		}
	}

	if aMFName != nil {
		amf.AMFName = aMFName
	}
	if servedGUAMIList != nil {
		amf.ServedGUAMIList = servedGUAMIList
	}

	if relativeAMFCapacity != nil {
		amf.RelativeAMFCapacity = relativeAMFCapacity
	}

	if pLMNSupportList != nil {
		amf.PLMNSupportList = pLMNSupportList
	}

	successList := []ngapType.AMFTNLAssociationSetupItem{}
	if aMFTNLAssociationToAddList != nil {
		// TODO: Establish TNL Association with AMF
		for _, item := range aMFTNLAssociationToAddList.List {
			tnlItem := amf.AddAMFTNLAssociationItem(item.AMFTNLAssociationAddress)
			tnlItem.TNLAddressWeightFactor = &item.TNLAddressWeightFactor.Value
			if item.TNLAssociationUsage != nil {
				tnlItem.TNLAssociationUsage = item.TNLAssociationUsage
			}
			setupItem := ngapType.AMFTNLAssociationSetupItem{
				AMFTNLAssociationAddress: item.AMFTNLAssociationAddress,
			}
			successList = append(successList, setupItem)
		}
	}
	if aMFTNLAssociationToRemoveList != nil {
		// TODO: Remove TNL Association with AMF
		for _, item := range aMFTNLAssociationToRemoveList.List {
			amf.DeleteAMFTNLAssociationItem(item.AMFTNLAssociationAddress)
		}
	}
	if aMFTNLAssociationToUpdateList != nil {
		// TODO: Update TNL Association with AMF
		for _, item := range aMFTNLAssociationToUpdateList.List {
			tnlItem := amf.FindAMFTNLAssociationItem(item.AMFTNLAssociationAddress)
			if tnlItem == nil {
				continue
			}
			if item.TNLAddressWeightFactor != nil {
				tnlItem.TNLAddressWeightFactor = &item.TNLAddressWeightFactor.Value
			}
			if item.TNLAssociationUsage != nil {
				tnlItem.TNLAssociationUsage = item.TNLAssociationUsage
			}
		}
	}

	var setupList *ngapType.AMFTNLAssociationSetupList
	if len(successList) > 0 {
		setupList = &ngapType.AMFTNLAssociationSetupList{
			List: successList,
		}
	}
	ngap_message.SendAMFConfigurationUpdateAcknowledge(amf, setupList, nil, nil)
}

func HandleRANConfigurationUpdateAcknowledge(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle RAN Configuration Update Acknowledge")

	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ngapLog.Error("SuccessfulOutcome is nil")
		return
	}

	rANConfigurationUpdateAcknowledge := successfulOutcome.Value.RANConfigurationUpdateAcknowledge
	if rANConfigurationUpdateAcknowledge == nil {
		ngapLog.Error("rANConfigurationUpdateAcknowledge is nil")
		return
	}

	for _, ie := range rANConfigurationUpdateAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			ngapLog.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

}

func HandleRANConfigurationUpdateFailure(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle RAN Configuration Update Failure")

	var cause *ngapType.Cause
	var timeToWait *ngapType.TimeToWait
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		ngapLog.Error("UnsuccessfulOutcome is nil")
		return
	}

	rANConfigurationUpdateFailure := unsuccessfulOutcome.Value.RANConfigurationUpdateFailure
	if rANConfigurationUpdateFailure == nil {
		ngapLog.Error("rANConfigurationUpdateFailure is nil")
		return
	}

	for _, ie := range rANConfigurationUpdateFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			ngapLog.Traceln("[NGAP] Decode IE Cause")
			cause = ie.Value.Cause
		case ngapType.ProtocolIEIDTimeToWait:
			ngapLog.Traceln("[NGAP] Decode IE TimeToWait")
			timeToWait = ie.Value.TimeToWait
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			ngapLog.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if cause != nil {
		printAndGetCause(cause)
	}

	printCriticalityDiagnostics(criticalityDiagnostics)

	var waitingTime int

	if timeToWait != nil {

		switch timeToWait.Value {
		case ngapType.TimeToWaitPresentV1s:
			waitingTime = 1
		case ngapType.TimeToWaitPresentV2s:
			waitingTime = 2
		case ngapType.TimeToWaitPresentV5s:
			waitingTime = 5
		case ngapType.TimeToWaitPresentV10s:
			waitingTime = 10
		case ngapType.TimeToWaitPresentV20s:
			waitingTime = 20
		case ngapType.TimeToWaitPresentV60s:
			waitingTime = 60
		}

	}

	if waitingTime != 0 {
		ngapLog.Infof("Wait at lease  %ds to resend RAN Configuration Update to same AMF[%s]", waitingTime, amf.SCTPAddr)
		n3iwf_context.N3IWFSelf().AMFReInitAvailableList[amf.SCTPAddr] = false
		handlerMessage := n3iwf_message.HandlerMessage{
			Event:    n3iwf_message.EventTimerSendRanConfigUpdateMessage,
			SCTPAddr: amf.SCTPAddr,
		}
		timer.StartTimer(waitingTime, func(msg interface{}) {
			n3iwf_message.SendMessage(msg.(n3iwf_message.HandlerMessage))
		}, handlerMessage)
		return
	}
	ngap_message.SendRANConfigurationUpdate(amf)
}

func HandleDownlinkRANConfigurationTransfer(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Downlink RAN Configuration Transfer")
}

func HandleDownlinkRANStatusTransfer(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Downlink RAN Status Transfer")
}

func HandleAMFStatusIndication(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle AMF Status Indication")
}

func HandleLocationReportingControl(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Location Reporting Control")
}

func HandleUETNLAReleaseRequest(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle UE TNLA Release Request")
}

func HandleOverloadStart(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Overload Start")

	var aMFOverloadResponse *ngapType.OverloadResponse
	var aMFTrafficLoadReductionIndication *ngapType.TrafficLoadReductionIndication
	var overloadStartNSSAIList *ngapType.OverloadStartNSSAIList

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("InitiatingMessage is nil")
		return
	}

	overloadStart := initiatingMessage.Value.OverloadStart
	if overloadStart == nil {
		ngapLog.Error("overloadStart is nil")
		return
	}

	for _, ie := range overloadStart.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFOverloadResponse:
			ngapLog.Traceln("[NGAP] Decode IE AMFOverloadResponse")
			aMFOverloadResponse = ie.Value.AMFOverloadResponse
		case ngapType.ProtocolIEIDAMFTrafficLoadReductionIndication:
			ngapLog.Traceln("[NGAP] Decode IE AMFTrafficLoadReductionIndication")
			aMFTrafficLoadReductionIndication = ie.Value.AMFTrafficLoadReductionIndication
		case ngapType.ProtocolIEIDOverloadStartNSSAIList:
			ngapLog.Traceln("[NGAP] Decode IE OverloadStartNSSAIList")
			overloadStartNSSAIList = ie.Value.OverloadStartNSSAIList
		}
	}
	// TODO: restrict rule about overload action
	amf.StartOverload(aMFOverloadResponse, aMFTrafficLoadReductionIndication, overloadStartNSSAIList)

}

func HandleOverloadStop(amf *n3iwf_context.N3IWFAMF, message *ngapType.NGAPPDU) {

	ngapLog.Infoln("[N3IWF] Handle Overload Stop")

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}
	// TODO: remove restrict about overload action
	amf.StopOverload()

}

func buildCriticalityDiagnostics(
	procedureCode *int64,
	triggeringMessage *aper.Enumerated,
	procedureCriticality *aper.Enumerated,
	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (criticalityDiagnostics ngapType.CriticalityDiagnostics) {

	if procedureCode != nil {
		criticalityDiagnostics.ProcedureCode = new(ngapType.ProcedureCode)
		criticalityDiagnostics.ProcedureCode.Value = *procedureCode
	}

	if triggeringMessage != nil {
		criticalityDiagnostics.TriggeringMessage = new(ngapType.TriggeringMessage)
		criticalityDiagnostics.TriggeringMessage.Value = *triggeringMessage
	}

	if procedureCriticality != nil {
		criticalityDiagnostics.ProcedureCriticality = new(ngapType.Criticality)
		criticalityDiagnostics.ProcedureCriticality.Value = *procedureCriticality
	}

	if iesCriticalityDiagnostics != nil {
		criticalityDiagnostics.IEsCriticalityDiagnostics = iesCriticalityDiagnostics
	}

	return criticalityDiagnostics
}

func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (item ngapType.CriticalityDiagnosticsIEItem) {

	item = ngapType.CriticalityDiagnosticsIEItem{
		IECriticality: ngapType.Criticality{
			Value: ieCriticality,
		},
		IEID: ngapType.ProtocolIEID{
			Value: ieID,
		},
		TypeOfError: ngapType.TypeOfError{
			Value: typeOfErr,
		},
	}

	return item
}

func buildCause(present int, value aper.Enumerated) (cause *ngapType.Cause) {
	cause = new(ngapType.Cause)
	cause.Present = present

	switch present {
	case ngapType.CausePresentRadioNetwork:
		cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
		cause.RadioNetwork.Value = value
	case ngapType.CausePresentTransport:
		cause.Transport = new(ngapType.CauseTransport)
		cause.Transport.Value = value
	case ngapType.CausePresentNas:
		cause.Nas = new(ngapType.CauseNas)
		cause.Nas.Value = value
	case ngapType.CausePresentProtocol:
		cause.Protocol = new(ngapType.CauseProtocol)
		cause.Protocol.Value = value
	case ngapType.CausePresentMisc:
		cause.Misc = new(ngapType.CauseMisc)
		cause.Misc.Value = value
	case ngapType.CausePresentNothing:
	}

	return
}

func printAndGetCause(cause *ngapType.Cause) (present int, value aper.Enumerated) {

	present = cause.Present
	switch cause.Present {
	case ngapType.CausePresentRadioNetwork:
		ngapLog.Warnf("Cause RadioNetwork[%d]", cause.RadioNetwork.Value)
		value = cause.RadioNetwork.Value
	case ngapType.CausePresentTransport:
		ngapLog.Warnf("Cause Transport[%d]", cause.Transport.Value)
		value = cause.Transport.Value
	case ngapType.CausePresentProtocol:
		ngapLog.Warnf("Cause Protocol[%d]", cause.Protocol.Value)
		value = cause.Protocol.Value
	case ngapType.CausePresentNas:
		ngapLog.Warnf("Cause Nas[%d]", cause.Nas.Value)
		value = cause.Nas.Value
	case ngapType.CausePresentMisc:
		ngapLog.Warnf("Cause Misc[%d]", cause.Misc.Value)
		value = cause.Misc.Value
	default:
		ngapLog.Errorf("Invalid Cause group[%d]", cause.Present)
	}
	return
}

func printCriticalityDiagnostics(criticalityDiagnostics *ngapType.CriticalityDiagnostics) {
	if criticalityDiagnostics == nil {
		return
	} else {
		iesCriticalityDiagnostics := criticalityDiagnostics.IEsCriticalityDiagnostics
		if iesCriticalityDiagnostics != nil {
			for index, item := range iesCriticalityDiagnostics.List {
				ngapLog.Warnf("Criticality IE item %d:", index+1)
				ngapLog.Warnf("IE ID: %d", item.IEID.Value)

				switch item.IECriticality.Value {
				case ngapType.CriticalityPresentReject:
					ngapLog.Warn("IE Criticality: Reject")
				case ngapType.CriticalityPresentIgnore:
					ngapLog.Warn("IE Criticality: Ignore")
				case ngapType.CriticalityPresentNotify:
					ngapLog.Warn("IE Criticality: Notify")
				}

				switch item.TypeOfError.Value {
				case ngapType.TypeOfErrorPresentNotUnderstood:
					ngapLog.Warn("Type of error: Not Understood")
				case ngapType.TypeOfErrorPresentMissing:
					ngapLog.Warn("Type of error: Missing")
				}
			}
		} else {
			ngapLog.Error("IEsCriticalityDiagnostics is nil")
		}
		return
	}
}

func getPDUSessionResourceReleaseResponseTransfer() []byte {
	data := ngapType.PDUSessionResourceReleaseResponseTransfer{}
	encodeData, _ := aper.MarshalWithParams(data, "valueExt")
	return encodeData
}
