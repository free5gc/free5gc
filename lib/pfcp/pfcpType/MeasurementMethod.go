//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type MeasurementMethod struct {
	Event bool
	Volum bool
	Durat bool
}

func (m *MeasurementMethod) MarshalBinary() (data []byte, err error) {}

func (m *MeasurementMethod) UnmarshalBinary(data []byte) error {}
