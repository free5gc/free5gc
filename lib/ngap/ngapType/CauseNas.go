//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	CauseNasPresentNormalRelease         aper.Enumerated = 0
	CauseNasPresentAuthenticationFailure aper.Enumerated = 1
	CauseNasPresentDeregister            aper.Enumerated = 2
	CauseNasPresentUnspecified           aper.Enumerated = 3
)

type CauseNas struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:3"`
}
