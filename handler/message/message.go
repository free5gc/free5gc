package message

//"free5gc/lib/openapi/models"
//"net"

type HttpResponseMessageType string

const (
	HttpResponseMessageResponse       HttpResponseMessageType = "HTTP Response"
	HttpResponseMessageResponseError  HttpResponseMessageType = "HTTP Response Error"
	HttpResponseMessageProblemDetails HttpResponseMessageType = "Problem Details"
)

type ChannelMessage struct {
	Event       string
	HttpChannel chan HttpResponseMessage // return Http response
	//NgapConn    net.Conn                 // NGAP Connection
	Value interface{} // input/request value
}

func NewHttpChannelMessage() ChannelMessage {
	msg := ChannelMessage{}
	msg.HttpChannel = make(chan HttpResponseMessage)
	return msg
}

type HttpResponseMessage struct {
	Type            HttpResponseMessageType
	Response        interface{}
	AdditionalValue []interface{}
}

/* Send HTTP Response to HTTP handler thread through HTTP channel, args[0] is response payload and args[1:] is Additional Value*/
func SendHttpResponseMessage(channel chan HttpResponseMessage, responseType HttpResponseMessageType, args ...interface{}) {
	responseMsg := HttpResponseMessage{}
	responseMsg.Type = responseType
	responseMsg.Response = args[0]
	if len(args) > 1 {
		responseMsg.AdditionalValue = args[1:]
	}
	channel <- responseMsg
}
