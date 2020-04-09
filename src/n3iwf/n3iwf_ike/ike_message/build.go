package ike_message

import (
	"encoding/binary"
	"net"
)

func BuildIKEHeader(
	initiatorSPI uint64,
	responsorSPI uint64,
	exchangeType uint8,
	flags uint8,
	messageID uint32) *IKEMessage {

	ikeMessage := new(IKEMessage)

	ikeMessage.InitiatorSPI = initiatorSPI
	ikeMessage.ResponderSPI = responsorSPI
	ikeMessage.Version = 0x20
	ikeMessage.ExchangeType = exchangeType
	ikeMessage.Flags = flags
	ikeMessage.MessageID = messageID

	return ikeMessage
}

func BuildNotification(protocolID uint8, notifyMessageType uint16, spi []byte, notificationData []byte) *Notification {
	notification := new(Notification)
	notification.ProtocolID = protocolID
	notification.NotifyMessageType = notifyMessageType
	notification.SPI = append(notification.SPI, spi...)
	notification.NotificationData = append(notification.NotificationData, notificationData...)
	return notification
}

func BuildCertificate(certificateEncode uint8, certificateData []byte) *Certificate {
	certificate := new(Certificate)
	certificate.CertificateEncoding = certificateEncode
	certificate.CertificateData = append(certificate.CertificateData, certificateData...)
	return certificate
}

func BuildEncrypted(nextPayload IKEType, encryptedData []byte) *Encrypted {
	encrypted := new(Encrypted)
	encrypted.NextPayload = uint8(nextPayload)
	encrypted.EncryptedData = append(encrypted.EncryptedData, encryptedData...)
	return encrypted
}

func BUildKeyExchange(diffiehellmanGroup uint16, keyExchangeData []byte) *KeyExchange {
	keyExchange := new(KeyExchange)
	keyExchange.DiffieHellmanGroup = diffiehellmanGroup
	keyExchange.KeyExchangeData = append(keyExchange.KeyExchangeData, keyExchangeData...)
	return keyExchange
}

func BuildIdentificationInitiator(idType uint8, idData []byte) *IdentificationInitiator {
	identification := new(IdentificationInitiator)
	identification.IDType = idType
	identification.IDData = append(identification.IDData, idData...)
	return identification
}

func BuildIdentificationResponder(idType uint8, idData []byte) *IdentificationResponder {
	identification := new(IdentificationResponder)
	identification.IDType = idType
	identification.IDData = append(identification.IDData, idData...)
	return identification
}

func BuildAuthentication(authenticationMethod uint8, authenticationData []byte) *Authentication {
	authentication := new(Authentication)
	authentication.AuthenticationMethod = authenticationMethod
	authentication.AuthenticationData = append(authentication.AuthenticationData, authenticationData...)
	return authentication
}

func BuildConfiguration(configurationType uint8, attributes []*IndividualConfigurationAttribute) *Configuration {
	configuration := new(Configuration)
	configuration.ConfigurationType = configurationType
	configuration.ConfigurationAttribute = append(configuration.ConfigurationAttribute, attributes...)
	return configuration
}

func BuildConfigurationAttribute(attributeType uint16, attributeValue []byte) *IndividualConfigurationAttribute {
	configurationAttribute := new(IndividualConfigurationAttribute)
	configurationAttribute.Type = attributeType
	configurationAttribute.Value = append(configurationAttribute.Value, attributeValue...)
	return configurationAttribute
}

func BuildNonce(nonceData []byte) *Nonce {
	nonce := new(Nonce)
	nonce.NonceData = append(nonce.NonceData, nonceData...)
	return nonce
}

func BuildTrafficSelectorInitiator(trafficSelectors []*IndividualTrafficSelector) *TrafficSelectorInitiator {
	trafficSelectorInitiator := new(TrafficSelectorInitiator)
	trafficSelectorInitiator.TrafficSelectors = append(trafficSelectorInitiator.TrafficSelectors, trafficSelectors...)
	return trafficSelectorInitiator
}

func BuildTrafficSelectorResponder(trafficSelectors []*IndividualTrafficSelector) *TrafficSelectorResponder {
	trafficSelectorResponder := new(TrafficSelectorResponder)
	trafficSelectorResponder.TrafficSelectors = append(trafficSelectorResponder.TrafficSelectors, trafficSelectors...)
	return trafficSelectorResponder
}

func BuildIndividualTrafficSelector(tsType uint8, ipProtocolID uint8, startPort uint16, endPort uint16, startAddr []byte, endAddr []byte) *IndividualTrafficSelector {
	trafficSelector := new(IndividualTrafficSelector)
	trafficSelector.TSType = tsType
	trafficSelector.IPProtocolID = ipProtocolID
	trafficSelector.StartPort = startPort
	trafficSelector.EndPort = endPort
	trafficSelector.StartAddress = append(trafficSelector.StartAddress, startAddr...)
	trafficSelector.EndAddress = append(trafficSelector.EndAddress, endAddr...)
	return trafficSelector
}

func BuildTransform(transformType uint8, transformID uint16, attributeType *uint16, attributeValue *uint16, variableLengthAttributeValue []byte) *Transform {
	transform := new(Transform)
	transform.TransformType = transformType
	transform.TransformID = transformID
	if attributeType != nil {
		transform.AttributePresent = true
		transform.AttributeType = *attributeType
		if attributeValue != nil {
			transform.AttributeFormat = AttributeFormatUseTV
			transform.AttributeValue = *attributeValue
		} else if len(variableLengthAttributeValue) != 0 {
			transform.AttributeFormat = AttributeFormatUseTLV
			transform.VariableLengthAttributeValue = append(transform.VariableLengthAttributeValue, variableLengthAttributeValue...)
		} else {
			return nil
		}
	} else {
		transform.AttributePresent = false
	}
	return transform
}

