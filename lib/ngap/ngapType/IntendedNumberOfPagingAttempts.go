//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type IntendedNumberOfPagingAttempts struct {
	Value int64 `aper:"valueExt,valueLB:1,valueUB:16"`
}
