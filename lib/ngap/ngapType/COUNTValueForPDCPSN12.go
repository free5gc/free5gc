//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type COUNTValueForPDCPSN12 struct {
	PDCPSN12     int64                                                  `aper:"valueLB:0,valueUB:4095"`
	HFNPDCPSN12  int64                                                  `aper:"valueLB:0,valueUB:1048575"`
	IEExtensions *ProtocolExtensionContainerCOUNTValueForPDCPSN12ExtIEs `aper:"optional"`
}
