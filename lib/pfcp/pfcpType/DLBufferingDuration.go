//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
)

type DLBufferingDuration struct {
	TimerUnit  uint8 // 0x11100000
	TimerValue uint8 // 0x00011111
}

func (d *DLBufferingDuration) MarshalBinary() (data []byte, err error) {}

func (d *DLBufferingDuration) UnmarshalBinary(data []byte) error {}
