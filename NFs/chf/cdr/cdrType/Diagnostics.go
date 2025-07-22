package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const (
	DiagnosticsPresentNothing int = iota /* No components present */
	DiagnosticsPresentGsm0408Cause
	DiagnosticsPresentGsm0902MapErrorValue
	DiagnosticsPresentItuTQ767Cause
	DiagnosticsPresentNetworkSpecificCause
	DiagnosticsPresentManufacturerSpecificCause
	DiagnosticsPresentPositionMethodFailureCause
	DiagnosticsPresentUnauthorizedLCSClientCause
	DiagnosticsPresentDiameterResultCodeAndExperimentalResult
)

type Diagnostics struct {
	Present                                 int                              /* Choice Type */
	Gsm0408Cause                            *int64                           `ber:"tagNum:0"`
	Gsm0902MapErrorValue                    *int64                           `ber:"tagNum:1"`
	ItuTQ767Cause                           *int64                           `ber:"tagNum:2"`
	NetworkSpecificCause                    *ManagementExtension             `ber:"tagNum:3"`
	ManufacturerSpecificCause               *ManagementExtension             `ber:"tagNum:4"`
	PositionMethodFailureCause              *PositionMethodFailureDiagnostic `ber:"tagNum:5"`
	UnauthorizedLCSClientCause              *UnauthorizedLCSClientDiagnostic `ber:"tagNum:6"`
	DiameterResultCodeAndExperimentalResult *int64                           `ber:"tagNum:7"`
}
