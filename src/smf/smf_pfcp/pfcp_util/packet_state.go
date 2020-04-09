package pfcp_util

type PacketState uint8

const (
	RECV_REQUEST  PacketState = 0
	SEND_REQUEST  PacketState = 1
	RECV_RESPONSE PacketState = 2
	SEND_RESPONSE PacketState = 3
	FINISH        PacketState = 4
)
