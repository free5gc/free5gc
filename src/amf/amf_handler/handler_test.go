package amf_handler_test

import (
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/ngap"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/amf_ngap"
	"free5gc/src/test/ngapTestpacket"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	go amf_handler.Handle()
	TestAmf.SctpSever()
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	message := ngapTestpacket.BuildNGSetupRequest()
	ngapMsg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	msg := amf_message.HandlerMessage{}
	msg.Event = amf_message.EventNGAPMessage
	msg.NgapAddr = TestAmf.Laddr.String()
	msg.Value = ngapMsg
	amf_message.SendMessage(msg)

	time.Sleep(100 * time.Millisecond)

}
