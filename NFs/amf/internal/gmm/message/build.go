package message

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/nas/nas_security"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
)

func BuildDLNASTransport(ue *context.AmfUe, accessType models.AccessType, payloadContainerType uint8, nasPdu []byte,
	pduSessionId uint8, cause *uint8, backoffTimerUint *uint8, backoffTimer uint8,
) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDLNASTransport)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	dLNASTransport := nasMessage.NewDLNASTransport(0)
	dLNASTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	dLNASTransport.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	dLNASTransport.SetMessageType(nas.MsgTypeDLNASTransport)
	dLNASTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(payloadContainerType)
	dLNASTransport.PayloadContainer.SetLen(uint16(len(nasPdu)))
	dLNASTransport.PayloadContainer.SetPayloadContainerContents(nasPdu)

	if pduSessionId != 0 {
		dLNASTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
		dLNASTransport.PduSessionID2Value.SetIei(nasMessage.DLNASTransportPduSessionID2ValueType)
		dLNASTransport.PduSessionID2Value.SetPduSessionID2Value(pduSessionId)
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

	return nas_security.Encode(ue, m, accessType)
}

func BuildNotification(ue *context.AmfUe, accessType models.AccessType) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeNotification)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	notification := nasMessage.NewNotification(0)
	notification.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	notification.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	notification.SetMessageType(nas.MsgTypeNotification)
	if accessType == models.AccessType__3_GPP_ACCESS {
		notification.SetAccessType(nasMessage.AccessType3GPP)
	} else {
		notification.SetAccessType(nasMessage.AccessTypeNon3GPP)
	}

	m.GmmMessage.Notification = notification

	return nas_security.Encode(ue, m, accessType)
}

func BuildIdentityRequest(ue *context.AmfUe, accessType models.AccessType, typeOfIdentity uint8) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeIdentityRequest)

	if ue.SecurityContextAvailable {
		m.SecurityHeader = nas.SecurityHeader{
			ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
			SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
		}
	}

	identityRequest := nasMessage.NewIdentityRequest(0)
	identityRequest.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	identityRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	identityRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	identityRequest.IdentityRequestMessageIdentity.SetMessageType(nas.MsgTypeIdentityRequest)
	identityRequest.SpareHalfOctetAndIdentityType.SetTypeOfIdentity(typeOfIdentity)

	m.GmmMessage.IdentityRequest = identityRequest

	return nas_security.Encode(ue, m, accessType)
}

func BuildAuthenticationRequest(ue *context.AmfUe, accessType models.AccessType) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationRequest)

	authenticationRequest := nasMessage.NewAuthenticationRequest(0)
	authenticationRequest.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	authenticationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationRequest.AuthenticationRequestMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationRequest)
	authenticationRequest.SpareHalfOctetAndNgksi = nasConvert.SpareHalfOctetAndNgksiToNas(ue.NgKsi)
	authenticationRequest.ABBA.SetLen(uint8(len(ue.ABBA)))
	authenticationRequest.ABBA.SetABBAContents(ue.ABBA)

	switch ue.AuthenticationCtx.AuthType {
	case models.AusfUeAuthenticationAuthType__5_G_AKA:
		var tmpArray [16]byte
		var av5gAka models.Av5gAka

		if err := mapstructure.Decode(ue.AuthenticationCtx.Var5gAuthData, &av5gAka); err != nil {
			logger.GmmLog.Error("Var5gAuthData Convert Type Error")
			return nil, err
		}

		rand, err := hex.DecodeString(av5gAka.Rand)
		if err != nil {
			return nil, err
		}
		authenticationRequest.AuthenticationParameterRAND = nasType.
			NewAuthenticationParameterRAND(nasMessage.AuthenticationRequestAuthenticationParameterRANDType)
		copy(tmpArray[:], rand[0:16])
		authenticationRequest.AuthenticationParameterRAND.SetRANDValue(tmpArray)

		autn, err := hex.DecodeString(av5gAka.Autn)
		if err != nil {
			return nil, err
		}
		authenticationRequest.AuthenticationParameterAUTN = nasType.
			NewAuthenticationParameterAUTN(nasMessage.AuthenticationRequestAuthenticationParameterAUTNType)
		authenticationRequest.AuthenticationParameterAUTN.SetLen(uint8(len(autn)))
		copy(tmpArray[:], autn[0:16])
		authenticationRequest.AuthenticationParameterAUTN.SetAUTN(tmpArray)
	case models.AusfUeAuthenticationAuthType_EAP_AKA_PRIME:
		eapMsg := ue.AuthenticationCtx.Var5gAuthData.(string)
		rawEapMsg, err := base64.StdEncoding.DecodeString(eapMsg)
		if err != nil {
			return nil, err
		}
		authenticationRequest.EAPMessage = nasType.NewEAPMessage(nasMessage.AuthenticationRequestEAPMessageType)
		authenticationRequest.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		authenticationRequest.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.AuthenticationRequest = authenticationRequest

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	return nas_security.Encode(ue, m, accessType)
}

