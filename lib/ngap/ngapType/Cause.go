//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	CausePresentNothing int = iota /* No components present */
	CausePresentRadioNetwork
	CausePresentTransport
	CausePresentNas
	CausePresentProtocol
	CausePresentMisc
	CausePresentChoiceExtensions
)

type Cause struct {
	Present          int
	RadioNetwork     *CauseRadioNetwork
	Transport        *CauseTransport
	Nas              *CauseNas
	Protocol         *CauseProtocol
	Misc             *CauseMisc
	ChoiceExtensions *ProtocolIESingleContainerCauseExtIEs
}
