package ike_message

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"free5gc/src/n3iwf/logger"
)

// Log
var ikeLog *logrus.Entry

func init() {
	ikeLog = logger.IKELog
}

type IKEMessage struct {
	InitiatorSPI uint64
	ResponderSPI uint64
	Version      uint8
	ExchangeType uint8
	Flags        uint8
	MessageID    uint32
	IKEPayload   []IKEPayloadType
}

func Encode(ikeMessage *IKEMessage) ([]byte, error) {
	ikeLog.Info("Encoding IKE message")

	if ikeMessage != nil {
		ikeLog.Info("[IKE] Start encoding IKE message")

		ikeMessageData := make([]byte, 28)

		binary.BigEndian.PutUint64(ikeMessageData[0:8], ikeMessage.InitiatorSPI)
		binary.BigEndian.PutUint64(ikeMessageData[8:16], ikeMessage.ResponderSPI)
		ikeMessageData[17] = ikeMessage.Version
		ikeMessageData[18] = ikeMessage.ExchangeType
		ikeMessageData[19] = ikeMessage.Flags
		binary.BigEndian.PutUint32(ikeMessageData[20:24], ikeMessage.MessageID)

		if len(ikeMessage.IKEPayload) > 0 {
			ikeMessageData[16] = byte(ikeMessage.IKEPayload[0].Type())
		} else {
			ikeMessageData[16] = NoNext
		}

		ikeMessagePayloadData, err := EncodePayload(ikeMessage.IKEPayload)
		if err != nil {
			return nil, fmt.Errorf("[IKE] Encode(): EncodePayload failed: %+v", err)
		}

		ikeMessageData = append(ikeMessageData, ikeMessagePayloadData...)
		binary.BigEndian.PutUint32(ikeMessageData[24:28], uint32(len(ikeMessageData)))

		ikeLog.Tracef("Encoded %d bytes", len(ikeMessageData))
		ikeLog.Tracef("IKE message data:\n%s", hex.Dump(ikeMessageData))

		return ikeMessageData, nil
	} else {
		return nil, errors.New("[IKE] Encode(): Input IKE message is nil")
	}
}

func EncodePayload(ikePayload []IKEPayloadType) ([]byte, error) {
	ikeLog.Info("Encoding IKE payloads")

	ikeMessagePayloadData := make([]byte, 0)

	for index, payload := range ikePayload {
		payloadData := make([]byte, 4)
		if (index + 1) < len(ikePayload) {
			payloadData[0] = uint8(ikePayload[index+1].Type())
		} else {
			if payload.Type() == TypeSK {
				payloadData[0] = payload.(*Encrypted).NextPayload
			} else {
				payloadData[0] = NoNext
			}
		}

		data, err := payload.marshal()
		if err != nil {
			return nil, fmt.Errorf("EncodePayload(): Failed to marshal payload: %+v", err)
		}

		payloadData = append(payloadData, data...)
		binary.BigEndian.PutUint16(payloadData[2:4], uint16(len(payloadData)))

		ikeMessagePayloadData = append(ikeMessagePayloadData, payloadData...)
	}

	return ikeMessagePayloadData, nil
}

func Decode(rawData []byte) (*IKEMessage, error) {
	// IKE message packet format this implementation referenced is
	// defined in RFC 7296, Section 3.1
	ikeLog.Info("Decoding IKE message")
	ikeLog.Tracef("Received IKE message:\n%s", hex.Dump(rawData))

	// bounds checking
	if len(rawData) < 28 {
		return nil, errors.New("[IKE] Decode(): Received broken IKE header")
	}
	ikeMessageLength := binary.BigEndian.Uint32(rawData[24:28])
	if ikeMessageLength < 28 {
		return nil, fmt.Errorf("[IKE] Decode(): Illegal IKE message length %d < header length 20", ikeMessageLength)
	}
	// len() return int, which is 64 bit on 64-bit host and 32 bit
	// on 32-bit host, so this implementation may potentially cause
	// problem on 32-bit machine
	if len(rawData) != int(ikeMessageLength) {
		return nil, errors.New("[IKE] Decode(): The length of received message not matchs the length specified in header")
	}

	nextPayload := rawData[16]

	ikeMessage := new(IKEMessage)

	ikeMessage.InitiatorSPI = binary.BigEndian.Uint64(rawData[:8])
	ikeMessage.ResponderSPI = binary.BigEndian.Uint64(rawData[8:16])
	ikeMessage.Version = rawData[17]
	ikeMessage.ExchangeType = rawData[18]
	ikeMessage.Flags = rawData[19]
	ikeMessage.MessageID = binary.BigEndian.Uint32(rawData[20:24])

	ikePayload, err := DecodePayload(nextPayload, rawData[28:])
	if err != nil {
		return nil, fmt.Errorf("[IKE] Decode(): DecodePayload failed: %+v", err)
	}
	ikeMessage.IKEPayload = append(ikeMessage.IKEPayload, ikePayload...)

	return ikeMessage, nil
}

