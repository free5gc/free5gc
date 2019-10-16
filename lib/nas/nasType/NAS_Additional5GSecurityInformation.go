//go:binary-only-package

package nasType

// Additional5GSecurityInformation 9.11.3.12
// RINMR Row, sBit, len = [0, 0], 2 , 1
// HDP Row, sBit, len = [0, 0], 1 , 1
type Additional5GSecurityInformation struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewAdditional5GSecurityInformation(iei uint8) (additional5GSecurityInformation *Additional5GSecurityInformation) {}

// Additional5GSecurityInformation 9.11.3.12
// Iei Row, sBit, len = [], 8, 8
func (a *Additional5GSecurityInformation) GetIei() (iei uint8) {}

// Additional5GSecurityInformation 9.11.3.12
// Iei Row, sBit, len = [], 8, 8
func (a *Additional5GSecurityInformation) SetIei(iei uint8) {}

// Additional5GSecurityInformation 9.11.3.12
// Len Row, sBit, len = [], 8, 8
func (a *Additional5GSecurityInformation) GetLen() (len uint8) {}

// Additional5GSecurityInformation 9.11.3.12
// Len Row, sBit, len = [], 8, 8
func (a *Additional5GSecurityInformation) SetLen(len uint8) {}

// Additional5GSecurityInformation 9.11.3.12
// RINMR Row, sBit, len = [0, 0], 2 , 1
func (a *Additional5GSecurityInformation) GetRINMR() (rINMR uint8) {}

// Additional5GSecurityInformation 9.11.3.12
// RINMR Row, sBit, len = [0, 0], 2 , 1
func (a *Additional5GSecurityInformation) SetRINMR(rINMR uint8) {}

// Additional5GSecurityInformation 9.11.3.12
// HDP Row, sBit, len = [0, 0], 1 , 1
func (a *Additional5GSecurityInformation) GetHDP() (hDP uint8) {}

// Additional5GSecurityInformation 9.11.3.12
// HDP Row, sBit, len = [0, 0], 1 , 1
func (a *Additional5GSecurityInformation) SetHDP(hDP uint8) {}
