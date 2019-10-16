//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"time"
)

type EndTime struct {
	EndTime time.Time
}

func (e *EndTime) MarshalBinary() (data []byte, err error) {}

func (e *EndTime) UnmarshalBinary(data []byte) error {}
