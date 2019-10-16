//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type RegistrationAccept struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.RegistrationAcceptMessageIdentity
	nasType.RegistrationResult5GS
	*nasType.GUTI5G
	*nasType.EquivalentPlmns
	*nasType.TAIList
	*nasType.AllowedNSSAI
	*nasType.RejectedNSSAI
	*nasType.ConfiguredNSSAI
	*nasType.NetworkFeatureSupport5GS
	*nasType.PDUSessionStatus
	*nasType.PDUSessionReactivationResult
	*nasType.PDUSessionReactivationResultErrorCause
	*nasType.LADNInformation
	*nasType.MICOIndication
	*nasType.NetworkSlicingIndication
	*nasType.ServiceAreaList
	*nasType.T3512Value
	*nasType.Non3GppDeregistrationTimerValue
	*nasType.T3502Value
	*nasType.EmergencyNumberList
	*nasType.ExtendedEmergencyNumberList
	*nasType.SORTransparentContainer
	*nasType.EAPMessage
	*nasType.NSSAIInclusionMode
	*nasType.OperatordefinedAccessCategoryDefinitions
	*nasType.NegotiatedDRXParameters
}

func NewRegistrationAccept(iei uint8) (registrationAccept *RegistrationAccept) {}

const (
	RegistrationAcceptGUTI5GType                                   uint8 = 0x77
	RegistrationAcceptEquivalentPlmnsType                          uint8 = 0x4A
	RegistrationAcceptTAIListType                                  uint8 = 0x54
	RegistrationAcceptAllowedNSSAIType                             uint8 = 0x15
	RegistrationAcceptRejectedNSSAIType                            uint8 = 0x11
	RegistrationAcceptConfiguredNSSAIType                          uint8 = 0x31
	RegistrationAcceptNetworkFeatureSupport5GSType                 uint8 = 0x21
	RegistrationAcceptPDUSessionStatusType                         uint8 = 0x50
	RegistrationAcceptPDUSessionReactivationResultType             uint8 = 0x26
	RegistrationAcceptPDUSessionReactivationResultErrorCauseType   uint8 = 0x72
	RegistrationAcceptLADNInformationType                          uint8 = 0x79
	RegistrationAcceptMICOIndicationType                           uint8 = 0x0B
	RegistrationAcceptNetworkSlicingIndicationType                 uint8 = 0x09
	RegistrationAcceptServiceAreaListType                          uint8 = 0x27
	RegistrationAcceptT3512ValueType                               uint8 = 0x5E
	RegistrationAcceptNon3GppDeregistrationTimerValueType          uint8 = 0x5D
	RegistrationAcceptT3502ValueType                               uint8 = 0x16
	RegistrationAcceptEmergencyNumberListType                      uint8 = 0x34
	RegistrationAcceptExtendedEmergencyNumberListType              uint8 = 0x7A
	RegistrationAcceptSORTransparentContainerType                  uint8 = 0x73
	RegistrationAcceptEAPMessageType                               uint8 = 0x78
	RegistrationAcceptNSSAIInclusionModeType                       uint8 = 0x0A
	RegistrationAcceptOperatordefinedAccessCategoryDefinitionsType uint8 = 0x76
	RegistrationAcceptNegotiatedDRXParametersType                  uint8 = 0x51
)

func (a *RegistrationAccept) EncodeRegistrationAccept(buffer *bytes.Buffer) {}

func (a *RegistrationAccept) DecodeRegistrationAccept(byteArray *[]byte) {}
