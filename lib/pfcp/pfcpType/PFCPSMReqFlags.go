//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type PFCPSMReqFlags struct {
	Qaurr bool
	Sndem bool
	Drobu bool
}

func (p *PFCPSMReqFlags) MarshalBinary() (data []byte, err error) {}

func (p *PFCPSMReqFlags) UnmarshalBinary(data []byte) error {}
