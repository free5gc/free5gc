//go:binary-only-package

package nasType

// SpareHalfOctetAndAccessType 9.11.3.11 9.5
// AccessType Row, sBit, len = [0, 0], 2 , 2
type SpareHalfOctetAndAccessType struct {
	Octet uint8
}

func NewSpareHalfOctetAndAccessType() (spareHalfOctetAndAccessType *SpareHalfOctetAndAccessType) {}

// SpareHalfOctetAndAccessType 9.11.3.11 9.5
// AccessType Row, sBit, len = [0, 0], 2 , 2
func (a *SpareHalfOctetAndAccessType) GetAccessType() (accessType uint8) {}

// SpareHalfOctetAndAccessType 9.11.3.11 9.5
// AccessType Row, sBit, len = [0, 0], 2 , 2
func (a *SpareHalfOctetAndAccessType) SetAccessType(accessType uint8) {}
