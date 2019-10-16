//go:binary-only-package

package nasType

// T3346Value 9.11.2.4
// GPRSTimer2Value Row, sBit, len = [0, 0], 8 , 8
type T3346Value struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewT3346Value(iei uint8) (t3346Value *T3346Value) {}

// T3346Value 9.11.2.4
// Iei Row, sBit, len = [], 8, 8
func (a *T3346Value) GetIei() (iei uint8) {}

// T3346Value 9.11.2.4
// Iei Row, sBit, len = [], 8, 8
func (a *T3346Value) SetIei(iei uint8) {}

// T3346Value 9.11.2.4
// Len Row, sBit, len = [], 8, 8
func (a *T3346Value) GetLen() (len uint8) {}

// T3346Value 9.11.2.4
// Len Row, sBit, len = [], 8, 8
func (a *T3346Value) SetLen(len uint8) {}

// T3346Value 9.11.2.4
// GPRSTimer2Value Row, sBit, len = [0, 0], 8 , 8
func (a *T3346Value) GetGPRSTimer2Value() (gPRSTimer2Value uint8) {}

// T3346Value 9.11.2.4
// GPRSTimer2Value Row, sBit, len = [0, 0], 8 , 8
func (a *T3346Value) SetGPRSTimer2Value(gPRSTimer2Value uint8) {}
