package gmm_message

import (
	"encoding/base64"
	"encoding/hex"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasConvert"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_nas/nas_security"
	"free5gc/src/amf/gmm/gmm_state"
	"free5gc/src/amf/logger"

	"github.com/mitchellh/mapstructure"
)

func BuildDLNASTransport(ue *amf_context.AmfUe, payloadContainerType uint8, nasPdu []byte, pduSessionId *uint8, cause *uint8, backoffTimerUint *uint8, backoffTimer uint8) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDLNASTransport)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	dLNASTransport := nasMessage.NewDLNASTransport(0)
	dLNASTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	dLNASTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	dLNASTransport.SetMessageType(nas.MsgTypeDLNASTransport)
	dLNASTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(payloadContainerType)
	dLNASTransport.PayloadContainer.SetLen(uint16(len(nasPdu)))
	dLNASTransport.PayloadContainer.SetPayloadContainerContents(nasPdu)

	if pduSessionId != nil {
		dLNASTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
		dLNASTransport.PduSessionID2Value.SetIei(nasMessage.DLNASTransportPduSessionID2ValueType)
		dLNASTransport.PduSessionID2Value.SetPduSessionID2Value(*pduSessionId)

	}
	if cause != nil {
		dLNASTransport.Cause5GMM = new(nasType.Cause5GMM)
		dLNASTransport.Cause5GMM.SetIei(nasMessage.DLNASTransportCause5GMMType)
		dLNASTransport.Cause5GMM.SetCauseValue(*cause)
	}
	if backoffTimerUint != nil {
		dLNASTransport.BackoffTimerValue = new(nasType.BackoffTimerValue)
		dLNASTransport.BackoffTimerValue.SetIei(nasMessage.DLNASTransportBackoffTimerValueType)
		dLNASTransport.BackoffTimerValue.SetLen(1)
		dLNASTransport.BackoffTimerValue.SetUnitTimerValue(*backoffTimerUint)
		dLNASTransport.BackoffTimerValue.SetTimerValue(backoffTimer)
	}

	m.GmmMessage.DLNASTransport = dLNASTransport

	return nas_security.Encode(ue, m)
}

func BuildNotification(ue *amf_context.AmfUe, accessType uint8) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeNotification)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	notification := nasMessage.NewNotification(0)
	notification.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypeIntegrityProtectedAndCiphered)
	notification.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	notification.SetMessageType(nas.MsgTypeNotification)
	notification.SetAccessType(accessType)

	m.GmmMessage.Notification = notification

	return nas_security.Encode(ue, m)
}

func BuildIdentityRequest(typeOfIdentity uint8) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeIdentityRequest)

	identityRequest := nasMessage.NewIdentityRequest(0)
	identityRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	identityRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	identityRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	identityRequest.IdentityRequestMessageIdentity.SetMessageType(nas.MsgTypeIdentityRequest)
	identityRequest.SpareHalfOctetAndIdentityType.SetTypeOfIdentity(typeOfIdentity)

	m.GmmMessage.IdentityRequest = identityRequest

	return m.PlainNasEncode()
}

