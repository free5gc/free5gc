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

type nasMessageNotificationData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeader                uint8
	inSpareHalfOctet1               uint8
	inNotificationMessageIdentity   uint8
	inAccessType                    uint8
	inSpareHalfOctet2               uint8
}

var nasMessageNotificationTable = []nasMessageNotificationData{
	{
		inExtendedProtocolDiscriminator: 0x01,
		inSecurityHeader:                0x08,
		inSpareHalfOctet1:               0x01,
		inNotificationMessageIdentity:   nas.MsgTypeNotification,
		inAccessType:                    0x01,
		inSpareHalfOctet2:               0x01,
	},
}

func TestNasTypeNewNotification(t *testing.T) {}

func TestNasTypeNewNotificationMessage(t *testing.T) {}
