//go:binary-only-package

package nasType

// MobileIdentity 9.11.3.4
// MobileIdentityContents Row, sBit, len = [0, 0], 8 , INF
type MobileIdentity struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewMobileIdentity(iei uint8) (mobileIdentity *MobileIdentity) {}

// MobileIdentity 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *MobileIdentity) GetIei() (iei uint8) {}

// MobileIdentity 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *MobileIdentity) SetIei(iei uint8) {}

// MobileIdentity 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *MobileIdentity) GetLen() (len uint16) {}

// MobileIdentity 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *MobileIdentity) SetLen(len uint16) {}

// MobileIdentity 9.11.3.4
// MobileIdentityContents Row, sBit, len = [0, 0], 8 , INF
func (a *MobileIdentity) GetMobileIdentityContents() (mobileIdentityContents []uint8) {}

// MobileIdentity 9.11.3.4
// MobileIdentityContents Row, sBit, len = [0, 0], 8 , INF
func (a *MobileIdentity) SetMobileIdentityContents(mobileIdentityContents []uint8) {}
