//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type ProcedureCode struct {
	Value int64 `aper:"valueLB:0,valueUB:255"`
}

const ProcedureCodeAMFConfigurationUpdate int64 = 0
const ProcedureCodeAMFStatusIndication int64 = 1
const ProcedureCodeCellTrafficTrace int64 = 2
const ProcedureCodeDeactivateTrace int64 = 3
const ProcedureCodeDownlinkNASTransport int64 = 4
const ProcedureCodeDownlinkNonUEAssociatedNRPPaTransport int64 = 5
const ProcedureCodeDownlinkRANConfigurationTransfer int64 = 6
const ProcedureCodeDownlinkRANStatusTransfer int64 = 7
const ProcedureCodeDownlinkUEAssociatedNRPPaTransport int64 = 8
const ProcedureCodeErrorIndication int64 = 9
const ProcedureCodeHandoverCancel int64 = 10
const ProcedureCodeHandoverNotification int64 = 11
const ProcedureCodeHandoverPreparation int64 = 12
const ProcedureCodeHandoverResourceAllocation int64 = 13
const ProcedureCodeInitialContextSetup int64 = 14
const ProcedureCodeInitialUEMessage int64 = 15
const ProcedureCodeLocationReportingControl int64 = 16
const ProcedureCodeLocationReportingFailureIndication int64 = 17
const ProcedureCodeLocationReport int64 = 18
const ProcedureCodeNASNonDeliveryIndication int64 = 19
const ProcedureCodeNGReset int64 = 20
const ProcedureCodeNGSetup int64 = 21
const ProcedureCodeOverloadStart int64 = 22
const ProcedureCodeOverloadStop int64 = 23
const ProcedureCodePaging int64 = 24
const ProcedureCodePathSwitchRequest int64 = 25
const ProcedureCodePDUSessionResourceModify int64 = 26
const ProcedureCodePDUSessionResourceModifyIndication int64 = 27
const ProcedureCodePDUSessionResourceRelease int64 = 28
const ProcedureCodePDUSessionResourceSetup int64 = 29
const ProcedureCodePDUSessionResourceNotify int64 = 30
const ProcedureCodePrivateMessage int64 = 31
const ProcedureCodePWSCancel int64 = 32
const ProcedureCodePWSFailureIndication int64 = 33
const ProcedureCodePWSRestartIndication int64 = 34
const ProcedureCodeRANConfigurationUpdate int64 = 35
const ProcedureCodeRerouteNASRequest int64 = 36
const ProcedureCodeRRCInactiveTransitionReport int64 = 37
const ProcedureCodeTraceFailureIndication int64 = 38
const ProcedureCodeTraceStart int64 = 39
const ProcedureCodeUEContextModification int64 = 40
const ProcedureCodeUEContextRelease int64 = 41
const ProcedureCodeUEContextReleaseRequest int64 = 42
const ProcedureCodeUERadioCapabilityCheck int64 = 43
const ProcedureCodeUERadioCapabilityInfoIndication int64 = 44
const ProcedureCodeUETNLABindingRelease int64 = 45
const ProcedureCodeUplinkNASTransport int64 = 46
const ProcedureCodeUplinkNonUEAssociatedNRPPaTransport int64 = 47
const ProcedureCodeUplinkRANConfigurationTransfer int64 = 48
const ProcedureCodeUplinkRANStatusTransfer int64 = 49
const ProcedureCodeUplinkUEAssociatedNRPPaTransport int64 = 50
const ProcedureCodeWriteReplaceWarning int64 = 51
