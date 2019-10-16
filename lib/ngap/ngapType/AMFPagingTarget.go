//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	AMFPagingTargetPresentNothing int = iota /* No components present */
	AMFPagingTargetPresentGlobalRANNodeID
	AMFPagingTargetPresentTAI
	AMFPagingTargetPresentChoiceExtensions
)

type AMFPagingTarget struct {
	Present          int
	GlobalRANNodeID  *GlobalRANNodeID `aper:"valueLB:0,valueUB:3"`
	TAI              *TAI             `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerAMFPagingTargetExtIEs
}
