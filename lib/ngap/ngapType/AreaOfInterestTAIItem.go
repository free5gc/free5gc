//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type AreaOfInterestTAIItem struct {
	TAI          TAI                                                    `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerAreaOfInterestTAIItemExtIEs `aper:"optional"`
}
