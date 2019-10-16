//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type TargeteNBID struct {
	GlobalENBID    GlobalNgENBID                                `aper:"valueExt"`
	SelectedEPSTAI EPSTAI                                       `aper:"valueExt"`
	IEExtensions   *ProtocolExtensionContainerTargeteNBIDExtIEs `aper:"optional"`
}
