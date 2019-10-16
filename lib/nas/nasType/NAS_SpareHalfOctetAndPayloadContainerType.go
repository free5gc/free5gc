//go:binary-only-package

package nasType

// SpareHalfOctetAndPayloadContainerType 9.11.3.40 9.5
// PayloadContainerType Row, sBit, len = [0, 0], 4 , 4
type SpareHalfOctetAndPayloadContainerType struct {
	Octet uint8
}

func NewSpareHalfOctetAndPayloadContainerType() (spareHalfOctetAndPayloadContainerType *SpareHalfOctetAndPayloadContainerType) {}

// SpareHalfOctetAndPayloadContainerType 9.11.3.40 9.5
// PayloadContainerType Row, sBit, len = [0, 0], 4 , 4
func (a *SpareHalfOctetAndPayloadContainerType) GetPayloadContainerType() (payloadContainerType uint8) {}

// SpareHalfOctetAndPayloadContainerType 9.11.3.40 9.5
// PayloadContainerType Row, sBit, len = [0, 0], 4 , 4
func (a *SpareHalfOctetAndPayloadContainerType) SetPayloadContainerType(payloadContainerType uint8) {}
