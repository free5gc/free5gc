//go:binary-only-package

package nasType

// ServiceAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type ServiceAcceptMessageIdentity struct {
	Octet uint8
}

func NewServiceAcceptMessageIdentity() (serviceAcceptMessageIdentity *ServiceAcceptMessageIdentity) {}

// ServiceAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ServiceAcceptMessageIdentity) GetMessageType() (messageType uint8) {}

// ServiceAcceptMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ServiceAcceptMessageIdentity) SetMessageType(messageType uint8) {}
