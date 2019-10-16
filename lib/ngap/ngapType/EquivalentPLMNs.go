//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct EquivalentPLMNs */
/* PLMNIdentity */
type EquivalentPLMNs struct {
	List []PLMNIdentity `aper:"sizeLB:1,sizeUB:15"`
}
