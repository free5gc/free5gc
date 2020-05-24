package handler_test

import (
	"free5gc/src/nrf/handler"
	nrf_message "free5gc/src/nrf/handler/message"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	go handler.Handle()
	msg := nrf_message.ChannelMessage{}
	msg.Event = nrf_message.EventNFDiscovery
	//msg.Value = ngapMsg
	handler.SendMessage(msg)

	time.Sleep(100 * time.Millisecond)

}
