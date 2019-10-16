//go:binary-only-package

package nasType

// ABBA 9.11.3.10
// ABBAContents Row, sBit, len = [0, 0], 8 , INF
type ABBA struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewABBA(iei uint8) (aBBA *ABBA) {}

// ABBA 9.11.3.10
// Iei Row, sBit, len = [], 8, 8
func (a *ABBA) GetIei() (iei uint8) {}

// ABBA 9.11.3.10
// Iei Row, sBit, len = [], 8, 8
func (a *ABBA) SetIei(iei uint8) {}

// ABBA 9.11.3.10
// Len Row, sBit, len = [], 8, 8
func (a *ABBA) GetLen() (len uint8) {}

// ABBA 9.11.3.10
// Len Row, sBit, len = [], 8, 8
func (a *ABBA) SetLen(len uint8) {}

// ABBA 9.11.3.10
// ABBAContents Row, sBit, len = [0, 0], 8 , INF
func (a *ABBA) GetABBAContents() (aBBAContents []uint8) {}

// ABBA 9.11.3.10
// ABBAContents Row, sBit, len = [0, 0], 8 , INF
func (a *ABBA) SetABBAContents(aBBAContents []uint8) {}
