//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type SuggestedBufferingPacketsCount struct {
	PacketCountValue uint8
}

func (s *SuggestedBufferingPacketsCount) MarshalBinary() (data []byte, err error) {}

func (s *SuggestedBufferingPacketsCount) UnmarshalBinary(data []byte) error {}
