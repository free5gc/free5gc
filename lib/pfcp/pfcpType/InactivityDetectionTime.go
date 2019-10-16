//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type InactivityDetectionTime struct {
	InactivityDetectionTime uint32
}

func (i *InactivityDetectionTime) MarshalBinary() (data []byte, err error) {}

func (i *InactivityDetectionTime) UnmarshalBinary(data []byte) error {}
