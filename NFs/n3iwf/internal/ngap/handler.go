package ngap

import (
	"encoding/binary"
	"math"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-gtp/gtpv1"

	"github.com/free5gc/aper"
	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/n3iwf/internal/nas/nas_security"
	"github.com/free5gc/n3iwf/internal/ngap/message"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/sctp"
)

func (s *Server) HandleNGSetupResponse(
	sctpAddr string,
	conn *sctp.SCTPConn,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle NG Setup Response")

	var amfName *ngapType.AMFName
	var servedGUAMIList *ngapType.ServedGUAMIList
	var relativeAMFCapacity *ngapType.RelativeAMFCapacity
	var plmnSupportList *ngapType.PLMNSupportList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	n3iwfCtx := s.Context()

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	successfulOutcome := pdu.SuccessfulOutcome
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDServedGUAMIList:
			ngapLog.Traceln("[NGAP] Decode IE ServedGUAMIList")
			servedGUAMIList = ie.Value.ServedGUAMIList
			if servedGUAMIList == nil {
				ngapLog.Errorf("ServedGUAMIList is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			ngapLog.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if len(iesCriticalityDiagnostics.List) != 0 {
		ngapLog.Traceln("[NGAP] Sending error indication to AMF, because some mandatory IEs were not included")

		cause := message.BuildCause(ngapType.CausePresentProtocol,
			ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)

		procedureCode := ngapType.ProcedureCodeNGSetup
		triggeringMessage := ngapType.TriggeringMessagePresentSuccessfulOutcome
		procedureCriticality := ngapType.CriticalityPresentReject

		criticalityDiagnostics := buildCriticalityDiagnostics(
			&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)

		message.SendErrorIndicationWithSctpConn(conn, nil, nil, cause, &criticalityDiagnostics)

		return
	}

	amfInfo := n3iwfCtx.NewN3iwfAmf(sctpAddr, conn)

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

func (s *Server) HandleNGSetupFailure(
	sctpAddr string,
	conn *sctp.SCTPConn,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle NG Setup Failure")

	var cause *ngapType.Cause
	var timeToWait *ngapType.TimeToWait
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	n3iwfCtx := s.Context()

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		ngapLog.Error("Unsuccessful Message is nil")
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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

		cause = message.BuildCause(
			ngapType.CausePresentProtocol,
			ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)

		procedureCode := ngapType.ProcedureCodeNGSetup
		triggeringMessage := ngapType.TriggeringMessagePresentUnsuccessfullOutcome
		procedureCriticality := ngapType.CriticalityPresentReject

		criticalityDiagnostics := buildCriticalityDiagnostics(
			&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)

		message.SendErrorIndicationWithSctpConn(conn, nil, nil, cause, &criticalityDiagnostics)
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
		n3iwfCtx.AMFReInitAvailableListStore(sctpAddr, false)
		time.AfterFunc(time.Duration(waitingTime)*time.Second, func() {
			n3iwfCtx.AMFReInitAvailableListStore(sctpAddr, true)
			message.SendNGSetupRequest(conn, n3iwfCtx)
		})
		return
	}
}

func (s *Server) HandleNGReset(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle NG Reset")

	var cause *ngapType.Cause
	var resetType *ngapType.ResetType
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	n3iwfCtx := s.Context()

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodeNGReset
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(
			&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		message.SendErrorIndication(amf, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	printAndGetCause(cause)

	switch resetType.Present {
	case ngapType.ResetTypePresentNGInterface:
		ngapLog.Trace("ResetType Present: NG Interface")
		// TODO: Release Uu Interface related to this amf(IPSec)
		// Remove all Ue
		if err := amf.RemoveAllRelatedUe(); err != nil {
			ngapLog.Errorf("RemoveAllRelatedUe error : %v", err)
		}
		message.SendNGResetAcknowledge(amf, nil, nil)
	case ngapType.ResetTypePresentPartOfNGInterface:
		ngapLog.Trace("ResetType Present: Part of NG Interface")

		partOfNGInterface := resetType.PartOfNGInterface
		if partOfNGInterface == nil {
			ngapLog.Error("PartOfNGInterface is nil")
			return
		}

		var ranUe n3iwf_context.RanUe

		for _, ueAssociatedLogicalNGConnectionItem := range partOfNGInterface.List {
			if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
				ngapLog.Tracef("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				ranUe, _ = n3iwfCtx.RanUePoolLoad(ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
			} else if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
				ngapLog.Tracef("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				ranUe = amf.FindUeByAmfUeNgapID(ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
			}

			if ranUe == nil {
				ngapLog.Warn("Cannot not find RanUE Context")
				if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
					ngapLog.Warnf("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				}
				if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
					ngapLog.Warnf("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				}
				continue
			}
			// TODO: Release Uu Interface (IPSec)
			if err := ranUe.Remove(); err != nil {
				ngapLog.Errorf("Remove RanUE Context error : %v", err)
			}
		}
		message.SendNGResetAcknowledge(amf, partOfNGInterface, nil)
	default:
		ngapLog.Warnf("Invalid ResetType[%d]", resetType.Present)
	}
}

func (s *Server) HandleNGResetAcknowledge(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle NG Reset Acknowledge")

	var uEAssociatedLogicalNGConnectionList *ngapType.UEAssociatedLogicalNGConnectionList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	successfulOutcome := pdu.SuccessfulOutcome
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
		ngapLog.Tracef("%d RanUE association(s) has been reset", len(uEAssociatedLogicalNGConnectionList.List))
		for i, item := range uEAssociatedLogicalNGConnectionList.List {
			if item.AMFUENGAPID != nil && item.RANUENGAPID != nil {
				ngapLog.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[%d]",
					i+1, item.AMFUENGAPID.Value, item.RANUENGAPID.Value)
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

func (s *Server) HandleInitialContextSetupRequest(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle Initial Context Setup Request")

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
	// var nasPDU *ngapType.NASPDU
	var emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var ranUe n3iwf_context.RanUe
	var ranUeCtx *n3iwf_context.RanUeSharedCtx

	n3iwfCtx := s.Context()

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Errorf("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDUESecurityCapabilities:
			ngapLog.Traceln("[NGAP] Decode IE UESecurityCapabilities")
			ueSecurityCapabilities = ie.Value.UESecurityCapabilities
			if ueSecurityCapabilities == nil {
				ngapLog.Errorf("UESecurityCapabilities is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDSecurityKey:
			ngapLog.Traceln("[NGAP] Decode IE SecurityKey")
			securityKey = ie.Value.SecurityKey
			if securityKey == nil {
				ngapLog.Errorf("SecurityKey is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
			// nasPDU = ie.Value.NASPDU
		case ngapType.ProtocolIEIDEmergencyFallbackIndicator:
			ngapLog.Traceln("[NGAP] Decode IE EmergencyFallbackIndicator")
			emergencyFallbackIndicator = ie.Value.EmergencyFallbackIndicator
			if emergencyFallbackIndicator != nil {
				ngapLog.Warnln("Not Supported IE [EmergencyFallbackIndicator]")
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		ngapLog.Traceln(
			"[NGAP] Sending unsuccessful outcome to AMF, because some mandatory IEs were not included")
		cause := message.BuildCause(ngapType.CausePresentProtocol,
			ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage)

		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)

		failedListCxtFail := new(ngapType.PDUSessionResourceFailedToSetupListCxtFail)
		for _, item := range pduSessionResourceSetupListCxtReq.List {
			transfer, err := message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
			if err != nil {
				ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", err)
			}
			message.AppendPDUSessionResourceFailedToSetupListCxtfail(
				failedListCxtFail, item.PDUSessionID.Value, transfer)
		}

		message.SendInitialContextSetupFailure(ranUe, *cause, failedListCxtFail, &criticalityDiagnostics)
		return
	}

	if (amfUeNgapID != nil) && (ranUeNgapID != nil) {
		// Find UE context
		var ok bool
		ranUe, ok = n3iwfCtx.RanUePoolLoad(ranUeNgapID.Value)
		if !ok {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
			// TODO: build cause and handle error
			// Cause: Unknown local UE NGAP ID
			return
		}
		ranUeCtx = ranUe.GetSharedCtx()
		if ranUeCtx.AmfUeNgapId != amfUeNgapID.Value {
			// TODO: build cause and handle error
			// Cause: Inconsistent remote UE NGAP ID
			return
		}
	}

	if ranUe == nil {
		ngapLog.Errorf("RAN UE context is nil")
		return
	}

	ranUeCtx.AmfUeNgapId = amfUeNgapID.Value
	ranUeCtx.RanUeNgapId = ranUeNgapID.Value

	if pduSessionResourceSetupListCxtReq != nil {
		if ueAggregateMaximumBitRate != nil {
			ranUeCtx.Ambr = ueAggregateMaximumBitRate
		} else {
			ngapLog.Errorln("IE[UEAggregateMaximumBitRate] is nil")
			cause := message.BuildCause(ngapType.CausePresentProtocol,
				ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage)

			criticalityDiagnosticsIEItem := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
				ngapType.ProtocolIEIDUEAggregateMaximumBitRate, ngapType.TypeOfErrorPresentMissing)
			iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, criticalityDiagnosticsIEItem)
			criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)

			failedListCxtFail := new(ngapType.PDUSessionResourceFailedToSetupListCxtFail)
			for _, item := range pduSessionResourceSetupListCxtReq.List {
				transfer, err := message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
				if err != nil {
					ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", err)
				}
				message.AppendPDUSessionResourceFailedToSetupListCxtfail(
					failedListCxtFail, item.PDUSessionID.Value, transfer)
			}

			message.SendInitialContextSetupFailure(ranUe, *cause, failedListCxtFail, &criticalityDiagnostics)
			return
		}

		setupListCxtRes := new(ngapType.PDUSessionResourceSetupListCxtRes)
		failedListCxtRes := new(ngapType.PDUSessionResourceFailedToSetupListCxtRes)

		// UE temporary data for PDU session setup response
		ranUeCtx.TemporaryPDUSessionSetupData.SetupListCxtRes = setupListCxtRes
		ranUeCtx.TemporaryPDUSessionSetupData.FailedListCxtRes = failedListCxtRes
		ranUeCtx.TemporaryPDUSessionSetupData.Index = 0
		ranUeCtx.TemporaryPDUSessionSetupData.UnactivatedPDUSession = nil
		ranUeCtx.TemporaryPDUSessionSetupData.NGAPProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup

		for _, item := range pduSessionResourceSetupListCxtReq.List {
			pduSessionID := item.PDUSessionID.Value
			// TODO: send NAS to UE
			// pduSessionNasPdu := item.NASPDU
			snssai := item.SNSSAI

			transfer := ngapType.PDUSessionResourceSetupRequestTransfer{}
			err := aper.UnmarshalWithParams(item.PDUSessionResourceSetupRequestTransfer, &transfer, "valueExt")
			if err != nil {
				ngapLog.Errorf("[PDUSessionID: %d] PDUSessionResourceSetupRequestTransfer Decode Error: %v\n",
					pduSessionID, err)
			}

			pduSession, err := ranUeCtx.CreatePDUSession(pduSessionID, snssai)
			if err != nil {
				ngapLog.Errorf("Create PDU Session Error: %v\n", err)

				cause := message.BuildCause(ngapType.CausePresentRadioNetwork,
					ngapType.CauseRadioNetworkPresentMultiplePDUSessionIDInstances)
				unsuccessfulTransfer, buildErr := message.
					BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
				if buildErr != nil {
					ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", buildErr)
				}
				message.AppendPDUSessionResourceFailedToSetupListCxtRes(
					failedListCxtRes, pduSessionID, unsuccessfulTransfer)
				continue
			}

			success, resTransfer := s.handlePDUSessionResourceSetupRequestTransfer(ranUe, pduSession, transfer)
			if success {
				// Append this PDU session to unactivated PDU session list
				ranUeCtx.TemporaryPDUSessionSetupData.UnactivatedPDUSession = append(
					ranUeCtx.TemporaryPDUSessionSetupData.UnactivatedPDUSession,
					pduSession)
			} else {
				// Delete the pdusession store in UE conext
				delete(ranUeCtx.PduSessionList, pduSessionID)
				message.
					AppendPDUSessionResourceFailedToSetupListCxtRes(failedListCxtRes, pduSessionID, resTransfer)
			}
		}
	}

	if oldAMF != nil {
		ngapLog.Debugf("Old AMF: %s\n", oldAMF.Value)
	}

	if guami != nil {
		ranUeCtx.Guami = guami
	}

	if allowedNSSAI != nil {
		ranUeCtx.AllowedNssai = allowedNSSAI
	}

	if maskedIMEISV != nil {
		ranUeCtx.MaskedIMEISV = maskedIMEISV
	}

	if ueRadioCapability != nil {
		ranUeCtx.RadioCapability = ueRadioCapability
	}

	if coreNetworkAssistanceInformation != nil {
		ranUeCtx.CoreNetworkAssistanceInformation = coreNetworkAssistanceInformation
	}

	if indexToRFSP != nil {
		ranUeCtx.IndexToRfsp = indexToRFSP.Value
	}

	if ueSecurityCapabilities != nil {
		ranUeCtx.SecurityCapabilities = ueSecurityCapabilities
	}

	// Send EAP Success to UE
	switch ue := ranUe.(type) {
	case *n3iwf_context.N3IWFRanUe:
		spi, ok := n3iwfCtx.IkeSpiLoad(ranUeCtx.RanUeNgapId)
		if !ok {
			ngapLog.Errorf("Cannot get spi from ngapid : %+v", ranUeCtx.RanUeNgapId)
			return
		}

		s.SendIkeEvt(n3iwf_context.NewSendEAPSuccessMsgEvt(
			spi, securityKey.Value.Bytes, len(ranUeCtx.PduSessionList),
		))
	default:
		ngapLog.Errorf("Unknown UE type: %T", ue)
	}
}

// handlePDUSessionResourceSetupRequestTransfer parse and store needed information from NGAP
// and setup user plane connection for UE
// Parameters:
// UE context :: a pointer to the UE's pdusession data structure ::
// SMF PDU session resource setup request transfer
// Return value:
// a status value indicate whether the handlling is "success" ::
// if failed, an unsuccessfulTransfer is set, otherwise, set to nil
func (s *Server) handlePDUSessionResourceSetupRequestTransfer(
	ranUe n3iwf_context.RanUe,
	pduSession *n3iwf_context.PDUSession,
	transfer ngapType.PDUSessionResourceSetupRequestTransfer,
) (bool, []byte) {
	var pduSessionAMBR *ngapType.PDUSessionAggregateMaximumBitRate
	var ulNGUUPTNLInformation *ngapType.UPTransportLayerInformation
	var pduSessionType *ngapType.PDUSessionType
	var securityIndication *ngapType.SecurityIndication
	var networkInstance *ngapType.NetworkInstance
	var qosFlowSetupRequestList *ngapType.QosFlowSetupRequestList
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	n3iwfCtx := s.Context()

	for _, ie := range transfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDPDUSessionAggregateMaximumBitRate:
			pduSessionAMBR = ie.Value.PDUSessionAggregateMaximumBitRate
		case ngapType.ProtocolIEIDULNGUUPTNLInformation:
			ulNGUUPTNLInformation = ie.Value.ULNGUUPTNLInformation
			if ulNGUUPTNLInformation == nil {
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionType:
			pduSessionType = ie.Value.PDUSessionType
			if pduSessionType == nil {
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDSecurityIndication:
			securityIndication = ie.Value.SecurityIndication
		case ngapType.ProtocolIEIDNetworkInstance:
			networkInstance = ie.Value.NetworkInstance
		case ngapType.ProtocolIEIDQosFlowSetupRequestList:
			qosFlowSetupRequestList = ie.Value.QosFlowSetupRequestList
			if qosFlowSetupRequestList == nil {
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	ngapLog := logger.NgapLog

	if len(iesCriticalityDiagnostics.List) > 0 {
		cause := message.BuildCause(ngapType.CausePresentProtocol,
			ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage)
		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)
		responseTransfer, err := message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(
			*cause, &criticalityDiagnostics)
		if err != nil {
			ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", err)
		}
		return false, responseTransfer
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
			cause := message.BuildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentSemanticError)
			responseTransfer, err := message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
			if err != nil {
				ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", err)
			}
			return false, responseTransfer
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
			cause := message.BuildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentSemanticError)
			responseTransfer, err := message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
			if err != nil {
				ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", err)
			}
			return false, responseTransfer
		}
	} else {
		pduSession.SecurityIntegrity = true
		pduSession.SecurityCipher = true
	}

	// TODO: apply qos rule
	for _, item := range qosFlowSetupRequestList.List {
		// QoS Flow
		qosFlow := new(n3iwf_context.QosFlow)
		qosFlow.Identifier = item.QosFlowIdentifier.Value
		qosFlow.Parameters = item.QosFlowLevelQosParameters
		pduSession.QosFlows[item.QosFlowIdentifier.Value] = qosFlow

		value := item.QosFlowIdentifier.Value
		if value < 0 || value > math.MaxUint8 {
			ngapLog.Errorf("handlePDUSessionResourceSetupRequestTransfer() "+
				"item.QosFlowIdentifier.Value exceeds uint8 range: %d", value)
			return false, nil
		}
		// QFI List
		pduSession.QFIList = append(pduSession.QFIList, uint8(value))
	}

	// Setup GTP tunnel with UPF
	// TODO: Support IPv6
	upfIPv4, _ := ngapConvert.IPAddressToString(ulNGUUPTNLInformation.GTPTunnel.TransportLayerAddress)
	if upfIPv4 != "" {
		gtpConnInfo := &n3iwf_context.GTPConnectionInfo{
			UPFIPAddr:    upfIPv4,
			OutgoingTEID: binary.BigEndian.Uint32(ulNGUUPTNLInformation.GTPTunnel.GTPTEID.Value),
		}

		// UPF UDP address
		upfAddr := upfIPv4 + gtpv1.GTPUPort
		upfUDPAddr, err := net.ResolveUDPAddr("udp", upfAddr)
		if err != nil {
			var responseTransfer []byte

			ngapLog.Errorf("Resolve UPF addr [%s] failed: %v", upfAddr, err)
			cause := message.BuildCause(ngapType.CausePresentTransport,
				ngapType.CauseTransportPresentTransportResourceUnavailable)
			responseTransfer, err = message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
			if err != nil {
				ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", err)
			}
			return false, responseTransfer
		}

		// UE TEID
		ueTEID := n3iwfCtx.NewTEID(ranUe)
		if ueTEID == 0 {
			var responseTransfer []byte

			ngapLog.Error("Invalid TEID (0).")
			cause := message.BuildCause(
				ngapType.CausePresentProtocol,
				ngapType.CauseProtocolPresentUnspecified)
			responseTransfer, err = message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
			if err != nil {
				ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", err)
			}
			return false, responseTransfer
		}

		// Setup GTP connection with UPF
		gtpConnInfo.UPFUDPAddr = upfUDPAddr
		gtpConnInfo.IncomingTEID = ueTEID

		pduSession.GTPConnInfo = gtpConnInfo
	} else {
		ngapLog.Error(
			"Cannot parse \"PDU session resource setup request transfer\" message \"UL NG-U UP TNL Information\"")
		cause := message.BuildCause(ngapType.CausePresentProtocol,
			ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)
		responseTransfer, err := message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
		if err != nil {
			ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", err)
		}
		return false, responseTransfer
	}

	return true, nil
}

func (s *Server) HandleUEContextModificationRequest(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle UE Context Modification Request")

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

	var ranUe n3iwf_context.RanUe
	var ranUeCtx *n3iwf_context.RanUeSharedCtx

	n3iwfCtx := s.Context()

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Errorf("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
		var ok bool
		ranUe, ok = n3iwfCtx.RanUePoolLoad(ranUeNgapID.Value)
		if !ok {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
			// TODO: build cause and handle error
			// Cause: Unknown local UE NGAP ID
			return
		}
		ranUeCtx = ranUe.GetSharedCtx()
		if ranUeCtx.AmfUeNgapId != amfUeNgapID.Value {
			// TODO: build cause and handle error
			// Cause: Inconsistent remote UE NGAP ID
			return
		}
	}

	if newAmfUeNgapID != nil {
		ngapLog.Debugf("New AmfUeNgapID[%d]\n", newAmfUeNgapID.Value)
		ranUeCtx.AmfUeNgapId = newAmfUeNgapID.Value
	}

	if ueAggregateMaximumBitRate != nil {
		ranUeCtx.Ambr = ueAggregateMaximumBitRate
		// TODO: use the received UE Aggregate Maximum Bit Rate for all non-GBR QoS flows
	}

	if ueSecurityCapabilities != nil {
		ranUeCtx.SecurityCapabilities = ueSecurityCapabilities
	}

	// TODO: use new security key to update security context

	if indexToRFSP != nil {
		ranUeCtx.IndexToRfsp = indexToRFSP.Value
	}

	message.SendUEContextModificationResponse(ranUe, nil)

	spi, ok := n3iwfCtx.IkeSpiLoad(ranUeCtx.RanUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get spi from ngapid : %+v", ranUeCtx.RanUeNgapId)
		return
	}

	s.SendIkeEvt(n3iwf_context.NewIKEContextUpdateEvt(spi, securityKey.Value.Bytes)) // Kn3iwf
}

func (s *Server) HandleUEContextReleaseCommand(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle UE Context Release Command")

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}

	var ueNgapIDs *ngapType.UENGAPIDs
	var cause *ngapType.Cause
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
	var ranUe n3iwf_context.RanUe

	n3iwfCtx := s.Context()

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
		var ok bool
		ranUe, ok = n3iwfCtx.RanUePoolLoad(ueNgapIDs.UENGAPIDPair.RANUENGAPID.Value)
		if !ok {
			ranUe = amf.FindUeByAmfUeNgapID(ueNgapIDs.UENGAPIDPair.AMFUENGAPID.Value)
		}
	case ngapType.UENGAPIDsPresentAMFUENGAPID:
		// TODO: find UE according to specific AMF
		// The implementation here may have error when N3IWF need to
		// connect multiple AMFs.
		// Use UEpool in AMF context can solve this problem
		ranUe = amf.FindUeByAmfUeNgapID(ueNgapIDs.AMFUENGAPID.Value)
	}

	if ranUe == nil {
		// TODO: send error indication(unknown local ngap ue id)
		return
	}

	if cause != nil {
		printAndGetCause(cause)
	}

	ranUe.GetSharedCtx().UeCtxRelState = n3iwf_context.UeCtxRelStateOngoing

	message.SendUEContextReleaseComplete(ranUe, nil)

	err := s.releaseIkeUeAndRanUe(ranUe)
	if err != nil {
		ngapLog.Warnf("HandleUEContextReleaseCommand(): %v", err)
	}
}

func (s *Server) releaseIkeUeAndRanUe(ranUe n3iwf_context.RanUe) error {
	n3iwfCtx := s.Context()
	ranUeNgapID := ranUe.GetSharedCtx().RanUeNgapId

	localSPI, ok := n3iwfCtx.IkeSpiLoad(ranUeNgapID)
	if ok {
		s.SendIkeEvt(n3iwf_context.NewIKEDeleteRequestEvt(localSPI))
	}

	if err := ranUe.Remove(); err != nil {
		return errors.Wrapf(err, "releaseIkeUeAndRanUe RanUeNgapId[%016x]", ranUeNgapID)
	}
	return nil
}

func (s *Server) HandleDownlinkNASTransport(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle Downlink NAS Transport")

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
	var ranUe n3iwf_context.RanUe

	n3iwfCtx := s.Context()

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Errorf("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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

	// if len(iesCriticalityDiagnostics.List) > 0 {
	// TODO: Send Error Indication
	// }

	if ranUeNgapID != nil {
		var ok bool
		ranUe, ok = n3iwfCtx.RanUePoolLoad(ranUeNgapID.Value)
		if !ok {
			ngapLog.Warnf("No UE Context[RanUeNgapID:%d]\n", ranUeNgapID.Value)
			return
		}
	}
	ranUeCtx := ranUe.GetSharedCtx()

	if amfUeNgapID != nil {
		if ranUeCtx.AmfUeNgapId == n3iwf_context.AmfUeNgapIdUnspecified {
			ngapLog.Tracef("Create new logical UE-associated NG-connection")
			ranUeCtx.AmfUeNgapId = amfUeNgapID.Value
		} else {
			if ranUeCtx.AmfUeNgapId != amfUeNgapID.Value {
				ngapLog.Warn("AMFUENGAPID unmatched")
				return
			}
		}
	}

	if oldAMF != nil {
		ngapLog.Debugf("Old AMF: %s\n", oldAMF.Value)
	}

	if indexToRFSP != nil {
		ranUeCtx.IndexToRfsp = indexToRFSP.Value
	}

	if ueAggregateMaximumBitRate != nil {
		ranUeCtx.Ambr = ueAggregateMaximumBitRate
	}

	if allowedNSSAI != nil {
		ranUeCtx.AllowedNssai = allowedNSSAI
	}

	if nasPDU != nil {
		switch ue := ranUe.(type) {
		case *n3iwf_context.N3IWFRanUe:
			// Send EAP5G NAS to UE
			spi, ok := n3iwfCtx.IkeSpiLoad(ue.RanUeNgapId)
			if !ok {
				ngapLog.Errorf("Cannot get SPI from RanUeNGAPId : %+v", ue.RanUeNgapId)
				return
			}

			if !ue.IsNASTCPConnEstablished {
				s.SendIkeEvt(n3iwf_context.NewSendEAPNASMsgEvt(spi, []byte(nasPDU.Value)))
			} else {
				// Using a "NAS message envelope" to transport a NAS message
				// over the non-3GPP access between the UE and the N3IWF
				nasEnv := nas_security.EncapNasMsgToEnvelope([]byte(nasPDU.Value))

				if ue.IsNASTCPConnEstablishedComplete {
					// Send to UE
					if n, err := ue.TCPConnection.Write(nasEnv); err != nil {
						ngapLog.Errorf("Writing via IPSec signalling SA failed: %v", err)
					} else {
						ngapLog.Trace("Forward NWu <- N2")
						ngapLog.Tracef("Wrote %d bytes", n)
					}
				} else {
					ue.TemporaryCachedNASMessage = nasEnv
				}
			}
		default:
			ngapLog.Errorf("Unknown UE type: %T", ue)
		}
	}
}

func (s *Server) HandlePDUSessionResourceSetupRequest(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle PDU Session Resource Setup Request")

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}

	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	var nasPDU *ngapType.NASPDU
	var pduSessionResourceSetupListSUReq *ngapType.PDUSessionResourceSetupListSUReq
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
	var pduSessionEstablishmentAccept *ngapType.NASPDU
	var ranUe n3iwf_context.RanUe
	var ranUeCtx *n3iwf_context.RanUeSharedCtx

	n3iwfCtx := s.Context()

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Errorf("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
		var ok bool
		ranUe, ok = n3iwfCtx.RanUePoolLoad(ranUeNgapID.Value)
		if !ok {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
			// TODO: build cause and handle error
			// Cause: Unknown local UE NGAP ID
			return
		}
		ranUeCtx = ranUe.GetSharedCtx()
		if ranUeCtx.AmfUeNgapId != amfUeNgapID.Value {
			// TODO: build cause and handle error
			// Cause: Inconsistent remote UE NGAP ID
			return
		}
	}

	if nasPDU != nil {
		n3iwfUe, ok := ranUe.(*n3iwf_context.N3IWFRanUe)
		if !ok {
			ngapLog.Errorln("HandlePDUSessionResourceSetupRequest(): [Type Assertion] RanUe -> N3iwfRanUe failed")
			return
		}
		if n3iwfUe.TCPConnection == nil {
			ngapLog.Error("No IPSec NAS signalling SA for this UE")
			return
		}

		// Using a "NAS message envelope" to transport a NAS message
		// over the non-3GPP access between the UE and the N3IWF
		nasEnv := nas_security.EncapNasMsgToEnvelope([]byte(nasPDU.Value))

		n, err := n3iwfUe.TCPConnection.Write(nasEnv)
		if err != nil {
			ngapLog.Errorf("Send NAS to UE failed: %v", err)
			return
		}
		ngapLog.Tracef("Wrote %d bytes", n)
	}

	tempPDUSessionSetupData := ranUeCtx.TemporaryPDUSessionSetupData
	tempPDUSessionSetupData.NGAPProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup

	if pduSessionResourceSetupListSUReq != nil {
		setupListSURes := new(ngapType.PDUSessionResourceSetupListSURes)
		failedListSURes := new(ngapType.PDUSessionResourceFailedToSetupListSURes)

		tempPDUSessionSetupData.SetupListSURes = setupListSURes
		tempPDUSessionSetupData.FailedListSURes = failedListSURes
		tempPDUSessionSetupData.Index = 0
		tempPDUSessionSetupData.UnactivatedPDUSession = nil
		tempPDUSessionSetupData.NGAPProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup

		for _, item := range pduSessionResourceSetupListSUReq.List {
			pduSessionID := item.PDUSessionID.Value
			pduSessionEstablishmentAccept = item.PDUSessionNASPDU
			snssai := item.SNSSAI

			transfer := ngapType.PDUSessionResourceSetupRequestTransfer{}
			err := aper.UnmarshalWithParams(item.PDUSessionResourceSetupRequestTransfer, &transfer, "valueExt")
			if err != nil {
				ngapLog.Errorf("[PDUSessionID: %d] PDUSessionResourceSetupRequestTransfer Decode Error: %v\n",
					pduSessionID, err)
			}

			pduSession, err := ranUeCtx.CreatePDUSession(pduSessionID, snssai)
			if err != nil {
				ngapLog.Errorf("Create PDU Session Error: %v\n", err)

				cause := message.BuildCause(ngapType.CausePresentRadioNetwork,
					ngapType.CauseRadioNetworkPresentMultiplePDUSessionIDInstances)
				unsuccessfulTransfer, buildErr := message.
					BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
				if buildErr != nil {
					ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %v\n", buildErr)
				}
				message.AppendPDUSessionResourceFailedToSetupListSURes(
					failedListSURes, pduSessionID, unsuccessfulTransfer)
				continue
			}

			// Process the message for AN
			success, resTransfer := s.handlePDUSessionResourceSetupRequestTransfer(
				ranUe, pduSession, transfer)
			if success {
				// Append this PDU session to unactivated PDU session list
				tempPDUSessionSetupData.UnactivatedPDUSession = append(
					tempPDUSessionSetupData.UnactivatedPDUSession,
					pduSession)
			} else {
				// Delete the pdusession store in UE conext
				delete(ranUeCtx.PduSessionList, pduSessionID)
				message.AppendPDUSessionResourceFailedToSetupListSURes(
					failedListSURes, pduSessionID, resTransfer)
			}
		}
	}

	if tempPDUSessionSetupData != nil && len(tempPDUSessionSetupData.UnactivatedPDUSession) != 0 {
		switch ue := ranUe.(type) {
		case *n3iwf_context.N3IWFRanUe:
			spi, ok := n3iwfCtx.IkeSpiLoad(ue.RanUeNgapId)
			if !ok {
				ngapLog.Errorf("Cannot get SPI from ranNgapID : %+v", ranUeNgapID)
				return
			}

			s.SendIkeEvt(n3iwf_context.NewCreatePDUSessionEvt(spi,
				len(ue.PduSessionList),
				ue.TemporaryPDUSessionSetupData),
			)

			// TS 23.501 4.12.5 Requested PDU Session Establishment via Untrusted non-3GPP Access
			// After all IPsec Child SAs are established, the N3IWF shall forward to UE via the signalling IPsec SA
			// the PDU Session Establishment Accept message
			nasEnv := nas_security.EncapNasMsgToEnvelope([]byte(pduSessionEstablishmentAccept.Value))

			// Cache the pduSessionEstablishmentAccept and forward to the UE after all CREATE_CHILD_SAs finish
			ue.TemporaryCachedNASMessage = nasEnv
		}
	}
}

func (s *Server) HandlePDUSessionResourceModifyRequest(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle PDU Session Resource Modify Request")

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}

	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	var pduSessionResourceModifyListModReq *ngapType.PDUSessionResourceModifyListModReq
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
	var ranUe n3iwf_context.RanUe
	var ranUeCtx *n3iwf_context.RanUeSharedCtx

	n3iwfCtx := s.Context()

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				ngapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModReq:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionResourceModifyListModReq")
			pduSessionResourceModifyListModReq = ie.Value.PDUSessionResourceModifyListModReq
			if pduSessionResourceModifyListModReq == nil {
				ngapLog.Error("PDUSessionResourceModifyListModReq is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)
		message.SendPDUSessionResourceModifyResponse(nil, nil, nil, &criticalityDiagnostics)
		return
	}

	if (amfUeNgapID != nil) && (ranUeNgapID != nil) {
		// Find UE context
		var ok bool
		ranUe, ok = n3iwfCtx.RanUePoolLoad(ranUeNgapID.Value)
		if !ok {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
			// TODO: build cause and send error indication
			// Cause: Unknown local UE NGAP ID
			return
		}
		ranUeCtx = ranUe.GetSharedCtx()
		if ranUeCtx.AmfUeNgapId != amfUeNgapID.Value {
			// TODO: build cause and send error indication
			// Cause: Inconsistent remote UE NGAP ID
			return
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
				ngapLog.Errorf(
					"[PDUSessionID: %d] PDUSessionResourceModifyRequestTransfer Decode Error: %v\n",
					pduSessionID, err)
			}

			if pduSession = ranUeCtx.FindPDUSession(pduSessionID); pduSession == nil {
				ngapLog.Errorf("[PDUSessionID: %d] Unknown PDU session ID", pduSessionID)

				cause := message.BuildCause(ngapType.CausePresentRadioNetwork,
					ngapType.CauseRadioNetworkPresentUnknownPDUSessionID)
				unsuccessfulTransfer, buildErr := message.
					BuildPDUSessionResourceModifyUnsuccessfulTransfer(*cause, nil)
				if buildErr != nil {
					ngapLog.Errorf("Build PDUSessionResourceModifyUnsuccessfulTransfer Error: %v\n", buildErr)
				}
				message.AppendPDUSessionResourceFailedToModifyListModRes(
					failedListModRes, pduSessionID, unsuccessfulTransfer)
				continue
			}

			success, resTransfer := s.handlePDUSessionResourceModifyRequestTransfer(
				pduSession, transfer)
			if success {
				message.AppendPDUSessionResourceModifyListModRes(responseList, pduSessionID, resTransfer)
			} else {
				message.AppendPDUSessionResourceFailedToModifyListModRes(
					failedListModRes, pduSessionID, resTransfer)
			}
		}
	}

	message.SendPDUSessionResourceModifyResponse(ranUe, responseList, failedListModRes, nil)
}

func (s *Server) handlePDUSessionResourceModifyRequestTransfer(
	pduSession *n3iwf_context.PDUSession,
	transfer ngapType.PDUSessionResourceModifyRequestTransfer,
) (
	success bool, responseTransfer []byte,
) {
	ngapLog := logger.NgapLog
	ngapLog.Trace("Handle PDU Session Resource Modify Request Transfer")

	var pduSessionAMBR *ngapType.PDUSessionAggregateMaximumBitRate
	var ulNGUUPTNLModifyList *ngapType.ULNGUUPTNLModifyList
	var networkInstance *ngapType.NetworkInstance
	var qosFlowAddOrModifyRequestList *ngapType.QosFlowAddOrModifyRequestList
	var qosFlowToReleaseList *ngapType.QosFlowListWithCause
	// var additionalULNGUUPTNLInformation *ngapType.UPTransportLayerInformation

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	// used for building response transfer
	var resDLNGUUPTNLInfo *ngapType.UPTransportLayerInformation
	var resULNGUUPTNLInfo *ngapType.UPTransportLayerInformation
	var resQosFlowAddOrModifyRequestList ngapType.QosFlowAddOrModifyResponseList
	var resQosFlowFailedToAddOrModifyList ngapType.QosFlowListWithCause

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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDQosFlowToReleaseList:
			ngapLog.Traceln("[NGAP] Decode IE QosFlowToReleaseList")
			qosFlowToReleaseList = ie.Value.QosFlowToReleaseList
			if qosFlowToReleaseList != nil && len(qosFlowToReleaseList.List) == 0 {
				ngapLog.Error("qosFlowToReleaseList should have at least one element")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDAdditionalULNGUUPTNLInformation:
			ngapLog.Traceln("[NGAP] Decode IE AdditionalULNGUUPTNLInformation")
			// additionalULNGUUPTNLInformation = ie.Value.AdditionalULNGUUPTNLInformation
		}
	}

	if len(iesCriticalityDiagnostics.List) != 0 {
		// build unsuccessful transfer
		cause := message.BuildCause(ngapType.CausePresentProtocol,
			ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)
		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)
		unsuccessfulTransfer, err := message.BuildPDUSessionResourceModifyUnsuccessfulTransfer(
			*cause, &criticalityDiagnostics)
		if err != nil {
			ngapLog.Errorf("Build PDUSessionResourceModifyUnsuccessfulTransfer Error: %v\n", err)
		}

		responseTransfer = unsuccessfulTransfer
		return success, responseTransfer
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

				cause := message.BuildCause(
					ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentUnkownQosFlowID)

				item := ngapType.QosFlowWithCauseItem{
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

	// if additionalULNGUUPTNLInformation != nil {
	// TODO: forward AdditionalULNGUUPTNLInfomation to S-NG-RAN
	// }

	encodeData, err := message.BuildPDUSessionResourceModifyResponseTransfer(
		resULNGUUPTNLInfo, resDLNGUUPTNLInfo, &resQosFlowAddOrModifyRequestList, &resQosFlowFailedToAddOrModifyList)
	if err != nil {
		ngapLog.Errorf("Build PDUSessionResourceModifyTransfer Error: %v\n", err)
	}

	success = true
	responseTransfer = encodeData

	return success, responseTransfer
}

func (s *Server) HandlePDUSessionResourceModifyConfirm(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle PDU Session Resource Modify Confirm")

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceModifyListModCfm *ngapType.PDUSessionResourceModifyListModCfm
	var pDUSessionResourceFailedToModifyListModCfm *ngapType.PDUSessionResourceFailedToModifyListModCfm
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
	var ranUe n3iwf_context.RanUe
	var ranUeCtx *n3iwf_context.RanUeSharedCtx

	n3iwfCtx := s.Context()

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	successfulOutcome := pdu.SuccessfulOutcome
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

	if rANUENGAPID != nil {
		var ok bool
		ranUe, ok = n3iwfCtx.RanUePoolLoad(rANUENGAPID.Value)
		if !ok {
			ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", rANUENGAPID.Value)
			return
		}
		ranUeCtx = ranUe.GetSharedCtx()
	}

	if aMFUENGAPID != nil {
		if ranUe != nil {
			if ranUeCtx.AmfUeNgapId != aMFUENGAPID.Value {
				ngapLog.Errorf("Inconsistent remote UE NGAP ID, AMFUENGAPID: %d, RanUe.AmfUeNgapId: %d",
					aMFUENGAPID.Value, ranUeCtx.AmfUeNgapId)
				return
			}
		} else {
			ranUe = amf.FindUeByAmfUeNgapID(aMFUENGAPID.Value)
			if ranUe == nil {
				ngapLog.Errorf("Inconsistent remote UE NGAP ID, AMFUENGAPID: %d",
					aMFUENGAPID.Value)
				return
			}
		}
	}

	if ranUe == nil {
		ngapLog.Warn("RANUENGAPID and  AMFUENGAPID are both nil")
		return
	}

	if pDUSessionResourceModifyListModCfm != nil {
		for _, item := range pDUSessionResourceModifyListModCfm.List {
			pduSessionId := item.PDUSessionID.Value
			ngapLog.Tracef("PDU Session Id[%d] in Pdu Session Resource Modification Confrim List", pduSessionId)
			sess, exist := ranUeCtx.PduSessionList[pduSessionId]
			if !exist {
				ngapLog.Warnf(
					"PDU Session Id[%d] is not exist in Ue[ranUeNgapId:%d]", pduSessionId, ranUeCtx.RanUeNgapId)
			} else {
				transfer := ngapType.PDUSessionResourceModifyConfirmTransfer{}
				err := aper.UnmarshalWithParams(item.PDUSessionResourceModifyConfirmTransfer, &transfer, "valueExt")
				if err != nil {
					ngapLog.Warnf(
						"[PDUSessionID: %d] PDUSessionResourceSetupRequestTransfer Decode Error: %v\n",
						pduSessionId, err)
				} else if transfer.QosFlowFailedToModifyList != nil {
					for _, flow := range transfer.QosFlowFailedToModifyList.List {
						ngapLog.Warnf(
							"Delete QFI[%d] due to Qos Flow Failure in Pdu Session Resource Modification Confrim List",
							flow.QosFlowIdentifier.Value)
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
			err := aper.UnmarshalWithParams(
				item.PDUSessionResourceModifyIndicationUnsuccessfulTransfer, &transfer, "valueExt")
			if err != nil {
				ngapLog.Warnf(
					"[PDUSessionID: %d] PDUSessionResourceModifyIndicationUnsuccessfulTransfer Decode Error: %v\n",
					pduSessionId, err)
			} else {
				printAndGetCause(&transfer.Cause)
			}
			ngapLog.Tracef(
				"Release PDU Session Id[%d] due to PDU Session Resource Modify Indication Unsuccessful", pduSessionId)
			delete(ranUeCtx.PduSessionList, pduSessionId)
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func (s *Server) HandlePDUSessionResourceReleaseCommand(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle PDU Session Resource Release Command")
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	// var rANPagingPriority *ngapType.RANPagingPriority
	// var nASPDU *ngapType.NASPDU
	var pDUSessionResourceToReleaseListRelCmd *ngapType.PDUSessionResourceToReleaseListRelCmd

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	n3iwfCtx := s.Context()

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				ngapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANPagingPriority:
			ngapLog.Traceln("[NGAP] Decode IE RANPagingPriority")
			// rANPagingPriority = ie.Value.RANPagingPriority
		case ngapType.ProtocolIEIDNASPDU:
			ngapLog.Traceln("[NGAP] Decode IE NASPDU")
			// nASPDU = ie.Value.NASPDU
		case ngapType.ProtocolIEIDPDUSessionResourceToReleaseListRelCmd:
			ngapLog.Traceln("[NGAP] Decode IE PDUSessionResourceToReleaseListRelCmd")
			pDUSessionResourceToReleaseListRelCmd = ie.Value.PDUSessionResourceToReleaseListRelCmd
			if pDUSessionResourceToReleaseListRelCmd == nil {
				ngapLog.Error("PDUSessionResourceToReleaseListRelCmd is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodePDUSessionResourceRelease
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(
			&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		message.SendErrorIndication(amf, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	ranUe, ok := n3iwfCtx.RanUePoolLoad(rANUENGAPID.Value)
	if !ok {
		ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", rANUENGAPID.Value)
		cause := message.BuildCause(ngapType.CausePresentRadioNetwork,
			ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID)
		message.SendErrorIndication(amf, nil, nil, cause, nil)
		return
	}
	ranUeCtx := ranUe.GetSharedCtx()

	if ranUeCtx.AmfUeNgapId != aMFUENGAPID.Value {
		ngapLog.Errorf("Inconsistent remote UE NGAP ID, AMFUENGAPID: %d, RanUe.AmfUeNgapId: %d",
			aMFUENGAPID.Value, ranUeCtx.AmfUeNgapId)
		cause := message.BuildCause(ngapType.CausePresentRadioNetwork,
			ngapType.CauseRadioNetworkPresentInconsistentRemoteUENGAPID)
		message.SendErrorIndication(amf, nil, &rANUENGAPID.Value, cause, nil)
		return
	}

	// if rANPagingPriority != nil {
	// n3iwf does not support paging
	// }

	releaseList := ngapType.PDUSessionResourceReleasedListRelRes{}
	var releaseIdList []int64
	for _, item := range pDUSessionResourceToReleaseListRelCmd.List {
		pduSessionId := item.PDUSessionID.Value
		transfer := ngapType.PDUSessionResourceReleaseCommandTransfer{}
		err := aper.UnmarshalWithParams(item.PDUSessionResourceReleaseCommandTransfer, &transfer, "valueExt")
		if err != nil {
			ngapLog.Warnf(
				"[PDUSessionID: %d] PDUSessionResourceReleaseCommandTransfer Decode Error: %v\n",
				pduSessionId, err)
		} else {
			printAndGetCause(&transfer.Cause)
		}
		ngapLog.Tracef("Release PDU Session Id[%d] due to PDU Session Resource Release Command", pduSessionId)
		delete(ranUeCtx.PduSessionList, pduSessionId)

		// response list
		releaseItem := ngapType.PDUSessionResourceReleasedItemRelRes{
			PDUSessionID: item.PDUSessionID,
			PDUSessionResourceReleaseResponseTransfer: getPDUSessionResourceReleaseResponseTransfer(),
		}
		releaseList.List = append(releaseList.List, releaseItem)

		releaseIdList = append(releaseIdList, pduSessionId)
	}

	localSPI, ok := n3iwfCtx.IkeSpiLoad(rANUENGAPID.Value)
	if !ok {
		ngapLog.Errorf("Cannot get SPI from RanUeNgapID : %+v", rANUENGAPID.Value)
		return
	}
	ranUe.GetSharedCtx().PduSessResRelState = n3iwf_context.PduSessResRelStateOngoing

	s.SendIkeEvt(n3iwf_context.NewSendChildSADeleteRequestEvt(localSPI, releaseIdList))

	ranUeCtx.PduSessionReleaseList = releaseList
	// if nASPDU != nil {
	// TODO: Send NAS to UE
	// }
}

func (s *Server) HandleErrorIndication(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle Error Indication")

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if amf == nil {
		ngapLog.Error("Corresponding AMF context not found")
		return
	}
	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := pdu.InitiatingMessage
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

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(rANUENGAPID.Value)
	if ok {
		err := s.releaseIkeUeAndRanUe(ranUe)
		if err != nil {
			ngapLog.Warnf("HandleErrorIndication(): %v", err)
		}
	}

	ranUe = amf.FindUeByAmfUeNgapID(aMFUENGAPID.Value)
	if ranUe != nil {
		err := s.releaseIkeUeAndRanUe(ranUe)
		if err != nil {
			ngapLog.Warnf("HandleErrorIndication(): %v", err)
		}
	}

	// TODO: handle error based on cause/criticalityDiagnostics
}

func (s *Server) HandleUERadioCapabilityCheckRequest(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle UE Radio Capability Check Request")
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var uERadioCapability *ngapType.UERadioCapability
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	n3iwfCtx := s.Context()

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				ngapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(
					ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
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
		criticalityDiagnostics := buildCriticalityDiagnostics(
			&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		message.SendErrorIndication(amf, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	ranUe, ok := n3iwfCtx.RanUePoolLoad(rANUENGAPID.Value)
	if !ok {
		ngapLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", rANUENGAPID.Value)
		cause := message.BuildCause(ngapType.CausePresentRadioNetwork,
			ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID)
		message.SendErrorIndication(amf, nil, nil, cause, nil)
		return
	}

	ranUe.GetSharedCtx().RadioCapability = uERadioCapability
}

func (s *Server) HandleAMFConfigurationUpdate(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle AMF Configuration Updaet")

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

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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
	message.SendAMFConfigurationUpdateAcknowledge(amf, setupList, nil, nil)
}

func (s *Server) HandleRANConfigurationUpdateAcknowledge(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle RAN Configuration Update Acknowledge")

	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	successfulOutcome := pdu.SuccessfulOutcome
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

func (s *Server) HandleRANConfigurationUpdateFailure(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle RAN Configuration Update Failure")

	var cause *ngapType.Cause
	var timeToWait *ngapType.TimeToWait
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	n3iwfCtx := s.Context()

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
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
		ngapLog.Infof("Wait at lease  %ds to resend RAN Configuration Update to same AMF[%s]",
			waitingTime, amf.SCTPAddr)
		n3iwfCtx.AMFReInitAvailableListStore(amf.SCTPAddr, false)
		time.AfterFunc(time.Duration(waitingTime)*time.Second, func() {
			ngapLog.Infof("Re-send Ran Configuration Update Message when waiting time expired")
			n3iwfCtx.AMFReInitAvailableListStore(amf.SCTPAddr, true)
			message.SendRANConfigurationUpdate(n3iwfCtx, amf)
		})
		return
	}
	message.SendRANConfigurationUpdate(n3iwfCtx, amf)
}

func (s *Server) HandleDownlinkRANConfigurationTransfer(
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle Downlink RAN Configuration Transfer")
}

func (s *Server) HandleDownlinkRANStatusTransfer(
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle Downlink RAN Status Transfer")
}

func (s *Server) HandleAMFStatusIndication(
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle AMF Status Indication")
}

func (s *Server) HandleLocationReportingControl(
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle Location Reporting Control")
}

func (s *Server) HandleUETNLAReleaseRequest(
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle UE TNLA Release Request")
}

func (s *Server) HandleOverloadStart(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle Overload Start")

	var aMFOverloadResponse *ngapType.OverloadResponse
	var aMFTrafficLoadReductionIndication *ngapType.TrafficLoadReductionIndication
	var overloadStartNSSAIList *ngapType.OverloadStartNSSAIList

	if amf == nil {
		ngapLog.Error("AMF Context is nil")
		return
	}

	if pdu == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := pdu.InitiatingMessage
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

func (s *Server) HandleOverloadStop(
	amf *n3iwf_context.N3IWFAMF,
	pdu *ngapType.NGAPPDU,
) {
	ngapLog := logger.NgapLog
	ngapLog.Infoln("Handle Overload Stop")

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
	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList,
) (
	criticalityDiagnostics ngapType.CriticalityDiagnostics,
) {
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

func buildCriticalityDiagnosticsIEItem(
	ieCriticality aper.Enumerated,
	ieID int64,
	typeOfErr aper.Enumerated,
) (
	item ngapType.CriticalityDiagnosticsIEItem,
) {
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

func printAndGetCause(
	cause *ngapType.Cause,
) (
	present int, value aper.Enumerated,
) {
	ngapLog := logger.NgapLog
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

func printCriticalityDiagnostics(
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	ngapLog := logger.NgapLog
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
	ngapLog := logger.NgapLog
	data := ngapType.PDUSessionResourceReleaseResponseTransfer{}
	encodeData, err := aper.MarshalWithParams(data, "valueExt")
	if err != nil {
		ngapLog.Errorf("aper MarshalWithParams error in getPDUSessionResourceReleaseResponseTransfer: %d", err)
	}
	return encodeData
}

func (s *Server) HandleEvent(ngapEvent n3iwf_context.NgapEvt) {
	ngapLog := logger.NgapLog
	ngapLog.Infof("NGAP event handle")

	switch ngapEvent.Type() {
	case n3iwf_context.UnmarshalEAP5GData:
		s.HandleUnmarshalEAP5GData(ngapEvent)
	case n3iwf_context.SendInitialUEMessage:
		s.HandleSendInitialUEMessage(ngapEvent)
	case n3iwf_context.SendPDUSessionResourceSetupResponse:
		s.HandleSendPDUSessionResourceSetupResponse(ngapEvent)
	case n3iwf_context.SendNASMsg:
		s.HandleSendNASMsg(ngapEvent)
	case n3iwf_context.StartTCPSignalNASMsg:
		s.HandleStartTCPSignalNASMsg(ngapEvent)
	case n3iwf_context.NASTCPConnEstablishedComplete:
		s.HandleNASTCPConnEstablishedComplete(ngapEvent)
	case n3iwf_context.SendUEContextRelease:
		s.HandleSendSendUEContextRelease(ngapEvent)
	case n3iwf_context.SendUEContextReleaseRequest:
		s.HandleSendUEContextReleaseRequest(ngapEvent)
	case n3iwf_context.SendUEContextReleaseComplete:
		s.HandleSendUEContextReleaseComplete(ngapEvent)
	case n3iwf_context.SendPDUSessionResourceRelease:
		s.HandleSendSendPDUSessionResourceRelease(ngapEvent)
	case n3iwf_context.SendPDUSessionResourceReleaseResponse:
		s.HandleSendPDUSessionResourceReleaseRes(ngapEvent)
	case n3iwf_context.GetNGAPContext:
		s.HandleGetNGAPContext(ngapEvent)
	case n3iwf_context.SendUplinkNASTransport:
		s.HandleSendUplinkNASTransport(ngapEvent)
	case n3iwf_context.SendInitialContextSetupResponse:
		s.HandleSendInitialContextSetupResponse(ngapEvent)
	default:
		ngapLog.Errorf("Undefine NGAP event type")
		return
	}
}

func (s *Server) HandleGetNGAPContext(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle HandleGetNGAPContext Event")

	evt := ngapEvent.(*n3iwf_context.GetNGAPContextEvt)
	ranUeNgapId := evt.RanUeNgapId
	ngapCxtReqNumlist := evt.NgapCxtReqNumlist

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	var ngapCxt []interface{}

	for _, num := range ngapCxtReqNumlist {
		switch num {
		case n3iwf_context.CxtTempPDUSessionSetupData:
			ngapCxt = append(ngapCxt, ranUe.GetSharedCtx().TemporaryPDUSessionSetupData)
		default:
			ngapLog.Errorf("Receive undefine NGAP Context Request number : %d", num)
		}
	}

	spi, ok := n3iwfCtx.IkeSpiLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get spi from ngapid : %+v", ranUeNgapId)
		return
	}

	s.SendIkeEvt(n3iwf_context.NewGetNGAPContextRepEvt(spi, ngapCxtReqNumlist, ngapCxt))
}

func (s *Server) HandleUnmarshalEAP5GData(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle UnmarshalEAP5GData Event")

	evt := ngapEvent.(*n3iwf_context.UnmarshalEAP5GDataEvt)
	spi := evt.LocalSPI
	eapVendorData := evt.EAPVendorData
	isInitialUE := evt.IsInitialUE

	n3iwfCtx := s.Context()

	anParameters, nasPDU, err := UnmarshalEAP5GData(eapVendorData)
	if err != nil {
		ngapLog.Errorf("Unmarshalling EAP-5G packet failed: %v", err)
		return
	}

	if !isInitialUE { // ikeSA.ikeUE == nil
		ngapLog.Debug("Select AMF with the following AN parameters:")
		if anParameters.GUAMI == nil {
			ngapLog.Debug("\tGUAMI: nil")
		} else {
			ngapLog.Debugf("\tGUAMI: PLMNIdentity[% x], "+
				"AMFRegionID[% x], AMFSetID[% x], AMFPointer[% x]",
				anParameters.GUAMI.PLMNIdentity, anParameters.GUAMI.AMFRegionID,
				anParameters.GUAMI.AMFSetID, anParameters.GUAMI.AMFPointer)
		}
		if anParameters.SelectedPLMNID == nil {
			ngapLog.Debug("\tSelectedPLMNID: nil")
		} else {
			ngapLog.Debugf("\tSelectedPLMNID: % v", anParameters.SelectedPLMNID.Value)
		}
		if anParameters.RequestedNSSAI == nil {
			ngapLog.Debug("\tRequestedNSSAI: nil")
		} else {
			ngapLog.Debugf("\tRequestedNSSAI:")
			for i := 0; i < len(anParameters.RequestedNSSAI.List); i++ {
				ngapLog.Debugf("\tRequestedNSSAI:")
				ngapLog.Debugf("\t\tSNSSAI %d:", i+1)
				ngapLog.Debugf("\t\t\tSST: % x", anParameters.RequestedNSSAI.List[i].SNSSAI.SST.Value)
				sd := anParameters.RequestedNSSAI.List[i].SNSSAI.SD
				if sd == nil {
					ngapLog.Debugf("\t\t\tSD: nil")
				} else {
					ngapLog.Debugf("\t\t\tSD: % x", sd.Value)
				}
			}
		}

		selectedAMF := n3iwfCtx.AMFSelection(anParameters.GUAMI, anParameters.SelectedPLMNID)
		if selectedAMF == nil {
			s.SendIkeEvt(n3iwf_context.NewSendEAP5GFailureMsgEvt(spi, n3iwf_context.ErrAMFSelection))
		} else {
			n3iwfUe := n3iwfCtx.NewN3iwfRanUe()
			n3iwfUe.AMF = selectedAMF
			if anParameters.EstablishmentCause != nil {
				value := uint64(anParameters.EstablishmentCause.Value)
				if value > uint64(math.MaxInt16) {
					ngapLog.Errorf("HandleUnmarshalEAP5GData() anParameters.EstablishmentCause.Value "+
						"exceeds int16: %+v", value)
					return
				} else {
					n3iwfUe.RRCEstablishmentCause = int16(value)
				}
			}

			s.SendIkeEvt(n3iwf_context.NewUnmarshalEAP5GDataResponseEvt(spi, n3iwfUe.RanUeNgapId, nasPDU))
		}
	} else {
		ranUeNgapId := evt.RanUeNgapId
		ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
		if !ok {
			ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
			return
		}
		message.SendUplinkNASTransport(ranUe, nasPDU)
	}
}

func (s *Server) HandleSendInitialUEMessage(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendInitialUEMessage Event")

	evt := ngapEvent.(*n3iwf_context.SendInitialUEMessageEvt)
	ranUeNgapId := evt.RanUeNgapId
	ipv4Addr := evt.IPv4Addr
	ipv4Port := evt.IPv4Port
	nasPDU := evt.NasPDU

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}
	ranUeCtx := ranUe.GetSharedCtx()

	ranUeCtx.IPAddrv4 = ipv4Addr
	ranUeCtx.PortNumber = int32(ipv4Port) // #nosec G115
	message.SendInitialUEMessage(ranUeCtx.AMF, ranUe, nasPDU)
}

func (s *Server) HandleSendPDUSessionResourceSetupResponse(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendPDUSessionResourceSetupResponse Event")

	evt := ngapEvent.(*n3iwf_context.SendPDUSessionResourceSetupResEvt)
	ranUeNgapId := evt.RanUeNgapId

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}
	ranUeCtx := ranUe.GetSharedCtx()

	temporaryPDUSessionSetupData := ranUeCtx.TemporaryPDUSessionSetupData

	if len(temporaryPDUSessionSetupData.UnactivatedPDUSession) != 0 {
		for index, pduSession := range temporaryPDUSessionSetupData.UnactivatedPDUSession {
			errStr := temporaryPDUSessionSetupData.FailedErrStr[index]
			if errStr != n3iwf_context.ErrNil {
				var cause ngapType.Cause
				switch errStr {
				case n3iwf_context.ErrTransportResourceUnavailable:
					cause = ngapType.Cause{
						Present: ngapType.CausePresentTransport,
						Transport: &ngapType.CauseTransport{
							Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
						},
					}
				default:
					ngapLog.Errorf("Undefine event error string : %+s", errStr.Error())
					return
				}

				transfer, err := message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
				if err != nil {
					ngapLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %v", err)
					continue
				}

				if temporaryPDUSessionSetupData.NGAPProcedureCode.Value == ngapType.ProcedureCodeInitialContextSetup {
					message.AppendPDUSessionResourceFailedToSetupListCxtRes(
						temporaryPDUSessionSetupData.FailedListCxtRes, pduSession.Id, transfer)
				} else {
					message.AppendPDUSessionResourceFailedToSetupListSURes(
						temporaryPDUSessionSetupData.FailedListSURes, pduSession.Id, transfer)
				}
			} else {
				var gtpAddr string
				switch ranUe.(type) {
				case *n3iwf_context.N3IWFRanUe:
					gtpAddr = s.Config().GetN3iwfGtpBindAddress()
				}

				// Append NGAP PDU session resource setup response transfer
				transfer, err := message.BuildPDUSessionResourceSetupResponseTransfer(
					pduSession, gtpAddr)
				if err != nil {
					ngapLog.Errorf("Build PDU session resource setup response transfer failed: %v", err)
					return
				}
				if temporaryPDUSessionSetupData.NGAPProcedureCode.Value == ngapType.ProcedureCodeInitialContextSetup {
					message.AppendPDUSessionResourceSetupListCxtRes(
						temporaryPDUSessionSetupData.SetupListCxtRes, pduSession.Id, transfer)
				} else {
					message.AppendPDUSessionResourceSetupListSURes(
						temporaryPDUSessionSetupData.SetupListSURes, pduSession.Id, transfer)
				}
			}
		}

		if temporaryPDUSessionSetupData.NGAPProcedureCode.Value == ngapType.ProcedureCodeInitialContextSetup {
			message.SendInitialContextSetupResponse(ranUe,
				temporaryPDUSessionSetupData.SetupListCxtRes,
				temporaryPDUSessionSetupData.FailedListCxtRes, nil)
		} else {
			message.SendPDUSessionResourceSetupResponse(ranUe,
				temporaryPDUSessionSetupData.SetupListSURes,
				temporaryPDUSessionSetupData.FailedListSURes, nil)
		}
	} else {
		message.SendInitialContextSetupResponse(ranUe, nil, nil, nil)
	}
}

func (s *Server) HandleSendNASMsg(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendNASMsg Event")

	evt := ngapEvent.(*n3iwf_context.SendNASMsgEvt)
	ranUeNgapId := evt.RanUeNgapId

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	n3iwfUe, ok := ranUe.(*n3iwf_context.N3IWFRanUe)
	if !ok {
		ngapLog.Errorln("HandleSendNASMsg(): [Type Assertion] RanUe -> N3iwfUe failed")
		return
	}

	if n, ikeErr := n3iwfUe.TCPConnection.Write(n3iwfUe.TemporaryCachedNASMessage); ikeErr != nil {
		ngapLog.Errorf("Writing via IPSec signalling SA failed: %v", ikeErr)
	} else {
		ngapLog.Tracef("Forward PDU Seesion Establishment Accept to UE. Wrote %d bytes", n)
		n3iwfUe.TemporaryCachedNASMessage = nil
	}
}

func (s *Server) HandleStartTCPSignalNASMsg(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle StartTCPSignalNASMsg Event")

	evt := ngapEvent.(*n3iwf_context.StartTCPSignalNASMsgEvt)
	ranUeNgapId := evt.RanUeNgapId

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	n3iwfUe, ok := ranUe.(*n3iwf_context.N3IWFRanUe)
	if !ok {
		ngapLog.Errorln("HandleStartTCPSignalNASMsg(): [Type Assertion] RanUe -> N3iwfUe failed")
		return
	}

	n3iwfUe.IsNASTCPConnEstablished = true
}

func (s *Server) HandleNASTCPConnEstablishedComplete(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle NASTCPConnEstablishedComplete Event")

	evt := ngapEvent.(*n3iwf_context.NASTCPConnEstablishedCompleteEvt)
	ranUeNgapId := evt.RanUeNgapId

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}
	n3iwfUe, ok := ranUe.(*n3iwf_context.N3IWFRanUe)
	if !ok {
		ngapLog.Errorln("HandleNASTCPConnEstablishedComplete(): [Type Assertion] RanUe -> N3iwfUe failed")
		return
	}

	n3iwfUe.IsNASTCPConnEstablishedComplete = true

	if n3iwfUe.TemporaryCachedNASMessage != nil {
		// Send to UE
		if n, err := n3iwfUe.TCPConnection.Write(n3iwfUe.TemporaryCachedNASMessage); err != nil {
			ngapLog.Errorf("Writing via IPSec signalling SA failed: %v", err)
		} else {
			ngapLog.Trace("Forward NWu <- N2")
			ngapLog.Tracef("Wrote %d bytes", n)
		}
		n3iwfUe.TemporaryCachedNASMessage = nil
	}
}

func (s *Server) HandleSendUEContextReleaseRequest(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendUEContextReleaseRequest Event")

	evt := ngapEvent.(*n3iwf_context.SendUEContextReleaseRequestEvt)

	ranUeNgapId := evt.RanUeNgapId
	errMsg := evt.ErrMsg

	var cause *ngapType.Cause
	switch errMsg {
	case n3iwf_context.ErrRadioConnWithUeLost:
		cause = message.BuildCause(ngapType.CausePresentRadioNetwork,
			ngapType.CauseRadioNetworkPresentRadioConnectionWithUeLost)
	case n3iwf_context.ErrNil:
	default:
		ngapLog.Errorf("Undefine event error string : %+s", errMsg.Error())
		return
	}

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	message.SendUEContextReleaseRequest(ranUe, *cause)
}

func (s *Server) HandleSendUEContextReleaseComplete(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendUEContextReleaseComplete Event")

	evt := ngapEvent.(*n3iwf_context.SendUEContextReleaseCompleteEvt)
	ranUeNgapId := evt.RanUeNgapId

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	if err := ranUe.Remove(); err != nil {
		ngapLog.Errorf("Delete RanUe Context error : %v", err)
	}
	message.SendUEContextReleaseComplete(ranUe, nil)
}

func (s *Server) HandleSendPDUSessionResourceReleaseRes(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendPDUSessionResourceReleaseResponse Event")

	evt := ngapEvent.(*n3iwf_context.SendPDUSessionResourceReleaseResEvt)
	ranUeNgapId := evt.RanUeNgapId

	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	message.SendPDUSessionResourceReleaseResponse(ranUe, ranUe.GetSharedCtx().PduSessionReleaseList, nil)
}

func (s *Server) HandleSendUplinkNASTransport(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendUplinkNASTransport Event")

	evt := ngapEvent.(*n3iwf_context.SendUplinkNASTransportEvt)
	ranUeNgapId := evt.RanUeNgapId
	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	message.SendUplinkNASTransport(ranUe, evt.Pdu)
}

func (s *Server) HandleSendInitialContextSetupResponse(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendInitialContextSetupResponse Event")

	evt := ngapEvent.(*n3iwf_context.SendInitialContextSetupRespEvt)
	ranUeNgapId := evt.RanUeNgapId
	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	message.SendInitialContextSetupResponse(ranUe, evt.ResponseList, evt.FailedList, evt.CriticalityDiagnostics)
}

func (s *Server) HandleSendSendUEContextRelease(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendSendUEContextRelease Event")

	evt := ngapEvent.(*n3iwf_context.SendUEContextReleaseEvt)
	ranUeNgapId := evt.RanUeNgapId
	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	if ranUe.GetSharedCtx().UeCtxRelState {
		if err := ranUe.Remove(); err != nil {
			ngapLog.Errorf("Delete RanUe Context error : %v", err)
		}
		message.SendUEContextReleaseComplete(ranUe, nil)
		ranUe.GetSharedCtx().UeCtxRelState = n3iwf_context.UeCtxRelStateNone
	} else {
		cause := message.BuildCause(ngapType.CausePresentRadioNetwork,
			ngapType.CauseRadioNetworkPresentRadioConnectionWithUeLost)
		message.SendUEContextReleaseRequest(ranUe, *cause)
		ranUe.GetSharedCtx().UeCtxRelState = n3iwf_context.UeCtxRelStateOngoing
	}
}

func (s *Server) HandleSendSendPDUSessionResourceRelease(
	ngapEvent n3iwf_context.NgapEvt,
) {
	ngapLog := logger.NgapLog
	ngapLog.Tracef("Handle SendSendPDUSessionResourceRelease Event")

	evt := ngapEvent.(*n3iwf_context.SendPDUSessionResourceReleaseEvt)
	ranUeNgapId := evt.RanUeNgapId
	deletPduIds := evt.DeletPduIds
	n3iwfCtx := s.Context()
	ranUe, ok := n3iwfCtx.RanUePoolLoad(ranUeNgapId)
	if !ok {
		ngapLog.Errorf("Cannot get RanUE from ranUeNgapId : %+v", ranUeNgapId)
		return
	}

	if ranUe.GetSharedCtx().PduSessResRelState {
		message.SendPDUSessionResourceReleaseResponse(ranUe, ranUe.GetSharedCtx().PduSessionReleaseList, nil)
		ranUe.GetSharedCtx().PduSessResRelState = n3iwf_context.PduSessResRelStateNone
	} else {
		for _, id := range deletPduIds {
			ranUe.GetSharedCtx().DeletePDUSession(id)
		}
		ranUe.GetSharedCtx().PduSessResRelState = n3iwf_context.PduSessResRelStateOngoing
	}
}
