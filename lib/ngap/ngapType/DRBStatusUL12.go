//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

type DRBStatusUL12 struct {
	ULCOUNTValue              COUNTValueForPDCPSN12                          `aper:"valueExt"`
	ReceiveStatusOfULPDCPSDUs *aper.BitString                                `aper:"sizeLB:1,sizeUB:2048,optional"`
	IEExtension               *ProtocolExtensionContainerDRBStatusUL12ExtIEs `aper:"optional"`
}
