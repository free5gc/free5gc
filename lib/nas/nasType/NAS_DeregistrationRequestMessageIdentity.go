//go:binary-only-package

package nasType

// DeregistrationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type DeregistrationRequestMessageIdentity struct {
	Octet uint8
}

func NewDeregistrationRequestMessageIdentity() (deregistrationRequestMessageIdentity *DeregistrationRequestMessageIdentity) {}

// DeregistrationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *DeregistrationRequestMessageIdentity) GetMessageType() (messageType uint8) {}

// DeregistrationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *DeregistrationRequestMessageIdentity) SetMessageType(messageType uint8) {}
