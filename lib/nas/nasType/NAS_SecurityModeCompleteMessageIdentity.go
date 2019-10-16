//go:binary-only-package

package nasType

// SecurityModeCompleteMessageIdentity 9.6
// MessageType Row, sBit, len = [0, 0], 8 , 8
type SecurityModeCompleteMessageIdentity struct {
	Octet uint8
}

func NewSecurityModeCompleteMessageIdentity() (securityModeCompleteMessageIdentity *SecurityModeCompleteMessageIdentity) {}

// SecurityModeCompleteMessageIdentity 9.6
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *SecurityModeCompleteMessageIdentity) GetMessageType() (messageType uint8) {}

// SecurityModeCompleteMessageIdentity 9.6
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *SecurityModeCompleteMessageIdentity) SetMessageType(messageType uint8) {}