func BuildAuthenticationRequest(ue *amf_context.AmfUe) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationRequest)

	authenticationRequest := nasMessage.NewAuthenticationRequest(0)
	authenticationRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	authenticationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationRequest.AuthenticationRequestMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationRequest)
	authenticationRequest.SpareHalfOctetAndNgksi = nasConvert.SpareHalfOctetAndNgksiToNas(ue.NgKsi)
	authenticationRequest.ABBA.SetLen(uint8(len(ue.ABBA)))
	authenticationRequest.ABBA.SetABBAContents(ue.ABBA)

	switch ue.AuthenticationCtx.AuthType {
	case models.AuthType__5_G_AKA:
		var tmpArray [16]byte
		var av5gAka models.Av5gAka

		if err := mapstructure.Decode(ue.AuthenticationCtx.Var5gAuthData, &av5gAka); err != nil {
			logger.GmmLog.Error("Var5gAuthData Convert Type Error")
			return nil, err
		}

		rand, _ := hex.DecodeString(av5gAka.Rand)
		authenticationRequest.AuthenticationParameterRAND = nasType.NewAuthenticationParameterRAND(nasMessage.AuthenticationRequestAuthenticationParameterRANDType)
		copy(tmpArray[:], rand[0:16])
		authenticationRequest.AuthenticationParameterRAND.SetRANDValue(tmpArray)

		autn, _ := hex.DecodeString(av5gAka.Autn)
		authenticationRequest.AuthenticationParameterAUTN = nasType.NewAuthenticationParameterAUTN(nasMessage.AuthenticationRequestAuthenticationParameterAUTNType)
		authenticationRequest.AuthenticationParameterAUTN.SetLen(uint8(len(autn)))
		copy(tmpArray[:], autn[0:16])
		authenticationRequest.AuthenticationParameterAUTN.SetAUTN(tmpArray)
	case models.AuthType_EAP_AKA_PRIME:
		eapMsg := ue.AuthenticationCtx.Var5gAuthData.(string)
		rawEapMsg, _ := base64.StdEncoding.DecodeString(eapMsg)
		authenticationRequest.EAPMessage = nasType.NewEAPMessage(nasMessage.AuthenticationRequestEAPMessageType)
		authenticationRequest.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		authenticationRequest.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.AuthenticationRequest = authenticationRequest

	return m.PlainNasEncode()
}

func BuildServiceAccept(ue *amf_context.AmfUe, pDUSessionStatus *[16]bool, reactivationResult *[16]bool, errPduSessionId, errCause []uint8) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeServiceAccept)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	serviceAccept := nasMessage.NewServiceAccept(0)
	serviceAccept.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	serviceAccept.SetSecurityHeaderType(nas.SecurityHeaderTypeIntegrityProtectedAndCiphered)
	serviceAccept.SetMessageType(nas.MsgTypeServiceAccept)
	if pDUSessionStatus != nil {
		serviceAccept.PDUSessionStatus = new(nasType.PDUSessionStatus)
		serviceAccept.PDUSessionStatus.SetIei(nasMessage.ServiceAcceptPDUSessionStatusType)
		serviceAccept.PDUSessionStatus.SetLen(2)
		serviceAccept.PDUSessionStatus.Buffer = nasConvert.PSIToBuf(*pDUSessionStatus)
	}
	if reactivationResult != nil {
		serviceAccept.PDUSessionReactivationResult = new(nasType.PDUSessionReactivationResult)
		serviceAccept.PDUSessionReactivationResult.SetIei(nasMessage.ServiceAcceptPDUSessionReactivationResultType)
		serviceAccept.PDUSessionReactivationResult.SetLen(2)
		serviceAccept.PDUSessionReactivationResult.Buffer = nasConvert.PSIToBuf(*reactivationResult)
	}
	if errPduSessionId != nil {
		serviceAccept.PDUSessionReactivationResultErrorCause = new(nasType.PDUSessionReactivationResultErrorCause)
		serviceAccept.PDUSessionReactivationResultErrorCause.SetIei(nasMessage.ServiceAcceptPDUSessionReactivationResultErrorCauseType)
		buf := nasConvert.PDUSessionReactivationResultErrorCauseToBuf(errPduSessionId, errCause)
		serviceAccept.PDUSessionReactivationResultErrorCause.SetLen(uint16(len(buf)))
		serviceAccept.PDUSessionReactivationResultErrorCause.Buffer = buf
	}
	m.GmmMessage.ServiceAccept = serviceAccept

	return nas_security.Encode(ue, m)
}

func BuildAuthenticationReject(ue *amf_context.AmfUe, eapMsg string) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationReject)

	authenticationReject := nasMessage.NewAuthenticationReject(0)
	authenticationReject.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	authenticationReject.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationReject.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationReject.AuthenticationRejectMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationReject)

	if eapMsg != "" {
		rawEapMsg, _ := base64.StdEncoding.DecodeString(eapMsg)
		authenticationReject.EAPMessage = nasType.NewEAPMessage(nasMessage.AuthenticationRejectEAPMessageType)
		authenticationReject.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		authenticationReject.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.AuthenticationReject = authenticationReject

	return m.PlainNasEncode()
}

