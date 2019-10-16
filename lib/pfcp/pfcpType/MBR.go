//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

type MBR struct {
	UlMbr uint64 // 40-bit data
	DlMbr uint64 // 40-bit data
}

func (m *MBR) MarshalBinary() (data []byte, err error) {}

func (m *MBR) UnmarshalBinary(data []byte) error {}
