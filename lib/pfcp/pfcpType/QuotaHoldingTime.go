//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type QuotaHoldingTime struct {
	QuotaHoldingTimeValue uint32
}

func (q *QuotaHoldingTime) MarshalBinary() (data []byte, err error) {}

func (q *QuotaHoldingTime) UnmarshalBinary(data []byte) error {}
