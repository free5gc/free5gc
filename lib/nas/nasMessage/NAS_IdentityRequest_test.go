//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas/logger"

	//"fmt"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

type nasMessageIdentityRequestData struct {
	inExtendedProtocolDiscriminator  uint8
	inSecurityHeader                 uint8
	inSpareHalfOctet1                uint8
	inIdentityRequestMessageIdentity uint8
	inIdentityType                   uint8
	inSpareHalfOctet2                uint8
}

var nasMessageIdentityRequestTable = []nasMessageIdentityRequestData{
	{
		inExtendedProtocolDiscriminator:  0x01,
		inSecurityHeader:                 0x08,
		inSpareHalfOctet1:                0x01,
		inIdentityRequestMessageIdentity: nas.MsgTypeIdentityRequest,
		inIdentityType:                   0x01,
		inSpareHalfOctet2:                0x01,
	},
}

func TestNasTypeNewIdentityRequest(t *testing.T) {}

func TestNasTypeNewIdentityRequestMessage(t *testing.T) {}
