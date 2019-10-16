//go:binary-only-package

package nasType

// IntegrityProtectionMaximumDataRate 9.11.4.7
// MaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink Row, sBit, len = [0, 0], 8 , 8
// MaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink Row, sBit, len = [1, 1], 8 , 8
type IntegrityProtectionMaximumDataRate struct {
	Iei   uint8
	Octet [2]uint8
}

func NewIntegrityProtectionMaximumDataRate(iei uint8) (integrityProtectionMaximumDataRate *IntegrityProtectionMaximumDataRate) {}

// IntegrityProtectionMaximumDataRate 9.11.4.7
// Iei Row, sBit, len = [], 8, 8
func (a *IntegrityProtectionMaximumDataRate) GetIei() (iei uint8) {}

// IntegrityProtectionMaximumDataRate 9.11.4.7
// Iei Row, sBit, len = [], 8, 8
func (a *IntegrityProtectionMaximumDataRate) SetIei(iei uint8) {}

// IntegrityProtectionMaximumDataRate 9.11.4.7
// MaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink Row, sBit, len = [0, 0], 8 , 8
func (a *IntegrityProtectionMaximumDataRate) GetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink() (maximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink uint8) {}

// IntegrityProtectionMaximumDataRate 9.11.4.7
// MaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink Row, sBit, len = [0, 0], 8 , 8
func (a *IntegrityProtectionMaximumDataRate) SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink(maximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink uint8) {}

// IntegrityProtectionMaximumDataRate 9.11.4.7
// MaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink Row, sBit, len = [1, 1], 8 , 8
func (a *IntegrityProtectionMaximumDataRate) GetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink() (maximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink uint8) {}

// IntegrityProtectionMaximumDataRate 9.11.4.7
// MaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink Row, sBit, len = [1, 1], 8 , 8
func (a *IntegrityProtectionMaximumDataRate) SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink(maximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink uint8) {}
