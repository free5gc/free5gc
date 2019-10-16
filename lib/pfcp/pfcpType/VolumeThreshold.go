//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type VolumeThreshold struct {
	Dlvol          bool
	Ulvol          bool
	Tovol          bool
	TotalVolume    uint64
	UplinkVolume   uint64
	DownlinkVolume uint64
}

func (v *VolumeThreshold) MarshalBinary() (data []byte, err error) {}

func (v *VolumeThreshold) UnmarshalBinary(data []byte) error {}
