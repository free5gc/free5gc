//go:binary-only-package

package nasType

// NetworkDaylightSavingTime 9.11.3.19
// value Row, sBit, len = [0, 0], 2 , 2
type NetworkDaylightSavingTime struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewNetworkDaylightSavingTime(iei uint8) (networkDaylightSavingTime *NetworkDaylightSavingTime) {}

// NetworkDaylightSavingTime 9.11.3.19
// Iei Row, sBit, len = [], 8, 8
func (a *NetworkDaylightSavingTime) GetIei() (iei uint8) {}

// NetworkDaylightSavingTime 9.11.3.19
// Iei Row, sBit, len = [], 8, 8
func (a *NetworkDaylightSavingTime) SetIei(iei uint8) {}

// NetworkDaylightSavingTime 9.11.3.19
// Len Row, sBit, len = [], 8, 8
func (a *NetworkDaylightSavingTime) GetLen() (len uint8) {}

// NetworkDaylightSavingTime 9.11.3.19
// Len Row, sBit, len = [], 8, 8
func (a *NetworkDaylightSavingTime) SetLen(len uint8) {}

// NetworkDaylightSavingTime 9.11.3.19
// value Row, sBit, len = [0, 0], 2 , 2
func (a *NetworkDaylightSavingTime) Getvalue() (value uint8) {}

// NetworkDaylightSavingTime 9.11.3.19
// value Row, sBit, len = [0, 0], 2 , 2
func (a *NetworkDaylightSavingTime) Setvalue(value uint8) {}