func BuildAuthenticationResult(ue *amf_context.AmfUe, eapSuccess bool, eapMsg string) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationResult)

	authenticationResult := nasMessage.NewAuthenticationResult(0)
	authenticationResult.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	authenticationResult.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationResult.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationResult.AuthenticationResultMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationResult)
	authenticationResult.SpareHalfOctetAndNgksi = nasConvert.SpareHalfOctetAndNgksiToNas(ue.NgKsi)
	rawEapMsg, _ := base64.StdEncoding.DecodeString(eapMsg)
	authenticationResult.EAPMessage.SetLen(uint16(len(rawEapMsg)))
	authenticationResult.EAPMessage.SetEAPMessage(rawEapMsg)

	if eapSuccess {
		authenticationResult.ABBA = nasType.NewABBA(nasMessage.AuthenticationResultABBAType)
		authenticationResult.ABBA.SetLen(uint8(len(ue.ABBA)))
		authenticationResult.ABBA.SetABBAContents(ue.ABBA)
	}

	m.GmmMessage.AuthenticationResult = authenticationResult

	return m.PlainNasEncode()
}

// T3346 Timer and EAP are not Supported
func BuildServiceReject(pDUSessionStatus *[16]bool, cause uint8) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeServiceReject)

	serviceReject := nasMessage.NewServiceReject(0)
	serviceReject.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	serviceReject.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	serviceReject.SetMessageType(nas.MsgTypeServiceReject)
	serviceReject.SetCauseValue(cause)
	if pDUSessionStatus != nil {
		serviceReject.PDUSessionStatus = new(nasType.PDUSessionStatus)
		serviceReject.PDUSessionStatus.SetIei(nasMessage.ServiceAcceptPDUSessionStatusType)
		serviceReject.PDUSessionStatus.SetLen(2)
		serviceReject.PDUSessionStatus.Buffer = nasConvert.PSIToBuf(*pDUSessionStatus)
	}

	m.GmmMessage.ServiceReject = serviceReject

	return m.PlainNasEncode()
}

// T3346 timer are not supported
func BuildRegistrationReject(ue *amf_context.AmfUe, cause5GMM uint8, eapMessage string) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationReject)

	registrationReject := nasMessage.NewRegistrationReject(0)
	registrationReject.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationReject.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	registrationReject.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	registrationReject.RegistrationRejectMessageIdentity.SetMessageType(nas.MsgTypeRegistrationReject)
	registrationReject.Cause5GMM.SetCauseValue(cause5GMM)

	if ue.T3502Value != 0 {
		registrationReject.T3502Value = nasType.NewT3502Value(nasMessage.RegistrationRejectT3502ValueType)
		registrationReject.T3502Value.SetLen(1)
		t3502 := nasConvert.GPRSTimer2ToNas(ue.T3502Value)
		registrationReject.T3502Value.SetGPRSTimer2Value(t3502)
	}

	if eapMessage != "" {
		registrationReject.EAPMessage = nasType.NewEAPMessage(nasMessage.RegistrationRejectEAPMessageType)
		rawEapMsg, _ := base64.StdEncoding.DecodeString(eapMessage)
		registrationReject.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		registrationReject.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.RegistrationReject = registrationReject

	return m.PlainNasEncode()
}

