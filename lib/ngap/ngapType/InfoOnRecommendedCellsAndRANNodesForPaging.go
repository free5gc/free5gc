//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type InfoOnRecommendedCellsAndRANNodesForPaging struct {
	RecommendedCellsForPaging  RecommendedCellsForPaging                                                   `aper:"valueExt"`
	RecommendRANNodesForPaging RecommendedRANNodesForPaging                                                `aper:"valueExt"`
	IEExtensions               *ProtocolExtensionContainerInfoOnRecommendedCellsAndRANNodesForPagingExtIEs `aper:"optional"`
}
