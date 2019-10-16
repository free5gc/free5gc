//go:binary-only-package

package nasType

// PayloadContainer 9.11.3.39
// PayloadContainerContents Row, sBit, len = [0, 0], 8 , INF
type PayloadContainer struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewPayloadContainer(iei uint8) (payloadContainer *PayloadContainer) {}

// PayloadContainer 9.11.3.39
// Iei Row, sBit, len = [], 8, 8
func (a *PayloadContainer) GetIei() (iei uint8) {}

// PayloadContainer 9.11.3.39
// Iei Row, sBit, len = [], 8, 8
func (a *PayloadContainer) SetIei(iei uint8) {}

// PayloadContainer 9.11.3.39
// Len Row, sBit, len = [], 8, 16
func (a *PayloadContainer) GetLen() (len uint16) {}

// PayloadContainer 9.11.3.39
// Len Row, sBit, len = [], 8, 16
func (a *PayloadContainer) SetLen(len uint16) {}

// PayloadContainer 9.11.3.39
// PayloadContainerContents Row, sBit, len = [0, 0], 8 , INF
func (a *PayloadContainer) GetPayloadContainerContents() (payloadContainerContents []uint8) {}

// PayloadContainer 9.11.3.39
// PayloadContainerContents Row, sBit, len = [0, 0], 8 , INF
func (a *PayloadContainer) SetPayloadContainerContents(payloadContainerContents []uint8) {}