// TS 24.501 8.2.25
func BuildSecurityModeCommand(ue *amf_context.AmfUe, eapSuccess bool, eapMessage string) ([]byte, error) {

	// Select enc/int algorithm based on ue security capability & amf's policy,
	self := amf_context.AMF_Self()
	ue.SelectSecurityAlg(self.SecurityAlgorithm.IntegrityOrder, self.SecurityAlgorithm.CipheringOrder)
	// Generate KnasEnc, KnasInt
	ue.DerivateAlgKey()

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeSecurityModeCommand)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext,
	}

	securityModeCommand := nasMessage.NewSecurityModeCommand(0)
	securityModeCommand.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	securityModeCommand.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext)
	securityModeCommand.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	securityModeCommand.SecurityModeCommandMessageIdentity.SetMessageType(nas.MsgTypeSecurityModeCommand)

	securityModeCommand.SelectedNASSecurityAlgorithms.SetTypeOfCipheringAlgorithm(ue.CipheringAlg)
	securityModeCommand.SelectedNASSecurityAlgorithms.SetTypeOfIntegrityProtectionAlgorithm(ue.IntegrityAlg)

	securityModeCommand.SpareHalfOctetAndNgksi = nasConvert.SpareHalfOctetAndNgksiToNas(ue.NgKsi)

	securityModeCommand.ReplayedUESecurityCapabilities.SetLen(ue.NasUESecurityCapability.GetLen())
	securityModeCommand.ReplayedUESecurityCapabilities.Buffer = ue.NasUESecurityCapability.Buffer

	if ue.Pei != "" {
		securityModeCommand.IMEISVRequest = nasType.NewIMEISVRequest(nasMessage.SecurityModeCommandIMEISVRequestType)
		securityModeCommand.IMEISVRequest.SetIMEISVRequestValue(nasMessage.IMEISVNotRequested)
	} else {
		securityModeCommand.IMEISVRequest = nasType.NewIMEISVRequest(nasMessage.SecurityModeCommandIMEISVRequestType)
		securityModeCommand.IMEISVRequest.SetIMEISVRequestValue(nasMessage.IMEISVRequested)
	}

	securityModeCommand.Additional5GSecurityInformation = nasType.NewAdditional5GSecurityInformation(nasMessage.SecurityModeCommandAdditional5GSecurityInformationType)
	securityModeCommand.Additional5GSecurityInformation.SetLen(1)
	securityModeCommand.Additional5GSecurityInformation.SetRINMR(1)
	if ue.Kamf != "" {
		securityModeCommand.Additional5GSecurityInformation.SetHDP(1)
	} else {
		securityModeCommand.Additional5GSecurityInformation.SetHDP(0)
	}

	if eapMessage != "" {
		securityModeCommand.EAPMessage = nasType.NewEAPMessage(nasMessage.SecurityModeCommandEAPMessageType)
		rawEapMsg, _ := base64.StdEncoding.DecodeString(eapMessage)
		securityModeCommand.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		securityModeCommand.EAPMessage.SetEAPMessage(rawEapMsg)

		if eapSuccess {
			securityModeCommand.ABBA = nasType.NewABBA(nasMessage.SecurityModeCommandABBAType)
			securityModeCommand.ABBA.SetLen(uint8(len(ue.ABBA)))
			securityModeCommand.ABBA.SetABBAContents(ue.ABBA)
		}
	}

	m.GmmMessage.SecurityModeCommand = securityModeCommand
	return nas_security.Encode(ue, m)
}

// T3346 timer are not supported
func BuildDeregistrationRequest(ue *amf_context.RanUe, accessType uint8, reRegistrationRequired bool, cause5GMM uint8) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationRequestUETerminatedDeregistration)

	deregistrationRequest := nasMessage.NewDeregistrationRequestUETerminatedDeregistration(0)
	deregistrationRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	deregistrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deregistrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	deregistrationRequest.DeregistrationRequestMessageIdentity.SetMessageType(nas.MsgTypeDeregistrationRequestUETerminatedDeregistration)

	deregistrationRequest.SpareHalfOctetAndDeregistrationType.SetAccessType(accessType)
	deregistrationRequest.SpareHalfOctetAndDeregistrationType.SetSwitchOff(0)
	if reRegistrationRequired {
		deregistrationRequest.SpareHalfOctetAndDeregistrationType.SetReRegistrationRequired(nasMessage.ReRegistrationRequired)
	} else {
		deregistrationRequest.SpareHalfOctetAndDeregistrationType.SetReRegistrationRequired(nasMessage.ReRegistrationNotRequired)
	}

	if cause5GMM != 0 {
		deregistrationRequest.Cause5GMM = nasType.NewCause5GMM(nasMessage.DeregistrationRequestUETerminatedDeregistrationCause5GMMType)
		deregistrationRequest.Cause5GMM.SetCauseValue(cause5GMM)
	}
	m.GmmMessage.DeregistrationRequestUETerminatedDeregistration = deregistrationRequest

	if ue != nil && ue.AmfUe != nil {
		m.SecurityHeader = nas.SecurityHeader{
			ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
			SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
		}
		m.GmmMessage.DeregistrationRequestUETerminatedDeregistration.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypeIntegrityProtectedAndCiphered)
		return nas_security.Encode(ue.AmfUe, m)
	}
	return m.PlainNasEncode()
}

