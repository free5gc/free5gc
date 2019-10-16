//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type MeasurementInformation struct {
	Radi bool
	Inam bool
	Mbqe bool
}

func (m *MeasurementInformation) MarshalBinary() (data []byte, err error) {}

func (m *MeasurementInformation) UnmarshalBinary(data []byte) error {}
