//go:binary-only-package

package nasType

// PDUSESSIONESTABLISHMENTREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONESTABLISHMENTREJECTMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONESTABLISHMENTREJECTMessageIdentity() (pDUSESSIONESTABLISHMENTREJECTMessageIdentity *PDUSESSIONESTABLISHMENTREJECTMessageIdentity) {}

// PDUSESSIONESTABLISHMENTREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONESTABLISHMENTREJECTMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONESTABLISHMENTREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONESTABLISHMENTREJECTMessageIdentity) SetMessageType(messageType uint8) {}
