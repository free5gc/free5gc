//go:binary-only-package

package nasType

// NASMessageContainer 9.11.3.33
// NASMessageContainerContents Row, sBit, len = [0, 0], 8 , INF
type NASMessageContainer struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewNASMessageContainer(iei uint8) (nASMessageContainer *NASMessageContainer) {}

// NASMessageContainer 9.11.3.33
// Iei Row, sBit, len = [], 8, 8
func (a *NASMessageContainer) GetIei() (iei uint8) {}

// NASMessageContainer 9.11.3.33
// Iei Row, sBit, len = [], 8, 8
func (a *NASMessageContainer) SetIei(iei uint8) {}

// NASMessageContainer 9.11.3.33
// Len Row, sBit, len = [], 8, 16
func (a *NASMessageContainer) GetLen() (len uint16) {}

// NASMessageContainer 9.11.3.33
// Len Row, sBit, len = [], 8, 16
func (a *NASMessageContainer) SetLen(len uint16) {}

// NASMessageContainer 9.11.3.33
// NASMessageContainerContents Row, sBit, len = [0, 0], 8 , INF
func (a *NASMessageContainer) GetNASMessageContainerContents() (nASMessageContainerContents []uint8) {}

// NASMessageContainer 9.11.3.33
// NASMessageContainerContents Row, sBit, len = [0, 0], 8 , INF
func (a *NASMessageContainer) SetNASMessageContainerContents(nASMessageContainerContents []uint8) {}
