//go:binary-only-package

package nasType

// UEStatus 9.11.3.56
// N1ModeReg Row, sBit, len = [0, 0], 2 , 1
// S1ModeReg Row, sBit, len = [0, 0], 1 , 1
type UEStatus struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewUEStatus(iei uint8) (uEStatus *UEStatus) {}

// UEStatus 9.11.3.56
// Iei Row, sBit, len = [], 8, 8
func (a *UEStatus) GetIei() (iei uint8) {}

// UEStatus 9.11.3.56
// Iei Row, sBit, len = [], 8, 8
func (a *UEStatus) SetIei(iei uint8) {}

// UEStatus 9.11.3.56
// Len Row, sBit, len = [], 8, 8
func (a *UEStatus) GetLen() (len uint8) {}

// UEStatus 9.11.3.56
// Len Row, sBit, len = [], 8, 8
func (a *UEStatus) SetLen(len uint8) {}

// UEStatus 9.11.3.56
// N1ModeReg Row, sBit, len = [0, 0], 2 , 1
func (a *UEStatus) GetN1ModeReg() (n1ModeReg uint8) {}

// UEStatus 9.11.3.56
// N1ModeReg Row, sBit, len = [0, 0], 2 , 1
func (a *UEStatus) SetN1ModeReg(n1ModeReg uint8) {}

// UEStatus 9.11.3.56
// S1ModeReg Row, sBit, len = [0, 0], 1 , 1
func (a *UEStatus) GetS1ModeReg() (s1ModeReg uint8) {}

// UEStatus 9.11.3.56
// S1ModeReg Row, sBit, len = [0, 0], 1 , 1
func (a *UEStatus) SetS1ModeReg(s1ModeReg uint8) {}
