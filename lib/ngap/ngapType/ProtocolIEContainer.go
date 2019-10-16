//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P0 */
/* PDUSessionResourceSetupRequestIEs */
type ProtocolIEContainerPDUSessionResourceSetupRequestIEs struct {
	List []PDUSessionResourceSetupRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P1 */
/* PDUSessionResourceSetupResponseIEs */
type ProtocolIEContainerPDUSessionResourceSetupResponseIEs struct {
	List []PDUSessionResourceSetupResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P2 */
/* PDUSessionResourceReleaseCommandIEs */
type ProtocolIEContainerPDUSessionResourceReleaseCommandIEs struct {
	List []PDUSessionResourceReleaseCommandIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P3 */
/* PDUSessionResourceReleaseResponseIEs */
type ProtocolIEContainerPDUSessionResourceReleaseResponseIEs struct {
	List []PDUSessionResourceReleaseResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P4 */
/* PDUSessionResourceModifyRequestIEs */
type ProtocolIEContainerPDUSessionResourceModifyRequestIEs struct {
	List []PDUSessionResourceModifyRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P5 */
/* PDUSessionResourceModifyResponseIEs */
type ProtocolIEContainerPDUSessionResourceModifyResponseIEs struct {
	List []PDUSessionResourceModifyResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P6 */
/* PDUSessionResourceNotifyIEs */
type ProtocolIEContainerPDUSessionResourceNotifyIEs struct {
	List []PDUSessionResourceNotifyIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P7 */
/* PDUSessionResourceModifyIndicationIEs */
type ProtocolIEContainerPDUSessionResourceModifyIndicationIEs struct {
	List []PDUSessionResourceModifyIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P8 */
/* PDUSessionResourceModifyConfirmIEs */
type ProtocolIEContainerPDUSessionResourceModifyConfirmIEs struct {
	List []PDUSessionResourceModifyConfirmIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P9 */
/* InitialContextSetupRequestIEs */
type ProtocolIEContainerInitialContextSetupRequestIEs struct {
	List []InitialContextSetupRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P10 */
/* InitialContextSetupResponseIEs */
type ProtocolIEContainerInitialContextSetupResponseIEs struct {
	List []InitialContextSetupResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P11 */
/* InitialContextSetupFailureIEs */
type ProtocolIEContainerInitialContextSetupFailureIEs struct {
	List []InitialContextSetupFailureIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P12 */
/* UEContextReleaseRequestIEs */
type ProtocolIEContainerUEContextReleaseRequestIEs struct {
	List []UEContextReleaseRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P13 */
/* UEContextReleaseCommandIEs */
type ProtocolIEContainerUEContextReleaseCommandIEs struct {
	List []UEContextReleaseCommandIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P14 */
/* UEContextReleaseCompleteIEs */
type ProtocolIEContainerUEContextReleaseCompleteIEs struct {
	List []UEContextReleaseCompleteIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P15 */
/* UEContextModificationRequestIEs */
type ProtocolIEContainerUEContextModificationRequestIEs struct {
	List []UEContextModificationRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P16 */
/* UEContextModificationResponseIEs */
type ProtocolIEContainerUEContextModificationResponseIEs struct {
	List []UEContextModificationResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P17 */
/* UEContextModificationFailureIEs */
type ProtocolIEContainerUEContextModificationFailureIEs struct {
	List []UEContextModificationFailureIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P18 */
/* RRCInactiveTransitionReportIEs */
type ProtocolIEContainerRRCInactiveTransitionReportIEs struct {
	List []RRCInactiveTransitionReportIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P19 */
/* HandoverRequiredIEs */
type ProtocolIEContainerHandoverRequiredIEs struct {
	List []HandoverRequiredIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P20 */
/* HandoverCommandIEs */
type ProtocolIEContainerHandoverCommandIEs struct {
	List []HandoverCommandIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P21 */
/* HandoverPreparationFailureIEs */
type ProtocolIEContainerHandoverPreparationFailureIEs struct {
	List []HandoverPreparationFailureIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P22 */
/* HandoverRequestIEs */
type ProtocolIEContainerHandoverRequestIEs struct {
	List []HandoverRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P23 */
/* HandoverRequestAcknowledgeIEs */
type ProtocolIEContainerHandoverRequestAcknowledgeIEs struct {
	List []HandoverRequestAcknowledgeIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P24 */
/* HandoverFailureIEs */
type ProtocolIEContainerHandoverFailureIEs struct {
	List []HandoverFailureIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P25 */
/* HandoverNotifyIEs */
type ProtocolIEContainerHandoverNotifyIEs struct {
	List []HandoverNotifyIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P26 */
/* PathSwitchRequestIEs */
type ProtocolIEContainerPathSwitchRequestIEs struct {
	List []PathSwitchRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P27 */
/* PathSwitchRequestAcknowledgeIEs */
type ProtocolIEContainerPathSwitchRequestAcknowledgeIEs struct {
	List []PathSwitchRequestAcknowledgeIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P28 */
/* PathSwitchRequestFailureIEs */
type ProtocolIEContainerPathSwitchRequestFailureIEs struct {
	List []PathSwitchRequestFailureIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P29 */
/* HandoverCancelIEs */
type ProtocolIEContainerHandoverCancelIEs struct {
	List []HandoverCancelIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P30 */
/* HandoverCancelAcknowledgeIEs */
type ProtocolIEContainerHandoverCancelAcknowledgeIEs struct {
	List []HandoverCancelAcknowledgeIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P31 */
/* UplinkRANStatusTransferIEs */
type ProtocolIEContainerUplinkRANStatusTransferIEs struct {
	List []UplinkRANStatusTransferIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P32 */
/* DownlinkRANStatusTransferIEs */
type ProtocolIEContainerDownlinkRANStatusTransferIEs struct {
	List []DownlinkRANStatusTransferIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P33 */
/* PagingIEs */
type ProtocolIEContainerPagingIEs struct {
	List []PagingIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P34 */
/* InitialUEMessageIEs */
type ProtocolIEContainerInitialUEMessageIEs struct {
	List []InitialUEMessageIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P35 */
/* DownlinkNASTransportIEs */
type ProtocolIEContainerDownlinkNASTransportIEs struct {
	List []DownlinkNASTransportIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P36 */
/* UplinkNASTransportIEs */
type ProtocolIEContainerUplinkNASTransportIEs struct {
	List []UplinkNASTransportIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P37 */
/* NASNonDeliveryIndicationIEs */
type ProtocolIEContainerNASNonDeliveryIndicationIEs struct {
	List []NASNonDeliveryIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P38 */
/* RerouteNASRequestIEs */
type ProtocolIEContainerRerouteNASRequestIEs struct {
	List []RerouteNASRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P39 */
/* NGSetupRequestIEs */
type ProtocolIEContainerNGSetupRequestIEs struct {
	List []NGSetupRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P40 */
/* NGSetupResponseIEs */
type ProtocolIEContainerNGSetupResponseIEs struct {
	List []NGSetupResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P41 */
/* NGSetupFailureIEs */
type ProtocolIEContainerNGSetupFailureIEs struct {
	List []NGSetupFailureIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P42 */
/* RANConfigurationUpdateIEs */
type ProtocolIEContainerRANConfigurationUpdateIEs struct {
	List []RANConfigurationUpdateIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P43 */
/* RANConfigurationUpdateAcknowledgeIEs */
type ProtocolIEContainerRANConfigurationUpdateAcknowledgeIEs struct {
	List []RANConfigurationUpdateAcknowledgeIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P44 */
/* RANConfigurationUpdateFailureIEs */
type ProtocolIEContainerRANConfigurationUpdateFailureIEs struct {
	List []RANConfigurationUpdateFailureIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P45 */
/* AMFConfigurationUpdateIEs */
type ProtocolIEContainerAMFConfigurationUpdateIEs struct {
	List []AMFConfigurationUpdateIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P46 */
/* AMFConfigurationUpdateAcknowledgeIEs */
type ProtocolIEContainerAMFConfigurationUpdateAcknowledgeIEs struct {
	List []AMFConfigurationUpdateAcknowledgeIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P47 */
/* AMFConfigurationUpdateFailureIEs */
type ProtocolIEContainerAMFConfigurationUpdateFailureIEs struct {
	List []AMFConfigurationUpdateFailureIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P48 */
/* AMFStatusIndicationIEs */
type ProtocolIEContainerAMFStatusIndicationIEs struct {
	List []AMFStatusIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P49 */
/* NGResetIEs */
type ProtocolIEContainerNGResetIEs struct {
	List []NGResetIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P50 */
/* NGResetAcknowledgeIEs */
type ProtocolIEContainerNGResetAcknowledgeIEs struct {
	List []NGResetAcknowledgeIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P51 */
/* ErrorIndicationIEs */
type ProtocolIEContainerErrorIndicationIEs struct {
	List []ErrorIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P52 */
/* OverloadStartIEs */
type ProtocolIEContainerOverloadStartIEs struct {
	List []OverloadStartIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P53 */
/* OverloadStopIEs */
type ProtocolIEContainerOverloadStopIEs struct {
	List []OverloadStopIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P54 */
/* UplinkRANConfigurationTransferIEs */
type ProtocolIEContainerUplinkRANConfigurationTransferIEs struct {
	List []UplinkRANConfigurationTransferIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P55 */
/* DownlinkRANConfigurationTransferIEs */
type ProtocolIEContainerDownlinkRANConfigurationTransferIEs struct {
	List []DownlinkRANConfigurationTransferIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P56 */
/* WriteReplaceWarningRequestIEs */
type ProtocolIEContainerWriteReplaceWarningRequestIEs struct {
	List []WriteReplaceWarningRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P57 */
/* WriteReplaceWarningResponseIEs */
type ProtocolIEContainerWriteReplaceWarningResponseIEs struct {
	List []WriteReplaceWarningResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P58 */
/* PWSCancelRequestIEs */
type ProtocolIEContainerPWSCancelRequestIEs struct {
	List []PWSCancelRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P59 */
/* PWSCancelResponseIEs */
type ProtocolIEContainerPWSCancelResponseIEs struct {
	List []PWSCancelResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P60 */
/* PWSRestartIndicationIEs */
type ProtocolIEContainerPWSRestartIndicationIEs struct {
	List []PWSRestartIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P61 */
/* PWSFailureIndicationIEs */
type ProtocolIEContainerPWSFailureIndicationIEs struct {
	List []PWSFailureIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P62 */
/* DownlinkUEAssociatedNRPPaTransportIEs */
type ProtocolIEContainerDownlinkUEAssociatedNRPPaTransportIEs struct {
	List []DownlinkUEAssociatedNRPPaTransportIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P63 */
/* UplinkUEAssociatedNRPPaTransportIEs */
type ProtocolIEContainerUplinkUEAssociatedNRPPaTransportIEs struct {
	List []UplinkUEAssociatedNRPPaTransportIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P64 */
/* DownlinkNonUEAssociatedNRPPaTransportIEs */
type ProtocolIEContainerDownlinkNonUEAssociatedNRPPaTransportIEs struct {
	List []DownlinkNonUEAssociatedNRPPaTransportIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P65 */
/* UplinkNonUEAssociatedNRPPaTransportIEs */
type ProtocolIEContainerUplinkNonUEAssociatedNRPPaTransportIEs struct {
	List []UplinkNonUEAssociatedNRPPaTransportIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P66 */
/* TraceStartIEs */
type ProtocolIEContainerTraceStartIEs struct {
	List []TraceStartIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P67 */
/* TraceFailureIndicationIEs */
type ProtocolIEContainerTraceFailureIndicationIEs struct {
	List []TraceFailureIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P68 */
/* DeactivateTraceIEs */
type ProtocolIEContainerDeactivateTraceIEs struct {
	List []DeactivateTraceIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P69 */
/* CellTrafficTraceIEs */
type ProtocolIEContainerCellTrafficTraceIEs struct {
	List []CellTrafficTraceIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P70 */
/* LocationReportingControlIEs */
type ProtocolIEContainerLocationReportingControlIEs struct {
	List []LocationReportingControlIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P71 */
/* LocationReportingFailureIndicationIEs */
type ProtocolIEContainerLocationReportingFailureIndicationIEs struct {
	List []LocationReportingFailureIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P72 */
/* LocationReportIEs */
type ProtocolIEContainerLocationReportIEs struct {
	List []LocationReportIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P73 */
/* UETNLABindingReleaseRequestIEs */
type ProtocolIEContainerUETNLABindingReleaseRequestIEs struct {
	List []UETNLABindingReleaseRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P74 */
/* UERadioCapabilityInfoIndicationIEs */
type ProtocolIEContainerUERadioCapabilityInfoIndicationIEs struct {
	List []UERadioCapabilityInfoIndicationIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P75 */
/* UERadioCapabilityCheckRequestIEs */
type ProtocolIEContainerUERadioCapabilityCheckRequestIEs struct {
	List []UERadioCapabilityCheckRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P76 */
/* UERadioCapabilityCheckResponseIEs */
type ProtocolIEContainerUERadioCapabilityCheckResponseIEs struct {
	List []UERadioCapabilityCheckResponseIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P77 */
/* PDUSessionResourceModifyRequestTransferIEs */
type ProtocolIEContainerPDUSessionResourceModifyRequestTransferIEs struct {
	List []PDUSessionResourceModifyRequestTransferIEs `aper:"sizeLB:0,sizeUB:65535"`
}

/* Sequence of = 35, FULL Name = struct ProtocolIE_Container_6449P78 */
/* PDUSessionResourceSetupRequestTransferIEs */
type ProtocolIEContainerPDUSessionResourceSetupRequestTransferIEs struct {
	List []PDUSessionResourceSetupRequestTransferIEs `aper:"sizeLB:0,sizeUB:65535"`
}
