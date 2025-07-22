package app

import (
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/free5gc/go-upf/pkg/factory"
)

func TestWaitRoutineStopped(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	cfg := &factory.Config{
		Version: "1.0.3",
		Pfcp: &factory.Pfcp{
			Addr:   "127.0.0.1",
			NodeID: "127.0.0.1",
		},
		Gtpu: &factory.Gtpu{
			Forwarder: "gtp5g",
			IfList: []factory.IfInfo{
				{
					Addr:   "127.0.0.1",
					Type:   "",
					Name:   "",
					IfName: "",
				},
			},
		},
		DnnList: []factory.DnnList{
			{
				Dnn:  "internet",
				Cidr: "10.60.0.1/24",
			},
		},
		Logger: &factory.Logger{
			Enable: true,
			Level:  "info",
		},
	}
	N := 10
	for i := 0; i < N; i++ {
		var wg sync.WaitGroup
		upf, err := NewApp(cfg)
		if err != nil {
			t.Fatal(err)
		}
		wg.Add(1)
		go func() {
			upf.Start()
			wg.Done()
		}()
		// Must wait for signal initialized
		time.Sleep(500 * time.Millisecond)
		err = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		if err != nil {
			t.Fatal(err)
		}
		wg.Wait()
	}
}
