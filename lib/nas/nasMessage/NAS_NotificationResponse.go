//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type NotificationResponse struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.NotificationResponseMessageIdentity
	*nasType.PDUSessionStatus
}

func NewNotificationResponse(iei uint8) (notificationResponse *NotificationResponse) {}

const (
	NotificationResponsePDUSessionStatusType uint8 = 0x50
)

func (a *NotificationResponse) EncodeNotificationResponse(buffer *bytes.Buffer) {}

func (a *NotificationResponse) DecodeNotificationResponse(byteArray *[]byte) {}
