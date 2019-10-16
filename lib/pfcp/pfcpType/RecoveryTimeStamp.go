//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"time"
)

type RecoveryTimeStamp struct {
	RecoveryTimeStamp time.Time
}

func (r *RecoveryTimeStamp) MarshalBinary() (data []byte, err error) {}

func (r *RecoveryTimeStamp) UnmarshalBinary(data []byte) error {}
