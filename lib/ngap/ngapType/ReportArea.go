//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	ReportAreaPresentCell aper.Enumerated = 0
)

type ReportArea struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
