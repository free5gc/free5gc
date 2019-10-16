//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type PacketDetectionRuleID struct {
	RuleId uint16
}

func (p *PacketDetectionRuleID) MarshalBinary() (data []byte, err error) {}

func (p *PacketDetectionRuleID) UnmarshalBinary(data []byte) error {}
