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

type nasMessagePDUSessionModificationCommandData struct {
	inExtendedProtocolDiscriminator                uint8
	inPDUSessionID                                 uint8
	inPTI                                          uint8
	inPDUSESSIONMODIFICATIONCOMMANDMessageIdentity uint8
	inCause5GSM                                    nasType.Cause5GSM
	inSessionAMBR                                  nasType.SessionAMBR
	inRQTimerValue                                 nasType.RQTimerValue
	inAlwaysonPDUSessionIndication                 nasType.AlwaysonPDUSessionIndication
	inAuthorizedQosRules                           nasType.AuthorizedQosRules
	inMappedEPSBearerContexts                      nasType.MappedEPSBearerContexts
	inAuthorizedQosFlowDescriptions                nasType.AuthorizedQosFlowDescriptions
	inExtendedProtocolConfigurationOptions         nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionModificationCommandTable = []nasMessagePDUSessionModificationCommandData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONMODIFICATIONCOMMANDMessageIdentity: 0x01,
		inCause5GSM: nasType.Cause5GSM{
			Iei:   nasMessage.PDUSessionModificationCommandCause5GSMType,
			Octet: 0x01,
		},
		inSessionAMBR: nasType.SessionAMBR{
			Iei:   nasMessage.PDUSessionModificationCommandSessionAMBRType,
			Len:   6,
			Octet: [6]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		},
		inRQTimerValue: nasType.RQTimerValue{
			Iei:   nasMessage.PDUSessionModificationCommandRQTimerValueType,
			Octet: 0x01,
		},
		inAlwaysonPDUSessionIndication: nasType.AlwaysonPDUSessionIndication{
			Octet: 0x80,
		},
		inAuthorizedQosRules: nasType.AuthorizedQosRules{
			Iei:    nasMessage.PDUSessionModificationCommandAuthorizedQosRulesType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inMappedEPSBearerContexts: nasType.MappedEPSBearerContexts{
			Iei:    nasMessage.PDUSessionModificationCommandMappedEPSBearerContextsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inAuthorizedQosFlowDescriptions: nasType.AuthorizedQosFlowDescriptions{
			Iei:    nasMessage.PDUSessionModificationCommandAuthorizedQosFlowDescriptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionModificationCommandExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionModificationCommand(t *testing.T) {}

func TestNasTypeNewPDUSessionModificationCommandMessage(t *testing.T) {}
