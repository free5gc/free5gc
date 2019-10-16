//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	BroadcastCompletedAreaListPresentNothing int = iota /* No components present */
	BroadcastCompletedAreaListPresentCellIDBroadcastEUTRA
	BroadcastCompletedAreaListPresentTAIBroadcastEUTRA
	BroadcastCompletedAreaListPresentEmergencyAreaIDBroadcastEUTRA
	BroadcastCompletedAreaListPresentCellIDBroadcastNR
	BroadcastCompletedAreaListPresentTAIBroadcastNR
	BroadcastCompletedAreaListPresentEmergencyAreaIDBroadcastNR
	BroadcastCompletedAreaListPresentChoiceExtensions
)

type BroadcastCompletedAreaList struct {
	Present                       int
	CellIDBroadcastEUTRA          *CellIDBroadcastEUTRA
	TAIBroadcastEUTRA             *TAIBroadcastEUTRA
	EmergencyAreaIDBroadcastEUTRA *EmergencyAreaIDBroadcastEUTRA
	CellIDBroadcastNR             *CellIDBroadcastNR
	TAIBroadcastNR                *TAIBroadcastNR
	EmergencyAreaIDBroadcastNR    *EmergencyAreaIDBroadcastNR
	ChoiceExtensions              *ProtocolIESingleContainerBroadcastCompletedAreaListExtIEs
}
