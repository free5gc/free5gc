//go:binary-only-package

package nasType

// UplinkDataStatus 9.11.3.57
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
type UplinkDataStatus struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewUplinkDataStatus(iei uint8) (uplinkDataStatus *UplinkDataStatus) {}

// UplinkDataStatus 9.11.3.57
// Iei Row, sBit, len = [], 8, 8
func (a *UplinkDataStatus) GetIei() (iei uint8) {}

// UplinkDataStatus 9.11.3.57
// Iei Row, sBit, len = [], 8, 8
func (a *UplinkDataStatus) SetIei(iei uint8) {}

// UplinkDataStatus 9.11.3.57
// Len Row, sBit, len = [], 8, 8
func (a *UplinkDataStatus) GetLen() (len uint8) {}

// UplinkDataStatus 9.11.3.57
// Len Row, sBit, len = [], 8, 8
func (a *UplinkDataStatus) SetLen(len uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI7 Row, sBit, len = [0, 0], 8 , 1
func (a *UplinkDataStatus) GetPSI7() (pSI7 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI7 Row, sBit, len = [0, 0], 8 , 1
func (a *UplinkDataStatus) SetPSI7(pSI7 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI6 Row, sBit, len = [0, 0], 7 , 1
func (a *UplinkDataStatus) GetPSI6() (pSI6 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI6 Row, sBit, len = [0, 0], 7 , 1
func (a *UplinkDataStatus) SetPSI6(pSI6 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI5 Row, sBit, len = [0, 0], 6 , 1
func (a *UplinkDataStatus) GetPSI5() (pSI5 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI5 Row, sBit, len = [0, 0], 6 , 1
func (a *UplinkDataStatus) SetPSI5(pSI5 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI4 Row, sBit, len = [0, 0], 5 , 1
func (a *UplinkDataStatus) GetPSI4() (pSI4 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI4 Row, sBit, len = [0, 0], 5 , 1
func (a *UplinkDataStatus) SetPSI4(pSI4 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI3 Row, sBit, len = [0, 0], 4 , 1
func (a *UplinkDataStatus) GetPSI3() (pSI3 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI3 Row, sBit, len = [0, 0], 4 , 1
func (a *UplinkDataStatus) SetPSI3(pSI3 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI2 Row, sBit, len = [0, 0], 3 , 1
func (a *UplinkDataStatus) GetPSI2() (pSI2 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI2 Row, sBit, len = [0, 0], 3 , 1
func (a *UplinkDataStatus) SetPSI2(pSI2 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI1 Row, sBit, len = [0, 0], 2 , 1
func (a *UplinkDataStatus) GetPSI1() (pSI1 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI1 Row, sBit, len = [0, 0], 2 , 1
func (a *UplinkDataStatus) SetPSI1(pSI1 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI0 Row, sBit, len = [0, 0], 1 , 1
func (a *UplinkDataStatus) GetPSI0() (pSI0 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI0 Row, sBit, len = [0, 0], 1 , 1
func (a *UplinkDataStatus) SetPSI0(pSI0 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI15 Row, sBit, len = [1, 1], 8 , 1
func (a *UplinkDataStatus) GetPSI15() (pSI15 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI15 Row, sBit, len = [1, 1], 8 , 1
func (a *UplinkDataStatus) SetPSI15(pSI15 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI14 Row, sBit, len = [1, 1], 7 , 1
func (a *UplinkDataStatus) GetPSI14() (pSI14 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI14 Row, sBit, len = [1, 1], 7 , 1
func (a *UplinkDataStatus) SetPSI14(pSI14 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI13 Row, sBit, len = [1, 1], 6 , 1
func (a *UplinkDataStatus) GetPSI13() (pSI13 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI13 Row, sBit, len = [1, 1], 6 , 1
func (a *UplinkDataStatus) SetPSI13(pSI13 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI12 Row, sBit, len = [1, 1], 5 , 1
func (a *UplinkDataStatus) GetPSI12() (pSI12 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI12 Row, sBit, len = [1, 1], 5 , 1
func (a *UplinkDataStatus) SetPSI12(pSI12 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI11 Row, sBit, len = [1, 1], 4 , 1
func (a *UplinkDataStatus) GetPSI11() (pSI11 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI11 Row, sBit, len = [1, 1], 4 , 1
func (a *UplinkDataStatus) SetPSI11(pSI11 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI10 Row, sBit, len = [1, 1], 3 , 1
func (a *UplinkDataStatus) GetPSI10() (pSI10 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI10 Row, sBit, len = [1, 1], 3 , 1
func (a *UplinkDataStatus) SetPSI10(pSI10 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI9 Row, sBit, len = [1, 1], 2 , 1
func (a *UplinkDataStatus) GetPSI9() (pSI9 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI9 Row, sBit, len = [1, 1], 2 , 1
func (a *UplinkDataStatus) SetPSI9(pSI9 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI8 Row, sBit, len = [1, 1], 1 , 1
func (a *UplinkDataStatus) GetPSI8() (pSI8 uint8) {}

// UplinkDataStatus 9.11.3.57
// PSI8 Row, sBit, len = [1, 1], 1 , 1
func (a *UplinkDataStatus) SetPSI8(pSI8 uint8) {}

// UplinkDataStatus 9.11.3.57
// Spare Row, sBit, len = [2, 2], 1 , INF
func (a *UplinkDataStatus) GetSpare() (spare []uint8) {}

// UplinkDataStatus 9.11.3.57
// Spare Row, sBit, len = [2, 2], 1 , INF
func (a *UplinkDataStatus) SetSpare(spare []uint8) {}
