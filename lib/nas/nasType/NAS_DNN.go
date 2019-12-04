//go:binary-only-package

package nasType

import (
	"free5gc/lib/util_3gpp"
)

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
type DNN struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewDNN(iei uint8) (dNN *DNN) {}

// DNN 9.11.2.1A
// Iei Row, sBit, len = [], 8, 8
func (a *DNN) GetIei() (iei uint8) {}

// DNN 9.11.2.1A
// Iei Row, sBit, len = [], 8, 8
func (a *DNN) SetIei(iei uint8) {}

// DNN 9.11.2.1A
// Len Row, sBit, len = [], 8, 8
func (a *DNN) GetLen() (len uint8) {}

// DNN 9.11.2.1A
// Len Row, sBit, len = [], 8, 8
func (a *DNN) SetLen(len uint8) {}

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
func (a *DNN) GetDNN() (dNN []uint8) {}

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
func (a *DNN) SetDNN(dNN []uint8) {}
