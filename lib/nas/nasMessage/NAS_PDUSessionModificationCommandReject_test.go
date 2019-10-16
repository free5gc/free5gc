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

type nasMessagePDUSessionModificationCommandRejectData struct {
	inExtendedProtocolDiscriminator                      uint8
	inPDUSessionID                                       uint8
	inPTI                                                uint8
	inPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity uint8
	inCause5GSM                                          nasType.Cause5GSM
	inExtendedProtocolConfigurationOptions               nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionModificationCommandRejectTable = []nasMessagePDUSessionModificationCommandRejectData{
	{
		inExtendedProtocolDiscriminator: nas.MsgTypePDUSessionModificationCommandReject,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity: 0x01,
		inCause5GSM: nasType.Cause5GSM{
			Iei:   0,
			Octet: 0x01,
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionModificationCommandRejectExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionModificationCommandReject(t *testing.T) {}

func TestNasTypeNewPDUSessionModificationCommandRejectMessage(t *testing.T) {}
