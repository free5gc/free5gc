//go:binary-only-package

package nasType

// ExtendedEmergencyNumberList 9.11.3.26
// EENL Row, sBit, len = [0, 0], 1 , 1
// EmergencyInformation Row, sBit, len = [0, 0], 8 , INF
type ExtendedEmergencyNumberList struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewExtendedEmergencyNumberList(iei uint8) (extendedEmergencyNumberList *ExtendedEmergencyNumberList) {}

// ExtendedEmergencyNumberList 9.11.3.26
// Iei Row, sBit, len = [], 8, 8
func (a *ExtendedEmergencyNumberList) GetIei() (iei uint8) {}

// ExtendedEmergencyNumberList 9.11.3.26
// Iei Row, sBit, len = [], 8, 8
func (a *ExtendedEmergencyNumberList) SetIei(iei uint8) {}

// ExtendedEmergencyNumberList 9.11.3.26
// Len Row, sBit, len = [], 8, 16
func (a *ExtendedEmergencyNumberList) GetLen() (len uint16) {}

// ExtendedEmergencyNumberList 9.11.3.26
// Len Row, sBit, len = [], 8, 16
func (a *ExtendedEmergencyNumberList) SetLen(len uint16) {}

// ExtendedEmergencyNumberList 9.11.3.26
// EENL Row, sBit, len = [0, 0], 1 , 1
func (a *ExtendedEmergencyNumberList) GetEENL() (eENL uint8) {}

// ExtendedEmergencyNumberList 9.11.3.26
// EENL Row, sBit, len = [0, 0], 1 , 1
func (a *ExtendedEmergencyNumberList) SetEENL(eENL uint8) {}

// ExtendedEmergencyNumberList 9.11.3.26
// EmergencyInformation Row, sBit, len = [0, 0], 8 , INF
func (a *ExtendedEmergencyNumberList) GetEmergencyInformation() (emergencyInformation []uint8) {}

// ExtendedEmergencyNumberList 9.11.3.26
// EmergencyInformation Row, sBit, len = [0, 0], 8 , INF
func (a *ExtendedEmergencyNumberList) SetEmergencyInformation(emergencyInformation []uint8) {}