func DecodePayload(nextPayload uint8, rawData []byte) ([]IKEPayloadType, error) {
	ikeLog.Info("Decoding IKE payloads")

	var ikePayload []IKEPayloadType

	for len(rawData) > 0 {
		// bounds checking
		ikeLog.Trace("[IKE] DecodePayload(): Decode 1 payload")
		if len(rawData) < 4 {
			return nil, errors.New("DecodePayload(): No sufficient bytes to decode next payload")
		}
		payloadLength := binary.BigEndian.Uint16(rawData[2:4])
		if payloadLength < 4 {
			return nil, fmt.Errorf("DecodePayload(): Illegal payload length %d < header length 4", payloadLength)
		}
		if len(rawData) < int(payloadLength) {
			return nil, errors.New("DecodePayload(): The length of received message not matchs the length specified in header")
		}

		criticalBit := (rawData[1] & 0x80) >> 7

		var payload IKEPayloadType

		switch nextPayload {
		case TypeSA:
			payload = new(SecurityAssociation)
		case TypeKE:
			payload = new(KeyExchange)
		case TypeIDi:
			payload = new(IdentificationInitiator)
		case TypeIDr:
			payload = new(IdentificationResponder)
		case TypeCERT:
			payload = new(Certificate)
		case TypeCERTreq:
			payload = new(CertificateRequest)
		case TypeAUTH:
			payload = new(Authentication)
		case TypeNiNr:
			payload = new(Nonce)
		case TypeN:
			payload = new(Notification)
		case TypeD:
			payload = new(Delete)
		case TypeV:
			payload = new(VendorID)
		case TypeTSi:
			payload = new(TrafficSelectorInitiator)
		case TypeTSr:
			payload = new(TrafficSelectorResponder)
		case TypeSK:
			encryptedPayload := new(Encrypted)
			encryptedPayload.NextPayload = rawData[0]
			payload = encryptedPayload
		case TypeCP:
			payload = new(Configuration)
		case TypeEAP:
			payload = new(EAP)
		default:
			if criticalBit == 0 {
				// Skip this payload
				nextPayload = rawData[0]
				rawData = rawData[payloadLength:]
				continue
			} else {
				// TODO: Reject this IKE message
				return nil, fmt.Errorf("Unknown payload type: %d", nextPayload)
			}
		}

		if err := payload.unmarshal(rawData[4:payloadLength]); err != nil {
			return nil, fmt.Errorf("DecodePayload(): Unmarshal payload failed: %+v", err)
		}

		ikePayload = append(ikePayload, payload)

		nextPayload = rawData[0]
		rawData = rawData[payloadLength:]
	}

	return ikePayload, nil
}

type IKEPayloadType interface {
	// Type specifies the IKE payload types
	Type() IKEType

	// Called by Encode() or Decode()
	marshal() ([]byte, error)
	unmarshal(rawData []byte) error
}

// Definition of Security Association

var _ IKEPayloadType = &SecurityAssociation{}

type SecurityAssociation struct {
	Proposals []*Proposal
}

type Proposal struct {
	ProposalNumber          uint8
	ProtocolID              uint8
	SPI                     []byte
	EncryptionAlgorithm     []*Transform
	PseudorandomFunction    []*Transform
	IntegrityAlgorithm      []*Transform
	DiffieHellmanGroup      []*Transform
	ExtendedSequenceNumbers []*Transform
}

type Transform struct {
	TransformType                uint8
	TransformID                  uint16
	AttributePresent             bool
	AttributeFormat              uint8
	AttributeType                uint16
	AttributeValue               uint16
	VariableLengthAttributeValue []byte
}

func (securityAssociation *SecurityAssociation) Type() IKEType { return TypeSA }

func (securityAssociation *SecurityAssociation) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][SecurityAssociation] marshal(): Start marshalling")

	securityAssociationData := make([]byte, 0)

	for proposalIndex, proposal := range securityAssociation.Proposals {
		proposalData := make([]byte, 8)

		if (proposalIndex + 1) < len(securityAssociation.Proposals) {
			proposalData[0] = 2
		} else {
			proposalData[0] = 0
		}

		proposalData[4] = proposal.ProposalNumber
		proposalData[5] = proposal.ProtocolID

		proposalData[6] = uint8(len(proposal.SPI))
		if len(proposal.SPI) > 0 {
			proposalData = append(proposalData, proposal.SPI...)
		}

		// combine all transforms
		var transformList []*Transform
		transformList = append(transformList, proposal.EncryptionAlgorithm...)
		transformList = append(transformList, proposal.PseudorandomFunction...)
		transformList = append(transformList, proposal.IntegrityAlgorithm...)
		transformList = append(transformList, proposal.DiffieHellmanGroup...)
		transformList = append(transformList, proposal.ExtendedSequenceNumbers...)

		if len(transformList) == 0 {
			return nil, errors.New("One proposal has no any transform")
		}
		proposalData[7] = uint8(len(transformList))

		proposalTransformData := make([]byte, 0)

		for transformIndex, transform := range transformList {
			transformData := make([]byte, 8)

			if (transformIndex + 1) < len(transformList) {
				transformData[0] = 3
			} else {
				transformData[0] = 0
			}

			transformData[4] = transform.TransformType
			binary.BigEndian.PutUint16(transformData[6:8], transform.TransformID)

			if transform.AttributePresent {
				attributeData := make([]byte, 4)

				if transform.AttributeFormat == 0 {
					// TLV
					if len(transform.VariableLengthAttributeValue) == 0 {
						return nil, errors.New("Attribute of one transform not specified")
					}
					attributeFormatAndType := ((uint16(transform.AttributeFormat) & 0x1) << 15) | transform.AttributeType
					binary.BigEndian.PutUint16(attributeData[0:2], attributeFormatAndType)
					binary.BigEndian.PutUint16(attributeData[2:4], uint16(len(transform.VariableLengthAttributeValue)))
					attributeData = append(attributeData, transform.VariableLengthAttributeValue...)
				} else {
					// TV
					attributeFormatAndType := ((uint16(transform.AttributeFormat) & 0x1) << 15) | transform.AttributeType
					binary.BigEndian.PutUint16(attributeData[0:2], attributeFormatAndType)
					binary.BigEndian.PutUint16(attributeData[2:4], transform.AttributeValue)
				}

				transformData = append(transformData, attributeData...)
			}

			binary.BigEndian.PutUint16(transformData[2:4], uint16(len(transformData)))

			proposalTransformData = append(proposalTransformData, transformData...)
		}

		proposalData = append(proposalData, proposalTransformData...)
		binary.BigEndian.PutUint16(proposalData[2:4], uint16(len(proposalData)))

		securityAssociationData = append(securityAssociationData, proposalData...)
	}

	return securityAssociationData, nil
}

