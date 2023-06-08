package nasTestpacket

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/logger"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
)

const (
	PDUSesModiReq    string = "PDU Session Modification Request"
	PDUSesModiCmp    string = "PDU Session Modification Complete"
	PDUSesModiCmdRej string = "PDU Session Modification Command Reject"
	PDUSesRelReq     string = "PDU Session Release Request"
	PDUSesRelCmp     string = "PDU Session Release Complete"
	PDUSesRelRej     string = "PDU Session Release Reject"
	PDUSesAuthCmp    string = "PDU Session Authentication Complete"
)

func GetRegistrationRequest(
	registrationType uint8,
	mobileIdentity nasType.MobileIdentity5GS,
	requestedNSSAI *nasType.RequestedNSSAI,
	ueSecurityCapability *nasType.UESecurityCapability,
	capability5GMM *nasType.Capability5GMM,
	nasMessageContainer []uint8,
	uplinkDataStatus *nasType.UplinkDataStatus,
) []byte {
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationRequest)

	registrationRequest := nasMessage.NewRegistrationRequest(0)
	registrationRequest.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	registrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0x00)
	registrationRequest.RegistrationRequestMessageIdentity.SetMessageType(nas.MsgTypeRegistrationRequest)
	registrationRequest.NgksiAndRegistrationType5GS.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	registrationRequest.NgksiAndRegistrationType5GS.SetNasKeySetIdentifiler(0x7)
	registrationRequest.NgksiAndRegistrationType5GS.SetFOR(1)
	registrationRequest.NgksiAndRegistrationType5GS.SetRegistrationType5GS(registrationType)
	registrationRequest.MobileIdentity5GS = mobileIdentity

	registrationRequest.UESecurityCapability = ueSecurityCapability
	registrationRequest.Capability5GMM = capability5GMM
	registrationRequest.RequestedNSSAI = requestedNSSAI
	registrationRequest.UplinkDataStatus = uplinkDataStatus

	if nasMessageContainer != nil {
		registrationRequest.NASMessageContainer = nasType.NewNASMessageContainer(
			nasMessage.RegistrationRequestNASMessageContainerType)
		registrationRequest.NASMessageContainer.SetLen(uint16(len(nasMessageContainer)))
		registrationRequest.NASMessageContainer.SetNASMessageContainerContents(nasMessageContainer)
	}

	m.GmmMessage.RegistrationRequest = registrationRequest

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetPduSessionEstablishmentRequest(pduSessionId uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionEstablishmentRequest)

	pduSessionEstablishmentRequest := nasMessage.NewPDUSessionEstablishmentRequest(0)
	pduSessionEstablishmentRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSSessionManagementMessage)
	pduSessionEstablishmentRequest.SetMessageType(nas.MsgTypePDUSessionEstablishmentRequest)
	pduSessionEstablishmentRequest.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionEstablishmentRequest.PTI.SetPTI(0x00)
	pduSessionEstablishmentRequest.IntegrityProtectionMaximumDataRate.
		SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink(0xff)
	pduSessionEstablishmentRequest.IntegrityProtectionMaximumDataRate.
		SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink(0xff)

	pduSessionEstablishmentRequest.PDUSessionType =
		nasType.NewPDUSessionType(nasMessage.PDUSessionEstablishmentRequestPDUSessionTypeType)
	pduSessionEstablishmentRequest.PDUSessionType.SetPDUSessionTypeValue(uint8(0x01)) //IPv4 type

	pduSessionEstablishmentRequest.SSCMode =
		nasType.NewSSCMode(nasMessage.PDUSessionEstablishmentRequestSSCModeType)
	pduSessionEstablishmentRequest.SSCMode.SetSSCMode(uint8(0x01)) //SSC Mode 1

	pduSessionEstablishmentRequest.ExtendedProtocolConfigurationOptions =
		nasType.NewExtendedProtocolConfigurationOptions(
			nasMessage.PDUSessionEstablishmentRequestExtendedProtocolConfigurationOptionsType)
	protocolConfigurationOptions := nasConvert.NewProtocolConfigurationOptions()
	protocolConfigurationOptions.AddIPAddressAllocationViaNASSignallingUL()
	protocolConfigurationOptions.AddDNSServerIPv4AddressRequest()
	protocolConfigurationOptions.AddDNSServerIPv6AddressRequest()
	pcoContents := protocolConfigurationOptions.Marshal()
	pcoContentsLength := len(pcoContents)
	pduSessionEstablishmentRequest.ExtendedProtocolConfigurationOptions.SetLen(uint16(pcoContentsLength))
	pduSessionEstablishmentRequest.ExtendedProtocolConfigurationOptions.
		SetExtendedProtocolConfigurationOptionsContents(pcoContents)

	m.GsmMessage.PDUSessionEstablishmentRequest = pduSessionEstablishmentRequest

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetUlNasTransport_PduSessionEstablishmentRequest(pduSessionId uint8, requestType uint8, dnnString string,
	sNssai *models.Snssai) []byte {

	pduSessionEstablishmentRequest := GetPduSessionEstablishmentRequest(pduSessionId)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(pduSessionId)
	ulNasTransport.RequestType = new(nasType.RequestType)
	ulNasTransport.RequestType.SetIei(nasMessage.ULNASTransportRequestTypeType)
	ulNasTransport.RequestType.SetRequestTypeValue(requestType)
	if dnnString != "" {
		ulNasTransport.DNN = new(nasType.DNN)
		ulNasTransport.DNN.SetIei(nasMessage.ULNASTransportDNNType)
		ulNasTransport.DNN.SetDNN(dnnString)
	}
	if sNssai != nil {
		var sdTemp [3]uint8
		sd, err := hex.DecodeString(sNssai.Sd)
		if err != nil {
			logger.NasMsgLog.Errorf("sNssai decode error: %+v", err)
		}
		copy(sdTemp[:], sd)
		ulNasTransport.SNSSAI = nasType.NewSNSSAI(nasMessage.ULNASTransportSNSSAIType)
		ulNasTransport.SNSSAI.SetLen(4)
		ulNasTransport.SNSSAI.SetSST(uint8(sNssai.Sst))
		ulNasTransport.SNSSAI.SetSD(sdTemp)
	}

	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(pduSessionEstablishmentRequest)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(pduSessionEstablishmentRequest)

	m.GmmMessage.ULNASTransport = ulNasTransport

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetUlNasTransport_PduSessionModificationRequest(pduSessionId uint8, requestType uint8, dnnString string,
	sNssai *models.Snssai) []byte {

	pduSessionModificationRequest := GetPduSessionModificationRequest(pduSessionId)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(pduSessionId)
	ulNasTransport.RequestType = new(nasType.RequestType)
	ulNasTransport.RequestType.SetIei(nasMessage.ULNASTransportRequestTypeType)
	ulNasTransport.RequestType.SetRequestTypeValue(requestType)
	if dnnString != "" {
		ulNasTransport.DNN = new(nasType.DNN)
		ulNasTransport.DNN.SetIei(nasMessage.ULNASTransportDNNType)
		ulNasTransport.DNN.SetDNN(dnnString)
	}
	if sNssai != nil {
		var sdTemp [3]uint8
		sd, err := hex.DecodeString(sNssai.Sd)
		if err != nil {
			logger.NasMsgLog.Errorf("sNssai SD decode error: %+v", err)
		}
		copy(sdTemp[:], sd)
		ulNasTransport.SNSSAI = nasType.NewSNSSAI(nasMessage.ULNASTransportSNSSAIType)
		ulNasTransport.SNSSAI.SetLen(4)
		ulNasTransport.SNSSAI.SetSST(uint8(sNssai.Sst))
		ulNasTransport.SNSSAI.SetSD(sdTemp)
	}

	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(pduSessionModificationRequest)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(pduSessionModificationRequest)

	m.GmmMessage.ULNASTransport = ulNasTransport

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetPduSessionModificationRequest(pduSessionId uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionModificationRequest)

	pduSessionModificationRequest := nasMessage.NewPDUSessionModificationRequest(0)
	pduSessionModificationRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSSessionManagementMessage)
	pduSessionModificationRequest.SetMessageType(nas.MsgTypePDUSessionModificationRequest)
	pduSessionModificationRequest.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionModificationRequest.PTI.SetPTI(0x00)
	// pduSessionModificationRequest.RequestedQosFlowDescriptions = nasType.NewRequestedQosFlowDescriptions(nasMessage.
	// PDUSessionModificationRequestRequestedQosFlowDescriptionsType)
	// pduSessionModificationRequest.RequestedQosFlowDescriptions.SetLen(6)
	// pduSessionModificationRequest.RequestedQosFlowDescriptions.SetQoSFlowDescriptions([]uint8{0x09, 0x20, 0x41, 0x01,
	// 0x01, 0x09})

	m.GsmMessage.PDUSessionModificationRequest = pduSessionModificationRequest

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}
func GetPduSessionModificationComplete(pduSessionId uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionModificationComplete)

	pduSessionModificationComplete := nasMessage.NewPDUSessionModificationComplete(0)
	pduSessionModificationComplete.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSSessionManagementMessage)
	pduSessionModificationComplete.SetMessageType(nas.MsgTypePDUSessionModificationComplete)
	pduSessionModificationComplete.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionModificationComplete.PTI.SetPTI(0x00)

	m.GsmMessage.PDUSessionModificationComplete = pduSessionModificationComplete

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}
func GetPduSessionModificationCommandReject(pduSessionId uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionModificationCommandReject)

	pduSessionModificationCommandReject := nasMessage.NewPDUSessionModificationCommandReject(0)
	pduSessionModificationCommandReject.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSSessionManagementMessage)
	pduSessionModificationCommandReject.SetMessageType(nas.MsgTypePDUSessionModificationCommandReject)
	pduSessionModificationCommandReject.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionModificationCommandReject.PTI.SetPTI(0x00)

	m.GsmMessage.PDUSessionModificationCommandReject = pduSessionModificationCommandReject

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetPduSessionReleaseRequest(pduSessionId uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseRequest)

	pduSessionReleaseRequest := nasMessage.NewPDUSessionReleaseRequest(0)
	pduSessionReleaseRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSSessionManagementMessage)
	pduSessionReleaseRequest.SetMessageType(nas.MsgTypePDUSessionReleaseRequest)
	pduSessionReleaseRequest.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionReleaseRequest.PTI.SetPTI(0x00)

	m.GsmMessage.PDUSessionReleaseRequest = pduSessionReleaseRequest

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetPduSessionReleaseComplete(pduSessionId uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseComplete)

	pduSessionReleaseComplete := nasMessage.NewPDUSessionReleaseComplete(0)
	pduSessionReleaseComplete.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSSessionManagementMessage)
	pduSessionReleaseComplete.SetMessageType(nas.MsgTypePDUSessionReleaseComplete)
	pduSessionReleaseComplete.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionReleaseComplete.PTI.SetPTI(0x00)

	m.GsmMessage.PDUSessionReleaseComplete = pduSessionReleaseComplete

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetPduSessionReleaseReject(pduSessionId uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseReject)

	pduSessionReleaseReject := nasMessage.NewPDUSessionReleaseReject(0)
	pduSessionReleaseReject.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSSessionManagementMessage)
	pduSessionReleaseReject.SetMessageType(nas.MsgTypePDUSessionReleaseReject)
	pduSessionReleaseReject.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionReleaseReject.PTI.SetPTI(0x00)

	m.GsmMessage.PDUSessionReleaseReject = pduSessionReleaseReject

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetPduSessionAuthenticationComplete(pduSessionId uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionAuthenticationComplete)

	pduSessionAuthenticaitonComplete := nasMessage.NewPDUSessionAuthenticationComplete(0)
	pduSessionAuthenticaitonComplete.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSSessionManagementMessage)
	pduSessionAuthenticaitonComplete.SetMessageType(nas.MsgTypePDUSessionAuthenticationComplete)
	pduSessionAuthenticaitonComplete.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionAuthenticaitonComplete.PTI.SetPTI(0x00)
	pduSessionAuthenticaitonComplete.EAPMessage.SetLen(6)
	pduSessionAuthenticaitonComplete.EAPMessage.SetEAPMessage([]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55})

	m.GsmMessage.PDUSessionAuthenticationComplete = pduSessionAuthenticaitonComplete

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetUlNasTransport_PduSessionCommonData(pduSessionId uint8, types string) []byte {

	var payload []byte
	switch types {
	case PDUSesModiReq:
		payload = GetPduSessionModificationRequest(pduSessionId)
	case PDUSesModiCmp:
		payload = GetPduSessionModificationComplete(pduSessionId)
	case PDUSesModiCmdRej:
		payload = GetPduSessionModificationCommandReject(pduSessionId)
	case PDUSesRelReq:
		payload = GetPduSessionReleaseRequest(pduSessionId)
	case PDUSesRelCmp:
		payload = GetPduSessionReleaseComplete(pduSessionId)
	case PDUSesRelRej:
		payload = GetPduSessionReleaseReject(pduSessionId)
	case PDUSesAuthCmp:
		payload = GetPduSessionAuthenticationComplete(pduSessionId)
	}

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(pduSessionId)

	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(payload)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(payload)

	m.GmmMessage.ULNASTransport = ulNasTransport

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetIdentityResponse(mobileIdentity nasType.MobileIdentity) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeIdentityResponse)

	identityResponse := nasMessage.NewIdentityResponse(0)
	identityResponse.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	identityResponse.IdentityResponseMessageIdentity.SetMessageType(nas.MsgTypeIdentityResponse)
	identityResponse.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	identityResponse.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	identityResponse.MobileIdentity = mobileIdentity

	m.GmmMessage.IdentityResponse = identityResponse

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetNotificationResponse(pDUSessionStatus []uint8) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeNotificationResponse)

	notificationResponse := nasMessage.NewNotificationResponse(0)
	notificationResponse.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	notificationResponse.SetMessageType(nas.MsgTypeNotificationResponse)
	notificationResponse.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	notificationResponse.PDUSessionStatus = new(nasType.PDUSessionStatus)
	notificationResponse.PDUSessionStatus.SetIei(nasMessage.NotificationResponsePDUSessionStatusType)
	notificationResponse.PDUSessionStatus.Buffer = pDUSessionStatus

	m.GmmMessage.NotificationResponse = notificationResponse

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetConfigurationUpdateComplete() []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeConfigurationUpdateComplete)

	configurationUpdateComplete := nasMessage.NewConfigurationUpdateComplete(0)
	configurationUpdateComplete.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	configurationUpdateComplete.SetSecurityHeaderType(0x00)
	configurationUpdateComplete.SetSpareHalfOctet(0x00)
	configurationUpdateComplete.SetMessageType(nas.MsgTypeConfigurationUpdateComplete)

	m.GmmMessage.ConfigurationUpdateComplete = configurationUpdateComplete

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetServiceRequest(serviceType uint8) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeServiceRequest)

	serviceRequest := nasMessage.NewServiceRequest(0)
	serviceRequest.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	serviceRequest.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	serviceRequest.SetMessageType(nas.MsgTypeServiceRequest)
	serviceRequest.SetServiceTypeValue(serviceType)
	serviceRequest.SetNasKeySetIdentifiler(0x01)
	serviceRequest.SetAMFSetID(uint16(0xFE) << 2)
	serviceRequest.SetAMFPointer(0)
	serviceRequest.TMSI5GS.Octet[0] = 0xf0
	serviceRequest.SetTypeOfIdentity(4) // 5G-S-TMSI
	serviceRequest.SetTMSI5G([4]uint8{0, 0, 0, 1})
	serviceRequest.TMSI5GS.SetLen(7)
	switch serviceType {
	case nasMessage.ServiceTypeMobileTerminatedServices:
		serviceRequest.AllowedPDUSessionStatus = new(nasType.AllowedPDUSessionStatus)
		serviceRequest.AllowedPDUSessionStatus.SetIei(nasMessage.ServiceRequestAllowedPDUSessionStatusType)
		serviceRequest.AllowedPDUSessionStatus.SetLen(2)
		serviceRequest.AllowedPDUSessionStatus.Buffer = []uint8{0x00, 0x08}
	case nasMessage.ServiceTypeData:
		serviceRequest.UplinkDataStatus = new(nasType.UplinkDataStatus)
		serviceRequest.UplinkDataStatus.SetIei(nasMessage.ServiceRequestUplinkDataStatusType)
		serviceRequest.UplinkDataStatus.SetLen(2)
		serviceRequest.UplinkDataStatus.Buffer = []uint8{0x00, 0x04}
	case nasMessage.ServiceTypeSignalling:
	}

	m.GmmMessage.ServiceRequest = serviceRequest

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetAuthenticationResponse(authenticationResponseParam []uint8, eapMsg string) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationResponse)

	authenticationResponse := nasMessage.NewAuthenticationResponse(0)
	authenticationResponse.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	authenticationResponse.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationResponse.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationResponse.AuthenticationResponseMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationResponse)

	if len(authenticationResponseParam) > 0 {
		authenticationResponse.AuthenticationResponseParameter = nasType.NewAuthenticationResponseParameter(
			nasMessage.AuthenticationResponseAuthenticationResponseParameterType)
		authenticationResponse.AuthenticationResponseParameter.SetLen(uint8(len(authenticationResponseParam)))
		copy(authenticationResponse.AuthenticationResponseParameter.Octet[:], authenticationResponseParam[0:16])
	} else if eapMsg != "" {
		rawEapMsg, err := base64.StdEncoding.DecodeString(eapMsg)
		if err != nil {
			logger.NasMsgLog.Warnf("EAP decode error: %+v", err)
		}
		authenticationResponse.EAPMessage = nasType.NewEAPMessage(nasMessage.AuthenticationResponseEAPMessageType)
		authenticationResponse.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		authenticationResponse.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.AuthenticationResponse = authenticationResponse

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetAuthenticationFailure(cause5GMM uint8, authenticationFailureParam []uint8) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationFailure)

	authenticationFailure := nasMessage.NewAuthenticationFailure(0)
	authenticationFailure.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	authenticationFailure.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationFailure.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationFailure.AuthenticationFailureMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationFailure)
	authenticationFailure.Cause5GMM.SetCauseValue(cause5GMM)

	if cause5GMM == nasMessage.Cause5GMMSynchFailure {
		authenticationFailure.AuthenticationFailureParameter = nasType.NewAuthenticationFailureParameter(
			nasMessage.AuthenticationFailureAuthenticationFailureParameterType)
		authenticationFailure.AuthenticationFailureParameter.SetLen(uint8(len(authenticationFailureParam)))
		copy(authenticationFailure.AuthenticationFailureParameter.Octet[:], authenticationFailureParam)
	}

	m.GmmMessage.AuthenticationFailure = authenticationFailure

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetRegistrationComplete(sorTransparentContainer []uint8) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationComplete)

	registrationComplete := nasMessage.NewRegistrationComplete(0)
	registrationComplete.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	registrationComplete.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	registrationComplete.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	registrationComplete.RegistrationCompleteMessageIdentity.SetMessageType(nas.MsgTypeRegistrationComplete)

	if sorTransparentContainer != nil {
		registrationComplete.SORTransparentContainer = nasType.NewSORTransparentContainer(
			nasMessage.RegistrationCompleteSORTransparentContainerType)
		registrationComplete.SORTransparentContainer.SetLen(uint16(len(sorTransparentContainer)))
		registrationComplete.SORTransparentContainer.SetSORContent(sorTransparentContainer)
	}

	m.GmmMessage.RegistrationComplete = registrationComplete

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

