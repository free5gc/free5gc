//go:binary-only-package

package nasType

// AlwaysonPDUSessionRequested 9.11.4.4
// Iei Row, sBit, len = [0, 0], 8 , 4
// APSR Row, sBit, len = [0, 0], 1 , 1
type AlwaysonPDUSessionRequested struct {
	Octet uint8
}

func NewAlwaysonPDUSessionRequested(iei uint8) (alwaysonPDUSessionRequested *AlwaysonPDUSessionRequested) {}

// AlwaysonPDUSessionRequested 9.11.4.4
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *AlwaysonPDUSessionRequested) GetIei() (iei uint8) {}

// AlwaysonPDUSessionRequested 9.11.4.4
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *AlwaysonPDUSessionRequested) SetIei(iei uint8) {}

// AlwaysonPDUSessionRequested 9.11.4.4
// APSR Row, sBit, len = [0, 0], 1 , 1
func (a *AlwaysonPDUSessionRequested) GetAPSR() (aPSR uint8) {}

// AlwaysonPDUSessionRequested 9.11.4.4
// APSR Row, sBit, len = [0, 0], 1 , 1
func (a *AlwaysonPDUSessionRequested) SetAPSR(aPSR uint8) {}
