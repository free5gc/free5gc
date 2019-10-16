//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"time"
)

type TimeOfLastPacket struct {
	TimeOfLastPacket time.Time
}

func (t *TimeOfLastPacket) MarshalBinary() (data []byte, err error) {}

func (t *TimeOfLastPacket) UnmarshalBinary(data []byte) error {}
