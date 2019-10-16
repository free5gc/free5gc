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

type nasMessageSecurityModeRejectData struct {
	inExtendedProtocolDiscriminator     uint8
	inSecurityHeader                    uint8
	inSpareHalfOctet                    uint8
	inSecurityModeRejectMessageIdentity uint8
	inCause5GMM                         nasType.Cause5GMM
}

var nasMessageSecurityModeRejectTable = []nasMessageSecurityModeRejectData{
	{
		inExtendedProtocolDiscriminator:     nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                    0x01,
		inSpareHalfOctet:                    0x01,
		inSecurityModeRejectMessageIdentity: nas.MsgTypeSecurityModeReject,
		inCause5GMM: nasType.Cause5GMM{
			Octet: 0x01,
		},
	},
}

func TestNasTypeNewSecurityModeReject(t *testing.T) {}

func TestNasTypeNewSecurityModeRejectMessage(t *testing.T) {}
