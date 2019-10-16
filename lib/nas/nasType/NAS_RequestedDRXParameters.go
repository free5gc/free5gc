//go:binary-only-package

package nasType

// RequestedDRXParameters 9.11.3.2A
// DRXValue Row, sBit, len = [0, 0], 4 , 4
type RequestedDRXParameters struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewRequestedDRXParameters(iei uint8) (requestedDRXParameters *RequestedDRXParameters) {}

// RequestedDRXParameters 9.11.3.2A
// Iei Row, sBit, len = [], 8, 8
func (a *RequestedDRXParameters) GetIei() (iei uint8) {}

// RequestedDRXParameters 9.11.3.2A
// Iei Row, sBit, len = [], 8, 8
func (a *RequestedDRXParameters) SetIei(iei uint8) {}

// RequestedDRXParameters 9.11.3.2A
// Len Row, sBit, len = [], 8, 8
func (a *RequestedDRXParameters) GetLen() (len uint8) {}

// RequestedDRXParameters 9.11.3.2A
// Len Row, sBit, len = [], 8, 8
func (a *RequestedDRXParameters) SetLen(len uint8) {}

// RequestedDRXParameters 9.11.3.2A
// DRXValue Row, sBit, len = [0, 0], 4 , 4
func (a *RequestedDRXParameters) GetDRXValue() (dRXValue uint8) {}

// RequestedDRXParameters 9.11.3.2A
// DRXValue Row, sBit, len = [0, 0], 4 , 4
func (a *RequestedDRXParameters) SetDRXValue(dRXValue uint8) {}
