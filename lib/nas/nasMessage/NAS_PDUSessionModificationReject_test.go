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

type nasMessagePDUSessionModificationRejectData struct {
	inExtendedProtocolDiscriminator               uint8
	inPDUSessionID                                uint8
	inPTI                                         uint8
	inPDUSESSIONMODIFICATIONREJECTMessageIdentity uint8
	inCause5GSM                                   nasType.Cause5GSM
	inBackoffTimerValue                           nasType.BackoffTimerValue
	inExtendedProtocolConfigurationOptions        nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionModificationRejectTable = []nasMessagePDUSessionModificationRejectData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONMODIFICATIONREJECTMessageIdentity: 0x01,
		inCause5GSM: nasType.Cause5GSM{
			Iei:   0,
			Octet: 0x01,
		},
		inBackoffTimerValue: nasType.BackoffTimerValue{
			Iei:   nasMessage.PDUSessionModificationRejectBackoffTimerValueType,
			Len:   2,
			Octet: 0x01,
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionModificationRejectExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionModificationReject(t *testing.T) {}

func TestNasTypeNewPDUSessionModificationRejectMessage(t *testing.T) {}
