package handler

import (
	"free5gc/lib/openapi/models"
	nrf_message "free5gc/src/nrf/handler/message"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/producer"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	MaxChannel int = 100000
)

var nrfChannel chan nrf_message.HandlerMessage
var HandlerLog *logrus.Entry

func init() {
	// init Pool
	HandlerLog = logger.HandlerLog
	nrfChannel = make(chan nrf_message.HandlerMessage, 20)
}

func SendMessage(msg nrf_message.HandlerMessage) {
	nrfChannel <- msg
	//go Handle()
}

func Handle() {
	for {
		select {
		case msg, ok := <-nrfChannel:
			if ok {
				switch msg.Event {
				case nrf_message.EventNFDiscovery:
					HandlerLog.Info("EventNFDiscovery")
				case nrf_message.EventNFManagement:
					logger.HandlerLog.Info("EventNFManagement")
				case nrf_message.EventNotificationNFRegisted:
					HandlerLog.Info("EventNotificationNFRegisted")
					url := msg.HTTPRequest.Params["url"]
					producer.HandleNotification(msg.ResponseChan, url, msg.HTTPRequest.Body.(models.NotificationData))
				case nrf_message.EventNotificationNFDeregisted:
					HandlerLog.Info("EventNotificationNFDeregisted")
					url := msg.HTTPRequest.Params["url"]
					producer.HandleNotification(msg.ResponseChan, url, msg.HTTPRequest.Body.(models.NotificationData))
				case nrf_message.EventNotificationNFProfileChanged:
					HandlerLog.Info("EventNotificationNFProfileChanged")
					url := msg.HTTPRequest.Params["url"]
					producer.HandleNotification(msg.ResponseChan, url, msg.HTTPRequest.Body.(models.NotificationData))
				case nrf_message.EventAccessToken:
					HandlerLog.Info("EventAccessToken")
				default:
					HandlerLog.Warnf("Event[%s] has not implemented", msg.Event)
				}
			} else {
				HandlerLog.Errorln("Channel closed!")
			}

		case <-time.After(time.Second * 1):

		}
	}
}
