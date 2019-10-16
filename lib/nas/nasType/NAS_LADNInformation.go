//go:binary-only-package

package nasType

// LADNInformation 9.11.3.30
// LADND Row, sBit, len = [0, 0], 8 , INF
type LADNInformation struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewLADNInformation(iei uint8) (lADNInformation *LADNInformation) {}

// LADNInformation 9.11.3.30
// Iei Row, sBit, len = [], 8, 8
func (a *LADNInformation) GetIei() (iei uint8) {}

// LADNInformation 9.11.3.30
// Iei Row, sBit, len = [], 8, 8
func (a *LADNInformation) SetIei(iei uint8) {}

// LADNInformation 9.11.3.30
// Len Row, sBit, len = [], 8, 16
func (a *LADNInformation) GetLen() (len uint16) {}

// LADNInformation 9.11.3.30
// Len Row, sBit, len = [], 8, 16
func (a *LADNInformation) SetLen(len uint16) {}

// LADNInformation 9.11.3.30
// LADND Row, sBit, len = [0, 0], 8 , INF
func (a *LADNInformation) GetLADND() (lADND []uint8) {}

// LADNInformation 9.11.3.30
// LADND Row, sBit, len = [0, 0], 8 , INF
func (a *LADNInformation) SetLADND(lADND []uint8) {}
