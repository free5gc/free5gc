//go:binary-only-package

package pfcpType

import (
	"fmt"
	"net"
)

type UEIPAddress struct {
	Ipv6d                    bool
	Sd                       bool
	V4                       bool
	V6                       bool
	Ipv4Address              net.IP
	Ipv6Address              net.IP
	Ipv6PrefixDelegationBits uint8
}

func (u *UEIPAddress) MarshalBinary() (data []byte, err error) {}

func (u *UEIPAddress) UnmarshalBinary(data []byte) error {}
