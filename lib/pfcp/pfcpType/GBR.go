//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

type GBR struct {
	UlGbr uint64 // 40-bit data
	DlGbr uint64 // 40-bit data
}

func (m *GBR) MarshalBinary() (data []byte, err error) {}

func (m *GBR) UnmarshalBinary(data []byte) error {}
