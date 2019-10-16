//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type SubsequentTimeThreshold struct {
	SubsequentTimeThreshold uint32
}

func (s *SubsequentTimeThreshold) MarshalBinary() (data []byte, err error) {}

func (s *SubsequentTimeThreshold) UnmarshalBinary(data []byte) error {}
