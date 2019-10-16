//go:binary-only-package

package nasType

// RegistrationCompleteMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type RegistrationCompleteMessageIdentity struct {
	Octet uint8
}

func NewRegistrationCompleteMessageIdentity() (registrationCompleteMessageIdentity *RegistrationCompleteMessageIdentity) {}

// RegistrationCompleteMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *RegistrationCompleteMessageIdentity) GetMessageType() (messageType uint8) {}

// RegistrationCompleteMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *RegistrationCompleteMessageIdentity) SetMessageType(messageType uint8) {}
