//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type LinkedURRID struct {
	LinkedUrrIdValue uint32
}

func (l *LinkedURRID) MarshalBinary() (data []byte, err error) {}

func (l *LinkedURRID) UnmarshalBinary(data []byte) error {}
