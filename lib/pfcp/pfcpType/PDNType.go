//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
)

const (
	PDNTypeIpv4 uint8 = iota + 1
	PDNTypeIpv6
	PDNTypeIpv4v6
	PDNTypeNonIp
	PDNTypeEthernet
)

type PDNType struct {
	PdnType uint8 // 0x00000111
}

func (p *PDNType) MarshalBinary() (data []byte, err error) {}

func (p *PDNType) UnmarshalBinary(data []byte) error {}
