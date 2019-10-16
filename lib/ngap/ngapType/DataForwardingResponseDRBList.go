//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct DataForwardingResponseDRBList */
/* DataForwardingResponseDRBItem */
type DataForwardingResponseDRBList struct {
	List []DataForwardingResponseDRBItem `aper:"valueExt,sizeLB:1,sizeUB:32"`
}
