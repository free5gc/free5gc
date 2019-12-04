package ngap_handler

import (
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_handler/n3iwf_message"
	"free5gc/src/n3iwf/n3iwf_ngap/ngap_message"
	"time"
)

var ngapLog *logrus.Entry

func init() {
	ngapLog = logger.NgapLog
}

func HandleEventSCTPConnect(sctpSessionID string) {
	ngapLog.Infoln("[N3IWF] Handle SCTP connect event")
	ngap_message.SendNGSetupRequest(sctpSessionID)
}

func HandleNGSetupResponse(sctpSessionID string, message *ngapType.NGAPPDU) {
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
		// TODO: Send error indication
	}

	amfInfo := n3iwfSelf.NewN3iwfAmf(sctpSessionID)

	if amfName != nil {
		amfInfo.AMFName = *amfName
	}

	if servedGUAMIList != nil {
		amfInfo.ServedGUAMIList = *servedGUAMIList
	}

	if relativeAMFCapacity != nil {
		amfInfo.RelativeAMFCapacity = *relativeAMFCapacity
	}

	if plmnSupportList != nil {
		amfInfo.PLMNSupportList = *plmnSupportList
	}

	if criticalityDiagnostics != nil {
		// TODO: handle criticalityDiagnostics
	}
}

func HandleNGSetupFailure(sctpSessionID string, message *ngapType.NGAPPDU) {
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
	}

	if cause != nil {
		printAndGetCause(cause)
	}

	if criticalityDiagnostics != nil {
		// TODO: Handle criticalityDiagnostics
	}

	var waittingTime time.Duration

	if timeToWait != nil {

		switch timeToWait.Value {
		case ngapType.TimeToWaitPresentV1s:
			waittingTime = 1
		case ngapType.TimeToWaitPresentV2s:
			waittingTime = 2
		case ngapType.TimeToWaitPresentV5s:
			waittingTime = 5
		case ngapType.TimeToWaitPresentV10s:
			waittingTime = 10
		case ngapType.TimeToWaitPresentV20s:
			waittingTime = 20
		case ngapType.TimeToWaitPresentV60s:
			waittingTime = 60
		}

	}

	if waittingTime != 0 {
		time.Sleep(waittingTime * time.Second)
	}

	// TODO: Limited retry mechanism
	handlerMessage := n3iwf_message.HandlerMessage{
		Event:         n3iwf_message.EventSCTPConnectMessage,
		SCTPSessionID: sctpSessionID,
	}
	n3iwf_message.SendMessage(handlerMessage)
}

func HandleNGReset(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle NG Reset")
}

func HandleNGResetAcknowledge(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle NG Reset Acknowledge")
}

