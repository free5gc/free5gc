//go:binary-only-package

package nasType

// ServiceRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type ServiceRequestMessageIdentity struct {
	Octet uint8
}

func NewServiceRequestMessageIdentity() (serviceRequestMessageIdentity *ServiceRequestMessageIdentity) {}

// ServiceRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ServiceRequestMessageIdentity) GetMessageType() (messageType uint8) {}

// ServiceRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ServiceRequestMessageIdentity) SetMessageType(messageType uint8) {}
