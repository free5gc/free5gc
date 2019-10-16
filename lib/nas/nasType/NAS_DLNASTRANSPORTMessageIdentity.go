//go:binary-only-package

package nasType

// DLNASTRANSPORTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type DLNASTRANSPORTMessageIdentity struct {
	Octet uint8
}

func NewDLNASTRANSPORTMessageIdentity() (dLNASTRANSPORTMessageIdentity *DLNASTRANSPORTMessageIdentity) {}

// DLNASTRANSPORTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *DLNASTRANSPORTMessageIdentity) GetMessageType() (messageType uint8) {}

// DLNASTRANSPORTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *DLNASTRANSPORTMessageIdentity) SetMessageType(messageType uint8) {}
