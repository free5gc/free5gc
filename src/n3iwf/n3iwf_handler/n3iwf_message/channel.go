package n3iwf_message

import (
	"sync"
)

var N3iwfChannel chan HandlerMessage
var mtx sync.Mutex

const (
	MaxChannel int = 100000
)

func init() {
	// init Pool
	N3iwfChannel = make(chan HandlerMessage, MaxChannel)
}

func SendMessage(msg HandlerMessage) {
	mtx.Lock()
	N3iwfChannel <- msg
	mtx.Unlock()
}
