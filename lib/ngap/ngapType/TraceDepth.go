//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	TraceDepthPresentMinimum                               aper.Enumerated = 0
	TraceDepthPresentMedium                                aper.Enumerated = 1
	TraceDepthPresentMaximum                               aper.Enumerated = 2
	TraceDepthPresentMinimumWithoutVendorSpecificExtension aper.Enumerated = 3
	TraceDepthPresentMediumWithoutVendorSpecificExtension  aper.Enumerated = 4
	TraceDepthPresentMaximumWithoutVendorSpecificExtension aper.Enumerated = 5
)

type TraceDepth struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:5"`
}
