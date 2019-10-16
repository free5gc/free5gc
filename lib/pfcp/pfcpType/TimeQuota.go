//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type TimeQuota struct {
	TimeQuotaValue uint32
}

func (t *TimeQuota) MarshalBinary() (data []byte, err error) {}

func (t *TimeQuota) UnmarshalBinary(data []byte) error {}
