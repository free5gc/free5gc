//go:binary-only-package

package nasType

// ULNASTRANSPORTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type ULNASTRANSPORTMessageIdentity struct {
	Octet uint8
}

func NewULNASTRANSPORTMessageIdentity() (uLNASTRANSPORTMessageIdentity *ULNASTRANSPORTMessageIdentity) {}

// ULNASTRANSPORTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ULNASTRANSPORTMessageIdentity) GetMessageType() (messageType uint8) {}

// ULNASTRANSPORTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ULNASTRANSPORTMessageIdentity) SetMessageType(messageType uint8) {}
