//go:binary-only-package

package nasType

// SpareHalfOctetAndDeregistrationType 9.11.3.20 9.5
// SwitchOff Row, sBit, len = [0, 0], 4 , 1
// ReRegistrationRequired Row, sBit, len = [0, 0], 3 , 1
// AccessType Row, sBit, len = [0, 0], 2 , 2
type SpareHalfOctetAndDeregistrationType struct {
	Octet uint8
}

func NewSpareHalfOctetAndDeregistrationType() (spareHalfOctetAndDeregistrationType *SpareHalfOctetAndDeregistrationType) {}

// SpareHalfOctetAndDeregistrationType 9.11.3.20 9.5
// SwitchOff Row, sBit, len = [0, 0], 4 , 1
func (a *SpareHalfOctetAndDeregistrationType) GetSwitchOff() (switchOff uint8) {}

// SpareHalfOctetAndDeregistrationType 9.11.3.20 9.5
// SwitchOff Row, sBit, len = [0, 0], 4 , 1
func (a *SpareHalfOctetAndDeregistrationType) SetSwitchOff(switchOff uint8) {}

// SpareHalfOctetAndDeregistrationType 9.11.3.20 9.5
// ReRegistrationRequired Row, sBit, len = [0, 0], 3 , 1
func (a *SpareHalfOctetAndDeregistrationType) GetReRegistrationRequired() (reRegistrationRequired uint8) {}

// SpareHalfOctetAndDeregistrationType 9.11.3.20 9.5
// ReRegistrationRequired Row, sBit, len = [0, 0], 3 , 1
func (a *SpareHalfOctetAndDeregistrationType) SetReRegistrationRequired(reRegistrationRequired uint8) {}

// SpareHalfOctetAndDeregistrationType 9.11.3.20 9.5
// AccessType Row, sBit, len = [0, 0], 2 , 2
func (a *SpareHalfOctetAndDeregistrationType) GetAccessType() (accessType uint8) {}

// SpareHalfOctetAndDeregistrationType 9.11.3.20 9.5
// AccessType Row, sBit, len = [0, 0], 2 , 2
func (a *SpareHalfOctetAndDeregistrationType) SetAccessType(accessType uint8) {}
