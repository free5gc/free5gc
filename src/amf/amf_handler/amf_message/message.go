package amf_message

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"net/http"
)

type HttpResponseMessageType string

type HandlerMessage struct {
	Event        Event
	HTTPRequest  *http_wrapper.Request
	ResponseChan chan HandlerResponseMessage // return Http response
	NgapAddr     string                      // NGAP Connection Addr
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

type EventN1N2MessageTransferValue struct {
	UeContextId string
	Request     models.N1N2MessageTransferRequest
	RequestUri  string
}

type EventN1N2MessageTransferStatusValue struct {
	UeContextId string
	RequestUri  string
}

type EventProvideDomainSelectionInfoValue struct {
	UeContextId        string
	RequestUri         string
	UeContextInfoClass string
}

type EventProvideLocationInfoValue struct {
	UeContextId string
	Request     models.RequestLocInfo
	RequestUri  string
}

type EventN1N2MessageSubscribeValue struct {
	UeContextId string
	Request     models.UeN1N2InfoSubscriptionCreateData
	RequestUri  string
}

type EventGMMT3560ValueForSecurityCommand struct {
	RanUe      *amf_context.RanUe
	EapSuccess bool
	EapMessage string
}

type EventGMMT3550Value struct {
	AmfUe                       *amf_context.AmfUe
	AccessType                  models.AccessType
	PDUSessionStatus            *[16]bool
	ReactivationResult          *[16]bool
	ErrPduSessionId             []uint8
	ErrCause                    []uint8
	PduSessionResourceSetupList *ngapType.PDUSessionResourceSetupListCxtReq
}

type EventGMMT3522Value struct {
	RanUe                  *amf_context.RanUe
	AccessType             uint8
	ReRegistrationRequired bool
	Cause5GMM              uint8
}
