//go:binary-only-package

package nasType

// T3502Value 9.11.2.4
// GPRSTimer2Value Row, sBit, len = [0, 0], 8 , 8
type T3502Value struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewT3502Value(iei uint8) (t3502Value *T3502Value) {}

// T3502Value 9.11.2.4
// Iei Row, sBit, len = [], 8, 8
func (a *T3502Value) GetIei() (iei uint8) {}

// T3502Value 9.11.2.4
// Iei Row, sBit, len = [], 8, 8
func (a *T3502Value) SetIei(iei uint8) {}

// T3502Value 9.11.2.4
// Len Row, sBit, len = [], 8, 8
func (a *T3502Value) GetLen() (len uint8) {}

// T3502Value 9.11.2.4
// Len Row, sBit, len = [], 8, 8
func (a *T3502Value) SetLen(len uint8) {}

// T3502Value 9.11.2.4
// GPRSTimer2Value Row, sBit, len = [0, 0], 8 , 8
func (a *T3502Value) GetGPRSTimer2Value() (gPRSTimer2Value uint8) {}

// T3502Value 9.11.2.4
// GPRSTimer2Value Row, sBit, len = [0, 0], 8 , 8
func (a *T3502Value) SetGPRSTimer2Value(gPRSTimer2Value uint8) {}
