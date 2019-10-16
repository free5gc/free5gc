//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type ForwardingPolicy struct {
	ForwardingPolicyIdentifierLength uint8
	ForwardingPolicyIdentifier       []byte
}

func (f *ForwardingPolicy) MarshalBinary() (data []byte, err error) {}

func (f *ForwardingPolicy) UnmarshalBinary(data []byte) error {}
