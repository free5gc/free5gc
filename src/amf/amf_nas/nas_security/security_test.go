package nas_security_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/amf_nas/nas_security"
	"free5gc/src/amf/amf_ngap/ngap_message"
	"reflect"
	"strings"
	"testing"
)

func init() {
	go amf_handler.Handle()

	TestAmf.SctpSever()

}

func TestMacCalculate(t *testing.T) {
	key, err := hex.DecodeString(strings.Repeat("1", 32))
	if err != nil {
		t.Error(err.Error())
	}
	count := []byte{0x00, 0x01, 0x02, 0x03}
	var bearer uint8 = 0
	var direction uint8 = 1
	msg := []byte("hello world")
	if err != nil {
		t.Error(err.Error())
	}
	mac1, err := nas_security.NasMacCalculate(amf_context.ALG_INTEGRITY_128_NIA2, key, count, bearer, direction, msg)
	if err != nil {
		t.Error(err.Error())
	}
	mac2, err := nas_security.NasMacCalculate(amf_context.ALG_INTEGRITY_128_NIA2, key, count, bearer, direction, msg)
	if err != nil {
		t.Error(err.Error())
	} else if !reflect.DeepEqual(mac1, mac2) {
		t.Errorf("mac1[0x%x]\nmac2[0x%x]", mac1, mac2)
	}
}
func TestSecurity(t *testing.T) {
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.DerivateAlgKey()
	ue.DLCount = 4
	m := getRegistrationComplete(nil)
	nasPdu, err := nas_security.Encode(ue, m)
	if err != nil {
		t.Error(err.Error())
	}
	ngap_message.SendDownlinkNasTransport(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasPdu)
	msg, err := ranDecode(ue, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, nasPdu)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(msg.GmmMessage.RegistrationComplete, m.GmmMessage.RegistrationComplete) {
		t.Errorf("Expect: %s\n Output: %s", TestAmf.Config.Sdump(m.GmmMessage.RegistrationComplete), TestAmf.Config.Sdump(msg.GmmMessage.RegistrationComplete))
	}

}

func getRegistrationComplete(sorTransparentContainer []uint8) *nas.Message {

	m := nas.NewMessage()
	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
	}
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationComplete)

	registrationComplete := nasMessage.NewRegistrationComplete(0)
	registrationComplete.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationComplete.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	registrationComplete.RegistrationCompleteMessageIdentity.SetMessageType(nas.MsgTypeRegistrationComplete)

	if sorTransparentContainer != nil {
		registrationComplete.SORTransparentContainer = nasType.NewSORTransparentContainer(nasMessage.RegistrationCompleteSORTransparentContainerType)
		registrationComplete.SORTransparentContainer.SetLen(uint16(len(sorTransparentContainer)))
		registrationComplete.SORTransparentContainer.SetSORContent(sorTransparentContainer)
	}

	m.GmmMessage.RegistrationComplete = registrationComplete

	return m
}

func ranDecode(ue *amf_context.AmfUe, securityHeaderType uint8, payload []byte) (msg *nas.Message, err error) {

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
		// sequenceNumber := payload[6]
		receivedMac32 := securityHeader[2:]
		// remove security Header except for sequece Number
		payload = payload[6:]

		var dlcount = make([]byte, 4)
		binary.BigEndian.PutUint16(dlcount, uint16((ue.DLCount-1)&0xffffff))
		if integrityProtected {
			mac32, err := nas_security.NasMacCalculate(ue.IntegrityAlg, ue.KnasInt, dlcount, amf_context.SECURITY_ONLY_ONE_BEARER,
				amf_context.SECURITY_DIRECTION_DOWNLINK, payload)
			if err != nil {
				ue.MacFailed = true
				return nil, err
			}
			if !reflect.DeepEqual(mac32, receivedMac32) {
				fmt.Printf("NAS MAC verification failed(0x%x != 0x%x)", mac32, receivedMac32)
				ue.MacFailed = true
			} else {
				fmt.Printf("cmac value: 0x%x\n", mac32)
			}
		}
		// remove sequece Number
		payload = payload[1:]

		if ciphering {
			// TODO: Support for ue has nas connection in both accessType

			if err = nas_security.NasEncrypt(ue.CipheringAlg, ue.KnasEnc, dlcount, amf_context.SECURITY_ONLY_ONE_BEARER,
				amf_context.SECURITY_DIRECTION_DOWNLINK, payload); err != nil {
				return
			}
		}
	}
	err = msg.PlainNasDecode(&payload)

	return
}
