//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"time"
)

type MonitoringTime struct {
	MonitoringTime time.Time
}

func (m *MonitoringTime) MarshalBinary() (data []byte, err error) {}

func (m *MonitoringTime) UnmarshalBinary(data []byte) error {}
