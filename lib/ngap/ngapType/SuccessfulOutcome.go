//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type SuccessfulOutcome struct {
	ProcedureCode ProcedureCode
	Criticality   Criticality
	Value         SuccessfulOutcomeValue `aper:"openType,referenceFieldName:ProcedureCode"`
}

const (
	SuccessfulOutcomePresentNothing int = iota /* No components present */
	SuccessfulOutcomePresentAMFConfigurationUpdateAcknowledge
	SuccessfulOutcomePresentHandoverCancelAcknowledge
	SuccessfulOutcomePresentHandoverCommand
	SuccessfulOutcomePresentHandoverRequestAcknowledge
	SuccessfulOutcomePresentInitialContextSetupResponse
	SuccessfulOutcomePresentNGResetAcknowledge
	SuccessfulOutcomePresentNGSetupResponse
	SuccessfulOutcomePresentPathSwitchRequestAcknowledge
	SuccessfulOutcomePresentPDUSessionResourceModifyResponse
	SuccessfulOutcomePresentPDUSessionResourceModifyConfirm
	SuccessfulOutcomePresentPDUSessionResourceReleaseResponse
	SuccessfulOutcomePresentPDUSessionResourceSetupResponse
	SuccessfulOutcomePresentPWSCancelResponse
	SuccessfulOutcomePresentRANConfigurationUpdateAcknowledge
	SuccessfulOutcomePresentUEContextModificationResponse
	SuccessfulOutcomePresentUEContextReleaseComplete
	SuccessfulOutcomePresentUERadioCapabilityCheckResponse
	SuccessfulOutcomePresentWriteReplaceWarningResponse
)

type SuccessfulOutcomeValue struct {
	Present                           int
	AMFConfigurationUpdateAcknowledge *AMFConfigurationUpdateAcknowledge `aper:"valueExt,referenceFieldValue:0"`
	HandoverCancelAcknowledge         *HandoverCancelAcknowledge         `aper:"valueExt,referenceFieldValue:10"`
	HandoverCommand                   *HandoverCommand                   `aper:"valueExt,referenceFieldValue:12"`
	HandoverRequestAcknowledge        *HandoverRequestAcknowledge        `aper:"valueExt,referenceFieldValue:13"`
	InitialContextSetupResponse       *InitialContextSetupResponse       `aper:"valueExt,referenceFieldValue:14"`
	NGResetAcknowledge                *NGResetAcknowledge                `aper:"valueExt,referenceFieldValue:20"`
	NGSetupResponse                   *NGSetupResponse                   `aper:"valueExt,referenceFieldValue:21"`
	PathSwitchRequestAcknowledge      *PathSwitchRequestAcknowledge      `aper:"valueExt,referenceFieldValue:25"`
	PDUSessionResourceModifyResponse  *PDUSessionResourceModifyResponse  `aper:"valueExt,referenceFieldValue:26"`
	PDUSessionResourceModifyConfirm   *PDUSessionResourceModifyConfirm   `aper:"valueExt,referenceFieldValue:27"`
	PDUSessionResourceReleaseResponse *PDUSessionResourceReleaseResponse `aper:"valueExt,referenceFieldValue:28"`
	PDUSessionResourceSetupResponse   *PDUSessionResourceSetupResponse   `aper:"valueExt,referenceFieldValue:29"`
	PWSCancelResponse                 *PWSCancelResponse                 `aper:"valueExt,referenceFieldValue:32"`
	RANConfigurationUpdateAcknowledge *RANConfigurationUpdateAcknowledge `aper:"valueExt,referenceFieldValue:35"`
	UEContextModificationResponse     *UEContextModificationResponse     `aper:"valueExt,referenceFieldValue:40"`
	UEContextReleaseComplete          *UEContextReleaseComplete          `aper:"valueExt,referenceFieldValue:41"`
	UERadioCapabilityCheckResponse    *UERadioCapabilityCheckResponse    `aper:"valueExt,referenceFieldValue:43"`
	WriteReplaceWarningResponse       *WriteReplaceWarningResponse       `aper:"valueExt,referenceFieldValue:51"`
}
