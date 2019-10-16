//go:binary-only-package

package nasType

// ServiceTypeAndNgksi 9.11.3.32 9.11.3.50
// ServiceTypeValue Row, sBit, len = [0, 0], 8 , 4
// TSC Row, sBit, len = [0, 0], 4 , 1
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 3 , 3
type ServiceTypeAndNgksi struct {
	Octet uint8
}

func NewServiceTypeAndNgksi() (serviceTypeAndNgksi *ServiceTypeAndNgksi) {}

// ServiceTypeAndNgksi 9.11.3.32 9.11.3.50
// ServiceTypeValue Row, sBit, len = [0, 0], 8 , 4
func (a *ServiceTypeAndNgksi) GetServiceTypeValue() (serviceTypeValue uint8) {}

// ServiceTypeAndNgksi 9.11.3.32 9.11.3.50
// ServiceTypeValue Row, sBit, len = [0, 0], 8 , 4
func (a *ServiceTypeAndNgksi) SetServiceTypeValue(serviceTypeValue uint8) {}

// ServiceTypeAndNgksi 9.11.3.32 9.11.3.50
// TSC Row, sBit, len = [0, 0], 4 , 1
func (a *ServiceTypeAndNgksi) GetTSC() (tSC uint8) {}

// ServiceTypeAndNgksi 9.11.3.32 9.11.3.50
// TSC Row, sBit, len = [0, 0], 4 , 1
func (a *ServiceTypeAndNgksi) SetTSC(tSC uint8) {}

// ServiceTypeAndNgksi 9.11.3.32 9.11.3.50
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 3 , 3
func (a *ServiceTypeAndNgksi) GetNasKeySetIdentifiler() (nasKeySetIdentifiler uint8) {}

// ServiceTypeAndNgksi 9.11.3.32 9.11.3.50
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 3 , 3
func (a *ServiceTypeAndNgksi) SetNasKeySetIdentifiler(nasKeySetIdentifiler uint8) {}
