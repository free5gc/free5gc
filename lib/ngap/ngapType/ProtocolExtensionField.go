//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type AllocationAndRetentionPriorityExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AllocationAndRetentionPriorityExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AllocationAndRetentionPriorityExtIEsPresentNothing int = iota /* No components present */
)

type AllocationAndRetentionPriorityExtIEsExtensionValue struct {
	Present int
}

type AllowedNSSAIItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AllowedNSSAIItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AllowedNSSAIItemExtIEsPresentNothing int = iota /* No components present */
)

type AllowedNSSAIItemExtIEsExtensionValue struct {
	Present int
}

type AMFTNLAssociationSetupItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AMFTNLAssociationSetupItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFTNLAssociationSetupItemExtIEsPresentNothing int = iota /* No components present */
)

type AMFTNLAssociationSetupItemExtIEsExtensionValue struct {
	Present int
}

type AMFTNLAssociationToAddItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AMFTNLAssociationToAddItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFTNLAssociationToAddItemExtIEsPresentNothing int = iota /* No components present */
)

type AMFTNLAssociationToAddItemExtIEsExtensionValue struct {
	Present int
}

type AMFTNLAssociationToRemoveItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AMFTNLAssociationToRemoveItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFTNLAssociationToRemoveItemExtIEsPresentNothing int = iota /* No components present */
)

type AMFTNLAssociationToRemoveItemExtIEsExtensionValue struct {
	Present int
}

type AMFTNLAssociationToUpdateItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AMFTNLAssociationToUpdateItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFTNLAssociationToUpdateItemExtIEsPresentNothing int = iota /* No components present */
)

type AMFTNLAssociationToUpdateItemExtIEsExtensionValue struct {
	Present int
}

type AreaOfInterestExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AreaOfInterestExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AreaOfInterestExtIEsPresentNothing int = iota /* No components present */
)

type AreaOfInterestExtIEsExtensionValue struct {
	Present int
}

type AreaOfInterestCellItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AreaOfInterestCellItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AreaOfInterestCellItemExtIEsPresentNothing int = iota /* No components present */
)

type AreaOfInterestCellItemExtIEsExtensionValue struct {
	Present int
}

type AreaOfInterestItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AreaOfInterestItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AreaOfInterestItemExtIEsPresentNothing int = iota /* No components present */
)

type AreaOfInterestItemExtIEsExtensionValue struct {
	Present int
}

type AreaOfInterestRANNodeItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AreaOfInterestRANNodeItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AreaOfInterestRANNodeItemExtIEsPresentNothing int = iota /* No components present */
)

type AreaOfInterestRANNodeItemExtIEsExtensionValue struct {
	Present int
}

type AreaOfInterestTAIItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AreaOfInterestTAIItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AreaOfInterestTAIItemExtIEsPresentNothing int = iota /* No components present */
)

type AreaOfInterestTAIItemExtIEsExtensionValue struct {
	Present int
}

type AssistanceDataForPagingExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AssistanceDataForPagingExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AssistanceDataForPagingExtIEsPresentNothing int = iota /* No components present */
)

type AssistanceDataForPagingExtIEsExtensionValue struct {
	Present int
}

type AssistanceDataForRecommendedCellsExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AssistanceDataForRecommendedCellsExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AssistanceDataForRecommendedCellsExtIEsPresentNothing int = iota /* No components present */
)

type AssistanceDataForRecommendedCellsExtIEsExtensionValue struct {
	Present int
}

type AssociatedQosFlowItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue AssociatedQosFlowItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AssociatedQosFlowItemExtIEsPresentNothing int = iota /* No components present */
)

type AssociatedQosFlowItemExtIEsExtensionValue struct {
	Present int
}

type BroadcastPLMNItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue BroadcastPLMNItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	BroadcastPLMNItemExtIEsPresentNothing int = iota /* No components present */
)

type BroadcastPLMNItemExtIEsExtensionValue struct {
	Present int
}

type CancelledCellsInEAIEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CancelledCellsInEAIEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CancelledCellsInEAIEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type CancelledCellsInEAIEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type CancelledCellsInEAINRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CancelledCellsInEAINRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CancelledCellsInEAINRItemExtIEsPresentNothing int = iota /* No components present */
)

type CancelledCellsInEAINRItemExtIEsExtensionValue struct {
	Present int
}

type CancelledCellsInTAIEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CancelledCellsInTAIEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CancelledCellsInTAIEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type CancelledCellsInTAIEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type CancelledCellsInTAINRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CancelledCellsInTAINRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CancelledCellsInTAINRItemExtIEsPresentNothing int = iota /* No components present */
)

type CancelledCellsInTAINRItemExtIEsExtensionValue struct {
	Present int
}

type CellIDBroadcastEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CellIDBroadcastEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CellIDBroadcastEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type CellIDBroadcastEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type CellIDBroadcastNRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CellIDBroadcastNRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CellIDBroadcastNRItemExtIEsPresentNothing int = iota /* No components present */
)

