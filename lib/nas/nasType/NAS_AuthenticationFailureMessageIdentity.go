//go:binary-only-package

package nasType

// AuthenticationFailureMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type AuthenticationFailureMessageIdentity struct {
	Octet uint8
}

func NewAuthenticationFailureMessageIdentity() (authenticationFailureMessageIdentity *AuthenticationFailureMessageIdentity) {}

// AuthenticationFailureMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *AuthenticationFailureMessageIdentity) GetMessageType() (messageType uint8) {}

// AuthenticationFailureMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *AuthenticationFailureMessageIdentity) SetMessageType(messageType uint8) {}
