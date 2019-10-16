//go:binary-only-package

package nasType

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// TSC Row, sBit, len = [0, 0], 8 , 1
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 7 , 3
// SwitchOff Row, sBit, len = [0, 0], 4 , 1
// ReRegistrationRequired Row, sBit, len = [0, 0], 3 , 1
// AccessType Row, sBit, len = [0, 0], 2 , 2
type NgksiAndDeregistrationType struct {
	Octet uint8
}

func NewNgksiAndDeregistrationType() (ngksiAndDeregistrationType *NgksiAndDeregistrationType) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// TSC Row, sBit, len = [0, 0], 8 , 1
func (a *NgksiAndDeregistrationType) GetTSC() (tSC uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// TSC Row, sBit, len = [0, 0], 8 , 1
func (a *NgksiAndDeregistrationType) SetTSC(tSC uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 7 , 3
func (a *NgksiAndDeregistrationType) GetNasKeySetIdentifiler() (nasKeySetIdentifiler uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 7 , 3
func (a *NgksiAndDeregistrationType) SetNasKeySetIdentifiler(nasKeySetIdentifiler uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// SwitchOff Row, sBit, len = [0, 0], 4 , 1
func (a *NgksiAndDeregistrationType) GetSwitchOff() (switchOff uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// SwitchOff Row, sBit, len = [0, 0], 4 , 1
func (a *NgksiAndDeregistrationType) SetSwitchOff(switchOff uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// ReRegistrationRequired Row, sBit, len = [0, 0], 3 , 1
func (a *NgksiAndDeregistrationType) GetReRegistrationRequired() (reRegistrationRequired uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// ReRegistrationRequired Row, sBit, len = [0, 0], 3 , 1
func (a *NgksiAndDeregistrationType) SetReRegistrationRequired(reRegistrationRequired uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// AccessType Row, sBit, len = [0, 0], 2 , 2
func (a *NgksiAndDeregistrationType) GetAccessType() (accessType uint8) {}

// NgksiAndDeregistrationType 9.11.3.20 9.11.3.32
// AccessType Row, sBit, len = [0, 0], 2 , 2
func (a *NgksiAndDeregistrationType) SetAccessType(accessType uint8) {}
