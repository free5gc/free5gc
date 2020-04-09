package smf_context

import "fmt"

type IDQueue struct {
	QueueType IDType
	Queue     []int
}

type IDType int

const (
	PDRType  IDType = iota
	FARType  
	BARType  
	TEIDType 
)

func NewIDQueue(idType IDType) (idQueue *IDQueue) {

	idQueue = &IDQueue{
		QueueType: idType,
	}
	idQueue.Queue = make([]int, 0)

	return
}

func (q *IDQueue) Push(item int) {
	q.Queue = append(q.Queue, item)
}

func (q *IDQueue) Pop() (id int, err error) {

	id = -1
	err = nil

	if !q.IsEmpty() {
		id = q.Queue[0]
		q.Queue = q.Queue[1:]
	} else {
		err = fmt.Errorf("Can't pop from empty id queue")
	}

	return
}

func (q *IDQueue) IsEmpty() (isEmpty bool) {
	isEmpty = (len(q.Queue) == 0)

	return
}
