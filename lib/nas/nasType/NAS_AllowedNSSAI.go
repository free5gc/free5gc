//go:binary-only-package

package nasType

// AllowedNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
type AllowedNSSAI struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewAllowedNSSAI(iei uint8) (allowedNSSAI *AllowedNSSAI) {}

// AllowedNSSAI 9.11.3.37
// Iei Row, sBit, len = [], 8, 8
func (a *AllowedNSSAI) GetIei() (iei uint8) {}

// AllowedNSSAI 9.11.3.37
// Iei Row, sBit, len = [], 8, 8
func (a *AllowedNSSAI) SetIei(iei uint8) {}

// AllowedNSSAI 9.11.3.37
// Len Row, sBit, len = [], 8, 8
func (a *AllowedNSSAI) GetLen() (len uint8) {}

// AllowedNSSAI 9.11.3.37
// Len Row, sBit, len = [], 8, 8
func (a *AllowedNSSAI) SetLen(len uint8) {}

// AllowedNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
func (a *AllowedNSSAI) GetSNSSAIValue() (sNSSAIValue []uint8) {}

// AllowedNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
func (a *AllowedNSSAI) SetSNSSAIValue(sNSSAIValue []uint8) {}
