//go:binary-only-package

package nasType

// Cause5GMM 9.11.3.2
// CauseValue Row, sBit, len = [0, 0], 8 , 8
type Cause5GMM struct {
	Iei   uint8
	Octet uint8
}

func NewCause5GMM(iei uint8) (cause5GMM *Cause5GMM) {}

// Cause5GMM 9.11.3.2
// Iei Row, sBit, len = [], 8, 8
func (a *Cause5GMM) GetIei() (iei uint8) {}

// Cause5GMM 9.11.3.2
// Iei Row, sBit, len = [], 8, 8
func (a *Cause5GMM) SetIei(iei uint8) {}

// Cause5GMM 9.11.3.2
// CauseValue Row, sBit, len = [0, 0], 8 , 8
func (a *Cause5GMM) GetCauseValue() (causeValue uint8) {}

// Cause5GMM 9.11.3.2
// CauseValue Row, sBit, len = [0, 0], 8 , 8
func (a *Cause5GMM) SetCauseValue(causeValue uint8) {}
