package n3iwf_message

import (
	"net"
)

type HandlerMessage struct {
	Event       Event
	UDPSendInfo *UDPSendInfoGroup // used only when Event == EventN1UDPMessage
	SCTPAddr    string            // used when Event == EventNGAPMessage || Event == EventSCTPConnectMessage
	UEInnerIP   string            // used when Event == EventN1TunnelCPMessage || Event == EventN1TunnelUPMessage
	TEID        uint32            // used only when Event == EventGTPMessage
	Value       interface{}
}

type UDPSendInfoGroup struct {
	ChannelID int
	Addr      *net.UDPAddr
}