type CellIDBroadcastNRItemExtIEsExtensionValue struct {
	Present int
}

type CellIDCancelledEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CellIDCancelledEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CellIDCancelledEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type CellIDCancelledEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type CellIDCancelledNRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CellIDCancelledNRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CellIDCancelledNRItemExtIEsPresentNothing int = iota /* No components present */
)

type CellIDCancelledNRItemExtIEsExtensionValue struct {
	Present int
}

type CellTypeExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CellTypeExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CellTypeExtIEsPresentNothing int = iota /* No components present */
)

type CellTypeExtIEsExtensionValue struct {
	Present int
}

type CompletedCellsInEAIEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CompletedCellsInEAIEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CompletedCellsInEAIEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type CompletedCellsInEAIEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type CompletedCellsInEAINRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CompletedCellsInEAINRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CompletedCellsInEAINRItemExtIEsPresentNothing int = iota /* No components present */
)

type CompletedCellsInEAINRItemExtIEsExtensionValue struct {
	Present int
}

type CompletedCellsInTAIEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CompletedCellsInTAIEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CompletedCellsInTAIEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type CompletedCellsInTAIEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type CompletedCellsInTAINRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CompletedCellsInTAINRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CompletedCellsInTAINRItemExtIEsPresentNothing int = iota /* No components present */
)

type CompletedCellsInTAINRItemExtIEsExtensionValue struct {
	Present int
}

type CoreNetworkAssistanceInformationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CoreNetworkAssistanceInformationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CoreNetworkAssistanceInformationExtIEsPresentNothing int = iota /* No components present */
)

type CoreNetworkAssistanceInformationExtIEsExtensionValue struct {
	Present int
}

type COUNTValueForPDCPSN12ExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue COUNTValueForPDCPSN12ExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	COUNTValueForPDCPSN12ExtIEsPresentNothing int = iota /* No components present */
)

type COUNTValueForPDCPSN12ExtIEsExtensionValue struct {
	Present int
}

type COUNTValueForPDCPSN18ExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue COUNTValueForPDCPSN18ExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	COUNTValueForPDCPSN18ExtIEsPresentNothing int = iota /* No components present */
)

type COUNTValueForPDCPSN18ExtIEsExtensionValue struct {
	Present int
}

type CriticalityDiagnosticsExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CriticalityDiagnosticsExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CriticalityDiagnosticsExtIEsPresentNothing int = iota /* No components present */
)

type CriticalityDiagnosticsExtIEsExtensionValue struct {
	Present int
}

type CriticalityDiagnosticsIEItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue CriticalityDiagnosticsIEItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CriticalityDiagnosticsIEItemExtIEsPresentNothing int = iota /* No components present */
)

type CriticalityDiagnosticsIEItemExtIEsExtensionValue struct {
	Present int
}

type DataForwardingResponseDRBItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue DataForwardingResponseDRBItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DataForwardingResponseDRBItemExtIEsPresentNothing int = iota /* No components present */
)

type DataForwardingResponseDRBItemExtIEsExtensionValue struct {
	Present int
}

type DRBsSubjectToStatusTransferItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue DRBsSubjectToStatusTransferItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DRBsSubjectToStatusTransferItemExtIEsPresentNothing int = iota /* No components present */
)

type DRBsSubjectToStatusTransferItemExtIEsExtensionValue struct {
	Present int
}

type DRBStatusDL12ExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue DRBStatusDL12ExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DRBStatusDL12ExtIEsPresentNothing int = iota /* No components present */
)

type DRBStatusDL12ExtIEsExtensionValue struct {
	Present int
}

type DRBStatusDL18ExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue DRBStatusDL18ExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DRBStatusDL18ExtIEsPresentNothing int = iota /* No components present */
)

type DRBStatusDL18ExtIEsExtensionValue struct {
	Present int
}

type DRBStatusUL12ExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue DRBStatusUL12ExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DRBStatusUL12ExtIEsPresentNothing int = iota /* No components present */
)

type DRBStatusUL12ExtIEsExtensionValue struct {
	Present int
}

type DRBStatusUL18ExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue DRBStatusUL18ExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DRBStatusUL18ExtIEsPresentNothing int = iota /* No components present */
)

type DRBStatusUL18ExtIEsExtensionValue struct {
	Present int
}

type DRBsToQosFlowsMappingItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue DRBsToQosFlowsMappingItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DRBsToQosFlowsMappingItemExtIEsPresentNothing int = iota /* No components present */
)

type DRBsToQosFlowsMappingItemExtIEsExtensionValue struct {
	Present int
}

type Dynamic5QIDescriptorExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue Dynamic5QIDescriptorExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	Dynamic5QIDescriptorExtIEsPresentNothing int = iota /* No components present */
)

