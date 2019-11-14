package smf_message

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/lib/pfcp/pfcpUdp"
)

type HandlerMessage struct {
	Event        Event
	HTTPRequest  *http_wrapper.Request
	PFCPRequest  *pfcpUdp.Message
	ResponseChan chan HandlerResponseMessage
}

type HandlerResponseMessage struct {
	HTTPResponse *http_wrapper.Response
	PFCPResponse *pfcpUdp.Message
}

type ResponseQueueItem struct {
	RspChan      chan HandlerResponseMessage
	ResponseBody models.UpdateSmContextResponse
}

func NewPfcpMessage(pfcpRequest *pfcpUdp.Message) (msg HandlerMessage) {
	msg = HandlerMessage{}
	msg.Event = PFCPMessage
	msg.ResponseChan = make(chan HandlerResponseMessage)
	msg.PFCPRequest = pfcpRequest
	return
}

func NewHandlerMessage(event Event, httpRequest *http_wrapper.Request) (msg HandlerMessage) {
	msg = HandlerMessage{}
	msg.Event = event
	msg.ResponseChan = make(chan HandlerResponseMessage)
	if httpRequest != nil {
		msg.HTTPRequest = httpRequest
	}
	return
}