func HandleInitialContextSetupRequest(sctpSessionID string, message *ngapType.NGAPPDU) {
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
		// TODO: Send Error Indication
	}

	if amfUeNgapID != nil {
		n3iwfUe = n3iwfSelf.FindUeByAmfUeNgapID(amfUeNgapID.Value)
		if n3iwfUe == nil {
			ngapLog.Warnf("No UE Context[AmfUeNgapID:%d]\n", amfUeNgapID.Value)
		}
	}

	if n3iwfUe == nil && ranUeNgapID != nil {
		n3iwfUe = n3iwfSelf.FindUeByRanUeNgapID(ranUeNgapID.Value)
		if n3iwfUe == nil {
			ngapLog.Warnf("No UE Context[RanUeNgapID:%d]\n", ranUeNgapID.Value)
			return
		}
	}

	n3iwfUe.AmfUeNgapId = amfUeNgapID.Value
	n3iwfUe.RanUeNgapId = ranUeNgapID.Value

	var responseList *ngapType.PDUSessionResourceSetupListCxtRes
	var failedListCxtRes *ngapType.PDUSessionResourceFailedToSetupListCxtRes
	var failedListCxtFail *ngapType.PDUSessionResourceFailedToSetupListCxtFail

	if pduSessionResourceSetupListCxtReq != nil {
		if ueAggregateMaximumBitRate != nil {
			n3iwfUe.Ambr = ueAggregateMaximumBitRate
		} else {
			ngapLog.Errorln("IE[UEAggregateMaximumBitRate] is nil")
			cause := ngapType.Cause{}
			cause.Present = ngapType.CausePresentProtocol
			cause.Protocol = &ngapType.CauseProtocol{
				Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
			}
			criticalityDiagnosticsIEItem := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDUEAggregateMaximumBitRate, ngapType.TypeOfErrorPresentMissing)
			iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, criticalityDiagnosticsIEItem)
			criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)

			failedListCxtFail = new(ngapType.PDUSessionResourceFailedToSetupListCxtFail)
			for _, item := range pduSessionResourceSetupListCxtReq.List {
				transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
				if err != nil {
					ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", err)
				}
				ngap_message.AppendPDUSessionResourceFailedToSetupListCxtfail(failedListCxtFail, item.PDUSessionID.Value, transfer)
			}

			ngap_message.SendInitialContextSetupFailure(sctpSessionID, n3iwfUe, cause, failedListCxtFail, &criticalityDiagnostics)
			return
		}

		var pduSession *n3iwf_context.PDUSession
		responseList = new(ngapType.PDUSessionResourceSetupListCxtRes)
		failedListCxtRes = new(ngapType.PDUSessionResourceFailedToSetupListCxtRes)

		for _, item := range pduSessionResourceSetupListCxtReq.List {
			pduSessionID := item.PDUSessionID.Value
			// TODO: send NAS to UE
			// pduSessionNasPdu := item.NASPDU
			snssai := item.SNSSAI
			transfer := ngapType.PDUSessionResourceSetupRequestTransfer{}
			err := aper.UnmarshalWithParams(item.PDUSessionResourceSetupRequestTransfer, transfer, "valueExt")
			if err != nil {
				ngapLog.Errorf("[PDUSessionID: %d] PDUSessionResourceSetupRequestTransfer Decode Error: %+v\n", pduSessionID, err)
			}

			if pduSession = n3iwfUe.FindPDUSession(pduSessionID); pduSession == nil {
				pduSession, err = n3iwfUe.CreatePDUSession(pduSessionID, snssai)
				if err != nil {
					ngapLog.Errorf("Create PDU Session Error: %+v\n", err)

					cause := ngapType.Cause{}
					cause.Present = ngapType.CausePresentRadioNetwork
					cause.RadioNetwork = &ngapType.CauseRadioNetwork{
						Value: ngapType.CauseRadioNetworkPresentMultiplePDUSessionIDInstances,
					}
					unsuccessfulTransfer, buildErr := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
					if buildErr != nil {
						ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", buildErr)
					}
					ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(failedListCxtRes, pduSessionID, unsuccessfulTransfer)
					continue
				}
			}

			success, resTransfer := handlePDUSessionResourceSetupRequestTransfer(pduSession, transfer)
			if success {
				ngap_message.AppendPDUSessionResourceSetupListCxtRes(responseList, pduSessionID, resTransfer)
			} else {
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

	// TODO: use Kn3iwf to generate security context

	if nasPDU != nil {
		// TODO: Send NAS to UE
	}

	ngap_message.SendInitialContextSetupResponse(sctpSessionID, n3iwfUe, responseList, failedListCxtRes, nil)
}

// TODO: finish handle PDUSessionResourceSetupRequestTransfer
func handlePDUSessionResourceSetupRequestTransfer(pduSession *n3iwf_context.PDUSession, transfer ngapType.PDUSessionResourceSetupRequestTransfer) (success bool, responseTransfer []byte) {

	var pduSessionType *ngapType.PDUSessionType
	var ulNGUUPTNLInformation *ngapType.UPTransportLayerInformation
	var qosFlowSetupRequestList *ngapType.QosFlowSetupRequestList
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	for _, ie := range transfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDPDUSessionType:
			pduSessionType = ie.Value.PDUSessionType
			if pduSessionType == nil {
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDULNGUUPTNLInformation:
			ulNGUUPTNLInformation = ie.Value.ULNGUUPTNLInformation
			if ulNGUUPTNLInformation == nil {
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
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
		cause := ngapType.Cause{
			Present: ngapType.CausePresentProtocol,
			Protocol: &ngapType.CauseProtocol{
				Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
			},
		}
		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)
		unsuccessfulTransfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, &criticalityDiagnostics)
		if err != nil {
			ngapLog.Errorf("Build PDUSessionResourceSetupUnsuccessfulTransfer Error: %+v\n", err)
		}
		responseTransfer = unsuccessfulTransfer
		return
	}

	pduSession.Type = *pduSessionType

	// TODO: configure gtpu connection
	tunnel := ulNGUUPTNLInformation.GTPTunnel
	pduSession.GTPEndpointIPv4, pduSession.GTPEndpointIPv6 = ngapConvert.IPAddressToString(tunnel.TransportLayerAddress)
	pduSession.TEID = hex.EncodeToString(tunnel.GTPTEID.Value)
	ngapLog.Debugf("PDU Session[%d]: get NG-U info[ipv4:%s, ipv6:%s, TEID:%s]\n", pduSession.Id, pduSession.GTPEndpointIPv4, pduSession.GTPEndpointIPv6, pduSession.TEID)

	// TODO: apply qos rule
	for _, item := range qosFlowSetupRequestList.List {
		qosFlow := new(n3iwf_context.QosFlow)
		qosFlow.Identifier = item.QosFlowIdentifier.Value
		qosFlow.Parameters = item.QosFlowLevelQosParameters
		pduSession.QosFlows[qosFlow.Identifier] = qosFlow
	}

	success = true

	// TODO: allocate N3 Tunnel Info

	encodeData, err := ngap_message.BuildPDUSessionResourceSetupResponseTransfer(pduSession)
	if err != nil {
		ngapLog.Errorf("Encode PDUSessionResourceSetupResponseTransfer Error: %+v\n", err)
	}
	responseTransfer = encodeData
	return
}

