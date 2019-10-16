//go:binary-only-package

package nasType

// PDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity() (pDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity *PDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity) {}

// PDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity) SetMessageType(messageType uint8) {}
