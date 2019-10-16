//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct CancelledCellsInEAI_EUTRA */
/* CancelledCellsInEAIEUTRAItem */
type CancelledCellsInEAIEUTRA struct {
	List []CancelledCellsInEAIEUTRAItem `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}
