package udr_message

import (
	"sync"
)

var UdrChannel chan HandlerMessage
var mtx sync.Mutex

const (
	MaxChannel int = 100000
)

func init() {
	// init Pool
	UdrChannel = make(chan HandlerMessage, MaxChannel)
}

func SendMessage(msg HandlerMessage) {
	mtx.Lock()
	UdrChannel <- msg
	mtx.Unlock()
}
