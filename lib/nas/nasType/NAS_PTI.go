//go:binary-only-package

package nasType

// PTI 9.6
// PTI Row, sBit, len = [0, 0], 8 , 8
type PTI struct {
	Octet uint8
}

func NewPTI() (pTI *PTI) {}

// PTI 9.6
// PTI Row, sBit, len = [0, 0], 8 , 8
func (a *PTI) GetPTI() (pTI uint8) {}

// PTI 9.6
// PTI Row, sBit, len = [0, 0], 8 , 8
func (a *PTI) SetPTI(pTI uint8) {}
