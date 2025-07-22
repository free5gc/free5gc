package buffnetlink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"

	"github.com/free5gc/go-upf/internal/report"
)

type testHandler struct {
	q map[uint64]map[uint16]chan []byte
}

func NewTestHandler() *testHandler {
	return &testHandler{q: make(map[uint64]map[uint16]chan []byte)}
}

func (h *testHandler) Close() {
	for _, s := range h.q {
		for _, q := range s {
			close(q)
		}
	}
}

func (h *testHandler) NotifySessReport(sr report.SessReport) {
	s, ok := h.q[sr.SEID]
	if !ok {
		return
	}
	for _, rep := range sr.Reports {
		switch r := rep.(type) {
		case report.DLDReport:
			if r.Action&report.APPLY_ACT_BUFF != 0 && len(r.BufPkt) > 0 {
				q, ok := s[r.PDRID]
				if !ok {
					qlen := 10
					s[r.PDRID] = make(chan []byte, qlen)
					q = s[r.PDRID]
				}
				q <- r.BufPkt
			}
		default:
		}
	}
}

func (h *testHandler) PopBufPkt(seid uint64, pdrid uint16) ([]byte, bool) {
	s, ok := h.q[seid]
	if !ok {
		return nil, false
	}
	q, ok := s[pdrid]
	if !ok {
		return nil, false
	}
	select {
	case pkt := <-q:
		return pkt, true
	default:
		return nil, false
	}
}

func TestServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	var wg sync.WaitGroup

	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = mux.Serve()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
	}()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		t.Fatal(err)
	}

	c := nl.NewClient(conn, mux)
	s, err := OpenServer(&wg, c, mux)
	if err != nil {
		t.Fatal(err)
	}

	f, err := genl.GetFamily(c, "gtp5g")
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()

	h := NewTestHandler()
	defer func() {
		h.Close()
		s.Close()
		mux.Close()
		wg.Wait()
	}()

	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW|syscall.SOCK_CLOEXEC, syscall.NETLINK_GENERIC)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		errClose := syscall.Close(fd)
		if errClose != nil {
			t.Fatal(errClose)
		}
	}()

	seid := uint64(6)
	h.q[seid] = make(map[uint16]chan []byte)
	s.Handle(h)

	pkt := []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x10,
		0x00,
		0x00, 0x00,
		0x0c, 0x00, 0x06, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x03, 0x00,
		0x06, 0x00, 0x05, 0x00, 0x03, 0x00, 0x00, 0x00,
		0x06, 0x00, 0x07, 0x00, 0x0c, 0x00, 0x00, 0x00,
		0x08, 0x00, 0x04, 0x00,
		0xee, 0xbb,
		0xdd, 0xcc,
	}

	binary.LittleEndian.PutUint16(pkt[4:6], f.ID)
	binary.LittleEndian.PutUint32(pkt[0:4], uint32(len(pkt)))

	N := 10
	for i := 0; i < N; i++ {
		addr := syscall.SockaddrNetlink{
			Family: syscall.AF_NETLINK,
			Groups: 1 << (f.Groups[0].ID - 1),
		}
		err = syscall.Sendmsg(fd, pkt, nil, &addr, 0)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)

		pdrid := uint16(3)
		pkt, ok := s.Pop(seid, pdrid)
		if !ok {
			t.Fatal("not found")
		}

		want := []byte{0xee, 0xbb, 0xdd, 0xcc}
		if !bytes.Equal(pkt, want) {
			t.Errorf("want %x; but got %x\n", want, pkt)
		}

		_, ok = s.Pop(seid, pdrid)
		if ok {
			t.Fatal("found")
		}
	}
}
