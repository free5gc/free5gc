//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	DRBStatusULPresentNothing int = iota /* No components present */
	DRBStatusULPresentDRBStatusUL12
	DRBStatusULPresentDRBStatusUL18
	DRBStatusULPresentChoiceExtensions
)

type DRBStatusUL struct {
	Present          int
	DRBStatusUL12    *DRBStatusUL12 `aper:"valueExt"`
	DRBStatusUL18    *DRBStatusUL18 `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerDRBStatusULExtIEs
}
