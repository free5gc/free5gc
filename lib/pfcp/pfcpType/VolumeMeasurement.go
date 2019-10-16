//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type VolumeMeasurement struct {
	Dlvol          bool
	Ulvol          bool
	Tovol          bool
	TotalVolume    uint64
	UplinkVolume   uint64
	DownlinkVolume uint64
}

func (v *VolumeMeasurement) MarshalBinary() (data []byte, err error) {}

func (v *VolumeMeasurement) UnmarshalBinary(data []byte) error {}
