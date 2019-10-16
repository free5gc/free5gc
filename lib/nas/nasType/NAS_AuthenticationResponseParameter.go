//go:binary-only-package

package nasType

// AuthenticationResponseParameter 9.11.3.17
// RES Row, sBit, len = [0, 15], 8 , 128
type AuthenticationResponseParameter struct {
	Iei   uint8
	Len   uint8
	Octet [16]uint8
}

func NewAuthenticationResponseParameter(iei uint8) (authenticationResponseParameter *AuthenticationResponseParameter) {}

// AuthenticationResponseParameter 9.11.3.17
// Iei Row, sBit, len = [], 8, 8
func (a *AuthenticationResponseParameter) GetIei() (iei uint8) {}

// AuthenticationResponseParameter 9.11.3.17
// Iei Row, sBit, len = [], 8, 8
func (a *AuthenticationResponseParameter) SetIei(iei uint8) {}

// AuthenticationResponseParameter 9.11.3.17
// Len Row, sBit, len = [], 8, 8
func (a *AuthenticationResponseParameter) GetLen() (len uint8) {}

// AuthenticationResponseParameter 9.11.3.17
// Len Row, sBit, len = [], 8, 8
func (a *AuthenticationResponseParameter) SetLen(len uint8) {}

// AuthenticationResponseParameter 9.11.3.17
// RES Row, sBit, len = [0, 15], 8 , 128
func (a *AuthenticationResponseParameter) GetRES() (rES [16]uint8) {}

// AuthenticationResponseParameter 9.11.3.17
// RES Row, sBit, len = [0, 15], 8 , 128
func (a *AuthenticationResponseParameter) SetRES(rES [16]uint8) {}
