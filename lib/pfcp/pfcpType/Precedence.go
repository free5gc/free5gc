//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type Precedence struct {
	PrecedenceValue uint32
}

func (p *Precedence) MarshalBinary() (data []byte, err error) {}

func (p *Precedence) UnmarshalBinary(data []byte) error {}
