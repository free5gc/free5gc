//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type UEPresenceInAreaOfInterestItem struct {
	LocationReportingReferenceID LocationReportingReferenceID
	UEPresence                   UEPresence
	IEExtensions                 *ProtocolExtensionContainerUEPresenceInAreaOfInterestItemExtIEs `aper:"optional"`
}
