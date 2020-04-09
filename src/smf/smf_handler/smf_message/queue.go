package smf_message

import (
	"free5gc/lib/http_wrapper"
)

var RspQueue *ResponseQueue

type ResponseQueue struct {
	RspQueue map[uint32]*ResponseQueueItem
}

func init() {
	RspQueue = NewQueue()
}

func NewQueue() *ResponseQueue {
	var rq ResponseQueue
	rq.RspQueue = make(map[uint32]*ResponseQueueItem)
	return &rq
}

func (rq ResponseQueue) PutItem(seqNum uint32, rspChan chan HandlerResponseMessage, response http_wrapper.Response) {

	Item := new(ResponseQueueItem)
	Item.RspChan = rspChan
	Item.Response = response
	rq.RspQueue[seqNum] = Item
}

func (rq ResponseQueue) GetItem(seqNum uint32) *ResponseQueueItem {
	return rq.RspQueue[seqNum]
}

func (rq ResponseQueue) DeleteItem(seqNum uint32) {
	delete(rq.RspQueue, seqNum)
}

func (rq ResponseQueue) CheckItemExist(seqNum uint32) (exist bool) {
	_, exist = rq.RspQueue[seqNum]
	return exist
}
