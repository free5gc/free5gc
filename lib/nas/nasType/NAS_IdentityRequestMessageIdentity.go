//go:binary-only-package

package nasType

// IdentityRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type IdentityRequestMessageIdentity struct {
	Octet uint8
}

func NewIdentityRequestMessageIdentity() (identityRequestMessageIdentity *IdentityRequestMessageIdentity) {}

// IdentityRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *IdentityRequestMessageIdentity) GetMessageType() (messageType uint8) {}

// IdentityRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *IdentityRequestMessageIdentity) SetMessageType(messageType uint8) {}