func BuildServiceAccept(ue *context.AmfUe, accessType models.AccessType, pDUSessionStatus *[16]bool,
	reactivationResult *[16]bool, errPduSessionId, errCause []uint8,
) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeServiceAccept)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	serviceAccept := nasMessage.NewServiceAccept(0)
	serviceAccept.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	serviceAccept.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
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
		serviceAccept.PDUSessionReactivationResultErrorCause.SetIei(
			nasMessage.ServiceAcceptPDUSessionReactivationResultErrorCauseType)
		buf := nasConvert.PDUSessionReactivationResultErrorCauseToBuf(errPduSessionId, errCause)
		serviceAccept.PDUSessionReactivationResultErrorCause.SetLen(uint16(len(buf)))
		serviceAccept.PDUSessionReactivationResultErrorCause.Buffer = buf
	}
	m.GmmMessage.ServiceAccept = serviceAccept

	return nas_security.Encode(ue, m, accessType)
}

func BuildAuthenticationReject(ue *context.AmfUe, accessType models.AccessType, eapMsg string) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationReject)

	authenticationReject := nasMessage.NewAuthenticationReject(0)
	authenticationReject.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	authenticationReject.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationReject.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationReject.AuthenticationRejectMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationReject)

	if eapMsg != "" {
		rawEapMsg, err := base64.StdEncoding.DecodeString(eapMsg)
		if err != nil {
			return nil, err
		}
		authenticationReject.EAPMessage = nasType.NewEAPMessage(nasMessage.AuthenticationRejectEAPMessageType)
		authenticationReject.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		authenticationReject.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.AuthenticationReject = authenticationReject

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	return nas_security.Encode(ue, m, accessType)
}

func BuildAuthenticationResult(ue *context.AmfUe, accessType models.AccessType, eapSuccess bool, eapMsg string,
) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationResult)

	authenticationResult := nasMessage.NewAuthenticationResult(0)
	authenticationResult.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	authenticationResult.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationResult.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationResult.AuthenticationResultMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationResult)
	authenticationResult.SpareHalfOctetAndNgksi = nasConvert.SpareHalfOctetAndNgksiToNas(ue.NgKsi)
	rawEapMsg, err := base64.StdEncoding.DecodeString(eapMsg)
	if err != nil {
		return nil, err
	}
	authenticationResult.EAPMessage.SetLen(uint16(len(rawEapMsg)))
	authenticationResult.EAPMessage.SetEAPMessage(rawEapMsg)

	if eapSuccess {
		authenticationResult.ABBA = nasType.NewABBA(nasMessage.AuthenticationResultABBAType)
		authenticationResult.ABBA.SetLen(uint8(len(ue.ABBA)))
		authenticationResult.ABBA.SetABBAContents(ue.ABBA)
	}

	m.GmmMessage.AuthenticationResult = authenticationResult

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	return nas_security.Encode(ue, m, accessType)
}

// T3346 Timer and EAP are not Supported
func BuildServiceReject(ue *context.AmfUe, accessType models.AccessType, pDUSessionStatus *[16]bool, cause uint8,
) ([]byte, error) {
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

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	return nas_security.Encode(ue, m, accessType)
}

// T3346 timer are not supported
func BuildRegistrationReject(ue *context.AmfUe, accessType models.AccessType, cause5GMM uint8, eapMessage string,
) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationReject)

	registrationReject := nasMessage.NewRegistrationReject(0)
	registrationReject.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationReject.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	registrationReject.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	registrationReject.RegistrationRejectMessageIdentity.SetMessageType(nas.MsgTypeRegistrationReject)
	registrationReject.Cause5GMM.SetCauseValue(cause5GMM)

	t3502Val := context.GetSelf().T3502Value
	if ue != nil {
		t3502Val = ue.T3502Value
	}
	if t3502Val != 0 {
		registrationReject.T3502Value = nasType.NewT3502Value(nasMessage.RegistrationRejectT3502ValueType)
		registrationReject.T3502Value.SetLen(1)
		t3502 := nasConvert.GPRSTimer2ToNas(t3502Val)
		registrationReject.T3502Value.SetGPRSTimer2Value(t3502)
	}

	if eapMessage != "" {
		registrationReject.EAPMessage = nasType.NewEAPMessage(nasMessage.RegistrationRejectEAPMessageType)
		rawEapMsg, err := base64.StdEncoding.DecodeString(eapMessage)
		if err != nil {
			return nil, err
		}
		registrationReject.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		registrationReject.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.RegistrationReject = registrationReject

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	return nas_security.Encode(ue, m, accessType)
}

