package message

import (
	"free5gc/lib/http_wrapper"
)

type HandlerMessage struct {
	Event        string
	HTTPRequest  *http_wrapper.Request
	ResponseChan chan HandlerResponseMessage // return Http response
	Value        interface{}                 // input/request value
}

type HandlerResponseMessage struct {
	HTTPResponse *http_wrapper.Response
}

func NewMessage(event string, httpRequest *http_wrapper.Request) (msg HandlerMessage) {
	msg = HandlerMessage{}
	msg.Event = event
	msg.ResponseChan = make(chan HandlerResponseMessage)
	if httpRequest != nil {
		msg.HTTPRequest = httpRequest
	}
	return msg
}
