//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type LastVisitedNGRANCellInformation struct {
	GlobalCellID                          NGRANCGI `aper:"valueLB:0,valueUB:2"`
	CellType                              CellType `aper:"valueExt"`
	TimeUEStayedInCell                    TimeUEStayedInCell
	TimeUEStayedInCellEnhancedGranularity *TimeUEStayedInCellEnhancedGranularity                           `aper:"optional"`
	HOCauseValue                          *Cause                                                           `aper:"valueLB:0,valueUB:5,optional"`
	IEExtensions                          *ProtocolExtensionContainerLastVisitedNGRANCellInformationExtIEs `aper:"optional"`
}
