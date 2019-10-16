//go:binary-only-package

package nasType

// NotificationMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type NotificationMessageIdentity struct {
	Octet uint8
}

func NewNotificationMessageIdentity() (notificationMessageIdentity *NotificationMessageIdentity) {}

// NotificationMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *NotificationMessageIdentity) GetMessageType() (messageType uint8) {}

// NotificationMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *NotificationMessageIdentity) SetMessageType(messageType uint8) {}
