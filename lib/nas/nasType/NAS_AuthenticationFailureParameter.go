//go:binary-only-package

package nasType

// AuthenticationFailureParameter 9.11.3.14
// AuthenticationFailureParameter Row, sBit, len = [0, 13], 8 , 112
type AuthenticationFailureParameter struct {
	Iei   uint8
	Len   uint8
	Octet [14]uint8
}

func NewAuthenticationFailureParameter(iei uint8) (authenticationFailureParameter *AuthenticationFailureParameter) {}

// AuthenticationFailureParameter 9.11.3.14
// Iei Row, sBit, len = [], 8, 8
func (a *AuthenticationFailureParameter) GetIei() (iei uint8) {}

// AuthenticationFailureParameter 9.11.3.14
// Iei Row, sBit, len = [], 8, 8
func (a *AuthenticationFailureParameter) SetIei(iei uint8) {}

// AuthenticationFailureParameter 9.11.3.14
// Len Row, sBit, len = [], 8, 8
func (a *AuthenticationFailureParameter) GetLen() (len uint8) {}

// AuthenticationFailureParameter 9.11.3.14
// Len Row, sBit, len = [], 8, 8
func (a *AuthenticationFailureParameter) SetLen(len uint8) {}

// AuthenticationFailureParameter 9.11.3.14
// AuthenticationFailureParameter Row, sBit, len = [0, 13], 8 , 112
func (a *AuthenticationFailureParameter) GetAuthenticationFailureParameter() (authenticationFailureParameter [14]uint8) {}

// AuthenticationFailureParameter 9.11.3.14
// AuthenticationFailureParameter Row, sBit, len = [0, 13], 8 , 112
func (a *AuthenticationFailureParameter) SetAuthenticationFailureParameter(authenticationFailureParameter [14]uint8) {}
