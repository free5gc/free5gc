//go:binary-only-package

package nasType

// RegistrationAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type RegistrationAcceptMessageIdentity struct {
	Octet uint8
}

func NewRegistrationAcceptMessageIdentity() (registrationAcceptMessageIdentity *RegistrationAcceptMessageIdentity) {}

// RegistrationAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *RegistrationAcceptMessageIdentity) GetMessageType() (messageType uint8) {}

// RegistrationAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *RegistrationAcceptMessageIdentity) SetMessageType(messageType uint8) {}