func (securityAssociation *SecurityAssociation) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][SecurityAssociation] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][SecurityAssociation] unmarshal(): Payload length %d bytes", len(rawData))

	for len(rawData) > 0 {
		ikeLog.Trace("[IKE][SecurityAssociation] unmarshal(): Unmarshal 1 proposal")
		// bounds checking
		if len(rawData) < 8 {
			return errors.New("Proposal: No sufficient bytes to decode next proposal")
		}
		proposalLength := binary.BigEndian.Uint16(rawData[2:4])
		if proposalLength < 8 {
			return errors.New("Proposal: Illegal payload length %d < header length 8")
		}
		if len(rawData) < int(proposalLength) {
			return errors.New("Proposal: The length of received message not matchs the length specified in header")
		}

		// Log whether this proposal is the last
		if rawData[0] == 0 {
			ikeLog.Trace("[IKE][SecurityAssociation] This proposal is the last")
		}
		// Log the number of transform in the proposal
		ikeLog.Tracef("[IKE][SecurityAssociation] This proposal contained %d transform", rawData[7])

		proposal := new(Proposal)
		var transformData []byte

		proposal.ProposalNumber = rawData[4]
		proposal.ProtocolID = rawData[5]

		spiSize := rawData[6]
		if spiSize > 0 {
			// bounds checking
			if len(rawData) < int(8+spiSize) {
				return errors.New("Proposal: No sufficient bytes for unmarshalling SPI of proposal")
			}
			proposal.SPI = append(proposal.SPI, rawData[8:8+spiSize]...)
		}

		transformData = rawData[8+spiSize : proposalLength]

		for len(transformData) > 0 {
			// bounds checking
			ikeLog.Trace("[IKE][SecurityAssociation] unmarshal(): Unmarshal 1 transform")
			if len(transformData) < 8 {
				return errors.New("Transform: No sufficient bytes to decode next transform")
			}
			transformLength := binary.BigEndian.Uint16(transformData[2:4])
			if transformLength < 8 {
				return errors.New("Transform: Illegal payload length %d < header length 8")
			}
			if len(transformData) < int(transformLength) {
				return errors.New("Transform: The length of received message not matchs the length specified in header")
			}

			// Log whether this transform is the last
			if transformData[0] == 0 {
				ikeLog.Trace("[IKE][SecurityAssociation] This transform is the last")
			}

			transform := new(Transform)

			transform.TransformType = transformData[4]
			transform.TransformID = binary.BigEndian.Uint16(transformData[6:8])
			if transformLength > 8 {
				transform.AttributePresent = true
				transform.AttributeFormat = ((transformData[8] & 0x80) >> 7)
				transform.AttributeType = binary.BigEndian.Uint16(transformData[8:10]) & 0x7f

				if transform.AttributeFormat == 0 {
					attributeLength := binary.BigEndian.Uint16(transformData[10:12])
					// bounds checking
					if (12 + attributeLength) != transformLength {
						return fmt.Errorf("Illegal attribute length %d not satisfies the transform length %d", attributeLength, transformLength)
					}
					copy(transform.VariableLengthAttributeValue, transformData[12:12+attributeLength])
				} else {
					transform.AttributeValue = binary.BigEndian.Uint16(transformData[10:12])
				}

			}

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
			}

			transformData = transformData[transformLength:]
		}

		securityAssociation.Proposals = append(securityAssociation.Proposals, proposal)

		rawData = rawData[proposalLength:]
	}

	return nil
}

// Definition of Key Exchange

var _ IKEPayloadType = &KeyExchange{}

type KeyExchange struct {
	DiffieHellmanGroup uint16
	KeyExchangeData    []byte
}

func (keyExchange *KeyExchange) Type() IKEType { return TypeKE }

func (keyExchange *KeyExchange) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][KeyExchange] marshal(): Start marshalling")

	keyExchangeData := make([]byte, 4)

	binary.BigEndian.PutUint16(keyExchangeData[0:2], keyExchange.DiffieHellmanGroup)
	keyExchangeData = append(keyExchangeData, keyExchange.KeyExchangeData...)

	return keyExchangeData, nil
}

func (keyExchange *KeyExchange) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][KeyExchange] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][KeyExchange] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][KeyExchange] unmarshal(): Unmarshal 1 key exchange data")
		// bounds checking
		if len(rawData) <= 4 {
			return errors.New("KeyExchange: No sufficient bytes to decode next key exchange data")
		}

		keyExchange.DiffieHellmanGroup = binary.BigEndian.Uint16(rawData[0:2])
		keyExchange.KeyExchangeData = append(keyExchange.KeyExchangeData, rawData[4:]...)
	}

	return nil
}

// Definition of Identification - Initiator

var _ IKEPayloadType = &IdentificationInitiator{}

type IdentificationInitiator struct {
	IDType uint8
	IDData []byte
}

func (identification *IdentificationInitiator) Type() IKEType { return TypeIDi }

func (identification *IdentificationInitiator) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Identification] marshal(): Start marshalling")

	identificationData := make([]byte, 4)

	identificationData[0] = identification.IDType
	identificationData = append(identificationData, identification.IDData...)

	return identificationData, nil
}

func (identification *IdentificationInitiator) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Identification] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Identification] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][Identification] unmarshal(): Unmarshal 1 identification")
		// bounds checking
		if len(rawData) <= 4 {
			return errors.New("Identification: No sufficient bytes to decode next identification")
		}

		identification.IDType = rawData[0]
		identification.IDData = append(identification.IDData, rawData[4:]...)
	}

	return nil
}

// Definition of Identification - Responder

var _ IKEPayloadType = &IdentificationResponder{}

type IdentificationResponder struct {
	IDType uint8
	IDData []byte
}

func (identification *IdentificationResponder) Type() IKEType { return TypeIDr }

func (identification *IdentificationResponder) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Identification] marshal(): Start marshalling")

	identificationData := make([]byte, 4)

	identificationData[0] = identification.IDType
	identificationData = append(identificationData, identification.IDData...)

	return identificationData, nil
}

func (identification *IdentificationResponder) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Identification] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Identification] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][Identification] unmarshal(): Unmarshal 1 identification")
		// bounds checking
		if len(rawData) <= 4 {
			return errors.New("Identification: No sufficient bytes to decode next identification")
		}

		identification.IDType = rawData[0]
		identification.IDData = append(identification.IDData, rawData[4:]...)
	}

	return nil
}

// Definition of Certificate

var _ IKEPayloadType = &Certificate{}

type Certificate struct {
	CertificateEncoding uint8
	CertificateData     []byte
}

func (certificate *Certificate) Type() IKEType { return TypeCERT }

func (certificate *Certificate) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Certificate] marshal(): Start marshalling")

	certificateData := make([]byte, 1)

	certificateData[0] = certificate.CertificateEncoding
	certificateData = append(certificateData, certificate.CertificateData...)

	return certificateData, nil
}

