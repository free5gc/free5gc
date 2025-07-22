package nas_security

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
)

func Encode(ue *context.AmfUe, msg *nas.Message, accessType models.AccessType) ([]byte, error) {
	if msg == nil {
		return nil, fmt.Errorf("NAS Message is nil")
	}

	// Plain NAS message
	if ue == nil || !ue.SecurityContextAvailable {
		if msg.GmmMessage == nil {
			return nil, fmt.Errorf("msg.GmmMessage is nil")
		}
		switch msgType := msg.GmmHeader.GetMessageType(); msgType {
		case nas.MsgTypeIdentityRequest:
			if msg.GmmMessage.IdentityRequest == nil {
				return nil,
					fmt.Errorf("identity Request (type unknown) is requierd security, but security context is not available")
			}
			if identityType := msg.GmmMessage.IdentityRequest.SpareHalfOctetAndIdentityType.GetTypeOfIdentity(); identityType !=
				nasMessage.MobileIdentity5GSTypeSuci {
				return nil,
					fmt.Errorf("identity Request (%d) is requierd security, but security context is not available", identityType)
			}
		case nas.MsgTypeAuthenticationRequest:
		case nas.MsgTypeAuthenticationResult:
		case nas.MsgTypeAuthenticationReject:
		case nas.MsgTypeRegistrationReject:
		case nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration:
		case nas.MsgTypeServiceReject:
		default:
			return nil, fmt.Errorf("NAS message type %d is requierd security, but security context is not available", msgType)
		}
		pdu, err := msg.PlainNasEncode()
		return pdu, err
	} else {
		// Security protected NAS Message
		// a security protected NAS message must be integrity protected, and ciphering is optional
		needCiphering := false
		switch msg.SecurityHeader.SecurityHeaderType {
		case nas.SecurityHeaderTypeIntegrityProtected:
			ue.NASLog.Debugln("Security header type: Integrity Protected")
		case nas.SecurityHeaderTypeIntegrityProtectedAndCiphered:
			ue.NASLog.Debugln("Security header type: Integrity Protected And Ciphered")
			needCiphering = true
		case nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext:
			ue.NASLog.Debugln("Security header type: Integrity Protected With New 5G Security Context")
			ue.ULCount.Set(0, 0)
			ue.DLCount.Set(0, 0)
		default:
			return nil, fmt.Errorf("wrong security header type: 0x%0x", msg.SecurityHeader.SecurityHeaderType)
		}

		// encode plain nas first
		payload, err := msg.PlainNasEncode()
		if err != nil {
			return nil, fmt.Errorf("plain NAS encode error: %+v", err)
		}

		ue.NASLog.Tracef("plain payload:\n%+v", hex.Dump(payload))
		if needCiphering {
			ue.NASLog.Debugf("Encrypt NAS message (algorithm: %+v, DLCount: 0x%0x)", ue.CipheringAlg, ue.DLCount.Get())
			ue.NASLog.Tracef("NAS ciphering key: %0x", ue.KnasEnc)
			if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.DLCount.Get(),
				GetBearerType(accessType), security.DirectionDownlink, payload); err != nil {
				return nil, fmt.Errorf("encrypt error: %+v", err)
			}
		}

		// add sequece number
		addsqn := []byte{}
		addsqn = append(addsqn, []byte{ue.DLCount.SQN()}...)
		addsqn = append(addsqn, payload...)
		payload = addsqn

		ue.NASLog.Debugf("Calculate NAS MAC (algorithm: %+v, DLCount: 0x%0x)", ue.IntegrityAlg, ue.DLCount.Get())
		ue.NASLog.Tracef("NAS integrity key: %0x", ue.KnasInt)
		mac32, err := security.NASMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.DLCount.Get(),
			GetBearerType(accessType), security.DirectionDownlink, payload)
		if err != nil {
			return nil, fmt.Errorf("MAC calcuate error: %+v", err)
		}
		// Add mac value
		ue.NASLog.Tracef("MAC: 0x%08x", mac32)
		addmac := []byte{}
		addmac = append(addmac, mac32...)
		addmac = append(addmac, payload...)
		payload = addmac

		// Add EPD and Security Type
		msgSecurityHeader := []byte{msg.SecurityHeader.ProtocolDiscriminator, msg.SecurityHeader.SecurityHeaderType}
		encodepayload := []byte{}
		encodepayload = append(encodepayload, msgSecurityHeader...)
		encodepayload = append(encodepayload, payload...)
		payload = encodepayload

		// Increase DL Count
		ue.DLCount.AddOne()
		return payload, nil
	}
}