type Dynamic5QIDescriptorExtIEsExtensionValue struct {
	Present int
}

type EmergencyAreaIDBroadcastEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue EmergencyAreaIDBroadcastEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	EmergencyAreaIDBroadcastEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type EmergencyAreaIDBroadcastEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type EmergencyAreaIDBroadcastNRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue EmergencyAreaIDBroadcastNRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	EmergencyAreaIDBroadcastNRItemExtIEsPresentNothing int = iota /* No components present */
)

type EmergencyAreaIDBroadcastNRItemExtIEsExtensionValue struct {
	Present int
}

type EmergencyAreaIDCancelledEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue EmergencyAreaIDCancelledEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	EmergencyAreaIDCancelledEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type EmergencyAreaIDCancelledEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type EmergencyAreaIDCancelledNRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue EmergencyAreaIDCancelledNRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	EmergencyAreaIDCancelledNRItemExtIEsPresentNothing int = iota /* No components present */
)

type EmergencyAreaIDCancelledNRItemExtIEsExtensionValue struct {
	Present int
}

type EmergencyFallbackIndicatorExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue EmergencyFallbackIndicatorExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	EmergencyFallbackIndicatorExtIEsPresentNothing int = iota /* No components present */
)

type EmergencyFallbackIndicatorExtIEsExtensionValue struct {
	Present int
}

type EPSTAIExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue EPSTAIExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	EPSTAIExtIEsPresentNothing int = iota /* No components present */
)

type EPSTAIExtIEsExtensionValue struct {
	Present int
}

type ERABInformationItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue ERABInformationItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ERABInformationItemExtIEsPresentNothing int = iota /* No components present */
)

type ERABInformationItemExtIEsExtensionValue struct {
	Present int
}

type EUTRACGIExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue EUTRACGIExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	EUTRACGIExtIEsPresentNothing int = iota /* No components present */
)

type EUTRACGIExtIEsExtensionValue struct {
	Present int
}

type ExpectedUEActivityBehaviourExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue ExpectedUEActivityBehaviourExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ExpectedUEActivityBehaviourExtIEsPresentNothing int = iota /* No components present */
)

type ExpectedUEActivityBehaviourExtIEsExtensionValue struct {
	Present int
}

type ExpectedUEBehaviourExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue ExpectedUEBehaviourExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ExpectedUEBehaviourExtIEsPresentNothing int = iota /* No components present */
)

type ExpectedUEBehaviourExtIEsExtensionValue struct {
	Present int
}

type ExpectedUEMovingTrajectoryItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue ExpectedUEMovingTrajectoryItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ExpectedUEMovingTrajectoryItemExtIEsPresentNothing int = iota /* No components present */
)

type ExpectedUEMovingTrajectoryItemExtIEsExtensionValue struct {
	Present int
}

type FiveGSTMSIExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue FiveGSTMSIExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	FiveGSTMSIExtIEsPresentNothing int = iota /* No components present */
)

type FiveGSTMSIExtIEsExtensionValue struct {
	Present int
}

type ForbiddenAreaInformationItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue ForbiddenAreaInformationItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ForbiddenAreaInformationItemExtIEsPresentNothing int = iota /* No components present */
)

type ForbiddenAreaInformationItemExtIEsExtensionValue struct {
	Present int
}

type GBRQosInformationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue GBRQosInformationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	GBRQosInformationExtIEsPresentNothing int = iota /* No components present */
)

type GBRQosInformationExtIEsExtensionValue struct {
	Present int
}

type GlobalGNBIDExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue GlobalGNBIDExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	GlobalGNBIDExtIEsPresentNothing int = iota /* No components present */
)

type GlobalGNBIDExtIEsExtensionValue struct {
	Present int
}

type GlobalN3IWFIDExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue GlobalN3IWFIDExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	GlobalN3IWFIDExtIEsPresentNothing int = iota /* No components present */
)

type GlobalN3IWFIDExtIEsExtensionValue struct {
	Present int
}

type GlobalNgENBIDExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue GlobalNgENBIDExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	GlobalNgENBIDExtIEsPresentNothing int = iota /* No components present */
)

type GlobalNgENBIDExtIEsExtensionValue struct {
	Present int
}

type GTPTunnelExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue GTPTunnelExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	GTPTunnelExtIEsPresentNothing int = iota /* No components present */
)

type GTPTunnelExtIEsExtensionValue struct {
	Present int
}

type GUAMIExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue GUAMIExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	GUAMIExtIEsPresentNothing int = iota /* No components present */
)

type GUAMIExtIEsExtensionValue struct {
	Present int
}

type HandoverCommandTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue HandoverCommandTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverCommandTransferExtIEsPresentNothing int = iota /* No components present */
)

type HandoverCommandTransferExtIEsExtensionValue struct {
	Present int
}

type HandoverPreparationUnsuccessfulTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue HandoverPreparationUnsuccessfulTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverPreparationUnsuccessfulTransferExtIEsPresentNothing int = iota /* No components present */
)

type HandoverPreparationUnsuccessfulTransferExtIEsExtensionValue struct {
	Present int
}

type HandoverRequestAcknowledgeTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue HandoverRequestAcknowledgeTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverRequestAcknowledgeTransferExtIEsPresentNothing int = iota /* No components present */
)

type HandoverRequestAcknowledgeTransferExtIEsExtensionValue struct {
	Present int
}

type HandoverRequiredTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue HandoverRequiredTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverRequiredTransferExtIEsPresentNothing int = iota /* No components present */
)

type HandoverRequiredTransferExtIEsExtensionValue struct {
	Present int
}

type HandoverResourceAllocationUnsuccessfulTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue HandoverResourceAllocationUnsuccessfulTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverResourceAllocationUnsuccessfulTransferExtIEsPresentNothing int = iota /* No components present */
)

type HandoverResourceAllocationUnsuccessfulTransferExtIEsExtensionValue struct {
	Present int
}

type InfoOnRecommendedCellsAndRANNodesForPagingExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue InfoOnRecommendedCellsAndRANNodesForPagingExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	InfoOnRecommendedCellsAndRANNodesForPagingExtIEsPresentNothing int = iota /* No components present */
)

type InfoOnRecommendedCellsAndRANNodesForPagingExtIEsExtensionValue struct {
	Present int
}

type LastVisitedCellItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue LastVisitedCellItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	LastVisitedCellItemExtIEsPresentNothing int = iota /* No components present */
)

type LastVisitedCellItemExtIEsExtensionValue struct {
	Present int
}

type LastVisitedNGRANCellInformationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue LastVisitedNGRANCellInformationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	LastVisitedNGRANCellInformationExtIEsPresentNothing int = iota /* No components present */
)

type LastVisitedNGRANCellInformationExtIEsExtensionValue struct {
	Present int
}

type LocationReportingRequestTypeExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue LocationReportingRequestTypeExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	LocationReportingRequestTypeExtIEsPresentNothing int = iota /* No components present */
)

type LocationReportingRequestTypeExtIEsExtensionValue struct {
	Present int
}

type MobilityRestrictionListExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue MobilityRestrictionListExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	MobilityRestrictionListExtIEsPresentNothing int = iota /* No components present */
)

type MobilityRestrictionListExtIEsExtensionValue struct {
	Present int
}

type MultipleTNLInformationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue MultipleTNLInformationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	MultipleTNLInformationExtIEsPresentNothing int = iota /* No components present */
)

type MultipleTNLInformationExtIEsExtensionValue struct {
	Present int
}

type NonDynamic5QIDescriptorExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue NonDynamic5QIDescriptorExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NonDynamic5QIDescriptorExtIEsPresentNothing int = iota /* No components present */
)

type NonDynamic5QIDescriptorExtIEsExtensionValue struct {
	Present int
}

type NRCGIExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue NRCGIExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NRCGIExtIEsPresentNothing int = iota /* No components present */
)

type NRCGIExtIEsExtensionValue struct {
	Present int
}

type OverloadStartNSSAIItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue OverloadStartNSSAIItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	OverloadStartNSSAIItemExtIEsPresentNothing int = iota /* No components present */
)

type OverloadStartNSSAIItemExtIEsExtensionValue struct {
	Present int
}

type PacketErrorRateExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PacketErrorRateExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PacketErrorRateExtIEsPresentNothing int = iota /* No components present */
)

type PacketErrorRateExtIEsExtensionValue struct {
	Present int
}

type PagingAttemptInformationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PagingAttemptInformationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PagingAttemptInformationExtIEsPresentNothing int = iota /* No components present */
)

type PagingAttemptInformationExtIEsExtensionValue struct {
	Present int
}

type PathSwitchRequestAcknowledgeTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PathSwitchRequestAcknowledgeTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PathSwitchRequestAcknowledgeTransferExtIEsPresentNothing int = iota /* No components present */
)

type PathSwitchRequestAcknowledgeTransferExtIEsExtensionValue struct {
	Present int
}

type PathSwitchRequestSetupFailedTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PathSwitchRequestSetupFailedTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PathSwitchRequestSetupFailedTransferExtIEsPresentNothing int = iota /* No components present */
)

type PathSwitchRequestSetupFailedTransferExtIEsExtensionValue struct {
	Present int
}

type PathSwitchRequestTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PathSwitchRequestTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PathSwitchRequestTransferExtIEsPresentNothing int = iota /* No components present */
)

type PathSwitchRequestTransferExtIEsExtensionValue struct {
	Present int
}

type PathSwitchRequestUnsuccessfulTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PathSwitchRequestUnsuccessfulTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PathSwitchRequestUnsuccessfulTransferExtIEsPresentNothing int = iota /* No components present */
)

type PathSwitchRequestUnsuccessfulTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionAggregateMaximumBitRateExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionAggregateMaximumBitRateExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionAggregateMaximumBitRateExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionAggregateMaximumBitRateExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceAdmittedItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceAdmittedItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceAdmittedItemExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceAdmittedItemExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceFailedToModifyItemModCfmExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceFailedToModifyItemModCfmExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceFailedToModifyItemModCfmExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceFailedToModifyItemModCfmExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceFailedToModifyItemModResExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceFailedToModifyItemModResExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceFailedToModifyItemModResExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceFailedToModifyItemModResExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceFailedToSetupItemCxtFailExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceFailedToSetupItemCxtFailExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceFailedToSetupItemCxtFailExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceFailedToSetupItemCxtFailExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceFailedToSetupItemCxtResExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceFailedToSetupItemCxtResExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceFailedToSetupItemCxtResExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceFailedToSetupItemCxtResExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceFailedToSetupItemHOAckExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceFailedToSetupItemHOAckExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceFailedToSetupItemHOAckExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceFailedToSetupItemHOAckExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceFailedToSetupItemPSReqExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceFailedToSetupItemPSReqExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceFailedToSetupItemPSReqExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceFailedToSetupItemPSReqExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceFailedToSetupItemSUResExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceFailedToSetupItemSUResExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceFailedToSetupItemSUResExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceFailedToSetupItemSUResExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceHandoverItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceHandoverItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceHandoverItemExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceHandoverItemExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceInformationItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceInformationItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceInformationItemExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceInformationItemExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceItemCxtRelCplExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceItemCxtRelCplExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceItemCxtRelCplExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceItemCxtRelCplExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceItemCxtRelReqExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceItemCxtRelReqExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceItemCxtRelReqExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceItemCxtRelReqExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceItemHORqdExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceItemHORqdExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceItemHORqdExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceItemHORqdExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyConfirmTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyConfirmTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyConfirmTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyConfirmTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyIndicationUnsuccessfulTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyIndicationUnsuccessfulTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyIndicationUnsuccessfulTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyIndicationUnsuccessfulTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyResponseTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyResponseTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyResponseTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyResponseTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyIndicationTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyIndicationTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyIndicationTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyIndicationTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyItemModCfmExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyItemModCfmExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyItemModCfmExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyItemModCfmExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyItemModIndExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyItemModIndExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyItemModIndExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyItemModIndExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyItemModReqExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyItemModReqExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyItemModReqExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyItemModReqExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyItemModResExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyItemModResExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyItemModResExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyItemModResExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceModifyUnsuccessfulTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceModifyUnsuccessfulTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyUnsuccessfulTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceModifyUnsuccessfulTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceNotifyItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceNotifyItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceNotifyItemExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceNotifyItemExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceNotifyReleasedTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceNotifyReleasedTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceNotifyReleasedTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceNotifyReleasedTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceNotifyTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceNotifyTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceNotifyTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceNotifyTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceReleaseCommandTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceReleaseCommandTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceReleaseCommandTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceReleaseCommandTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceReleasedItemNotExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceReleasedItemNotExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceReleasedItemNotExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceReleasedItemNotExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceReleasedItemPSAckExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceReleasedItemPSAckExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceReleasedItemPSAckExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceReleasedItemPSAckExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceReleasedItemPSFailExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceReleasedItemPSFailExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceReleasedItemPSFailExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceReleasedItemPSFailExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceReleasedItemRelResExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceReleasedItemRelResExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceReleasedItemRelResExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceReleasedItemRelResExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceReleaseResponseTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceReleaseResponseTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceReleaseResponseTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceReleaseResponseTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceSetupItemCxtReqExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceSetupItemCxtReqExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupItemCxtReqExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceSetupItemCxtReqExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceSetupItemCxtResExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceSetupItemCxtResExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupItemCxtResExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceSetupItemCxtResExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceSetupItemHOReqExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceSetupItemHOReqExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupItemHOReqExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceSetupItemHOReqExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceSetupItemSUReqExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceSetupItemSUReqExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupItemSUReqExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceSetupItemSUReqExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceSetupItemSUResExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceSetupItemSUResExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupItemSUResExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceSetupItemSUResExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceSetupResponseTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceSetupResponseTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupResponseTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceSetupResponseTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceSetupUnsuccessfulTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceSetupUnsuccessfulTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupUnsuccessfulTransferExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceSetupUnsuccessfulTransferExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceSwitchedItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceSwitchedItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSwitchedItemExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceSwitchedItemExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceToBeSwitchedDLItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceToBeSwitchedDLItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceToBeSwitchedDLItemExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceToBeSwitchedDLItemExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceToReleaseItemHOCmdExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceToReleaseItemHOCmdExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceToReleaseItemHOCmdExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceToReleaseItemHOCmdExtIEsExtensionValue struct {
	Present int
}

