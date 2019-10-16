//go:binary-only-package

package nasType

// PDUSESSIONMODIFICATIONREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONMODIFICATIONREQUESTMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONMODIFICATIONREQUESTMessageIdentity() (pDUSESSIONMODIFICATIONREQUESTMessageIdentity *PDUSESSIONMODIFICATIONREQUESTMessageIdentity) {}

// PDUSESSIONMODIFICATIONREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONMODIFICATIONREQUESTMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONMODIFICATIONREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONMODIFICATIONREQUESTMessageIdentity) SetMessageType(messageType uint8) {}
