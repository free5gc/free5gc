//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type DroppedDLTrafficThreshold struct {
	Dlby                        bool
	Dlpa                        bool
	DownlinkPackets             uint64
	NumberOfBytesOfDownlinkData uint64
}

func (d *DroppedDLTrafficThreshold) MarshalBinary() (data []byte, err error) {}

func (d *DroppedDLTrafficThreshold) UnmarshalBinary(data []byte) error {}
