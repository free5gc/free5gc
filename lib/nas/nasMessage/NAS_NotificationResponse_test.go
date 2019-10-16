//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas/logger"

	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

type nasMessageNotificationResponseData struct {
	inExtendedProtocolDiscriminator       uint8
	inSecurityHeader                      uint8
	inSpareHalfOctet                      uint8
	inNotificationResponseMessageIdentity uint8
	inPDUSessionStatus                    nasType.PDUSessionStatus
}

var nasMessageNotificationResponseTable = []nasMessageNotificationResponseData{
	{
		inExtendedProtocolDiscriminator:       0x01,
		inSecurityHeader:                      0x08,
		inSpareHalfOctet:                      0x01,
		inNotificationResponseMessageIdentity: 0x01,
		inPDUSessionStatus: nasType.PDUSessionStatus{
			Iei:    nasMessage.NotificationResponsePDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewNotificationResponse(t *testing.T) {}

func TestNasTypeNewNotificationResponseMessage(t *testing.T) {}
