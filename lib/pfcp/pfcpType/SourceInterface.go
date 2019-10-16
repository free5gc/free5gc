//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
)

const (
	SourceInterfaceAccess uint8 = iota
	SourceInterfaceCore
	SourceInterfaceSgiLanN6Lan
	SourceInterfaceCpFunction
)

type SourceInterface struct {
	InterfaceValue uint8 // 0x00001111
}

func (s *SourceInterface) MarshalBinary() (data []byte, err error) {}

func (s *SourceInterface) UnmarshalBinary(data []byte) error {}
