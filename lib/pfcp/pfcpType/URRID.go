//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type URRID struct {
	UrrIdValue uint32
}

func (u *URRID) MarshalBinary() (data []byte, err error) {}

func (u *URRID) UnmarshalBinary(data []byte) error {}
