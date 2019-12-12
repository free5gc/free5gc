package ngap_message

import (
	"github.com/sirupsen/logrus"

	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_ngap/n3iwf_sctp"
)

var ngaplog *logrus.Entry

func init() {
	ngaplog = logger.NgapLog
}

func SendToAmf(sessionID string, packet []byte) {
	if ok := n3iwf_sctp.Send(sessionID, packet); !ok {
		// TODO: Feature: retry sending
	}
}

func SendNGSetupRequest(sessionID string) {
	ngaplog.Infoln("[N3IWF] Send NG Setup Request")
	pkt, err := BuildNGSetupRequest()
	if err != nil {
		ngaplog.Errorf("Build NGSetup Request failed: %+v\n", err)
		return
	}

	SendToAmf(sessionID, pkt)
}

func SendNGReset() {
	ngaplog.Infoln("[N3IWF] Send NG Reset")
}

func SendNGResetAcknowledge() {
	ngaplog.Infoln("[N3IWF] Send NG Reset Acknowledge")
}

func SendInitialContextSetupResponse(
	sessionID string,
	ue *n3iwf_context.N3IWFUe,
	responseList *ngapType.PDUSessionResourceSetupListCxtRes,
	failedList *ngapType.PDUSessionResourceFailedToSetupListCxtRes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Infoln("[N3IWF] Send Initial Context Setup Response")

	if responseList != nil && len(responseList.List) > n3iwf_context.MaxNumOfPDUSessions {
		ngaplog.Errorln("Pdu List out of range")
		return
	}

	if failedList != nil && len(failedList.List) > n3iwf_context.MaxNumOfPDUSessions {
		ngaplog.Errorln("Pdu List out of range")
		return
	}

	pkt, err := BuildInitialContextSetupResponse(ue, responseList, failedList, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build Initial Context Setup Response failed : %+v\n", err)
		return
	}

	SendToAmf(sessionID, pkt)
}

func SendInitialContextSetupFailure(
	sessionID string,
	ue *n3iwf_context.N3IWFUe,
	cause ngapType.Cause,
	failedList *ngapType.PDUSessionResourceFailedToSetupListCxtFail,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Infoln("[N3IWF] Send Initial Context Setup Failure")

	if failedList != nil && len(failedList.List) > n3iwf_context.MaxNumOfPDUSessions {
		ngaplog.Errorln("Pdu List out of range")
		return
	}

	pkt, err := BuildInitialContextSetupFailure(ue, cause, failedList, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build Initial Context Setup Failure failed : %+v\n", err)
		return
	}

	SendToAmf(sessionID, pkt)
}

func SendUEContextModificationResponse(
	sessionID string,
	ue *n3iwf_context.N3IWFUe,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Infoln("[N3IWF] Send UE Context Modification Response")

	pkt, err := BuildUEContextModificationResponse(ue, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build UE Context Modification Response failed : %+v\n", err)
		return
	}

	SendToAmf(sessionID, pkt)
}

func SendUEContextModificationFailure(
	sessionID string,
	ue *n3iwf_context.N3IWFUe,
	cause ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Infoln("[N3IWF] Send UE Context Modification Failure")

	pkt, err := BuildUEContextModificationFailure(ue, cause, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build UE Context Modification Failure failed : %+v\n", err)
		return
	}

	SendToAmf(sessionID, pkt)
}

func SendUEContextReleaseComplete(
	sessionID string,
	ue *n3iwf_context.N3IWFUe,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Infoln("[N3IWF] Send UE Context Release Complete")

	pkt, err := BuildUEContextReleaseComplete(ue, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build UE Context Release Complete failed : %+v\n", err)
		return
	}

	SendToAmf(sessionID, pkt)
}

func SendUEContextReleaseRequest(
	sessionID string,
	ue *n3iwf_context.N3IWFUe, cause ngapType.Cause) {

	ngaplog.Infoln("[N3IWF] Send UE Context Release Request")

	pkt, err := BuildUEContextReleaseRequest(ue, cause)
	if err != nil {
		ngaplog.Errorf("Build UE Context Release Request failed : %+v\n", err)
		return
	}

	SendToAmf(sessionID, pkt)
}

func SendInitialUEMessage() {
	ngaplog.Infoln("[N3IWF] Send Initial UE Message")
}

func SendUplinkNASTransport(
	sessionID string,
	ue *n3iwf_context.N3IWFUe,
	nasPdu []byte) {

	ngaplog.Infoln("[N3IWF] Send Uplink NAS Transport")

	if len(nasPdu) == 0 {
		ngaplog.Errorln("NAS Pdu is nil")
		return
	}

	pkt, err := BuildUplinkNASTransport(ue, nasPdu)
	if err != nil {
		ngaplog.Errorf("Build Uplink NAS Transport failed : %+v\n", err)
		return
	}

	SendToAmf(sessionID, pkt)
}

func SendNASNonDeliveryIndication() {
	ngaplog.Infoln("[N3IWF] Send NAS NonDelivery Indication")
}

func SendRerouteNASRequest() {
	ngaplog.Infoln("[N3IWF] Send Reroute NAS Request")
}

func SendPDUSessionResourceSetupResponse() {
	ngaplog.Infoln("[N3IWF] Send PDU Session Resource Setup Response")
}

func SendPDUSessionResourceModifyResponse() {
	ngaplog.Infoln("[N3IWF] Send PDU Session Resource Modify Response")
}

func SendPDUSessionResourceModifyIndication() {
	ngaplog.Infoln("[N3IWF] Send PDU Session Resource Modify Indication")
}

func SendPDUSessionResourceNotify() {
	ngaplog.Infoln("[N3IWF] Send PDU Session Resource Notify")
}

func SendPDUSessionResourceReleaseResponse() {
	ngaplog.Infoln("[N3IWF] Send PDU Session Resource Release Response")
}

func SendErrorIndication() {
	ngaplog.Infoln("[N3IWF] Send Error Indication")
}

func SendUERadioCapabilityInfoIndication() {
	ngaplog.Infoln("[N3IWF] Send UE Radio Capability Info Indication")
}

func SendUERadioCapabilityCheckResponse() {
	ngaplog.Infoln("[N3IWF] Send UE Radio Capability Check Response")
}

func SendAMFConfigurationUpdateAcknowledge() {
	ngaplog.Infoln("[N3IWF] Send AMF Configuration Update Acknowledge")
}

func SendAMFConfigurationUpdateFailure() {
	ngaplog.Infoln("[N3IWF] Send AMF Configuration Update Failure")
}

func SendRANConfigurationUpdate() {
	ngaplog.Infoln("[N3IWF] Send RAN Configuration Update")
}

func SendUplinkRANConfigurationTransfer() {
	ngaplog.Infoln("[N3IWF] Send Uplink RAN Configuration Transfer")
}

func SendUplinkRANStatusTransfer() {
	ngaplog.Infoln("[N3IWF] Send Uplink RAN Status Transfer")
}

func SendLocationReportingFailureIndication() {
	ngaplog.Infoln("[N3IWF] Send Location Reporting Failure Indication")
}

func SendLocationReport() {
	ngaplog.Infoln("[N3IWF] Send Location Report")
}

func SendRRCInactiveTransitionReport() {
	ngaplog.Infoln("[N3IWF] Send RRC Inactive Transition Report")
}
