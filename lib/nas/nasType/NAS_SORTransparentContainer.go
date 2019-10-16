//go:binary-only-package

package nasType

// SORTransparentContainer 9.11.3.51
// SORContent Row, sBit, len = [0, 0], 8 , INF
type SORTransparentContainer struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewSORTransparentContainer(iei uint8) (sORTransparentContainer *SORTransparentContainer) {}

// SORTransparentContainer 9.11.3.51
// Iei Row, sBit, len = [], 8, 8
func (a *SORTransparentContainer) GetIei() (iei uint8) {}

// SORTransparentContainer 9.11.3.51
// Iei Row, sBit, len = [], 8, 8
func (a *SORTransparentContainer) SetIei(iei uint8) {}

// SORTransparentContainer 9.11.3.51
// Len Row, sBit, len = [], 8, 16
func (a *SORTransparentContainer) GetLen() (len uint16) {}

// SORTransparentContainer 9.11.3.51
// Len Row, sBit, len = [], 8, 16
func (a *SORTransparentContainer) SetLen(len uint16) {}

// SORTransparentContainer 9.11.3.51
// SORContent Row, sBit, len = [0, 0], 8 , INF
func (a *SORTransparentContainer) GetSORContent() (sORContent []uint8) {}

// SORTransparentContainer 9.11.3.51
// SORContent Row, sBit, len = [0, 0], 8 , INF
func (a *SORTransparentContainer) SetSORContent(sORContent []uint8) {}
