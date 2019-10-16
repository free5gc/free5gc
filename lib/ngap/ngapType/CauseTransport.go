//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	CauseTransportPresentTransportResourceUnavailable aper.Enumerated = 0
	CauseTransportPresentUnspecified                  aper.Enumerated = 1
)

type CauseTransport struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
