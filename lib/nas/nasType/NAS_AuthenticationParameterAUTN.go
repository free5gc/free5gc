//go:binary-only-package

package nasType

// AuthenticationParameterAUTN 9.11.3.15
// AUTN Row, sBit, len = [0, 15], 8 , 128
type AuthenticationParameterAUTN struct {
	Iei   uint8
	Len   uint8
	Octet [16]uint8
}

func NewAuthenticationParameterAUTN(iei uint8) (authenticationParameterAUTN *AuthenticationParameterAUTN) {}

// AuthenticationParameterAUTN 9.11.3.15
// Iei Row, sBit, len = [], 8, 8
func (a *AuthenticationParameterAUTN) GetIei() (iei uint8) {}

// AuthenticationParameterAUTN 9.11.3.15
// Iei Row, sBit, len = [], 8, 8
func (a *AuthenticationParameterAUTN) SetIei(iei uint8) {}

// AuthenticationParameterAUTN 9.11.3.15
// Len Row, sBit, len = [], 8, 8
func (a *AuthenticationParameterAUTN) GetLen() (len uint8) {}

// AuthenticationParameterAUTN 9.11.3.15
// Len Row, sBit, len = [], 8, 8
func (a *AuthenticationParameterAUTN) SetLen(len uint8) {}

// AuthenticationParameterAUTN 9.11.3.15
// AUTN Row, sBit, len = [0, 15], 8 , 128
func (a *AuthenticationParameterAUTN) GetAUTN() (aUTN [16]uint8) {}

// AuthenticationParameterAUTN 9.11.3.15
// AUTN Row, sBit, len = [0, 15], 8 , 128
func (a *AuthenticationParameterAUTN) SetAUTN(aUTN [16]uint8) {}
