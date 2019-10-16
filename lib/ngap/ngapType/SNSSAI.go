//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type SNSSAI struct {
	SST          SST
	SD           *SD                                     `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerSNSSAIExtIEs `aper:"optional"`
}
