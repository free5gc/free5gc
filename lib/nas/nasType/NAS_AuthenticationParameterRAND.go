//go:binary-only-package

package nasType

// AuthenticationParameterRAND 9.11.3.16
// RANDValue Row, sBit, len = [0, 15], 8 , 128
type AuthenticationParameterRAND struct {
	Iei   uint8
	Octet [16]uint8
}

func NewAuthenticationParameterRAND(iei uint8) (authenticationParameterRAND *AuthenticationParameterRAND) {}

// AuthenticationParameterRAND 9.11.3.16
// Iei Row, sBit, len = [], 8, 8
func (a *AuthenticationParameterRAND) GetIei() (iei uint8) {}

// AuthenticationParameterRAND 9.11.3.16
// Iei Row, sBit, len = [], 8, 8
func (a *AuthenticationParameterRAND) SetIei(iei uint8) {}

// AuthenticationParameterRAND 9.11.3.16
// RANDValue Row, sBit, len = [0, 15], 8 , 128
func (a *AuthenticationParameterRAND) GetRANDValue() (rANDValue [16]uint8) {}

// AuthenticationParameterRAND 9.11.3.16
// RANDValue Row, sBit, len = [0, 15], 8 , 128
func (a *AuthenticationParameterRAND) SetRANDValue(rANDValue [16]uint8) {}
