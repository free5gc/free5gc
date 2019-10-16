//go:binary-only-package

package nasType

// AdditionalGUTI 9.11.3.4
// Spare Row, sBit, len = [0, 0], 4 , 1
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
// MCCDigit2 Row, sBit, len = [1, 1], 8 , 4
// MCCDigit1 Row, sBit, len = [1, 1], 4 , 4
// MNCDigit3 Row, sBit, len = [2, 2], 8 , 4
// MCCDigit3 Row, sBit, len = [2, 2], 4 , 4
// MNCDigit2 Row, sBit, len = [3, 3], 8 , 4
// MNCDigit1 Row, sBit, len = [3, 3], 4 , 4
// AMFRegionID Row, sBit, len = [4, 4], 8 , 8
// AMFSetID Row, sBit, len = [5, 6], 8 , 10
// AMFPointer Row, sBit, len = [6, 6], 6 , 6
// TMSI5G Row, sBit, len = [7, 10], 8 , 32
type AdditionalGUTI struct {
	Iei   uint8
	Len   uint16
	Octet [11]uint8
}

func NewAdditionalGUTI(iei uint8) (additionalGUTI *AdditionalGUTI) {}

// AdditionalGUTI 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *AdditionalGUTI) GetIei() (iei uint8) {}

// AdditionalGUTI 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *AdditionalGUTI) SetIei(iei uint8) {}

// AdditionalGUTI 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *AdditionalGUTI) GetLen() (len uint16) {}

// AdditionalGUTI 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *AdditionalGUTI) SetLen(len uint16) {}

// AdditionalGUTI 9.11.3.4
// Spare Row, sBit, len = [0, 0], 4 , 1
func (a *AdditionalGUTI) GetSpare() (spare uint8) {}

// AdditionalGUTI 9.11.3.4
// Spare Row, sBit, len = [0, 0], 4 , 1
func (a *AdditionalGUTI) SetSpare(spare uint8) {}

// AdditionalGUTI 9.11.3.4
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *AdditionalGUTI) GetTypeOfIdentity() (typeOfIdentity uint8) {}

// AdditionalGUTI 9.11.3.4
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *AdditionalGUTI) SetTypeOfIdentity(typeOfIdentity uint8) {}

// AdditionalGUTI 9.11.3.4
// MCCDigit2 Row, sBit, len = [1, 1], 8 , 4
func (a *AdditionalGUTI) GetMCCDigit2() (mCCDigit2 uint8) {}

// AdditionalGUTI 9.11.3.4
// MCCDigit2 Row, sBit, len = [1, 1], 8 , 4
func (a *AdditionalGUTI) SetMCCDigit2(mCCDigit2 uint8) {}

// AdditionalGUTI 9.11.3.4
// MCCDigit1 Row, sBit, len = [1, 1], 4 , 4
func (a *AdditionalGUTI) GetMCCDigit1() (mCCDigit1 uint8) {}

// AdditionalGUTI 9.11.3.4
// MCCDigit1 Row, sBit, len = [1, 1], 4 , 4
func (a *AdditionalGUTI) SetMCCDigit1(mCCDigit1 uint8) {}

// AdditionalGUTI 9.11.3.4
// MNCDigit3 Row, sBit, len = [2, 2], 8 , 4
func (a *AdditionalGUTI) GetMNCDigit3() (mNCDigit3 uint8) {}

// AdditionalGUTI 9.11.3.4
// MNCDigit3 Row, sBit, len = [2, 2], 8 , 4
func (a *AdditionalGUTI) SetMNCDigit3(mNCDigit3 uint8) {}

// AdditionalGUTI 9.11.3.4
// MCCDigit3 Row, sBit, len = [2, 2], 4 , 4
func (a *AdditionalGUTI) GetMCCDigit3() (mCCDigit3 uint8) {}

// AdditionalGUTI 9.11.3.4
// MCCDigit3 Row, sBit, len = [2, 2], 4 , 4
func (a *AdditionalGUTI) SetMCCDigit3(mCCDigit3 uint8) {}

// AdditionalGUTI 9.11.3.4
// MNCDigit2 Row, sBit, len = [3, 3], 8 , 4
func (a *AdditionalGUTI) GetMNCDigit2() (mNCDigit2 uint8) {}

// AdditionalGUTI 9.11.3.4
// MNCDigit2 Row, sBit, len = [3, 3], 8 , 4
func (a *AdditionalGUTI) SetMNCDigit2(mNCDigit2 uint8) {}

// AdditionalGUTI 9.11.3.4
// MNCDigit1 Row, sBit, len = [3, 3], 4 , 4
func (a *AdditionalGUTI) GetMNCDigit1() (mNCDigit1 uint8) {}

// AdditionalGUTI 9.11.3.4
// MNCDigit1 Row, sBit, len = [3, 3], 4 , 4
func (a *AdditionalGUTI) SetMNCDigit1(mNCDigit1 uint8) {}

// AdditionalGUTI 9.11.3.4
// AMFRegionID Row, sBit, len = [4, 4], 8 , 8
func (a *AdditionalGUTI) GetAMFRegionID() (aMFRegionID uint8) {}

// AdditionalGUTI 9.11.3.4
// AMFRegionID Row, sBit, len = [4, 4], 8 , 8
func (a *AdditionalGUTI) SetAMFRegionID(aMFRegionID uint8) {}

// AdditionalGUTI 9.11.3.4
// AMFSetID Row, sBit, len = [5, 6], 8 , 10
func (a *AdditionalGUTI) GetAMFSetID() (aMFSetID uint16) {}

// AdditionalGUTI 9.11.3.4
// AMFSetID Row, sBit, len = [5, 6], 8 , 10
func (a *AdditionalGUTI) SetAMFSetID(aMFSetID uint16) {}

// AdditionalGUTI 9.11.3.4
// AMFPointer Row, sBit, len = [6, 6], 6 , 6
func (a *AdditionalGUTI) GetAMFPointer() (aMFPointer uint8) {}

// AdditionalGUTI 9.11.3.4
// AMFPointer Row, sBit, len = [6, 6], 6 , 6
func (a *AdditionalGUTI) SetAMFPointer(aMFPointer uint8) {}

// AdditionalGUTI 9.11.3.4
// TMSI5G Row, sBit, len = [7, 10], 8 , 32
func (a *AdditionalGUTI) GetTMSI5G() (tMSI5G [4]uint8) {}

// AdditionalGUTI 9.11.3.4
// TMSI5G Row, sBit, len = [7, 10], 8 , 32
func (a *AdditionalGUTI) SetTMSI5G(tMSI5G [4]uint8) {}
