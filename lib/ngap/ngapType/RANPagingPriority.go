//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type RANPagingPriority struct {
	Value int64 `aper:"valueLB:1,valueUB:256"`
}
