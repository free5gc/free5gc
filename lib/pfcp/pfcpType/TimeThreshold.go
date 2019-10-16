//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type TimeThreshold struct {
	TimeThreshold uint32
}

func (t *TimeThreshold) MarshalBinary() (data []byte, err error) {}

func (t *TimeThreshold) UnmarshalBinary(data []byte) error {}
