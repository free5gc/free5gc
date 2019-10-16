//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct RATRestrictions */
/* RATRestrictionsItem */
type RATRestrictions struct {
	List []RATRestrictionsItem `aper:"valueExt,sizeLB:16,sizeUB:16"`
}
