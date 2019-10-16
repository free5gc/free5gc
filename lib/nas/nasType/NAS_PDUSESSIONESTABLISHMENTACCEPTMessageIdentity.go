//go:binary-only-package

package nasType

// PDUSESSIONESTABLISHMENTACCEPTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONESTABLISHMENTACCEPTMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONESTABLISHMENTACCEPTMessageIdentity() (pDUSESSIONESTABLISHMENTACCEPTMessageIdentity *PDUSESSIONESTABLISHMENTACCEPTMessageIdentity) {}

// PDUSESSIONESTABLISHMENTACCEPTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONESTABLISHMENTACCEPTMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONESTABLISHMENTACCEPTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONESTABLISHMENTACCEPTMessageIdentity) SetMessageType(messageType uint8) {}
