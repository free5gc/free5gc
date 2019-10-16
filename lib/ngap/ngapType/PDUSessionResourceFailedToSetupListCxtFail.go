//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceFailedToSetupListCxtFail */
/* PDUSessionResourceFailedToSetupItemCxtFail */
type PDUSessionResourceFailedToSetupListCxtFail struct {
	List []PDUSessionResourceFailedToSetupItemCxtFail `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
