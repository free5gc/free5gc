//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
	"net"
)

type UserPlaneIPResourceInformation struct {
	Assosi          bool
	Assoni          bool
	Teidri          uint8 // 0x00011100
	V6              bool
	V4              bool
	TeidRange       uint8
	Ipv4Address     net.IP
	Ipv6Address     net.IP
	NetworkInstance []byte
	SourceInterface uint8 // 0x00001111
}

func (u *UserPlaneIPResourceInformation) MarshalBinary() (data []byte, err error) {}

func (u *UserPlaneIPResourceInformation) UnmarshalBinary(data []byte) error {}
