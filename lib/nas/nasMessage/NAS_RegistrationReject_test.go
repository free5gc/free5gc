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

type nasMessageRegistrationRejectData struct {
	inExtendedProtocolDiscriminator     uint8
	inSecurityHeader                    uint8
	inSpareHalfOctet                    uint8
	inRegistrationRejectMessageIdentity uint8
	inCause5GMM                         nasType.Cause5GMM
	inT3346Value                        nasType.T3346Value
	inT3502Value                        nasType.T3502Value
	inEAPMessage                        nasType.EAPMessage
}

var nasMessageRegistrationRejectTable = []nasMessageRegistrationRejectData{
	{
		inExtendedProtocolDiscriminator:     nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                    0x01,
		inSpareHalfOctet:                    0x01,
		inRegistrationRejectMessageIdentity: nas.MsgTypeRegistrationReject,
		inCause5GMM: nasType.Cause5GMM{
			Octet: 0x01,
		},
		inT3346Value: nasType.T3346Value{
			Iei:   nasMessage.RegistrationRejectT3346ValueType,
			Len:   2,
			Octet: 0x01,
		},
		inT3502Value: nasType.T3502Value{
			Iei:   nasMessage.RegistrationRejectT3502ValueType,
			Len:   2,
			Octet: 0x01,
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.RegistrationRejectEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewRegistrationReject(t *testing.T) {}

func TestNasTypeNewRegistrationRejectMessage(t *testing.T) {}
