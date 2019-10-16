//go:binary-only-package

package nasType

// NegotiatedDRXParameters 9.11.3.2A
// DRXValue Row, sBit, len = [0, 0], 4 , 4
type NegotiatedDRXParameters struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewNegotiatedDRXParameters(iei uint8) (negotiatedDRXParameters *NegotiatedDRXParameters) {}

// NegotiatedDRXParameters 9.11.3.2A
// Iei Row, sBit, len = [], 8, 8
func (a *NegotiatedDRXParameters) GetIei() (iei uint8) {}

// NegotiatedDRXParameters 9.11.3.2A
// Iei Row, sBit, len = [], 8, 8
func (a *NegotiatedDRXParameters) SetIei(iei uint8) {}

// NegotiatedDRXParameters 9.11.3.2A
// Len Row, sBit, len = [], 8, 8
func (a *NegotiatedDRXParameters) GetLen() (len uint8) {}

// NegotiatedDRXParameters 9.11.3.2A
// Len Row, sBit, len = [], 8, 8
func (a *NegotiatedDRXParameters) SetLen(len uint8) {}

// NegotiatedDRXParameters 9.11.3.2A
// DRXValue Row, sBit, len = [0, 0], 4 , 4
func (a *NegotiatedDRXParameters) GetDRXValue() (dRXValue uint8) {}

// NegotiatedDRXParameters 9.11.3.2A
// DRXValue Row, sBit, len = [0, 0], 4 , 4
func (a *NegotiatedDRXParameters) SetDRXValue(dRXValue uint8) {}
