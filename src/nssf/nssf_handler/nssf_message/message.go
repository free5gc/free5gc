package nssf_message

import (
	"free5gc/lib/http_wrapper"
)

type HandlerMessage struct {
	Event        Event
	HttpRequest  *http_wrapper.Request
	ResponseChan chan HandlerResponseMessage
}

type HandlerResponseMessage struct {
	HttpResponse *http_wrapper.Response
}

func NewMessage(event Event, httpRequest *http_wrapper.Request) (msg HandlerMessage) {
	msg = HandlerMessage{}
	msg.Event = event
	msg.ResponseChan = make(chan HandlerResponseMessage)
	if httpRequest != nil {
		msg.HttpRequest = httpRequest
	}
	return
}