// TS 24.501 8.2.25
func BuildSecurityModeCommand(ue *context.AmfUe, accessType models.AccessType, eapSuccess bool, eapMessage string) (
	[]byte, error,
) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeSecurityModeCommand)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext,
	}

	securityModeCommand := nasMessage.NewSecurityModeCommand(0)
	securityModeCommand.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	securityModeCommand.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	securityModeCommand.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	securityModeCommand.SecurityModeCommandMessageIdentity.SetMessageType(nas.MsgTypeSecurityModeCommand)

	securityModeCommand.SelectedNASSecurityAlgorithms.SetTypeOfCipheringAlgorithm(ue.CipheringAlg)
	securityModeCommand.SelectedNASSecurityAlgorithms.SetTypeOfIntegrityProtectionAlgorithm(ue.IntegrityAlg)

	securityModeCommand.SpareHalfOctetAndNgksi = nasConvert.SpareHalfOctetAndNgksiToNas(ue.NgKsi)

	securityModeCommand.ReplayedUESecurityCapabilities.SetLen(ue.UESecurityCapability.GetLen())
	securityModeCommand.ReplayedUESecurityCapabilities.Buffer = ue.UESecurityCapability.Buffer

	if ue.Pei != "" {
		securityModeCommand.IMEISVRequest = nasType.NewIMEISVRequest(nasMessage.SecurityModeCommandIMEISVRequestType)
		securityModeCommand.IMEISVRequest.SetIMEISVRequestValue(nasMessage.IMEISVNotRequested)
	} else {
		securityModeCommand.IMEISVRequest = nasType.NewIMEISVRequest(nasMessage.SecurityModeCommandIMEISVRequestType)
		securityModeCommand.IMEISVRequest.SetIMEISVRequestValue(nasMessage.IMEISVRequested)
	}

	securityModeCommand.Additional5GSecurityInformation = nasType.
		NewAdditional5GSecurityInformation(nasMessage.SecurityModeCommandAdditional5GSecurityInformationType)
	securityModeCommand.Additional5GSecurityInformation.SetLen(1)
	if ue.RetransmissionOfInitialNASMsg {
		securityModeCommand.Additional5GSecurityInformation.SetRINMR(1)
	} else {
		securityModeCommand.Additional5GSecurityInformation.SetRINMR(0)
	}

	if ue.RegistrationType5GS == nasMessage.RegistrationType5GSPeriodicRegistrationUpdating ||
		ue.RegistrationType5GS == nasMessage.RegistrationType5GSMobilityRegistrationUpdating {
		securityModeCommand.Additional5GSecurityInformation.SetHDP(1)
	} else {
		securityModeCommand.Additional5GSecurityInformation.SetHDP(0)
	}

	if eapMessage != "" {
		securityModeCommand.EAPMessage = nasType.NewEAPMessage(nasMessage.SecurityModeCommandEAPMessageType)
		rawEapMsg, err := base64.StdEncoding.DecodeString(eapMessage)
		if err != nil {
			return nil, err
		}
		securityModeCommand.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		securityModeCommand.EAPMessage.SetEAPMessage(rawEapMsg)

		if eapSuccess {
			securityModeCommand.ABBA = nasType.NewABBA(nasMessage.SecurityModeCommandABBAType)
			securityModeCommand.ABBA.SetLen(uint8(len(ue.ABBA)))
			securityModeCommand.ABBA.SetABBAContents(ue.ABBA)
		}
	}

	ue.SecurityContextAvailable = true
	m.GmmMessage.SecurityModeCommand = securityModeCommand
	payload, err := nas_security.Encode(ue, m, accessType)
	if err != nil {
		ue.SecurityContextAvailable = false
		return nil, err
	} else {
		return payload, nil
	}
}

