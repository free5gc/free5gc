//go:binary-only-package

package nasType

// GUTI5G 9.11.3.4
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
type GUTI5G struct {
	Iei   uint8
	Len   uint16
	Octet [11]uint8
}

func NewGUTI5G(iei uint8) (gUTI5G *GUTI5G) {}

// GUTI5G 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *GUTI5G) GetIei() (iei uint8) {}

// GUTI5G 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *GUTI5G) SetIei(iei uint8) {}

// GUTI5G 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *GUTI5G) GetLen() (len uint16) {}

// GUTI5G 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *GUTI5G) SetLen(len uint16) {}

// GUTI5G 9.11.3.4
// Spare Row, sBit, len = [0, 0], 4 , 1
func (a *GUTI5G) GetSpare() (spare uint8) {}

// GUTI5G 9.11.3.4
// Spare Row, sBit, len = [0, 0], 4 , 1
func (a *GUTI5G) SetSpare(spare uint8) {}

// GUTI5G 9.11.3.4
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *GUTI5G) GetTypeOfIdentity() (typeOfIdentity uint8) {}

// GUTI5G 9.11.3.4
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *GUTI5G) SetTypeOfIdentity(typeOfIdentity uint8) {}

// GUTI5G 9.11.3.4
// MCCDigit2 Row, sBit, len = [1, 1], 8 , 4
func (a *GUTI5G) GetMCCDigit2() (mCCDigit2 uint8) {}

// GUTI5G 9.11.3.4
// MCCDigit2 Row, sBit, len = [1, 1], 8 , 4
func (a *GUTI5G) SetMCCDigit2(mCCDigit2 uint8) {}

// GUTI5G 9.11.3.4
// MCCDigit1 Row, sBit, len = [1, 1], 4 , 4
func (a *GUTI5G) GetMCCDigit1() (mCCDigit1 uint8) {}

// GUTI5G 9.11.3.4
// MCCDigit1 Row, sBit, len = [1, 1], 4 , 4
func (a *GUTI5G) SetMCCDigit1(mCCDigit1 uint8) {}

// GUTI5G 9.11.3.4
// MNCDigit3 Row, sBit, len = [2, 2], 8 , 4
func (a *GUTI5G) GetMNCDigit3() (mNCDigit3 uint8) {}

// GUTI5G 9.11.3.4
// MNCDigit3 Row, sBit, len = [2, 2], 8 , 4
func (a *GUTI5G) SetMNCDigit3(mNCDigit3 uint8) {}

// GUTI5G 9.11.3.4
// MCCDigit3 Row, sBit, len = [2, 2], 4 , 4
func (a *GUTI5G) GetMCCDigit3() (mCCDigit3 uint8) {}

// GUTI5G 9.11.3.4
// MCCDigit3 Row, sBit, len = [2, 2], 4 , 4
func (a *GUTI5G) SetMCCDigit3(mCCDigit3 uint8) {}

// GUTI5G 9.11.3.4
// MNCDigit2 Row, sBit, len = [3, 3], 8 , 4
func (a *GUTI5G) GetMNCDigit2() (mNCDigit2 uint8) {}

// GUTI5G 9.11.3.4
// MNCDigit2 Row, sBit, len = [3, 3], 8 , 4
func (a *GUTI5G) SetMNCDigit2(mNCDigit2 uint8) {}

// GUTI5G 9.11.3.4
// MNCDigit1 Row, sBit, len = [3, 3], 4 , 4
func (a *GUTI5G) GetMNCDigit1() (mNCDigit1 uint8) {}

// GUTI5G 9.11.3.4
// MNCDigit1 Row, sBit, len = [3, 3], 4 , 4
func (a *GUTI5G) SetMNCDigit1(mNCDigit1 uint8) {}

// GUTI5G 9.11.3.4
// AMFRegionID Row, sBit, len = [4, 4], 8 , 8
func (a *GUTI5G) GetAMFRegionID() (aMFRegionID uint8) {}

// GUTI5G 9.11.3.4
// AMFRegionID Row, sBit, len = [4, 4], 8 , 8
func (a *GUTI5G) SetAMFRegionID(aMFRegionID uint8) {}

// GUTI5G 9.11.3.4
// AMFSetID Row, sBit, len = [5, 6], 8 , 10
func (a *GUTI5G) GetAMFSetID() (aMFSetID uint16) {}

// GUTI5G 9.11.3.4
// AMFSetID Row, sBit, len = [5, 6], 8 , 10
func (a *GUTI5G) SetAMFSetID(aMFSetID uint16) {}

// GUTI5G 9.11.3.4
// AMFPointer Row, sBit, len = [6, 6], 6 , 6
func (a *GUTI5G) GetAMFPointer() (aMFPointer uint8) {}

// GUTI5G 9.11.3.4
// AMFPointer Row, sBit, len = [6, 6], 6 , 6
func (a *GUTI5G) SetAMFPointer(aMFPointer uint8) {}

// GUTI5G 9.11.3.4
// TMSI5G Row, sBit, len = [7, 10], 8 , 32
func (a *GUTI5G) GetTMSI5G() (tMSI5G [4]uint8) {}

// GUTI5G 9.11.3.4
// TMSI5G Row, sBit, len = [7, 10], 8 , 32
func (a *GUTI5G) SetTMSI5G(tMSI5G [4]uint8) {}
