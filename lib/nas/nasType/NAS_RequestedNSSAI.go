//go:binary-only-package

package nasType

// RequestedNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
type RequestedNSSAI struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewRequestedNSSAI(iei uint8) (requestedNSSAI *RequestedNSSAI) {}

// RequestedNSSAI 9.11.3.37
// Iei Row, sBit, len = [], 8, 8
func (a *RequestedNSSAI) GetIei() (iei uint8) {}

// RequestedNSSAI 9.11.3.37
// Iei Row, sBit, len = [], 8, 8
func (a *RequestedNSSAI) SetIei(iei uint8) {}

// RequestedNSSAI 9.11.3.37
// Len Row, sBit, len = [], 8, 8
func (a *RequestedNSSAI) GetLen() (len uint8) {}

// RequestedNSSAI 9.11.3.37
// Len Row, sBit, len = [], 8, 8
func (a *RequestedNSSAI) SetLen(len uint8) {}

// RequestedNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
func (a *RequestedNSSAI) GetSNSSAIValue() (sNSSAIValue []uint8) {}

// RequestedNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
func (a *RequestedNSSAI) SetSNSSAIValue(sNSSAIValue []uint8) {}