// T3346 timer are not supported
func BuildDeregistrationRequest(ue *context.RanUe, accessType uint8, reRegistrationRequired bool,
	cause5GMM uint8,
) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationRequestUETerminatedDeregistration)

	deregistrationRequest := nasMessage.NewDeregistrationRequestUETerminatedDeregistration(0)
	deregistrationRequest.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	deregistrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deregistrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	deregistrationRequest.SetMessageType(nas.MsgTypeDeregistrationRequestUETerminatedDeregistration)

	deregistrationRequest.SetAccessType(accessType)
	deregistrationRequest.SetSwitchOff(0)
	if reRegistrationRequired {
		deregistrationRequest.SetReRegistrationRequired(nasMessage.ReRegistrationRequired)
	} else {
		deregistrationRequest.SetReRegistrationRequired(nasMessage.ReRegistrationNotRequired)
	}

	if cause5GMM != 0 {
		deregistrationRequest.Cause5GMM = nasType.NewCause5GMM(
			nasMessage.DeregistrationRequestUETerminatedDeregistrationCause5GMMType)
		deregistrationRequest.Cause5GMM.SetCauseValue(cause5GMM)
	}
	m.GmmMessage.DeregistrationRequestUETerminatedDeregistration = deregistrationRequest

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	var anType models.AccessType
	switch accessType {
	case 0x01:
		anType = models.AccessType__3_GPP_ACCESS
	case 0x02:
		anType = models.AccessType_NON_3_GPP_ACCESS
	}
	if ue != nil {
		return nas_security.Encode(ue.AmfUe, m, anType)
	} else {
		return nas_security.Encode(nil, m, anType)
	}
}

func BuildDeregistrationAccept(ue *context.AmfUe, accessType models.AccessType) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration)

	deregistrationAccept := nasMessage.NewDeregistrationAcceptUEOriginatingDeregistration(0)
	deregistrationAccept.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	deregistrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deregistrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	deregistrationAccept.SetMessageType(nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration)

	m.GmmMessage.DeregistrationAcceptUEOriginatingDeregistration = deregistrationAccept

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	return nas_security.Encode(ue, m, accessType)
}

