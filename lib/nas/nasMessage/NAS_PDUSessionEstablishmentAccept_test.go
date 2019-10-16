//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessagePDUSessionEstablishmentAcceptData struct {
	inExtendedProtocolDiscriminator                uint8
	inPDUSessionID                                 uint8
	inPTI                                          uint8
	inPDUSESSIONESTABLISHMENTACCEPTMessageIdentity uint8
	inSelectedSSCModeAndSelectedPDUSessionType     nasType.SelectedSSCModeAndSelectedPDUSessionType
	inAuthorizedQosRules                           nasType.AuthorizedQosRules
	inSessionAMBR                                  nasType.SessionAMBR
	inCause5GSM                                    nasType.Cause5GSM
	inPDUAddress                                   nasType.PDUAddress
	inRQTimerValue                                 nasType.RQTimerValue
	inSNSSAI                                       nasType.SNSSAI
	inAlwaysonPDUSessionIndication                 nasType.AlwaysonPDUSessionIndication
	inMappedEPSBearerContexts                      nasType.MappedEPSBearerContexts
	inEAPMessage                                   nasType.EAPMessage
	inAuthorizedQosFlowDescriptions                nasType.AuthorizedQosFlowDescriptions
	inExtendedProtocolConfigurationOptions         nasType.ExtendedProtocolConfigurationOptions
	inDNN                                          nasType.DNN
}

var nasMessagePDUSessionEstablishmentAcceptTable = []nasMessagePDUSessionEstablishmentAcceptData{
	{
		inExtendedProtocolDiscriminator: nas.MsgTypePDUSessionEstablishmentAccept,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONESTABLISHMENTACCEPTMessageIdentity: 0x01,
		inSelectedSSCModeAndSelectedPDUSessionType: nasType.SelectedSSCModeAndSelectedPDUSessionType{
			Octet: 0x01,
		},
		inAuthorizedQosRules: nasType.AuthorizedQosRules{
			Iei:    0,
			Len:    1,
			Buffer: []uint8{0x01},
		},
		inSessionAMBR: nasType.SessionAMBR{
			Iei:   0,
			Len:   6,
			Octet: [6]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		},
		inCause5GSM: nasType.Cause5GSM{
			Iei:   nasMessage.PDUSessionEstablishmentAcceptCause5GSMType,
			Octet: 0x01,
		},
		inPDUAddress: nasType.PDUAddress{
			Iei:   nasMessage.PDUSessionEstablishmentAcceptPDUAddressType,
			Len:   13,
			Octet: [13]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C},
		},
		inRQTimerValue: nasType.RQTimerValue{
			Iei:   nasMessage.PDUSessionEstablishmentAcceptRQTimerValueType,
			Octet: 0x01,
		},
		inSNSSAI: nasType.SNSSAI{
			Iei:   nasMessage.PDUSessionEstablishmentAcceptSNSSAIType,
			Len:   2,
			Octet: [8]uint8{0x01, 0x01},
		},
		inAlwaysonPDUSessionIndication: nasType.AlwaysonPDUSessionIndication{
			Octet: 0x80,
		},
		inMappedEPSBearerContexts: nasType.MappedEPSBearerContexts{
			Iei:    nasMessage.PDUSessionEstablishmentAcceptMappedEPSBearerContextsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.PDUSessionEstablishmentAcceptEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inAuthorizedQosFlowDescriptions: nasType.AuthorizedQosFlowDescriptions{
			Iei:    nasMessage.PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionEstablishmentAcceptExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inDNN: nasType.DNN{
			Iei:    nasMessage.ULNASTransportDNNType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionEstablishmentAccept(t *testing.T) {}

func TestNasTypeNewPDUSessionEstablishmentAcceptMessage(t *testing.T) {}