func (certificate *Certificate) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Certificate] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Certificate] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][Certificate] unmarshal(): Unmarshal 1 certificate")
		// bounds checking
		if len(rawData) <= 1 {
			return errors.New("Certificate: No sufficient bytes to decode next certificate")
		}

		certificate.CertificateEncoding = rawData[0]
		certificate.CertificateData = append(certificate.CertificateData, rawData[1:]...)
	}

	return nil
}

// Definition of Certificate Request

var _ IKEPayloadType = &CertificateRequest{}

type CertificateRequest struct {
	CertificateEncoding    uint8
	CertificationAuthority []byte
}

func (certificateRequest *CertificateRequest) Type() IKEType { return TypeCERTreq }

func (certificateRequest *CertificateRequest) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][CertificateRequest] marshal(): Start marshalling")

	certificateRequestData := make([]byte, 1)

	certificateRequestData[0] = certificateRequest.CertificateEncoding
	certificateRequestData = append(certificateRequestData, certificateRequest.CertificationAuthority...)

	return certificateRequestData, nil
}

func (certificateRequest *CertificateRequest) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][CertificateRequest] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][CertificateRequest] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][CertificateRequest] unmarshal(): Unmarshal 1 certificate request")
		// bounds checking
		if len(rawData) <= 1 {
			return errors.New("CertificateRequest: No sufficient bytes to decode next certificate request")
		}

		certificateRequest.CertificateEncoding = rawData[0]
		certificateRequest.CertificationAuthority = append(certificateRequest.CertificationAuthority, rawData[1:]...)
	}

	return nil
}

// Definition of Authentication

var _ IKEPayloadType = &Authentication{}

type Authentication struct {
	AuthenticationMethod uint8
	AuthenticationData   []byte
}

func (authentication *Authentication) Type() IKEType { return TypeAUTH }

func (authentication *Authentication) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Authentication] marshal(): Start marshalling")

	authenticationData := make([]byte, 4)

	authenticationData[0] = authentication.AuthenticationMethod
	authenticationData = append(authenticationData, authentication.AuthenticationData...)

	return authenticationData, nil
}

func (authentication *Authentication) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Authentication] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Authentication] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][Authentication] unmarshal(): Unmarshal 1 authentication")
		// bounds checking
		if len(rawData) <= 4 {
			return errors.New("Authentication: No sufficient bytes to decode next authentication")
		}

		authentication.AuthenticationMethod = rawData[0]
		authentication.AuthenticationData = append(authentication.AuthenticationData, rawData[4:]...)
	}

	return nil
}

// Definition of Nonce

var _ IKEPayloadType = &Nonce{}

type Nonce struct {
	NonceData []byte
}

func (nonce *Nonce) Type() IKEType { return TypeNiNr }

func (nonce *Nonce) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Nonce] marshal(): Start marshalling")

	nonceData := make([]byte, 0)
	nonceData = append(nonceData, nonce.NonceData...)

	return nonceData, nil
}

func (nonce *Nonce) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Nonce] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Nonce] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][Nonce] unmarshal(): Unmarshal 1 nonce")
		nonce.NonceData = append(nonce.NonceData, rawData...)
	}

	return nil
}

// Definition of Notification

var _ IKEPayloadType = &Notification{}

type Notification struct {
	ProtocolID        uint8
	NotifyMessageType uint16
	SPI               []byte
	NotificationData  []byte
}

func (notification *Notification) Type() IKEType { return TypeN }

func (notification *Notification) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Notification] marshal(): Start marshalling")

	notificationData := make([]byte, 4)

	notificationData[0] = notification.ProtocolID
	notificationData[1] = uint8(len(notification.SPI))
	binary.BigEndian.PutUint16(notificationData[2:4], notification.NotifyMessageType)

	notificationData = append(notificationData, notification.SPI...)
	notificationData = append(notificationData, notification.NotificationData...)

	return notificationData, nil
}

func (notification *Notification) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Notification] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Notification] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][Notification] unmarshal(): Unmarshal 1 notification")
		// bounds checking
		if len(rawData) < 4 {
			return errors.New("Notification: No sufficient bytes to decode next notification")
		}
		spiSize := rawData[1]
		if len(rawData) < int(4+spiSize) {
			return errors.New("Notification: No sufficient bytes to get SPI according to the length specified in header")
		}

		notification.ProtocolID = rawData[0]
		notification.NotifyMessageType = binary.BigEndian.Uint16(rawData[2:4])

		notification.SPI = append(notification.SPI, rawData[4:4+spiSize]...)
		notification.NotificationData = append(notification.NotificationData, rawData[4+spiSize:]...)
	}

	return nil
}

// Definition of Delete

var _ IKEPayloadType = &Delete{}

type Delete struct {
	ProtocolID  uint8
	SPISize     uint8
	NumberOfSPI uint16
	SPIs        []byte
}

func (delete *Delete) Type() IKEType { return TypeD }

func (delete *Delete) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Delete] marshal(): Start marshalling")

	if len(delete.SPIs) != (int(delete.SPISize) * int(delete.NumberOfSPI)) {
		return nil, fmt.Errorf("Total bytes of all SPIs not correct")
	}

	deleteData := make([]byte, 4)

	deleteData[0] = delete.ProtocolID
	deleteData[1] = delete.SPISize
	binary.BigEndian.PutUint16(deleteData[2:4], delete.NumberOfSPI)

	deleteData = append(deleteData, delete.SPIs...)

	return deleteData, nil
}

func (delete *Delete) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Delete] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Delete] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][Delete] unmarshal(): Unmarshal 1 delete")
		// bounds checking
		if len(rawData) <= 4 {
			return errors.New("Delete: No sufficient bytes to decode next delete")
		}
		spiSize := rawData[1]
		numberOfSPI := binary.BigEndian.Uint16(rawData[2:4])
		if len(rawData) < (4 + (int(spiSize) * int(numberOfSPI))) {
			return errors.New("Delete: No Sufficient bytes to get SPIs according to the length specified in header")
		}

		delete.ProtocolID = rawData[0]
		delete.SPISize = spiSize
		delete.NumberOfSPI = numberOfSPI

		delete.SPIs = append(delete.SPIs, rawData[4:]...)
	}

	return nil
}

// Definition of Vendor ID

var _ IKEPayloadType = &VendorID{}

