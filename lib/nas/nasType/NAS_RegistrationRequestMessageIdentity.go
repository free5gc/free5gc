//go:binary-only-package

package nasType

// RegistrationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type RegistrationRequestMessageIdentity struct {
	Octet uint8
}

func NewRegistrationRequestMessageIdentity() (registrationRequestMessageIdentity *RegistrationRequestMessageIdentity) {}

// RegistrationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *RegistrationRequestMessageIdentity) GetMessageType() (messageType uint8) {}

// RegistrationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *RegistrationRequestMessageIdentity) SetMessageType(messageType uint8) {}
