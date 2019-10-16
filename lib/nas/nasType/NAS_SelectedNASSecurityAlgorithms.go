//go:binary-only-package

package nasType

// SelectedNASSecurityAlgorithms 9.11.3.34
// TypeOfCipheringAlgorithm Row, sBit, len = [0, 0], 8 , 4
// TypeOfIntegrityProtectionAlgorithm Row, sBit, len = [0, 0], 4 , 4
type SelectedNASSecurityAlgorithms struct {
	Iei   uint8
	Octet uint8
}

func NewSelectedNASSecurityAlgorithms(iei uint8) (selectedNASSecurityAlgorithms *SelectedNASSecurityAlgorithms) {}

// SelectedNASSecurityAlgorithms 9.11.3.34
// Iei Row, sBit, len = [], 8, 8
func (a *SelectedNASSecurityAlgorithms) GetIei() (iei uint8) {}

// SelectedNASSecurityAlgorithms 9.11.3.34
// Iei Row, sBit, len = [], 8, 8
func (a *SelectedNASSecurityAlgorithms) SetIei(iei uint8) {}

// SelectedNASSecurityAlgorithms 9.11.3.34
// TypeOfCipheringAlgorithm Row, sBit, len = [0, 0], 8 , 4
func (a *SelectedNASSecurityAlgorithms) GetTypeOfCipheringAlgorithm() (typeOfCipheringAlgorithm uint8) {}

// SelectedNASSecurityAlgorithms 9.11.3.34
// TypeOfCipheringAlgorithm Row, sBit, len = [0, 0], 8 , 4
func (a *SelectedNASSecurityAlgorithms) SetTypeOfCipheringAlgorithm(typeOfCipheringAlgorithm uint8) {}

// SelectedNASSecurityAlgorithms 9.11.3.34
// TypeOfIntegrityProtectionAlgorithm Row, sBit, len = [0, 0], 4 , 4
func (a *SelectedNASSecurityAlgorithms) GetTypeOfIntegrityProtectionAlgorithm() (typeOfIntegrityProtectionAlgorithm uint8) {}

// SelectedNASSecurityAlgorithms 9.11.3.34
// TypeOfIntegrityProtectionAlgorithm Row, sBit, len = [0, 0], 4 , 4
func (a *SelectedNASSecurityAlgorithms) SetTypeOfIntegrityProtectionAlgorithm(typeOfIntegrityProtectionAlgorithm uint8) {}