type PDUSessionResourceToReleaseItemRelCmdExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PDUSessionResourceToReleaseItemRelCmdExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceToReleaseItemRelCmdExtIEsPresentNothing int = iota /* No components present */
)

type PDUSessionResourceToReleaseItemRelCmdExtIEsExtensionValue struct {
	Present int
}

type PLMNSupportItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue PLMNSupportItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PLMNSupportItemExtIEsPresentNothing int = iota /* No components present */
)

type PLMNSupportItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowAcceptedItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowAcceptedItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowAcceptedItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowAcceptedItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowAddOrModifyRequestItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowAddOrModifyRequestItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowAddOrModifyRequestItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowAddOrModifyRequestItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowAddOrModifyResponseItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowAddOrModifyResponseItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowAddOrModifyResponseItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowAddOrModifyResponseItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowInformationItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowInformationItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowInformationItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowInformationItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowLevelQosParametersExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowLevelQosParametersExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowLevelQosParametersExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowLevelQosParametersExtIEsExtensionValue struct {
	Present int
}

type QosFlowItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowModifyConfirmItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowModifyConfirmItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowModifyConfirmItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowModifyConfirmItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowNotifyItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowNotifyItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowNotifyItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowNotifyItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowPerTNLInformationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowPerTNLInformationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowPerTNLInformationExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowPerTNLInformationExtIEsExtensionValue struct {
	Present int
}

type QosFlowSetupRequestItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowSetupRequestItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowSetupRequestItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowSetupRequestItemExtIEsExtensionValue struct {
	Present int
}

type QosFlowSetupResponseItemHOReqAckExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowSetupResponseItemHOReqAckExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowSetupResponseItemHOReqAckExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowSetupResponseItemHOReqAckExtIEsExtensionValue struct {
	Present int
}

type QosFlowSetupResponseItemSUResExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowSetupResponseItemSUResExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowSetupResponseItemSUResExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowSetupResponseItemSUResExtIEsExtensionValue struct {
	Present int
}

type QosFlowToBeForwardedItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue QosFlowToBeForwardedItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosFlowToBeForwardedItemExtIEsPresentNothing int = iota /* No components present */
)

type QosFlowToBeForwardedItemExtIEsExtensionValue struct {
	Present int
}

type RANStatusTransferTransparentContainerExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue RANStatusTransferTransparentContainerExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RANStatusTransferTransparentContainerExtIEsPresentNothing int = iota /* No components present */
)

type RANStatusTransferTransparentContainerExtIEsExtensionValue struct {
	Present int
}

type RATRestrictionsItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue RATRestrictionsItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RATRestrictionsItemExtIEsPresentNothing int = iota /* No components present */
)

type RATRestrictionsItemExtIEsExtensionValue struct {
	Present int
}

type RecommendedCellsForPagingExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue RecommendedCellsForPagingExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RecommendedCellsForPagingExtIEsPresentNothing int = iota /* No components present */
)

type RecommendedCellsForPagingExtIEsExtensionValue struct {
	Present int
}

type RecommendedCellItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue RecommendedCellItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RecommendedCellItemExtIEsPresentNothing int = iota /* No components present */
)

type RecommendedCellItemExtIEsExtensionValue struct {
	Present int
}

type RecommendedRANNodesForPagingExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue RecommendedRANNodesForPagingExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RecommendedRANNodesForPagingExtIEsPresentNothing int = iota /* No components present */
)

type RecommendedRANNodesForPagingExtIEsExtensionValue struct {
	Present int
}

type RecommendedRANNodeItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue RecommendedRANNodeItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RecommendedRANNodeItemExtIEsPresentNothing int = iota /* No components present */
)

type RecommendedRANNodeItemExtIEsExtensionValue struct {
	Present int
}

type SecurityContextExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SecurityContextExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SecurityContextExtIEsPresentNothing int = iota /* No components present */
)

type SecurityContextExtIEsExtensionValue struct {
	Present int
}

type SecurityIndicationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SecurityIndicationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SecurityIndicationExtIEsPresentNothing int = iota /* No components present */
)

type SecurityIndicationExtIEsExtensionValue struct {
	Present int
}

type SecurityResultExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SecurityResultExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SecurityResultExtIEsPresentNothing int = iota /* No components present */
)

type SecurityResultExtIEsExtensionValue struct {
	Present int
}

type ServedGUAMIItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue ServedGUAMIItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ServedGUAMIItemExtIEsPresentNothing int = iota /* No components present */
)

type ServedGUAMIItemExtIEsExtensionValue struct {
	Present int
}

type ServiceAreaInformationItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue ServiceAreaInformationItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ServiceAreaInformationItemExtIEsPresentNothing int = iota /* No components present */
)

type ServiceAreaInformationItemExtIEsExtensionValue struct {
	Present int
}

type SingleTNLInformationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SingleTNLInformationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SingleTNLInformationExtIEsPresentNothing int = iota /* No components present */
)

type SingleTNLInformationExtIEsExtensionValue struct {
	Present int
}

type SliceOverloadItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SliceOverloadItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SliceOverloadItemExtIEsPresentNothing int = iota /* No components present */
)

type SliceOverloadItemExtIEsExtensionValue struct {
	Present int
}

type SliceSupportItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SliceSupportItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SliceSupportItemExtIEsPresentNothing int = iota /* No components present */
)

type SliceSupportItemExtIEsExtensionValue struct {
	Present int
}

type SNSSAIExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SNSSAIExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SNSSAIExtIEsPresentNothing int = iota /* No components present */
)

type SNSSAIExtIEsExtensionValue struct {
	Present int
}

type SONConfigurationTransferExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SONConfigurationTransferExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SONConfigurationTransferExtIEsPresentNothing int = iota /* No components present */
)

type SONConfigurationTransferExtIEsExtensionValue struct {
	Present int
}

type SONInformationReplyExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SONInformationReplyExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SONInformationReplyExtIEsPresentNothing int = iota /* No components present */
)

type SONInformationReplyExtIEsExtensionValue struct {
	Present int
}

type SourceNGRANNodeToTargetNGRANNodeTransparentContainerExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SourceNGRANNodeToTargetNGRANNodeTransparentContainerExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SourceNGRANNodeToTargetNGRANNodeTransparentContainerExtIEsPresentNothing int = iota /* No components present */
)

type SourceNGRANNodeToTargetNGRANNodeTransparentContainerExtIEsExtensionValue struct {
	Present int
}

type SourceRANNodeIDExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SourceRANNodeIDExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SourceRANNodeIDExtIEsPresentNothing int = iota /* No components present */
)

type SourceRANNodeIDExtIEsExtensionValue struct {
	Present int
}

type SupportedTAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue SupportedTAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SupportedTAItemExtIEsPresentNothing int = iota /* No components present */
)

type SupportedTAItemExtIEsExtensionValue struct {
	Present int
}

type TAIExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TAIExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TAIExtIEsPresentNothing int = iota /* No components present */
)

type TAIExtIEsExtensionValue struct {
	Present int
}

type TAIBroadcastEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TAIBroadcastEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TAIBroadcastEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type TAIBroadcastEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type TAIBroadcastNRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TAIBroadcastNRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TAIBroadcastNRItemExtIEsPresentNothing int = iota /* No components present */
)

type TAIBroadcastNRItemExtIEsExtensionValue struct {
	Present int
}

type TAICancelledEUTRAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TAICancelledEUTRAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TAICancelledEUTRAItemExtIEsPresentNothing int = iota /* No components present */
)

type TAICancelledEUTRAItemExtIEsExtensionValue struct {
	Present int
}

type TAICancelledNRItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TAICancelledNRItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TAICancelledNRItemExtIEsPresentNothing int = iota /* No components present */
)

type TAICancelledNRItemExtIEsExtensionValue struct {
	Present int
}

type TAIListForInactiveItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TAIListForInactiveItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TAIListForInactiveItemExtIEsPresentNothing int = iota /* No components present */
)

type TAIListForInactiveItemExtIEsExtensionValue struct {
	Present int
}

type TAIListForPagingItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TAIListForPagingItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TAIListForPagingItemExtIEsPresentNothing int = iota /* No components present */
)

type TAIListForPagingItemExtIEsExtensionValue struct {
	Present int
}

type TargeteNBIDExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TargeteNBIDExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TargeteNBIDExtIEsPresentNothing int = iota /* No components present */
)

type TargeteNBIDExtIEsExtensionValue struct {
	Present int
}

type TargetNGRANNodeToSourceNGRANNodeTransparentContainerExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TargetNGRANNodeToSourceNGRANNodeTransparentContainerExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TargetNGRANNodeToSourceNGRANNodeTransparentContainerExtIEsPresentNothing int = iota /* No components present */
)

type TargetNGRANNodeToSourceNGRANNodeTransparentContainerExtIEsExtensionValue struct {
	Present int
}

type TargetRANNodeIDExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TargetRANNodeIDExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TargetRANNodeIDExtIEsPresentNothing int = iota /* No components present */
)

type TargetRANNodeIDExtIEsExtensionValue struct {
	Present int
}

type TNLAssociationItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TNLAssociationItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TNLAssociationItemExtIEsPresentNothing int = iota /* No components present */
)

type TNLAssociationItemExtIEsExtensionValue struct {
	Present int
}

type TNLInformationItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TNLInformationItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TNLInformationItemExtIEsPresentNothing int = iota /* No components present */
)

type TNLInformationItemExtIEsExtensionValue struct {
	Present int
}

type TNLMappingItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TNLMappingItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TNLMappingItemExtIEsPresentNothing int = iota /* No components present */
)

type TNLMappingItemExtIEsExtensionValue struct {
	Present int
}

type TraceActivationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue TraceActivationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TraceActivationExtIEsPresentNothing int = iota /* No components present */
)

type TraceActivationExtIEsExtensionValue struct {
	Present int
}

type UEAggregateMaximumBitRateExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UEAggregateMaximumBitRateExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEAggregateMaximumBitRateExtIEsPresentNothing int = iota /* No components present */
)

type UEAggregateMaximumBitRateExtIEsExtensionValue struct {
	Present int
}

type UEAssociatedLogicalNGConnectionItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UEAssociatedLogicalNGConnectionItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEAssociatedLogicalNGConnectionItemExtIEsPresentNothing int = iota /* No components present */
)

type UEAssociatedLogicalNGConnectionItemExtIEsExtensionValue struct {
	Present int
}

type UENGAPIDPairExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UENGAPIDPairExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UENGAPIDPairExtIEsPresentNothing int = iota /* No components present */
)

type UENGAPIDPairExtIEsExtensionValue struct {
	Present int
}

type UEPresenceInAreaOfInterestItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UEPresenceInAreaOfInterestItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEPresenceInAreaOfInterestItemExtIEsPresentNothing int = iota /* No components present */
)

type UEPresenceInAreaOfInterestItemExtIEsExtensionValue struct {
	Present int
}

type UERadioCapabilityForPagingExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UERadioCapabilityForPagingExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UERadioCapabilityForPagingExtIEsPresentNothing int = iota /* No components present */
)

type UERadioCapabilityForPagingExtIEsExtensionValue struct {
	Present int
}

type UESecurityCapabilitiesExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UESecurityCapabilitiesExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UESecurityCapabilitiesExtIEsPresentNothing int = iota /* No components present */
)

type UESecurityCapabilitiesExtIEsExtensionValue struct {
	Present int
}

type ULNGUUPTNLModifyItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue ULNGUUPTNLModifyItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ULNGUUPTNLModifyItemExtIEsPresentNothing int = iota /* No components present */
)

type ULNGUUPTNLModifyItemExtIEsExtensionValue struct {
	Present int
}

type UnavailableGUAMIItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UnavailableGUAMIItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UnavailableGUAMIItemExtIEsPresentNothing int = iota /* No components present */
)

type UnavailableGUAMIItemExtIEsExtensionValue struct {
	Present int
}

type UserLocationInformationEUTRAExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UserLocationInformationEUTRAExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UserLocationInformationEUTRAExtIEsPresentNothing int = iota /* No components present */
)

type UserLocationInformationEUTRAExtIEsExtensionValue struct {
	Present int
}

type UserLocationInformationN3IWFExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UserLocationInformationN3IWFExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UserLocationInformationN3IWFExtIEsPresentNothing int = iota /* No components present */
)

type UserLocationInformationN3IWFExtIEsExtensionValue struct {
	Present int
}

type UserLocationInformationNRExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UserLocationInformationNRExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UserLocationInformationNRExtIEsPresentNothing int = iota /* No components present */
)

type UserLocationInformationNRExtIEsExtensionValue struct {
	Present int
}

type UserPlaneSecurityInformationExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue UserPlaneSecurityInformationExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UserPlaneSecurityInformationExtIEsPresentNothing int = iota /* No components present */
)

type UserPlaneSecurityInformationExtIEsExtensionValue struct {
	Present int
}

type XnExtTLAItemExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue XnExtTLAItemExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	XnExtTLAItemExtIEsPresentNothing int = iota /* No components present */
)

type XnExtTLAItemExtIEsExtensionValue struct {
	Present int
}

type XnTNLConfigurationInfoExtIEs struct {
	Id             ProtocolExtensionID
	Criticality    Criticality
	ExtensionValue XnTNLConfigurationInfoExtIEsExtensionValue `aper:"openType,referenceFieldName:Id"`
}

const (
	XnTNLConfigurationInfoExtIEsPresentNothing int = iota /* No components present */
)

type XnTNLConfigurationInfoExtIEsExtensionValue struct {
	Present int
}
