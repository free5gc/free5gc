//go:binary-only-package

package nasType

// SelectedSSCModeAndSelectedPDUSessionType 9.11.4.11 9.11.4.16
// SSCMode Row, sBit, len = [0, 0], 7 , 3
// PDUSessionType  Row, sBit, len = [0, 0], 3 , 3
type SelectedSSCModeAndSelectedPDUSessionType struct {
	Octet uint8
}

func NewSelectedSSCModeAndSelectedPDUSessionType() (selectedSSCModeAndSelectedPDUSessionType *SelectedSSCModeAndSelectedPDUSessionType) {}

// SelectedSSCModeAndSelectedPDUSessionType 9.11.4.11 9.11.4.16
// SSCMode Row, sBit, len = [0, 0], 7 , 3
func (a *SelectedSSCModeAndSelectedPDUSessionType) GetSSCMode() (sSCMode uint8) {}

// SelectedSSCModeAndSelectedPDUSessionType 9.11.4.11 9.11.4.16
// SSCMode Row, sBit, len = [0, 0], 7 , 3
func (a *SelectedSSCModeAndSelectedPDUSessionType) SetSSCMode(sSCMode uint8) {}

// SelectedSSCModeAndSelectedPDUSessionType 9.11.4.11 9.11.4.16
// PDUSessionType Row, sBit, len = [0, 0], 3 , 3
func (a *SelectedSSCModeAndSelectedPDUSessionType) GetPDUSessionType() (pDUSessionType uint8) {}

// SelectedSSCModeAndSelectedPDUSessionType 9.11.4.11 9.11.4.16
// PDUSessionType Row, sBit, len = [0, 0], 3 , 3
func (a *SelectedSSCModeAndSelectedPDUSessionType) SetPDUSessionType(pDUSessionType uint8) {}
