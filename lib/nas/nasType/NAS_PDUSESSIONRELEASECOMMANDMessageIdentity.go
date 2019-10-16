//go:binary-only-package

package nasType

// PDUSESSIONRELEASECOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONRELEASECOMMANDMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONRELEASECOMMANDMessageIdentity() (pDUSESSIONRELEASECOMMANDMessageIdentity *PDUSESSIONRELEASECOMMANDMessageIdentity) {}

// PDUSESSIONRELEASECOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONRELEASECOMMANDMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONRELEASECOMMANDMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONRELEASECOMMANDMessageIdentity) SetMessageType(messageType uint8) {}
