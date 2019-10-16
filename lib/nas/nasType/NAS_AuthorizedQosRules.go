//go:binary-only-package

package nasType

// AuthorizedQosRules 9.11.4.13
// QosRule Row, sBit, len = [0, 0], 3 , INF
type AuthorizedQosRules struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewAuthorizedQosRules(iei uint8) (authorizedQosRules *AuthorizedQosRules) {}

// AuthorizedQosRules 9.11.4.13
// Iei Row, sBit, len = [], 8, 8
func (a *AuthorizedQosRules) GetIei() (iei uint8) {}

// AuthorizedQosRules 9.11.4.13
// Iei Row, sBit, len = [], 8, 8
func (a *AuthorizedQosRules) SetIei(iei uint8) {}

// AuthorizedQosRules 9.11.4.13
// Len Row, sBit, len = [], 8, 16
func (a *AuthorizedQosRules) GetLen() (len uint16) {}

// AuthorizedQosRules 9.11.4.13
// Len Row, sBit, len = [], 8, 16
func (a *AuthorizedQosRules) SetLen(len uint16) {}

// AuthorizedQosRules 9.11.4.13
// QosRule Row, sBit, len = [0, 0], 3 , INF
func (a *AuthorizedQosRules) GetQosRule() (qosRule []uint8) {}

// AuthorizedQosRules 9.11.4.13
// QosRule Row, sBit, len = [0, 0], 3 , INF
func (a *AuthorizedQosRules) SetQosRule(qosRule []uint8) {}
