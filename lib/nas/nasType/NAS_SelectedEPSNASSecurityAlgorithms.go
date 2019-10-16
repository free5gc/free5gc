//go:binary-only-package

package nasType

// SelectedEPSNASSecurityAlgorithms 9.11.3.25
// TypeOfCipheringAlgorithm Row, sBit, len = [0, 0], 7 , 3
// TypeOfIntegrityProtectionAlgorithm Row, sBit, len = [0, 0], 3 , 3
type SelectedEPSNASSecurityAlgorithms struct {
	Iei   uint8
	Octet uint8
}

func NewSelectedEPSNASSecurityAlgorithms(iei uint8) (selectedEPSNASSecurityAlgorithms *SelectedEPSNASSecurityAlgorithms) {}

// SelectedEPSNASSecurityAlgorithms 9.11.3.25
// Iei Row, sBit, len = [], 8, 8
func (a *SelectedEPSNASSecurityAlgorithms) GetIei() (iei uint8) {}

// SelectedEPSNASSecurityAlgorithms 9.11.3.25
// Iei Row, sBit, len = [], 8, 8
func (a *SelectedEPSNASSecurityAlgorithms) SetIei(iei uint8) {}

// SelectedEPSNASSecurityAlgorithms 9.11.3.25
// TypeOfCipheringAlgorithm Row, sBit, len = [0, 0], 7 , 3
func (a *SelectedEPSNASSecurityAlgorithms) GetTypeOfCipheringAlgorithm() (typeOfCipheringAlgorithm uint8) {}

// SelectedEPSNASSecurityAlgorithms 9.11.3.25
// TypeOfCipheringAlgorithm Row, sBit, len = [0, 0], 7 , 3
func (a *SelectedEPSNASSecurityAlgorithms) SetTypeOfCipheringAlgorithm(typeOfCipheringAlgorithm uint8) {}

// SelectedEPSNASSecurityAlgorithms 9.11.3.25
// TypeOfIntegrityProtectionAlgorithm Row, sBit, len = [0, 0], 3 , 3
func (a *SelectedEPSNASSecurityAlgorithms) GetTypeOfIntegrityProtectionAlgorithm() (typeOfIntegrityProtectionAlgorithm uint8) {}

// SelectedEPSNASSecurityAlgorithms 9.11.3.25
// TypeOfIntegrityProtectionAlgorithm Row, sBit, len = [0, 0], 3 , 3
func (a *SelectedEPSNASSecurityAlgorithms) SetTypeOfIntegrityProtectionAlgorithm(typeOfIntegrityProtectionAlgorithm uint8) {}
