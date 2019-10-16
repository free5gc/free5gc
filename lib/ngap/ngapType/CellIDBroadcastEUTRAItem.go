//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type CellIDBroadcastEUTRAItem struct {
	EUTRACGI     EUTRACGI                                                  `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCellIDBroadcastEUTRAItemExtIEs `aper:"optional"`
}
