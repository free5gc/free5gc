//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
)

type GateStatus struct {
	UlGate uint8 // 0x00001100
	DlGate uint8 // 0x00000011
}

func (g *GateStatus) MarshalBinary() (data []byte, err error) {}

func (g *GateStatus) UnmarshalBinary(data []byte) error {}
