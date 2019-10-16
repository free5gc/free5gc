package pcf_message

import (
	"sync"
)

var PCFChannel chan ChannelMessage
var mtx sync.Mutex

const (
	MaxChannel int = 100000
)

func init() {
	// init Pool
	PCFChannel = make(chan ChannelMessage, MaxChannel)
}

func SendMessage(msg ChannelMessage) {
	mtx.Lock()
	PCFChannel <- msg
	mtx.Unlock()
}
