//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	DRBStatusDLPresentNothing int = iota /* No components present */
	DRBStatusDLPresentDRBStatusDL12
	DRBStatusDLPresentDRBStatusDL18
	DRBStatusDLPresentChoiceExtensions
)

type DRBStatusDL struct {
	Present          int
	DRBStatusDL12    *DRBStatusDL12 `aper:"valueExt"`
	DRBStatusDL18    *DRBStatusDL18 `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerDRBStatusDLExtIEs
}
