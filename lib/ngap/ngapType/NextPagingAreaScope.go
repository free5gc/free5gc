//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	NextPagingAreaScopePresentSame    aper.Enumerated = 0
	NextPagingAreaScopePresentChanged aper.Enumerated = 1
)

type NextPagingAreaScope struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
