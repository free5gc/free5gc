//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceHandoverList */
/* PDUSessionResourceHandoverItem */
type PDUSessionResourceHandoverList struct {
	List []PDUSessionResourceHandoverItem `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
