//go:binary-only-package

package nasType

// RejectedNSSAI 9.11.3.46
// RejectedNSSAIContents Row, sBit, len = [0, 0], 0 , INF
type RejectedNSSAI struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewRejectedNSSAI(iei uint8) (rejectedNSSAI *RejectedNSSAI) {}

// RejectedNSSAI 9.11.3.46
// Iei Row, sBit, len = [], 8, 8
func (a *RejectedNSSAI) GetIei() (iei uint8) {}

// RejectedNSSAI 9.11.3.46
// Iei Row, sBit, len = [], 8, 8
func (a *RejectedNSSAI) SetIei(iei uint8) {}

// RejectedNSSAI 9.11.3.46
// Len Row, sBit, len = [], 8, 8
func (a *RejectedNSSAI) GetLen() (len uint8) {}

// RejectedNSSAI 9.11.3.46
// Len Row, sBit, len = [], 8, 8
func (a *RejectedNSSAI) SetLen(len uint8) {}

// RejectedNSSAI 9.11.3.46
// RejectedNSSAIContents Row, sBit, len = [0, 0], 0 , INF
func (a *RejectedNSSAI) GetRejectedNSSAIContents() (rejectedNSSAIContents []uint8) {}

// RejectedNSSAI 9.11.3.46
// RejectedNSSAIContents Row, sBit, len = [0, 0], 0 , INF
func (a *RejectedNSSAI) SetRejectedNSSAIContents(rejectedNSSAIContents []uint8) {}
