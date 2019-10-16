//go:binary-only-package

package nasType

// RequestedQosFlowDescriptions 9.11.4.12
// QoSFlowDescriptions Row, sBit, len = [0, 0], 8 , INF
type RequestedQosFlowDescriptions struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewRequestedQosFlowDescriptions(iei uint8) (requestedQosFlowDescriptions *RequestedQosFlowDescriptions) {}

// RequestedQosFlowDescriptions 9.11.4.12
// Iei Row, sBit, len = [], 8, 8
func (a *RequestedQosFlowDescriptions) GetIei() (iei uint8) {}

// RequestedQosFlowDescriptions 9.11.4.12
// Iei Row, sBit, len = [], 8, 8
func (a *RequestedQosFlowDescriptions) SetIei(iei uint8) {}

// RequestedQosFlowDescriptions 9.11.4.12
// Len Row, sBit, len = [], 8, 8
func (a *RequestedQosFlowDescriptions) GetLen() (len uint8) {}

// RequestedQosFlowDescriptions 9.11.4.12
// Len Row, sBit, len = [], 8, 8
func (a *RequestedQosFlowDescriptions) SetLen(len uint8) {}

// RequestedQosFlowDescriptions 9.11.4.12
// QoSFlowDescriptions Row, sBit, len = [0, 0], 8 , INF
func (a *RequestedQosFlowDescriptions) GetQoSFlowDescriptions() (qoSFlowDescriptions []uint8) {}

// RequestedQosFlowDescriptions 9.11.4.12
// QoSFlowDescriptions Row, sBit, len = [0, 0], 8 , INF
func (a *RequestedQosFlowDescriptions) SetQoSFlowDescriptions(qoSFlowDescriptions []uint8) {}
