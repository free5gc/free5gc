//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type ReportingTriggers struct {
	Liusa bool
	Droth bool
	Stopt bool
	Start bool
	Quhti bool
	Timth bool
	Volth bool
	Perio bool
	Evequ bool
	Eveth bool
	Macar bool
	Envcl bool
	Timqu bool
	Volqu bool
}

func (r *ReportingTriggers) MarshalBinary() (data []byte, err error) {}

func (r *ReportingTriggers) UnmarshalBinary(data []byte) error {}
