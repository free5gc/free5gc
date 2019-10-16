//go:binary-only-package

package nasType

// SpareHalfOctetAndIdentityType 9.11.3.3 9.5
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
type SpareHalfOctetAndIdentityType struct {
	Octet uint8
}

func NewSpareHalfOctetAndIdentityType() (spareHalfOctetAndIdentityType *SpareHalfOctetAndIdentityType) {}

// SpareHalfOctetAndIdentityType 9.11.3.3 9.5
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *SpareHalfOctetAndIdentityType) GetTypeOfIdentity() (typeOfIdentity uint8) {}

// SpareHalfOctetAndIdentityType 9.11.3.3 9.5
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *SpareHalfOctetAndIdentityType) SetTypeOfIdentity(typeOfIdentity uint8) {}
