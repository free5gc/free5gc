//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type FARID struct {
	FarIdValue uint32
}

func (f *FARID) MarshalBinary() (data []byte, err error) {}

func (f *FARID) UnmarshalBinary(data []byte) error {}
