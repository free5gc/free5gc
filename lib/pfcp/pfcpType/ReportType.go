//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type ReportType struct {
	Upir bool
	Erir bool
	Usar bool
	Dldr bool
}

func (r *ReportType) MarshalBinary() (data []byte, err error) {}

func (r *ReportType) UnmarshalBinary(data []byte) error {}