/*
payload either a security protected 5GS NAS message or a plain 5GS NAS message which
format is followed TS 24.501 9.1.1
*/
func Decode(ue *context.AmfUe, accessType models.AccessType, payload []byte,
	initialMessage bool,
) (msg *nas.Message, integrityProtected bool, err error) {
	if ue == nil {
		return nil, false, fmt.Errorf("amfUe is nil")
	}
	if payload == nil {
		return nil, false, fmt.Errorf("NAS payload is empty")
	}
	if len(payload) < 2 {
		return nil, false, fmt.Errorf("NAS payload is too short")
	}

	ulCountNew := ue.ULCount

	msg = new(nas.Message)
	msg.ProtocolDiscriminator = payload[0]
	msg.SecurityHeaderType = nas.GetSecurityHeaderType(payload) & 0x0f
	ue.NASLog.Traceln("securityHeaderType is ", msg.SecurityHeaderType)
	if msg.SecurityHeaderType != nas.SecurityHeaderTypePlainNas { // Security protected NAS message
		// Extended protocol discriminator	V 1
		// Security header type				V 1/2
		// Spare half octet					V 1/2
		// Message authentication code		V 4
		// Sequence number					V 1
		// Plain 5GS NAS message			V 3-n
		if len(payload) < (1 + 1 + 4 + 1 + 3) {
			return nil, false, fmt.Errorf("NAS payload is too short")
		}
		securityHeader := payload[0:6]
		ue.NASLog.Traceln("securityHeader is ", securityHeader)
		sequenceNumber := payload[6]
		ue.NASLog.Traceln("sequenceNumber", sequenceNumber)
		msg.SequenceNumber = sequenceNumber
		receivedMac32 := securityHeader[2:]
		msg.MessageAuthenticationCode = binary.BigEndian.Uint32(receivedMac32)
		// remove security Header except for sequece Number
		payload = payload[6:]

		// a security protected NAS message must be integrity protected, and ciphering is optional
		ciphered := false
		switch msg.SecurityHeaderType {
		case nas.SecurityHeaderTypeIntegrityProtected:
			ue.NASLog.Debugln("Security header type: Integrity Protected")
		case nas.SecurityHeaderTypeIntegrityProtectedAndCiphered:
			ue.NASLog.Debugln("Security header type: Integrity Protected And Ciphered")
			ciphered = true
		case nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext:
			ue.NASLog.Debugln("Security header type: Integrity Protected And Ciphered With New 5G Security Context")
			ciphered = true
			ulCountNew.Set(0, 0)
		default:
			return nil, false, fmt.Errorf("wrong security header type: 0x%0x", msg.SecurityHeader.SecurityHeaderType)
		}

		if ciphered && !ue.SecurityContextAvailable {
			return nil, false, fmt.Errorf("NAS message is ciphered, but UE Security Context is not Available")
		}

		if ue.SecurityContextAvailable {
			if ulCountNew.SQN() > sequenceNumber {
				ue.NASLog.Debugf("set ULCount overflow")
				ulCountNew.SetOverflow(ulCountNew.Overflow() + 1)
			}
			ulCountNew.SetSQN(sequenceNumber)

			ue.NASLog.Debugf("Calculate NAS MAC (algorithm: %+v, ULCount: 0x%0x)", ue.IntegrityAlg, ulCountNew.Get())
			ue.NASLog.Tracef("NAS integrity key0x: %0x", ue.KnasInt)
			var mac32 []byte
			mac32, err = security.NASMacCalculate(ue.IntegrityAlg, ue.KnasInt, ulCountNew.Get(),
				GetBearerType(accessType), security.DirectionUplink, payload)
			if err != nil {
				return nil, false, fmt.Errorf("MAC calcuate error: %+v", err)
			}

			if !reflect.DeepEqual(mac32, receivedMac32) {
				ue.NASLog.Warnf("NAS MAC verification failed(received: 0x%08x, expected: 0x%08x)", receivedMac32, mac32)
			} else {
				ue.NASLog.Tracef("cmac value: 0x%08x", mac32)
				integrityProtected = true
			}
		} else {
			ue.NASLog.Debugln("UE Security Context is not Available, so skip MAC verify")
		}

		if ciphered {
			if !integrityProtected {
				return nil, false, fmt.Errorf("NAS message is ciphered, but MAC verification failed")
			}
			ue.NASLog.Debugf("Decrypt NAS message (algorithm: %+v, ULCount: 0x%0x)", ue.CipheringAlg, ulCountNew.Get())
			ue.NASLog.Tracef("NAS ciphering key: %0x", ue.KnasEnc)
			// decrypt payload without sequence number (payload[1])
			if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ulCountNew.Get(), GetBearerType(accessType),
				security.DirectionUplink, payload[1:]); err != nil {
				return nil, false, fmt.Errorf("decrypt error: %+v", err)
			}
		}

		// remove sequece Number
		payload = payload[1:]
	}

	err = msg.PlainNasDecode(&payload)
	if err != nil {
		return nil, false, err
	}

	msgTypeText := func() string {
		if msg.GmmMessage == nil {
			return "Non GMM message"
		} else {
			return fmt.Sprintf(" message type %d", msg.GmmHeader.GetMessageType())
		}
	}
	errNoSecurityContext := func() error {
		return fmt.Errorf("UE Security Context is not Available, %s", msgTypeText())
	}
	errWrongSecurityHeader := func() error {
		return fmt.Errorf("wrong security header type: 0x%0x, %s", msg.SecurityHeader.SecurityHeaderType, msgTypeText())
	}
	errMacVerificationFailed := func() error {
		return fmt.Errorf("MAC verification failed, %s", msgTypeText())
	}

	if msg.GmmMessage == nil {
		if !ue.SecurityContextAvailable {
			return nil, false, errNoSecurityContext()
		}
		if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedAndCiphered {
			return nil, false, errWrongSecurityHeader()
		}
		if !integrityProtected {
			return nil, false, errMacVerificationFailed()
		}
	} else {
		switch msg.GmmHeader.GetMessageType() {
		case nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration, nas.MsgTypeRegistrationRequest:
			if initialMessage {
				if msg.SecurityHeaderType == nas.SecurityHeaderTypeIntegrityProtectedAndCiphered ||
					msg.SecurityHeaderType == nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext {
					return nil, false, errWrongSecurityHeader()
				}
			} else {
				if ue.SecurityContextAvailable {
					if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedAndCiphered {
						return nil, false, errWrongSecurityHeader()
					}
					if !integrityProtected {
						return nil, false, errMacVerificationFailed()
					}
				}
			}
		case nas.MsgTypeServiceRequest:
			if initialMessage {
				if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtected {
					return nil, false, errWrongSecurityHeader()
				}
			} else {
				if !ue.SecurityContextAvailable {
					return nil, false, errNoSecurityContext()
				}
				if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedAndCiphered {
					return nil, false, errWrongSecurityHeader()
				}
				if !integrityProtected {
					return nil, false, errMacVerificationFailed()
				}
			}
		case nas.MsgTypeIdentityResponse:
			mobileIdentityContents := msg.IdentityResponse.MobileIdentity.GetMobileIdentityContents()
			if len(mobileIdentityContents) >= 1 &&
				nasConvert.GetTypeOfIdentity(mobileIdentityContents[0]) == nasMessage.MobileIdentity5GSTypeSuci {
				// Identity is SUCI
				if ue.SecurityContextAvailable {
					if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedAndCiphered {
						return nil, false, errWrongSecurityHeader()
					}
					if !integrityProtected {
						return nil, false, errMacVerificationFailed()
					}
				}
			} else {
				// Identity is not SUCI
				if !ue.SecurityContextAvailable {
					return nil, false, errNoSecurityContext()
				}
				if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedAndCiphered {
					return nil, false, errWrongSecurityHeader()
				}
				if !integrityProtected {
					return nil, false, errMacVerificationFailed()
				}
			}
		case nas.MsgTypeAuthenticationResponse,
			nas.MsgTypeAuthenticationFailure,
			nas.MsgTypeSecurityModeReject,
			nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration:
			if ue.SecurityContextAvailable {
				if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedAndCiphered {
					return nil, false, errWrongSecurityHeader()
				}
				if !integrityProtected {
					return nil, false, errMacVerificationFailed()
				}
			}
		case nas.MsgTypeSecurityModeComplete:
			if !ue.SecurityContextAvailable {
				return nil, false, errNoSecurityContext()
			}
			if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext {
				return nil, false, errWrongSecurityHeader()
			}
			if !integrityProtected {
				return nil, false, errMacVerificationFailed()
			}
		default:
			if !ue.SecurityContextAvailable {
				return nil, false, errNoSecurityContext()
			}
			if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedAndCiphered {
				return nil, false, errWrongSecurityHeader()
			}
			if !integrityProtected {
				return nil, false, errMacVerificationFailed()
			}
		}
	}

	if integrityProtected {
		ue.ULCount = ulCountNew
	}
	return msg, integrityProtected, nil
}

