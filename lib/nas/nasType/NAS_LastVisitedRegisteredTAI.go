//go:binary-only-package

package nasType

// LastVisitedRegisteredTAI 9.11.3.8
// MCCDigit2 Row, sBit, len = [0, 0], 8 , 4
// MCCDigit1 Row, sBit, len = [0, 0], 4 , 4
// MNCDigit3 Row, sBit, len = [1, 1], 8 , 4
// MCCDigit3 Row, sBit, len = [1, 1], 4 , 4
// MNCDigit2 Row, sBit, len = [2, 2], 8 , 4
// MNCDigit1 Row, sBit, len = [2, 2], 4 , 4
// TAC Row, sBit, len = [3, 5], 8 , 24
type LastVisitedRegisteredTAI struct {
	Iei   uint8
	Octet [7]uint8
}

func NewLastVisitedRegisteredTAI(iei uint8) (lastVisitedRegisteredTAI *LastVisitedRegisteredTAI) {}

// LastVisitedRegisteredTAI 9.11.3.8
// Iei Row, sBit, len = [], 8, 8
func (a *LastVisitedRegisteredTAI) GetIei() (iei uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// Iei Row, sBit, len = [], 8, 8
func (a *LastVisitedRegisteredTAI) SetIei(iei uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MCCDigit2 Row, sBit, len = [0, 0], 8 , 4
func (a *LastVisitedRegisteredTAI) GetMCCDigit2() (mCCDigit2 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MCCDigit2 Row, sBit, len = [0, 0], 8 , 4
func (a *LastVisitedRegisteredTAI) SetMCCDigit2(mCCDigit2 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MCCDigit1 Row, sBit, len = [0, 0], 4 , 4
func (a *LastVisitedRegisteredTAI) GetMCCDigit1() (mCCDigit1 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MCCDigit1 Row, sBit, len = [0, 0], 4 , 4
func (a *LastVisitedRegisteredTAI) SetMCCDigit1(mCCDigit1 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MNCDigit3 Row, sBit, len = [1, 1], 8 , 4
func (a *LastVisitedRegisteredTAI) GetMNCDigit3() (mNCDigit3 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MNCDigit3 Row, sBit, len = [1, 1], 8 , 4
func (a *LastVisitedRegisteredTAI) SetMNCDigit3(mNCDigit3 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MCCDigit3 Row, sBit, len = [1, 1], 4 , 4
func (a *LastVisitedRegisteredTAI) GetMCCDigit3() (mCCDigit3 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MCCDigit3 Row, sBit, len = [1, 1], 4 , 4
func (a *LastVisitedRegisteredTAI) SetMCCDigit3(mCCDigit3 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MNCDigit2 Row, sBit, len = [2, 2], 8 , 4
func (a *LastVisitedRegisteredTAI) GetMNCDigit2() (mNCDigit2 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MNCDigit2 Row, sBit, len = [2, 2], 8 , 4
func (a *LastVisitedRegisteredTAI) SetMNCDigit2(mNCDigit2 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MNCDigit1 Row, sBit, len = [2, 2], 4 , 4
func (a *LastVisitedRegisteredTAI) GetMNCDigit1() (mNCDigit1 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// MNCDigit1 Row, sBit, len = [2, 2], 4 , 4
func (a *LastVisitedRegisteredTAI) SetMNCDigit1(mNCDigit1 uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// TAC Row, sBit, len = [3, 5], 8 , 24
func (a *LastVisitedRegisteredTAI) GetTAC() (tAC [3]uint8) {}

// LastVisitedRegisteredTAI 9.11.3.8
// TAC Row, sBit, len = [3, 5], 8 , 24
func (a *LastVisitedRegisteredTAI) SetTAC(tAC [3]uint8) {}
