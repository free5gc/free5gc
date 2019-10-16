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

type nasMessagePDUSessionAuthenticationResultData struct {
	inExtendedProtocolDiscriminator                 uint8
	inPDUSessionID                                  uint8
	inPTI                                           uint8
	inPDUSESSIONAUTHENTICATIONRESULTMessageIdentity uint8
	inEAPMessage                                    nasType.EAPMessage
	inExtendedProtocolConfigurationOptions          nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionAuthenticationResultTable = []nasMessagePDUSessionAuthenticationResultData{
	{
		inExtendedProtocolDiscriminator: nas.MsgTypePDUSessionAuthenticationResult,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONAUTHENTICATIONRESULTMessageIdentity: 0x01,
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.PDUSessionAuthenticationResultEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionAuthenticationResultExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionAuthenticationResult(t *testing.T) {}

func TestNasTypeNewPDUSessionAuthenticationResultMessage(t *testing.T) {}
