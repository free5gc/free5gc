//go:binary-only-package

package nasType

// TMSI5GS 9.11.3.4
// Spare Row, sBit, len = [0, 0], 4 , 1
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
// AMFSetID Row, sBit, len = [1, 1], 8 , 8
// AMFSetID Row, sBit, len = [1, 2], 8 , 10
// AMFPointer Row, sBit, len = [2, 2], 6 , 6
// TMSI5G Row, sBit, len = [3, 6], 8 , 32
type TMSI5GS struct {
	Iei   uint8
	Len   uint16
	Octet [7]uint8
}

func NewTMSI5GS(iei uint8) (tMSI5GS *TMSI5GS) {}

// TMSI5GS 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *TMSI5GS) GetIei() (iei uint8) {}

// TMSI5GS 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *TMSI5GS) SetIei(iei uint8) {}

// TMSI5GS 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *TMSI5GS) GetLen() (len uint16) {}

// TMSI5GS 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *TMSI5GS) SetLen(len uint16) {}

// TMSI5GS 9.11.3.4
// Spare Row, sBit, len = [0, 0], 4 , 1
func (a *TMSI5GS) GetSpare() (spare uint8) {}

// TMSI5GS 9.11.3.4
// Spare Row, sBit, len = [0, 0], 4 , 1
func (a *TMSI5GS) SetSpare(spare uint8) {}

// TMSI5GS 9.11.3.4
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *TMSI5GS) GetTypeOfIdentity() (typeOfIdentity uint8) {}

// TMSI5GS 9.11.3.4
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *TMSI5GS) SetTypeOfIdentity(typeOfIdentity uint8) {}

// TMSI5GS 9.11.3.4
// AMFSetID Row, sBit, len = [1, 2], 8 , 10
func (a *TMSI5GS) GetAMFSetID() (aMFSetID uint16) {}

// TMSI5GS 9.11.3.4
// AMFSetID Row, sBit, len = [1, 2], 8 , 10
func (a *TMSI5GS) SetAMFSetID(aMFSetID uint16) {}

// TMSI5GS 9.11.3.4
// AMFPointer Row, sBit, len = [2, 2], 6 , 6
func (a *TMSI5GS) GetAMFPointer() (aMFPointer uint8) {}

// TMSI5GS 9.11.3.4
// AMFPointer Row, sBit, len = [2, 2], 6 , 6
func (a *TMSI5GS) SetAMFPointer(aMFPointer uint8) {}

// TMSI5GS 9.11.3.4
// TMSI5G Row, sBit, len = [3, 6], 8 , 32
func (a *TMSI5GS) GetTMSI5G() (tMSI5G [4]uint8) {}

// TMSI5GS 9.11.3.4
// TMSI5G Row, sBit, len = [3, 6], 8 , 32
func (a *TMSI5GS) SetTMSI5G(tMSI5G [4]uint8) {}
