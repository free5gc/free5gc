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

type nasMessagePDUSessionEstablishmentRequestData struct {
	inExtendedProtocolDiscriminator                 uint8
	inPDUSessionID                                  uint8
	inPTI                                           uint8
	inPDUSESSIONESTABLISHMENTREQUESTMessageIdentity uint8
	inIntegrityProtectionMaximumDataRate            nasType.IntegrityProtectionMaximumDataRate
	inPDUSessionType                                nasType.PDUSessionType
	inSSCMode                                       nasType.SSCMode
	inCapability5GSM                                nasType.Capability5GSM
	inMaximumNumberOfSupportedPacketFilters         nasType.MaximumNumberOfSupportedPacketFilters
	inAlwaysonPDUSessionRequested                   nasType.AlwaysonPDUSessionRequested
	inSMPDUDNRequestContainer                       nasType.SMPDUDNRequestContainer
	inExtendedProtocolConfigurationOptions          nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionEstablishmentRequestTable = []nasMessagePDUSessionEstablishmentRequestData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONESTABLISHMENTREQUESTMessageIdentity: 0x01,
		inIntegrityProtectionMaximumDataRate: nasType.IntegrityProtectionMaximumDataRate{
			Iei:   0,
			Octet: [2]uint8{0x01, 0x01},
		},
		inPDUSessionType: nasType.PDUSessionType{
			Octet: 0x90,
		},
		inSSCMode: nasType.SSCMode{
			Octet: 0xA0,
		},
		inCapability5GSM: nasType.Capability5GSM{
			Iei:   nasMessage.PDUSessionEstablishmentRequestCapability5GSMType,
			Len:   2,
			Octet: [13]uint8{0x01, 0x01},
		},
		inMaximumNumberOfSupportedPacketFilters: nasType.MaximumNumberOfSupportedPacketFilters{
			Iei:   nasMessage.PDUSessionEstablishmentRequestMaximumNumberOfSupportedPacketFiltersType,
			Octet: [2]uint8{0x01, 0x01},
		},
		inAlwaysonPDUSessionRequested: nasType.AlwaysonPDUSessionRequested{
			Octet: 0xB0,
		},
		inSMPDUDNRequestContainer: nasType.SMPDUDNRequestContainer{
			Iei:    nasMessage.PDUSessionEstablishmentRequestSMPDUDNRequestContainerType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionEstablishmentRequestExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionEstablishmentRequest(t *testing.T) {}

func TestNasTypeNewPDUSessionEstablishmentRequestMessage(t *testing.T) {}
