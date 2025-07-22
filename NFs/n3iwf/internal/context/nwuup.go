package context

import gtpQoSMsg "github.com/free5gc/n3iwf/internal/gtp/message"

type NwuupEventType int64

// NWuup Event Type
const (
	NwuupForwardDL NwuupEventType = iota
)

type NwuupEvt interface {
	Type() NwuupEventType
}

type NwuupForwardDLEvt struct {
	Packet gtpQoSMsg.QoSTPDUPacket
}

func (nwuupForwardDLEvt *NwuupForwardDLEvt) Type() NwuupEventType {
	return NwuupForwardDL
}

func NewNwuupForwardDLEvt(packet gtpQoSMsg.QoSTPDUPacket) *NwuupForwardDLEvt {
	return &NwuupForwardDLEvt{
		Packet: packet,
	}
}
