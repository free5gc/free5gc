//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type Notification struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.NotificationMessageIdentity
	nasType.SpareHalfOctetAndAccessType
}

func NewNotification(iei uint8) (notification *Notification) {}

func (a *Notification) EncodeNotification(buffer *bytes.Buffer) {}

func (a *Notification) DecodeNotification(byteArray *[]byte) {}
