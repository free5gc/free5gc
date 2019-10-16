//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	CellSizePresentVerysmall aper.Enumerated = 0
	CellSizePresentSmall     aper.Enumerated = 1
	CellSizePresentMedium    aper.Enumerated = 2
	CellSizePresentLarge     aper.Enumerated = 3
)

type CellSize struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:3"`
}
