package ngap_message

import (
	"github.com/sirupsen/logrus"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapSctp"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_util"
	"free5gc/src/amf/logger"
)

var ngaplog *logrus.Entry

func init() {
	ngaplog = logger.NgapLog
}

func SendToRan(ran *amf_context.AmfRan, packet []byte) {

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	if len(packet) == 0 {
		ngaplog.Error("packet len is 0")
		return
	}

	ngaplog.Debugf("[NGAP] Send To Ran [IP: %s]", ran.Conn.RemoteAddr().String())

	ngapSctp.SendMsg(ran.Conn, packet)
}

func SendToRanUe(ue *amf_context.RanUe, packet []byte) {

	var ran *amf_context.AmfRan

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if ran = ue.Ran; ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	if ue.AmfUe == nil {
		ngaplog.Warn("AmfUe is nil")
	}

	SendToRan(ran, packet)
}

func NasSendToRan(ue *amf_context.AmfUe, packet []byte) {

	if ue == nil {
		ngaplog.Error("AmfUe is nil")
		return
	}

	ranUe := ue.RanUe[ue.GetAnType()]
	if ranUe == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	SendToRanUe(ranUe, packet)
}

func SendNGSetupResponse(ran *amf_context.AmfRan) {

	ngaplog.Info("[AMF] Send NG-Setup response")

	pkt, err := BuildNGSetupResponse()
	if err != nil {
		ngaplog.Errorf("Build NGSetupResponse failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendNGSetupFailure(ran *amf_context.AmfRan, cause ngapType.Cause) {

	ngaplog.Info("[AMF] Send NG-Setup failure")

	if cause.Present == ngapType.CausePresentNothing {
		ngaplog.Errorf("Cause present is nil")
		return
	}

	pkt, err := BuildNGSetupFailure(cause)
	if err != nil {
		ngaplog.Errorf("Build NGSetupFailure failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

// partOfNGInterface: if reset type is "reset all", set it to nil TS 38.413 9.2.6.11
func SendNGReset(ran *amf_context.AmfRan, cause ngapType.Cause, partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList) {

	ngaplog.Info("[AMF] Send NG Reset")

	pkt, err := BuildNGReset(cause, partOfNGInterface)
	if err != nil {
		ngaplog.Errorf("Build NGReset failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendNGResetAcknowledge(ran *amf_context.AmfRan, partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send NG Reset Acknowledge")

	if partOfNGInterface != nil && len(partOfNGInterface.List) == 0 {
		ngaplog.Error("length of partOfNGInterface is 0")
		return
	}

	pkt, err := BuildNGResetAcknowledge(partOfNGInterface, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build NGResetAcknowledge failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendDownlinkNasTransport(ue *amf_context.RanUe, nasPdu []byte) {

	ngaplog.Info("[AMF] Send Downlink Nas Transport")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(nasPdu) == 0 {
		ngaplog.Errorf("[Send DownlinkNasTransport] Error: nasPdu is nil")
	}

	pkt, err := BuildDownlinkNasTransport(ue, nasPdu)
	if err != nil {
		ngaplog.Errorf("Build NGResetAcknowledge failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendPDUSessionResourceReleaseCommand(ue *amf_context.RanUe, nasPdu []byte, pduSessionResourceReleasedList ngapType.PDUSessionResourceToReleaseListRelCmd) {

	ngaplog.Info("[AMF] Send PDU Session Resource Release Command")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(nasPdu) == 0 {
		ngaplog.Errorf("[Send PDUSessionResourceReleaseCommand] Error: nasPdu is nil")
		return
	}

	pkt, err := BuildPDUSessionResourceReleaseCommand(ue, nasPdu, pduSessionResourceReleasedList)
	if err != nil {
		ngaplog.Errorf("Build PDUSessionResourceReleaseCommand failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendUEContextReleaseCommand(ue *amf_context.RanUe, action amf_context.RelAction, causePresent int, cause aper.Enumerated) {

	ngaplog.Info("[AMF] Send UE Context Release Command")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildUEContextReleaseCommand(ue, causePresent, cause)
	if err != nil {
		ngaplog.Errorf("Build UEContextReleaseCommand failed : %s", err.Error())
		return
	}
	ue.ReleaseAction = action
	if ue.AmfUe != nil && ue.Ran != nil {
		ue.AmfUe.ReleaseCause[ue.Ran.AnType] = &amf_context.CauseAll{
			NgapCause: &models.NgApCause{
				Group: int32(causePresent),
				Value: int32(cause),
			},
		}
	}
	SendToRanUe(ue, pkt)
}

func SendErrorIndication(ran *amf_context.AmfRan, amfUeNgapId, ranUeNgapId *int64, cause *ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send Error Indication")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildErrorIndication(amfUeNgapId, ranUeNgapId, cause, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build ErrorIndication failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendUERadioCapabilityCheckRequest(ue *amf_context.RanUe) {

	ngaplog.Info("[AMF] Send UE Radio Capability Check Request")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildUERadioCapabilityCheckRequest(ue)
	if err != nil {
		ngaplog.Errorf("Build UERadioCapabilityCheckRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendHandoverCancelAcknowledge(ue *amf_context.RanUe, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send Handover Cancel Acknowledge")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildHandoverCancelAcknowledge(ue, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build HandoverCancelAcknowledge failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// nasPDU: from nas layer
// pduSessionResourceSetupRequestList: provided by AMF, and transfer data is from SMF
func SendPDUSessionResourceSetupRequest(ue *amf_context.RanUe, nasPdu []byte, pduSessionResourceSetupRequestList ngapType.PDUSessionResourceSetupListSUReq) {

	ngaplog.Info("[AMF] Send PDU Session Resource Setup Request")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(pduSessionResourceSetupRequestList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPDUSessionResourceSetupRequest(ue, nasPdu, pduSessionResourceSetupRequestList)
	if err != nil {
		ngaplog.Errorf("Build PDUSessionResourceSetupRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// pduSessionResourceModifyConfirmList: provided by AMF, and transfer data is return from SMF
// pduSessionResourceFailedToModifyList: provided by AMF, and transfer data is return from SMF
func SendPDUSessionResourceModifyConfirm(
	ue *amf_context.RanUe,
	pduSessionResourceModifyConfirmList ngapType.PDUSessionResourceModifyListModCfm,
	pduSessionResourceFailedToModifyList ngapType.PDUSessionResourceFailedToModifyListModCfm,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send PDU Session Resource Modify Confirm")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(pduSessionResourceModifyConfirmList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	if len(pduSessionResourceFailedToModifyList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPDUSessionResourceModifyConfirm(ue, pduSessionResourceModifyConfirmList, pduSessionResourceFailedToModifyList, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build PDUSessionResourceModifyConfirm failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// pduSessionResourceModifyRequestList: from SMF
func SendPDUSessionResourceModifyRequest(ue *amf_context.RanUe, pduSessionResourceModifyRequestList ngapType.PDUSessionResourceModifyListModReq) {

	ngaplog.Info("[AMF] Send PDU Session Resource Modify Request")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(pduSessionResourceModifyRequestList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPDUSessionResourceModifyRequest(ue, pduSessionResourceModifyRequestList)
	if err != nil {
		ngaplog.Errorf("Build PDUSessionResourceModifyRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendInitialContextSetupRequest(
	amfUe *amf_context.AmfUe,
	anType models.AccessType,
	nasPdu []byte,
	oldAmf *string,
	pduSessionResourceSetupRequestList *ngapType.PDUSessionResourceSetupListCxtReq,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	coreNetworkAssistanceInfo *ngapType.CoreNetworkAssistanceInformation,
	emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator) {

	ngaplog.Info("[AMF] Send Initial Context Setup Request")

	if amfUe == nil {
		ngaplog.Error("AmfUe is nil")
		return
	}

	if pduSessionResourceSetupRequestList != nil {
		if len(pduSessionResourceSetupRequestList.List) > amf_context.MaxNumOfPDUSessions {
			ngaplog.Error("Pdu List out of range")
			return
		}
	}

	pkt, err := BuildInitialContextSetupRequest(amfUe, anType, nasPdu, oldAmf, pduSessionResourceSetupRequestList,
		rrcInactiveTransitionReportRequest, coreNetworkAssistanceInfo, emergencyFallbackIndicator)
	if err != nil {
		ngaplog.Errorf("Build InitialContextSetupRequest failed : %s", err.Error())
		return
	}
	NasSendToRan(amfUe, pkt)
}

func SendUEContextModificationRequest(
	amfUe *amf_context.AmfUe,
	anType models.AccessType,
	oldAmfUeNgapID *int64,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	coreNetworkAssistanceInfo *ngapType.CoreNetworkAssistanceInformation,
	mobilityRestrictionList *ngapType.MobilityRestrictionList,
	emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator) {

	ngaplog.Info("[AMF] Send UE Context Modification Request")

	if amfUe == nil {
		ngaplog.Error("AmfUe is nil")
		return
	}

	pkt, err := BuildUEContextModificationRequest(amfUe, anType, oldAmfUeNgapID, rrcInactiveTransitionReportRequest, coreNetworkAssistanceInfo, mobilityRestrictionList, emergencyFallbackIndicator)
	if err != nil {
		ngaplog.Errorf("Build UEContextModificationRequest failed : %s", err.Error())
		return
	}
	NasSendToRan(amfUe, pkt)
}

// pduSessionResourceHandoverList: provided by amf and transfer is return from smf
// pduSessionResourceToReleaseList: provided by amf and transfer is return from smf
// criticalityDiagnostics = criticalityDiagonstics IE in receiver node's error indication when received node can't comprehend the IE or missing IE
func SendHandoverCommand(
	sourceUe *amf_context.RanUe,
	pduSessionResourceHandoverList ngapType.PDUSessionResourceHandoverList,
	pduSessionResourceToReleaseList ngapType.PDUSessionResourceToReleaseListHOCmd,
	container ngapType.TargetToSourceTransparentContainer,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send Handover Command")

	if sourceUe == nil {
		ngaplog.Error("SourceUe is nil")
		return
	}

	if len(pduSessionResourceHandoverList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	if len(pduSessionResourceToReleaseList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildHandoverCommand(sourceUe, pduSessionResourceHandoverList, pduSessionResourceToReleaseList, container, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build HandoverCommand failed : %s", err.Error())
		return
	}
	SendToRanUe(sourceUe, pkt)
}

// cause = initiate the Handover Cancel procedure with the appropriate value for the Cause IE.
// criticalityDiagnostics = criticalityDiagonstics IE in receiver node's error indication when received node can't comprehend the IE or missing IE
func SendHandoverPreparationFailure(sourceUe *amf_context.RanUe, cause ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send Handover Preparation Failure")

	if sourceUe == nil {
		ngaplog.Error("SourceUe is nil")
		return
	}
	amfUe := sourceUe.AmfUe
	if amfUe == nil {
		ngaplog.Error("amfUe is nil")
		return
	}
	amfUe.OnGoing[sourceUe.Ran.AnType].Procedure = amf_context.OnGoingProcedureNothing
	pkt, err := BuildHandoverPreparationFailure(sourceUe, cause, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build HandoverPreparationFailure failed : %s", err.Error())
		return
	}
	SendToRanUe(sourceUe, pkt)
}

/*The PGW-C+SMF (V-SMF in the case of home-routed roaming scenario only) sends
a Nsmf_PDUSession_CreateSMContext Response(N2 SM Information (PDU Session ID, cause code)) to the AMF.*/
// Cause is from SMF
// pduSessionResourceSetupList provided by AMF, and the transfer data is from SMF
// sourceToTargetTransparentContainer is received from S-RAN
// nsci: new security context indicator, if amfUe has updated security context, set nsci to true, otherwise set to false
// N2 handover in same AMF
func SendHandoverRequest(sourceUe *amf_context.RanUe, targetRan *amf_context.AmfRan, cause ngapType.Cause, pduSessionResourceSetupListHOReq ngapType.PDUSessionResourceSetupListHOReq,
	sourceToTargetTransparentContainer ngapType.SourceToTargetTransparentContainer, nsci bool) {

	ngaplog.Info("[AMF] Send Handover Request")

	if sourceUe == nil {
		ngaplog.Error("sourceUe is nil")
		return
	}
	amfUe := sourceUe.AmfUe
	if amfUe == nil {
		ngaplog.Error("amfUe is nil")
		return
	}
	if targetRan == nil {
		ngaplog.Error("targetRan is nil")
		return
	}

	if sourceUe.TargetUe != nil {
		ngaplog.Error("Handover Required Duplicated")
		return
	}

	if len(pduSessionResourceSetupListHOReq.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	if len(sourceToTargetTransparentContainer.Value) == 0 {
		ngaplog.Error("Source To Target TransparentContainer is nil")
		return
	}

	targetUe := targetRan.NewRanUe()

	ngaplog.Tracef("Source : AMF_UE_NGAP_ID[%d], RAN_UE_NGAP_ID[%d]", sourceUe.AmfUeNgapId, sourceUe.RanUeNgapId)
	ngaplog.Tracef("Target : AMF_UE_NGAP_ID[%d], RAN_UE_NGAP_ID[Unknown]", targetUe.AmfUeNgapId)
	amf_context.AttachSourceUeTargetUe(sourceUe, targetUe)

	pkt, err := BuildHandoverRequest(targetUe, cause, pduSessionResourceSetupListHOReq, sourceToTargetTransparentContainer, nsci)
	if err != nil {
		ngaplog.Errorf("Build HandoverRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(targetUe, pkt)
}

// pduSessionResourceSwitchedList: provided by AMF, and the transfer data is from SMF
// pduSessionResourceReleasedList: provided by AMF, and the transfer data is from SMF
// newSecurityContextIndicator: if AMF has activated a new 5G NAS security context, set it to true, otherwise set to false
// coreNetworkAssistanceInformation: provided by AMF, based on collection of UE behaviour statistics and/or other available
// information about the expected UE behaviour. TS 23.501 5.4.6, 5.4.6.2
// rrcInactiveTransitionReportRequest: configured by amf
// criticalityDiagnostics: from received node when received not comprehended IE or missing IE
func SendPathSwitchRequestAcknowledge(
	ue *amf_context.RanUe,
	pduSessionResourceSwitchedList ngapType.PDUSessionResourceSwitchedList,
	pduSessionResourceReleasedList ngapType.PDUSessionResourceReleasedListPSAck,
	newSecurityContextIndicator bool,
	coreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send Path Switch Request Acknowledge")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(pduSessionResourceSwitchedList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	if len(pduSessionResourceReleasedList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPathSwitchRequestAcknowledge(ue, pduSessionResourceSwitchedList, pduSessionResourceReleasedList,
		newSecurityContextIndicator, coreNetworkAssistanceInformation, rrcInactiveTransitionReportRequest, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build PathSwitchRequestAcknowledge failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// pduSessionResourceReleasedList: provided by AMF, and the transfer data is from SMF
// criticalityDiagnostics: from received node when received not comprehended IE or missing IE
func SendPathSwitchRequestFailure(
	ran *amf_context.AmfRan,
	amfUeNgapId,
	ranUeNgapId int64,
	pduSessionResourceReleasedList *ngapType.PDUSessionResourceReleasedListPSFail,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send Path Switch Request Failure")

	if pduSessionResourceReleasedList != nil && len(pduSessionResourceReleasedList.List) > amf_context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPathSwitchRequestFailure(amfUeNgapId, ranUeNgapId, pduSessionResourceReleasedList, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build PathSwitchRequestFailure failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

//ranStatusTransferTransparentContainer from Uplink Ran Configuration Transfer
func SendDownlinkRanStatusTransfer(ue *amf_context.RanUe, ranStatusTransferTransparentContainer ngapType.RANStatusTransferTransparentContainer) {

	ngaplog.Info("[AMF] Send Downlink Ran Status Transfer")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(ranStatusTransferTransparentContainer.DRBsSubjectToStatusTransferList.List) > amf_context.MaxNumOfDRBs {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildDownlinkRanStatusTransfer(ue, ranStatusTransferTransparentContainer)
	if err != nil {
		ngaplog.Errorf("Build DownlinkRanStatusTransfer failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// anType indicate amfUe send this msg for which accessType
// Paging Priority: is included only if the AMF receives an Namf_Communication_N1N2MessageTransfer message with an ARP value associated with
// priority services (e.g., MPS, MCS), as configured by the operator. (TS 23.502 4.2.3.3, TS 23.501 5.22.3)
// pagingOriginNon3GPP: TS 23.502 4.2.3.3 step 4b: If the UE is simultaneously registered over 3GPP and non-3GPP accesses in the same PLMN,
// the UE is in CM-IDLE state in both 3GPP access and non-3GPP access, and the PDU Session ID in step 3a
// is associated with non-3GPP access, the AMF sends a Paging message with associated access "non-3GPP" to
// NG-RAN node(s) via 3GPP access.
// more paging policy with 3gpp/non-3gpp access is described in TS 23.501 5.6.8
func SendPaging(ue *amf_context.AmfUe, ngapBuf []byte) {

	// var pagingPriority *ngapType.PagingPriority
	if ue == nil {
		ngaplog.Error("AmfUe is nil")
		return
	}
	ue.LastPagingPkg = ngapBuf
	/* Start T3513 */
	amf_util.StartT3513(ue)
	// if ppi != nil {
	// pagingPriority = new(ngapType.PagingPriority)
	// pagingPriority.Value = aper.Enumerated(*ppi)
	// }
	// pkt, err := BuildPaging(ue, pagingPriority, pagingOriginNon3GPP)
	// if err != nil {
	// 	ngaplog.Errorf("Build Paging failed : %s", err.Error())
	// }
	taiList := ue.RegistrationArea[models.AccessType__3_GPP_ACCESS]
	for _, ran := range amf_context.AMF_Self().AmfRanPool {
		for _, item := range ran.SupportedTAList {
			if amf_context.InTaiList(item.Tai, taiList) {
				ngaplog.Infof("[AMF] Send Paging to TAI(%+v, Tac:%+v) for Ue[%s]", item.Tai.PlmnId, item.Tai.Tac, ue.Supi)
				SendToRan(ran, ngapBuf)
				break
			}
		}
	}

}

// TS 23.502 4.2.2.2.3
// anType: indicate amfUe send this msg for which accessType
// amfUeNgapID: initial AMF get it from target AMF
// ngapMessage: initial UE Message to reroute
// allowedNSSAI: provided by AMF, and AMF get it from NSSF (4.2.2.2.3 step 4b)
func SendRerouteNasRequest(ue *amf_context.AmfUe, anType models.AccessType, amfUeNgapID *int64, ngapMessage []byte, allowedNSSAI *ngapType.AllowedNSSAI) {

	ngaplog.Info("[AMF] Send Reroute Nas Request")

	if ue == nil {
		ngaplog.Error("AmfUe is nil")
		return
	}

	if len(ngapMessage) == 0 {
		ngaplog.Error("Ngap Message is nil")
		return
	}

	pkt, err := BuildRerouteNasRequest(ue, anType, amfUeNgapID, ngapMessage, allowedNSSAI)
	if err != nil {
		ngaplog.Errorf("Build RerouteNasRequest failed : %s", err.Error())
		return
	}
	NasSendToRan(ue, pkt)
}

// criticality ->from received node when received node can't comprehend the IE or missing IE
func SendRanConfigurationUpdateAcknowledge(ran *amf_context.AmfRan, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send Ran Configuration Update Acknowledge")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildRanConfigurationUpdateAcknowledge(criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build RanConfigurationUpdateAcknowledge failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

// criticality ->from received node when received node can't comprehend the IE or missing IE
// If the AMF cannot accept the update,
// it shall respond with a RAN CONFIGURATION UPDATE FAILURE message and appropriate cause value.
func SendRanConfigurationUpdateFailure(ran *amf_context.AmfRan, cause ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[AMF] Send Ran Configuration Update Failure")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildRanConfigurationUpdateFailure(cause, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build RanConfigurationUpdateFailure failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

//An AMF shall be able to instruct other peer CP NFs, subscribed to receive such a notification,
//that it will be unavailable on this AMF and its corresponding target AMF(s).
//If CP NF does not subscribe to receive AMF unavailable notification, the CP NF may attempt
//forwarding the transaction towards the old AMF and detect that the AMF is unavailable. When
//it detects unavailable, it marks the AMF and its associated GUAMI(s) as unavailable.
//Defined in 23.501 5.21.2.2.2
func SendAMFStatusIndication(ran *amf_context.AmfRan, unavailableGUAMIList ngapType.UnavailableGUAMIList) {

	ngaplog.Info("[AMF] Send AMF Status Indication")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	if len(unavailableGUAMIList.List) > amf_context.MaxNumOfServedGuamiList {
		ngaplog.Error("GUAMI List out of range")
		return
	}

	pkt, err := BuildAMFStatusIndication(unavailableGUAMIList)
	if err != nil {
		ngaplog.Errorf("Build AMFStatusIndication failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

// TS 23.501 5.19.5.2
// amfOverloadResponse: the required behaviour of NG-RAN, provided by AMF
// amfTrafficLoadReductionIndication(int 1~99): indicates the percentage of the type, set to 0 if does not need this ie
// of traffic relative to the instantaneous incoming rate at the NG-RAN node, provided by AMF
// overloadStartNSSAIList: overload slices, provide by AMF
func SendOverloadStart(
	ran *amf_context.AmfRan,
	amfOverloadResponse *ngapType.OverloadResponse,
	amfTrafficLoadReductionIndication int64,
	overloadStartNSSAIList *ngapType.OverloadStartNSSAIList) {

	ngaplog.Info("[AMF] Send Overload Start")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	if amfTrafficLoadReductionIndication != 0 && (amfTrafficLoadReductionIndication < 1 || amfTrafficLoadReductionIndication > 99) {
		ngaplog.Error("AmfTrafficLoadReductionIndication out of range (should be 1 ~ 99)")
		return
	}

	if overloadStartNSSAIList != nil && len(overloadStartNSSAIList.List) > amf_context.MaxNumOfSlice {
		ngaplog.Error("NSSAI List out of range")
		return
	}

	pkt, err := BuildOverloadStart(amfOverloadResponse, amfTrafficLoadReductionIndication, overloadStartNSSAIList)
	if err != nil {
		ngaplog.Errorf("Build OverloadStart failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendOverloadStop(ran *amf_context.AmfRan) {

	ngaplog.Info("[AMF] Send Overload Stop")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildOverloadStop()
	if err != nil {
		ngaplog.Errorf("Build OverloadStop failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

// sONConfigurationTransfer = sONConfigurationTransfer from uplink Ran Configuration Transfer
func SendDownlinkRanConfigurationTransfer(ran *amf_context.AmfRan, sONConfigurationTransfer *ngapType.SONConfigurationTransfer) {

	ngaplog.Info("[AMF] Send Downlink Ran Configuration Transfer")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildDownlinkRanConfigurationTransfer(sONConfigurationTransfer)
	if err != nil {
		ngaplog.Errorf("Build DownlinkRanConfigurationTransfer failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

//NRPPa PDU is by pass
//NRPPa PDU is from LMF define in 4.13.5.6
func SendDownlinkNonUEAssociatedNRPPATransport(ue *amf_context.RanUe, nRPPaPDU ngapType.NRPPaPDU) {

	ngaplog.Info("[AMF] Send Downlink Non UE Associated NRPPA Transport")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(nRPPaPDU.Value) == 0 {
		ngaplog.Error("length of NRPPA-PDU is 0")
		return
	}

	pkt, err := BuildDownlinkNonUEAssociatedNRPPATransport(ue, nRPPaPDU)
	if err != nil {
		ngaplog.Errorf("Build DownlinkNonUEAssociatedNRPPATransport failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendDeactivateTrace(amfUe *amf_context.AmfUe, anType models.AccessType) {

	ngaplog.Info("[AMF] Send Deactivate Trace")

	if amfUe == nil {
		ngaplog.Error("AmfUe is nil")
		return
	}

	ranUe := amfUe.RanUe[anType]
	if ranUe == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildDeactivateTrace(amfUe, anType)
	if err != nil {
		ngaplog.Errorf("Build DeactivateTrace failed : %s", err.Error())
		return
	}
	SendToRanUe(ranUe, pkt)
}

// AOI List is from SMF
// The SMF may subscribe to the UE mobility event notification from the AMF (e.g. location reporting, UE moving into or out of Area Of Interest) TS 23.502 4.3.2.2.1 Step.17
// The Location Reporting Control message shall identify the UE for which reports are requested and may include Reporting Type, Location Reporting Level, Area Of Interest and Request Reference ID TS 23.502 4.10 LocationReportingProcedure
// The AMF may request the NG-RAN location reporting with event reporting type (e.g. UE location or UE presence in Area of Interest), reporting mode and its related parameters (e.g. number of reporting) TS 23.501 5.4.7
// Location Reference ID To Be Cancelled IE shall be present if the Event Type IE is set to "Stop UE presence in the area of interest". otherwise set it to 0
func SendLocationReportingControl(
	ue *amf_context.RanUe,
	AOIList *ngapType.AreaOfInterestList,
	LocationReportingReferenceIDToBeCancelled int64,
	eventType ngapType.EventType) {

	ngaplog.Info("[AMF] Send Location Reporting Control")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if AOIList != nil && len(AOIList.List) > amf_context.MaxNumOfAOI {
		ngaplog.Error("AOI List out of range")
		return
	}

	if eventType.Value == ngapType.EventTypePresentStopUePresenceInAreaOfInterest {
		if LocationReportingReferenceIDToBeCancelled < 1 || LocationReportingReferenceIDToBeCancelled > 64 {
			ngaplog.Error("LocationReportingReferenceIDToBeCancelled out of range (should be 1 ~ 64)")
			return
		}
	}

	pkt, err := BuildLocationReportingControl(ue, AOIList, LocationReportingReferenceIDToBeCancelled, eventType)
	if err != nil {
		ngaplog.Errorf("Build LocationReportingControl failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendUETNLABindingReleaseRequest(ue *amf_context.RanUe) {

	ngaplog.Info("[AMF] Send UE TNLA Binging Release Request")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildUETNLABindingReleaseRequest(ue)
	if err != nil {
		ngaplog.Errorf("Build UETNLABindingReleaseRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// Weight Factor associated with each of the TNL association within the AMF
func SendAMFConfigurationUpdate(ran *amf_context.AmfRan, tNLassociationUsage ngapType.TNLAssociationUsage, tNLAddressWeightFactor ngapType.TNLAddressWeightFactor) {

	ngaplog.Info("[AMF] Send AMF Configuration Update")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildAMFConfigurationUpdate(tNLassociationUsage, tNLAddressWeightFactor)
	if err != nil {
		ngaplog.Errorf("Build AMFConfigurationUpdate failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

//NRPPa PDU is a pdu from LMF to RAN defined in TS 23.502 4.13.5.5 step 3
//NRPPa PDU is by pass
func SendDownlinkUEAssociatedNRPPaTransport(ue *amf_context.RanUe, nRPPaPDU ngapType.NRPPaPDU) {

	ngaplog.Info("[AMF] Send Downlink UE Associated NRPPa Transport")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(nRPPaPDU.Value) == 0 {
		ngaplog.Error("length of NRPPA-PDU is 0")
		return
	}

	pkt, err := BuildDownlinkUEAssociatedNRPPaTransport(ue, nRPPaPDU)
	if err != nil {
		ngaplog.Errorf("Build DownlinkUEAssociatedNRPPaTransport failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}
