//go:binary-only-package

package pfcpType

import (
	"fmt"
	"math/bits"
)

type DownlinkDataServiceInformation struct {
	Qfii                        bool
	Ppi                         bool
	PagingPolicyIndicationValue uint8 // 0x00111111
	Qfi                         uint8 // 0x00111111
}

func (d *DownlinkDataServiceInformation) MarshalBinary() (data []byte, err error) {}

func (d *DownlinkDataServiceInformation) UnmarshalBinary(data []byte) error {}
