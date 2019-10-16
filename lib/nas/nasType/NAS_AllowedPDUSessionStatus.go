//go:binary-only-package

package nasType

// AllowedPDUSessionStatus 9.11.3.13
// PSI7 Row, sBit, len = [0, 0], 8 , 1
// PSI6 Row, sBit, len = [0, 0], 7 , 1
// PSI5 Row, sBit, len = [0, 0], 6 , 1
// PSI4 Row, sBit, len = [0, 0], 5 , 1
// PSI3 Row, sBit, len = [0, 0], 4 , 1
// PSI2 Row, sBit, len = [0, 0], 3 , 1
// PSI1 Row, sBit, len = [0, 0], 2 , 1
// PSI0 Row, sBit, len = [0, 0], 1 , 1
// PSI15 Row, sBit, len = [1, 1], 8 , 1
// PSI14 Row, sBit, len = [1, 1], 7 , 1
// PSI13 Row, sBit, len = [1, 1], 6 , 1
// PSI12 Row, sBit, len = [1, 1], 5 , 1
// PSI11 Row, sBit, len = [1, 1], 4 , 1
// PSI10 Row, sBit, len = [1, 1], 3 , 1
// PSI9 Row, sBit, len = [1, 1], 2 , 1
// PSI8 Row, sBit, len = [1, 1], 1 , 1
// Spare Row, sBit, len = [2, 2], 1 , INF
type AllowedPDUSessionStatus struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewAllowedPDUSessionStatus(iei uint8) (allowedPDUSessionStatus *AllowedPDUSessionStatus) {}

// AllowedPDUSessionStatus 9.11.3.13
// Iei Row, sBit, len = [], 8, 8
func (a *AllowedPDUSessionStatus) GetIei() (iei uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// Iei Row, sBit, len = [], 8, 8
func (a *AllowedPDUSessionStatus) SetIei(iei uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// Len Row, sBit, len = [], 8, 8
func (a *AllowedPDUSessionStatus) GetLen() (len uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// Len Row, sBit, len = [], 8, 8
func (a *AllowedPDUSessionStatus) SetLen(len uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI7 Row, sBit, len = [0, 0], 8 , 1
func (a *AllowedPDUSessionStatus) GetPSI7() (pSI7 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI7 Row, sBit, len = [0, 0], 8 , 1
func (a *AllowedPDUSessionStatus) SetPSI7(pSI7 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI6 Row, sBit, len = [0, 0], 7 , 1
func (a *AllowedPDUSessionStatus) GetPSI6() (pSI6 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI6 Row, sBit, len = [0, 0], 7 , 1
func (a *AllowedPDUSessionStatus) SetPSI6(pSI6 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI5 Row, sBit, len = [0, 0], 6 , 1
func (a *AllowedPDUSessionStatus) GetPSI5() (pSI5 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI5 Row, sBit, len = [0, 0], 6 , 1
func (a *AllowedPDUSessionStatus) SetPSI5(pSI5 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI4 Row, sBit, len = [0, 0], 5 , 1
func (a *AllowedPDUSessionStatus) GetPSI4() (pSI4 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI4 Row, sBit, len = [0, 0], 5 , 1
func (a *AllowedPDUSessionStatus) SetPSI4(pSI4 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI3 Row, sBit, len = [0, 0], 4 , 1
func (a *AllowedPDUSessionStatus) GetPSI3() (pSI3 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI3 Row, sBit, len = [0, 0], 4 , 1
func (a *AllowedPDUSessionStatus) SetPSI3(pSI3 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI2 Row, sBit, len = [0, 0], 3 , 1
func (a *AllowedPDUSessionStatus) GetPSI2() (pSI2 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI2 Row, sBit, len = [0, 0], 3 , 1
func (a *AllowedPDUSessionStatus) SetPSI2(pSI2 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI1 Row, sBit, len = [0, 0], 2 , 1
func (a *AllowedPDUSessionStatus) GetPSI1() (pSI1 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI1 Row, sBit, len = [0, 0], 2 , 1
func (a *AllowedPDUSessionStatus) SetPSI1(pSI1 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI0 Row, sBit, len = [0, 0], 1 , 1
func (a *AllowedPDUSessionStatus) GetPSI0() (pSI0 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI0 Row, sBit, len = [0, 0], 1 , 1
func (a *AllowedPDUSessionStatus) SetPSI0(pSI0 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI15 Row, sBit, len = [1, 1], 8 , 1
func (a *AllowedPDUSessionStatus) GetPSI15() (pSI15 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI15 Row, sBit, len = [1, 1], 8 , 1
func (a *AllowedPDUSessionStatus) SetPSI15(pSI15 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI14 Row, sBit, len = [1, 1], 7 , 1
func (a *AllowedPDUSessionStatus) GetPSI14() (pSI14 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI14 Row, sBit, len = [1, 1], 7 , 1
func (a *AllowedPDUSessionStatus) SetPSI14(pSI14 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI13 Row, sBit, len = [1, 1], 6 , 1
func (a *AllowedPDUSessionStatus) GetPSI13() (pSI13 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI13 Row, sBit, len = [1, 1], 6 , 1
func (a *AllowedPDUSessionStatus) SetPSI13(pSI13 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI12 Row, sBit, len = [1, 1], 5 , 1
func (a *AllowedPDUSessionStatus) GetPSI12() (pSI12 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI12 Row, sBit, len = [1, 1], 5 , 1
func (a *AllowedPDUSessionStatus) SetPSI12(pSI12 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI11 Row, sBit, len = [1, 1], 4 , 1
func (a *AllowedPDUSessionStatus) GetPSI11() (pSI11 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI11 Row, sBit, len = [1, 1], 4 , 1
func (a *AllowedPDUSessionStatus) SetPSI11(pSI11 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI10 Row, sBit, len = [1, 1], 3 , 1
func (a *AllowedPDUSessionStatus) GetPSI10() (pSI10 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI10 Row, sBit, len = [1, 1], 3 , 1
func (a *AllowedPDUSessionStatus) SetPSI10(pSI10 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI9 Row, sBit, len = [1, 1], 2 , 1
func (a *AllowedPDUSessionStatus) GetPSI9() (pSI9 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI9 Row, sBit, len = [1, 1], 2 , 1
func (a *AllowedPDUSessionStatus) SetPSI9(pSI9 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI8 Row, sBit, len = [1, 1], 1 , 1
func (a *AllowedPDUSessionStatus) GetPSI8() (pSI8 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// PSI8 Row, sBit, len = [1, 1], 1 , 1
func (a *AllowedPDUSessionStatus) SetPSI8(pSI8 uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// Spare Row, sBit, len = [2, 2], 1 , INF
func (a *AllowedPDUSessionStatus) GetSpare() (spare []uint8) {}

// AllowedPDUSessionStatus 9.11.3.13
// Spare Row, sBit, len = [2, 2], 1 , INF
func (a *AllowedPDUSessionStatus) SetSpare(spare []uint8) {}
