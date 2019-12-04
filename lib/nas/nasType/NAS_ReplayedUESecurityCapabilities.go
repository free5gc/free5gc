//go:binary-only-package

package nasType

// ReplayedUESecurityCapabilities 9.11.3.54
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
type ReplayedUESecurityCapabilities struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewReplayedUESecurityCapabilities(iei uint8) (replayedUESecurityCapabilities *ReplayedUESecurityCapabilities) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// Iei Row, sBit, len = [], 8, 8
func (a *ReplayedUESecurityCapabilities) GetIei() (iei uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// Iei Row, sBit, len = [], 8, 8
func (a *ReplayedUESecurityCapabilities) SetIei(iei uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// Len Row, sBit, len = [], 8, 8
func (a *ReplayedUESecurityCapabilities) GetLen() (len uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// Len Row, sBit, len = [], 8, 8
func (a *ReplayedUESecurityCapabilities) SetLen(len uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA0_5G Row, sBit, len = [0, 0], 8 , 1
func (a *ReplayedUESecurityCapabilities) GetEA0_5G() (eA0_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA0_5G Row, sBit, len = [0, 0], 8 , 1
func (a *ReplayedUESecurityCapabilities) SetEA0_5G(eA0_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA1_128_5G Row, sBit, len = [0, 0], 7 , 1
func (a *ReplayedUESecurityCapabilities) GetEA1_128_5G() (eA1_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA1_128_5G Row, sBit, len = [0, 0], 7 , 1
func (a *ReplayedUESecurityCapabilities) SetEA1_128_5G(eA1_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA2_128_5G Row, sBit, len = [0, 0], 6 , 1
func (a *ReplayedUESecurityCapabilities) GetEA2_128_5G() (eA2_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA2_128_5G Row, sBit, len = [0, 0], 6 , 1
func (a *ReplayedUESecurityCapabilities) SetEA2_128_5G(eA2_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA3_128_5G Row, sBit, len = [0, 0], 5 , 1
func (a *ReplayedUESecurityCapabilities) GetEA3_128_5G() (eA3_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA3_128_5G Row, sBit, len = [0, 0], 5 , 1
func (a *ReplayedUESecurityCapabilities) SetEA3_128_5G(eA3_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA4_5G Row, sBit, len = [0, 0], 4 , 1
func (a *ReplayedUESecurityCapabilities) GetEA4_5G() (eA4_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA4_5G Row, sBit, len = [0, 0], 4 , 1
func (a *ReplayedUESecurityCapabilities) SetEA4_5G(eA4_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA5_5G Row, sBit, len = [0, 0], 3 , 1
func (a *ReplayedUESecurityCapabilities) GetEA5_5G() (eA5_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA5_5G Row, sBit, len = [0, 0], 3 , 1
func (a *ReplayedUESecurityCapabilities) SetEA5_5G(eA5_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA6_5G Row, sBit, len = [0, 0], 2 , 1
func (a *ReplayedUESecurityCapabilities) GetEA6_5G() (eA6_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA6_5G Row, sBit, len = [0, 0], 2 , 1
func (a *ReplayedUESecurityCapabilities) SetEA6_5G(eA6_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA7_5G Row, sBit, len = [0, 0], 1 , 1
func (a *ReplayedUESecurityCapabilities) GetEA7_5G() (eA7_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EA7_5G Row, sBit, len = [0, 0], 1 , 1
func (a *ReplayedUESecurityCapabilities) SetEA7_5G(eA7_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA0_5G Row, sBit, len = [1, 1], 8 , 1
func (a *ReplayedUESecurityCapabilities) GetIA0_5G() (iA0_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA0_5G Row, sBit, len = [1, 1], 8 , 1
func (a *ReplayedUESecurityCapabilities) SetIA0_5G(iA0_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA1_128_5G Row, sBit, len = [1, 1], 7 , 1
func (a *ReplayedUESecurityCapabilities) GetIA1_128_5G() (iA1_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA1_128_5G Row, sBit, len = [1, 1], 7 , 1
func (a *ReplayedUESecurityCapabilities) SetIA1_128_5G(iA1_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA2_128_5G Row, sBit, len = [1, 1], 6 , 1
func (a *ReplayedUESecurityCapabilities) GetIA2_128_5G() (iA2_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA2_128_5G Row, sBit, len = [1, 1], 6 , 1
func (a *ReplayedUESecurityCapabilities) SetIA2_128_5G(iA2_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA3_128_5G Row, sBit, len = [1, 1], 5 , 1
func (a *ReplayedUESecurityCapabilities) GetIA3_128_5G() (iA3_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA3_128_5G Row, sBit, len = [1, 1], 5 , 1
func (a *ReplayedUESecurityCapabilities) SetIA3_128_5G(iA3_128_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA4_5G Row, sBit, len = [1, 1], 4 , 1
func (a *ReplayedUESecurityCapabilities) GetIA4_5G() (iA4_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA4_5G Row, sBit, len = [1, 1], 4 , 1
func (a *ReplayedUESecurityCapabilities) SetIA4_5G(iA4_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA5_5G Row, sBit, len = [1, 1], 3 , 1
func (a *ReplayedUESecurityCapabilities) GetIA5_5G() (iA5_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA5_5G Row, sBit, len = [1, 1], 3 , 1
func (a *ReplayedUESecurityCapabilities) SetIA5_5G(iA5_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA6_5G Row, sBit, len = [1, 1], 2 , 1
func (a *ReplayedUESecurityCapabilities) GetIA6_5G() (iA6_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA6_5G Row, sBit, len = [1, 1], 2 , 1
func (a *ReplayedUESecurityCapabilities) SetIA6_5G(iA6_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA7_5G Row, sBit, len = [1, 1], 1 , 1
func (a *ReplayedUESecurityCapabilities) GetIA7_5G() (iA7_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// IA7_5G Row, sBit, len = [1, 1], 1 , 1
func (a *ReplayedUESecurityCapabilities) SetIA7_5G(iA7_5G uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA0 Row, sBit, len = [2, 2], 8 , 1
func (a *ReplayedUESecurityCapabilities) GetEEA0() (eEA0 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA0 Row, sBit, len = [2, 2], 8 , 1
func (a *ReplayedUESecurityCapabilities) SetEEA0(eEA0 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA1_128 Row, sBit, len = [2, 2], 7 , 1
func (a *ReplayedUESecurityCapabilities) GetEEA1_128() (eEA1_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA1_128 Row, sBit, len = [2, 2], 7 , 1
func (a *ReplayedUESecurityCapabilities) SetEEA1_128(eEA1_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA2_128 Row, sBit, len = [2, 2], 6 , 1
func (a *ReplayedUESecurityCapabilities) GetEEA2_128() (eEA2_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA2_128 Row, sBit, len = [2, 2], 6 , 1
func (a *ReplayedUESecurityCapabilities) SetEEA2_128(eEA2_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA3_128 Row, sBit, len = [2, 2], 5 , 1
func (a *ReplayedUESecurityCapabilities) GetEEA3_128() (eEA3_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA3_128 Row, sBit, len = [2, 2], 5 , 1
func (a *ReplayedUESecurityCapabilities) SetEEA3_128(eEA3_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA4 Row, sBit, len = [2, 2], 4 , 1
func (a *ReplayedUESecurityCapabilities) GetEEA4() (eEA4 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA4 Row, sBit, len = [2, 2], 4 , 1
func (a *ReplayedUESecurityCapabilities) SetEEA4(eEA4 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA5 Row, sBit, len = [2, 2], 3 , 1
func (a *ReplayedUESecurityCapabilities) GetEEA5() (eEA5 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA5 Row, sBit, len = [2, 2], 3 , 1
func (a *ReplayedUESecurityCapabilities) SetEEA5(eEA5 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA6 Row, sBit, len = [2, 2], 2 , 1
func (a *ReplayedUESecurityCapabilities) GetEEA6() (eEA6 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA6 Row, sBit, len = [2, 2], 2 , 1
func (a *ReplayedUESecurityCapabilities) SetEEA6(eEA6 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA7 Row, sBit, len = [2, 2], 1 , 1
func (a *ReplayedUESecurityCapabilities) GetEEA7() (eEA7 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EEA7 Row, sBit, len = [2, 2], 1 , 1
func (a *ReplayedUESecurityCapabilities) SetEEA7(eEA7 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA0 Row, sBit, len = [3, 3], 8 , 1
func (a *ReplayedUESecurityCapabilities) GetEIA0() (eIA0 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA0 Row, sBit, len = [3, 3], 8 , 1
func (a *ReplayedUESecurityCapabilities) SetEIA0(eIA0 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA1_128 Row, sBit, len = [3, 3], 7 , 1
func (a *ReplayedUESecurityCapabilities) GetEIA1_128() (eIA1_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA1_128 Row, sBit, len = [3, 3], 7 , 1
func (a *ReplayedUESecurityCapabilities) SetEIA1_128(eIA1_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA2_128 Row, sBit, len = [3, 3], 6 , 1
func (a *ReplayedUESecurityCapabilities) GetEIA2_128() (eIA2_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA2_128 Row, sBit, len = [3, 3], 6 , 1
func (a *ReplayedUESecurityCapabilities) SetEIA2_128(eIA2_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA3_128 Row, sBit, len = [3, 3], 5 , 1
func (a *ReplayedUESecurityCapabilities) GetEIA3_128() (eIA3_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA3_128 Row, sBit, len = [3, 3], 5 , 1
func (a *ReplayedUESecurityCapabilities) SetEIA3_128(eIA3_128 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA4 Row, sBit, len = [3, 3], 4 , 1
func (a *ReplayedUESecurityCapabilities) GetEIA4() (eIA4 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA4 Row, sBit, len = [3, 3], 4 , 1
func (a *ReplayedUESecurityCapabilities) SetEIA4(eIA4 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA5 Row, sBit, len = [3, 3], 3 , 1
func (a *ReplayedUESecurityCapabilities) GetEIA5() (eIA5 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA5 Row, sBit, len = [3, 3], 3 , 1
func (a *ReplayedUESecurityCapabilities) SetEIA5(eIA5 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA6 Row, sBit, len = [3, 3], 2 , 1
func (a *ReplayedUESecurityCapabilities) GetEIA6() (eIA6 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA6 Row, sBit, len = [3, 3], 2 , 1
func (a *ReplayedUESecurityCapabilities) SetEIA6(eIA6 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA7 Row, sBit, len = [3, 3], 1 , 1
func (a *ReplayedUESecurityCapabilities) GetEIA7() (eIA7 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// EIA7 Row, sBit, len = [3, 3], 1 , 1
func (a *ReplayedUESecurityCapabilities) SetEIA7(eIA7 uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// Spare Row, sBit, len = [4, 7], 8 , 32
func (a *ReplayedUESecurityCapabilities) GetSpare() (spare [4]uint8) {}

// ReplayedUESecurityCapabilities 9.11.3.54
// Spare Row, sBit, len = [4, 7], 8 , 32
func (a *ReplayedUESecurityCapabilities) SetSpare(spare [4]uint8) {}
