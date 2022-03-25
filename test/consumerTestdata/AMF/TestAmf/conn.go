package TestAmf

import (
	"fmt"
	"net"
	"time"
)

type TestAddr struct {
	Net   string
	Value string
}

func (addr *TestAddr) Network() string {
	return addr.Net
}
func (addr *TestAddr) String() string {
	return addr.Value
}

type TestConn struct {
	RAddr *TestAddr
	LAddr *TestAddr
}

// Read implements the Conn Read method.
func (c *TestConn) Read(b []byte) (int, error) {
	return 0, fmt.Errorf("Read not support for Test\n")
}

// Write implements the Conn Write method.
func (c *TestConn) Write(b []byte) (int, error) {
	return 0, fmt.Errorf("Write not support for Test\n")
}

// Close closes the connection.
func (c *TestConn) Close() error {
	return fmt.Errorf("Close not support for Test\n")
}

func (c *TestConn) LocalAddr() net.Addr {
	return c.LAddr
}

func (c *TestConn) RemoteAddr() net.Addr {
	return c.RAddr
}

func (c *TestConn) SetDeadline(t time.Time) error {
	return fmt.Errorf("SetDeadline not support for Test\n")
}

func (c *TestConn) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("SetReadDeadline not support for Test\n")
}

func (c *TestConn) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("SetWriteDeadline not support for Test\n")
}
