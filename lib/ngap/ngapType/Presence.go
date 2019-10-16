//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	PresencePresentOptional    aper.Enumerated = 0
	PresencePresentConditional aper.Enumerated = 1
	PresencePresentMandatory   aper.Enumerated = 2
)

type Presence struct {
	Value aper.Enumerated `aper:"valueLB:0,valueUB:2"`
}
