//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type EthernetInactivityTimer struct {
	EthernetInactivityTimer uint32
}

func (e *EthernetInactivityTimer) MarshalBinary() (data []byte, err error) {}

func (e *EthernetInactivityTimer) UnmarshalBinary(data []byte) error {}
