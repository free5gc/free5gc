//go:binary-only-package

package nasType

// PDUSESSIONRELEASEREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONRELEASEREQUESTMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONRELEASEREQUESTMessageIdentity() (pDUSESSIONRELEASEREQUESTMessageIdentity *PDUSESSIONRELEASEREQUESTMessageIdentity) {}

// PDUSESSIONRELEASEREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONRELEASEREQUESTMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONRELEASEREQUESTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONRELEASEREQUESTMessageIdentity) SetMessageType(messageType uint8) {}
