//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type URSEQN struct {
	UrseqnValue uint32
}

func (u *URSEQN) MarshalBinary() (data []byte, err error) {}

func (u *URSEQN) UnmarshalBinary(data []byte) error {}
