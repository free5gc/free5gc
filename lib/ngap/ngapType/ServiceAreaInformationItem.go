//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type ServiceAreaInformationItem struct {
	PLMNIdentity   PLMNIdentity
	AllowedTACs    *AllowedTACs                                                `aper:"optional"`
	NotAllowedTACs *NotAllowedTACs                                             `aper:"optional"`
	IEExtensions   *ProtocolExtensionContainerServiceAreaInformationItemExtIEs `aper:"optional"`
}
