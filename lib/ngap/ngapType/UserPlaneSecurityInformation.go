//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type UserPlaneSecurityInformation struct {
	SecurityResult     SecurityResult                                                `aper:"valueExt"`
	SecurityIndication SecurityIndication                                            `aper:"valueExt"`
	IEExtensions       *ProtocolExtensionContainerUserPlaneSecurityInformationExtIEs `aper:"optional"`
}
