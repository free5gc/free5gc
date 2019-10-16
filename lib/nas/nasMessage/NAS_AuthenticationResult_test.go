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

type nasMessageAuthenticationResultData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeaderType            uint8
	inMessageType                   uint8
	inTsc                           uint8
	inNASKeySetIdentifier           uint8
	inEAPLen                        uint16
	inEAPMessage                    []uint8
	inABBA                          nasType.ABBA
}

var aBBATestData = []nasType.ABBA{
	{Iei: nasMessage.AuthenticationResultABBAType, Len: 2, Buffer: []byte{0x00, 0x00}},
	//{Iei: 0x81, Len: 2, Buffer: []byte{0x00, 0x00}},
}

var nasMessageAuthenticationResultTable = []nasMessageAuthenticationResultData{
	{inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inSecurityHeaderType:  0x01,
		inMessageType:         nas.MsgTypeAuthenticationResult,
		inTsc:                 0x01,
		inNASKeySetIdentifier: 0x01,
		inEAPLen:              0x02,
		inEAPMessage:          []uint8{0x10, 0x11},
		inABBA:                aBBATestData[0]},
	/*{inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
	inSecurityHeaderType:  0x01,
	inMessageType:         nas.MsgTypeAuthenticationResult,
	inTsc:                 0x01,
	inNASKeySetIdentifier: 0x01,
	inEAPLen:              0x02,
	inEAPMessage:          []uint8{0x10, 0x11},
	inABBA:                aBBATestData[1]},*/
}

func TestNasTypeNewAuthenticationResult(t *testing.T) {}

func TestNasTypeNewAuthenticationResultMessage(t *testing.T) {}
