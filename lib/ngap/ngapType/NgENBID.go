//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	NgENBIDPresentNothing int = iota /* No components present */
	NgENBIDPresentMacroNgENBID
	NgENBIDPresentShortMacroNgENBID
	NgENBIDPresentLongMacroNgENBID
	NgENBIDPresentChoiceExtensions
)

type NgENBID struct {
	Present           int
	MacroNgENBID      *aper.BitString `aper:"sizeLB:20,sizeUB:20"`
	ShortMacroNgENBID *aper.BitString `aper:"sizeLB:18,sizeUB:18"`
	LongMacroNgENBID  *aper.BitString `aper:"sizeLB:21,sizeUB:21"`
	ChoiceExtensions  *ProtocolIESingleContainerNgENBIDExtIEs
}
