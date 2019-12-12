package gmm_message

import (
	"free5gc/lib/nas/nasType"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/amf_ngap/ngap_message"
	"free5gc/src/amf/amf_util"
	"free5gc/src/amf/logger"
)

// backOffTimerUint = 7 means backoffTimer is null
func SendDLNASTransport(ue *amf_context.RanUe, payloadContainerType uint8, nasPdu []byte, pduSessionId *int32, cause uint8, backOffTimerUint *uint8, backOffTimer uint8) {

	logger.GmmLog.Info("[NAS] Send DL NAS Transport")
	var causePtr *uint8
	if cause != 0 {
		causePtr = &cause
	}
	pduSessionID := uint8(*pduSessionId)
	nasMsg, err := BuildDLNASTransport(ue.AmfUe, payloadContainerType, nasPdu, &pduSessionID, causePtr, backOffTimerUint, backOffTimer)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}

func SendNotification(ue *amf_context.RanUe, nasMsg []byte) {

	logger.GmmLog.Info("[NAS] Send Notification")

	amfUe := ue.AmfUe
	if amfUe == nil {
		logger.GmmLog.Error("AmfUe is nil")
		return
	}
	amf_util.StartT3565(ue)
	amfUe.LastNotificationPkg = nasMsg
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}

func SendIdentityRequest(ue *amf_context.RanUe, typeOfIdentity uint8) {

	logger.GmmLog.Info("[NAS] Send Identity Request")

	nasMsg, err := BuildIdentityRequest(typeOfIdentity)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}

func SendAuthenticationRequest(ue *amf_context.RanUe) {

	amfUe := ue.AmfUe
	if amfUe == nil {
		logger.GmmLog.Error("AmfUe is nil")
		return
	}

	logger.GmmLog.Infof("[NAS] Send Authentication Request[Retry: %d]", amfUe.T3560RetryTimes)

	if amfUe.AuthenticationCtx == nil {
		logger.GmmLog.Error("Authentication Context of UE is nil")
		return
	}

	nasMsg, err := BuildAuthenticationRequest(amfUe)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)

	amf_util.StartT3560(ue, amf_message.EventGMMT3560ForAuthenticationRequest, nil, nil)
}

func SendServiceAccept(ue *amf_context.RanUe, pDUSessionStatus *[16]bool, reactivationResult *[16]bool, errPduSessionId, errCause []uint8) {

	logger.GmmLog.Info("[NAS] Send Service Accept")

	nasMsg, err := BuildServiceAccept(ue.AmfUe, pDUSessionStatus, reactivationResult, errPduSessionId, errCause)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}

func SendConfigurationUpdateCommand(amfUe *amf_context.AmfUe, accessType models.AccessType, networkSlicingIndication *nasType.NetworkSlicingIndication) {

	logger.GmmLog.Info("[NAS] Configuration Update Command")

	nasMsg, err := BuildConfigurationUpdateCommand(amfUe, accessType, networkSlicingIndication)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(amfUe.RanUe[accessType], nasMsg)
}

func SendAuthenticationReject(ue *amf_context.RanUe, eapMsg string) {

	logger.GmmLog.Info("[NAS] Send Authentication Reject")

	nasMsg, err := BuildAuthenticationReject(ue.AmfUe, eapMsg)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}

func SendAuthenticationResult(ue *amf_context.RanUe, eapSuccess bool, eapMsg string) {

	logger.GmmLog.Info("[NAS] Send Authentication Result")

	if ue.AmfUe == nil {
		logger.GmmLog.Errorf("AmfUe is nil")
		return
	}

	nasMsg, err := BuildAuthenticationResult(ue.AmfUe, eapSuccess, eapMsg)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}
func SendServiceReject(ue *amf_context.RanUe, pDUSessionStatus *[16]bool, cause uint8) {

	logger.GmmLog.Info("[NAS] Send Service Reject")

	nasMsg, err := BuildServiceReject(pDUSessionStatus, cause)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}

// T3502: This IE may be included to indicate a value for timer T3502 during the initial registration
// eapMessage: if the REGISTRATION REJECT message is used to convey EAP-failure message
func SendRegistrationReject(ue *amf_context.RanUe, cause5GMM uint8, eapMessage string) {

	logger.GmmLog.Info("[NAS] Send Registration Reject")

	nasMsg, err := BuildRegistrationReject(ue.AmfUe, cause5GMM, eapMessage)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}

// eapSuccess: only used when authType is EAP-AKA', set the value to false if authType is not EAP-AKA'
// eapMessage: only used when authType is EAP-AKA', set the value to "" if authType is not EAP-AKA'
func SendSecurityModeCommand(ue *amf_context.RanUe, eapSuccess bool, eapMessage string) {

	logger.GmmLog.Info("[NAS] Send Security Mode Command")

	nasMsg, err := BuildSecurityModeCommand(ue.AmfUe, eapSuccess, eapMessage)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)

	amf_util.StartT3560(ue, amf_message.EventGMMT3560ForSecurityModeCommand, &eapSuccess, &eapMessage)
}

func SendDeregistrationRequest(ue *amf_context.RanUe, accessType uint8, reRegistrationRequired bool, cause5GMM uint8) {

	logger.GmmLog.Info("[NAS] Send Deregistration Request")

	nasMsg, err := BuildDeregistrationRequest(ue, accessType, reRegistrationRequired, cause5GMM)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)

	amf_util.StartT3522(ue, &accessType, &reRegistrationRequired, &cause5GMM)
}

func SendDeregistrationAccept(ue *amf_context.RanUe) {

	logger.GmmLog.Info("[NAS] Send Deregistration Accept")

	nasMsg, err := BuildDeregistrationAccept()
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}

func SendRegistrationAccept(
	ue *amf_context.AmfUe,
	anType models.AccessType,
	pDUSessionStatus *[16]bool,
	reactivationResult *[16]bool,
	errPduSessionId, errCause []uint8,
	pduSessionResourceSetupList *ngapType.PDUSessionResourceSetupListCxtReq) {

	logger.GmmLog.Info("[NAS] Send Registration Accept")

	nasMsg, err := BuildRegistrationAccept(ue, anType, pDUSessionStatus, reactivationResult, errPduSessionId, errCause)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendInitialContextSetupRequest(ue, anType, nasMsg, nil, pduSessionResourceSetupList, nil, nil, nil)
	amf_util.StartT3550(ue, anType, pDUSessionStatus, reactivationResult, errPduSessionId, errCause, pduSessionResourceSetupList)
}

func SendStatus5GMM(ue *amf_context.RanUe, cause uint8) {

	logger.GmmLog.Info("[NAS] Send Status 5GMM")

	nasMsg, err := BuildStatus5GMM(cause)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg)
}
