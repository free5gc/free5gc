//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type DownlinkDataNotificationDelay struct {
	DelayValue uint8
}

func (d *DownlinkDataNotificationDelay) MarshalBinary() (data []byte, err error) {}

func (d *DownlinkDataNotificationDelay) UnmarshalBinary(data []byte) error {}
