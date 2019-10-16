//go:binary-only-package

package nasType

// EmergencyNumberList 9.11.3.23
// Lengthof1EmergencyNumberInformation Row, sBit, len = [0, 0], 8 , 8
// EmergencyServiceCategoryValue Row, sBit, len = [1, 1], 5 , 5
// EmergencyInformation Row, sBit, len = [0, 0], 8 , INF
type EmergencyNumberList struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewEmergencyNumberList(iei uint8) (emergencyNumberList *EmergencyNumberList) {}

// EmergencyNumberList 9.11.3.23
// Iei Row, sBit, len = [], 8, 8
func (a *EmergencyNumberList) GetIei() (iei uint8) {}

// EmergencyNumberList 9.11.3.23
// Iei Row, sBit, len = [], 8, 8
func (a *EmergencyNumberList) SetIei(iei uint8) {}

// EmergencyNumberList 9.11.3.23
// Len Row, sBit, len = [], 8, 8
func (a *EmergencyNumberList) GetLen() (len uint8) {}

// EmergencyNumberList 9.11.3.23
// Len Row, sBit, len = [], 8, 8
func (a *EmergencyNumberList) SetLen(len uint8) {}

// EmergencyNumberList 9.11.3.23
// Lengthof1EmergencyNumberInformation Row, sBit, len = [0, 0], 8 , 8
func (a *EmergencyNumberList) GetLengthof1EmergencyNumberInformation() (lengthof1EmergencyNumberInformation uint8) {}

// EmergencyNumberList 9.11.3.23
// Lengthof1EmergencyNumberInformation Row, sBit, len = [0, 0], 8 , 8
func (a *EmergencyNumberList) SetLengthof1EmergencyNumberInformation(lengthof1EmergencyNumberInformation uint8) {}

// EmergencyNumberList 9.11.3.23
// EmergencyServiceCategoryValue Row, sBit, len = [1, 1], 5 , 5
func (a *EmergencyNumberList) GetEmergencyServiceCategoryValue() (emergencyServiceCategoryValue uint8) {}

// EmergencyNumberList 9.11.3.23
// EmergencyServiceCategoryValue Row, sBit, len = [1, 1], 5 , 5
func (a *EmergencyNumberList) SetEmergencyServiceCategoryValue(emergencyServiceCategoryValue uint8) {}

// EmergencyNumberList 9.11.3.23
// EmergencyInformation Row, sBit, len = [0, 0], 8 , INF
func (a *EmergencyNumberList) GetEmergencyInformation() (emergencyInformation []uint8) {}

// EmergencyNumberList 9.11.3.23
// EmergencyInformation Row, sBit, len = [0, 0], 8 , INF
func (a *EmergencyNumberList) SetEmergencyInformation(emergencyInformation []uint8) {}
