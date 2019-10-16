//go:binary-only-package

package nasType

// AuthenticationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type AuthenticationRequestMessageIdentity struct {
	Octet uint8
}

func NewAuthenticationRequestMessageIdentity() (authenticationRequestMessageIdentity *AuthenticationRequestMessageIdentity) {}

// AuthenticationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *AuthenticationRequestMessageIdentity) GetMessageType() (messageType uint8) {}

// AuthenticationRequestMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *AuthenticationRequestMessageIdentity) SetMessageType(messageType uint8) {}
