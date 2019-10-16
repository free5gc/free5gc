//go:binary-only-package

package nasType

// SMPDUDNRequestContainer 9.11.4.15
// DNSpecificIdentity Row, sBit, len = [0, 0], 8 , INF
type SMPDUDNRequestContainer struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewSMPDUDNRequestContainer(iei uint8) (sMPDUDNRequestContainer *SMPDUDNRequestContainer) {}

// SMPDUDNRequestContainer 9.11.4.15
// Iei Row, sBit, len = [], 8, 8
func (a *SMPDUDNRequestContainer) GetIei() (iei uint8) {}

// SMPDUDNRequestContainer 9.11.4.15
// Iei Row, sBit, len = [], 8, 8
func (a *SMPDUDNRequestContainer) SetIei(iei uint8) {}

// SMPDUDNRequestContainer 9.11.4.15
// Len Row, sBit, len = [], 8, 8
func (a *SMPDUDNRequestContainer) GetLen() (len uint8) {}

// SMPDUDNRequestContainer 9.11.4.15
// Len Row, sBit, len = [], 8, 8
func (a *SMPDUDNRequestContainer) SetLen(len uint8) {}

// SMPDUDNRequestContainer 9.11.4.15
// DNSpecificIdentity Row, sBit, len = [0, 0], 8 , INF
func (a *SMPDUDNRequestContainer) GetDNSpecificIdentity() (dNSpecificIdentity []uint8) {}

// SMPDUDNRequestContainer 9.11.4.15
// DNSpecificIdentity Row, sBit, len = [0, 0], 8 , INF
func (a *SMPDUDNRequestContainer) SetDNSpecificIdentity(dNSpecificIdentity []uint8) {}
