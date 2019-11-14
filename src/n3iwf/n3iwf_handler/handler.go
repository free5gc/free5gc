package n3iwf_handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_handler/n3iwf_message"
)

var HandlerLog *logrus.Entry

func init() {
	// init pool
	HandlerLog = logger.HandlerLog
}

func Handle() {
	for {
		msg, ok := <-n3iwf_message.N3iwfChannel
		if ok {
			fmt.Println(string(msg.Value.([]byte)))
		}
	}
}
