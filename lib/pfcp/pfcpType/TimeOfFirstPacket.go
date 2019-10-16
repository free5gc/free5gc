//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"time"
)

type TimeOfFirstPacket struct {
	TimeOfFirstPacket time.Time
}

func (t *TimeOfFirstPacket) MarshalBinary() (data []byte, err error) {}

func (t *TimeOfFirstPacket) UnmarshalBinary(data []byte) error {}
