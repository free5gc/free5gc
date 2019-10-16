//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type SDFFilter struct {
	Bid                     bool
	Fl                      bool
	Spi                     bool
	Ttc                     bool
	Fd                      bool
	LengthOfFlowDescription uint16
	FlowDescription         []byte
	TosTrafficClass         []byte
	SecurityParameterIndex  []byte
	FlowLabel               []byte
	SdfFilterId             uint32
}

func (s *SDFFilter) MarshalBinary() (data []byte, err error) {}

func (s *SDFFilter) UnmarshalBinary(data []byte) error {}
