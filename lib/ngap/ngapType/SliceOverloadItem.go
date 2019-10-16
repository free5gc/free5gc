//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type SliceOverloadItem struct {
	SNSSAI       SNSSAI                                             `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerSliceOverloadItemExtIEs `aper:"optional"`
}
