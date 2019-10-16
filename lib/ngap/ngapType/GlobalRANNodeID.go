//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	GlobalRANNodeIDPresentNothing int = iota /* No components present */
	GlobalRANNodeIDPresentGlobalGNBID
	GlobalRANNodeIDPresentGlobalNgENBID
	GlobalRANNodeIDPresentGlobalN3IWFID
	GlobalRANNodeIDPresentChoiceExtensions
)

type GlobalRANNodeID struct {
	Present          int
	GlobalGNBID      *GlobalGNBID   `aper:"valueExt"`
	GlobalNgENBID    *GlobalNgENBID `aper:"valueExt"`
	GlobalN3IWFID    *GlobalN3IWFID `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerGlobalRANNodeIDExtIEs
}
