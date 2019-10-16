//go:binary-only-package

package nasType

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// TSC Row, sBit, len = [0, 0], 8 , 1
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 7 , 3
// FOR  Row, sBit, len = [0, 0], 4 , 1
// RegistrationType5GS Row, sBit, len = [0, 0], 3 , 3
type NgksiAndRegistrationType5GS struct {
	Octet uint8
}

func NewNgksiAndRegistrationType5GS() (ngksiAndRegistrationType5GS *NgksiAndRegistrationType5GS) {}

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// TSC Row, sBit, len = [0, 0], 8 , 1
func (a *NgksiAndRegistrationType5GS) GetTSC() (tSC uint8) {}

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// TSC Row, sBit, len = [0, 0], 8 , 1
func (a *NgksiAndRegistrationType5GS) SetTSC(tSC uint8) {}

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 7 , 3
func (a *NgksiAndRegistrationType5GS) GetNasKeySetIdentifiler() (nasKeySetIdentifiler uint8) {}

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 7 , 3
func (a *NgksiAndRegistrationType5GS) SetNasKeySetIdentifiler(nasKeySetIdentifiler uint8) {}

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// FOR Row, sBit, len = [0, 0], 4 , 1
func (a *NgksiAndRegistrationType5GS) GetFOR() (fOR uint8) {}

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// FOR Row, sBit, len = [0, 0], 4 , 1
func (a *NgksiAndRegistrationType5GS) SetFOR(fOR uint8) {}

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// RegistrationType5GS Row, sBit, len = [0, 0], 3 , 3
func (a *NgksiAndRegistrationType5GS) GetRegistrationType5GS() (registrationType5GS uint8) {}

// NgksiAndRegistrationType5GS 9.11.3.7 9.11.3.32
// RegistrationType5GS Row, sBit, len = [0, 0], 3 , 3
func (a *NgksiAndRegistrationType5GS) SetRegistrationType5GS(registrationType5GS uint8) {}
