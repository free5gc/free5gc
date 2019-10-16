package udm_message

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"net/http"
)

type HttpResponseMessageType string

type HandlerMessage struct {
	Event        Event
	HTTPRequest  *http_wrapper.Request
	ResponseChan chan HandlerResponseMessage //return Http response
	Value        interface{}                 // input/request value
}

func NewHandlerMessage(event Event, httpRequest *http_wrapper.Request) (msg HandlerMessage) {
	msg = HandlerMessage{}
	msg.Event = event
	msg.ResponseChan = make(chan HandlerResponseMessage)
	if httpRequest != nil {
		msg.HTTPRequest = httpRequest
	}
	return msg
}

type HandlerResponseMessage struct {
	HTTPResponse *http_wrapper.Response
}

/* Send HTTP Response to HTTP handler thread through HTTP channel, args[0] is response payload and args[1:] is Additional Value*/
func SendHttpResponseMessage(channel chan HandlerResponseMessage, header http.Header, status int, body interface{}) {
	responseMsg := HandlerResponseMessage{}
	responseMsg.HTTPResponse = http_wrapper.NewResponse(status, header, body)

	channel <- responseMsg
}

type EventGenerateAuthDataValue struct {
	request    models.AuthenticationInfoRequest
	RequestUri string
}

type EventConfirmAuthValue struct {
	request    models.AuthEvent
	RequestUri string
	Supi       string
}
