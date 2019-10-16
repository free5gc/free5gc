package nssf_handler

import (
	"github.com/sirupsen/logrus"

	. "free5gc/lib/openapi/models"
	"free5gc/src/nssf/logger"
	"free5gc/src/nssf/nssf_handler/nssf_message"
	"free5gc/src/nssf/nssf_producer"
	"free5gc/src/nssf/plugin"
)

const (
	MaxChannel int = 100
)

var nssfChannel chan nssf_message.HandlerMessage
var HandlerLog *logrus.Entry

func init() {
	// init Pool
	HandlerLog = logger.HandlerLog
	nssfChannel = make(chan nssf_message.HandlerMessage, MaxChannel)
}

func SendMessage(msg nssf_message.HandlerMessage) {
	nssfChannel <- msg
}

func Handle() {
	for {
		msg, ok := <-nssfChannel
		if ok {
			switch msg.Event {
			case nssf_message.NSSelectionGet:
				query := msg.HttpRequest.Query
				nssf_producer.NSSelectionGet(msg.ResponseChan, query)
			case nssf_message.NSSAIAvailabilityPut:
				nfId := msg.HttpRequest.Params["nfId"]
				nssaiAvailabilityInfo := msg.HttpRequest.Body.(NssaiAvailabilityInfo)
				nssf_producer.NSSAIAvailabilityPut(msg.ResponseChan, nfId, nssaiAvailabilityInfo)
			case nssf_message.NSSAIAvailabilityPatch:
				nfId := msg.HttpRequest.Params["nfId"]
				patchDocument := msg.HttpRequest.Body.(plugin.PatchDocument)
				nssf_producer.NSSAIAvailabilityPatch(msg.ResponseChan, nfId, patchDocument)
			case nssf_message.NSSAIAvailabilityDelete:
				nfId := msg.HttpRequest.Params["nfId"]
				nssf_producer.NSSAIAvailabilityDelete(msg.ResponseChan, nfId)
			case nssf_message.NSSAIAvailabilityPost:
				nssfEventSubscriptionCreateData := msg.HttpRequest.Body.(NssfEventSubscriptionCreateData)
				nssf_producer.NSSAIAvailabilityPost(msg.ResponseChan, nssfEventSubscriptionCreateData)
			case nssf_message.NSSAIAvailabilityUnsubscribe:
				subscriptionId := msg.HttpRequest.Params["subscriptionId"]
				nssf_producer.NSSAIAvailabilityUnsubscribe(msg.ResponseChan, subscriptionId)
			default:
				HandlerLog.Warnf("Event[%d] has not implemented", int(msg.Event))
			}
		} else {
			HandlerLog.Errorln("Channel closed!")
			break
		}
	}
}
