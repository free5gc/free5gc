//go:binary-only-package

package nasType

// PDUSESSIONMODIFICATIONCOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONMODIFICATIONCOMMANDMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONMODIFICATIONCOMMANDMessageIdentity() (pDUSESSIONMODIFICATIONCOMMANDMessageIdentity *PDUSESSIONMODIFICATIONCOMMANDMessageIdentity) {}

// PDUSESSIONMODIFICATIONCOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONMODIFICATIONCOMMANDMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONMODIFICATIONCOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONMODIFICATIONCOMMANDMessageIdentity) SetMessageType(messageType uint8) {}
