//go:binary-only-package

package nasType

// PDUSessionType 9.11.4.11
// Iei Row, sBit, len = [0, 0], 8 , 4
// Spare Row, sBit, len = [0, 0], 4 , 1
// PDUSessionTypeValue Row, sBit, len = [0, 0], 3 , 3
type PDUSessionType struct {
	Octet uint8
}

func NewPDUSessionType(iei uint8) (pDUSessionType *PDUSessionType) {}

// PDUSessionType 9.11.4.11
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *PDUSessionType) GetIei() (iei uint8) {}

// PDUSessionType 9.11.4.11
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *PDUSessionType) SetIei(iei uint8) {}

// PDUSessionType 9.11.4.11
// Spare Row, sBit, len = [0, 0], 4 , 1
func (a *PDUSessionType) GetSpare() (spare uint8) {}

// PDUSessionType 9.11.4.11
// Spare Row, sBit, len = [0, 0], 4 , 1
func (a *PDUSessionType) SetSpare(spare uint8) {}

// PDUSessionType 9.11.4.11
// PDUSessionTypeValue Row, sBit, len = [0, 0], 3 , 3
func (a *PDUSessionType) GetPDUSessionTypeValue() (pDUSessionTypeValue uint8) {}

// PDUSessionType 9.11.4.11
// PDUSessionTypeValue Row, sBit, len = [0, 0], 3 , 3
func (a *PDUSessionType) SetPDUSessionTypeValue(pDUSessionTypeValue uint8) {}
