//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"time"
)

type StartTime struct {
	StartTime time.Time
}

func (s *StartTime) MarshalBinary() (data []byte, err error) {}

func (s *StartTime) UnmarshalBinary(data []byte) error {}
