//go:binary-only-package

package nasType

// AuthenticationResponseMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type AuthenticationResponseMessageIdentity struct {
	Octet uint8
}

func NewAuthenticationResponseMessageIdentity() (authenticationResponseMessageIdentity *AuthenticationResponseMessageIdentity) {}

// AuthenticationResponseMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *AuthenticationResponseMessageIdentity) GetMessageType() (messageType uint8) {}

// AuthenticationResponseMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *AuthenticationResponseMessageIdentity) SetMessageType(messageType uint8) {}
