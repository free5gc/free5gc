//go:binary-only-package

package nasType

// SecurityModeRejectMessageIdentity 9.6
// MessageType Row, sBit, len = [0, 0], 8 , 8
type SecurityModeRejectMessageIdentity struct {
	Octet uint8
}

func NewSecurityModeRejectMessageIdentity() (securityModeRejectMessageIdentity *SecurityModeRejectMessageIdentity) {}

// SecurityModeRejectMessageIdentity 9.6
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *SecurityModeRejectMessageIdentity) GetMessageType() (messageType uint8) {}

// SecurityModeRejectMessageIdentity 9.6
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *SecurityModeRejectMessageIdentity) SetMessageType(messageType uint8) {}