type VendorID struct {
	VendorIDData []byte
}

func (vendorID *VendorID) Type() IKEType { return TypeV }

func (vendorID *VendorID) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][VendorID] marshal(): Start marshalling")
	return vendorID.VendorIDData, nil
}

func (vendorID *VendorID) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][VendorID] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][VendorID] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][VendorID] unmarshal(): Unmarshal 1 vendor ID")
		vendorID.VendorIDData = append(vendorID.VendorIDData, rawData...)
	}

	return nil
}

// Definition of Traffic Selector - Initiator

var _ IKEPayloadType = &TrafficSelectorInitiator{}

type TrafficSelectorInitiator struct {
	TrafficSelectors []*IndividualTrafficSelector
}

type IndividualTrafficSelector struct {
	TSType       uint8
	IPProtocolID uint8
	StartPort    uint16
	EndPort      uint16
	StartAddress []byte
	EndAddress   []byte
}

func (trafficSelector *TrafficSelectorInitiator) Type() IKEType { return TypeTSi }

func (trafficSelector *TrafficSelectorInitiator) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][TrafficSelector] marshal(): Start marshalling")

	if len(trafficSelector.TrafficSelectors) > 0 {
		trafficSelectorData := make([]byte, 4)
		trafficSelectorData[0] = uint8(len(trafficSelector.TrafficSelectors))

		for _, individualTrafficSelector := range trafficSelector.TrafficSelectors {
			if individualTrafficSelector.TSType == TS_IPV4_ADDR_RANGE {
				// Address length checking
				if len(individualTrafficSelector.StartAddress) != 4 {
					ikeLog.Errorf("Address length %d", len(individualTrafficSelector.StartAddress))
					return nil, errors.New("TrafficSelector: Start IPv4 address length is not correct")
				}
				if len(individualTrafficSelector.EndAddress) != 4 {
					return nil, errors.New("TrafficSelector: End IPv4 address length is not correct")
				}

				individualTrafficSelectorData := make([]byte, 8)

				individualTrafficSelectorData[0] = individualTrafficSelector.TSType
				individualTrafficSelectorData[1] = individualTrafficSelector.IPProtocolID
				binary.BigEndian.PutUint16(individualTrafficSelectorData[4:6], individualTrafficSelector.StartPort)
				binary.BigEndian.PutUint16(individualTrafficSelectorData[6:8], individualTrafficSelector.EndPort)

				individualTrafficSelectorData = append(individualTrafficSelectorData, individualTrafficSelector.StartAddress...)
				individualTrafficSelectorData = append(individualTrafficSelectorData, individualTrafficSelector.EndAddress...)

				binary.BigEndian.PutUint16(individualTrafficSelectorData[2:4], uint16(len(individualTrafficSelectorData)))

				trafficSelectorData = append(trafficSelectorData, individualTrafficSelectorData...)
			} else if individualTrafficSelector.TSType == TS_IPV6_ADDR_RANGE {
				// Address length checking
				if len(individualTrafficSelector.StartAddress) != 16 {
					return nil, errors.New("TrafficSelector: Start IPv6 address length is not correct")
				}
				if len(individualTrafficSelector.EndAddress) != 16 {
					return nil, errors.New("TrafficSelector: End IPv6 address length is not correct")
				}

				individualTrafficSelectorData := make([]byte, 8)

				individualTrafficSelectorData[0] = individualTrafficSelector.TSType
				individualTrafficSelectorData[1] = individualTrafficSelector.IPProtocolID
				binary.BigEndian.PutUint16(individualTrafficSelectorData[4:6], individualTrafficSelector.StartPort)
				binary.BigEndian.PutUint16(individualTrafficSelectorData[6:8], individualTrafficSelector.EndPort)

				individualTrafficSelectorData = append(individualTrafficSelectorData, individualTrafficSelector.StartAddress...)
				individualTrafficSelectorData = append(individualTrafficSelectorData, individualTrafficSelector.EndAddress...)

				binary.BigEndian.PutUint16(individualTrafficSelectorData[2:4], uint16(len(individualTrafficSelectorData)))

				trafficSelectorData = append(trafficSelectorData, individualTrafficSelectorData...)
			} else {
				return nil, errors.New("TrafficSelector: Unsupported traffic selector type")
			}
		}

		return trafficSelectorData, nil
	} else {
		return nil, errors.New("TrafficSelector: Contains no traffic selector for marshalling message")
	}
}

func (trafficSelector *TrafficSelectorInitiator) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][TrafficSelector] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][TrafficSelector] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][TrafficSelector] unmarshal(): Unmarshal 1 traffic selector")
		// bounds checking
		if len(rawData) < 4 {
			return errors.New("TrafficSelector: No sufficient bytes to get number of traffic selector in header")
		}

		numberOfSPI := rawData[0]

		rawData = rawData[4:]

		for ; numberOfSPI > 0; numberOfSPI-- {
			// bounds checking
			if len(rawData) < 4 {
				return errors.New("TrafficSelector: No sufficient bytes to decode next individual traffic selector length in header")
			}
			trafficSelectorType := rawData[0]
			if trafficSelectorType == TS_IPV4_ADDR_RANGE {
				selectorLength := binary.BigEndian.Uint16(rawData[2:4])
				if selectorLength != 16 {
					return errors.New("TrafficSelector: A TS_IPV4_ADDR_RANGE type traffic selector should has length 16 bytes")
				}
				if len(rawData) < int(selectorLength) {
					return errors.New("TrafficSelector: No sufficient bytes to decode next individual traffic selector")
				}

				individualTrafficSelector := &IndividualTrafficSelector{}

				individualTrafficSelector.TSType = rawData[0]
				individualTrafficSelector.IPProtocolID = rawData[1]
				individualTrafficSelector.StartPort = binary.BigEndian.Uint16(rawData[4:6])
				individualTrafficSelector.EndPort = binary.BigEndian.Uint16(rawData[6:8])

				individualTrafficSelector.StartAddress = append(individualTrafficSelector.StartAddress, rawData[8:12]...)
				individualTrafficSelector.EndAddress = append(individualTrafficSelector.EndAddress, rawData[12:16]...)

				trafficSelector.TrafficSelectors = append(trafficSelector.TrafficSelectors, individualTrafficSelector)

				rawData = rawData[16:]
			} else if trafficSelectorType == TS_IPV6_ADDR_RANGE {
				selectorLength := binary.BigEndian.Uint16(rawData[2:4])
				if selectorLength != 40 {
					return errors.New("TrafficSelector: A TS_IPV6_ADDR_RANGE type traffic selector should has length 40 bytes")
				}
				if len(rawData) < int(selectorLength) {
					return errors.New("TrafficSelector: No sufficient bytes to decode next individual traffic selector")
				}

				individualTrafficSelector := &IndividualTrafficSelector{}

				individualTrafficSelector.TSType = rawData[0]
				individualTrafficSelector.IPProtocolID = rawData[1]
				individualTrafficSelector.StartPort = binary.BigEndian.Uint16(rawData[4:6])
				individualTrafficSelector.EndPort = binary.BigEndian.Uint16(rawData[6:8])

				individualTrafficSelector.StartAddress = append(individualTrafficSelector.StartAddress, rawData[8:24]...)
				individualTrafficSelector.EndAddress = append(individualTrafficSelector.EndAddress, rawData[24:40]...)

				trafficSelector.TrafficSelectors = append(trafficSelector.TrafficSelectors, individualTrafficSelector)

				rawData = rawData[40:]
			} else {
				return errors.New("TrafficSelector: Unsupported traffic selector type")
			}
		}
	}

	return nil
}

