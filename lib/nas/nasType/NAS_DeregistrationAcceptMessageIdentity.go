//go:binary-only-package

package nasType

// DeregistrationAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type DeregistrationAcceptMessageIdentity struct {
	Octet uint8
}

func NewDeregistrationAcceptMessageIdentity() (deregistrationAcceptMessageIdentity *DeregistrationAcceptMessageIdentity) {}

// DeregistrationAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *DeregistrationAcceptMessageIdentity) GetMessageType() (messageType uint8) {}

// DeregistrationAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *DeregistrationAcceptMessageIdentity) SetMessageType(messageType uint8) {}
