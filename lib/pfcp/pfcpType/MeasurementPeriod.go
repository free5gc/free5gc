//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type MeasurementPeriod struct {
	MeasurementPeriod uint32
}

func (m *MeasurementPeriod) MarshalBinary() (data []byte, err error) {}

func (m *MeasurementPeriod) UnmarshalBinary(data []byte) error {}
