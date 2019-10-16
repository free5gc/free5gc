//go:binary-only-package

package nasType

// AuthorizedQosFlowDescriptions 9.11.4.12
// QoSFlowDescription Row, sBit, len = [0, 0], 8 , INF
type AuthorizedQosFlowDescriptions struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewAuthorizedQosFlowDescriptions(iei uint8) (authorizedQosFlowDescriptions *AuthorizedQosFlowDescriptions) {}

// AuthorizedQosFlowDescriptions 9.11.4.12
// Iei Row, sBit, len = [], 8, 8
func (a *AuthorizedQosFlowDescriptions) GetIei() (iei uint8) {}

// AuthorizedQosFlowDescriptions 9.11.4.12
// Iei Row, sBit, len = [], 8, 8
func (a *AuthorizedQosFlowDescriptions) SetIei(iei uint8) {}

// AuthorizedQosFlowDescriptions 9.11.4.12
// Len Row, sBit, len = [], 8, 8
func (a *AuthorizedQosFlowDescriptions) GetLen() (len uint8) {}

// AuthorizedQosFlowDescriptions 9.11.4.12
// Len Row, sBit, len = [], 8, 8
func (a *AuthorizedQosFlowDescriptions) SetLen(len uint8) {}

// AuthorizedQosFlowDescriptions 9.11.4.12
// QoSFlowDescription Row, sBit, len = [0, 0], 8 , INF
func (a *AuthorizedQosFlowDescriptions) GetQoSFlowDescription() (qoSFlowDescription []uint8) {}

// AuthorizedQosFlowDescriptions 9.11.4.12
// QoSFlowDescription Row, sBit, len = [0, 0], 8 , INF
func (a *AuthorizedQosFlowDescriptions) SetQoSFlowDescription(qoSFlowDescription []uint8) {}
