//go:binary-only-package

package pfcpUdp

import (
	"net"

	"free5gc/lib/pfcp"
)

const (
	PFCP_PORT        = 8805
	PFCP_MAX_UDP_LEN = 2048
)

type PfcpServer struct {
	Addr string
	Conn *net.UDPConn
}

func NewPfcpServer(addr string) (PfcpServer, error) {}

func (p *PfcpServer) Listen() error {}

func (p *PfcpServer) ReadFrom(msg *pfcp.Message) (*net.UDPAddr, error) {}

func (p *PfcpServer) WriteTo(msg pfcp.Message, addr *net.UDPAddr) error {}

func (p *PfcpServer) Close() error {}

// Send a PFCP message and close UDP connection
func SendPfcpMessage(msg pfcp.Message, srcAddr *net.UDPAddr, dstAddr *net.UDPAddr) error {}

// Receive a PFCP message and close UDP connection
func ReceivePfcpMessage(msg *pfcp.Message, srcAddr *net.UDPAddr, dstAddr *net.UDPAddr) error {}
