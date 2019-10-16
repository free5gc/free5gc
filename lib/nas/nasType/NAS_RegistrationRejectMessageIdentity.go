//go:binary-only-package

package nasType

// RegistrationRejectMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type RegistrationRejectMessageIdentity struct {
	Octet uint8
}

func NewRegistrationRejectMessageIdentity() (registrationRejectMessageIdentity *RegistrationRejectMessageIdentity) {}

// RegistrationRejectMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *RegistrationRejectMessageIdentity) GetMessageType() (messageType uint8) {}

// RegistrationRejectMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *RegistrationRejectMessageIdentity) SetMessageType(messageType uint8) {}
