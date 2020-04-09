package n3iwf_handler

import (
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_data_relay"
	"free5gc/src/n3iwf/n3iwf_ike"
	"free5gc/src/n3iwf/n3iwf_ngap/ngap_message"
	"time"

	"github.com/sirupsen/logrus"

	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_handler/n3iwf_message"
	"free5gc/src/n3iwf/n3iwf_ngap"
	"free5gc/src/n3iwf/n3iwf_ngap/ngap_handler"
)

var handlerLog *logrus.Entry

func init() {
	// init pool
	handlerLog = logger.HandlerLog
}

func Handle() {
	for {
		select {
		case msg, ok := <-n3iwf_message.N3iwfChannel:
			if ok {
				switch msg.Event {
				case n3iwf_message.EventSCTPConnectMessage:
					ngap_handler.HandleEventSCTPConnect(msg.SCTPAddr)
				case n3iwf_message.EventNGAPMessage:
					n3iwf_ngap.Dispatch(msg.SCTPAddr, msg.Value.([]byte))
				case n3iwf_message.EventTimerSendRanConfigUpdateMessage:
					handlerLog.Infof("Re-send Ran Configuration Update Message when waiting time expired")
					self := n3iwf_context.N3IWFSelf()
					self.AMFReInitAvailableList[msg.SCTPAddr] = true
					ngap_message.SendRANConfigurationUpdate(self.AMFPool[msg.SCTPAddr])
				case n3iwf_message.EventN1UDPMessage:
					n3iwf_ike.Dispatch(msg.UDPSendInfo, msg.Value.([]byte))
				case n3iwf_message.EventN1TunnelUPMessage:
					self := n3iwf_context.N3IWFSelf()
					ue, ok := self.AllocatedUEIPAddress[msg.UEInnerIP]
					if !ok {
						handlerLog.Error("UE context not found")
						continue
					}
					n3iwf_data_relay.ForwardUPTrafficFromN1(ue, msg.Value.([]byte))
				case n3iwf_message.EventGTPMessage:
					self := n3iwf_context.N3IWFSelf()
					ue, ok := self.AllocatedUETEID[msg.TEID]
					if !ok {
						handlerLog.Error("UE context not found")
						continue
					}
					n3iwf_data_relay.ForwardUPTrafficFromN3(ue, msg.Value.([]byte))
				case n3iwf_message.EventN1TunnelCPMessage:
					self := n3iwf_context.N3IWFSelf()
					ue, ok := self.AllocatedUEIPAddress[msg.UEInnerIP]
					if !ok {
						handlerLog.Error("UE context not found")
						continue
					}
					n3iwf_data_relay.ForwardCPTrafficFromN1(ue, msg.Value.([]byte))
				}
			}
		case <-time.After(1 * time.Second):
		}
	}
}
