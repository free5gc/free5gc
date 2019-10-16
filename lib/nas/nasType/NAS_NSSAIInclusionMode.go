//go:binary-only-package

package nasType

// NSSAIInclusionMode 9.11.3.37A
// Iei Row, sBit, len = [0, 0], 8 , 4
// NSSAIInclusionMode Row, sBit, len = [0, 0], 2 , 2
type NSSAIInclusionMode struct {
	Octet uint8
}

func NewNSSAIInclusionMode(iei uint8) (nSSAIInclusionMode *NSSAIInclusionMode) {}

// NSSAIInclusionMode 9.11.3.37A
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *NSSAIInclusionMode) GetIei() (iei uint8) {}

// NSSAIInclusionMode 9.11.3.37A
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *NSSAIInclusionMode) SetIei(iei uint8) {}

// NSSAIInclusionMode 9.11.3.37A
// NSSAIInclusionMode Row, sBit, len = [0, 0], 2 , 2
func (a *NSSAIInclusionMode) GetNSSAIInclusionMode() (nSSAIInclusionMode uint8) {}

// NSSAIInclusionMode 9.11.3.37A
// NSSAIInclusionMode Row, sBit, len = [0, 0], 2 , 2
func (a *NSSAIInclusionMode) SetNSSAIInclusionMode(nSSAIInclusionMode uint8) {}
