//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
)

const (
	DestinationInterfaceAccess uint8 = iota
	DestinationInterfaceCore
	DestinationInterfaceSgiLanN6Lan
	DestinationInterfaceCpFunction
	DestinationInterfaceLiFunction
)

type DestinationInterface struct {
	InterfaceValue uint8 // 0x00001111
}

func (d *DestinationInterface) MarshalBinary() (data []byte, err error) {}

func (d *DestinationInterface) UnmarshalBinary(data []byte) error {}
