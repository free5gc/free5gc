package ausf_message

import (
	"fmt"
	"free5gc/lib/http_wrapper"
	"net/http"
)

type HttpResponseMessageType string

type HandlerMessage struct {
	Event        Event
	HTTPRequest  *http_wrapper.Request
	ResponseChan chan HandlerResponseMessage // return Http response
	Value        interface{}                 // input/request value
}

type HandlerResponseMessage struct {
	HTTPResponse *http_wrapper.Response
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

/* Send HTTP Response to HTTP handler thread through HTTP channel, args[0] is response payload and args[1:] is Additional Value*/
func SendHttpResponseMessage(channel chan HandlerResponseMessage, header http.Header, status int, body interface{}) {
	responseMsg := HandlerResponseMessage{}
	responseMsg.HTTPResponse = http_wrapper.NewResponse(status, header, body)
	fmt.Println("resp: ", responseMsg.HTTPResponse)
	channel <- responseMsg
}
