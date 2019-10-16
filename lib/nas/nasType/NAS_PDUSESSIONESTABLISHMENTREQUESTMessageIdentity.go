//go:binary-only-package

package nasType

// PDUSESSIONESTABLISHMENTREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONESTABLISHMENTREQUESTMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONESTABLISHMENTREQUESTMessageIdentity() (pDUSESSIONESTABLISHMENTREQUESTMessageIdentity *PDUSESSIONESTABLISHMENTREQUESTMessageIdentity) {}

// PDUSESSIONESTABLISHMENTREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONESTABLISHMENTREQUESTMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONESTABLISHMENTREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONESTABLISHMENTREQUESTMessageIdentity) SetMessageType(messageType uint8) {}