// DecodePlainNas is used to decode plain nas.
// If nas pdu is ciphered, this function will return error message.
// return value is: *nas.Message
func DecodePlainNasNoIntegrityCheck(payload []byte) (*nas.Message, error) {
	const SecurityHeaderTypeMask uint8 = 0x0f

	if len(payload) == 0 {
		return nil, fmt.Errorf("nas payload is empty")
	}

	msg := new(nas.Message)
	msg.SecurityHeaderType = nas.GetSecurityHeaderType(payload) & SecurityHeaderTypeMask
	if msg.SecurityHeaderType == nas.SecurityHeaderTypeIntegrityProtectedAndCiphered ||
		msg.SecurityHeaderType == nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext {
		return nil, fmt.Errorf("nas payload is ciphered")
	}

	if msg.SecurityHeaderType != nas.SecurityHeaderTypePlainNas {
		// remove security Header
		if len(payload) < 7 {
			return nil, fmt.Errorf("nas payload is too short")
		}
		payload = payload[7:]
	}

	err := msg.PlainNasDecode(&payload)
	return msg, err
}

func GetBearerType(accessType models.AccessType) uint8 {
	switch accessType {
	case models.AccessType__3_GPP_ACCESS:
		return security.Bearer3GPP
	case models.AccessType_NON_3_GPP_ACCESS:
		return security.BearerNon3GPP
	default:
		return security.OnlyOneBearer
	}
}
