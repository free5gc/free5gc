//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type SubsequentVolumeThreshold struct {
	Dlvol          bool
	Ulvol          bool
	Tovol          bool
	TotalVolume    uint64
	UplinkVolume   uint64
	DownlinkVolume uint64
}

func (s *SubsequentVolumeThreshold) MarshalBinary() (data []byte, err error) {}

func (s *SubsequentVolumeThreshold) UnmarshalBinary(data []byte) error {}
