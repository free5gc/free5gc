//go:binary-only-package

package nasType

// RegistrationResult5GS 9.11.3.6
// SMSAllowed Row, sBit, len = [0, 0], 4 , 1
// RegistrationResultValue5GS Row, sBit, len = [0, 0], 3 , 3
type RegistrationResult5GS struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewRegistrationResult5GS(iei uint8) (registrationResult5GS *RegistrationResult5GS) {}

// RegistrationResult5GS 9.11.3.6
// Iei Row, sBit, len = [], 8, 8
func (a *RegistrationResult5GS) GetIei() (iei uint8) {}

// RegistrationResult5GS 9.11.3.6
// Iei Row, sBit, len = [], 8, 8
func (a *RegistrationResult5GS) SetIei(iei uint8) {}

// RegistrationResult5GS 9.11.3.6
// Len Row, sBit, len = [], 8, 8
func (a *RegistrationResult5GS) GetLen() (len uint8) {}

// RegistrationResult5GS 9.11.3.6
// Len Row, sBit, len = [], 8, 8
func (a *RegistrationResult5GS) SetLen(len uint8) {}

// RegistrationResult5GS 9.11.3.6
// SMSAllowed Row, sBit, len = [0, 0], 4 , 1
func (a *RegistrationResult5GS) GetSMSAllowed() (sMSAllowed uint8) {}

// RegistrationResult5GS 9.11.3.6
// SMSAllowed Row, sBit, len = [0, 0], 4 , 1
func (a *RegistrationResult5GS) SetSMSAllowed(sMSAllowed uint8) {}

// RegistrationResult5GS 9.11.3.6
// RegistrationResultValue5GS Row, sBit, len = [0, 0], 3 , 3
func (a *RegistrationResult5GS) GetRegistrationResultValue5GS() (registrationResultValue5GS uint8) {}

// RegistrationResult5GS 9.11.3.6
// RegistrationResultValue5GS Row, sBit, len = [0, 0], 3 , 3
func (a *RegistrationResult5GS) SetRegistrationResultValue5GS(registrationResultValue5GS uint8) {}
