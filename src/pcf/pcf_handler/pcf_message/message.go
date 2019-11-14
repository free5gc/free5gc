package pcf_message

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"net/http"
)

type HttpResponseMessageType string

const (
	HttpResponseMessageResponse       HttpResponseMessageType = "HTTP Response"
	HttpResponseMessageAddition       HttpResponseMessageType = "HTTP Response Addition"
	HttpResponseMessageResponseError  HttpResponseMessageType = "HTTP Response Error"
	HttpResponseMessageProblemDetails HttpResponseMessageType = "Problem Details"
)

type ChannelMessage struct {
	Event       string
	HTTPRequest *http_wrapper.Request
	HttpChannel chan HttpResponseMessage // return Http response
	Value       interface{}              // input/request value
}

func NewHttpChannelMessage(event string, httpRequest *http_wrapper.Request) (msg ChannelMessage) {
	msg = ChannelMessage{}
	msg.Event = event
	msg.HttpChannel = make(chan HttpResponseMessage)
	if httpRequest != nil {
		msg.HTTPRequest = httpRequest
	}
	return msg
}

type HttpResponseMessage struct {
	HTTPResponse *http_wrapper.Response
}

/* Send HTTP Response to HTTP handler thread through HTTP channel, args[0] is response payload and args[1:] is Additional Value*/
func SendHttpResponseMessage(channel chan HttpResponseMessage, header http.Header, status int, body interface{}) {
	responseMsg := HttpResponseMessage{}
	responseMsg.HTTPResponse = http_wrapper.NewResponse(status, header, body)
	channel <- responseMsg
}

/*
type EventCreateBDTPolicyCreateValue struct {
	Request    models.BdtReqData
	RequestUri string
}

type EventGetBDTPolicyCreateValue struct {
	RequestUri string
}

type EventUpdateBDTPolicyCreateValue struct {
	Request    models.BdtPolicyDataPatch
	RequestUri string
}

type EventPostAppSessionsCreateValue struct {
	Request    models.AppSessionContext
	RequestUri string
}
*/
type EventGetAppSessionCreateValue struct {
	RequestUri string
}

type EventDeleteAppSessionCreateValue struct {
	RequestUri string
}

type EventModAppSessionCreateValue struct {
	Request    models.AppSessionContextUpdateData
	RequestUri string
}

type EventDeleteEventsSubscCreateValue struct {
	RequestUri string
}

type EventUpdateEventsSubscCreateValue struct {
	Request    models.EventsSubscReqData
	RequestUri string
}

type EventCreateAMPolicyCreateValue struct {
	Request    models.PolicyAssociationRequest
	RequestUri string
}

type EventGetAMPolicyCreateValue struct {
	RequestUri string
}

type EventUpdateAMPolicyCreateValue struct {
	Request    models.PolicyAssociationUpdateRequest
	RequestUri string
}

type EventDeleteAMPolicyCreateValue struct {
	RequestUri string
}
type EventCreateSMPolicyCreateValue struct {
	Request    models.SmPolicyContextData
	RequestUri string
}

type EventGetSMPolicyCreateValue struct {
	RequestUri string
}

type EventUpdateSMPolicyCreateValue struct {
	Request    models.SmPolicyUpdateContextData
	RequestUri string
}

type EventDeleteSMPolicyCreateValue struct {
	RequestUri string
}
type EventNotifySMPolicyCreateValue struct {
	Request    models.PolicyDataChangeNotification
	RequestUri string
}