func BuildRegistrationAccept(
	ue *context.AmfUe,
	anType models.AccessType,
	pDUSessionStatus *[16]bool,
	reactivationResult *[16]bool,
	errPduSessionId, errCause []uint8,
) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationAccept)

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	registrationAccept := nasMessage.NewRegistrationAccept(0)
	registrationAccept.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	registrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	registrationAccept.RegistrationAcceptMessageIdentity.SetMessageType(nas.MsgTypeRegistrationAccept)

	registrationAccept.RegistrationResult5GS.SetLen(1)
	registrationResult := uint8(0)
	if anType == models.AccessType__3_GPP_ACCESS {
		registrationResult |= nasMessage.AccessType3GPP
		if ue.State[models.AccessType_NON_3_GPP_ACCESS].Is(context.Registered) {
			registrationResult |= nasMessage.AccessTypeNon3GPP
		}
	} else {
		registrationResult |= nasMessage.AccessTypeNon3GPP
		if ue.State[models.AccessType__3_GPP_ACCESS].Is(context.Registered) {
			registrationResult |= nasMessage.AccessType3GPP
		}
	}
	registrationAccept.RegistrationResult5GS.SetRegistrationResultValue5GS(registrationResult)
	// TODO: set smsAllowed value of RegistrationResult5GS if need

	if ue.Guti != "" {
		gutiNas, err := nasConvert.GutiToNasWithError(ue.Guti)
		if err != nil {
			return nil, fmt.Errorf("encode GUTI failed: %w", err)
		}
		registrationAccept.GUTI5G = &gutiNas
		registrationAccept.GUTI5G.SetIei(nasMessage.RegistrationAcceptGUTI5GType)
	}

	amfSelf := context.GetSelf()
	if len(amfSelf.PlmnSupportList) > 1 {
		registrationAccept.EquivalentPlmns = nasType.NewEquivalentPlmns(nasMessage.RegistrationAcceptEquivalentPlmnsType)
		var buf []uint8
		for _, plmnSupportItem := range amfSelf.PlmnSupportList {
			buf = append(buf, nasConvert.PlmnIDToNas(*plmnSupportItem.PlmnId)...)
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
		for _, allowedSnssai := range ue.AllowedNssai[anType] {
			buf = append(buf, nasConvert.SnssaiToNas(*allowedSnssai.AllowedSnssai)...)
		}
		registrationAccept.AllowedNSSAI.SetLen(uint8(len(buf)))
		registrationAccept.AllowedNSSAI.SetSNSSAIValue(buf)
	}

	if ue.NetworkSliceInfo != nil {
		if len(ue.NetworkSliceInfo.RejectedNssaiInPlmn) != 0 || len(ue.NetworkSliceInfo.RejectedNssaiInTa) != 0 {
			rejectedNssaiNas := nasConvert.RejectedNssaiToNas(
				ue.NetworkSliceInfo.RejectedNssaiInPlmn, ue.NetworkSliceInfo.RejectedNssaiInTa)
			registrationAccept.RejectedNSSAI = &rejectedNssaiNas
			registrationAccept.RejectedNSSAI.SetIei(nasMessage.RegistrationAcceptRejectedNSSAIType)
		}
	}

	if includeConfiguredNssaiCheck(ue) {
		registrationAccept.ConfiguredNSSAI = nasType.NewConfiguredNSSAI(nasMessage.RegistrationAcceptConfiguredNSSAIType)
		var buf []uint8
		for _, snssai := range ue.ConfiguredNssai {
			buf = append(buf, nasConvert.SnssaiToNas(*snssai.ConfiguredSnssai)...)
		}
		registrationAccept.ConfiguredNSSAI.SetLen(uint8(len(buf)))
		registrationAccept.ConfiguredNSSAI.SetSNSSAIValue(buf)
	}

	// 5gs network feature support
	if c := factory.AmfConfig.GetNasIENetworkFeatureSupport5GS(); c != nil && c.Enable {
		registrationAccept.NetworkFeatureSupport5GS = nasType.
			NewNetworkFeatureSupport5GS(nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType)
		registrationAccept.NetworkFeatureSupport5GS.SetLen(c.Length)
		if anType == models.AccessType__3_GPP_ACCESS {
			registrationAccept.SetIMSVoPS3GPP(c.ImsVoPS)
		} else {
			registrationAccept.SetIMSVoPSN3GPP(c.ImsVoPS)
		}
		registrationAccept.SetEMC(c.Emc)
		registrationAccept.SetEMF(c.Emf)
		registrationAccept.SetIWKN26(c.IwkN26)
		registrationAccept.SetMPSI(c.Mpsi)
		registrationAccept.SetEMCN(c.EmcN3)
		registrationAccept.SetMCSI(c.Mcsi)
	}

	if pDUSessionStatus != nil {
		registrationAccept.PDUSessionStatus = nasType.NewPDUSessionStatus(nasMessage.RegistrationAcceptPDUSessionStatusType)
		registrationAccept.PDUSessionStatus.SetLen(2)
		registrationAccept.PDUSessionStatus.Buffer = nasConvert.PSIToBuf(*pDUSessionStatus)
	}

	if reactivationResult != nil {
		registrationAccept.PDUSessionReactivationResult = nasType.
			NewPDUSessionReactivationResult(nasMessage.RegistrationAcceptPDUSessionReactivationResultType)
		registrationAccept.PDUSessionReactivationResult.SetLen(2)
		registrationAccept.PDUSessionReactivationResult.Buffer = nasConvert.PSIToBuf(*reactivationResult)
	}

	if errPduSessionId != nil {
		registrationAccept.PDUSessionReactivationResultErrorCause = nasType.NewPDUSessionReactivationResultErrorCause(
			nasMessage.RegistrationAcceptPDUSessionReactivationResultErrorCauseType)
		buf := nasConvert.PDUSessionReactivationResultErrorCauseToBuf(errPduSessionId, errCause)
		registrationAccept.PDUSessionReactivationResultErrorCause.SetLen(uint16(len(buf)))
		registrationAccept.PDUSessionReactivationResultErrorCause.Buffer = buf
	}

	if ue.LadnInfo != nil {
		registrationAccept.LADNInformation = nasType.NewLADNInformation(nasMessage.RegistrationAcceptLADNInformationType)
		buf := make([]uint8, 0)
		for _, ladn := range ue.LadnInfo {
			ladnNas := nasConvert.LadnToNas(ladn.Dnn, ladn.TaiList)
			buf = append(buf, ladnNas...)
		}
		registrationAccept.LADNInformation.SetLen(uint16(len(buf)))
		registrationAccept.LADNInformation.SetLADND(buf)
	}

	if ue.NetworkSlicingSubscriptionChanged {
		registrationAccept.NetworkSlicingIndication = nasType.
			NewNetworkSlicingIndication(nasMessage.RegistrationAcceptNetworkSlicingIndicationType)
		registrationAccept.NetworkSlicingIndication.SetNSSCI(1)
		registrationAccept.NetworkSlicingIndication.SetDCNI(0)
		ue.NetworkSlicingSubscriptionChanged = false // reset the value
	}

	if anType == models.AccessType__3_GPP_ACCESS && ue.AmPolicyAssociation != nil &&
		ue.AmPolicyAssociation.ServAreaRes != nil {
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
		registrationAccept.Non3GppDeregistrationTimerValue = nasType.
			NewNon3GppDeregistrationTimerValue(nasMessage.RegistrationAcceptNon3GppDeregistrationTimerValueType)
		registrationAccept.Non3GppDeregistrationTimerValue.SetLen(1)
		timerValue := nasConvert.GPRSTimer2ToNas(ue.Non3gppDeregTimerValue)
		registrationAccept.Non3GppDeregistrationTimerValue.SetGPRSTimer2Value(timerValue)
	}

	if ue.T3502Value != 0 {
		registrationAccept.T3502Value = nasType.NewT3502Value(nasMessage.RegistrationAcceptT3502ValueType)
		registrationAccept.T3502Value.SetLen(1)
		t3502 := nasConvert.GPRSTimer2ToNas(ue.T3502Value)
		registrationAccept.T3502Value.SetGPRSTimer2Value(t3502)
	}

	if ue.UESpecificDRX != nasMessage.DRXValueNotSpecified {
		registrationAccept.NegotiatedDRXParameters = nasType.
			NewNegotiatedDRXParameters(nasMessage.RegistrationAcceptNegotiatedDRXParametersType)
		registrationAccept.NegotiatedDRXParameters.SetLen(1)
		registrationAccept.NegotiatedDRXParameters.SetDRXValue(ue.UESpecificDRX)
	}

	m.GmmMessage.RegistrationAccept = registrationAccept

	return nas_security.Encode(ue, m, anType)
}

