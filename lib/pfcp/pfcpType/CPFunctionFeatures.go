//go:binary-only-package

package pfcpType

import (
	"fmt"
)

const (
	CpFunctionFeaturesLoad uint8 = 1
	CpFunctionFeaturesOvrl uint8 = 1 << 1
)

type CPFunctionFeatures struct {
	SupportedFeatures uint8
}

func (c *CPFunctionFeatures) MarshalBinary() (data []byte, err error) {}

func (c *CPFunctionFeatures) UnmarshalBinary(data []byte) error {}
