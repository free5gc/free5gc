//go:binary-only-package

package nasType

// AdditionalInformation 9.11.2.1
// AdditionalInformationValue Row, sBit, len = [0, 0], 8 , INF
type AdditionalInformation struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewAdditionalInformation(iei uint8) (additionalInformation *AdditionalInformation) {}

// AdditionalInformation 9.11.2.1
// Iei Row, sBit, len = [], 8, 8
func (a *AdditionalInformation) GetIei() (iei uint8) {}

// AdditionalInformation 9.11.2.1
// Iei Row, sBit, len = [], 8, 8
func (a *AdditionalInformation) SetIei(iei uint8) {}

// AdditionalInformation 9.11.2.1
// Len Row, sBit, len = [], 8, 8
func (a *AdditionalInformation) GetLen() (len uint8) {}

// AdditionalInformation 9.11.2.1
// Len Row, sBit, len = [], 8, 8
func (a *AdditionalInformation) SetLen(len uint8) {}

// AdditionalInformation 9.11.2.1
// AdditionalInformationValue Row, sBit, len = [0, 0], 8 , INF
func (a *AdditionalInformation) GetAdditionalInformationValue() (additionalInformationValue []uint8) {}

// AdditionalInformation 9.11.2.1
// AdditionalInformationValue Row, sBit, len = [0, 0], 8 , INF
func (a *AdditionalInformation) SetAdditionalInformationValue(additionalInformationValue []uint8) {}
