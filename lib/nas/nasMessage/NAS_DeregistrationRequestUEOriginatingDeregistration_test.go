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

type nasMessageDeregistrationRequestUEOriginatingDeregistrationData struct {
	inExtendedProtocolDiscriminator        uint8
	inSecurityHeaderType                   uint8
	inSpareHalfOctet                       uint8
	inDeregistrationRequestMessageIdentity uint8
	inNgksiAndDeregistrationType           nasType.NgksiAndDeregistrationType
	inMobileIdentity5GS                    nasType.MobileIdentity5GS
}

var nasMessageDeregistrationRequestUEOriginatingDeregistrationTable = []nasMessageDeregistrationRequestUEOriginatingDeregistrationData{
	{
		inExtendedProtocolDiscriminator:        nasMessage.Epd5GSSessionManagementMessage,
		inSecurityHeaderType:                   0x01,
		inSpareHalfOctet:                       0x01,
		inDeregistrationRequestMessageIdentity: 0x01,
		inNgksiAndDeregistrationType: nasType.NgksiAndDeregistrationType{
			Octet: 0xFF,
		},
		inMobileIdentity5GS: nasType.MobileIdentity5GS{
			Iei:    0,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewDeregistrationRequestUEOriginatingDeregistration(t *testing.T) {}

func TestNasTypeNewDeregistrationRequestUEOriginatingDeregistrationMessage(t *testing.T) {}