func BuildDeregistrationAccept() ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration)

	deregistrationAccept := nasMessage.NewDeregistrationAcceptUEOriginatingDeregistration(0)
	deregistrationAccept.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	deregistrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deregistrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	deregistrationAccept.DeregistrationAcceptMessageIdentity.SetMessageType(nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration)

	m.GmmMessage.DeregistrationAcceptUEOriginatingDeregistration = deregistrationAccept

	return m.PlainNasEncode()
}

func BuildRegistrationAccept(
	ue *amf_context.AmfUe,
	anType models.AccessType,
	pDUSessionStatus *[16]bool,
	reactivationResult *[16]bool,
	errPduSessionId, errCause []uint8) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationAccept)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	registrationAccept := nasMessage.NewRegistrationAccept(0)
	registrationAccept.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypeIntegrityProtectedAndCiphered)
	registrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	registrationAccept.RegistrationAcceptMessageIdentity.SetMessageType(nas.MsgTypeRegistrationAccept)

	registrationAccept.RegistrationResult5GS.SetLen(1)
	registrationResult := uint8(0)
	if anType == models.AccessType__3_GPP_ACCESS {
		registrationResult |= nasMessage.AccessType3GPP
		if ue.Sm[models.AccessType_NON_3_GPP_ACCESS].Check(gmm_state.REGISTERED) {
			registrationResult |= nasMessage.AccessTypeNon3GPP
		}
	} else {
		registrationResult |= nasMessage.AccessTypeNon3GPP
		if ue.Sm[models.AccessType__3_GPP_ACCESS].Check(gmm_state.REGISTERED) {
			registrationResult |= nasMessage.AccessType3GPP
		}
	}
	registrationAccept.RegistrationResult5GS.SetRegistrationResultValue5GS(registrationResult)
	// TODO: set smsAllowed value of RegistrationResult5GS if need

	if ue.Guti != "" {
		gutiNas := nasConvert.GutiToNas(ue.Guti)
		registrationAccept.GUTI5G = &gutiNas
		registrationAccept.GUTI5G.SetIei(nasMessage.RegistrationAcceptGUTI5GType)
	}

	amfSelf := amf_context.AMF_Self()
	if len(amfSelf.PlmnSupportList) > 1 {
		registrationAccept.EquivalentPlmns = nasType.NewEquivalentPlmns(nasMessage.RegistrationAcceptEquivalentPlmnsType)
		var buf []uint8
		for _, plmnSupportItem := range amfSelf.PlmnSupportList {
			buf = append(buf, nasConvert.PlmnIDToNas(plmnSupportItem.PlmnId)...)
		}
		registrationAccept.EquivalentPlmns.SetLen(uint8(len(buf)))
		copy(registrationAccept.EquivalentPlmns.Octet[:], buf)
	}

	if len(ue.RegistrationArea[anType]) > 0 {
		registrationAccept.TAIList = nasType.NewTAIList(nasMessage.RegistrationAcceptTAIListType)
		taiListNas := nasConvert.TaiListToNas(ue.RegistrationArea[anType])
		registrationAccept.TAIList.SetLen(uint8(len(taiListNas)))
		registrationAccept.TAIList.SetPartialTrackingAreaIdentityList(taiListNas)
	}

	if len(ue.AllowedNssai[anType]) > 0 {
		registrationAccept.AllowedNSSAI = nasType.NewAllowedNSSAI(nasMessage.RegistrationAcceptAllowedNSSAIType)
		var buf []uint8
		for _, snssai := range ue.AllowedNssai[anType] {
			buf = append(buf, nasConvert.SnssaiToNas(snssai)...)
		}
		registrationAccept.AllowedNSSAI.SetLen(uint8(len(buf)))
		registrationAccept.AllowedNSSAI.SetSNSSAIValue(buf)
	}

	if len(ue.RejectedNssai[anType]) > 0 {
		rejectedNssaiNas := nasConvert.RejectedNssaiToNas(ue.RejectedNssai[anType], ue.RejectCause)
		registrationAccept.RejectedNSSAI = &rejectedNssaiNas
		registrationAccept.RejectedNSSAI.SetIei(nasMessage.RegistrationAcceptRejectedNSSAIType)
	}

	if len(ue.ConfiguredNssai[anType]) > 0 {
		registrationAccept.ConfiguredNSSAI = nasType.NewConfiguredNSSAI(nasMessage.RegistrationAcceptConfiguredNSSAIType)
		var buf []uint8
		for _, snssai := range ue.ConfiguredNssai[anType] {
			buf = append(buf, nasConvert.SnssaiToNas(snssai)...)
		}
		registrationAccept.ConfiguredNSSAI.SetLen(uint8(len(buf)))
		registrationAccept.ConfiguredNSSAI.SetSNSSAIValue(buf)
	}

	// TODO: 5gs network feature support

	if pDUSessionStatus != nil {
		registrationAccept.PDUSessionStatus = nasType.NewPDUSessionStatus(nasMessage.RegistrationAcceptPDUSessionStatusType)
		registrationAccept.PDUSessionStatus.SetLen(2)
		registrationAccept.PDUSessionStatus.Buffer = nasConvert.PSIToBuf(*pDUSessionStatus)
	}

	if reactivationResult != nil {
		registrationAccept.PDUSessionReactivationResult = nasType.NewPDUSessionReactivationResult(nasMessage.RegistrationAcceptPDUSessionReactivationResultType)
		registrationAccept.PDUSessionReactivationResult.SetLen(2)
		registrationAccept.PDUSessionReactivationResult.Buffer = nasConvert.PSIToBuf(*reactivationResult)
	}

	if errPduSessionId != nil {
		registrationAccept.PDUSessionReactivationResultErrorCause = nasType.NewPDUSessionReactivationResultErrorCause(nasMessage.RegistrationAcceptPDUSessionReactivationResultErrorCauseType)
		buf := nasConvert.PDUSessionReactivationResultErrorCauseToBuf(errPduSessionId, errCause)
		registrationAccept.PDUSessionReactivationResultErrorCause.SetLen(uint16(len(buf)))
		registrationAccept.PDUSessionReactivationResultErrorCause.Buffer = buf
	}

	if len(ue.LadnInfo) > 0 {
		registrationAccept.LADNInformation = nasType.NewLADNInformation(nasMessage.RegistrationAcceptLADNInformationType)
		var buf []uint8
		for _, ladn := range ue.LadnInfo {
			ladnNas := nasConvert.LadnToNas(ladn)
			buf = append(buf, ladnNas...)
		}
		registrationAccept.LADNInformation.SetLen(uint16(len(buf)))
		registrationAccept.LADNInformation.SetLADND(buf)
	}

	if ue.NetworkSlicingSubscriptionChanged {
		registrationAccept.NetworkSlicingIndication = nasType.NewNetworkSlicingIndication(nasMessage.RegistrationAcceptNetworkSlicingIndicationType)
		registrationAccept.NetworkSlicingIndication.SetNSSCI(1)
		registrationAccept.NetworkSlicingIndication.SetDCNI(0)
		ue.NetworkSlicingSubscriptionChanged = false // reset the value
	}

	if anType == models.AccessType__3_GPP_ACCESS && ue.AmPolicyAssociation != nil && ue.AmPolicyAssociation.ServAreaRes != nil {
		registrationAccept.ServiceAreaList = nasType.NewServiceAreaList(nasMessage.RegistrationAcceptServiceAreaListType)
		partialServiceAreaList := nasConvert.PartialServiceAreaListToNas(ue.PlmnId, *ue.AmPolicyAssociation.ServAreaRes)
		registrationAccept.ServiceAreaList.SetLen(uint8(len(partialServiceAreaList)))
		registrationAccept.ServiceAreaList.SetPartialServiceAreaList(partialServiceAreaList)
	}

	if anType == models.AccessType__3_GPP_ACCESS && ue.T3512Value != 0 {
		registrationAccept.T3512Value = nasType.NewT3512Value(nasMessage.RegistrationAcceptT3512ValueType)
		registrationAccept.T3512Value.SetLen(1)
		t3512 := nasConvert.GPRSTimer3ToNas(ue.T3512Value)
		registrationAccept.T3512Value.Octet = t3512
	}

	if anType == models.AccessType_NON_3_GPP_ACCESS {
		registrationAccept.Non3GppDeregistrationTimerValue = nasType.NewNon3GppDeregistrationTimerValue(nasMessage.RegistrationAcceptNon3GppDeregistrationTimerValueType)
		registrationAccept.Non3GppDeregistrationTimerValue.SetLen(1)
		timerValue := nasConvert.GPRSTimer2ToNas(ue.Non3gppDeregistrationTimerValue)
		registrationAccept.Non3GppDeregistrationTimerValue.SetGPRSTimer2Value(timerValue)
	}

	if ue.T3502Value != 0 {
		registrationAccept.T3502Value = nasType.NewT3502Value(nasMessage.RegistrationAcceptT3502ValueType)
		registrationAccept.T3502Value.SetLen(1)
		t3502 := nasConvert.GPRSTimer2ToNas(ue.T3502Value)
		registrationAccept.T3502Value.SetGPRSTimer2Value(t3502)
	}

	if ue.UESpecificDRX != nasMessage.DRXValueNotSpecified {
		registrationAccept.NegotiatedDRXParameters = nasType.NewNegotiatedDRXParameters(nasMessage.RegistrationAcceptNegotiatedDRXParametersType)
		registrationAccept.NegotiatedDRXParameters.SetLen(1)
		registrationAccept.NegotiatedDRXParameters.SetDRXValue(ue.UESpecificDRX)
	}

	m.GmmMessage.RegistrationAccept = registrationAccept

	return nas_security.Encode(ue, m)
}

