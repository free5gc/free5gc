//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
)

type FailedRuleID struct {
	RuleIdType  uint8 // 0x00001111
	RuleIdValue []byte
}

func (f *FailedRuleID) MarshalBinary() (data []byte, err error) {}

func (f *FailedRuleID) UnmarshalBinary(data []byte) error {}
