//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type EventID struct {
	EventId uint32
}

func (e *EventID) MarshalBinary() (data []byte, err error) {}

func (e *EventID) UnmarshalBinary(data []byte) error {}
