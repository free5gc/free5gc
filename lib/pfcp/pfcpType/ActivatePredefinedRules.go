//go:binary-only-package

package pfcpType

type ActivatePredefinedRules struct {
	PredefinedRulesName []byte
}

func (a *ActivatePredefinedRules) MarshalBinary() (data []byte, err error) {}

func (a *ActivatePredefinedRules) UnmarshalBinary(data []byte) error {}
