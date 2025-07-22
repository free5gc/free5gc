package message

// Used in AN-Parameter field for IE types
const (
	ANParametersTypeGUAMI              = 1
	ANParametersTypeSelectedPLMNID     = 2
	ANParametersTypeRequestedNSSAI     = 3
	ANParametersTypeEstablishmentCause = 4
)

// Used for checking if AN-Parameter length field is legal
const (
	ANParametersLenGUAMI    = 6
	ANParametersLenPLMNID   = 3
	ANParametersLenEstCause = 1
)

// Used in IE Establishment Cause field for cause types
const (
	EstablishmentCauseEmergency          = 0
	EstablishmentCauseHighPriorityAccess = 1
	EstablishmentCauseMO_Signalling      = 3
	EstablishmentCauseMO_Data            = 4
	EstablishmentCauseMPS_PriorityAccess = 8
	EstablishmentCauseMCS_PriorityAccess = 9
)
