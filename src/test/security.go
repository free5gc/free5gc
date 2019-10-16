package test

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/aead/cmac"
	"free5gc/lib/nas"
	"reflect"
)

// TS 33501 Annex A.8 Algorithm distinguisher For Knas_int Knas_enc
const (
	N_NAS_ENC_ALG uint8 = 0x01
	N_NAS_INT_ALG uint8 = 0x02
	N_RRC_ENC_ALG uint8 = 0x03
	N_RRC_INT_ALG uint8 = 0x04
	N_UP_ENC_alg  uint8 = 0x05
	N_UP_INT_alg  uint8 = 0x06
)

// TS 33501 Annex D Algorithm identifier values For Knas_int
const (
	ALG_INTEGRITY_128_NIA0 uint8 = 0x00 // NULL
	ALG_INTEGRITY_128_NIA1 uint8 = 0x01 // 128-Snow3G
	ALG_INTEGRITY_128_NIA2 uint8 = 0x02 // 128-AES
	ALG_INTEGRITY_128_NIA3 uint8 = 0x03 // 128-ZUC
)

// TS 33501 Annex D Algorithm identifier values For Knas_enc
const (
	ALG_CIPHERING_128_NEA0 uint8 = 0x00 // NULL
	ALG_CIPHERING_128_NEA1 uint8 = 0x01 // 128-Snow3G
	ALG_CIPHERING_128_NEA2 uint8 = 0x02 // 128-AES
	ALG_CIPHERING_128_NEA3 uint8 = 0x03 // 128-ZUC
)

// 1bit
const (
	SECURITY_DIRECTION_UPLINK   uint8 = 0x00
	SECURITY_DIRECTION_DOWNLINK uint8 = 0x01
)

// 5bits
const (
	SECURITY_ONLY_ONE_BEARER uint8 = 0x00
	SECURITY_BEARER_3GPP     uint8 = 0x01
	SECURITY_BEARER_NON_3GPP uint8 = 0x02
)

// TS 33501 Annex A.0 Access type distinguisher For Kgnb Kn3iwf
const (
	ACCESS_TYPE_3GPP     uint8 = 0x01
	ACCESS_TYPE_NON_3GPP uint8 = 0x02
)

func NASEncode(ue *RanUeContext, msg *nas.Message) (payload []byte, err error) {
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
		fmt.Println("NasPdu Security: Plain Nas")
		return msg.PlainNasEncode()
	case nas.SecurityHeaderTypeIntegrityProtected:
		fmt.Println("NasPdu Security: Integrity Protected")
		integrityProtected = true
	case nas.SecurityHeaderTypeIntegrityProtectedAndCiphered:
		fmt.Println("NasPdu Security: Integrity Protected And Ciphered")
		integrityProtected = true
		ciphering = true
	case nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext:
		fmt.Println("NasPdu Security: Integrity Protected With New 5gNasSecurityContext")
		integrityProtected = true
		newSecurityContext = true
	case nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext:
		fmt.Println("NasPdu Security: Integrity Protected And Ciphered WithNew 5gNasSecurityContext")
		integrityProtected = true
		ciphering = true
		newSecurityContext = true
	default:
		return nil, fmt.Errorf("Security Type[%d] is not be implemented", msg.SecurityHeader.SecurityHeaderType)
	}

	if newSecurityContext {
		ue.ULCount = 0
		ue.DLOverflow = 0
		ue.DLCountSQN = 0
	}
	if ue.CipheringAlg == ALG_CIPHERING_128_NEA0 {
		ciphering = false
	}
	if ue.IntegrityAlg == ALG_INTEGRITY_128_NIA0 {
		integrityProtected = false
	}
	if ciphering || integrityProtected {
		securityHeader := []byte{msg.SecurityHeader.ProtocolDiscriminator, msg.SecurityHeaderType}
		sequenceNumber := uint8(ue.ULCount & 0xff)
		payload, err = msg.PlainNasEncode()
		if err != nil {
			return
		}
		if ciphering {
			// TODO: Support for ue has nas connection in both accessType
			if err = NasEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.GetSecurityULCount(), SECURITY_ONLY_ONE_BEARER,
				SECURITY_DIRECTION_UPLINK, payload); err != nil {
				return
			}

		}
		// add sequece number
		payload = append([]byte{sequenceNumber}, payload[:]...)
		mac32 := make([]byte, 4)
		if integrityProtected {
			mac32, err = NasMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.GetSecurityULCount(), SECURITY_ONLY_ONE_BEARER, SECURITY_DIRECTION_UPLINK, payload)
			if err != nil {
				return
			}
		}
		// Add mac value
		payload = append(mac32, payload[:]...)
		// Add EPD and Security Type
		payload = append(securityHeader, payload[:]...)

		// Increase UL Count
		ue.ULCount = (ue.ULCount + 1) & 0xffffff

	} else {
		return msg.PlainNasEncode()
	}
	return
}

func NASDecode(ue *RanUeContext, securityHeaderType uint8, payload []byte) (msg *nas.Message, err error) {

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

	if newSecurityContext {
		ue.DLOverflow = 0
		ue.DLCountSQN = 0
	}
	if ue.CipheringAlg == ALG_CIPHERING_128_NEA0 {
		ciphering = false
	}
	if ue.IntegrityAlg == ALG_INTEGRITY_128_NIA0 {
		integrityProtected = false
	}
	if ciphering || integrityProtected {
		securityHeader := payload[0:6]
		sequenceNumber := payload[6]
		receivedMac32 := securityHeader[2:]
		// remove security Header except for sequece Number
		payload = payload[6:]

		// Caculate ul count
		if ue.DLCountSQN > sequenceNumber {
			ue.DLOverflow++
		}
		ue.DLCountSQN = sequenceNumber
		if integrityProtected {
			// ToDo: use real mac calculate
			mac32, err := NasMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.GetSecurityDLCount(), SECURITY_ONLY_ONE_BEARER,
				SECURITY_DIRECTION_DOWNLINK, payload)
			if err != nil {
				return nil, err
			}
			if !reflect.DeepEqual(mac32, receivedMac32) {
				return nil, fmt.Errorf("NAS MAC verification failed(0x%x != 0x%x)", mac32, receivedMac32)
			}
		}
		// remove sequece Number
		payload = payload[1:]

		if ciphering {
			// TODO: Support for ue has nas connection in both accessType
			if err = NasEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.GetSecurityDLCount(), SECURITY_ONLY_ONE_BEARER,
				SECURITY_DIRECTION_DOWNLINK, payload); err != nil {
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
	case ALG_CIPHERING_128_NEA1:
		return fmt.Errorf("NEA3 not implement yet.")
	case ALG_CIPHERING_128_NEA2:
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

	case ALG_CIPHERING_128_NEA3:
		return fmt.Errorf("NEA3 not implement yet.")
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
	case ALG_INTEGRITY_128_NIA1:
		return nil, fmt.Errorf("NEA3 not implement yet.")
	case ALG_INTEGRITY_128_NIA2:
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

	case ALG_INTEGRITY_128_NIA3:
		return nil, fmt.Errorf("NEA3 not implement yet.")
	default:
		return nil, fmt.Errorf("Unknown Algorithm Identity[%d]", AlgoID)
	}

}
