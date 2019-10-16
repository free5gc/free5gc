//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type CriticalityDiagnostics struct {
	ProcedureCode             *ProcedureCode                                          `aper:"optional"`
	TriggeringMessage         *TriggeringMessage                                      `aper:"optional"`
	ProcedureCriticality      *Criticality                                            `aper:"optional"`
	IEsCriticalityDiagnostics *CriticalityDiagnosticsIEList                           `aper:"optional"`
	IEExtensions              *ProtocolExtensionContainerCriticalityDiagnosticsExtIEs `aper:"optional"`
}