// Definition of Traffic Selector - Responder

var _ IKEPayloadType = &TrafficSelectorResponder{}

type TrafficSelectorResponder struct {
	TrafficSelectors []*IndividualTrafficSelector
}

func (trafficSelector *TrafficSelectorResponder) Type() IKEType { return TypeTSr }

func (trafficSelector *TrafficSelectorResponder) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][TrafficSelector] marshal(): Start marshalling")

	if len(trafficSelector.TrafficSelectors) > 0 {
		trafficSelectorData := make([]byte, 4)
		trafficSelectorData[0] = uint8(len(trafficSelector.TrafficSelectors))

		for _, individualTrafficSelector := range trafficSelector.TrafficSelectors {
			if individualTrafficSelector.TSType == TS_IPV4_ADDR_RANGE {
				// Address length checking
				if len(individualTrafficSelector.StartAddress) != 4 {
					return nil, errors.New("TrafficSelector: Start IPv4 address length is not correct")
				}
				if len(individualTrafficSelector.EndAddress) != 4 {
					return nil, errors.New("TrafficSelector: End IPv4 address length is not correct")
				}

				individualTrafficSelectorData := make([]byte, 8)

				individualTrafficSelectorData[0] = individualTrafficSelector.TSType
				individualTrafficSelectorData[1] = individualTrafficSelector.IPProtocolID
				binary.BigEndian.PutUint16(individualTrafficSelectorData[4:6], individualTrafficSelector.StartPort)
				binary.BigEndian.PutUint16(individualTrafficSelectorData[6:8], individualTrafficSelector.EndPort)

				individualTrafficSelectorData = append(individualTrafficSelectorData, individualTrafficSelector.StartAddress...)
				individualTrafficSelectorData = append(individualTrafficSelectorData, individualTrafficSelector.EndAddress...)

				binary.BigEndian.PutUint16(individualTrafficSelectorData[2:4], uint16(len(individualTrafficSelectorData)))

				trafficSelectorData = append(trafficSelectorData, individualTrafficSelectorData...)
			} else if individualTrafficSelector.TSType == TS_IPV6_ADDR_RANGE {
				// Address length checking
				if len(individualTrafficSelector.StartAddress) != 16 {
					return nil, errors.New("TrafficSelector: Start IPv6 address length is not correct")
				}
				if len(individualTrafficSelector.EndAddress) != 16 {
					return nil, errors.New("TrafficSelector: End IPv6 address length is not correct")
				}

				individualTrafficSelectorData := make([]byte, 8)

				individualTrafficSelectorData[0] = individualTrafficSelector.TSType
				individualTrafficSelectorData[1] = individualTrafficSelector.IPProtocolID
				binary.BigEndian.PutUint16(individualTrafficSelectorData[4:6], individualTrafficSelector.StartPort)
				binary.BigEndian.PutUint16(individualTrafficSelectorData[6:8], individualTrafficSelector.EndPort)

				individualTrafficSelectorData = append(individualTrafficSelectorData, individualTrafficSelector.StartAddress...)
				individualTrafficSelectorData = append(individualTrafficSelectorData, individualTrafficSelector.EndAddress...)

				binary.BigEndian.PutUint16(individualTrafficSelectorData[2:4], uint16(len(individualTrafficSelectorData)))

				trafficSelectorData = append(trafficSelectorData, individualTrafficSelectorData...)
			} else {
				return nil, errors.New("TrafficSelector: Unsupported traffic selector type")
			}
		}

		return trafficSelectorData, nil
	} else {
		return nil, errors.New("TrafficSelector: Contains no traffic selector for marshalling message")
	}
}

