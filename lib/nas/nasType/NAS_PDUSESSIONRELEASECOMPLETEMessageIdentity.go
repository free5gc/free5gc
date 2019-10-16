//go:binary-only-package

package nasType

// PDUSESSIONRELEASECOMPLETEMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONRELEASECOMPLETEMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONRELEASECOMPLETEMessageIdentity() (pDUSESSIONRELEASECOMPLETEMessageIdentity *PDUSESSIONRELEASECOMPLETEMessageIdentity) {}

// PDUSESSIONRELEASECOMPLETEMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONRELEASECOMPLETEMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONRELEASECOMPLETEMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONRELEASECOMPLETEMessageIdentity) SetMessageType(messageType uint8) {}
