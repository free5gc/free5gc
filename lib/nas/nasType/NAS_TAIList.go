//go:binary-only-package

package nasType

// TAIList 9.11.3.9
// PartialTrackingAreaIdentityList Row, sBit, len = [0, 0], 8 , INF
type TAIList struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewTAIList(iei uint8) (tAIList *TAIList) {}

// TAIList 9.11.3.9
// Iei Row, sBit, len = [], 8, 8
func (a *TAIList) GetIei() (iei uint8) {}

// TAIList 9.11.3.9
// Iei Row, sBit, len = [], 8, 8
func (a *TAIList) SetIei(iei uint8) {}

// TAIList 9.11.3.9
// Len Row, sBit, len = [], 8, 8
func (a *TAIList) GetLen() (len uint8) {}

// TAIList 9.11.3.9
// Len Row, sBit, len = [], 8, 8
func (a *TAIList) SetLen(len uint8) {}

// TAIList 9.11.3.9
// PartialTrackingAreaIdentityList Row, sBit, len = [0, 0], 8 , INF
func (a *TAIList) GetPartialTrackingAreaIdentityList() (partialTrackingAreaIdentityList []uint8) {}

// TAIList 9.11.3.9
// PartialTrackingAreaIdentityList Row, sBit, len = [0, 0], 8 , INF
func (a *TAIList) SetPartialTrackingAreaIdentityList(partialTrackingAreaIdentityList []uint8) {}
