//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
)

type HeaderEnrichment struct {
	HeaderType               uint8 // 0x00011111
	LengthOfHeaderFieldName  uint8
	HeaderFieldName          []byte
	LengthOfHeaderFieldValue uint8
	HeaderFieldValue         []byte
}

func (h *HeaderEnrichment) MarshalBinary() (data []byte, err error) {}

func (h *HeaderEnrichment) UnmarshalBinary(data []byte) error {}
