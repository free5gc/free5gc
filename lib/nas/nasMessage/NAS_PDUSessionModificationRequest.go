//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionModificationRequest struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONMODIFICATIONREQUESTMessageIdentity
	*nasType.Capability5GSM
	*nasType.Cause5GSM
	*nasType.MaximumNumberOfSupportedPacketFilters
	*nasType.AlwaysonPDUSessionRequested
	*nasType.IntegrityProtectionMaximumDataRate
	*nasType.RequestedQosRules
	*nasType.RequestedQosFlowDescriptions
	*nasType.MappedEPSBearerContexts
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionModificationRequest(iei uint8) (pDUSessionModificationRequest *PDUSessionModificationRequest) {}

const (
	PDUSessionModificationRequestCapability5GSMType                        uint8 = 0x28
	PDUSessionModificationRequestCause5GSMType                             uint8 = 0x59
	PDUSessionModificationRequestMaximumNumberOfSupportedPacketFiltersType uint8 = 0x55
	PDUSessionModificationRequestAlwaysonPDUSessionRequestedType           uint8 = 0x0B
	PDUSessionModificationRequestIntegrityProtectionMaximumDataRateType    uint8 = 0x13
	PDUSessionModificationRequestRequestedQosRulesType                     uint8 = 0x7A
	PDUSessionModificationRequestRequestedQosFlowDescriptionsType          uint8 = 0x79
	PDUSessionModificationRequestMappedEPSBearerContextsType               uint8 = 0x7F
	PDUSessionModificationRequestExtendedProtocolConfigurationOptionsType  uint8 = 0x7B
)

func (a *PDUSessionModificationRequest) EncodePDUSessionModificationRequest(buffer *bytes.Buffer) {}

func (a *PDUSessionModificationRequest) DecodePDUSessionModificationRequest(byteArray *[]byte) {}