func (trafficSelector *TrafficSelectorResponder) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][TrafficSelector] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][TrafficSelector] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][TrafficSelector] unmarshal(): Unmarshal 1 traffic selector")
		// bounds checking
		if len(rawData) < 4 {
			return errors.New("TrafficSelector: No sufficient bytes to get number of traffic selector in header")
		}

		numberOfSPI := rawData[0]

		rawData = rawData[4:]

		for ; numberOfSPI > 0; numberOfSPI-- {
			// bounds checking
			if len(rawData) < 4 {
				return errors.New("TrafficSelector: No sufficient bytes to decode next individual traffic selector length in header")
			}
			trafficSelectorType := rawData[0]
			if trafficSelectorType == TS_IPV4_ADDR_RANGE {
				selectorLength := binary.BigEndian.Uint16(rawData[2:4])
				if selectorLength != 16 {
					return errors.New("TrafficSelector: A TS_IPV4_ADDR_RANGE type traffic selector should has length 16 bytes")
				}
				if len(rawData) < int(selectorLength) {
					return errors.New("TrafficSelector: No sufficient bytes to decode next individual traffic selector")
				}

				individualTrafficSelector := &IndividualTrafficSelector{}

				individualTrafficSelector.TSType = rawData[0]
				individualTrafficSelector.IPProtocolID = rawData[1]
				individualTrafficSelector.StartPort = binary.BigEndian.Uint16(rawData[4:6])
				individualTrafficSelector.EndPort = binary.BigEndian.Uint16(rawData[6:8])

				individualTrafficSelector.StartAddress = append(individualTrafficSelector.StartAddress, rawData[8:12]...)
				individualTrafficSelector.EndAddress = append(individualTrafficSelector.EndAddress, rawData[12:16]...)

				trafficSelector.TrafficSelectors = append(trafficSelector.TrafficSelectors, individualTrafficSelector)

				rawData = rawData[16:]
			} else if trafficSelectorType == TS_IPV6_ADDR_RANGE {
				selectorLength := binary.BigEndian.Uint16(rawData[2:4])
				if selectorLength != 40 {
					return errors.New("TrafficSelector: A TS_IPV6_ADDR_RANGE type traffic selector should has length 40 bytes")
				}
				if len(rawData) < int(selectorLength) {
					return errors.New("TrafficSelector: No sufficient bytes to decode next individual traffic selector")
				}

				individualTrafficSelector := &IndividualTrafficSelector{}

				individualTrafficSelector.TSType = rawData[0]
				individualTrafficSelector.IPProtocolID = rawData[1]
				individualTrafficSelector.StartPort = binary.BigEndian.Uint16(rawData[4:6])
				individualTrafficSelector.EndPort = binary.BigEndian.Uint16(rawData[6:8])

				individualTrafficSelector.StartAddress = append(individualTrafficSelector.StartAddress, rawData[8:24]...)
				individualTrafficSelector.EndAddress = append(individualTrafficSelector.EndAddress, rawData[24:40]...)

				trafficSelector.TrafficSelectors = append(trafficSelector.TrafficSelectors, individualTrafficSelector)

				rawData = rawData[40:]
			} else {
				return errors.New("TrafficSelector: Unsupported traffic selector type")
			}
		}
	}

	return nil
}

// Definition of Encrypted Payload

var _ IKEPayloadType = &Encrypted{}

type Encrypted struct {
	NextPayload   uint8
	EncryptedData []byte
}

func (encrypted *Encrypted) Type() IKEType { return TypeSK }

func (encrypted *Encrypted) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Encrypted] marshal(): Start marshalling")

	if len(encrypted.EncryptedData) == 0 {
		ikeLog.Warn("[IKE][Encrypted] The encrypted data is empty")
	}

	return encrypted.EncryptedData, nil
}

func (encrypted *Encrypted) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Encrypted] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Encrypted] unmarshal(): Payload length %d bytes", len(rawData))
	encrypted.EncryptedData = append(encrypted.EncryptedData, rawData...)
	return nil
}

// Definition of Configuration

var _ IKEPayloadType = &Configuration{}

type Configuration struct {
	ConfigurationType      uint8
	ConfigurationAttribute []*IndividualConfigurationAttribute
}

type IndividualConfigurationAttribute struct {
	Type  uint16
	Value []byte
}

func (configuration *Configuration) Type() IKEType { return TypeCP }

func (configuration *Configuration) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][Configuration] marshal(): Start marshalling")

	configurationData := make([]byte, 4)

	configurationData[0] = configuration.ConfigurationType

	for _, attribute := range configuration.ConfigurationAttribute {
		individualConfigurationAttributeData := make([]byte, 4)

		binary.BigEndian.PutUint16(individualConfigurationAttributeData[0:2], (attribute.Type & 0x7fff))
		binary.BigEndian.PutUint16(individualConfigurationAttributeData[2:4], uint16(len(attribute.Value)))

		individualConfigurationAttributeData = append(individualConfigurationAttributeData, attribute.Value...)

		configurationData = append(configurationData, individualConfigurationAttributeData...)
	}

	return configurationData, nil
}

func (configuration *Configuration) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][Configuration] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][Configuration] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][Configuration] unmarshal(): Unmarshal 1 configuration")
		// bounds checking
		if len(rawData) <= 4 {
			return errors.New("Configuration: No sufficient bytes to decode next configuration")
		}
		configuration.ConfigurationType = rawData[0]

		configurationAttributeData := rawData[4:]

		for len(configurationAttributeData) > 0 {
			ikeLog.Trace("[IKE][Configuration] unmarshal(): Unmarshal 1 configuration attribute")
			// bounds checking
			if len(configurationAttributeData) < 4 {
				return errors.New("ConfigurationAttribute: No sufficient bytes to decode next configuration attribute")
			}
			length := binary.BigEndian.Uint16(configurationAttributeData[2:4])
			if len(configurationAttributeData) < int(4+length) {
				return errors.New("ConfigurationAttribute: TLV attribute length error")
			}

			individualConfigurationAttribute := new(IndividualConfigurationAttribute)

			individualConfigurationAttribute.Type = binary.BigEndian.Uint16(configurationAttributeData[0:2])
			configurationAttributeData = configurationAttributeData[4:]
			individualConfigurationAttribute.Value = append(individualConfigurationAttribute.Value, configurationAttributeData[:length]...)
			configurationAttributeData = configurationAttributeData[length:]

			configuration.ConfigurationAttribute = append(configuration.ConfigurationAttribute, individualConfigurationAttribute)
		}
	}

	return nil
}

// Definition of IKE EAP

var _ IKEPayloadType = &EAP{}

type EAP struct {
	Code        uint8
	Identifier  uint8
	EAPTypeData []EAPTypeFormat
}

func (eap *EAP) Type() IKEType { return TypeEAP }

func (eap *EAP) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][EAP] marshal(): Start marshalling")

	eapData := make([]byte, 4)

	eapData[0] = eap.Code
	eapData[1] = eap.Identifier

	if len(eap.EAPTypeData) > 0 {
		eapTypeData, err := eap.EAPTypeData[0].marshal()
		if err != nil {
			return nil, fmt.Errorf("EAP: EAP type data marshal failed: %+v", err)
		}

		eapData = append(eapData, eapTypeData...)
	}

	binary.BigEndian.PutUint16(eapData[2:4], uint16(len(eapData)))

	return eapData, nil
}

