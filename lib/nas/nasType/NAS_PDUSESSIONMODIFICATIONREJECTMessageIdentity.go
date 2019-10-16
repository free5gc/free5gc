//go:binary-only-package

package nasType

// PDUSESSIONMODIFICATIONREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONMODIFICATIONREJECTMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONMODIFICATIONREJECTMessageIdentity() (pDUSESSIONMODIFICATIONREJECTMessageIdentity *PDUSESSIONMODIFICATIONREJECTMessageIdentity) {}

// PDUSESSIONMODIFICATIONREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONMODIFICATIONREJECTMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONMODIFICATIONREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONMODIFICATIONREJECTMessageIdentity) SetMessageType(messageType uint8) {}
