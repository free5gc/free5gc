//go:binary-only-package

package nasMessage

const (
	ULNASTransportRequestTypeInitialRequest              uint8 = 1
	ULNASTransportRequestTypeExistingPduSession          uint8 = 2
	ULNASTransportRequestTypeInitialEmergencyRequest     uint8 = 3
	ULNASTransportRequestTypeExistingEmergencyPduSession uint8 = 4
	ULNASTransportRequestTypeExistingReserved            uint8 = 7
)

const (
	PayloadContainerTypeN1SMInfo          uint8 = 0x01
	PayloadContainerTypeSMS               uint8 = 0x02
	PayloadContainerTypeLPP               uint8 = 0x03
	PayloadContainerTypeSOR               uint8 = 0x04
	PayloadContainerTypeUEPolicy          uint8 = 0x05
	PayloadContainerTypeUEParameterUpdate uint8 = 0x06
	PayloadContainerTypeMultiplePayload   uint8 = 0x0f
)

const (
	Cause5GSMInsufficientResources                                       uint8 = 0x1a
	Cause5GSMMissingOrUnknownDNN                                         uint8 = 0x1b
	Cause5GSMUnknownPDUSessionType                                       uint8 = 0x1c
	Cause5GSMUserAuthenticationOrAuthorizationFailed                     uint8 = 0x1d
	Cause5GSMRequestRejectedUnspecified                                  uint8 = 0x1f
	Cause5GSMServiceOptionTemporarilyOutOfOrder                          uint8 = 0x22
	Cause5GSMPTIAlreadyInUse                                             uint8 = 0x23
	Cause5GSMRegularDeactivation                                         uint8 = 0x24
	Cause5GSMReactivationRequested                                       uint8 = 0x27
	Cause5GSMInvalidPDUSessionIdentity                                   uint8 = 0x2b
	Cause5GSMSemanticErrorsInPacketFilter                                uint8 = 0x2c
	Cause5GSMSyntacticalErrorInPacketFilter                              uint8 = 0x2d
	Cause5GSMOutOfLADNServiceArea                                        uint8 = 0x2e
	Cause5GSMPTIMismatch                                                 uint8 = 0x2f
	Cause5GSMPDUSessionTypeIPv4OnlyAllowed                               uint8 = 0x32
	Cause5GSMPDUSessionTypeIPv6OnlyAllowed                               uint8 = 0x33
	Cause5GSMPDUSessionDoesNotExist                                      uint8 = 0x36
	Cause5GSMInsufficientResourcesForSpecificSliceAndDNN                 uint8 = 0x43
	Cause5GSMNotSupportedSSCMode                                         uint8 = 0x44
	Cause5GSMInsufficientResourcesForSpecificSlice                       uint8 = 0x45
	Cause5GSMMissingOrUnknownDNNInASlice                                 uint8 = 0x46
	Cause5GSMInvalidPTIValue                                             uint8 = 0x51
	Cause5GSMMaximumDataRatePerUEForUserPlaneIntegrityProtectionIsTooLow uint8 = 0x52
	Cause5GSMSemanticErrorInTheQoSOperation                              uint8 = 0x53
	Cause5GSMSyntacticalErrorInTheQoSOperation                           uint8 = 0x54
	Cause5GSMInvalidMappedEPSBearerIdentity                              uint8 = 0x55
	Cause5GSMSemanticallyIncorrectMessage                                uint8 = 0x5f
	Cause5GSMInvalidMandatoryInformation                                 uint8 = 0x60
	Cause5GSMMessageTypeNonExistentOrNotImplemented                      uint8 = 0x61
	Cause5GSMMessageTypeNotCompatibleWithTheProtocolState                uint8 = 0x62
	Cause5GSMInformationElementNonExistentOrNotImplemented               uint8 = 0x63
	Cause5GSMConditionalIEError                                          uint8 = 0x64
	Cause5GSMMessageNotCompatibleWithTheProtocolState                    uint8 = 0x65
	Cause5GSMProtocolErrorUnspecified                                    uint8 = 0x6f
)

const (
	Cause5GMMIllegalUE                                      uint8 = 0x03
	Cause5GMMPEINotAccepted                                 uint8 = 0x05
	Cause5GMMIllegalME                                      uint8 = 0x06
	Cause5GMM5GSServicesNotAllowed                          uint8 = 0x07
	Cause5GMMUEIdentityCannotBeDerivedByTheNetwork          uint8 = 0x09
	Cause5GMMImplicitlyDeregistered                         uint8 = 0x0a
	Cause5GMMPLMNNotAllowed                                 uint8 = 0x0b
	Cause5GMMTrackingAreaNotAllowed                         uint8 = 0x0c
	Cause5GMMRoamingNotAllowedInThisTrackingArea            uint8 = 0x0d
	Cause5GMMNoSuitableCellsInTrackingArea                  uint8 = 0x0f
	Cause5GMMMACFailure                                     uint8 = 0x14
	Cause5GMMSynchFailure                                   uint8 = 0x15
	Cause5GMMCongestion                                     uint8 = 0x16
	Cause5GMMUESecurityCapabilitiesMismatch                 uint8 = 0x17
	Cause5GMMSecurityModeRejectedUnspecified                uint8 = 0x18
	Cause5GMMNon5GAuthenticationUnacceptable                uint8 = 0x1a
	Cause5GMMN1ModeNotAllowed                               uint8 = 0x1b
	Cause5GMMRestrictedServiceArea                          uint8 = 0x1c
	Cause5GMMLADNNotAvailable                               uint8 = 0x2b
	Cause5GMMMaximumNumberOfPDUSessionsReached              uint8 = 0x41
	Cause5GMMInsufficientResourcesForSpecificSliceAndDNN    uint8 = 0x43
	Cause5GMMInsufficientResourcesForSpecificSlice          uint8 = 0x45
	Cause5GMMngKSIAlreadyInUse                              uint8 = 0x47
	Cause5GMMNon3GPPAccessTo5GCNNotAllowed                  uint8 = 0x48
	Cause5GMMServingNetworkNotAuthorized                    uint8 = 0x49
	Cause5GMMPayloadWasNotForwarded                         uint8 = 0x5a
	Cause5GMMDNNNotSupportedOrNotSubscribedInTheSlice       uint8 = 0x5b
	Cause5GMMInsufficientUserPlaneResourcesForThePDUSession uint8 = 0x5c
	Cause5GMMSemanticallyIncorrectMessage                   uint8 = 0x5f
	Cause5GMMInvalidMandatoryInformation                    uint8 = 0x60
	Cause5GMMMessageTypeNonExistentOrNotImplemented         uint8 = 0x61
	Cause5GMMMessageTypeNotCompatibleWithTheProtocolState   uint8 = 0x62
	Cause5GMMInformationElementNonExistentOrNotImplemented  uint8 = 0x63
	Cause5GMMConditionalIEError                             uint8 = 0x64
	Cause5GMMMessageNotCompatibleWithTheProtocolState       uint8 = 0x65
	Cause5GMMProtocolErrorUnspecified                       uint8 = 0x6f
)

