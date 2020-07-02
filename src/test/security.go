package test

import (
	"fmt"
	"free5gc/lib/nas"
	"free5gc/lib/nas/security"
	"reflect"
)

func NASEncode(ue *RanUeContext, msg *nas.Message, securityContextAvailable bool, newSecurityContext bool) (payload []byte, err error) {
	var sequenceNumber uint8
	if ue == nil {
		err = fmt.Errorf("amfUe is nil")
		return
	}
	if msg == nil {
		err = fmt.Errorf("Nas Message is empty")
		return
	}

	if !securityContextAvailable {
		return msg.PlainNasEncode()
	} else {
		if newSecurityContext {
			ue.ULCount = 0
			ue.DLOverflow = 0
			ue.DLCountSQN = 0
		}

		sequenceNumber = uint8(ue.ULCount & 0xff)
		payload, err = msg.PlainNasEncode()
		if err != nil {
			return
		}

		// TODO: Support for ue has nas connection in both accessType
		if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.GetSecurityULCount(), security.SecurityBearer3GPP,
			security.SecurityDirectionUplink, payload); err != nil {
			return
		}
		// add sequece number
		payload = append([]byte{sequenceNumber}, payload[:]...)
		mac32 := make([]byte, 4)

		mac32, err = security.NASMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.GetSecurityULCount(), security.SecurityBearer3GPP, security.SecurityDirectionUplink, payload)
		if err != nil {
			return
		}

		// Add mac value
		payload = append(mac32, payload[:]...)
		// Add EPD and Security Type
		msgSecurityHeader := []byte{msg.SecurityHeader.ProtocolDiscriminator, msg.SecurityHeader.SecurityHeaderType}
		payload = append(msgSecurityHeader, payload[:]...)

		// Increase UL Count
		ue.ULCount = (ue.ULCount + 1) & 0xffffff
	}
	return
}

func NASDecode(ue *RanUeContext, securityHeaderType uint8, payload []byte) (msg *nas.Message, err error) {
	if ue == nil {
		err = fmt.Errorf("amfUe is nil")
		return
	}
	if payload == nil {
		err = fmt.Errorf("Nas payload is empty")
		return
	}

	msg = new(nas.Message)

	if securityHeaderType == nas.SecurityHeaderTypePlainNas {
		err = msg.PlainNasDecode(&payload)
		return
	} else if ue.IntegrityAlg == security.AlgIntegrity128NIA0 {
		fmt.Println("decode payload is ", payload)
		// remove header
		payload = payload[3:]

		if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.GetSecurityULCount(), security.SecurityBearer3GPP,
			security.SecurityDirectionDownlink, payload); err != nil {
			return nil, err
		}

		err = msg.PlainNasDecode(&payload)
		return
	} else {
		if securityHeaderType == nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext || securityHeaderType == nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext {
			ue.DLOverflow = 0
			ue.DLCountSQN = 0
		}

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
		// ToDo: use real mac calculate
		if ue.IntegrityAlg != security.AlgIntegrity128NIA0 {
			mac32, err := security.NASMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.GetSecurityULCount(), security.SecurityBearer3GPP,
				security.SecurityDirectionDownlink, payload)
			if err != nil {
				return nil, err
			}
			if !reflect.DeepEqual(mac32, receivedMac32) {
				fmt.Printf("NAS MAC verification failed(0x%x != 0x%x)", mac32, receivedMac32)
			} else {
				fmt.Printf("cmac value: 0x%x\n", mac32)
			}
		}

		// remove sequece Number
		payload = payload[1:]

		// TODO: Support for ue has nas connection in both accessType
		if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.GetSecurityULCount(), security.SecurityBearer3GPP,
			security.SecurityDirectionUplink, payload); err != nil {
			return nil, err
		}
	}
	err = msg.PlainNasDecode(&payload)
	fmt.Println("err", err)
	return

}