func HandleUEContextModificationRequest(sctpSessionID string, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle UE Context Modification Request")

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
		// TODO: send error indication
		return
	}

	if ranUeNgapID != nil {
		n3iwfUe = n3iwfSelf.FindUeByRanUeNgapID(ranUeNgapID.Value)
	}

	if n3iwfUe == nil && amfUeNgapID != nil {
		n3iwfUe = n3iwfSelf.FindUeByAmfUeNgapID(amfUeNgapID.Value)
	}

	if n3iwfUe == nil {
		// TODO: send error indication
		return
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

	ngap_message.SendUEContextModificationResponse(sctpSessionID, n3iwfUe, nil)
}

func HandleUEContextReleaseCommand(sctpSessionID string, message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle UE Context Release Command")

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
			n3iwfUe = n3iwfSelf.FindUeByAmfUeNgapID(ueNgapIDs.UENGAPIDPair.AMFUENGAPID.Value)
		}
	case ngapType.UENGAPIDsPresentAMFUENGAPID:
		n3iwfUe = n3iwfSelf.FindUeByAmfUeNgapID(ueNgapIDs.AMFUENGAPID.Value)
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

	ngap_message.SendUEContextReleaseComplete(sctpSessionID, n3iwfUe, nil)
}

func HandleDownlinkNASTransport(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Downlink NAS Transport")
}

func HandlePDUSessionResourceSetupRequest(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle PDU Session Resource Setup Request")
}

func HandlePDUSessionResourceModifyRequest(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle PDU Session Resource Modify Request")
}

func HandlePDUSessionResourceModifyConfirm(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle PDU Session Resource Modify Confirm")
}

func HandlePDUSessionResourceReleaseCommand(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle PDU Session Resource Release Command")
}

func HandleErrorIndication(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Error Indication")
}

func HandleUERadioCapabilityCheckRequest(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle UE Radio Capability Check Request")
}

func HandleAMFConfigurationUpdate(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle AMF Configuration Updaet")
}

func HandleRANConfigurationUpdateAcknowledge(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle RAN Configuration Update Acknowledge")
}

func HandleRANConfigurationUpdateFailure(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle RAN Configuration Update Failure")
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

func HandleOverloadStart(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Overload Start")
}

func HandleOverloadStop(message *ngapType.NGAPPDU) {
	ngapLog.Infoln("[N3IWF] Handle Overload Stop")
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
