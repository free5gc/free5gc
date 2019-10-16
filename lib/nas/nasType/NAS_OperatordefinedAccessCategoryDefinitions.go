//go:binary-only-package

package nasType

// OperatordefinedAccessCategoryDefinitions 9.11.3.38
// OperatorDefinedAccessCategoryDefintiion Row, sBit, len = [0, 0], 8 , INF
type OperatordefinedAccessCategoryDefinitions struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewOperatordefinedAccessCategoryDefinitions(iei uint8) (operatordefinedAccessCategoryDefinitions *OperatordefinedAccessCategoryDefinitions) {}

// OperatordefinedAccessCategoryDefinitions 9.11.3.38
// Iei Row, sBit, len = [], 8, 8
func (a *OperatordefinedAccessCategoryDefinitions) GetIei() (iei uint8) {}

// OperatordefinedAccessCategoryDefinitions 9.11.3.38
// Iei Row, sBit, len = [], 8, 8
func (a *OperatordefinedAccessCategoryDefinitions) SetIei(iei uint8) {}

// OperatordefinedAccessCategoryDefinitions 9.11.3.38
// Len Row, sBit, len = [], 8, 16
func (a *OperatordefinedAccessCategoryDefinitions) GetLen() (len uint16) {}

// OperatordefinedAccessCategoryDefinitions 9.11.3.38
// Len Row, sBit, len = [], 8, 16
func (a *OperatordefinedAccessCategoryDefinitions) SetLen(len uint16) {}

// OperatordefinedAccessCategoryDefinitions 9.11.3.38
// OperatorDefinedAccessCategoryDefintiion Row, sBit, len = [0, 0], 8 , INF
func (a *OperatordefinedAccessCategoryDefinitions) GetOperatorDefinedAccessCategoryDefintiion() (operatorDefinedAccessCategoryDefintiion []uint8) {}

// OperatordefinedAccessCategoryDefinitions 9.11.3.38
// OperatorDefinedAccessCategoryDefintiion Row, sBit, len = [0, 0], 8 , INF
func (a *OperatordefinedAccessCategoryDefinitions) SetOperatorDefinedAccessCategoryDefintiion(operatorDefinedAccessCategoryDefintiion []uint8) {}
