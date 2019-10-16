//go:binary-only-package

package nasType

// PDUSESSIONRELEASEREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONRELEASEREJECTMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONRELEASEREJECTMessageIdentity() (pDUSESSIONRELEASEREJECTMessageIdentity *PDUSESSIONRELEASEREJECTMessageIdentity) {}

// PDUSESSIONRELEASEREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONRELEASEREJECTMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONRELEASEREJECTMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONRELEASEREJECTMessageIdentity) SetMessageType(messageType uint8) {}
