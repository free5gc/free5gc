//go:binary-only-package

package pfcpType

import (
	"fmt"
)

type ApplyAction struct {
	Dupl bool
	Nocp bool
	Buff bool
	Forw bool
	Drop bool
}

func (a *ApplyAction) MarshalBinary() (data []byte, err error) {}

func (a *ApplyAction) UnmarshalBinary(data []byte) error {}
