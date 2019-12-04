package n3iwf_handler

import (
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
		msg, ok := <-n3iwf_message.N3iwfChannel
		if ok {
			switch msg.Event {
			case n3iwf_message.EventSCTPConnectMessage:
				ngap_handler.HandleEventSCTPConnect(msg.SCTPSessionID)
			case n3iwf_message.EventNGAPMessage:
				n3iwf_ngap.Dispatch(msg.SCTPSessionID, msg.Value.([]byte))
			}
		}
	}
}
