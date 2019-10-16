//go:binary-only-package

package nasType

// PDUSessionReactivationResultErrorCause 9.11.3.43
// PDUSessionIDAndCauseValue Row, sBit, len = [0, 0], 8 , INF
type PDUSessionReactivationResultErrorCause struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewPDUSessionReactivationResultErrorCause(iei uint8) (pDUSessionReactivationResultErrorCause *PDUSessionReactivationResultErrorCause) {}

// PDUSessionReactivationResultErrorCause 9.11.3.43
// Iei Row, sBit, len = [], 8, 8
func (a *PDUSessionReactivationResultErrorCause) GetIei() (iei uint8) {}

// PDUSessionReactivationResultErrorCause 9.11.3.43
// Iei Row, sBit, len = [], 8, 8
func (a *PDUSessionReactivationResultErrorCause) SetIei(iei uint8) {}

// PDUSessionReactivationResultErrorCause 9.11.3.43
// Len Row, sBit, len = [], 8, 16
func (a *PDUSessionReactivationResultErrorCause) GetLen() (len uint16) {}

// PDUSessionReactivationResultErrorCause 9.11.3.43
// Len Row, sBit, len = [], 8, 16
func (a *PDUSessionReactivationResultErrorCause) SetLen(len uint16) {}

// PDUSessionReactivationResultErrorCause 9.11.3.43
// PDUSessionIDAndCauseValue Row, sBit, len = [0, 0], 8 , INF
func (a *PDUSessionReactivationResultErrorCause) GetPDUSessionIDAndCauseValue() (pDUSessionIDAndCauseValue []uint8) {}

// PDUSessionReactivationResultErrorCause 9.11.3.43
// PDUSessionIDAndCauseValue Row, sBit, len = [0, 0], 8 , INF
func (a *PDUSessionReactivationResultErrorCause) SetPDUSessionIDAndCauseValue(pDUSessionIDAndCauseValue []uint8) {}
