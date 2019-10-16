//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type COUNTValueForPDCPSN18 struct {
	PDCPSN18     int64                                                  `aper:"valueLB:0,valueUB:262143"`
	HFNPDCPSN18  int64                                                  `aper:"valueLB:0,valueUB:16383"`
	IEExtensions *ProtocolExtensionContainerCOUNTValueForPDCPSN18ExtIEs `aper:"optional"`
}
