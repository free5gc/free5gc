//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

type RedirectInformation struct {
	RedirectAddressType         uint8 // 0x00001111
	RedirectServerAddressLength uint16
	RedirectServerAddress       []byte
}

func (r *RedirectInformation) MarshalBinary() (data []byte, err error) {}

func (r *RedirectInformation) UnmarshalBinary(data []byte) error {}
