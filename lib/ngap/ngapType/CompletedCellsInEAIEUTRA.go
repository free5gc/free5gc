//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct CompletedCellsInEAI_EUTRA */
/* CompletedCellsInEAIEUTRAItem */
type CompletedCellsInEAIEUTRA struct {
	List []CompletedCellsInEAIEUTRAItem `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}
