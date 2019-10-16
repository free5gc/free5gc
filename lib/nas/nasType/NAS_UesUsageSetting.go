//go:binary-only-package

package nasType

// UesUsageSetting 9.11.3.55
// UesUsageSetting Row, sBit, len = [0, 0], 1 , 1
type UesUsageSetting struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewUesUsageSetting(iei uint8) (uesUsageSetting *UesUsageSetting) {}

// UesUsageSetting 9.11.3.55
// Iei Row, sBit, len = [], 8, 8
func (a *UesUsageSetting) GetIei() (iei uint8) {}

// UesUsageSetting 9.11.3.55
// Iei Row, sBit, len = [], 8, 8
func (a *UesUsageSetting) SetIei(iei uint8) {}

// UesUsageSetting 9.11.3.55
// Len Row, sBit, len = [], 8, 8
func (a *UesUsageSetting) GetLen() (len uint8) {}

// UesUsageSetting 9.11.3.55
// Len Row, sBit, len = [], 8, 8
func (a *UesUsageSetting) SetLen(len uint8) {}

// UesUsageSetting 9.11.3.55
// UesUsageSetting Row, sBit, len = [0, 0], 1 , 1
func (a *UesUsageSetting) GetUesUsageSetting() (uesUsageSetting uint8) {}

// UesUsageSetting 9.11.3.55
// UesUsageSetting Row, sBit, len = [0, 0], 1 , 1
func (a *UesUsageSetting) SetUesUsageSetting(uesUsageSetting uint8) {}
