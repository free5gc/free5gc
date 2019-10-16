//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type RegistrationRequest struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.RegistrationRequestMessageIdentity
	nasType.NgksiAndRegistrationType5GS
	nasType.MobileIdentity5GS
	*nasType.NoncurrentNativeNASKeySetIdentifier
	*nasType.Capability5GMM
	*nasType.UESecurityCapability
	*nasType.RequestedNSSAI
	*nasType.LastVisitedRegisteredTAI
	*nasType.S1UENetworkCapability
	*nasType.UplinkDataStatus
	*nasType.PDUSessionStatus
	*nasType.MICOIndication
	*nasType.UEStatus
	*nasType.AdditionalGUTI
	*nasType.AllowedPDUSessionStatus
	*nasType.UesUsageSetting
	*nasType.RequestedDRXParameters
	*nasType.EPSNASMessageContainer
	*nasType.LADNIndication
	*nasType.PayloadContainer
	*nasType.NetworkSlicingIndication
	*nasType.UpdateType5GS
	*nasType.NASMessageContainer
}

func NewRegistrationRequest(iei uint8) (registrationRequest *RegistrationRequest) {}

const (
	RegistrationRequestNoncurrentNativeNASKeySetIdentifierType uint8 = 0x0C
	RegistrationRequestCapability5GMMType                      uint8 = 0x10
	RegistrationRequestUESecurityCapabilityType                uint8 = 0x2E
	RegistrationRequestRequestedNSSAIType                      uint8 = 0x2F
	RegistrationRequestLastVisitedRegisteredTAIType            uint8 = 0x52
	RegistrationRequestS1UENetworkCapabilityType               uint8 = 0x17
	RegistrationRequestUplinkDataStatusType                    uint8 = 0x40
	RegistrationRequestPDUSessionStatusType                    uint8 = 0x50
	RegistrationRequestMICOIndicationType                      uint8 = 0x0B
	RegistrationRequestUEStatusType                            uint8 = 0x2B
	RegistrationRequestAdditionalGUTIType                      uint8 = 0x77
	RegistrationRequestAllowedPDUSessionStatusType             uint8 = 0x25
	RegistrationRequestUesUsageSettingType                     uint8 = 0x18
	RegistrationRequestRequestedDRXParametersType              uint8 = 0x51
	RegistrationRequestEPSNASMessageContainerType              uint8 = 0x70
	RegistrationRequestLADNIndicationType                      uint8 = 0x74
	RegistrationRequestPayloadContainerType                    uint8 = 0x7B
	RegistrationRequestNetworkSlicingIndicationType            uint8 = 0x09
	RegistrationRequestUpdateType5GSType                       uint8 = 0x53
	RegistrationRequestNASMessageContainerType                 uint8 = 0x71
)

func (a *RegistrationRequest) EncodeRegistrationRequest(buffer *bytes.Buffer) {}

func (a *RegistrationRequest) DecodeRegistrationRequest(byteArray *[]byte) {}
