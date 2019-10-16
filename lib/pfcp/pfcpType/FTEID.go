//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"net"
)

type FTEID struct {
	Chid        bool
	Ch          bool
	V6          bool
	V4          bool
	Teid        uint32
	Ipv4Address net.IP
	Ipv6Address net.IP
	ChooseId    uint8
}

func (f *FTEID) MarshalBinary() (data []byte, err error) {}

func (f *FTEID) UnmarshalBinary(data []byte) error {}
