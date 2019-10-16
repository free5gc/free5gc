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

type nasMessagePDUSessionEstablishmentRejectData struct {
	inExtendedProtocolDiscriminator                uint8
	inPDUSessionID                                 uint8
	inPTI                                          uint8
	inPDUSESSIONESTABLISHMENTREJECTMessageIdentity uint8
	inCause5GSM                                    nasType.Cause5GSM
	inBackoffTimerValue                            nasType.BackoffTimerValue
	inAllowedSSCMode                               nasType.AllowedSSCMode
	inEAPMessage                                   nasType.EAPMessage
	inExtendedProtocolConfigurationOptions         nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionEstablishmentRejectTable = []nasMessagePDUSessionEstablishmentRejectData{
	{
		inExtendedProtocolDiscriminator: nas.MsgTypePDUSessionEstablishmentReject,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONESTABLISHMENTREJECTMessageIdentity: 0x01,
		inCause5GSM: nasType.Cause5GSM{
			Iei:   0,
			Octet: 0x01,
		},
		inBackoffTimerValue: nasType.BackoffTimerValue{
			Iei:   nasMessage.PDUSessionEstablishmentRejectBackoffTimerValueType,
			Len:   2,
			Octet: 0x01,
		},
		inAllowedSSCMode: nasType.AllowedSSCMode{
			Octet: 0xF0,
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.PDUSessionEstablishmentRejectEAPMessageType,
			Len:    1,
			Buffer: []uint8{0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionEstablishmentRejectExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionEstablishmentReject(t *testing.T) {}

func TestNasTypeNewPDUSessionEstablishmentRejectMessage(t *testing.T) {}
