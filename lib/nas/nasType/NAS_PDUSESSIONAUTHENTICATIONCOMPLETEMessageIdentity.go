//go:binary-only-package

package nasType

// PDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type PDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity struct {
	Octet uint8
}

func NewPDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity() (pDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity *PDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity) {}

// PDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity) GetMessageType() (messageType uint8) {}

// PDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *PDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity) SetMessageType(messageType uint8) {}
