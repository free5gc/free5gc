//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type EventThreshold struct {
	EventThreshold uint32
}

func (e *EventThreshold) MarshalBinary() (data []byte, err error) {}

func (e *EventThreshold) UnmarshalBinary(data []byte) error {}
