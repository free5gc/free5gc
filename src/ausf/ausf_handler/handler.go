package ausf_handler

import (
	// "fmt"
	"github.com/sirupsen/logrus"
	"free5gc/lib/openapi/models"
	"free5gc/src/ausf/ausf_handler/ausf_message"
	"free5gc/src/ausf/ausf_producer"
	"free5gc/src/ausf/logger"
	"time"
)

const (
	MaxChannel int = 20
)

var ausfChannel chan ausf_message.HandlerMessage
var HandlerLog *logrus.Entry

func init() {
	HandlerLog = logger.HandlerLog
	ausfChannel = make(chan ausf_message.HandlerMessage, MaxChannel)
}

func SendMessage(msg ausf_message.HandlerMessage) {
	ausfChannel <- msg
}

func Handle() {
	for {
		select {
		case msg, ok := <-ausfChannel:
			if ok {
				switch msg.Event {
				case ausf_message.EventUeAuthPost:
					ausf_producer.HandleUeAuthPostRequest(msg.ResponseChan, msg.HTTPRequest.Body.(models.AuthenticationInfo))
				case ausf_message.EventAuth5gAkaComfirm:
					authCtxId := msg.HTTPRequest.Params["authCtxId"]
					ausf_producer.HandleAuth5gAkaComfirmRequest(msg.ResponseChan, authCtxId, msg.HTTPRequest.Body.(models.ConfirmationData))
				case ausf_message.EventEapAuthComfirm:
					authCtxId := msg.HTTPRequest.Params["authCtxId"]
					ausf_producer.HandleEapAuthComfirmRequest(msg.ResponseChan, authCtxId, msg.HTTPRequest.Body.(models.EapSession))
				default:
					HandlerLog.Warnf("AUSF Event[%d] has not implemented", msg.Event)
				}
			} else {
				HandlerLog.Errorln("AUSF Channel closed!")
			}

		case <-time.After(time.Second * 1):

		}
	}
}
