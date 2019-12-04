package nas_security

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/aead/cmac"
	"free5gc/lib/nas"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/logger"
	"reflect"
)

func Encode(ue *amf_context.AmfUe, msg *nas.Message) (payload []byte, err error) {
	integrityProtected := false
	newSecurityContext := false
	ciphering := false
	if ue == nil {
		err = fmt.Errorf("amfUe is nil")
		return
	}
	if msg == nil {
		err = fmt.Errorf("Nas Message is empty")
		return
	}
	switch msg.SecurityHeader.SecurityHeaderType {
	case nas.SecurityHeaderTypePlainNas:
		logger.NasLog.Infoln("NasPdu Security: Plain Nas")
		return msg.PlainNasEncode()
	case nas.SecurityHeaderTypeIntegrityProtected:
		logger.NasLog.Infoln("NasPdu Security: Integrity Protected")
		integrityProtected = true
	case nas.SecurityHeaderTypeIntegrityProtectedAndCiphered:
		logger.NasLog.Infoln("NasPdu Security: Integrity Protected And Ciphered")
		integrityProtected = true
		ciphering = true
	case nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext:
		logger.NasLog.Infoln("NasPdu Security: Integrity Protected With New 5gNasSecurityContext")
		integrityProtected = true
		newSecurityContext = true
	case nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext:
		logger.NasLog.Infoln("NasPdu Security: Integrity Protected And Ciphered WithNew 5gNasSecurityContext")
		integrityProtected = true
		ciphering = true
		newSecurityContext = true
	default:
		return nil, fmt.Errorf("Security Type[%d] is not be implemented", msg.SecurityHeader.SecurityHeaderType)
	}

	if newSecurityContext {
		ue.DLCount = 0
		ue.ULCountOverflow = 0
		ue.ULCountSQN = 0
	}
	if ue.CipheringAlg == amf_context.ALG_CIPHERING_128_NEA0 {
		ciphering = false
	}
	if ue.IntegrityAlg == amf_context.ALG_INTEGRITY_128_NIA0 {
		integrityProtected = false
	}
	if ciphering || integrityProtected {
		securityHeader := []byte{msg.SecurityHeader.ProtocolDiscriminator, msg.SecurityHeaderType}
		sequenceNumber := uint8(ue.DLCount & 0xff)

		payload, err = msg.PlainNasEncode()
		if err != nil {
			return
		}
		if ciphering {
			// TODO: Support for ue has nas connection in both accessType
			if err = NasEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.GetSecurityDLCount(), amf_context.SECURITY_ONLY_ONE_BEARER,
				amf_context.SECURITY_DIRECTION_DOWNLINK, payload); err != nil {
				return
			}
		}
		// add sequece number
		payload = append([]byte{sequenceNumber}, payload[:]...)
		mac32 := make([]byte, 4)
		if integrityProtected {
			mac32, err = NasMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.GetSecurityDLCount(), amf_context.SECURITY_ONLY_ONE_BEARER, amf_context.SECURITY_DIRECTION_DOWNLINK, payload)
			if err != nil {
				return
			}
		}
		// Add mac value
		payload = append(mac32, payload[:]...)
		// Add EPD and Security Type
		payload = append(securityHeader, payload[:]...)

		// Increase DL Count
		ue.DLCount = (ue.DLCount + 1) & 0xffffff

		ue.SecurityContextAvailable = true
	} else {
		err = fmt.Errorf("NEA0 & NIA0 are illegal.")
		//  err = msg.PlainNasEncode()
		return
	}
	return
}

