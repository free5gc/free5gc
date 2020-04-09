package udp_server

import (
	"sync"
	"testing"

	"free5gc/src/n3iwf/n3iwf_handler"
)

func TestServer(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(2)

	go Run()
	go n3iwf_handler.Handle()

	wg.Wait()

}
