//go:binary-only-package

package nasType

// EAPMessage 9.11.2.2
// EAPMessage Row, sBit, len = [0, 0], 8 , INF
type EAPMessage struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewEAPMessage(iei uint8) (eAPMessage *EAPMessage) {}

// EAPMessage 9.11.2.2
// Iei Row, sBit, len = [], 8, 8
func (a *EAPMessage) GetIei() (iei uint8) {}

// EAPMessage 9.11.2.2
// Iei Row, sBit, len = [], 8, 8
func (a *EAPMessage) SetIei(iei uint8) {}

// EAPMessage 9.11.2.2
// Len Row, sBit, len = [], 8, 16
func (a *EAPMessage) GetLen() (len uint16) {}

// EAPMessage 9.11.2.2
// Len Row, sBit, len = [], 8, 16
func (a *EAPMessage) SetLen(len uint16) {}

// EAPMessage 9.11.2.2
// EAPMessage Row, sBit, len = [0, 0], 8 , INF
func (a *EAPMessage) GetEAPMessage() (eAPMessage []uint8) {}

// EAPMessage 9.11.2.2
// EAPMessage Row, sBit, len = [0, 0], 8 , INF
func (a *EAPMessage) SetEAPMessage(eAPMessage []uint8) {}
