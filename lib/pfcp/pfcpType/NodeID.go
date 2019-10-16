//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
	"net"
)

const (
	NodeIdTypeIpv4Address uint8 = iota
	NodeIdTypeIpv6Address
	NodeIdTypeFqdn
)

type NodeID struct {
	NodeIdType  uint8 // 0x00001111
	NodeIdValue []byte
}

func (n *NodeID) MarshalBinary() (data []byte, err error) {}

func (n *NodeID) UnmarshalBinary(data []byte) error {}

func (n *NodeID) ResolveNodeIdToIp() net.IP {}
