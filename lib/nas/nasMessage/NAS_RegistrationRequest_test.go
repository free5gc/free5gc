//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"fmt"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessageRegistrationRequestData struct {
	inExtendedProtocolDiscriminator       uint8
	inSecurityHeader                      uint8
	inSpareHalfOctet                      uint8
	inRegistrationRequestMessageIdentity  uint8
	inNgksi                               uint8
	inRegistrationType5GS                 uint8
	inMobileIdentity5GS                   nasType.MobileIdentity5GS
	inNoncurrentNativeNASKeySetIdentifier nasType.NoncurrentNativeNASKeySetIdentifier
	inCapability5GMM                      nasType.Capability5GMM
	inUESecurityCapability                nasType.UESecurityCapability
	inRequestedNSSAI                      nasType.RequestedNSSAI
	inLastVisitedRegisteredTAI            nasType.LastVisitedRegisteredTAI
	inS1UENetworkCapability               nasType.S1UENetworkCapability
	inUplinkDataStatus                    nasType.UplinkDataStatus
	inPDUSessionStatus                    nasType.PDUSessionStatus
	inMICOIndication                      nasType.MICOIndication
	inUEStatus                            nasType.UEStatus
	inAdditionalGUTI                      nasType.AdditionalGUTI
	inAllowedPDUSessionStatus             nasType.AllowedPDUSessionStatus
	inUesUsageSetting                     nasType.UesUsageSetting
	inRequestedDRXParameters              nasType.RequestedDRXParameters
	inEPSNASMessageContainer              nasType.EPSNASMessageContainer
	inLADNIndication                      nasType.LADNIndication
	inPayloadContainer                    nasType.PayloadContainer
	inNetworkSlicingIndication            nasType.NetworkSlicingIndication
	inUpdateType5GS                       nasType.UpdateType5GS
	inNASMessageContainer                 nasType.NASMessageContainer
}

var nasMessageRegistrationRequestTable = []nasMessageRegistrationRequestData{
	{
		inExtendedProtocolDiscriminator:      nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                     0x01,
		inSpareHalfOctet:                     0x01,
		inRegistrationRequestMessageIdentity: nas.MsgTypeRegistrationRequest,
		inNgksi:                              0x01,
		inRegistrationType5GS:                0x01,
		inMobileIdentity5GS: nasType.MobileIdentity5GS{
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inNoncurrentNativeNASKeySetIdentifier: nasType.NoncurrentNativeNASKeySetIdentifier{
			Octet: 0xC0,
		},
		inCapability5GMM: nasType.Capability5GMM{
			Iei:   nasMessage.RegistrationRequestCapability5GMMType,
			Len:   13,
			Octet: [13]uint8{0x01},
		},
		inUESecurityCapability: nasType.UESecurityCapability{
			Iei:    nasMessage.RegistrationRequestUESecurityCapabilityType,
			Len:    8,
			Buffer: []uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
		},
		inRequestedNSSAI: nasType.RequestedNSSAI{
			Iei:    nasMessage.RegistrationRequestRequestedNSSAIType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inLastVisitedRegisteredTAI: nasType.LastVisitedRegisteredTAI{
			Iei:   nasMessage.RegistrationRequestLastVisitedRegisteredTAIType,
			Octet: [7]uint8{0x01, 0x01},
		},
		inS1UENetworkCapability: nasType.S1UENetworkCapability{
			Iei:    nasMessage.RegistrationRequestS1UENetworkCapabilityType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inUplinkDataStatus: nasType.UplinkDataStatus{
			Iei:    nasMessage.RegistrationRequestUplinkDataStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPDUSessionStatus: nasType.PDUSessionStatus{
			Iei:    nasMessage.RegistrationRequestPDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inMICOIndication: nasType.MICOIndication{
			Octet: 0xB0,
		},
		inUEStatus: nasType.UEStatus{
			Iei:   nasMessage.RegistrationRequestUEStatusType,
			Len:   2,
			Octet: 0x01,
		},
		inAdditionalGUTI: nasType.AdditionalGUTI{
			Iei:   nasMessage.RegistrationRequestAdditionalGUTIType,
			Len:   11,
			Octet: [11]uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
		},
		inAllowedPDUSessionStatus: nasType.AllowedPDUSessionStatus{
			Iei:    nasMessage.RegistrationRequestAllowedPDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inUesUsageSetting: nasType.UesUsageSetting{
			Iei:   nasMessage.RegistrationRequestUesUsageSettingType,
			Len:   2,
			Octet: 0x01,
		},
		inRequestedDRXParameters: nasType.RequestedDRXParameters{
			Iei:   nasMessage.RegistrationRequestRequestedDRXParametersType,
			Len:   2,
			Octet: 0x01,
		},
		inEPSNASMessageContainer: nasType.EPSNASMessageContainer{
			Iei:    nasMessage.RegistrationRequestEPSNASMessageContainerType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inLADNIndication: nasType.LADNIndication{
			Iei:    nasMessage.RegistrationRequestLADNIndicationType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPayloadContainer: nasType.PayloadContainer{
			Iei:    nasMessage.RegistrationRequestPayloadContainerType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inNetworkSlicingIndication: nasType.NetworkSlicingIndication{
			Octet: 0x90,
		},
		inUpdateType5GS: nasType.UpdateType5GS{
			Iei:   nasMessage.RegistrationRequestUpdateType5GSType,
			Len:   2,
			Octet: 0x01,
		},
		inNASMessageContainer: nasType.NASMessageContainer{
			Iei:    nasMessage.RegistrationRequestNASMessageContainerType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewRegistrationRequest(t *testing.T) {}

func TestNasTypeNewRegistrationRequestMessage(t *testing.T) {}
