//go:binary-only-package

package nasType

// AuthenticationResultMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type AuthenticationResultMessageIdentity struct {
	Octet uint8
}

func NewAuthenticationResultMessageIdentity() (authenticationResultMessageIdentity *AuthenticationResultMessageIdentity) {}

// AuthenticationResultMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *AuthenticationResultMessageIdentity) GetMessageType() (messageType uint8) {}

// AuthenticationResultMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *AuthenticationResultMessageIdentity) SetMessageType(messageType uint8) {}
