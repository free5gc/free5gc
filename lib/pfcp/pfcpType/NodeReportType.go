//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type NodeReportType struct {
	Upfr bool
}

func (n *NodeReportType) MarshalBinary() (data []byte, err error) {}

func (n *NodeReportType) UnmarshalBinary(data []byte) error {}
