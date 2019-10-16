//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"net"
)

type FSEID struct {
	V4          bool
	V6          bool
	Seid        uint64
	Ipv4Address net.IP
	Ipv6Address net.IP
}

func (f *FSEID) MarshalBinary() (data []byte, err error) {}

func (f *FSEID) UnmarshalBinary(data []byte) error {}
