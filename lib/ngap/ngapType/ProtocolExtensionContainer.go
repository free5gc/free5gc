//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P0 */
/* AllocationAndRetentionPriorityExtIEs */
type ProtocolExtensionContainerAllocationAndRetentionPriorityExtIEs struct {
	List []AllocationAndRetentionPriorityExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P1 */
/* AllowedNSSAIItemExtIEs */
type ProtocolExtensionContainerAllowedNSSAIItemExtIEs struct {
	List []AllowedNSSAIItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P2 */
/* AMFTNLAssociationSetupItemExtIEs */
type ProtocolExtensionContainerAMFTNLAssociationSetupItemExtIEs struct {
	List []AMFTNLAssociationSetupItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P3 */
/* AMFTNLAssociationToAddItemExtIEs */
type ProtocolExtensionContainerAMFTNLAssociationToAddItemExtIEs struct {
	List []AMFTNLAssociationToAddItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P4 */
/* AMFTNLAssociationToRemoveItemExtIEs */
type ProtocolExtensionContainerAMFTNLAssociationToRemoveItemExtIEs struct {
	List []AMFTNLAssociationToRemoveItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P5 */
/* AMFTNLAssociationToUpdateItemExtIEs */
type ProtocolExtensionContainerAMFTNLAssociationToUpdateItemExtIEs struct {
	List []AMFTNLAssociationToUpdateItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P6 */
/* AreaOfInterestExtIEs */
type ProtocolExtensionContainerAreaOfInterestExtIEs struct {
	List []AreaOfInterestExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P7 */
/* AreaOfInterestCellItemExtIEs */
type ProtocolExtensionContainerAreaOfInterestCellItemExtIEs struct {
	List []AreaOfInterestCellItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P8 */
/* AreaOfInterestItemExtIEs */
type ProtocolExtensionContainerAreaOfInterestItemExtIEs struct {
	List []AreaOfInterestItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P9 */
/* AreaOfInterestRANNodeItemExtIEs */
type ProtocolExtensionContainerAreaOfInterestRANNodeItemExtIEs struct {
	List []AreaOfInterestRANNodeItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P10 */
/* AreaOfInterestTAIItemExtIEs */
type ProtocolExtensionContainerAreaOfInterestTAIItemExtIEs struct {
	List []AreaOfInterestTAIItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P11 */
/* AssistanceDataForPagingExtIEs */
type ProtocolExtensionContainerAssistanceDataForPagingExtIEs struct {
	List []AssistanceDataForPagingExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P12 */
/* AssistanceDataForRecommendedCellsExtIEs */
type ProtocolExtensionContainerAssistanceDataForRecommendedCellsExtIEs struct {
	List []AssistanceDataForRecommendedCellsExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P13 */
/* AssociatedQosFlowItemExtIEs */
type ProtocolExtensionContainerAssociatedQosFlowItemExtIEs struct {
	List []AssociatedQosFlowItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P14 */
/* BroadcastPLMNItemExtIEs */
type ProtocolExtensionContainerBroadcastPLMNItemExtIEs struct {
	List []BroadcastPLMNItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P15 */
/* CancelledCellsInEAIEUTRAItemExtIEs */
type ProtocolExtensionContainerCancelledCellsInEAIEUTRAItemExtIEs struct {
	List []CancelledCellsInEAIEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P16 */
/* CancelledCellsInEAINRItemExtIEs */
type ProtocolExtensionContainerCancelledCellsInEAINRItemExtIEs struct {
	List []CancelledCellsInEAINRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P17 */
/* CancelledCellsInTAIEUTRAItemExtIEs */
type ProtocolExtensionContainerCancelledCellsInTAIEUTRAItemExtIEs struct {
	List []CancelledCellsInTAIEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P18 */
/* CancelledCellsInTAINRItemExtIEs */
type ProtocolExtensionContainerCancelledCellsInTAINRItemExtIEs struct {
	List []CancelledCellsInTAINRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P19 */
/* CellIDBroadcastEUTRAItemExtIEs */
type ProtocolExtensionContainerCellIDBroadcastEUTRAItemExtIEs struct {
	List []CellIDBroadcastEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P20 */
/* CellIDBroadcastNRItemExtIEs */
type ProtocolExtensionContainerCellIDBroadcastNRItemExtIEs struct {
	List []CellIDBroadcastNRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P21 */
/* CellIDCancelledEUTRAItemExtIEs */
type ProtocolExtensionContainerCellIDCancelledEUTRAItemExtIEs struct {
	List []CellIDCancelledEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P22 */
/* CellIDCancelledNRItemExtIEs */
type ProtocolExtensionContainerCellIDCancelledNRItemExtIEs struct {
	List []CellIDCancelledNRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P23 */
/* CellTypeExtIEs */
type ProtocolExtensionContainerCellTypeExtIEs struct {
	List []CellTypeExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P24 */
/* CompletedCellsInEAIEUTRAItemExtIEs */
type ProtocolExtensionContainerCompletedCellsInEAIEUTRAItemExtIEs struct {
	List []CompletedCellsInEAIEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P25 */
/* CompletedCellsInEAINRItemExtIEs */
type ProtocolExtensionContainerCompletedCellsInEAINRItemExtIEs struct {
	List []CompletedCellsInEAINRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P26 */
/* CompletedCellsInTAIEUTRAItemExtIEs */
type ProtocolExtensionContainerCompletedCellsInTAIEUTRAItemExtIEs struct {
	List []CompletedCellsInTAIEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P27 */
/* CompletedCellsInTAINRItemExtIEs */
type ProtocolExtensionContainerCompletedCellsInTAINRItemExtIEs struct {
	List []CompletedCellsInTAINRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P28 */
/* CoreNetworkAssistanceInformationExtIEs */
type ProtocolExtensionContainerCoreNetworkAssistanceInformationExtIEs struct {
	List []CoreNetworkAssistanceInformationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P29 */
/* COUNTValueForPDCPSN12ExtIEs */
type ProtocolExtensionContainerCOUNTValueForPDCPSN12ExtIEs struct {
	List []COUNTValueForPDCPSN12ExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P30 */
/* COUNTValueForPDCPSN18ExtIEs */
type ProtocolExtensionContainerCOUNTValueForPDCPSN18ExtIEs struct {
	List []COUNTValueForPDCPSN18ExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P31 */
/* CriticalityDiagnosticsExtIEs */
type ProtocolExtensionContainerCriticalityDiagnosticsExtIEs struct {
	List []CriticalityDiagnosticsExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P32 */
/* CriticalityDiagnosticsIEItemExtIEs */
type ProtocolExtensionContainerCriticalityDiagnosticsIEItemExtIEs struct {
	List []CriticalityDiagnosticsIEItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P33 */
/* DataForwardingResponseDRBItemExtIEs */
type ProtocolExtensionContainerDataForwardingResponseDRBItemExtIEs struct {
	List []DataForwardingResponseDRBItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P34 */
/* DRBsSubjectToStatusTransferItemExtIEs */
type ProtocolExtensionContainerDRBsSubjectToStatusTransferItemExtIEs struct {
	List []DRBsSubjectToStatusTransferItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P35 */
/* DRBStatusDL12ExtIEs */
type ProtocolExtensionContainerDRBStatusDL12ExtIEs struct {
	List []DRBStatusDL12ExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P36 */
/* DRBStatusDL18ExtIEs */
type ProtocolExtensionContainerDRBStatusDL18ExtIEs struct {
	List []DRBStatusDL18ExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P37 */
/* DRBStatusUL12ExtIEs */
type ProtocolExtensionContainerDRBStatusUL12ExtIEs struct {
	List []DRBStatusUL12ExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P38 */
/* DRBStatusUL18ExtIEs */
type ProtocolExtensionContainerDRBStatusUL18ExtIEs struct {
	List []DRBStatusUL18ExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P39 */
/* DRBsToQosFlowsMappingItemExtIEs */
type ProtocolExtensionContainerDRBsToQosFlowsMappingItemExtIEs struct {
	List []DRBsToQosFlowsMappingItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P40 */
/* Dynamic5QIDescriptorExtIEs */
type ProtocolExtensionContainerDynamic5QIDescriptorExtIEs struct {
	List []Dynamic5QIDescriptorExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P41 */
/* EmergencyAreaIDBroadcastEUTRAItemExtIEs */
type ProtocolExtensionContainerEmergencyAreaIDBroadcastEUTRAItemExtIEs struct {
	List []EmergencyAreaIDBroadcastEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P42 */
/* EmergencyAreaIDBroadcastNRItemExtIEs */
type ProtocolExtensionContainerEmergencyAreaIDBroadcastNRItemExtIEs struct {
	List []EmergencyAreaIDBroadcastNRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P43 */
/* EmergencyAreaIDCancelledEUTRAItemExtIEs */
type ProtocolExtensionContainerEmergencyAreaIDCancelledEUTRAItemExtIEs struct {
	List []EmergencyAreaIDCancelledEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P44 */
/* EmergencyAreaIDCancelledNRItemExtIEs */
type ProtocolExtensionContainerEmergencyAreaIDCancelledNRItemExtIEs struct {
	List []EmergencyAreaIDCancelledNRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P45 */
/* EmergencyFallbackIndicatorExtIEs */
type ProtocolExtensionContainerEmergencyFallbackIndicatorExtIEs struct {
	List []EmergencyFallbackIndicatorExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P46 */
/* EPSTAIExtIEs */
type ProtocolExtensionContainerEPSTAIExtIEs struct {
	List []EPSTAIExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P47 */
/* ERABInformationItemExtIEs */
type ProtocolExtensionContainerERABInformationItemExtIEs struct {
	List []ERABInformationItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P48 */
/* EUTRACGIExtIEs */
type ProtocolExtensionContainerEUTRACGIExtIEs struct {
	List []EUTRACGIExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P49 */
/* ExpectedUEActivityBehaviourExtIEs */
type ProtocolExtensionContainerExpectedUEActivityBehaviourExtIEs struct {
	List []ExpectedUEActivityBehaviourExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P50 */
/* ExpectedUEBehaviourExtIEs */
type ProtocolExtensionContainerExpectedUEBehaviourExtIEs struct {
	List []ExpectedUEBehaviourExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P51 */
/* ExpectedUEMovingTrajectoryItemExtIEs */
type ProtocolExtensionContainerExpectedUEMovingTrajectoryItemExtIEs struct {
	List []ExpectedUEMovingTrajectoryItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P52 */
/* FiveGSTMSIExtIEs */
type ProtocolExtensionContainerFiveGSTMSIExtIEs struct {
	List []FiveGSTMSIExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P53 */
/* ForbiddenAreaInformationItemExtIEs */
type ProtocolExtensionContainerForbiddenAreaInformationItemExtIEs struct {
	List []ForbiddenAreaInformationItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P54 */
/* GBRQosInformationExtIEs */
type ProtocolExtensionContainerGBRQosInformationExtIEs struct {
	List []GBRQosInformationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P55 */
/* GlobalGNBIDExtIEs */
type ProtocolExtensionContainerGlobalGNBIDExtIEs struct {
	List []GlobalGNBIDExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P56 */
/* GlobalN3IWFIDExtIEs */
type ProtocolExtensionContainerGlobalN3IWFIDExtIEs struct {
	List []GlobalN3IWFIDExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P57 */
/* GlobalNgENBIDExtIEs */
type ProtocolExtensionContainerGlobalNgENBIDExtIEs struct {
	List []GlobalNgENBIDExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P58 */
/* GTPTunnelExtIEs */
type ProtocolExtensionContainerGTPTunnelExtIEs struct {
	List []GTPTunnelExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P59 */
/* GUAMIExtIEs */
type ProtocolExtensionContainerGUAMIExtIEs struct {
	List []GUAMIExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P60 */
/* HandoverCommandTransferExtIEs */
type ProtocolExtensionContainerHandoverCommandTransferExtIEs struct {
	List []HandoverCommandTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P61 */
/* HandoverPreparationUnsuccessfulTransferExtIEs */
type ProtocolExtensionContainerHandoverPreparationUnsuccessfulTransferExtIEs struct {
	List []HandoverPreparationUnsuccessfulTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P62 */
/* HandoverRequestAcknowledgeTransferExtIEs */
type ProtocolExtensionContainerHandoverRequestAcknowledgeTransferExtIEs struct {
	List []HandoverRequestAcknowledgeTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P63 */
/* HandoverRequiredTransferExtIEs */
type ProtocolExtensionContainerHandoverRequiredTransferExtIEs struct {
	List []HandoverRequiredTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P64 */
/* HandoverResourceAllocationUnsuccessfulTransferExtIEs */
type ProtocolExtensionContainerHandoverResourceAllocationUnsuccessfulTransferExtIEs struct {
	List []HandoverResourceAllocationUnsuccessfulTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P65 */
/* InfoOnRecommendedCellsAndRANNodesForPagingExtIEs */
type ProtocolExtensionContainerInfoOnRecommendedCellsAndRANNodesForPagingExtIEs struct {
	List []InfoOnRecommendedCellsAndRANNodesForPagingExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P66 */
/* LastVisitedCellItemExtIEs */
type ProtocolExtensionContainerLastVisitedCellItemExtIEs struct {
	List []LastVisitedCellItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P67 */
/* LastVisitedNGRANCellInformationExtIEs */
type ProtocolExtensionContainerLastVisitedNGRANCellInformationExtIEs struct {
	List []LastVisitedNGRANCellInformationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P68 */
/* LocationReportingRequestTypeExtIEs */
type ProtocolExtensionContainerLocationReportingRequestTypeExtIEs struct {
	List []LocationReportingRequestTypeExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P69 */
/* MobilityRestrictionListExtIEs */
type ProtocolExtensionContainerMobilityRestrictionListExtIEs struct {
	List []MobilityRestrictionListExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P70 */
/* MultipleTNLInformationExtIEs */
type ProtocolExtensionContainerMultipleTNLInformationExtIEs struct {
	List []MultipleTNLInformationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P71 */
/* NonDynamic5QIDescriptorExtIEs */
type ProtocolExtensionContainerNonDynamic5QIDescriptorExtIEs struct {
	List []NonDynamic5QIDescriptorExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P72 */
/* NRCGIExtIEs */
type ProtocolExtensionContainerNRCGIExtIEs struct {
	List []NRCGIExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P73 */
/* OverloadStartNSSAIItemExtIEs */
type ProtocolExtensionContainerOverloadStartNSSAIItemExtIEs struct {
	List []OverloadStartNSSAIItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P74 */
/* PacketErrorRateExtIEs */
type ProtocolExtensionContainerPacketErrorRateExtIEs struct {
	List []PacketErrorRateExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P75 */
/* PagingAttemptInformationExtIEs */
type ProtocolExtensionContainerPagingAttemptInformationExtIEs struct {
	List []PagingAttemptInformationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P76 */
/* PathSwitchRequestAcknowledgeTransferExtIEs */
type ProtocolExtensionContainerPathSwitchRequestAcknowledgeTransferExtIEs struct {
	List []PathSwitchRequestAcknowledgeTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P77 */
/* PathSwitchRequestSetupFailedTransferExtIEs */
type ProtocolExtensionContainerPathSwitchRequestSetupFailedTransferExtIEs struct {
	List []PathSwitchRequestSetupFailedTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P78 */
/* PathSwitchRequestTransferExtIEs */
type ProtocolExtensionContainerPathSwitchRequestTransferExtIEs struct {
	List []PathSwitchRequestTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P79 */
/* PathSwitchRequestUnsuccessfulTransferExtIEs */
type ProtocolExtensionContainerPathSwitchRequestUnsuccessfulTransferExtIEs struct {
	List []PathSwitchRequestUnsuccessfulTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P80 */
/* PDUSessionAggregateMaximumBitRateExtIEs */
type ProtocolExtensionContainerPDUSessionAggregateMaximumBitRateExtIEs struct {
	List []PDUSessionAggregateMaximumBitRateExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P81 */
/* PDUSessionResourceAdmittedItemExtIEs */
type ProtocolExtensionContainerPDUSessionResourceAdmittedItemExtIEs struct {
	List []PDUSessionResourceAdmittedItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P82 */
/* PDUSessionResourceFailedToModifyItemModCfmExtIEs */
type ProtocolExtensionContainerPDUSessionResourceFailedToModifyItemModCfmExtIEs struct {
	List []PDUSessionResourceFailedToModifyItemModCfmExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P83 */
/* PDUSessionResourceFailedToModifyItemModResExtIEs */
type ProtocolExtensionContainerPDUSessionResourceFailedToModifyItemModResExtIEs struct {
	List []PDUSessionResourceFailedToModifyItemModResExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P84 */
/* PDUSessionResourceFailedToSetupItemCxtFailExtIEs */
type ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemCxtFailExtIEs struct {
	List []PDUSessionResourceFailedToSetupItemCxtFailExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P85 */
/* PDUSessionResourceFailedToSetupItemCxtResExtIEs */
type ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemCxtResExtIEs struct {
	List []PDUSessionResourceFailedToSetupItemCxtResExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P86 */
/* PDUSessionResourceFailedToSetupItemHOAckExtIEs */
type ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemHOAckExtIEs struct {
	List []PDUSessionResourceFailedToSetupItemHOAckExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P87 */
/* PDUSessionResourceFailedToSetupItemPSReqExtIEs */
type ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemPSReqExtIEs struct {
	List []PDUSessionResourceFailedToSetupItemPSReqExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P88 */
/* PDUSessionResourceFailedToSetupItemSUResExtIEs */
type ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemSUResExtIEs struct {
	List []PDUSessionResourceFailedToSetupItemSUResExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P89 */
/* PDUSessionResourceHandoverItemExtIEs */
type ProtocolExtensionContainerPDUSessionResourceHandoverItemExtIEs struct {
	List []PDUSessionResourceHandoverItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P90 */
/* PDUSessionResourceInformationItemExtIEs */
type ProtocolExtensionContainerPDUSessionResourceInformationItemExtIEs struct {
	List []PDUSessionResourceInformationItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P91 */
/* PDUSessionResourceItemCxtRelCplExtIEs */
type ProtocolExtensionContainerPDUSessionResourceItemCxtRelCplExtIEs struct {
	List []PDUSessionResourceItemCxtRelCplExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P92 */
/* PDUSessionResourceItemCxtRelReqExtIEs */
type ProtocolExtensionContainerPDUSessionResourceItemCxtRelReqExtIEs struct {
	List []PDUSessionResourceItemCxtRelReqExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P93 */
/* PDUSessionResourceItemHORqdExtIEs */
type ProtocolExtensionContainerPDUSessionResourceItemHORqdExtIEs struct {
	List []PDUSessionResourceItemHORqdExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P94 */
/* PDUSessionResourceModifyConfirmTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyConfirmTransferExtIEs struct {
	List []PDUSessionResourceModifyConfirmTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P95 */
/* PDUSessionResourceModifyIndicationUnsuccessfulTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyIndicationUnsuccessfulTransferExtIEs struct {
	List []PDUSessionResourceModifyIndicationUnsuccessfulTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P96 */
/* PDUSessionResourceModifyResponseTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyResponseTransferExtIEs struct {
	List []PDUSessionResourceModifyResponseTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P97 */
/* PDUSessionResourceModifyIndicationTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyIndicationTransferExtIEs struct {
	List []PDUSessionResourceModifyIndicationTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P98 */
/* PDUSessionResourceModifyItemModCfmExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyItemModCfmExtIEs struct {
	List []PDUSessionResourceModifyItemModCfmExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P99 */
/* PDUSessionResourceModifyItemModIndExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyItemModIndExtIEs struct {
	List []PDUSessionResourceModifyItemModIndExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P100 */
/* PDUSessionResourceModifyItemModReqExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyItemModReqExtIEs struct {
	List []PDUSessionResourceModifyItemModReqExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P101 */
/* PDUSessionResourceModifyItemModResExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyItemModResExtIEs struct {
	List []PDUSessionResourceModifyItemModResExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P102 */
/* PDUSessionResourceModifyUnsuccessfulTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceModifyUnsuccessfulTransferExtIEs struct {
	List []PDUSessionResourceModifyUnsuccessfulTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P103 */
/* PDUSessionResourceNotifyItemExtIEs */
type ProtocolExtensionContainerPDUSessionResourceNotifyItemExtIEs struct {
	List []PDUSessionResourceNotifyItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P104 */
/* PDUSessionResourceNotifyReleasedTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceNotifyReleasedTransferExtIEs struct {
	List []PDUSessionResourceNotifyReleasedTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P105 */
/* PDUSessionResourceNotifyTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceNotifyTransferExtIEs struct {
	List []PDUSessionResourceNotifyTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P106 */
/* PDUSessionResourceReleaseCommandTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceReleaseCommandTransferExtIEs struct {
	List []PDUSessionResourceReleaseCommandTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P107 */
/* PDUSessionResourceReleasedItemNotExtIEs */
type ProtocolExtensionContainerPDUSessionResourceReleasedItemNotExtIEs struct {
	List []PDUSessionResourceReleasedItemNotExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P108 */
/* PDUSessionResourceReleasedItemPSAckExtIEs */
type ProtocolExtensionContainerPDUSessionResourceReleasedItemPSAckExtIEs struct {
	List []PDUSessionResourceReleasedItemPSAckExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P109 */
/* PDUSessionResourceReleasedItemPSFailExtIEs */
type ProtocolExtensionContainerPDUSessionResourceReleasedItemPSFailExtIEs struct {
	List []PDUSessionResourceReleasedItemPSFailExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P110 */
/* PDUSessionResourceReleasedItemRelResExtIEs */
type ProtocolExtensionContainerPDUSessionResourceReleasedItemRelResExtIEs struct {
	List []PDUSessionResourceReleasedItemRelResExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P111 */
/* PDUSessionResourceReleaseResponseTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceReleaseResponseTransferExtIEs struct {
	List []PDUSessionResourceReleaseResponseTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P112 */
/* PDUSessionResourceSetupItemCxtReqExtIEs */
type ProtocolExtensionContainerPDUSessionResourceSetupItemCxtReqExtIEs struct {
	List []PDUSessionResourceSetupItemCxtReqExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P113 */
/* PDUSessionResourceSetupItemCxtResExtIEs */
type ProtocolExtensionContainerPDUSessionResourceSetupItemCxtResExtIEs struct {
	List []PDUSessionResourceSetupItemCxtResExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P114 */
/* PDUSessionResourceSetupItemHOReqExtIEs */
type ProtocolExtensionContainerPDUSessionResourceSetupItemHOReqExtIEs struct {
	List []PDUSessionResourceSetupItemHOReqExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P115 */
/* PDUSessionResourceSetupItemSUReqExtIEs */
type ProtocolExtensionContainerPDUSessionResourceSetupItemSUReqExtIEs struct {
	List []PDUSessionResourceSetupItemSUReqExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P116 */
/* PDUSessionResourceSetupItemSUResExtIEs */
type ProtocolExtensionContainerPDUSessionResourceSetupItemSUResExtIEs struct {
	List []PDUSessionResourceSetupItemSUResExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P117 */
/* PDUSessionResourceSetupResponseTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceSetupResponseTransferExtIEs struct {
	List []PDUSessionResourceSetupResponseTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P118 */
/* PDUSessionResourceSetupUnsuccessfulTransferExtIEs */
type ProtocolExtensionContainerPDUSessionResourceSetupUnsuccessfulTransferExtIEs struct {
	List []PDUSessionResourceSetupUnsuccessfulTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P119 */
/* PDUSessionResourceSwitchedItemExtIEs */
type ProtocolExtensionContainerPDUSessionResourceSwitchedItemExtIEs struct {
	List []PDUSessionResourceSwitchedItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P120 */
/* PDUSessionResourceToBeSwitchedDLItemExtIEs */
type ProtocolExtensionContainerPDUSessionResourceToBeSwitchedDLItemExtIEs struct {
	List []PDUSessionResourceToBeSwitchedDLItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P121 */
/* PDUSessionResourceToReleaseItemHOCmdExtIEs */
type ProtocolExtensionContainerPDUSessionResourceToReleaseItemHOCmdExtIEs struct {
	List []PDUSessionResourceToReleaseItemHOCmdExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P122 */
/* PDUSessionResourceToReleaseItemRelCmdExtIEs */
type ProtocolExtensionContainerPDUSessionResourceToReleaseItemRelCmdExtIEs struct {
	List []PDUSessionResourceToReleaseItemRelCmdExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P123 */
/* PLMNSupportItemExtIEs */
type ProtocolExtensionContainerPLMNSupportItemExtIEs struct {
	List []PLMNSupportItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P124 */
/* QosFlowAcceptedItemExtIEs */
type ProtocolExtensionContainerQosFlowAcceptedItemExtIEs struct {
	List []QosFlowAcceptedItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P125 */
/* QosFlowAddOrModifyRequestItemExtIEs */
type ProtocolExtensionContainerQosFlowAddOrModifyRequestItemExtIEs struct {
	List []QosFlowAddOrModifyRequestItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P126 */
/* QosFlowAddOrModifyResponseItemExtIEs */
type ProtocolExtensionContainerQosFlowAddOrModifyResponseItemExtIEs struct {
	List []QosFlowAddOrModifyResponseItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P127 */
/* QosFlowInformationItemExtIEs */
type ProtocolExtensionContainerQosFlowInformationItemExtIEs struct {
	List []QosFlowInformationItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P128 */
/* QosFlowLevelQosParametersExtIEs */
type ProtocolExtensionContainerQosFlowLevelQosParametersExtIEs struct {
	List []QosFlowLevelQosParametersExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P129 */
/* QosFlowItemExtIEs */
type ProtocolExtensionContainerQosFlowItemExtIEs struct {
	List []QosFlowItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P130 */
/* QosFlowModifyConfirmItemExtIEs */
type ProtocolExtensionContainerQosFlowModifyConfirmItemExtIEs struct {
	List []QosFlowModifyConfirmItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P131 */
/* QosFlowNotifyItemExtIEs */
type ProtocolExtensionContainerQosFlowNotifyItemExtIEs struct {
	List []QosFlowNotifyItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P132 */
/* QosFlowPerTNLInformationExtIEs */
type ProtocolExtensionContainerQosFlowPerTNLInformationExtIEs struct {
	List []QosFlowPerTNLInformationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P133 */
/* QosFlowSetupRequestItemExtIEs */
type ProtocolExtensionContainerQosFlowSetupRequestItemExtIEs struct {
	List []QosFlowSetupRequestItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P134 */
/* QosFlowSetupResponseItemHOReqAckExtIEs */
type ProtocolExtensionContainerQosFlowSetupResponseItemHOReqAckExtIEs struct {
	List []QosFlowSetupResponseItemHOReqAckExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P135 */
/* QosFlowSetupResponseItemSUResExtIEs */
type ProtocolExtensionContainerQosFlowSetupResponseItemSUResExtIEs struct {
	List []QosFlowSetupResponseItemSUResExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P136 */
/* QosFlowToBeForwardedItemExtIEs */
type ProtocolExtensionContainerQosFlowToBeForwardedItemExtIEs struct {
	List []QosFlowToBeForwardedItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P137 */
/* RANStatusTransferTransparentContainerExtIEs */
type ProtocolExtensionContainerRANStatusTransferTransparentContainerExtIEs struct {
	List []RANStatusTransferTransparentContainerExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P138 */
/* RATRestrictionsItemExtIEs */
type ProtocolExtensionContainerRATRestrictionsItemExtIEs struct {
	List []RATRestrictionsItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P139 */
/* RecommendedCellsForPagingExtIEs */
type ProtocolExtensionContainerRecommendedCellsForPagingExtIEs struct {
	List []RecommendedCellsForPagingExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P140 */
/* RecommendedCellItemExtIEs */
type ProtocolExtensionContainerRecommendedCellItemExtIEs struct {
	List []RecommendedCellItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P141 */
/* RecommendedRANNodesForPagingExtIEs */
type ProtocolExtensionContainerRecommendedRANNodesForPagingExtIEs struct {
	List []RecommendedRANNodesForPagingExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P142 */
/* RecommendedRANNodeItemExtIEs */
type ProtocolExtensionContainerRecommendedRANNodeItemExtIEs struct {
	List []RecommendedRANNodeItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P143 */
/* SecurityContextExtIEs */
type ProtocolExtensionContainerSecurityContextExtIEs struct {
	List []SecurityContextExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P144 */
/* SecurityIndicationExtIEs */
type ProtocolExtensionContainerSecurityIndicationExtIEs struct {
	List []SecurityIndicationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P145 */
/* SecurityResultExtIEs */
type ProtocolExtensionContainerSecurityResultExtIEs struct {
	List []SecurityResultExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P146 */
/* ServedGUAMIItemExtIEs */
type ProtocolExtensionContainerServedGUAMIItemExtIEs struct {
	List []ServedGUAMIItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P147 */
/* ServiceAreaInformationItemExtIEs */
type ProtocolExtensionContainerServiceAreaInformationItemExtIEs struct {
	List []ServiceAreaInformationItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P148 */
/* SingleTNLInformationExtIEs */
type ProtocolExtensionContainerSingleTNLInformationExtIEs struct {
	List []SingleTNLInformationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P149 */
/* SliceOverloadItemExtIEs */
type ProtocolExtensionContainerSliceOverloadItemExtIEs struct {
	List []SliceOverloadItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P150 */
/* SliceSupportItemExtIEs */
type ProtocolExtensionContainerSliceSupportItemExtIEs struct {
	List []SliceSupportItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P151 */
/* SNSSAIExtIEs */
type ProtocolExtensionContainerSNSSAIExtIEs struct {
	List []SNSSAIExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P152 */
/* SONConfigurationTransferExtIEs */
type ProtocolExtensionContainerSONConfigurationTransferExtIEs struct {
	List []SONConfigurationTransferExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P153 */
/* SONInformationReplyExtIEs */
type ProtocolExtensionContainerSONInformationReplyExtIEs struct {
	List []SONInformationReplyExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P154 */
/* SourceNGRANNodeToTargetNGRANNodeTransparentContainerExtIEs */
type ProtocolExtensionContainerSourceNGRANNodeToTargetNGRANNodeTransparentContainerExtIEs struct {
	List []SourceNGRANNodeToTargetNGRANNodeTransparentContainerExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P155 */
/* SourceRANNodeIDExtIEs */
type ProtocolExtensionContainerSourceRANNodeIDExtIEs struct {
	List []SourceRANNodeIDExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P156 */
/* SupportedTAItemExtIEs */
type ProtocolExtensionContainerSupportedTAItemExtIEs struct {
	List []SupportedTAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P157 */
/* TAIExtIEs */
type ProtocolExtensionContainerTAIExtIEs struct {
	List []TAIExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P158 */
/* TAIBroadcastEUTRAItemExtIEs */
type ProtocolExtensionContainerTAIBroadcastEUTRAItemExtIEs struct {
	List []TAIBroadcastEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P159 */
/* TAIBroadcastNRItemExtIEs */
type ProtocolExtensionContainerTAIBroadcastNRItemExtIEs struct {
	List []TAIBroadcastNRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P160 */
/* TAICancelledEUTRAItemExtIEs */
type ProtocolExtensionContainerTAICancelledEUTRAItemExtIEs struct {
	List []TAICancelledEUTRAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P161 */
/* TAICancelledNRItemExtIEs */
type ProtocolExtensionContainerTAICancelledNRItemExtIEs struct {
	List []TAICancelledNRItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P162 */
/* TAIListForInactiveItemExtIEs */
type ProtocolExtensionContainerTAIListForInactiveItemExtIEs struct {
	List []TAIListForInactiveItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P163 */
/* TAIListForPagingItemExtIEs */
type ProtocolExtensionContainerTAIListForPagingItemExtIEs struct {
	List []TAIListForPagingItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P164 */
/* TargeteNBIDExtIEs */
type ProtocolExtensionContainerTargeteNBIDExtIEs struct {
	List []TargeteNBIDExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P165 */
/* TargetNGRANNodeToSourceNGRANNodeTransparentContainerExtIEs */
type ProtocolExtensionContainerTargetNGRANNodeToSourceNGRANNodeTransparentContainerExtIEs struct {
	List []TargetNGRANNodeToSourceNGRANNodeTransparentContainerExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P166 */
/* TargetRANNodeIDExtIEs */
type ProtocolExtensionContainerTargetRANNodeIDExtIEs struct {
	List []TargetRANNodeIDExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P167 */
/* TNLAssociationItemExtIEs */
type ProtocolExtensionContainerTNLAssociationItemExtIEs struct {
	List []TNLAssociationItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P168 */
/* TNLInformationItemExtIEs */
type ProtocolExtensionContainerTNLInformationItemExtIEs struct {
	List []TNLInformationItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P169 */
/* TNLMappingItemExtIEs */
type ProtocolExtensionContainerTNLMappingItemExtIEs struct {
	List []TNLMappingItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P170 */
/* TraceActivationExtIEs */
type ProtocolExtensionContainerTraceActivationExtIEs struct {
	List []TraceActivationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P171 */
/* UEAggregateMaximumBitRateExtIEs */
type ProtocolExtensionContainerUEAggregateMaximumBitRateExtIEs struct {
	List []UEAggregateMaximumBitRateExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P172 */
/* UEAssociatedLogicalNGConnectionItemExtIEs */
type ProtocolExtensionContainerUEAssociatedLogicalNGConnectionItemExtIEs struct {
	List []UEAssociatedLogicalNGConnectionItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P173 */
/* UENGAPIDPairExtIEs */
type ProtocolExtensionContainerUENGAPIDPairExtIEs struct {
	List []UENGAPIDPairExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P174 */
/* UEPresenceInAreaOfInterestItemExtIEs */
type ProtocolExtensionContainerUEPresenceInAreaOfInterestItemExtIEs struct {
	List []UEPresenceInAreaOfInterestItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P175 */
/* UERadioCapabilityForPagingExtIEs */
type ProtocolExtensionContainerUERadioCapabilityForPagingExtIEs struct {
	List []UERadioCapabilityForPagingExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P176 */
/* UESecurityCapabilitiesExtIEs */
type ProtocolExtensionContainerUESecurityCapabilitiesExtIEs struct {
	List []UESecurityCapabilitiesExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P177 */
/* ULNGUUPTNLModifyItemExtIEs */
type ProtocolExtensionContainerULNGUUPTNLModifyItemExtIEs struct {
	List []ULNGUUPTNLModifyItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P178 */
/* UnavailableGUAMIItemExtIEs */
type ProtocolExtensionContainerUnavailableGUAMIItemExtIEs struct {
	List []UnavailableGUAMIItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P179 */
/* UserLocationInformationEUTRAExtIEs */
type ProtocolExtensionContainerUserLocationInformationEUTRAExtIEs struct {
	List []UserLocationInformationEUTRAExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P180 */
/* UserLocationInformationN3IWFExtIEs */
type ProtocolExtensionContainerUserLocationInformationN3IWFExtIEs struct {
	List []UserLocationInformationN3IWFExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P181 */
/* UserLocationInformationNRExtIEs */
type ProtocolExtensionContainerUserLocationInformationNRExtIEs struct {
	List []UserLocationInformationNRExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P182 */
/* UserPlaneSecurityInformationExtIEs */
type ProtocolExtensionContainerUserPlaneSecurityInformationExtIEs struct {
	List []UserPlaneSecurityInformationExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P183 */
/* XnExtTLAItemExtIEs */
type ProtocolExtensionContainerXnExtTLAItemExtIEs struct {
	List []XnExtTLAItemExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolExtensionContainer_6500P184 */
/* XnTNLConfigurationInfoExtIEs */
type ProtocolExtensionContainerXnTNLConfigurationInfoExtIEs struct {
	List []XnTNLConfigurationInfoExtIEs `aper:"sizeLB:1,sizeUB:65535"`
}