func (eap *EAP) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][EAP] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][EAP] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		ikeLog.Trace("[IKE][EAP] unmarshal(): Unmarshal 1 EAP")
		// bounds checking
		if len(rawData) < 4 {
			return errors.New("EAP: No sufficient bytes to decode next EAP payload")
		}
		eapPayloadLength := binary.BigEndian.Uint16(rawData[2:4])
		if eapPayloadLength < 4 {
			return errors.New("EAP: Payload length specified in the header is too small for EAP")
		}
		if len(rawData) != int(eapPayloadLength) {
			return errors.New("EAP: Received payload length not matches the length specified in header")
		}

		eap.Code = rawData[0]
		eap.Identifier = rawData[1]

		// EAP Success or Failed
		if eapPayloadLength == 4 {
			return nil
		}

		eapType := rawData[4]
		var eapTypeData EAPTypeFormat

		switch eapType {
		case EAPTypeIdentity:
			eapTypeData = new(EAPIdentity)
		case EAPTypeNotification:
			eapTypeData = new(EAPNotification)
		case EAPTypeNak:
			eapTypeData = new(EAPNak)
		case EAPTypeExpanded:
			eapTypeData = new(EAPExpanded)
		default:
			// TODO: Create unsupprted type to handle it
			return errors.New("EAP: Not supported EAP type")
		}

		if err := eapTypeData.unmarshal(rawData[4:]); err != nil {
			return fmt.Errorf("EAP: Unamrshal EAP type data failed: %+v", err)
		}

		eap.EAPTypeData = append(eap.EAPTypeData, eapTypeData)

	}

	return nil
}

type EAPTypeFormat interface {
	// Type specifies EAP types
	Type() EAPType

	// Called by EAP.marshal() or EAP.unmarshal()
	marshal() ([]byte, error)
	unmarshal(rawData []byte) error
}

// Definition of EAP Identity

var _ EAPTypeFormat = &EAPIdentity{}

type EAPIdentity struct {
	IdentityData []byte
}

func (eapIdentity *EAPIdentity) Type() EAPType { return EAPTypeIdentity }

func (eapIdentity *EAPIdentity) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][EAP][Identity] marshal(): Start marshalling")

	if len(eapIdentity.IdentityData) == 0 {
		return nil, errors.New("EAPIdentity: EAP identity is empty")
	}

	eapIdentityData := []byte{EAPTypeIdentity}
	eapIdentityData = append(eapIdentityData, eapIdentity.IdentityData...)

	return eapIdentityData, nil
}

func (eapIdentity *EAPIdentity) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][EAP][Identity] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][EAP][Identity] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 1 {
		eapIdentity.IdentityData = append(eapIdentity.IdentityData, rawData[1:]...)
	}

	return nil
}

// Definition of EAP Notification

var _ EAPTypeFormat = &EAPNotification{}

type EAPNotification struct {
	NotificationData []byte
}

func (eapNotification *EAPNotification) Type() EAPType { return EAPTypeNotification }

func (eapNotification *EAPNotification) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][EAP][Notification] marshal(): Start marshalling")

	if len(eapNotification.NotificationData) == 0 {
		return nil, errors.New("EAPNotification: EAP notification is empty")
	}

	eapNotificationData := []byte{EAPTypeNotification}
	eapNotificationData = append(eapNotificationData, eapNotification.NotificationData...)

	return eapNotificationData, nil
}

func (eapNotification *EAPNotification) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][EAP][Notification] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][EAP][Notification] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 1 {
		eapNotification.NotificationData = append(eapNotification.NotificationData, rawData[1:]...)
	}

	return nil
}

// Definition of EAP Nak

var _ EAPTypeFormat = &EAPNak{}

type EAPNak struct {
	NakData []byte
}

func (eapNak *EAPNak) Type() EAPType { return EAPTypeNak }

func (eapNak *EAPNak) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][EAP][Nak] marshal(): Start marshalling")

	if len(eapNak.NakData) == 0 {
		return nil, errors.New("EAPNak: EAP nak is empty")
	}

	eapNakData := []byte{EAPTypeNak}
	eapNakData = append(eapNakData, eapNak.NakData...)

	return eapNakData, nil
}

func (eapNak *EAPNak) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][EAP][Nak] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][EAP][Nak] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 1 {
		eapNak.NakData = append(eapNak.NakData, rawData[1:]...)
	}

	return nil
}

// Definition of EAP expanded

var _ EAPTypeFormat = &EAPExpanded{}

type EAPExpanded struct {
	VendorID   uint32
	VendorType uint32
	VendorData []byte
}

func (eapExpanded *EAPExpanded) Type() EAPType { return EAPTypeExpanded }

func (eapExpanded *EAPExpanded) marshal() ([]byte, error) {
	ikeLog.Info("[IKE][EAP][Expanded] marshal(): Start marshalling")

	eapExpandedData := make([]byte, 8)

	vendorID := eapExpanded.VendorID & 0x00ffffff
	typeAndVendorID := (uint32(EAPTypeExpanded)<<24 | vendorID)

	binary.BigEndian.PutUint32(eapExpandedData[0:4], typeAndVendorID)
	binary.BigEndian.PutUint32(eapExpandedData[4:8], eapExpanded.VendorType)

	if len(eapExpanded.VendorData) == 0 {
		ikeLog.Warn("[IKE][EAP][Expanded] marshal(): EAP vendor data field is empty")
		return eapExpandedData, nil
	}

	eapExpandedData = append(eapExpandedData, eapExpanded.VendorData...)

	return eapExpandedData, nil
}

func (eapExpanded *EAPExpanded) unmarshal(rawData []byte) error {
	ikeLog.Info("[IKE][EAP][Expanded] unmarshal(): Start unmarshalling received bytes")
	ikeLog.Tracef("[IKE][EAP][Expanded] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		if len(rawData) < 8 {
			return errors.New("EAPExpanded: No sufficient bytes to decode the EAP expanded type")
		}

		typeAndVendorID := binary.BigEndian.Uint32(rawData[0:4])
		eapExpanded.VendorID = typeAndVendorID & 0x00ffffff

		eapExpanded.VendorType = binary.BigEndian.Uint32(rawData[4:8])

		if len(rawData) > 8 {
			eapExpanded.VendorData = append(eapExpanded.VendorData, rawData[8:]...)
		}
	}

	return nil
}
