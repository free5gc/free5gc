//go:binary-only-package

package pfcpType

import (
	"fmt"
)

const (
	OuterHeaderRemovalGtpUUdpIpv4 uint8 = iota
	OuterHeaderRemovalGtpUUdpIpv6
	OuterHeaderRemovalUdpIpv4
	OuterHeaderRemovalUdpIpv6
)

type OuterHeaderRemoval struct {
	OuterHeaderRemovalDescription uint8
}

func (o *OuterHeaderRemoval) MarshalBinary() (data []byte, err error) {}

func (o *OuterHeaderRemoval) UnmarshalBinary(data []byte) error {}