func BuildStatus5GMM(cause uint8) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeStatus5GMM)

	status5GMM := nasMessage.NewStatus5GMM(0)
	status5GMM.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	status5GMM.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	status5GMM.SetMessageType(nas.MsgTypeStatus5GMM)
	status5GMM.SetCauseValue(cause)

	m.GmmMessage.Status5GMM = status5GMM

	return m.PlainNasEncode()
}

func BuildConfigurationUpdateCommand(ue *amf_context.AmfUe, anType models.AccessType, networkSlicingIndication *nasType.NetworkSlicingIndication) ([]byte, error) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeConfigurationUpdateCommand)

	configurationUpdateCommand := nasMessage.NewConfigurationUpdateCommand(0)
	configurationUpdateCommand.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	configurationUpdateCommand.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	configurationUpdateCommand.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	configurationUpdateCommand.ConfigurationUpdateCommandMessageIdentity.SetMessageType(nas.MsgTypeConfigurationUpdateCommand)

	if ue.ConfigurationUpdateIndication.Octet != 0 {
		configurationUpdateCommand.ConfigurationUpdateIndication = nasType.NewConfigurationUpdateIndication(nasMessage.ConfigurationUpdateCommandConfigurationUpdateIndicationType)
		configurationUpdateCommand.ConfigurationUpdateIndication = &ue.ConfigurationUpdateIndication
	}

	if networkSlicingIndication != nil {
		configurationUpdateCommand.NetworkSlicingIndication = nasType.NewNetworkSlicingIndication(nasMessage.ConfigurationUpdateCommandNetworkSlicingIndicationType)
		configurationUpdateCommand.NetworkSlicingIndication = networkSlicingIndication
	}

	if ue.Guti != "" {
		gutiNas := nasConvert.GutiToNas(ue.Guti)
		configurationUpdateCommand.GUTI5G = &gutiNas
		configurationUpdateCommand.GUTI5G.SetIei(nasMessage.ConfigurationUpdateCommandGUTI5GType)
	}

	if len(ue.RegistrationArea[anType]) > 0 {
		configurationUpdateCommand.TAIList = nasType.NewTAIList(nasMessage.ConfigurationUpdateCommandTAIListType)
		taiListNas := nasConvert.TaiListToNas(ue.RegistrationArea[anType])
		configurationUpdateCommand.TAIList.SetLen(uint8(len(taiListNas)))
		configurationUpdateCommand.TAIList.SetPartialTrackingAreaIdentityList(taiListNas)
	}

	if len(ue.AllowedNssai[anType]) > 0 {
		configurationUpdateCommand.AllowedNSSAI = nasType.NewAllowedNSSAI(nasMessage.ConfigurationUpdateCommandAllowedNSSAIType)
		var buf []uint8
		for _, snssai := range ue.AllowedNssai[anType] {
			buf = append(buf, nasConvert.SnssaiToNas(snssai)...)
		}
		configurationUpdateCommand.AllowedNSSAI.SetLen(uint8(len(buf)))
		configurationUpdateCommand.AllowedNSSAI.SetSNSSAIValue(buf)
	}

	if len(ue.ConfiguredNssai[anType]) > 0 {
		configurationUpdateCommand.ConfiguredNSSAI = nasType.NewConfiguredNSSAI(nasMessage.ConfigurationUpdateCommandConfiguredNSSAIType)
		var buf []uint8
		for _, snssai := range ue.ConfiguredNssai[anType] {
			buf = append(buf, nasConvert.SnssaiToNas(snssai)...)
		}
		configurationUpdateCommand.ConfiguredNSSAI.SetLen(uint8(len(buf)))
		configurationUpdateCommand.ConfiguredNSSAI.SetSNSSAIValue(buf)
	}

	if len(ue.RejectedNssai[anType]) > 0 {
		rejectedNssaiNas := nasConvert.RejectedNssaiToNas(ue.RejectedNssai[anType], ue.RejectCause)
		configurationUpdateCommand.RejectedNSSAI = &rejectedNssaiNas
		configurationUpdateCommand.RejectedNSSAI.SetIei(nasMessage.ConfigurationUpdateCommandRejectedNSSAIType)
	}

	// TODO: service area list, UniversalTimeAndLocalTimeZone

	amfSelf := amf_context.AMF_Self()
	if amfSelf.NetworkName.Full != "" {
		fullNetworkName := nasConvert.FullNetworkNameToNas(amfSelf.NetworkName.Full)
		configurationUpdateCommand.FullNameForNetwork = &fullNetworkName
		configurationUpdateCommand.FullNameForNetwork.SetIei(nasMessage.ConfigurationUpdateCommandFullNameForNetworkType)
	}

	if amfSelf.NetworkName.Short != "" {
		shortNetworkName := nasConvert.ShortNetworkNameToNas(amfSelf.NetworkName.Short)
		configurationUpdateCommand.ShortNameForNetwork = &shortNetworkName
		configurationUpdateCommand.ShortNameForNetwork.SetIei(nasMessage.ConfigurationUpdateCommandShortNameForNetworkType)
	}

	if ue.TimeZone != "" {
		localTimeZone := nasConvert.LocalTimeZoneToNas(ue.TimeZone)
		localTimeZone.SetIei(nasMessage.ConfigurationUpdateCommandLocalTimeZoneType)
		configurationUpdateCommand.LocalTimeZone = nasType.NewLocalTimeZone(nasMessage.ConfigurationUpdateCommandLocalTimeZoneType)
		configurationUpdateCommand.LocalTimeZone = &localTimeZone
	}

	if ue.TimeZone != "" {
		daylightSavingTime := nasConvert.DaylightSavingTimeToNas(ue.TimeZone)
		daylightSavingTime.SetIei(nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType)
		configurationUpdateCommand.NetworkDaylightSavingTime = nasType.NewNetworkDaylightSavingTime(nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType)
		configurationUpdateCommand.NetworkDaylightSavingTime = &daylightSavingTime
	}

	if len(ue.LadnInfo) > 0 {
		configurationUpdateCommand.LADNInformation = nasType.NewLADNInformation(nasMessage.ConfigurationUpdateCommandLADNInformationType)
		var buf []uint8
		for _, ladn := range ue.LadnInfo {
			ladnNas := nasConvert.LadnToNas(ladn)
			buf = append(buf, ladnNas...)
		}
		configurationUpdateCommand.LADNInformation.SetLen(uint16(len(buf)))
		configurationUpdateCommand.LADNInformation.SetLADND(buf)
	}

	m.GmmMessage.ConfigurationUpdateCommand = configurationUpdateCommand

	return m.PlainNasEncode()
}
