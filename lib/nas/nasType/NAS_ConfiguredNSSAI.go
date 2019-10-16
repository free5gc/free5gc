//go:binary-only-package

package nasType

// ConfiguredNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
type ConfiguredNSSAI struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewConfiguredNSSAI(iei uint8) (configuredNSSAI *ConfiguredNSSAI) {}

// ConfiguredNSSAI 9.11.3.37
// Iei Row, sBit, len = [], 8, 8
func (a *ConfiguredNSSAI) GetIei() (iei uint8) {}

// ConfiguredNSSAI 9.11.3.37
// Iei Row, sBit, len = [], 8, 8
func (a *ConfiguredNSSAI) SetIei(iei uint8) {}

// ConfiguredNSSAI 9.11.3.37
// Len Row, sBit, len = [], 8, 8
func (a *ConfiguredNSSAI) GetLen() (len uint8) {}

// ConfiguredNSSAI 9.11.3.37
// Len Row, sBit, len = [], 8, 8
func (a *ConfiguredNSSAI) SetLen(len uint8) {}

// ConfiguredNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
func (a *ConfiguredNSSAI) GetSNSSAIValue() (sNSSAIValue []uint8) {}

// ConfiguredNSSAI 9.11.3.37
// SNSSAIValue Row, sBit, len = [0, 0], 0 , INF
func (a *ConfiguredNSSAI) SetSNSSAIValue(sNSSAIValue []uint8) {}
