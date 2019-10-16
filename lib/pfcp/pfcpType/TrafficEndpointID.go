//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type TrafficEndpointID struct {
	TrafficEndpointIdValue uint8
}

func (t *TrafficEndpointID) MarshalBinary() (data []byte, err error) {}

func (t *TrafficEndpointID) UnmarshalBinary(data []byte) error {}
