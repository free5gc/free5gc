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

type nasMessageSecurityProtected5GSNASMessageData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeader                uint8
	inSpareHalfOctet                uint8
	inMessageAuthenticationCode     nasType.MessageAuthenticationCode
	inSequenceNumber                nasType.SequenceNumber
	inPlain5GSNASMessage            nasType.Plain5GSNASMessage
}

var nasMessageSecurityProtected5GSNASMessageTable = []nasMessageSecurityProtected5GSNASMessageData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                0x01,
		inSpareHalfOctet:                0x01,
		inMessageAuthenticationCode: nasType.MessageAuthenticationCode{
			Octet: [4]uint8{0x01, 0x01, 0x01, 0x01},
		},
		inSequenceNumber: nasType.SequenceNumber{
			Octet: 0x01,
		},
		inPlain5GSNASMessage: nasType.Plain5GSNASMessage{},
	},
}

func TestNasTypeNewSecurityProtected5GSNASMessage(t *testing.T) {}

func TestNasTypeNewSecurityProtected5GSNASMessageMessage(t *testing.T) {}