func BuildProposal(proposalNumber uint8, protocolID uint8, spi []byte) *Proposal {
	proposal := new(Proposal)
	proposal.ProposalNumber = proposalNumber
	proposal.ProtocolID = protocolID
	proposal.SPI = append(proposal.SPI, spi...)
	return proposal
}

func AppendTransformToProposal(proposal *Proposal, transform *Transform) bool {
	if proposal == nil {
		return false
	} else {
		switch transform.TransformType {
		case TypeEncryptionAlgorithm:
			proposal.EncryptionAlgorithm = append(proposal.EncryptionAlgorithm, transform)
		case TypePseudorandomFunction:
			proposal.PseudorandomFunction = append(proposal.PseudorandomFunction, transform)
		case TypeIntegrityAlgorithm:
			proposal.IntegrityAlgorithm = append(proposal.IntegrityAlgorithm, transform)
		case TypeDiffieHellmanGroup:
			proposal.DiffieHellmanGroup = append(proposal.DiffieHellmanGroup, transform)
		case TypeExtendedSequenceNumbers:
			proposal.ExtendedSequenceNumbers = append(proposal.ExtendedSequenceNumbers, transform)
		default:
			return false
		}
		return true
	}
}

func BuildSecurityAssociation(proposals []*Proposal) *SecurityAssociation {
	securityAssociation := new(SecurityAssociation)
	securityAssociation.Proposals = append(securityAssociation.Proposals, proposals...)
	return securityAssociation
}

func BuildEAP(code uint8, identifier uint8, eapTypeData EAPTypeFormat) *EAP {
	eap := new(EAP)
	eap.Code = code
	eap.Identifier = identifier
	eap.EAPTypeData = append(eap.EAPTypeData, eapTypeData)
	return eap
}

func BuildEAPSuccess(identifier uint8) *EAP {
	eap := new(EAP)
	eap.Code = EAPCodeSuccess
	eap.Identifier = identifier
	return eap
}

func BuildEAPfailure(identifier uint8) *EAP {
	eap := new(EAP)
	eap.Code = EAPCodeFailure
	eap.Identifier = identifier
	return eap
}

func BuildEAPExpanded(vendorID uint32, vendorType uint32, vendorData []byte) *EAPExpanded {
	eapExpanded := new(EAPExpanded)
	eapExpanded.VendorID = vendorID
	eapExpanded.VendorType = vendorType
	eapExpanded.VendorData = append(eapExpanded.VendorData, vendorData...)
	return eapExpanded
}

func BuildEAP5GStart(identifier uint8) *EAP {
	vendorData := []byte{EAP5GType5GStart, EAP5GSpareValue}
	eapTypeData := BuildEAPExpanded(VendorID3GPP, VendorTypeEAP5G, vendorData)
	return BuildEAP(EAPCodeRequest, identifier, eapTypeData)
}

func BuildEAP5GNAS(identifier uint8, nasPDU []byte) *EAP {
	if len(nasPDU) == 0 {
		ikeLog.Error("[IKE] BuildEAP5GNAS(): NASPDU is nil")
		return nil
	}

	header := make([]byte, 4)

	// Message ID
	header[0] = EAP5GType5GNAS

	// NASPDU length (2 octets)
	binary.BigEndian.PutUint16(header[2:4], uint16(len(nasPDU)))

	vendorData := append(header, nasPDU...)
	eapTypeData := BuildEAPExpanded(VendorID3GPP, VendorTypeEAP5G, vendorData)
	return BuildEAP(EAPCodeRequest, identifier, eapTypeData)
}

func BuildNotify5G_QOS_INFO(pduSessionID uint8, qfiList []uint8, isDefault bool) *Notification {
	notifyData := make([]byte, 1)

	// Append PDU session ID
	notifyData = append(notifyData, pduSessionID)

	// Append QFI list length
	notifyData = append(notifyData, uint8(len(qfiList)))

	// Append QFI list
	notifyData = append(notifyData, qfiList...)

	// Append default and differentiated service flags
	var defaultAndDifferentiatedServiceFlags uint8
	if isDefault {
		defaultAndDifferentiatedServiceFlags |= NotifyType5G_QOS_INFOBitDCSICheck
	}
	notifyData = append(notifyData, defaultAndDifferentiatedServiceFlags)

	// Assign length
	notifyData[0] = uint8(len(notifyData))

	return BuildNotification(TypeNone, Vendor3GPPNotifyType5G_QOS_INFO, nil, notifyData)
}

func BuildNotifyNAS_IP4_ADDRESS(nasIPAddr string) *Notification {
	if nasIPAddr == "" {
		return nil
	} else {
		ipAddrByte := net.ParseIP(nasIPAddr).To4()
		return BuildNotification(TypeNone, Vendor3GPPNotifyTypeNAS_IP4_ADDRESS, nil, ipAddrByte)
	}
}

func BuildNotifyUP_IP4_ADDRESS(upIPAddr string) *Notification {
	if upIPAddr == "" {
		return nil
	} else {
		ipAddrByte := net.ParseIP(upIPAddr).To4()
		return BuildNotification(TypeNone, Vendor3GPPNotifyTypeUP_IP4_ADDRESS, nil, ipAddrByte)
	}
}

func BuildNotifyNAS_TCP_PORT(port uint16) *Notification {
	if port == 0 {
		return nil
	} else {
		portData := make([]byte, 2)
		binary.BigEndian.PutUint16(portData, port)
		return BuildNotification(TypeNone, Vendor3GPPNotifyTypeNAS_TCP_PORT, nil, portData)
	}
}
