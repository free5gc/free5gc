//go:binary-only-package

package pfcpType

import (
	"fmt"
	"net"
)

type RemoteGTPUPeer struct {
	V4          bool
	V6          bool
	Ipv4Address net.IP
	Ipv6Address net.IP
}

func (r *RemoteGTPUPeer) MarshalBinary() (data []byte, err error) {}

func (r *RemoteGTPUPeer) UnmarshalBinary(data []byte) error {}
