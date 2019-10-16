package amf_message

import (
	"sync"
)

var AmfChannel chan HandlerMessage
var mtx sync.Mutex

const (
	MaxChannel int = 100000
)

func init() {
	// init Pool
	AmfChannel = make(chan HandlerMessage, MaxChannel)
}

func SendMessage(msg HandlerMessage) {
	mtx.Lock()
	AmfChannel <- msg
	mtx.Unlock()
}
