//go:binary-only-package

package nasType

// NetworkSlicingIndication 9.11.3.36
// Iei Row, sBit, len = [0, 0], 8 , 4
// DCNI Row, sBit, len = [0, 0], 2 , 1
// NSSCI Row, sBit, len = [0, 0], 1 , 1
type NetworkSlicingIndication struct {
	Octet uint8
}

func NewNetworkSlicingIndication(iei uint8) (networkSlicingIndication *NetworkSlicingIndication) {}

// NetworkSlicingIndication 9.11.3.36
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *NetworkSlicingIndication) GetIei() (iei uint8) {}

// NetworkSlicingIndication 9.11.3.36
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *NetworkSlicingIndication) SetIei(iei uint8) {}

// NetworkSlicingIndication 9.11.3.36
// DCNI Row, sBit, len = [0, 0], 2 , 1
func (a *NetworkSlicingIndication) GetDCNI() (dCNI uint8) {}

// NetworkSlicingIndication 9.11.3.36
// DCNI Row, sBit, len = [0, 0], 2 , 1
func (a *NetworkSlicingIndication) SetDCNI(dCNI uint8) {}

// NetworkSlicingIndication 9.11.3.36
// NSSCI Row, sBit, len = [0, 0], 1 , 1
func (a *NetworkSlicingIndication) GetNSSCI() (nSSCI uint8) {}

// NetworkSlicingIndication 9.11.3.36
// NSSCI Row, sBit, len = [0, 0], 1 , 1
func (a *NetworkSlicingIndication) SetNSSCI(nSSCI uint8) {}
