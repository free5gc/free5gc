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

type nasMessageDeregistrationRequestUETerminatedDeregistrationData struct {
	inExtendedProtocolDiscriminator        uint8
	inSecurityHeaderType                   uint8
	inSpareHalfOctet1                      uint8
	inDeregistrationRequestMessageIdentity uint8
	inSpareHalfOctetAndDeregistrationType  nasType.SpareHalfOctetAndDeregistrationType
	inCause5GMM                            nasType.Cause5GMM
	inT3346Value                           nasType.T3346Value
}

var nasMessageDeregistrationRequestUETerminatedDeregistrationTable = []nasMessageDeregistrationRequestUETerminatedDeregistrationData{
	{
		inExtendedProtocolDiscriminator:        nasMessage.Epd5GSSessionManagementMessage,
		inSecurityHeaderType:                   0x01,
		inSpareHalfOctet1:                      0x01,
		inDeregistrationRequestMessageIdentity: 0x01,
		inSpareHalfOctetAndDeregistrationType: nasType.SpareHalfOctetAndDeregistrationType{
			Octet: 0x0F,
		},
		inCause5GMM: nasType.Cause5GMM{
			Iei:   nasMessage.DeregistrationRequestUETerminatedDeregistrationCause5GMMType,
			Octet: 0x01,
		},
		inT3346Value: nasType.T3346Value{
			Iei:   nasMessage.DeregistrationRequestUETerminatedDeregistrationT3346ValueType,
			Len:   2,
			Octet: 0x01,
		},
	},
}

func TestNasTypeNewDeregistrationRequestUETerminatedDeregistration(t *testing.T) {}

func TestNasTypeNewDeregistrationRequestUETerminatedDeregistrationMessage(t *testing.T) {}
