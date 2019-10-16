//go:binary-only-package

package nasType

// SpareHalfOctetAndNgksi 9.11.3.32 9.5
// SpareHalfOctet Row, sBit, len = [0, 0], 8 , 4
// TSC Row, sBit, len = [0, 0], 4 , 1
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 3 , 3
type SpareHalfOctetAndNgksi struct {
	Octet uint8
}

func NewSpareHalfOctetAndNgksi() (spareHalfOctetAndNgksi *SpareHalfOctetAndNgksi) {}

// SpareHalfOctetAndNgksi 9.11.3.32 9.5
// SpareHalfOctet Row, sBit, len = [0, 0], 8 , 4
func (a *SpareHalfOctetAndNgksi) GetSpareHalfOctet() (spareHalfOctet uint8) {}

// SpareHalfOctetAndNgksi 9.11.3.32 9.5
// SpareHalfOctet Row, sBit, len = [0, 0], 8 , 4
func (a *SpareHalfOctetAndNgksi) SetSpareHalfOctet(spareHalfOctet uint8) {}

// SpareHalfOctetAndNgksi 9.11.3.32 9.5
// TSC Row, sBit, len = [0, 0], 4 , 1
func (a *SpareHalfOctetAndNgksi) GetTSC() (tSC uint8) {}

// SpareHalfOctetAndNgksi 9.11.3.32 9.5
// TSC Row, sBit, len = [0, 0], 4 , 1
func (a *SpareHalfOctetAndNgksi) SetTSC(tSC uint8) {}

// SpareHalfOctetAndNgksi 9.11.3.32 9.5
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 3 , 3
func (a *SpareHalfOctetAndNgksi) GetNasKeySetIdentifiler() (nasKeySetIdentifiler uint8) {}

// SpareHalfOctetAndNgksi 9.11.3.32 9.5
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 3 , 3
func (a *SpareHalfOctetAndNgksi) SetNasKeySetIdentifiler(nasKeySetIdentifiler uint8) {}
