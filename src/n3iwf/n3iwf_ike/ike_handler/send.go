package ike_handler

import (
	"free5gc/src/n3iwf/n3iwf_handler/n3iwf_message"
	"free5gc/src/n3iwf/n3iwf_ike/ike_message"
	"free5gc/src/n3iwf/n3iwf_ike/udp_server"
)

func SendIKEMessageToUE(ueSendInfo *n3iwf_message.UDPSendInfoGroup, message *ike_message.IKEMessage) {
	ikeLog.Trace("[IKE] Send IKE message to UE")
	ikeLog.Trace("[IKE] Encoding...")
	pkt, err := ike_message.Encode(message)
	if err != nil {
		ikeLog.Errorln(err)
		return
	}
	// As specified in RFC 7296 section 3.1, the IKE message send from/to UDP port 4500
	// should prepend a 4 bytes zero
	if ueSendInfo.ChannelID == udp_server.ChannelIDForPort4500 {
		prependZero := make([]byte, 4)
		pkt = append(prependZero, pkt...)
	}
	ikeLog.Trace("[IKE] Sending...")
	udp_server.Send(ueSendInfo, pkt)
}
