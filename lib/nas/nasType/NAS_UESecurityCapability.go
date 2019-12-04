//go:binary-only-package

package nasType

// UESecurityCapability 9.11.3.54
// EA0_5G Row, sBit, len = [0, 0], 8 , 1
// EA1_128_5G Row, sBit, len = [0, 0], 7 , 1
// EA2_128_5G Row, sBit, len = [0, 0], 6 , 1
// EA3_128_5G Row, sBit, len = [0, 0], 5 , 1
// EA4_5G Row, sBit, len = [0, 0], 4 , 1
// EA5_5G Row, sBit, len = [0, 0], 3 , 1
// EA6_5G Row, sBit, len = [0, 0], 2 , 1
// EA7_5G Row, sBit, len = [0, 0], 1 , 1
// IA0_5G Row, sBit, len = [1, 1], 8 , 1
// IA1_128_5G Row, sBit, len = [1, 1], 7 , 1
// IA2_128_5G Row, sBit, len = [1, 1], 6 , 1
// IA3_128_5G Row, sBit, len = [1, 1], 5 , 1
// IA4_5G Row, sBit, len = [1, 1], 4 , 1
// IA5_5G Row, sBit, len = [1, 1], 3 , 1
// IA6_5G Row, sBit, len = [1, 1], 2 , 1
// IA7_5G Row, sBit, len = [1, 1], 1 , 1
// EEA0 Row, sBit, len = [2, 2], 8 , 1
// EEA1_128 Row, sBit, len = [2, 2], 7 , 1
// EEA2_128 Row, sBit, len = [2, 2], 6 , 1
// EEA3_128 Row, sBit, len = [2, 2], 5 , 1
// EEA4 Row, sBit, len = [2, 2], 4 , 1
// EEA5 Row, sBit, len = [2, 2], 3 , 1
// EEA6 Row, sBit, len = [2, 2], 2 , 1
// EEA7 Row, sBit, len = [2, 2], 1 , 1
// EIA0 Row, sBit, len = [3, 3], 8 , 1
// EIA1_128 Row, sBit, len = [3, 3], 7 , 1
// EIA2_128 Row, sBit, len = [3, 3], 6 , 1
// EIA3_128 Row, sBit, len = [3, 3], 5 , 1
// EIA4 Row, sBit, len = [3, 3], 4 , 1
// EIA5 Row, sBit, len = [3, 3], 3 , 1
// EIA6 Row, sBit, len = [3, 3], 2 , 1
// EIA7 Row, sBit, len = [3, 3], 1 , 1
// Spare Row, sBit, len = [4, 7], 8 , 32
type UESecurityCapability struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewUESecurityCapability(iei uint8) (uESecurityCapability *UESecurityCapability) {}

// UESecurityCapability 9.11.3.54
// Iei Row, sBit, len = [], 8, 8
func (a *UESecurityCapability) GetIei() (iei uint8) {}

// UESecurityCapability 9.11.3.54
// Iei Row, sBit, len = [], 8, 8
func (a *UESecurityCapability) SetIei(iei uint8) {}

// UESecurityCapability 9.11.3.54
// Len Row, sBit, len = [], 8, 8
func (a *UESecurityCapability) GetLen() (len uint8) {}

// UESecurityCapability 9.11.3.54
// Len Row, sBit, len = [], 8, 8
func (a *UESecurityCapability) SetLen(len uint8) {}

