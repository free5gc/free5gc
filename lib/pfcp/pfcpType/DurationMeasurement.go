//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type DurationMeasurement struct {
	DurationValue uint32
}

func (d *DurationMeasurement) MarshalBinary() (data []byte, err error) {}

func (d *DurationMeasurement) UnmarshalBinary(data []byte) error {}