// TS 24.501 8.2.26.
func GetSecurityModeComplete(nasMessageContainer []uint8) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeSecurityModeComplete)

	securityModeComplete := nasMessage.NewSecurityModeComplete(0)
	securityModeComplete.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	// TODO: modify security header type if need security protected
	securityModeComplete.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	securityModeComplete.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	securityModeComplete.SecurityModeCompleteMessageIdentity.SetMessageType(nas.MsgTypeSecurityModeComplete)

	securityModeComplete.IMEISV = nasType.NewIMEISV(nasMessage.SecurityModeCompleteIMEISVType)
	securityModeComplete.IMEISV.SetLen(9)
	securityModeComplete.SetOddEvenIdic(0)
	securityModeComplete.SetTypeOfIdentity(nasMessage.MobileIdentity5GSTypeImeisv)
	securityModeComplete.SetIdentityDigit1(1)
	securityModeComplete.SetIdentityDigitP_1(1)
	securityModeComplete.SetIdentityDigitP(1)

	if nasMessageContainer != nil {
		securityModeComplete.NASMessageContainer = nasType.NewNASMessageContainer(
			nasMessage.SecurityModeCompleteNASMessageContainerType)
		securityModeComplete.NASMessageContainer.SetLen(uint16(len(nasMessageContainer)))
		securityModeComplete.NASMessageContainer.SetNASMessageContainerContents(nasMessageContainer)
	}

	m.GmmMessage.SecurityModeComplete = securityModeComplete

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetSecurityModeReject(cause5GMM uint8) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeSecurityModeReject)

	securityModeReject := nasMessage.NewSecurityModeReject(0)
	securityModeReject.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	securityModeReject.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	securityModeReject.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	securityModeReject.SecurityModeRejectMessageIdentity.SetMessageType(nas.MsgTypeSecurityModeReject)

	securityModeReject.Cause5GMM.SetCauseValue(cause5GMM)

	m.GmmMessage.SecurityModeReject = securityModeReject

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetDeregistrationRequest(accessType uint8, switchOff uint8, ngKsi uint8,
	mobileIdentity5GS nasType.MobileIdentity5GS) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)

	deregistrationRequest := nasMessage.NewDeregistrationRequestUEOriginatingDeregistration(0)
	deregistrationRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	deregistrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deregistrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	deregistrationRequest.DeregistrationRequestMessageIdentity.SetMessageType(
		nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)

	deregistrationRequest.NgksiAndDeregistrationType.SetAccessType(accessType)
	deregistrationRequest.NgksiAndDeregistrationType.SetSwitchOff(switchOff)
	deregistrationRequest.NgksiAndDeregistrationType.SetReRegistrationRequired(0)
	deregistrationRequest.NgksiAndDeregistrationType.SetTSC(ngKsi)
	deregistrationRequest.NgksiAndDeregistrationType.SetNasKeySetIdentifiler(ngKsi)
	deregistrationRequest.MobileIdentity5GS.SetLen(mobileIdentity5GS.GetLen())
	deregistrationRequest.MobileIdentity5GS.SetMobileIdentity5GSContents(mobileIdentity5GS.GetMobileIdentity5GSContents())

	m.GmmMessage.DeregistrationRequestUEOriginatingDeregistration = deregistrationRequest

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetDeregistrationAccept() []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration)

	deregistrationAccept := nasMessage.NewDeregistrationAcceptUETerminatedDeregistration(0)
	deregistrationAccept.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	deregistrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deregistrationAccept.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	deregistrationAccept.DeregistrationAcceptMessageIdentity.SetMessageType(
		nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration)

	m.GmmMessage.DeregistrationAcceptUETerminatedDeregistration = deregistrationAccept

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetStatus5GMM(cause uint8) []byte {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeStatus5GMM)

	status5GMM := nasMessage.NewStatus5GMM(0)
	status5GMM.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	status5GMM.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	status5GMM.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	status5GMM.STATUSMessageIdentity5GMM.SetMessageType(nas.MsgTypeStatus5GMM)
	status5GMM.Cause5GMM.SetCauseValue(cause)

	m.GmmMessage.Status5GMM = status5GMM

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetStatus5GSM(pduSessionId uint8, cause uint8) []byte {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypeStatus5GSM)

	status5GSM := nasMessage.NewStatus5GSM(0)
	status5GSM.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	status5GSM.STATUSMessageIdentity5GSM.SetMessageType(nas.MsgTypeStatus5GSM)
	status5GSM.PDUSessionID.SetPDUSessionID(pduSessionId)
	status5GSM.PTI.SetPTI(0x00)
	status5GSM.Cause5GSM.SetCauseValue(cause)

	m.GsmMessage.Status5GSM = status5GSM

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetUlNasTransport_Status5GSM(pduSessionId uint8, cause uint8) []byte {

	payload := GetStatus5GSM(pduSessionId, cause)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(pduSessionId)

	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(payload)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(payload)

	m.GmmMessage.ULNASTransport = ulNasTransport

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetUlNasTransport_PduSessionReleaseRequest(pduSessionId uint8) []byte {

	pduSessionReleaseRequest := GetPduSessionReleaseRequest(pduSessionId)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(pduSessionId)

	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(pduSessionReleaseRequest)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(pduSessionReleaseRequest)

	m.GmmMessage.ULNASTransport = ulNasTransport

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}

func GetUlNasTransport_PduSessionReleaseComplete(pduSessionId uint8, requestType uint8, dnnString string,
	sNssai *models.Snssai) []byte {

	pduSessionReleaseRequest := GetPduSessionReleaseComplete(pduSessionId)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
		nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(pduSessionId)
	ulNasTransport.RequestType = new(nasType.RequestType)
	ulNasTransport.RequestType.SetIei(nasMessage.ULNASTransportRequestTypeType)
	ulNasTransport.RequestType.SetRequestTypeValue(requestType)
	if dnnString != "" {
		ulNasTransport.DNN = new(nasType.DNN)
		ulNasTransport.DNN.SetIei(nasMessage.ULNASTransportDNNType)
		ulNasTransport.DNN.SetDNN(dnnString)
	}
	if sNssai != nil {
		var sdTemp [3]uint8
		sd, err := hex.DecodeString(sNssai.Sd)
		if err != nil {
			logger.NasMsgLog.Warnf("sNssai SD decode error: %+v", err)
		}
		copy(sdTemp[:], sd)
		ulNasTransport.SNSSAI = nasType.NewSNSSAI(nasMessage.ULNASTransportSNSSAIType)
		ulNasTransport.SNSSAI.SetLen(4)
		ulNasTransport.SNSSAI.SetSST(uint8(sNssai.Sst))
		ulNasTransport.SNSSAI.SetSD(sdTemp)
	}

	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(pduSessionReleaseRequest)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(pduSessionReleaseRequest)

	m.GmmMessage.ULNASTransport = ulNasTransport

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data.Bytes()
}
