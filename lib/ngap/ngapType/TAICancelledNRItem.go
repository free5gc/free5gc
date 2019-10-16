//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type TAICancelledNRItem struct {
	TAI                   TAI `aper:"valueExt"`
	CancelledCellsInTAINR CancelledCellsInTAINR
	IEExtensions          *ProtocolExtensionContainerTAICancelledNRItemExtIEs `aper:"optional"`
}
