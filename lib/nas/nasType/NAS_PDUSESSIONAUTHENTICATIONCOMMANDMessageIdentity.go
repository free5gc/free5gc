//go:binary-only-package

package nasType

// PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity() (pDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity *PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity) {}

// PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity) SetMessageType(messageType uint8) {}
