//go:binary-only-package

package nasType

// MessageAuthenticationCode MAC 9.8
// MAC Row, sBit, len = [0, 3], 8 , 32
type MessageAuthenticationCode struct {
	Octet [4]uint8
}

func NewMessageAuthenticationCode() (messageAuthenticationCode *MessageAuthenticationCode) {}

// MessageAuthenticationCode MAC 9.8
// MAC Row, sBit, len = [0, 3], 8 , 32
func (a *MessageAuthenticationCode) GetMAC() (mAC [4]uint8) {}

// MessageAuthenticationCode MAC 9.8
// MAC Row, sBit, len = [0, 3], 8 , 32
func (a *MessageAuthenticationCode) SetMAC(mAC [4]uint8) {}
