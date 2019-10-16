//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type VolumeQuota struct {
	Dlvol          bool
	Ulvol          bool
	Tovol          bool
	TotalVolume    uint64
	UplinkVolume   uint64
	DownlinkVolume uint64
}

func (v *VolumeQuota) MarshalBinary() (data []byte, err error) {}

func (v *VolumeQuota) UnmarshalBinary(data []byte) error {}
