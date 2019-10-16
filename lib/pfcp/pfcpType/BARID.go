//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type BARID struct {
	BarIdValue uint8
}

func (b *BARID) MarshalBinary() (data []byte, err error) {}

func (b *BARID) UnmarshalBinary(data []byte) error {}