// TS 24.501 9.11.3.7
const (
	RegistrationType5GSInitialRegistration          uint8 = 0x01
	RegistrationType5GSMobilityRegistrationUpdating uint8 = 0x02
	RegistrationType5GSPeriodicRegistrationUpdating uint8 = 0x03
	RegistrationType5GSEmergencyRegistration        uint8 = 0x04
	RegistrationType5GSReserved                     uint8 = 0x07
)

// TS 24.501 9.11.3.7
const (
	FollowOnRequestNoPending uint8 = 0x00
	FollowOnRequestPending   uint8 = 0x01
)

const (
	MobileIdentity5GSTypeNoIdentity uint8 = 0x00
	MobileIdentity5GSTypeSuci       uint8 = 0x01
	MobileIdentity5GSType5gGuti     uint8 = 0x02
	MobileIdentity5GSTypeImei       uint8 = 0x03
	MobileIdentity5GSType5gSTmsi    uint8 = 0x04
	MobileIdentity5GSTypeImeisv     uint8 = 0x05
)

// TS 24.501 9.11.3.2A
const (
	DRXValueNotSpecified  uint8 = 0x00
	DRXcycleParameterT32  uint8 = 0x01
	DRXcycleParameterT64  uint8 = 0x02
	DRXcycleParameterT128 uint8 = 0x03
	DRXcycleParameterT256 uint8 = 0x04
)

// TS 24.501 9.11.3.32
const (
	TypeOfSecurityContextFlagNative uint8 = 0x00
	TypeOfSecurityContextFlagMapped uint8 = 0x01
)

// TS 24.501 9.11.3.32
const (
	NasKeySetIdentifierNoKeyIsAvailable int32 = 0x07
)

// TS 24.501 9.11.3.11
const (
	AccessType3GPP    uint8 = 0x01
	AccessTypeNon3GPP uint8 = 0x02
	AccessTypeBoth    uint8 = 0x03
)

// TS 24.501 9.11.3.50
const (
	ServiceTypeSignalling                uint8 = 0x00
	ServiceTypeData                      uint8 = 0x01
	ServiceTypeMobileTerminatedServices  uint8 = 0x02
	ServiceTypeEmergencyServices         uint8 = 0x03
	ServiceTypeEmergencyServicesFallback uint8 = 0x04
	ServiceTypeHighPriorityAccess        uint8 = 0x05
)

// TS 24.501 9.11.3.20
const (
	ReRegistrationNotRequired uint8 = 0x00
	ReRegistrationRequired    uint8 = 0x01
)

// TS 24.501 9.11.3.28 TS 24.008 10.5.5.10
const (
	IMEISVNotRequested uint8 = 0x00
	IMEISVRequested    uint8 = 0x01
)

// TS 24.501 9.11.3.6
const (
	RegistrationResult5GS3GPPAccess           uint8 = 0x01
	RegistrationResult5GSNon3GPPAccess        uint8 = 0x02
	RegistrationResult5GS3GPPandNon3GPPAccess uint8 = 0x03
)

// TS 24.501 9.11.3.6
const (
	SMSOverNasNotAllowed uint8 = 0x00
	SMSOverNasAllowed    uint8 = 0x01
)

// TS 24.501 9.11.3.46
const (
	SnssaiNotAvailableInCurrentPlmn             uint8 = 0x00
	SnssaiNotAvailableInCurrentRegistrationArea uint8 = 0x01
)

// TS 24.008 10.5.7.4a
const (
	GPRSTimer3UnitMultiplesOf10Minutes uint8 = 0x00
	GPRSTimer3UnitMultiplesOf1Hour     uint8 = 0x01
	GPRSTimer3UnitMultiplesOf10Hours   uint8 = 0x02
	GPRSTimer3UnitMultiplesOf2Seconds  uint8 = 0x03
	GPRSTimer3UnitMultiplesOf30Seconds uint8 = 0x04
	GPRSTimer3UnitMultiplesOf1Minute   uint8 = 0x05
)

// TS 24.501 9.11.3.9A
const (
	NGRanRadioCapabilityUpdateNotNeeded uint8 = 0x00
	NGRanRadioCapabilityUpdateNeeded    uint8 = 0x01
)

// TS 24.501 9.11.3.49
const (
	AllowedTypeAllowedArea    uint8 = 0x00
	AllowedTypeNonAllowedArea uint8 = 0x01
)
