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

type nasMessagePDUSessionReleaseCompleteData struct {
	inExtendedProtocolDiscriminator            uint8
	inPDUSessionID                             uint8
	inPTI                                      uint8
	inPDUSESSIONRELEASECOMPLETEMessageIdentity uint8
	inCause5GSM                                nasType.Cause5GSM
	inExtendedProtocolConfigurationOptions     nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionReleaseCompleteTable = []nasMessagePDUSessionReleaseCompleteData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONRELEASECOMPLETEMessageIdentity: 0x01,
		inCause5GSM: nasType.Cause5GSM{
			Iei:   nasMessage.PDUSessionReleaseCompleteCause5GSMType,
			Octet: 0x01,
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionReleaseCompleteExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionReleaseComplete(t *testing.T) {}

func TestNasTypeNewPDUSessionReleaseCompleteMessage(t *testing.T) {}
