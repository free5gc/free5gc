//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type UnsuccessfulOutcome struct {
	ProcedureCode ProcedureCode
	Criticality   Criticality
	Value         UnsuccessfulOutcomeValue `aper:"openType,referenceFieldName:ProcedureCode"`
}

const (
	UnsuccessfulOutcomePresentNothing int = iota /* No components present */
	UnsuccessfulOutcomePresentAMFConfigurationUpdateFailure
	UnsuccessfulOutcomePresentHandoverPreparationFailure
	UnsuccessfulOutcomePresentHandoverFailure
	UnsuccessfulOutcomePresentInitialContextSetupFailure
	UnsuccessfulOutcomePresentNGSetupFailure
	UnsuccessfulOutcomePresentPathSwitchRequestFailure
	UnsuccessfulOutcomePresentRANConfigurationUpdateFailure
	UnsuccessfulOutcomePresentUEContextModificationFailure
)

type UnsuccessfulOutcomeValue struct {
	Present                       int
	AMFConfigurationUpdateFailure *AMFConfigurationUpdateFailure `aper:"valueExt,referenceFieldValue:0"`
	HandoverPreparationFailure    *HandoverPreparationFailure    `aper:"valueExt,referenceFieldValue:12"`
	HandoverFailure               *HandoverFailure               `aper:"valueExt,referenceFieldValue:13"`
	InitialContextSetupFailure    *InitialContextSetupFailure    `aper:"valueExt,referenceFieldValue:14"`
	NGSetupFailure                *NGSetupFailure                `aper:"valueExt,referenceFieldValue:21"`
	PathSwitchRequestFailure      *PathSwitchRequestFailure      `aper:"valueExt,referenceFieldValue:25"`
	RANConfigurationUpdateFailure *RANConfigurationUpdateFailure `aper:"valueExt,referenceFieldValue:35"`
	UEContextModificationFailure  *UEContextModificationFailure  `aper:"valueExt,referenceFieldValue:40"`
}
