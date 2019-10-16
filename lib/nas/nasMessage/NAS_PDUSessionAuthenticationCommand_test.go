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

type nasMessagePDUSessionAuthenticationCommandData struct {
	inExtendedProtocolDiscriminator                  uint8
	inPDUSessionID                                   uint8
	inPTI                                            uint8
	inPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity uint8
	inEAPMessage                                     nasType.EAPMessage
	inExtendedProtocolConfigurationOptions           nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionAuthenticationCommandTable = []nasMessagePDUSessionAuthenticationCommandData{
	{
		inExtendedProtocolDiscriminator: nas.MsgTypePDUSessionAuthenticationCommand,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity: 0x01,
		inEAPMessage: nasType.EAPMessage{
			Iei:    0,
			Len:    4,
			Buffer: []uint8{0x01, 0x01, 0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionAuthenticationCommandExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionAuthenticationCommand(t *testing.T) {}

func TestNasTypeNewPDUSessionAuthenticationCommandMessage(t *testing.T) {}
