//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type OffendingIE struct {
	TypeOfOffendingIe uint16
}

func (o *OffendingIE) MarshalBinary() (data []byte, err error) {}

func (o *OffendingIE) UnmarshalBinary(data []byte) error {}
