//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct QosFlowList */
/* QosFlowItem */
type QosFlowList struct {
	List []QosFlowItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}
