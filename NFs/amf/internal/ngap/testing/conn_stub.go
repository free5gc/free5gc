package testing

import (
	"fmt"
	"net"
	"time"
)

type SctpConnStub struct {
	MsgList [][]byte
}

func (c *SctpConnStub) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (c *SctpConnStub) Write(b []byte) (n int, err error) {
	c.MsgList = append(c.MsgList, b)
	return 0, nil
}

func (c *SctpConnStub) Close() error {
	return fmt.Errorf("close error")
}

func (c *SctpConnStub) SetDeadline(time.Time) error {
	return fmt.Errorf("SetDeadline error")
}

func (c *SctpConnStub) SetReadDeadline(time.Time) error {
	return fmt.Errorf("SetReadDeadline error")
}

func (c *SctpConnStub) SetWriteDeadline(time.Time) error {
	return fmt.Errorf("SetWriteDeadline error")
}

func (c *SctpConnStub) LocalAddr() net.Addr {
	addr := new(testAddr)
	return addr
}

func (c *SctpConnStub) RemoteAddr() net.Addr {
	addr := new(testAddr)
	return addr
}

type testAddr struct{}

func (fa *testAddr) Network() string {
	return "sctp"
}

func (fa *testAddr) String() string {
	return "127.0.0.1"
}
