//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessageConfigurationUpdateCompleteData struct {
	inExtendedProtocolDiscriminator              uint8
	inSecurityHeaderType                         uint8
	inSpareHalfOctet                             uint8
	inConfigurationUpdateCompleteMessageIdentity uint8
}

var nasMessageConfigurationUpdateCompleteTable = []nasMessageConfigurationUpdateCompleteData{
	{
		inExtendedProtocolDiscriminator:              nasMessage.Epd5GSSessionManagementMessage,
		inSecurityHeaderType:                         0x01,
		inSpareHalfOctet:                             0x01,
		inConfigurationUpdateCompleteMessageIdentity: nas.MsgTypeConfigurationUpdateComplete,
	},
}

func TestNasTypeNewConfigurationUpdateComplete(t *testing.T) {}

func TestNasTypeNewConfigurationUpdateCompleteMessage(t *testing.T) {}
