//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessagePDUSessionModificationRequestData struct {
	inExtendedProtocolDiscriminator                uint8
	inPDUSessionID                                 uint8
	inPTI                                          uint8
	inPDUSessionModificationRequestMessageIdentity uint8
	inCapability5GSM                               nasType.Capability5GSM
	inCause5GSM                                    nasType.Cause5GSM
	inMaximumNumberOfSupportedPacketFilters        nasType.MaximumNumberOfSupportedPacketFilters
	inAlwaysonPDUSessionRequested                  nasType.AlwaysonPDUSessionRequested
	inIntegrityProtectionMaximumDataRate           nasType.IntegrityProtectionMaximumDataRate
	inRequestedQosRules                            nasType.RequestedQosRules
	inRequestedQosFlowDescriptions                 nasType.RequestedQosFlowDescriptions
	inMappedEPSBearerContexts                      nasType.MappedEPSBearerContexts
	inExtendedProtocolConfigurationOptions         nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionModificationRequestTable = []nasMessagePDUSessionModificationRequestData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSessionModificationRequestMessageIdentity: 0x01,
		inCapability5GSM: nasType.Capability5GSM{
			Iei:   nasMessage.PDUSessionModificationRequestCapability5GSMType,
			Len:   2,
			Octet: [13]uint8{0x01, 0x01},
		},
		inCause5GSM: nasType.Cause5GSM{
			Iei:   nasMessage.PDUSessionModificationRequestCause5GSMType,
			Octet: 0x01,
		},
		inMaximumNumberOfSupportedPacketFilters: nasType.MaximumNumberOfSupportedPacketFilters{
			Iei:   nasMessage.PDUSessionModificationRequestMaximumNumberOfSupportedPacketFiltersType,
			Octet: [2]uint8{0x01, 0x01},
		},
		inAlwaysonPDUSessionRequested: nasType.AlwaysonPDUSessionRequested{
			Octet: 0xB0,
		},
		inIntegrityProtectionMaximumDataRate: nasType.IntegrityProtectionMaximumDataRate{
			Iei:   nasMessage.PDUSessionModificationRequestIntegrityProtectionMaximumDataRateType,
			Octet: [2]uint8{0x01, 0x01},
		},
		inRequestedQosRules: nasType.RequestedQosRules{
			Iei:    nasMessage.PDUSessionModificationRequestRequestedQosRulesType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inRequestedQosFlowDescriptions: nasType.RequestedQosFlowDescriptions{
			Iei:    nasMessage.PDUSessionModificationRequestRequestedQosFlowDescriptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inMappedEPSBearerContexts: nasType.MappedEPSBearerContexts{
			Iei:    nasMessage.PDUSessionModificationRequestMappedEPSBearerContextsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionModificationRequestExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionModificationRequest(t *testing.T) {}

func TestNasTypeNewPDUSessionModificationRequestMessage(t *testing.T) {}
