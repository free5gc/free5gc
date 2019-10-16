//go:binary-only-package

package nasType

// RequestedQosRules 9.11.4.13
// QoSRules Row, sBit, len = [0, 0], 8 , INF
type RequestedQosRules struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewRequestedQosRules(iei uint8) (requestedQosRules *RequestedQosRules) {}

// RequestedQosRules 9.11.4.13
// Iei Row, sBit, len = [], 8, 8
func (a *RequestedQosRules) GetIei() (iei uint8) {}

// RequestedQosRules 9.11.4.13
// Iei Row, sBit, len = [], 8, 8
func (a *RequestedQosRules) SetIei(iei uint8) {}

// RequestedQosRules 9.11.4.13
// Len Row, sBit, len = [], 8, 8
func (a *RequestedQosRules) GetLen() (len uint8) {}

// RequestedQosRules 9.11.4.13
// Len Row, sBit, len = [], 8, 8
func (a *RequestedQosRules) SetLen(len uint8) {}

// RequestedQosRules 9.11.4.13
// QoSRules Row, sBit, len = [0, 0], 8 , INF
func (a *RequestedQosRules) GetQoSRules() (qoSRules []uint8) {}

// RequestedQosRules 9.11.4.13
// QoSRules Row, sBit, len = [0, 0], 8 , INF
func (a *RequestedQosRules) SetQoSRules(qoSRules []uint8) {}