func Decode(ue *amf_context.AmfUe, securityHeaderType uint8, payload []byte) (msg *nas.Message, err error) {

	integrityProtected := false
	newSecurityContext := false
	ciphering := false
	if ue == nil {
		err = fmt.Errorf("amfUe is nil")
		return
	}
	if payload == nil {
		err = fmt.Errorf("Nas payload is empty")
		return
	}

	switch securityHeaderType {
	case nas.SecurityHeaderTypePlainNas:
	case nas.SecurityHeaderTypeIntegrityProtected:
		integrityProtected = true
	case nas.SecurityHeaderTypeIntegrityProtectedAndCiphered:
		integrityProtected = true
		ciphering = true
	case nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext:
		integrityProtected = true
		newSecurityContext = true
	case nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext:
		integrityProtected = true
		ciphering = true
		newSecurityContext = true
	default:
		return nil, fmt.Errorf("Security Type[%d] is not be implemented", securityHeaderType)
	}
	msg = new(nas.Message)

	if !ue.SecurityContextAvailable {
		integrityProtected = false
		newSecurityContext = false
		ciphering = false
	}
	if newSecurityContext {
		ue.ULCountOverflow = 0
		ue.ULCountSQN = 0
	}
	if ue.CipheringAlg == amf_context.ALG_CIPHERING_128_NEA0 {
		ciphering = false
	}
	if ue.IntegrityAlg == amf_context.ALG_INTEGRITY_128_NIA0 {
		integrityProtected = false
	}
	if ciphering || integrityProtected {
		securityHeader := payload[0:6]
		sequenceNumber := payload[6]
		receivedMac32 := securityHeader[2:]
		// remove security Header except for sequece Number
		payload = payload[6:]

		// Caculate ul count
		if ue.ULCountSQN > sequenceNumber {
			ue.ULCountOverflow++
		}
		ue.ULCountSQN = sequenceNumber
		if integrityProtected {
			// ToDo: use real mac calculate
			mac32, err := NasMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.GetSecurityULCount(), amf_context.SECURITY_ONLY_ONE_BEARER,
				amf_context.SECURITY_DIRECTION_UPLINK, payload)
			if err != nil {
				ue.MacFailed = true
				return nil, err
			}
			if !reflect.DeepEqual(mac32, receivedMac32) {
				logger.NasLog.Warnf("NAS MAC verification failed(0x%x != 0x%x)", mac32, receivedMac32)
				ue.MacFailed = true
			}
		}
		// remove sequece Number
		payload = payload[1:]

		if ciphering {
			// TODO: Support for ue has nas connection in both accessType
			if err = NasEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.GetSecurityULCount(), amf_context.SECURITY_ONLY_ONE_BEARER,
				amf_context.SECURITY_DIRECTION_UPLINK, payload); err != nil {
				return
			}
		}
	}
	err = msg.PlainNasDecode(&payload)

	return
}

func NasEncrypt(AlgoID uint8, KnasEnc []byte, Count []byte, Bearer uint8, Direction uint8, plainText []byte) error {

	if len(KnasEnc) != 16 {
		return fmt.Errorf("Size of KnasEnc[%d] != 16 bytes)", len(KnasEnc))
	}
	if Bearer > 0x1f {
		return fmt.Errorf("Bearer is beyond 5 bits")
	}
	if Direction > 1 {
		return fmt.Errorf("Direction is beyond 1 bits")
	}
	if plainText == nil {
		return fmt.Errorf("Nas Payload is nil")
	}

	switch AlgoID {
	case amf_context.ALG_CIPHERING_128_NEA1:
		logger.NgapLog.Errorf("NEA1 not implement yet.")
		return nil
	case amf_context.ALG_CIPHERING_128_NEA2:
		// Couter[0..32] | BEARER[0..4] | DIRECTION[0] | 0^26 | 0^64
		CouterBlk := make([]byte, 16)
		//First 32 bits are count
		copy(CouterBlk, Count)
		//Put Bearer and direction together
		CouterBlk[4] = (Bearer << 3) | (Direction << 2)

		block, err := aes.NewCipher(KnasEnc)
		if err != nil {
			return err
		}

		ciphertext := make([]byte, len(plainText))

		stream := cipher.NewCTR(block, CouterBlk)
		stream.XORKeyStream(ciphertext, plainText)
		// override plainText with cipherText
		copy(plainText, ciphertext)
		return nil

	case amf_context.ALG_CIPHERING_128_NEA3:
		logger.NgapLog.Errorf("NEA3 not implement yet.")
		return nil
	default:
		return fmt.Errorf("Unknown Algorithm Identity[%d]", AlgoID)
	}

}

func NasMacCalculate(AlgoID uint8, KnasInt []byte, Count []byte, Bearer uint8, Direction uint8, msg []byte) ([]byte, error) {
	if len(KnasInt) != 16 {
		return nil, fmt.Errorf("Size of KnasEnc[%d] != 16 bytes)", len(KnasInt))
	}
	if Bearer > 0x1f {
		return nil, fmt.Errorf("Bearer is beyond 5 bits")
	}
	if Direction > 1 {
		return nil, fmt.Errorf("Direction is beyond 1 bits")
	}
	if msg == nil {
		return nil, fmt.Errorf("Nas Payload is nil")
	}

	switch AlgoID {
	case amf_context.ALG_INTEGRITY_128_NIA1:
		logger.NgapLog.Errorf("NEA1 not implement yet.")
		return nil, nil
	case amf_context.ALG_INTEGRITY_128_NIA2:
		// Couter[0..32] | BEARER[0..4] | DIRECTION[0] | 0^26
		m := make([]byte, len(msg)+8)
		//First 32 bits are count
		copy(m, Count)
		//Put Bearer and direction together
		m[4] = (Bearer << 3) | (Direction << 2)

		block, err := aes.NewCipher(KnasInt)
		if err != nil {
			return nil, err
		}

		copy(m[8:], msg)

		cmac, err := cmac.Sum(m, block, 16)
		if err != nil {
			return nil, err
		}
		// only get the most significant 32 bits to be mac value
		return cmac[:4], nil

	case amf_context.ALG_INTEGRITY_128_NIA3:
		logger.NgapLog.Errorf("NEA3 not implement yet.")
		return nil, nil
	default:
		return nil, fmt.Errorf("Unknown Algorithm Identity[%d]", AlgoID)
	}

}
