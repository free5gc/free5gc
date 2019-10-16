//go:binary-only-package

package nasType

// LADNIndication 9.11.3.29
// LADNDNNValue Row, sBit, len = [0, 0], 8 , INF
type LADNIndication struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewLADNIndication(iei uint8) (lADNIndication *LADNIndication) {}

// LADNIndication 9.11.3.29
// Iei Row, sBit, len = [], 8, 8
func (a *LADNIndication) GetIei() (iei uint8) {}

// LADNIndication 9.11.3.29
// Iei Row, sBit, len = [], 8, 8
func (a *LADNIndication) SetIei(iei uint8) {}

// LADNIndication 9.11.3.29
// Len Row, sBit, len = [], 8, 16
func (a *LADNIndication) GetLen() (len uint16) {}

// LADNIndication 9.11.3.29
// Len Row, sBit, len = [], 8, 16
func (a *LADNIndication) SetLen(len uint16) {}

// LADNIndication 9.11.3.29
// LADNDNNValue Row, sBit, len = [0, 0], 8 , INF
func (a *LADNIndication) GetLADNDNNValue() (lADNDNNValue []uint8) {}

// LADNIndication 9.11.3.29
// LADNDNNValue Row, sBit, len = [0, 0], 8 , INF
func (a *LADNIndication) SetLADNDNNValue(lADNDNNValue []uint8) {}
