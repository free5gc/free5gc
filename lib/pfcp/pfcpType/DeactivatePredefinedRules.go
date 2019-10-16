//go:binary-only-package

package pfcpType

type DeactivatePredefinedRules struct {
	PredefinedRulesName []byte
}

func (d *DeactivatePredefinedRules) MarshalBinary() (data []byte, err error) {}

func (d *DeactivatePredefinedRules) UnmarshalBinary(data []byte) error {}
