//go:binary-only-package

package nasType

// Capability5GMM 9.11.3.1
// LPP  Row, sBit, len = [0, 0], 3 , 1
// HOAttach Row, sBit, len = [0, 0], 2 , 1
// S1Mode Row, sBit, len = [0, 0], 1 , 1
// Spare Row, sBit, len = [1, 12], 8 , 96
type Capability5GMM struct {
	Iei   uint8
	Len   uint8
	Octet [13]uint8
}

func NewCapability5GMM(iei uint8) (capability5GMM *Capability5GMM) {}

// Capability5GMM 9.11.3.1
// Iei Row, sBit, len = [], 8, 8
func (a *Capability5GMM) GetIei() (iei uint8) {}

// Capability5GMM 9.11.3.1
// Iei Row, sBit, len = [], 8, 8
func (a *Capability5GMM) SetIei(iei uint8) {}

// Capability5GMM 9.11.3.1
// Len Row, sBit, len = [], 8, 8
func (a *Capability5GMM) GetLen() (len uint8) {}

// Capability5GMM 9.11.3.1
// Len Row, sBit, len = [], 8, 8
func (a *Capability5GMM) SetLen(len uint8) {}

// Capability5GMM 9.11.3.1
// LPP Row, sBit, len = [0, 0], 3 , 1
func (a *Capability5GMM) GetLPP() (lPP uint8) {}

// Capability5GMM 9.11.3.1
// LPP Row, sBit, len = [0, 0], 3 , 1
func (a *Capability5GMM) SetLPP(lPP uint8) {}

// Capability5GMM 9.11.3.1
// HOAttach Row, sBit, len = [0, 0], 2 , 1
func (a *Capability5GMM) GetHOAttach() (hOAttach uint8) {}

// Capability5GMM 9.11.3.1
// HOAttach Row, sBit, len = [0, 0], 2 , 1
func (a *Capability5GMM) SetHOAttach(hOAttach uint8) {}

// Capability5GMM 9.11.3.1
// S1Mode Row, sBit, len = [0, 0], 1 , 1
func (a *Capability5GMM) GetS1Mode() (s1Mode uint8) {}

// Capability5GMM 9.11.3.1
// S1Mode Row, sBit, len = [0, 0], 1 , 1
func (a *Capability5GMM) SetS1Mode(s1Mode uint8) {}

// Capability5GMM 9.11.3.1
// Spare Row, sBit, len = [1, 12], 8 , 96
func (a *Capability5GMM) GetSpare() (spare [12]uint8) {}

// Capability5GMM 9.11.3.1
// Spare Row, sBit, len = [1, 12], 8 , 96
func (a *Capability5GMM) SetSpare(spare [12]uint8) {}
