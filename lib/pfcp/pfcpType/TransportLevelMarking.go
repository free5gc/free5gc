//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type TransportLevelMarking struct {
	TosTrafficClass []byte
}

func (t *TransportLevelMarking) MarshalBinary() (data []byte, err error) {}

func (t *TransportLevelMarking) UnmarshalBinary(data []byte) error {}
