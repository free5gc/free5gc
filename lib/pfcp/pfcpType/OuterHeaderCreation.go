//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	OuterHeaderCreationGtpUUdpIpv4 uint16 = 1
	OuterHeaderCreationGtpUUdpIpv6 uint16 = 1 << 1
	OuterHeaderCreationUdpIpv4     uint16 = 1 << 2
	OuterHeaderCreationUdpIpv6     uint16 = 1 << 3
)

type OuterHeaderCreation struct {
	OuterHeaderCreationDescription uint16
	Teid                           uint32
	Ipv4Address                    net.IP
	Ipv6Address                    net.IP
	PortNumber                     uint16
}

func (o *OuterHeaderCreation) MarshalBinary() (data []byte, err error) {}

func (o *OuterHeaderCreation) UnmarshalBinary(data []byte) error {}
