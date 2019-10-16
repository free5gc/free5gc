//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas/logger"

	//"fmt"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

type nasMessageIdentityResponseData struct {
	inExtendedProtocolDiscriminator   uint8
	inSecurityHeader                  uint8
	inSpareHalfOctet                  uint8
	inIdentityResponseMessageIdentity uint8
	inMobileIdentity                  nasType.MobileIdentity
}

var nasMessageIdentityResponseTable = []nasMessageIdentityResponseData{
	{
		inExtendedProtocolDiscriminator:   0x01,
		inSecurityHeader:                  0x08,
		inSpareHalfOctet:                  0x01,
		inIdentityResponseMessageIdentity: nas.MsgTypeIdentityResponse,
		inMobileIdentity: nasType.MobileIdentity{
			Iei:    0,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewIdentityResponse(t *testing.T) {}

func TestNasTypeNewIdentityResponseMessage(t *testing.T) {}
