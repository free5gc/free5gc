//go:binary-only-package

package nasType

// ServiceAreaList 9.11.3.49
// PartialServiceAreaList Row, sBit, len = [0, 0], 8 , INF
type ServiceAreaList struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewServiceAreaList(iei uint8) (serviceAreaList *ServiceAreaList) {}

// ServiceAreaList 9.11.3.49
// Iei Row, sBit, len = [], 8, 8
func (a *ServiceAreaList) GetIei() (iei uint8) {}

// ServiceAreaList 9.11.3.49
// Iei Row, sBit, len = [], 8, 8
func (a *ServiceAreaList) SetIei(iei uint8) {}

// ServiceAreaList 9.11.3.49
// Len Row, sBit, len = [], 8, 8
func (a *ServiceAreaList) GetLen() (len uint8) {}

// ServiceAreaList 9.11.3.49
// Len Row, sBit, len = [], 8, 8
func (a *ServiceAreaList) SetLen(len uint8) {}

// ServiceAreaList 9.11.3.49
// PartialServiceAreaList Row, sBit, len = [0, 0], 8 , INF
func (a *ServiceAreaList) GetPartialServiceAreaList() (partialServiceAreaList []uint8) {}

// ServiceAreaList 9.11.3.49
// PartialServiceAreaList Row, sBit, len = [0, 0], 8 , INF
func (a *ServiceAreaList) SetPartialServiceAreaList(partialServiceAreaList []uint8) {}
