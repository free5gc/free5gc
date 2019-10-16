//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

type DLBufferingSuggestedPacketCount struct {
	PacketCountValue uint16
}

func (d *DLBufferingSuggestedPacketCount) MarshalBinary() (data []byte, err error) {}

func (d *DLBufferingSuggestedPacketCount) UnmarshalBinary(data []byte) error {}
