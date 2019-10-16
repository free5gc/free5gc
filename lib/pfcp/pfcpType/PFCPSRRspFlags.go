//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type PFCPSRRspFlags struct {
	Drobu bool
}

func (p *PFCPSRRspFlags) MarshalBinary() (data []byte, err error) {}

func (p *PFCPSRRspFlags) UnmarshalBinary(data []byte) error {}
