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

type nasMessagePDUSessionReleaseRejectData struct {
	inExtendedProtocolDiscriminator          uint8
	inPDUSessionID                           uint8
	inPTI                                    uint8
	inPDUSESSIONRELEASEREJECTMessageIdentity uint8
	inCause5GSM                              nasType.Cause5GSM
	inExtendedProtocolConfigurationOptions   nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionReleaseRejectTable = []nasMessagePDUSessionReleaseRejectData{
	{
		inExtendedProtocolDiscriminator:          nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID:                           0x01,
		inPTI:                                    0x01,
		inPDUSESSIONRELEASEREJECTMessageIdentity: 0x01,
		inCause5GSM: nasType.Cause5GSM{
			Octet: 0x01,
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionReleaseRejectExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionReleaseReject(t *testing.T) {}

func TestNasTypeNewPDUSessionReleaseRejectMessage(t *testing.T) {}
