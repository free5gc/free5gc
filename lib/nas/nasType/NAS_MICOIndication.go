//go:binary-only-package

package nasType

// MICOIndication 9.11.3.31
// Iei Row, sBit, len = [0, 0], 8 , 4
// RAAI Row, sBit, len = [0, 0], 1 , 1
type MICOIndication struct {
	Octet uint8
}

func NewMICOIndication(iei uint8) (mICOIndication *MICOIndication) {}

// MICOIndication 9.11.3.31
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *MICOIndication) GetIei() (iei uint8) {}

// MICOIndication 9.11.3.31
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *MICOIndication) SetIei(iei uint8) {}

// MICOIndication 9.11.3.31
// RAAI Row, sBit, len = [0, 0], 1 , 1
func (a *MICOIndication) GetRAAI() (rAAI uint8) {}

// MICOIndication 9.11.3.31
// RAAI Row, sBit, len = [0, 0], 1 , 1
func (a *MICOIndication) SetRAAI(rAAI uint8) {}
