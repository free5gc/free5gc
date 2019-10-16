//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type ProtocolExtensionID struct {
	Value int64 `aper:"valueLB:0,valueUB:65535"`
}