func includeConfiguredNssaiCheck(ue *context.AmfUe) bool {
	if len(ue.ConfiguredNssai) == 0 {
		return false
	}

	registrationRequest := ue.RegistrationRequest
	if registrationRequest.RequestedNSSAI == nil {
		return true
	}
	if ue.NetworkSliceInfo != nil && len(ue.NetworkSliceInfo.RejectedNssaiInPlmn) != 0 {
		return true
	}
	if registrationRequest.NetworkSlicingIndication != nil && registrationRequest.NetworkSlicingIndication.GetDCNI() == 1 {
		return true
	}
	return false
}

func BuildStatus5GMM(ue *context.AmfUe, accessType models.AccessType, cause uint8) ([]byte, error) {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeStatus5GMM)

	status5GMM := nasMessage.NewStatus5GMM(0)
	status5GMM.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	status5GMM.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	status5GMM.SetMessageType(nas.MsgTypeStatus5GMM)
	status5GMM.SetCauseValue(cause)

	m.GmmMessage.Status5GMM = status5GMM

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	return nas_security.Encode(ue, m, accessType)
}

// Fllowed by TS 24.501 - 5.4.4 Generic UE configuration update procedure - 5.4.4.1 General
func BuildConfigurationUpdateCommand(ue *context.AmfUe, anType models.AccessType,
	flags *context.ConfigurationUpdateCommandFlags,
) ([]byte, error, bool) {
	needTimer := false
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeConfigurationUpdateCommand)

	configurationUpdateCommand := nasMessage.NewConfigurationUpdateCommand(0)
	configurationUpdateCommand.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	configurationUpdateCommand.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	configurationUpdateCommand.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	configurationUpdateCommand.SetMessageType(nas.MsgTypeConfigurationUpdateCommand)

	if flags.NeedNetworkSlicingIndication {
		configurationUpdateCommand.NetworkSlicingIndication = nasType.
			NewNetworkSlicingIndication(nasMessage.ConfigurationUpdateCommandNetworkSlicingIndicationType)
		configurationUpdateCommand.NetworkSlicingIndication.SetNSSCI(0x01)
	}

	if flags.NeedGUTI {
		if ue.Guti != "" {
			gutiNas, err := nasConvert.GutiToNasWithError(ue.Guti)
			if err != nil {
				return nil, fmt.Errorf("encode GUTI failed: %w", err), needTimer
			}
			configurationUpdateCommand.GUTI5G = &gutiNas
			configurationUpdateCommand.GUTI5G.SetIei(nasMessage.ConfigurationUpdateCommandGUTI5GType)
		} else {
			logger.GmmLog.Warnf("Require 5G-GUTI, but got nothing.")
		}
	}

	if flags.NeedAllowedNSSAI {
		if len(ue.AllowedNssai[anType]) > 0 {
			configurationUpdateCommand.AllowedNSSAI = nasType.
				NewAllowedNSSAI(nasMessage.ConfigurationUpdateCommandAllowedNSSAIType)

			var buf []uint8
			for _, allowedSnssai := range ue.AllowedNssai[anType] {
				buf = append(buf, nasConvert.SnssaiToNas(*allowedSnssai.AllowedSnssai)...)
			}
			configurationUpdateCommand.AllowedNSSAI.SetLen(uint8(len(buf)))
			configurationUpdateCommand.AllowedNSSAI.SetSNSSAIValue(buf)
		} else {
			logger.GmmLog.Warnf("Require Allowed NSSAI, but got nothing.")
		}
	}

	if flags.NeedConfiguredNSSAI {
		if len(ue.ConfiguredNssai) > 0 {
			configurationUpdateCommand.ConfiguredNSSAI = nasType.
				NewConfiguredNSSAI(nasMessage.ConfigurationUpdateCommandConfiguredNSSAIType)

			var buf []uint8
			for _, snssai := range ue.ConfiguredNssai {
				buf = append(buf, nasConvert.SnssaiToNas(*snssai.ConfiguredSnssai)...)
			}
			configurationUpdateCommand.ConfiguredNSSAI.SetLen(uint8(len(buf)))
			configurationUpdateCommand.ConfiguredNSSAI.SetSNSSAIValue(buf)
		} else {
			logger.GmmLog.Warnf("Require Configured NSSAI, but got nothing.")
		}
	}

	if flags.NeedRejectNSSAI {
		if ue.NetworkSliceInfo != nil &&
			(len(ue.NetworkSliceInfo.RejectedNssaiInPlmn) != 0 || len(ue.NetworkSliceInfo.RejectedNssaiInTa) != 0) {
			rejectedNssaiNas := nasConvert.RejectedNssaiToNas(
				ue.NetworkSliceInfo.RejectedNssaiInPlmn, ue.NetworkSliceInfo.RejectedNssaiInTa)
			configurationUpdateCommand.RejectedNSSAI = &rejectedNssaiNas
			configurationUpdateCommand.RejectedNSSAI.SetIei(nasMessage.ConfigurationUpdateCommandRejectedNSSAIType)
		} else {
			logger.GmmLog.Warnf("Require Rejected NSSAI, but got nothing.")
		}
	}

	if flags.NeedTaiList && anType == models.AccessType__3_GPP_ACCESS {
		if len(ue.RegistrationArea[anType]) > 0 {
			configurationUpdateCommand.TAIList = nasType.NewTAIList(nasMessage.ConfigurationUpdateCommandTAIListType)
			taiListNas := nasConvert.TaiListToNas(ue.RegistrationArea[anType])
			configurationUpdateCommand.TAIList.SetLen(uint8(len(taiListNas)))
			configurationUpdateCommand.TAIList.SetPartialTrackingAreaIdentityList(taiListNas)
		} else {
			logger.GmmLog.Warnf("Require TAI List, but got nothing.")
		}
	}

	if flags.NeedServiceAreaList && anType == models.AccessType__3_GPP_ACCESS {
		if ue.AmPolicyAssociation != nil && ue.AmPolicyAssociation.ServAreaRes != nil {
			configurationUpdateCommand.ServiceAreaList = nasType.
				NewServiceAreaList(nasMessage.ConfigurationUpdateCommandServiceAreaListType)
			partialServiceAreaList := nasConvert.
				PartialServiceAreaListToNas(ue.PlmnId, *ue.AmPolicyAssociation.ServAreaRes)
			configurationUpdateCommand.ServiceAreaList.SetLen(uint8(len(partialServiceAreaList)))
			configurationUpdateCommand.ServiceAreaList.SetPartialServiceAreaList(partialServiceAreaList)
		} else {
			logger.GmmLog.Warnf("Require Service Area List, but got nothing.")
		}
	}

	if flags.NeedLadnInformation && anType == models.AccessType__3_GPP_ACCESS {
		if len(ue.LadnInfo) > 0 {
			configurationUpdateCommand.LADNInformation = nasType.
				NewLADNInformation(nasMessage.ConfigurationUpdateCommandLADNInformationType)
			var buf []uint8
			for _, ladn := range ue.LadnInfo {
				ladnNas := nasConvert.LadnToNas(ladn.Dnn, ladn.TaiList)
				buf = append(buf, ladnNas...)
			}
			configurationUpdateCommand.LADNInformation.SetLen(uint16(len(buf)))
			configurationUpdateCommand.LADNInformation.SetLADND(buf)
		} else {
			logger.GmmLog.Warnf("Require LADN Information, but got nothing.")
		}
	}

	amfSelf := context.GetSelf()

	if flags.NeedNITZ {
		// Full network name
		if amfSelf.NetworkName.Full != "" {
			fullNetworkName := nasConvert.FullNetworkNameToNas(amfSelf.NetworkName.Full)
			configurationUpdateCommand.FullNameForNetwork = &fullNetworkName
			configurationUpdateCommand.FullNameForNetwork.SetIei(nasMessage.ConfigurationUpdateCommandFullNameForNetworkType)
		} else {
			logger.GmmLog.Warnf("Require Full Network Name, but got nothing.")
		}
		// Short network name
		if amfSelf.NetworkName.Short != "" {
			shortNetworkName := nasConvert.ShortNetworkNameToNas(amfSelf.NetworkName.Short)
			configurationUpdateCommand.ShortNameForNetwork = &shortNetworkName
			configurationUpdateCommand.ShortNameForNetwork.SetIei(nasMessage.ConfigurationUpdateCommandShortNameForNetworkType)
		} else {
			logger.GmmLog.Warnf("Require Short Network Name, but got nothing.")
		}
		// Universal Time and Local Time Zone
		now := time.Now()
		universalTimeAndLocalTimeZone := nasConvert.EncodeUniversalTimeAndLocalTimeZoneToNas(now)
		universalTimeAndLocalTimeZone.SetIei(nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType)
		configurationUpdateCommand.UniversalTimeAndLocalTimeZone = &universalTimeAndLocalTimeZone

		if ue.TimeZone != amfSelf.TimeZone {
			ue.TimeZone = amfSelf.TimeZone
			// Local Time Zone
			localTimeZone := nasConvert.EncodeLocalTimeZoneToNas(ue.TimeZone)
			localTimeZone.SetIei(nasMessage.ConfigurationUpdateCommandLocalTimeZoneType)
			configurationUpdateCommand.LocalTimeZone = nasType.
				NewLocalTimeZone(nasMessage.ConfigurationUpdateCommandLocalTimeZoneType)
			configurationUpdateCommand.LocalTimeZone = &localTimeZone
			// Daylight Saving Time
			daylightSavingTime := nasConvert.EncodeDaylightSavingTimeToNas(ue.TimeZone)
			daylightSavingTime.SetIei(nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType)
			configurationUpdateCommand.NetworkDaylightSavingTime = nasType.
				NewNetworkDaylightSavingTime(nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType)
			configurationUpdateCommand.NetworkDaylightSavingTime = &daylightSavingTime
		}
	}

	configurationUpdateCommand.ConfigurationUpdateIndication = nasType.
		NewConfigurationUpdateIndication(nasMessage.ConfigurationUpdateCommandConfigurationUpdateIndicationType)
	if configurationUpdateCommand.GUTI5G != nil ||
		configurationUpdateCommand.TAIList != nil ||
		configurationUpdateCommand.AllowedNSSAI != nil ||
		configurationUpdateCommand.LADNInformation != nil ||
		configurationUpdateCommand.ServiceAreaList != nil ||
		configurationUpdateCommand.MICOIndication != nil ||
		configurationUpdateCommand.ConfiguredNSSAI != nil ||
		configurationUpdateCommand.RejectedNSSAI != nil ||
		configurationUpdateCommand.NetworkSlicingIndication != nil ||
		configurationUpdateCommand.OperatordefinedAccessCategoryDefinitions != nil ||
		configurationUpdateCommand.SMSIndication != nil {
		// TS 24.501 - 5.4.4.2 Generic UE configuration update procedure initiated by the network
		// Acknowledgement shall be requested for all parameters except when only NITZ is included
		configurationUpdateCommand.ConfigurationUpdateIndication.SetACK(uint8(1))
		needTimer = true
	}
	if configurationUpdateCommand.MICOIndication != nil {
		// Allowed NSSAI and Configured NSSAI are optional to request to perform the registration procedure
		configurationUpdateCommand.ConfigurationUpdateIndication.SetRED(uint8(1))
	}

	// Check if the Configuration Update Command is vaild
	if configurationUpdateCommand.ConfigurationUpdateIndication.GetACK() == uint8(0) &&
		configurationUpdateCommand.ConfigurationUpdateIndication.GetRED() == uint8(0) &&
		(configurationUpdateCommand.FullNameForNetwork == nil &&
			configurationUpdateCommand.ShortNameForNetwork == nil &&
			configurationUpdateCommand.UniversalTimeAndLocalTimeZone == nil &&
			configurationUpdateCommand.LocalTimeZone == nil &&
			configurationUpdateCommand.NetworkDaylightSavingTime == nil) {
		return nil, fmt.Errorf("configuration update command is invalid"), false
	}

	m.GmmMessage.ConfigurationUpdateCommand = configurationUpdateCommand

	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}

	b, err := nas_security.Encode(ue, m, anType)
	if err != nil {
		return nil, fmt.Errorf("BuildConfigurationUpdateCommand() err: %v", err), false
	}
	return b, err, needTimer
}