// UESecurityCapability 9.11.3.54
// EA0_5G Row, sBit, len = [0, 0], 8 , 1
func (a *UESecurityCapability) GetEA0_5G() (eA0_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA0_5G Row, sBit, len = [0, 0], 8 , 1
func (a *UESecurityCapability) SetEA0_5G(eA0_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA1_128_5G Row, sBit, len = [0, 0], 7 , 1
func (a *UESecurityCapability) GetEA1_128_5G() (eA1_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA1_128_5G Row, sBit, len = [0, 0], 7 , 1
func (a *UESecurityCapability) SetEA1_128_5G(eA1_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA2_128_5G Row, sBit, len = [0, 0], 6 , 1
func (a *UESecurityCapability) GetEA2_128_5G() (eA2_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA2_128_5G Row, sBit, len = [0, 0], 6 , 1
func (a *UESecurityCapability) SetEA2_128_5G(eA2_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA3_128_5G Row, sBit, len = [0, 0], 5 , 1
func (a *UESecurityCapability) GetEA3_128_5G() (eA3_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA3_128_5G Row, sBit, len = [0, 0], 5 , 1
func (a *UESecurityCapability) SetEA3_128_5G(eA3_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA4_5G Row, sBit, len = [0, 0], 4 , 1
func (a *UESecurityCapability) GetEA4_5G() (eA4_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA4_5G Row, sBit, len = [0, 0], 4 , 1
func (a *UESecurityCapability) SetEA4_5G(eA4_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA5_5G Row, sBit, len = [0, 0], 3 , 1
func (a *UESecurityCapability) GetEA5_5G() (eA5_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA5_5G Row, sBit, len = [0, 0], 3 , 1
func (a *UESecurityCapability) SetEA5_5G(eA5_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA6_5G Row, sBit, len = [0, 0], 2 , 1
func (a *UESecurityCapability) GetEA6_5G() (eA6_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA6_5G Row, sBit, len = [0, 0], 2 , 1
func (a *UESecurityCapability) SetEA6_5G(eA6_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA7_5G Row, sBit, len = [0, 0], 1 , 1
func (a *UESecurityCapability) GetEA7_5G() (eA7_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EA7_5G Row, sBit, len = [0, 0], 1 , 1
func (a *UESecurityCapability) SetEA7_5G(eA7_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA0_5G Row, sBit, len = [1, 1], 8 , 1
func (a *UESecurityCapability) GetIA0_5G() (iA0_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA0_5G Row, sBit, len = [1, 1], 8 , 1
func (a *UESecurityCapability) SetIA0_5G(iA0_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA1_128_5G Row, sBit, len = [1, 1], 7 , 1
func (a *UESecurityCapability) GetIA1_128_5G() (iA1_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA1_128_5G Row, sBit, len = [1, 1], 7 , 1
func (a *UESecurityCapability) SetIA1_128_5G(iA1_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA2_128_5G Row, sBit, len = [1, 1], 6 , 1
func (a *UESecurityCapability) GetIA2_128_5G() (iA2_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA2_128_5G Row, sBit, len = [1, 1], 6 , 1
func (a *UESecurityCapability) SetIA2_128_5G(iA2_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA3_128_5G Row, sBit, len = [1, 1], 5 , 1
func (a *UESecurityCapability) GetIA3_128_5G() (iA3_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA3_128_5G Row, sBit, len = [1, 1], 5 , 1
func (a *UESecurityCapability) SetIA3_128_5G(iA3_128_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA4_5G Row, sBit, len = [1, 1], 4 , 1
func (a *UESecurityCapability) GetIA4_5G() (iA4_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA4_5G Row, sBit, len = [1, 1], 4 , 1
func (a *UESecurityCapability) SetIA4_5G(iA4_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA5_5G Row, sBit, len = [1, 1], 3 , 1
func (a *UESecurityCapability) GetIA5_5G() (iA5_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA5_5G Row, sBit, len = [1, 1], 3 , 1
func (a *UESecurityCapability) SetIA5_5G(iA5_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA6_5G Row, sBit, len = [1, 1], 2 , 1
func (a *UESecurityCapability) GetIA6_5G() (iA6_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA6_5G Row, sBit, len = [1, 1], 2 , 1
func (a *UESecurityCapability) SetIA6_5G(iA6_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA7_5G Row, sBit, len = [1, 1], 1 , 1
func (a *UESecurityCapability) GetIA7_5G() (iA7_5G uint8) {}

// UESecurityCapability 9.11.3.54
// IA7_5G Row, sBit, len = [1, 1], 1 , 1
func (a *UESecurityCapability) SetIA7_5G(iA7_5G uint8) {}

// UESecurityCapability 9.11.3.54
// EEA0 Row, sBit, len = [2, 2], 8 , 1
func (a *UESecurityCapability) GetEEA0() (eEA0 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA0 Row, sBit, len = [2, 2], 8 , 1
func (a *UESecurityCapability) SetEEA0(eEA0 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA1_128 Row, sBit, len = [2, 2], 7 , 1
func (a *UESecurityCapability) GetEEA1_128() (eEA1_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA1_128 Row, sBit, len = [2, 2], 7 , 1
func (a *UESecurityCapability) SetEEA1_128(eEA1_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA2_128 Row, sBit, len = [2, 2], 6 , 1
func (a *UESecurityCapability) GetEEA2_128() (eEA2_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA2_128 Row, sBit, len = [2, 2], 6 , 1
func (a *UESecurityCapability) SetEEA2_128(eEA2_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA3_128 Row, sBit, len = [2, 2], 5 , 1
func (a *UESecurityCapability) GetEEA3_128() (eEA3_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA3_128 Row, sBit, len = [2, 2], 5 , 1
func (a *UESecurityCapability) SetEEA3_128(eEA3_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA4 Row, sBit, len = [2, 2], 4 , 1
func (a *UESecurityCapability) GetEEA4() (eEA4 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA4 Row, sBit, len = [2, 2], 4 , 1
func (a *UESecurityCapability) SetEEA4(eEA4 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA5 Row, sBit, len = [2, 2], 3 , 1
func (a *UESecurityCapability) GetEEA5() (eEA5 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA5 Row, sBit, len = [2, 2], 3 , 1
func (a *UESecurityCapability) SetEEA5(eEA5 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA6 Row, sBit, len = [2, 2], 2 , 1
func (a *UESecurityCapability) GetEEA6() (eEA6 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA6 Row, sBit, len = [2, 2], 2 , 1
func (a *UESecurityCapability) SetEEA6(eEA6 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA7 Row, sBit, len = [2, 2], 1 , 1
func (a *UESecurityCapability) GetEEA7() (eEA7 uint8) {}

// UESecurityCapability 9.11.3.54
// EEA7 Row, sBit, len = [2, 2], 1 , 1
func (a *UESecurityCapability) SetEEA7(eEA7 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA0 Row, sBit, len = [3, 3], 8 , 1
func (a *UESecurityCapability) GetEIA0() (eIA0 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA0 Row, sBit, len = [3, 3], 8 , 1
func (a *UESecurityCapability) SetEIA0(eIA0 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA1_128 Row, sBit, len = [3, 3], 7 , 1
func (a *UESecurityCapability) GetEIA1_128() (eIA1_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA1_128 Row, sBit, len = [3, 3], 7 , 1
func (a *UESecurityCapability) SetEIA1_128(eIA1_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA2_128 Row, sBit, len = [3, 3], 6 , 1
func (a *UESecurityCapability) GetEIA2_128() (eIA2_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA2_128 Row, sBit, len = [3, 3], 6 , 1
func (a *UESecurityCapability) SetEIA2_128(eIA2_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA3_128 Row, sBit, len = [3, 3], 5 , 1
func (a *UESecurityCapability) GetEIA3_128() (eIA3_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA3_128 Row, sBit, len = [3, 3], 5 , 1
func (a *UESecurityCapability) SetEIA3_128(eIA3_128 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA4 Row, sBit, len = [3, 3], 4 , 1
func (a *UESecurityCapability) GetEIA4() (eIA4 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA4 Row, sBit, len = [3, 3], 4 , 1
func (a *UESecurityCapability) SetEIA4(eIA4 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA5 Row, sBit, len = [3, 3], 3 , 1
func (a *UESecurityCapability) GetEIA5() (eIA5 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA5 Row, sBit, len = [3, 3], 3 , 1
func (a *UESecurityCapability) SetEIA5(eIA5 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA6 Row, sBit, len = [3, 3], 2 , 1
func (a *UESecurityCapability) GetEIA6() (eIA6 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA6 Row, sBit, len = [3, 3], 2 , 1
func (a *UESecurityCapability) SetEIA6(eIA6 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA7 Row, sBit, len = [3, 3], 1 , 1
func (a *UESecurityCapability) GetEIA7() (eIA7 uint8) {}

// UESecurityCapability 9.11.3.54
// EIA7 Row, sBit, len = [3, 3], 1 , 1
func (a *UESecurityCapability) SetEIA7(eIA7 uint8) {}

// UESecurityCapability 9.11.3.54
// Spare Row, sBit, len = [4, 7], 8 , 32
func (a *UESecurityCapability) GetSpare() (spare [4]uint8) {}

// UESecurityCapability 9.11.3.54
// Spare Row, sBit, len = [4, 7], 8 , 32
func (a *UESecurityCapability) SetSpare(spare [4]uint8) {}
