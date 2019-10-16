//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

type AMFPagingTargetExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       AMFPagingTargetExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFPagingTargetExtIEsPresentNothing int = iota /* No components present */
)

type AMFPagingTargetExtIEsValue struct {
	Present int
}

type BroadcastCancelledAreaListExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       BroadcastCancelledAreaListExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	BroadcastCancelledAreaListExtIEsPresentNothing int = iota /* No components present */
)

type BroadcastCancelledAreaListExtIEsValue struct {
	Present int
}

type BroadcastCompletedAreaListExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       BroadcastCompletedAreaListExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	BroadcastCompletedAreaListExtIEsPresentNothing int = iota /* No components present */
)

type BroadcastCompletedAreaListExtIEsValue struct {
	Present int
}

type CauseExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       CauseExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CauseExtIEsPresentNothing int = iota /* No components present */
)

type CauseExtIEsValue struct {
	Present int
}

type CellIDListForRestartExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       CellIDListForRestartExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CellIDListForRestartExtIEsPresentNothing int = iota /* No components present */
)

type CellIDListForRestartExtIEsValue struct {
	Present int
}

type CPTransportLayerInformationExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       CPTransportLayerInformationExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CPTransportLayerInformationExtIEsPresentNothing int = iota /* No components present */
)

type CPTransportLayerInformationExtIEsValue struct {
	Present int
}

type DRBStatusDLExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       DRBStatusDLExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DRBStatusDLExtIEsPresentNothing int = iota /* No components present */
)

type DRBStatusDLExtIEsValue struct {
	Present int
}

type DRBStatusULExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       DRBStatusULExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DRBStatusULExtIEsPresentNothing int = iota /* No components present */
)

type DRBStatusULExtIEsValue struct {
	Present int
}

type GlobalRANNodeIDExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       GlobalRANNodeIDExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	GlobalRANNodeIDExtIEsPresentNothing int = iota /* No components present */
)

type GlobalRANNodeIDExtIEsValue struct {
	Present int
}

type GNBIDExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       GNBIDExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	GNBIDExtIEsPresentNothing int = iota /* No components present */
)

type GNBIDExtIEsValue struct {
	Present int
}

type LastVisitedCellInformationExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       LastVisitedCellInformationExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	LastVisitedCellInformationExtIEsPresentNothing int = iota /* No components present */
)

type LastVisitedCellInformationExtIEsValue struct {
	Present int
}

type N3IWFIDExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       N3IWFIDExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	N3IWFIDExtIEsPresentNothing int = iota /* No components present */
)

type N3IWFIDExtIEsValue struct {
	Present int
}

type NgENBIDExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       NgENBIDExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NgENBIDExtIEsPresentNothing int = iota /* No components present */
)

type NgENBIDExtIEsValue struct {
	Present int
}

type NGRANCGIExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       NGRANCGIExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NGRANCGIExtIEsPresentNothing int = iota /* No components present */
)

type NGRANCGIExtIEsValue struct {
	Present int
}

type OverloadResponseExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       OverloadResponseExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	OverloadResponseExtIEsPresentNothing int = iota /* No components present */
)

type OverloadResponseExtIEsValue struct {
	Present int
}

type PWSFailedCellIDListExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PWSFailedCellIDListExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PWSFailedCellIDListExtIEsPresentNothing int = iota /* No components present */
)

type PWSFailedCellIDListExtIEsValue struct {
	Present int
}

type QosCharacteristicsExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       QosCharacteristicsExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	QosCharacteristicsExtIEsPresentNothing int = iota /* No components present */
)

type QosCharacteristicsExtIEsValue struct {
	Present int
}

type ResetTypeExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       ResetTypeExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ResetTypeExtIEsPresentNothing int = iota /* No components present */
)

type ResetTypeExtIEsValue struct {
	Present int
}

type SONInformationExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       SONInformationExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	SONInformationExtIEsPresentNothing int = iota /* No components present */
)

type SONInformationExtIEsValue struct {
	Present int
}

type TargetIDExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       TargetIDExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TargetIDExtIEsPresentNothing int = iota /* No components present */
)

type TargetIDExtIEsValue struct {
	Present int
}

type UEIdentityIndexValueExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UEIdentityIndexValueExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEIdentityIndexValueExtIEsPresentNothing int = iota /* No components present */
)

type UEIdentityIndexValueExtIEsValue struct {
	Present int
}

type UENGAPIDsExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UENGAPIDsExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UENGAPIDsExtIEsPresentNothing int = iota /* No components present */
)

type UENGAPIDsExtIEsValue struct {
	Present int
}

type UEPagingIdentityExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UEPagingIdentityExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEPagingIdentityExtIEsPresentNothing int = iota /* No components present */
)

type UEPagingIdentityExtIEsValue struct {
	Present int
}

type UPTNLInformationExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UPTNLInformationExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UPTNLInformationExtIEsPresentNothing int = iota /* No components present */
)

type UPTNLInformationExtIEsValue struct {
	Present int
}

type UPTransportLayerInformationExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UPTransportLayerInformationExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UPTransportLayerInformationExtIEsPresentNothing int = iota /* No components present */
)

type UPTransportLayerInformationExtIEsValue struct {
	Present int
}

type UserLocationInformationExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UserLocationInformationExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UserLocationInformationExtIEsPresentNothing int = iota /* No components present */
)

type UserLocationInformationExtIEsValue struct {
	Present int
}

type WarningAreaListExtIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       WarningAreaListExtIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	WarningAreaListExtIEsPresentNothing int = iota /* No components present */
)

type WarningAreaListExtIEsValue struct {
	Present int
}

type PDUSessionResourceSetupRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceSetupRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupRequestIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceSetupRequestIEsPresentAMFUENGAPID
	PDUSessionResourceSetupRequestIEsPresentRANUENGAPID
	PDUSessionResourceSetupRequestIEsPresentRANPagingPriority
	PDUSessionResourceSetupRequestIEsPresentNASPDU
	PDUSessionResourceSetupRequestIEsPresentPDUSessionResourceSetupListSUReq
)

type PDUSessionResourceSetupRequestIEsValue struct {
	Present                          int
	AMFUENGAPID                      *AMFUENGAPID                      `aper:"referenceFieldValue:10"`
	RANUENGAPID                      *RANUENGAPID                      `aper:"referenceFieldValue:85"`
	RANPagingPriority                *RANPagingPriority                `aper:"referenceFieldValue:83"`
	NASPDU                           *NASPDU                           `aper:"referenceFieldValue:38"`
	PDUSessionResourceSetupListSUReq *PDUSessionResourceSetupListSUReq `aper:"referenceFieldValue:74"`
}

type PDUSessionResourceSetupResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceSetupResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupResponseIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceSetupResponseIEsPresentAMFUENGAPID
	PDUSessionResourceSetupResponseIEsPresentRANUENGAPID
	PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceSetupListSURes
	PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceFailedToSetupListSURes
	PDUSessionResourceSetupResponseIEsPresentCriticalityDiagnostics
)

type PDUSessionResourceSetupResponseIEsValue struct {
	Present                                  int
	AMFUENGAPID                              *AMFUENGAPID                              `aper:"referenceFieldValue:10"`
	RANUENGAPID                              *RANUENGAPID                              `aper:"referenceFieldValue:85"`
	PDUSessionResourceSetupListSURes         *PDUSessionResourceSetupListSURes         `aper:"referenceFieldValue:75"`
	PDUSessionResourceFailedToSetupListSURes *PDUSessionResourceFailedToSetupListSURes `aper:"referenceFieldValue:58"`
	CriticalityDiagnostics                   *CriticalityDiagnostics                   `aper:"valueExt,referenceFieldValue:19"`
}

type PDUSessionResourceReleaseCommandIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceReleaseCommandIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceReleaseCommandIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceReleaseCommandIEsPresentAMFUENGAPID
	PDUSessionResourceReleaseCommandIEsPresentRANUENGAPID
	PDUSessionResourceReleaseCommandIEsPresentRANPagingPriority
	PDUSessionResourceReleaseCommandIEsPresentNASPDU
	PDUSessionResourceReleaseCommandIEsPresentPDUSessionResourceToReleaseListRelCmd
)

type PDUSessionResourceReleaseCommandIEsValue struct {
	Present                               int
	AMFUENGAPID                           *AMFUENGAPID                           `aper:"referenceFieldValue:10"`
	RANUENGAPID                           *RANUENGAPID                           `aper:"referenceFieldValue:85"`
	RANPagingPriority                     *RANPagingPriority                     `aper:"referenceFieldValue:83"`
	NASPDU                                *NASPDU                                `aper:"referenceFieldValue:38"`
	PDUSessionResourceToReleaseListRelCmd *PDUSessionResourceToReleaseListRelCmd `aper:"referenceFieldValue:79"`
}

type PDUSessionResourceReleaseResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceReleaseResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceReleaseResponseIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceReleaseResponseIEsPresentAMFUENGAPID
	PDUSessionResourceReleaseResponseIEsPresentRANUENGAPID
	PDUSessionResourceReleaseResponseIEsPresentPDUSessionResourceReleasedListRelRes
	PDUSessionResourceReleaseResponseIEsPresentUserLocationInformation
	PDUSessionResourceReleaseResponseIEsPresentCriticalityDiagnostics
)

type PDUSessionResourceReleaseResponseIEsValue struct {
	Present                              int
	AMFUENGAPID                          *AMFUENGAPID                          `aper:"referenceFieldValue:10"`
	RANUENGAPID                          *RANUENGAPID                          `aper:"referenceFieldValue:85"`
	PDUSessionResourceReleasedListRelRes *PDUSessionResourceReleasedListRelRes `aper:"referenceFieldValue:70"`
	UserLocationInformation              *UserLocationInformation              `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
	CriticalityDiagnostics               *CriticalityDiagnostics               `aper:"valueExt,referenceFieldValue:19"`
}

type PDUSessionResourceModifyRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceModifyRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyRequestIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceModifyRequestIEsPresentAMFUENGAPID
	PDUSessionResourceModifyRequestIEsPresentRANUENGAPID
	PDUSessionResourceModifyRequestIEsPresentRANPagingPriority
	PDUSessionResourceModifyRequestIEsPresentPDUSessionResourceModifyListModReq
)

type PDUSessionResourceModifyRequestIEsValue struct {
	Present                            int
	AMFUENGAPID                        *AMFUENGAPID                        `aper:"referenceFieldValue:10"`
	RANUENGAPID                        *RANUENGAPID                        `aper:"referenceFieldValue:85"`
	RANPagingPriority                  *RANPagingPriority                  `aper:"referenceFieldValue:83"`
	PDUSessionResourceModifyListModReq *PDUSessionResourceModifyListModReq `aper:"referenceFieldValue:64"`
}

type PDUSessionResourceModifyResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceModifyResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyResponseIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceModifyResponseIEsPresentAMFUENGAPID
	PDUSessionResourceModifyResponseIEsPresentRANUENGAPID
	PDUSessionResourceModifyResponseIEsPresentPDUSessionResourceModifyListModRes
	PDUSessionResourceModifyResponseIEsPresentPDUSessionResourceFailedToModifyListModRes
	PDUSessionResourceModifyResponseIEsPresentUserLocationInformation
	PDUSessionResourceModifyResponseIEsPresentCriticalityDiagnostics
)

type PDUSessionResourceModifyResponseIEsValue struct {
	Present                                    int
	AMFUENGAPID                                *AMFUENGAPID                                `aper:"referenceFieldValue:10"`
	RANUENGAPID                                *RANUENGAPID                                `aper:"referenceFieldValue:85"`
	PDUSessionResourceModifyListModRes         *PDUSessionResourceModifyListModRes         `aper:"referenceFieldValue:65"`
	PDUSessionResourceFailedToModifyListModRes *PDUSessionResourceFailedToModifyListModRes `aper:"referenceFieldValue:54"`
	UserLocationInformation                    *UserLocationInformation                    `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
	CriticalityDiagnostics                     *CriticalityDiagnostics                     `aper:"valueExt,referenceFieldValue:19"`
}

type PDUSessionResourceNotifyIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceNotifyIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceNotifyIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceNotifyIEsPresentAMFUENGAPID
	PDUSessionResourceNotifyIEsPresentRANUENGAPID
	PDUSessionResourceNotifyIEsPresentPDUSessionResourceNotifyList
	PDUSessionResourceNotifyIEsPresentPDUSessionResourceReleasedListNot
	PDUSessionResourceNotifyIEsPresentUserLocationInformation
)

type PDUSessionResourceNotifyIEsValue struct {
	Present                           int
	AMFUENGAPID                       *AMFUENGAPID                       `aper:"referenceFieldValue:10"`
	RANUENGAPID                       *RANUENGAPID                       `aper:"referenceFieldValue:85"`
	PDUSessionResourceNotifyList      *PDUSessionResourceNotifyList      `aper:"referenceFieldValue:66"`
	PDUSessionResourceReleasedListNot *PDUSessionResourceReleasedListNot `aper:"referenceFieldValue:67"`
	UserLocationInformation           *UserLocationInformation           `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
}

type PDUSessionResourceModifyIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceModifyIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyIndicationIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceModifyIndicationIEsPresentAMFUENGAPID
	PDUSessionResourceModifyIndicationIEsPresentRANUENGAPID
	PDUSessionResourceModifyIndicationIEsPresentPDUSessionResourceModifyListModInd
)

type PDUSessionResourceModifyIndicationIEsValue struct {
	Present                            int
	AMFUENGAPID                        *AMFUENGAPID                        `aper:"referenceFieldValue:10"`
	RANUENGAPID                        *RANUENGAPID                        `aper:"referenceFieldValue:85"`
	PDUSessionResourceModifyListModInd *PDUSessionResourceModifyListModInd `aper:"referenceFieldValue:63"`
}

type PDUSessionResourceModifyConfirmIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceModifyConfirmIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyConfirmIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceModifyConfirmIEsPresentAMFUENGAPID
	PDUSessionResourceModifyConfirmIEsPresentRANUENGAPID
	PDUSessionResourceModifyConfirmIEsPresentPDUSessionResourceModifyListModCfm
	PDUSessionResourceModifyConfirmIEsPresentPDUSessionResourceFailedToModifyListModCfm
	PDUSessionResourceModifyConfirmIEsPresentCriticalityDiagnostics
)

type PDUSessionResourceModifyConfirmIEsValue struct {
	Present                                    int
	AMFUENGAPID                                *AMFUENGAPID                                `aper:"referenceFieldValue:10"`
	RANUENGAPID                                *RANUENGAPID                                `aper:"referenceFieldValue:85"`
	PDUSessionResourceModifyListModCfm         *PDUSessionResourceModifyListModCfm         `aper:"referenceFieldValue:62"`
	PDUSessionResourceFailedToModifyListModCfm *PDUSessionResourceFailedToModifyListModCfm `aper:"referenceFieldValue:131"`
	CriticalityDiagnostics                     *CriticalityDiagnostics                     `aper:"valueExt,referenceFieldValue:19"`
}

type InitialContextSetupRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       InitialContextSetupRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	InitialContextSetupRequestIEsPresentNothing int = iota /* No components present */
	InitialContextSetupRequestIEsPresentAMFUENGAPID
	InitialContextSetupRequestIEsPresentRANUENGAPID
	InitialContextSetupRequestIEsPresentOldAMF
	InitialContextSetupRequestIEsPresentUEAggregateMaximumBitRate
	InitialContextSetupRequestIEsPresentCoreNetworkAssistanceInformation
	InitialContextSetupRequestIEsPresentGUAMI
	InitialContextSetupRequestIEsPresentPDUSessionResourceSetupListCxtReq
	InitialContextSetupRequestIEsPresentAllowedNSSAI
	InitialContextSetupRequestIEsPresentUESecurityCapabilities
	InitialContextSetupRequestIEsPresentSecurityKey
	InitialContextSetupRequestIEsPresentTraceActivation
	InitialContextSetupRequestIEsPresentMobilityRestrictionList
	InitialContextSetupRequestIEsPresentUERadioCapability
	InitialContextSetupRequestIEsPresentIndexToRFSP
	InitialContextSetupRequestIEsPresentMaskedIMEISV
	InitialContextSetupRequestIEsPresentNASPDU
	InitialContextSetupRequestIEsPresentEmergencyFallbackIndicator
	InitialContextSetupRequestIEsPresentRRCInactiveTransitionReportRequest
	InitialContextSetupRequestIEsPresentUERadioCapabilityForPaging
)

type InitialContextSetupRequestIEsValue struct {
	Present                            int
	AMFUENGAPID                        *AMFUENGAPID                        `aper:"referenceFieldValue:10"`
	RANUENGAPID                        *RANUENGAPID                        `aper:"referenceFieldValue:85"`
	OldAMF                             *AMFName                            `aper:"referenceFieldValue:48"`
	UEAggregateMaximumBitRate          *UEAggregateMaximumBitRate          `aper:"valueExt,referenceFieldValue:110"`
	CoreNetworkAssistanceInformation   *CoreNetworkAssistanceInformation   `aper:"valueExt,referenceFieldValue:18"`
	GUAMI                              *GUAMI                              `aper:"valueExt,referenceFieldValue:28"`
	PDUSessionResourceSetupListCxtReq  *PDUSessionResourceSetupListCxtReq  `aper:"referenceFieldValue:71"`
	AllowedNSSAI                       *AllowedNSSAI                       `aper:"referenceFieldValue:0"`
	UESecurityCapabilities             *UESecurityCapabilities             `aper:"valueExt,referenceFieldValue:119"`
	SecurityKey                        *SecurityKey                        `aper:"referenceFieldValue:94"`
	TraceActivation                    *TraceActivation                    `aper:"valueExt,referenceFieldValue:108"`
	MobilityRestrictionList            *MobilityRestrictionList            `aper:"valueExt,referenceFieldValue:36"`
	UERadioCapability                  *UERadioCapability                  `aper:"referenceFieldValue:117"`
	IndexToRFSP                        *IndexToRFSP                        `aper:"referenceFieldValue:31"`
	MaskedIMEISV                       *MaskedIMEISV                       `aper:"referenceFieldValue:34"`
	NASPDU                             *NASPDU                             `aper:"referenceFieldValue:38"`
	EmergencyFallbackIndicator         *EmergencyFallbackIndicator         `aper:"valueExt,referenceFieldValue:24"`
	RRCInactiveTransitionReportRequest *RRCInactiveTransitionReportRequest `aper:"referenceFieldValue:91"`
	UERadioCapabilityForPaging         *UERadioCapabilityForPaging         `aper:"valueExt,referenceFieldValue:118"`
}

type InitialContextSetupResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       InitialContextSetupResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	InitialContextSetupResponseIEsPresentNothing int = iota /* No components present */
	InitialContextSetupResponseIEsPresentAMFUENGAPID
	InitialContextSetupResponseIEsPresentRANUENGAPID
	InitialContextSetupResponseIEsPresentPDUSessionResourceSetupListCxtRes
	InitialContextSetupResponseIEsPresentPDUSessionResourceFailedToSetupListCxtRes
	InitialContextSetupResponseIEsPresentCriticalityDiagnostics
)

type InitialContextSetupResponseIEsValue struct {
	Present                                   int
	AMFUENGAPID                               *AMFUENGAPID                               `aper:"referenceFieldValue:10"`
	RANUENGAPID                               *RANUENGAPID                               `aper:"referenceFieldValue:85"`
	PDUSessionResourceSetupListCxtRes         *PDUSessionResourceSetupListCxtRes         `aper:"referenceFieldValue:72"`
	PDUSessionResourceFailedToSetupListCxtRes *PDUSessionResourceFailedToSetupListCxtRes `aper:"referenceFieldValue:55"`
	CriticalityDiagnostics                    *CriticalityDiagnostics                    `aper:"valueExt,referenceFieldValue:19"`
}

type InitialContextSetupFailureIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       InitialContextSetupFailureIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	InitialContextSetupFailureIEsPresentNothing int = iota /* No components present */
	InitialContextSetupFailureIEsPresentAMFUENGAPID
	InitialContextSetupFailureIEsPresentRANUENGAPID
	InitialContextSetupFailureIEsPresentPDUSessionResourceFailedToSetupListCxtFail
	InitialContextSetupFailureIEsPresentCause
	InitialContextSetupFailureIEsPresentCriticalityDiagnostics
)

type InitialContextSetupFailureIEsValue struct {
	Present                                    int
	AMFUENGAPID                                *AMFUENGAPID                                `aper:"referenceFieldValue:10"`
	RANUENGAPID                                *RANUENGAPID                                `aper:"referenceFieldValue:85"`
	PDUSessionResourceFailedToSetupListCxtFail *PDUSessionResourceFailedToSetupListCxtFail `aper:"referenceFieldValue:132"`
	Cause                                      *Cause                                      `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	CriticalityDiagnostics                     *CriticalityDiagnostics                     `aper:"valueExt,referenceFieldValue:19"`
}

type UEContextReleaseRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UEContextReleaseRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEContextReleaseRequestIEsPresentNothing int = iota /* No components present */
	UEContextReleaseRequestIEsPresentAMFUENGAPID
	UEContextReleaseRequestIEsPresentRANUENGAPID
	UEContextReleaseRequestIEsPresentPDUSessionResourceListCxtRelReq
	UEContextReleaseRequestIEsPresentCause
)

type UEContextReleaseRequestIEsValue struct {
	Present                         int
	AMFUENGAPID                     *AMFUENGAPID                     `aper:"referenceFieldValue:10"`
	RANUENGAPID                     *RANUENGAPID                     `aper:"referenceFieldValue:85"`
	PDUSessionResourceListCxtRelReq *PDUSessionResourceListCxtRelReq `aper:"referenceFieldValue:133"`
	Cause                           *Cause                           `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
}

type UEContextReleaseCommandIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UEContextReleaseCommandIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEContextReleaseCommandIEsPresentNothing int = iota /* No components present */
	UEContextReleaseCommandIEsPresentUENGAPIDs
	UEContextReleaseCommandIEsPresentCause
)

type UEContextReleaseCommandIEsValue struct {
	Present   int
	UENGAPIDs *UENGAPIDs `aper:"referenceFieldValue:114,valueLB:0,valueUB:2"`
	Cause     *Cause     `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
}

type UEContextReleaseCompleteIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UEContextReleaseCompleteIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEContextReleaseCompleteIEsPresentNothing int = iota /* No components present */
	UEContextReleaseCompleteIEsPresentAMFUENGAPID
	UEContextReleaseCompleteIEsPresentRANUENGAPID
	UEContextReleaseCompleteIEsPresentUserLocationInformation
	UEContextReleaseCompleteIEsPresentInfoOnRecommendedCellsAndRANNodesForPaging
	UEContextReleaseCompleteIEsPresentPDUSessionResourceListCxtRelCpl
	UEContextReleaseCompleteIEsPresentCriticalityDiagnostics
)

type UEContextReleaseCompleteIEsValue struct {
	Present                                    int
	AMFUENGAPID                                *AMFUENGAPID                                `aper:"referenceFieldValue:10"`
	RANUENGAPID                                *RANUENGAPID                                `aper:"referenceFieldValue:85"`
	UserLocationInformation                    *UserLocationInformation                    `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
	InfoOnRecommendedCellsAndRANNodesForPaging *InfoOnRecommendedCellsAndRANNodesForPaging `aper:"valueExt,referenceFieldValue:32"`
	PDUSessionResourceListCxtRelCpl            *PDUSessionResourceListCxtRelCpl            `aper:"referenceFieldValue:60"`
	CriticalityDiagnostics                     *CriticalityDiagnostics                     `aper:"valueExt,referenceFieldValue:19"`
}

type UEContextModificationRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UEContextModificationRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEContextModificationRequestIEsPresentNothing int = iota /* No components present */
	UEContextModificationRequestIEsPresentAMFUENGAPID
	UEContextModificationRequestIEsPresentRANUENGAPID
	UEContextModificationRequestIEsPresentRANPagingPriority
	UEContextModificationRequestIEsPresentSecurityKey
	UEContextModificationRequestIEsPresentIndexToRFSP
	UEContextModificationRequestIEsPresentUEAggregateMaximumBitRate
	UEContextModificationRequestIEsPresentUESecurityCapabilities
	UEContextModificationRequestIEsPresentCoreNetworkAssistanceInformation
	UEContextModificationRequestIEsPresentEmergencyFallbackIndicator
	UEContextModificationRequestIEsPresentNewAMFUENGAPID
	UEContextModificationRequestIEsPresentRRCInactiveTransitionReportRequest
)

type UEContextModificationRequestIEsValue struct {
	Present                            int
	AMFUENGAPID                        *AMFUENGAPID                        `aper:"referenceFieldValue:10"`
	RANUENGAPID                        *RANUENGAPID                        `aper:"referenceFieldValue:85"`
	RANPagingPriority                  *RANPagingPriority                  `aper:"referenceFieldValue:83"`
	SecurityKey                        *SecurityKey                        `aper:"referenceFieldValue:94"`
	IndexToRFSP                        *IndexToRFSP                        `aper:"referenceFieldValue:31"`
	UEAggregateMaximumBitRate          *UEAggregateMaximumBitRate          `aper:"valueExt,referenceFieldValue:110"`
	UESecurityCapabilities             *UESecurityCapabilities             `aper:"valueExt,referenceFieldValue:119"`
	CoreNetworkAssistanceInformation   *CoreNetworkAssistanceInformation   `aper:"valueExt,referenceFieldValue:18"`
	EmergencyFallbackIndicator         *EmergencyFallbackIndicator         `aper:"valueExt,referenceFieldValue:24"`
	NewAMFUENGAPID                     *AMFUENGAPID                        `aper:"referenceFieldValue:40"`
	RRCInactiveTransitionReportRequest *RRCInactiveTransitionReportRequest `aper:"referenceFieldValue:91"`
}

type UEContextModificationResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UEContextModificationResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEContextModificationResponseIEsPresentNothing int = iota /* No components present */
	UEContextModificationResponseIEsPresentAMFUENGAPID
	UEContextModificationResponseIEsPresentRANUENGAPID
	UEContextModificationResponseIEsPresentRRCState
	UEContextModificationResponseIEsPresentUserLocationInformation
	UEContextModificationResponseIEsPresentCriticalityDiagnostics
)

type UEContextModificationResponseIEsValue struct {
	Present                 int
	AMFUENGAPID             *AMFUENGAPID             `aper:"referenceFieldValue:10"`
	RANUENGAPID             *RANUENGAPID             `aper:"referenceFieldValue:85"`
	RRCState                *RRCState                `aper:"referenceFieldValue:92"`
	UserLocationInformation *UserLocationInformation `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
	CriticalityDiagnostics  *CriticalityDiagnostics  `aper:"valueExt,referenceFieldValue:19"`
}

type UEContextModificationFailureIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UEContextModificationFailureIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UEContextModificationFailureIEsPresentNothing int = iota /* No components present */
	UEContextModificationFailureIEsPresentAMFUENGAPID
	UEContextModificationFailureIEsPresentRANUENGAPID
	UEContextModificationFailureIEsPresentCause
	UEContextModificationFailureIEsPresentCriticalityDiagnostics
)

type UEContextModificationFailureIEsValue struct {
	Present                int
	AMFUENGAPID            *AMFUENGAPID            `aper:"referenceFieldValue:10"`
	RANUENGAPID            *RANUENGAPID            `aper:"referenceFieldValue:85"`
	Cause                  *Cause                  `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type RRCInactiveTransitionReportIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       RRCInactiveTransitionReportIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RRCInactiveTransitionReportIEsPresentNothing int = iota /* No components present */
	RRCInactiveTransitionReportIEsPresentAMFUENGAPID
	RRCInactiveTransitionReportIEsPresentRANUENGAPID
	RRCInactiveTransitionReportIEsPresentRRCState
	RRCInactiveTransitionReportIEsPresentUserLocationInformation
)

type RRCInactiveTransitionReportIEsValue struct {
	Present                 int
	AMFUENGAPID             *AMFUENGAPID             `aper:"referenceFieldValue:10"`
	RANUENGAPID             *RANUENGAPID             `aper:"referenceFieldValue:85"`
	RRCState                *RRCState                `aper:"referenceFieldValue:92"`
	UserLocationInformation *UserLocationInformation `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
}

type HandoverRequiredIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverRequiredIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverRequiredIEsPresentNothing int = iota /* No components present */
	HandoverRequiredIEsPresentAMFUENGAPID
	HandoverRequiredIEsPresentRANUENGAPID
	HandoverRequiredIEsPresentHandoverType
	HandoverRequiredIEsPresentCause
	HandoverRequiredIEsPresentTargetID
	HandoverRequiredIEsPresentDirectForwardingPathAvailability
	HandoverRequiredIEsPresentPDUSessionResourceListHORqd
	HandoverRequiredIEsPresentSourceToTargetTransparentContainer
)

type HandoverRequiredIEsValue struct {
	Present                            int
	AMFUENGAPID                        *AMFUENGAPID                        `aper:"referenceFieldValue:10"`
	RANUENGAPID                        *RANUENGAPID                        `aper:"referenceFieldValue:85"`
	HandoverType                       *HandoverType                       `aper:"referenceFieldValue:29"`
	Cause                              *Cause                              `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	TargetID                           *TargetID                           `aper:"referenceFieldValue:105,valueLB:0,valueUB:2"`
	DirectForwardingPathAvailability   *DirectForwardingPathAvailability   `aper:"referenceFieldValue:22"`
	PDUSessionResourceListHORqd        *PDUSessionResourceListHORqd        `aper:"referenceFieldValue:61"`
	SourceToTargetTransparentContainer *SourceToTargetTransparentContainer `aper:"referenceFieldValue:101"`
}

type HandoverCommandIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverCommandIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverCommandIEsPresentNothing int = iota /* No components present */
	HandoverCommandIEsPresentAMFUENGAPID
	HandoverCommandIEsPresentRANUENGAPID
	HandoverCommandIEsPresentHandoverType
	HandoverCommandIEsPresentNASSecurityParametersFromNGRAN
	HandoverCommandIEsPresentPDUSessionResourceHandoverList
	HandoverCommandIEsPresentPDUSessionResourceToReleaseListHOCmd
	HandoverCommandIEsPresentTargetToSourceTransparentContainer
	HandoverCommandIEsPresentCriticalityDiagnostics
)

type HandoverCommandIEsValue struct {
	Present                              int
	AMFUENGAPID                          *AMFUENGAPID                          `aper:"referenceFieldValue:10"`
	RANUENGAPID                          *RANUENGAPID                          `aper:"referenceFieldValue:85"`
	HandoverType                         *HandoverType                         `aper:"referenceFieldValue:29"`
	NASSecurityParametersFromNGRAN       *NASSecurityParametersFromNGRAN       `aper:"referenceFieldValue:39"`
	PDUSessionResourceHandoverList       *PDUSessionResourceHandoverList       `aper:"referenceFieldValue:59"`
	PDUSessionResourceToReleaseListHOCmd *PDUSessionResourceToReleaseListHOCmd `aper:"referenceFieldValue:78"`
	TargetToSourceTransparentContainer   *TargetToSourceTransparentContainer   `aper:"referenceFieldValue:106"`
	CriticalityDiagnostics               *CriticalityDiagnostics               `aper:"valueExt,referenceFieldValue:19"`
}

type HandoverPreparationFailureIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverPreparationFailureIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverPreparationFailureIEsPresentNothing int = iota /* No components present */
	HandoverPreparationFailureIEsPresentAMFUENGAPID
	HandoverPreparationFailureIEsPresentRANUENGAPID
	HandoverPreparationFailureIEsPresentCause
	HandoverPreparationFailureIEsPresentCriticalityDiagnostics
)

type HandoverPreparationFailureIEsValue struct {
	Present                int
	AMFUENGAPID            *AMFUENGAPID            `aper:"referenceFieldValue:10"`
	RANUENGAPID            *RANUENGAPID            `aper:"referenceFieldValue:85"`
	Cause                  *Cause                  `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type HandoverRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverRequestIEsPresentNothing int = iota /* No components present */
	HandoverRequestIEsPresentAMFUENGAPID
	HandoverRequestIEsPresentHandoverType
	HandoverRequestIEsPresentCause
	HandoverRequestIEsPresentUEAggregateMaximumBitRate
	HandoverRequestIEsPresentCoreNetworkAssistanceInformation
	HandoverRequestIEsPresentUESecurityCapabilities
	HandoverRequestIEsPresentSecurityContext
	HandoverRequestIEsPresentNewSecurityContextInd
	HandoverRequestIEsPresentNASC
	HandoverRequestIEsPresentPDUSessionResourceSetupListHOReq
	HandoverRequestIEsPresentAllowedNSSAI
	HandoverRequestIEsPresentTraceActivation
	HandoverRequestIEsPresentMaskedIMEISV
	HandoverRequestIEsPresentSourceToTargetTransparentContainer
	HandoverRequestIEsPresentMobilityRestrictionList
	HandoverRequestIEsPresentLocationReportingRequestType
	HandoverRequestIEsPresentRRCInactiveTransitionReportRequest
	HandoverRequestIEsPresentGUAMI
)

type HandoverRequestIEsValue struct {
	Present                            int
	AMFUENGAPID                        *AMFUENGAPID                        `aper:"referenceFieldValue:10"`
	HandoverType                       *HandoverType                       `aper:"referenceFieldValue:29"`
	Cause                              *Cause                              `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	UEAggregateMaximumBitRate          *UEAggregateMaximumBitRate          `aper:"valueExt,referenceFieldValue:110"`
	CoreNetworkAssistanceInformation   *CoreNetworkAssistanceInformation   `aper:"valueExt,referenceFieldValue:18"`
	UESecurityCapabilities             *UESecurityCapabilities             `aper:"valueExt,referenceFieldValue:119"`
	SecurityContext                    *SecurityContext                    `aper:"valueExt,referenceFieldValue:93"`
	NewSecurityContextInd              *NewSecurityContextInd              `aper:"referenceFieldValue:41"`
	NASC                               *NASPDU                             `aper:"referenceFieldValue:37"`
	PDUSessionResourceSetupListHOReq   *PDUSessionResourceSetupListHOReq   `aper:"referenceFieldValue:73"`
	AllowedNSSAI                       *AllowedNSSAI                       `aper:"referenceFieldValue:0"`
	TraceActivation                    *TraceActivation                    `aper:"valueExt,referenceFieldValue:108"`
	MaskedIMEISV                       *MaskedIMEISV                       `aper:"referenceFieldValue:34"`
	SourceToTargetTransparentContainer *SourceToTargetTransparentContainer `aper:"referenceFieldValue:101"`
	MobilityRestrictionList            *MobilityRestrictionList            `aper:"valueExt,referenceFieldValue:36"`
	LocationReportingRequestType       *LocationReportingRequestType       `aper:"valueExt,referenceFieldValue:33"`
	RRCInactiveTransitionReportRequest *RRCInactiveTransitionReportRequest `aper:"referenceFieldValue:91"`
	GUAMI                              *GUAMI                              `aper:"valueExt,referenceFieldValue:28"`
}

type HandoverRequestAcknowledgeIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverRequestAcknowledgeIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverRequestAcknowledgeIEsPresentNothing int = iota /* No components present */
	HandoverRequestAcknowledgeIEsPresentAMFUENGAPID
	HandoverRequestAcknowledgeIEsPresentRANUENGAPID
	HandoverRequestAcknowledgeIEsPresentPDUSessionResourceAdmittedList
	HandoverRequestAcknowledgeIEsPresentPDUSessionResourceFailedToSetupListHOAck
	HandoverRequestAcknowledgeIEsPresentTargetToSourceTransparentContainer
	HandoverRequestAcknowledgeIEsPresentCriticalityDiagnostics
)

type HandoverRequestAcknowledgeIEsValue struct {
	Present                                  int
	AMFUENGAPID                              *AMFUENGAPID                              `aper:"referenceFieldValue:10"`
	RANUENGAPID                              *RANUENGAPID                              `aper:"referenceFieldValue:85"`
	PDUSessionResourceAdmittedList           *PDUSessionResourceAdmittedList           `aper:"referenceFieldValue:53"`
	PDUSessionResourceFailedToSetupListHOAck *PDUSessionResourceFailedToSetupListHOAck `aper:"referenceFieldValue:56"`
	TargetToSourceTransparentContainer       *TargetToSourceTransparentContainer       `aper:"referenceFieldValue:106"`
	CriticalityDiagnostics                   *CriticalityDiagnostics                   `aper:"valueExt,referenceFieldValue:19"`
}

type HandoverFailureIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverFailureIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverFailureIEsPresentNothing int = iota /* No components present */
	HandoverFailureIEsPresentAMFUENGAPID
	HandoverFailureIEsPresentCause
	HandoverFailureIEsPresentCriticalityDiagnostics
)

type HandoverFailureIEsValue struct {
	Present                int
	AMFUENGAPID            *AMFUENGAPID            `aper:"referenceFieldValue:10"`
	Cause                  *Cause                  `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type HandoverNotifyIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverNotifyIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverNotifyIEsPresentNothing int = iota /* No components present */
	HandoverNotifyIEsPresentAMFUENGAPID
	HandoverNotifyIEsPresentRANUENGAPID
	HandoverNotifyIEsPresentUserLocationInformation
)

type HandoverNotifyIEsValue struct {
	Present                 int
	AMFUENGAPID             *AMFUENGAPID             `aper:"referenceFieldValue:10"`
	RANUENGAPID             *RANUENGAPID             `aper:"referenceFieldValue:85"`
	UserLocationInformation *UserLocationInformation `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
}

type PathSwitchRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PathSwitchRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PathSwitchRequestIEsPresentNothing int = iota /* No components present */
	PathSwitchRequestIEsPresentRANUENGAPID
	PathSwitchRequestIEsPresentSourceAMFUENGAPID
	PathSwitchRequestIEsPresentUserLocationInformation
	PathSwitchRequestIEsPresentUESecurityCapabilities
	PathSwitchRequestIEsPresentPDUSessionResourceToBeSwitchedDLList
	PathSwitchRequestIEsPresentPDUSessionResourceFailedToSetupListPSReq
)

type PathSwitchRequestIEsValue struct {
	Present                                  int
	RANUENGAPID                              *RANUENGAPID                              `aper:"referenceFieldValue:85"`
	SourceAMFUENGAPID                        *AMFUENGAPID                              `aper:"referenceFieldValue:100"`
	UserLocationInformation                  *UserLocationInformation                  `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
	UESecurityCapabilities                   *UESecurityCapabilities                   `aper:"valueExt,referenceFieldValue:119"`
	PDUSessionResourceToBeSwitchedDLList     *PDUSessionResourceToBeSwitchedDLList     `aper:"referenceFieldValue:76"`
	PDUSessionResourceFailedToSetupListPSReq *PDUSessionResourceFailedToSetupListPSReq `aper:"referenceFieldValue:57"`
}

type PathSwitchRequestAcknowledgeIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PathSwitchRequestAcknowledgeIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PathSwitchRequestAcknowledgeIEsPresentNothing int = iota /* No components present */
	PathSwitchRequestAcknowledgeIEsPresentAMFUENGAPID
	PathSwitchRequestAcknowledgeIEsPresentRANUENGAPID
	PathSwitchRequestAcknowledgeIEsPresentUESecurityCapabilities
	PathSwitchRequestAcknowledgeIEsPresentSecurityContext
	PathSwitchRequestAcknowledgeIEsPresentNewSecurityContextInd
	PathSwitchRequestAcknowledgeIEsPresentPDUSessionResourceSwitchedList
	PathSwitchRequestAcknowledgeIEsPresentPDUSessionResourceReleasedListPSAck
	PathSwitchRequestAcknowledgeIEsPresentAllowedNSSAI
	PathSwitchRequestAcknowledgeIEsPresentCoreNetworkAssistanceInformation
	PathSwitchRequestAcknowledgeIEsPresentRRCInactiveTransitionReportRequest
	PathSwitchRequestAcknowledgeIEsPresentCriticalityDiagnostics
)

type PathSwitchRequestAcknowledgeIEsValue struct {
	Present                             int
	AMFUENGAPID                         *AMFUENGAPID                         `aper:"referenceFieldValue:10"`
	RANUENGAPID                         *RANUENGAPID                         `aper:"referenceFieldValue:85"`
	UESecurityCapabilities              *UESecurityCapabilities              `aper:"valueExt,referenceFieldValue:119"`
	SecurityContext                     *SecurityContext                     `aper:"valueExt,referenceFieldValue:93"`
	NewSecurityContextInd               *NewSecurityContextInd               `aper:"referenceFieldValue:41"`
	PDUSessionResourceSwitchedList      *PDUSessionResourceSwitchedList      `aper:"referenceFieldValue:77"`
	PDUSessionResourceReleasedListPSAck *PDUSessionResourceReleasedListPSAck `aper:"referenceFieldValue:68"`
	AllowedNSSAI                        *AllowedNSSAI                        `aper:"referenceFieldValue:0"`
	CoreNetworkAssistanceInformation    *CoreNetworkAssistanceInformation    `aper:"valueExt,referenceFieldValue:18"`
	RRCInactiveTransitionReportRequest  *RRCInactiveTransitionReportRequest  `aper:"referenceFieldValue:91"`
	CriticalityDiagnostics              *CriticalityDiagnostics              `aper:"valueExt,referenceFieldValue:19"`
}

type PathSwitchRequestFailureIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PathSwitchRequestFailureIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PathSwitchRequestFailureIEsPresentNothing int = iota /* No components present */
	PathSwitchRequestFailureIEsPresentAMFUENGAPID
	PathSwitchRequestFailureIEsPresentRANUENGAPID
	PathSwitchRequestFailureIEsPresentPDUSessionResourceReleasedListPSFail
	PathSwitchRequestFailureIEsPresentCriticalityDiagnostics
)

type PathSwitchRequestFailureIEsValue struct {
	Present                              int
	AMFUENGAPID                          *AMFUENGAPID                          `aper:"referenceFieldValue:10"`
	RANUENGAPID                          *RANUENGAPID                          `aper:"referenceFieldValue:85"`
	PDUSessionResourceReleasedListPSFail *PDUSessionResourceReleasedListPSFail `aper:"referenceFieldValue:69"`
	CriticalityDiagnostics               *CriticalityDiagnostics               `aper:"valueExt,referenceFieldValue:19"`
}

type HandoverCancelIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverCancelIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverCancelIEsPresentNothing int = iota /* No components present */
	HandoverCancelIEsPresentAMFUENGAPID
	HandoverCancelIEsPresentRANUENGAPID
	HandoverCancelIEsPresentCause
)

type HandoverCancelIEsValue struct {
	Present     int
	AMFUENGAPID *AMFUENGAPID `aper:"referenceFieldValue:10"`
	RANUENGAPID *RANUENGAPID `aper:"referenceFieldValue:85"`
	Cause       *Cause       `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
}

type HandoverCancelAcknowledgeIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       HandoverCancelAcknowledgeIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	HandoverCancelAcknowledgeIEsPresentNothing int = iota /* No components present */
	HandoverCancelAcknowledgeIEsPresentAMFUENGAPID
	HandoverCancelAcknowledgeIEsPresentRANUENGAPID
	HandoverCancelAcknowledgeIEsPresentCriticalityDiagnostics
)

type HandoverCancelAcknowledgeIEsValue struct {
	Present                int
	AMFUENGAPID            *AMFUENGAPID            `aper:"referenceFieldValue:10"`
	RANUENGAPID            *RANUENGAPID            `aper:"referenceFieldValue:85"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type UplinkRANStatusTransferIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UplinkRANStatusTransferIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UplinkRANStatusTransferIEsPresentNothing int = iota /* No components present */
	UplinkRANStatusTransferIEsPresentAMFUENGAPID
	UplinkRANStatusTransferIEsPresentRANUENGAPID
	UplinkRANStatusTransferIEsPresentRANStatusTransferTransparentContainer
)

type UplinkRANStatusTransferIEsValue struct {
	Present                               int
	AMFUENGAPID                           *AMFUENGAPID                           `aper:"referenceFieldValue:10"`
	RANUENGAPID                           *RANUENGAPID                           `aper:"referenceFieldValue:85"`
	RANStatusTransferTransparentContainer *RANStatusTransferTransparentContainer `aper:"valueExt,referenceFieldValue:84"`
}

type DownlinkRANStatusTransferIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       DownlinkRANStatusTransferIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DownlinkRANStatusTransferIEsPresentNothing int = iota /* No components present */
	DownlinkRANStatusTransferIEsPresentAMFUENGAPID
	DownlinkRANStatusTransferIEsPresentRANUENGAPID
	DownlinkRANStatusTransferIEsPresentRANStatusTransferTransparentContainer
)

type DownlinkRANStatusTransferIEsValue struct {
	Present                               int
	AMFUENGAPID                           *AMFUENGAPID                           `aper:"referenceFieldValue:10"`
	RANUENGAPID                           *RANUENGAPID                           `aper:"referenceFieldValue:85"`
	RANStatusTransferTransparentContainer *RANStatusTransferTransparentContainer `aper:"valueExt,referenceFieldValue:84"`
}

type PagingIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PagingIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PagingIEsPresentNothing int = iota /* No components present */
	PagingIEsPresentUEPagingIdentity
	PagingIEsPresentPagingDRX
	PagingIEsPresentTAIListForPaging
	PagingIEsPresentPagingPriority
	PagingIEsPresentUERadioCapabilityForPaging
	PagingIEsPresentPagingOrigin
	PagingIEsPresentAssistanceDataForPaging
)

type PagingIEsValue struct {
	Present                    int
	UEPagingIdentity           *UEPagingIdentity           `aper:"referenceFieldValue:115,valueLB:0,valueUB:1"`
	PagingDRX                  *PagingDRX                  `aper:"referenceFieldValue:50"`
	TAIListForPaging           *TAIListForPaging           `aper:"referenceFieldValue:103"`
	PagingPriority             *PagingPriority             `aper:"referenceFieldValue:52"`
	UERadioCapabilityForPaging *UERadioCapabilityForPaging `aper:"valueExt,referenceFieldValue:118"`
	PagingOrigin               *PagingOrigin               `aper:"referenceFieldValue:51"`
	AssistanceDataForPaging    *AssistanceDataForPaging    `aper:"valueExt,referenceFieldValue:11"`
}

type InitialUEMessageIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       InitialUEMessageIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	InitialUEMessageIEsPresentNothing int = iota /* No components present */
	InitialUEMessageIEsPresentRANUENGAPID
	InitialUEMessageIEsPresentNASPDU
	InitialUEMessageIEsPresentUserLocationInformation
	InitialUEMessageIEsPresentRRCEstablishmentCause
	InitialUEMessageIEsPresentFiveGSTMSI
	InitialUEMessageIEsPresentAMFSetID
	InitialUEMessageIEsPresentUEContextRequest
	InitialUEMessageIEsPresentAllowedNSSAI
)

type InitialUEMessageIEsValue struct {
	Present                 int
	RANUENGAPID             *RANUENGAPID             `aper:"referenceFieldValue:85"`
	NASPDU                  *NASPDU                  `aper:"referenceFieldValue:38"`
	UserLocationInformation *UserLocationInformation `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
	RRCEstablishmentCause   *RRCEstablishmentCause   `aper:"referenceFieldValue:90"`
	FiveGSTMSI              *FiveGSTMSI              `aper:"valueExt,referenceFieldValue:26"`
	AMFSetID                *AMFSetID                `aper:"referenceFieldValue:3"`
	UEContextRequest        *UEContextRequest        `aper:"referenceFieldValue:112"`
	AllowedNSSAI            *AllowedNSSAI            `aper:"referenceFieldValue:0"`
}

type DownlinkNASTransportIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       DownlinkNASTransportIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DownlinkNASTransportIEsPresentNothing int = iota /* No components present */
	DownlinkNASTransportIEsPresentAMFUENGAPID
	DownlinkNASTransportIEsPresentRANUENGAPID
	DownlinkNASTransportIEsPresentOldAMF
	DownlinkNASTransportIEsPresentRANPagingPriority
	DownlinkNASTransportIEsPresentNASPDU
	DownlinkNASTransportIEsPresentMobilityRestrictionList
	DownlinkNASTransportIEsPresentIndexToRFSP
	DownlinkNASTransportIEsPresentUEAggregateMaximumBitRate
	DownlinkNASTransportIEsPresentAllowedNSSAI
)

type DownlinkNASTransportIEsValue struct {
	Present                   int
	AMFUENGAPID               *AMFUENGAPID               `aper:"referenceFieldValue:10"`
	RANUENGAPID               *RANUENGAPID               `aper:"referenceFieldValue:85"`
	OldAMF                    *AMFName                   `aper:"referenceFieldValue:48"`
	RANPagingPriority         *RANPagingPriority         `aper:"referenceFieldValue:83"`
	NASPDU                    *NASPDU                    `aper:"referenceFieldValue:38"`
	MobilityRestrictionList   *MobilityRestrictionList   `aper:"valueExt,referenceFieldValue:36"`
	IndexToRFSP               *IndexToRFSP               `aper:"referenceFieldValue:31"`
	UEAggregateMaximumBitRate *UEAggregateMaximumBitRate `aper:"valueExt,referenceFieldValue:110"`
	AllowedNSSAI              *AllowedNSSAI              `aper:"referenceFieldValue:0"`
}

type UplinkNASTransportIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UplinkNASTransportIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UplinkNASTransportIEsPresentNothing int = iota /* No components present */
	UplinkNASTransportIEsPresentAMFUENGAPID
	UplinkNASTransportIEsPresentRANUENGAPID
	UplinkNASTransportIEsPresentNASPDU
	UplinkNASTransportIEsPresentUserLocationInformation
)

type UplinkNASTransportIEsValue struct {
	Present                 int
	AMFUENGAPID             *AMFUENGAPID             `aper:"referenceFieldValue:10"`
	RANUENGAPID             *RANUENGAPID             `aper:"referenceFieldValue:85"`
	NASPDU                  *NASPDU                  `aper:"referenceFieldValue:38"`
	UserLocationInformation *UserLocationInformation `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
}

type NASNonDeliveryIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       NASNonDeliveryIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NASNonDeliveryIndicationIEsPresentNothing int = iota /* No components present */
	NASNonDeliveryIndicationIEsPresentAMFUENGAPID
	NASNonDeliveryIndicationIEsPresentRANUENGAPID
	NASNonDeliveryIndicationIEsPresentNASPDU
	NASNonDeliveryIndicationIEsPresentCause
)

type NASNonDeliveryIndicationIEsValue struct {
	Present     int
	AMFUENGAPID *AMFUENGAPID `aper:"referenceFieldValue:10"`
	RANUENGAPID *RANUENGAPID `aper:"referenceFieldValue:85"`
	NASPDU      *NASPDU      `aper:"referenceFieldValue:38"`
	Cause       *Cause       `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
}

type RerouteNASRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       RerouteNASRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RerouteNASRequestIEsPresentNothing int = iota /* No components present */
	RerouteNASRequestIEsPresentRANUENGAPID
	RerouteNASRequestIEsPresentAMFUENGAPID
	RerouteNASRequestIEsPresentNGAPMessage
	RerouteNASRequestIEsPresentAMFSetID
	RerouteNASRequestIEsPresentAllowedNSSAI
)

type RerouteNASRequestIEsValue struct {
	Present      int
	RANUENGAPID  *RANUENGAPID      `aper:"referenceFieldValue:85"`
	AMFUENGAPID  *AMFUENGAPID      `aper:"referenceFieldValue:10"`
	NGAPMessage  *aper.OctetString `aper:"referenceFieldValue:42"`
	AMFSetID     *AMFSetID         `aper:"referenceFieldValue:3"`
	AllowedNSSAI *AllowedNSSAI     `aper:"referenceFieldValue:0"`
}

type NGSetupRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       NGSetupRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NGSetupRequestIEsPresentNothing int = iota /* No components present */
	NGSetupRequestIEsPresentGlobalRANNodeID
	NGSetupRequestIEsPresentRANNodeName
	NGSetupRequestIEsPresentSupportedTAList
	NGSetupRequestIEsPresentDefaultPagingDRX
)

type NGSetupRequestIEsValue struct {
	Present          int
	GlobalRANNodeID  *GlobalRANNodeID `aper:"referenceFieldValue:27,valueLB:0,valueUB:3"`
	RANNodeName      *RANNodeName     `aper:"referenceFieldValue:82"`
	SupportedTAList  *SupportedTAList `aper:"referenceFieldValue:102"`
	DefaultPagingDRX *PagingDRX       `aper:"referenceFieldValue:21"`
}

type NGSetupResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       NGSetupResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NGSetupResponseIEsPresentNothing int = iota /* No components present */
	NGSetupResponseIEsPresentAMFName
	NGSetupResponseIEsPresentServedGUAMIList
	NGSetupResponseIEsPresentRelativeAMFCapacity
	NGSetupResponseIEsPresentPLMNSupportList
	NGSetupResponseIEsPresentCriticalityDiagnostics
)

type NGSetupResponseIEsValue struct {
	Present                int
	AMFName                *AMFName                `aper:"referenceFieldValue:1"`
	ServedGUAMIList        *ServedGUAMIList        `aper:"referenceFieldValue:96"`
	RelativeAMFCapacity    *RelativeAMFCapacity    `aper:"referenceFieldValue:86"`
	PLMNSupportList        *PLMNSupportList        `aper:"referenceFieldValue:80"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type NGSetupFailureIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       NGSetupFailureIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NGSetupFailureIEsPresentNothing int = iota /* No components present */
	NGSetupFailureIEsPresentCause
	NGSetupFailureIEsPresentTimeToWait
	NGSetupFailureIEsPresentCriticalityDiagnostics
)

type NGSetupFailureIEsValue struct {
	Present                int
	Cause                  *Cause                  `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	TimeToWait             *TimeToWait             `aper:"referenceFieldValue:107"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type RANConfigurationUpdateIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       RANConfigurationUpdateIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RANConfigurationUpdateIEsPresentNothing int = iota /* No components present */
	RANConfigurationUpdateIEsPresentRANNodeName
	RANConfigurationUpdateIEsPresentSupportedTAList
	RANConfigurationUpdateIEsPresentDefaultPagingDRX
)

type RANConfigurationUpdateIEsValue struct {
	Present          int
	RANNodeName      *RANNodeName     `aper:"referenceFieldValue:82"`
	SupportedTAList  *SupportedTAList `aper:"referenceFieldValue:102"`
	DefaultPagingDRX *PagingDRX       `aper:"referenceFieldValue:21"`
}

type RANConfigurationUpdateAcknowledgeIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       RANConfigurationUpdateAcknowledgeIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RANConfigurationUpdateAcknowledgeIEsPresentNothing int = iota /* No components present */
	RANConfigurationUpdateAcknowledgeIEsPresentCriticalityDiagnostics
)

type RANConfigurationUpdateAcknowledgeIEsValue struct {
	Present                int
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type RANConfigurationUpdateFailureIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       RANConfigurationUpdateFailureIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	RANConfigurationUpdateFailureIEsPresentNothing int = iota /* No components present */
	RANConfigurationUpdateFailureIEsPresentCause
	RANConfigurationUpdateFailureIEsPresentTimeToWait
	RANConfigurationUpdateFailureIEsPresentCriticalityDiagnostics
)

type RANConfigurationUpdateFailureIEsValue struct {
	Present                int
	Cause                  *Cause                  `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	TimeToWait             *TimeToWait             `aper:"referenceFieldValue:107"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type AMFConfigurationUpdateIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       AMFConfigurationUpdateIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFConfigurationUpdateIEsPresentNothing int = iota /* No components present */
	AMFConfigurationUpdateIEsPresentAMFName
	AMFConfigurationUpdateIEsPresentServedGUAMIList
	AMFConfigurationUpdateIEsPresentRelativeAMFCapacity
	AMFConfigurationUpdateIEsPresentPLMNSupportList
	AMFConfigurationUpdateIEsPresentAMFTNLAssociationToAddList
	AMFConfigurationUpdateIEsPresentAMFTNLAssociationToRemoveList
	AMFConfigurationUpdateIEsPresentAMFTNLAssociationToUpdateList
)

type AMFConfigurationUpdateIEsValue struct {
	Present                       int
	AMFName                       *AMFName                       `aper:"referenceFieldValue:1"`
	ServedGUAMIList               *ServedGUAMIList               `aper:"referenceFieldValue:96"`
	RelativeAMFCapacity           *RelativeAMFCapacity           `aper:"referenceFieldValue:86"`
	PLMNSupportList               *PLMNSupportList               `aper:"referenceFieldValue:80"`
	AMFTNLAssociationToAddList    *AMFTNLAssociationToAddList    `aper:"referenceFieldValue:6"`
	AMFTNLAssociationToRemoveList *AMFTNLAssociationToRemoveList `aper:"referenceFieldValue:7"`
	AMFTNLAssociationToUpdateList *AMFTNLAssociationToUpdateList `aper:"referenceFieldValue:8"`
}

type AMFConfigurationUpdateAcknowledgeIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       AMFConfigurationUpdateAcknowledgeIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFConfigurationUpdateAcknowledgeIEsPresentNothing int = iota /* No components present */
	AMFConfigurationUpdateAcknowledgeIEsPresentAMFTNLAssociationSetupList
	AMFConfigurationUpdateAcknowledgeIEsPresentAMFTNLAssociationFailedToSetupList
	AMFConfigurationUpdateAcknowledgeIEsPresentCriticalityDiagnostics
)

type AMFConfigurationUpdateAcknowledgeIEsValue struct {
	Present                            int
	AMFTNLAssociationSetupList         *AMFTNLAssociationSetupList `aper:"referenceFieldValue:5"`
	AMFTNLAssociationFailedToSetupList *TNLAssociationList         `aper:"referenceFieldValue:4"`
	CriticalityDiagnostics             *CriticalityDiagnostics     `aper:"valueExt,referenceFieldValue:19"`
}

type AMFConfigurationUpdateFailureIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       AMFConfigurationUpdateFailureIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFConfigurationUpdateFailureIEsPresentNothing int = iota /* No components present */
	AMFConfigurationUpdateFailureIEsPresentCause
	AMFConfigurationUpdateFailureIEsPresentTimeToWait
	AMFConfigurationUpdateFailureIEsPresentCriticalityDiagnostics
)

type AMFConfigurationUpdateFailureIEsValue struct {
	Present                int
	Cause                  *Cause                  `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	TimeToWait             *TimeToWait             `aper:"referenceFieldValue:107"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type AMFStatusIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       AMFStatusIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	AMFStatusIndicationIEsPresentNothing int = iota /* No components present */
	AMFStatusIndicationIEsPresentUnavailableGUAMIList
)

type AMFStatusIndicationIEsValue struct {
	Present              int
	UnavailableGUAMIList *UnavailableGUAMIList `aper:"referenceFieldValue:120"`
}

type NGResetIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       NGResetIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NGResetIEsPresentNothing int = iota /* No components present */
	NGResetIEsPresentCause
	NGResetIEsPresentResetType
)

type NGResetIEsValue struct {
	Present   int
	Cause     *Cause     `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	ResetType *ResetType `aper:"referenceFieldValue:88,valueLB:0,valueUB:2"`
}

type NGResetAcknowledgeIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       NGResetAcknowledgeIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	NGResetAcknowledgeIEsPresentNothing int = iota /* No components present */
	NGResetAcknowledgeIEsPresentUEAssociatedLogicalNGConnectionList
	NGResetAcknowledgeIEsPresentCriticalityDiagnostics
)

type NGResetAcknowledgeIEsValue struct {
	Present                             int
	UEAssociatedLogicalNGConnectionList *UEAssociatedLogicalNGConnectionList `aper:"referenceFieldValue:111"`
	CriticalityDiagnostics              *CriticalityDiagnostics              `aper:"valueExt,referenceFieldValue:19"`
}

type ErrorIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       ErrorIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	ErrorIndicationIEsPresentNothing int = iota /* No components present */
	ErrorIndicationIEsPresentAMFUENGAPID
	ErrorIndicationIEsPresentRANUENGAPID
	ErrorIndicationIEsPresentCause
	ErrorIndicationIEsPresentCriticalityDiagnostics
)

type ErrorIndicationIEsValue struct {
	Present                int
	AMFUENGAPID            *AMFUENGAPID            `aper:"referenceFieldValue:10"`
	RANUENGAPID            *RANUENGAPID            `aper:"referenceFieldValue:85"`
	Cause                  *Cause                  `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
	CriticalityDiagnostics *CriticalityDiagnostics `aper:"valueExt,referenceFieldValue:19"`
}

type OverloadStartIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       OverloadStartIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	OverloadStartIEsPresentNothing int = iota /* No components present */
	OverloadStartIEsPresentAMFOverloadResponse
	OverloadStartIEsPresentAMFTrafficLoadReductionIndication
	OverloadStartIEsPresentOverloadStartNSSAIList
)

type OverloadStartIEsValue struct {
	Present                           int
	AMFOverloadResponse               *OverloadResponse               `aper:"referenceFieldValue:2,valueLB:0,valueUB:1"`
	AMFTrafficLoadReductionIndication *TrafficLoadReductionIndication `aper:"referenceFieldValue:9"`
	OverloadStartNSSAIList            *OverloadStartNSSAIList         `aper:"referenceFieldValue:49"`
}

type OverloadStopIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       OverloadStopIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	OverloadStopIEsPresentNothing int = iota /* No components present */
)

type OverloadStopIEsValue struct {
	Present int
}

type UplinkRANConfigurationTransferIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UplinkRANConfigurationTransferIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UplinkRANConfigurationTransferIEsPresentNothing int = iota /* No components present */
	UplinkRANConfigurationTransferIEsPresentSONConfigurationTransferUL
)

type UplinkRANConfigurationTransferIEsValue struct {
	Present                    int
	SONConfigurationTransferUL *SONConfigurationTransfer `aper:"valueExt,referenceFieldValue:99"`
}

type DownlinkRANConfigurationTransferIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       DownlinkRANConfigurationTransferIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DownlinkRANConfigurationTransferIEsPresentNothing int = iota /* No components present */
	DownlinkRANConfigurationTransferIEsPresentSONConfigurationTransferDL
)

type DownlinkRANConfigurationTransferIEsValue struct {
	Present                    int
	SONConfigurationTransferDL *SONConfigurationTransfer `aper:"valueExt,referenceFieldValue:98"`
}

type WriteReplaceWarningRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       WriteReplaceWarningRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	WriteReplaceWarningRequestIEsPresentNothing int = iota /* No components present */
	WriteReplaceWarningRequestIEsPresentMessageIdentifier
	WriteReplaceWarningRequestIEsPresentSerialNumber
	WriteReplaceWarningRequestIEsPresentWarningAreaList
	WriteReplaceWarningRequestIEsPresentRepetitionPeriod
	WriteReplaceWarningRequestIEsPresentNumberOfBroadcastsRequested
	WriteReplaceWarningRequestIEsPresentWarningType
	WriteReplaceWarningRequestIEsPresentWarningSecurityInfo
	WriteReplaceWarningRequestIEsPresentDataCodingScheme
	WriteReplaceWarningRequestIEsPresentWarningMessageContents
	WriteReplaceWarningRequestIEsPresentConcurrentWarningMessageInd
	WriteReplaceWarningRequestIEsPresentWarningAreaCoordinates
)

type WriteReplaceWarningRequestIEsValue struct {
	Present                     int
	MessageIdentifier           *MessageIdentifier           `aper:"referenceFieldValue:35"`
	SerialNumber                *SerialNumber                `aper:"referenceFieldValue:95"`
	WarningAreaList             *WarningAreaList             `aper:"referenceFieldValue:122,valueLB:0,valueUB:4"`
	RepetitionPeriod            *RepetitionPeriod            `aper:"referenceFieldValue:87"`
	NumberOfBroadcastsRequested *NumberOfBroadcastsRequested `aper:"referenceFieldValue:47"`
	WarningType                 *WarningType                 `aper:"referenceFieldValue:125"`
	WarningSecurityInfo         *WarningSecurityInfo         `aper:"referenceFieldValue:124"`
	DataCodingScheme            *DataCodingScheme            `aper:"referenceFieldValue:20"`
	WarningMessageContents      *WarningMessageContents      `aper:"referenceFieldValue:123"`
	ConcurrentWarningMessageInd *ConcurrentWarningMessageInd `aper:"referenceFieldValue:17"`
	WarningAreaCoordinates      *WarningAreaCoordinates      `aper:"referenceFieldValue:141"`
}

type WriteReplaceWarningResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       WriteReplaceWarningResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	WriteReplaceWarningResponseIEsPresentNothing int = iota /* No components present */
	WriteReplaceWarningResponseIEsPresentMessageIdentifier
	WriteReplaceWarningResponseIEsPresentSerialNumber
	WriteReplaceWarningResponseIEsPresentBroadcastCompletedAreaList
	WriteReplaceWarningResponseIEsPresentCriticalityDiagnostics
)

type WriteReplaceWarningResponseIEsValue struct {
	Present                    int
	MessageIdentifier          *MessageIdentifier          `aper:"referenceFieldValue:35"`
	SerialNumber               *SerialNumber               `aper:"referenceFieldValue:95"`
	BroadcastCompletedAreaList *BroadcastCompletedAreaList `aper:"referenceFieldValue:13,valueLB:0,valueUB:6"`
	CriticalityDiagnostics     *CriticalityDiagnostics     `aper:"valueExt,referenceFieldValue:19"`
}

type PWSCancelRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PWSCancelRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PWSCancelRequestIEsPresentNothing int = iota /* No components present */
	PWSCancelRequestIEsPresentMessageIdentifier
	PWSCancelRequestIEsPresentSerialNumber
	PWSCancelRequestIEsPresentWarningAreaList
	PWSCancelRequestIEsPresentCancelAllWarningMessages
)

type PWSCancelRequestIEsValue struct {
	Present                  int
	MessageIdentifier        *MessageIdentifier        `aper:"referenceFieldValue:35"`
	SerialNumber             *SerialNumber             `aper:"referenceFieldValue:95"`
	WarningAreaList          *WarningAreaList          `aper:"referenceFieldValue:122,valueLB:0,valueUB:4"`
	CancelAllWarningMessages *CancelAllWarningMessages `aper:"referenceFieldValue:14"`
}

type PWSCancelResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PWSCancelResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PWSCancelResponseIEsPresentNothing int = iota /* No components present */
	PWSCancelResponseIEsPresentMessageIdentifier
	PWSCancelResponseIEsPresentSerialNumber
	PWSCancelResponseIEsPresentBroadcastCancelledAreaList
	PWSCancelResponseIEsPresentCriticalityDiagnostics
)

type PWSCancelResponseIEsValue struct {
	Present                    int
	MessageIdentifier          *MessageIdentifier          `aper:"referenceFieldValue:35"`
	SerialNumber               *SerialNumber               `aper:"referenceFieldValue:95"`
	BroadcastCancelledAreaList *BroadcastCancelledAreaList `aper:"referenceFieldValue:12,valueLB:0,valueUB:6"`
	CriticalityDiagnostics     *CriticalityDiagnostics     `aper:"valueExt,referenceFieldValue:19"`
}

type PWSRestartIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PWSRestartIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PWSRestartIndicationIEsPresentNothing int = iota /* No components present */
	PWSRestartIndicationIEsPresentCellIDListForRestart
	PWSRestartIndicationIEsPresentGlobalRANNodeID
	PWSRestartIndicationIEsPresentTAIListForRestart
	PWSRestartIndicationIEsPresentEmergencyAreaIDListForRestart
)

type PWSRestartIndicationIEsValue struct {
	Present                       int
	CellIDListForRestart          *CellIDListForRestart          `aper:"referenceFieldValue:16,valueLB:0,valueUB:2"`
	GlobalRANNodeID               *GlobalRANNodeID               `aper:"referenceFieldValue:27,valueLB:0,valueUB:3"`
	TAIListForRestart             *TAIListForRestart             `aper:"referenceFieldValue:104"`
	EmergencyAreaIDListForRestart *EmergencyAreaIDListForRestart `aper:"referenceFieldValue:23"`
}

type PWSFailureIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PWSFailureIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PWSFailureIndicationIEsPresentNothing int = iota /* No components present */
	PWSFailureIndicationIEsPresentPWSFailedCellIDList
	PWSFailureIndicationIEsPresentGlobalRANNodeID
)

type PWSFailureIndicationIEsValue struct {
	Present             int
	PWSFailedCellIDList *PWSFailedCellIDList `aper:"referenceFieldValue:81,valueLB:0,valueUB:2"`
	GlobalRANNodeID     *GlobalRANNodeID     `aper:"referenceFieldValue:27,valueLB:0,valueUB:3"`
}

type DownlinkUEAssociatedNRPPaTransportIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       DownlinkUEAssociatedNRPPaTransportIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DownlinkUEAssociatedNRPPaTransportIEsPresentNothing int = iota /* No components present */
	DownlinkUEAssociatedNRPPaTransportIEsPresentAMFUENGAPID
	DownlinkUEAssociatedNRPPaTransportIEsPresentRANUENGAPID
	DownlinkUEAssociatedNRPPaTransportIEsPresentRoutingID
	DownlinkUEAssociatedNRPPaTransportIEsPresentNRPPaPDU
)

type DownlinkUEAssociatedNRPPaTransportIEsValue struct {
	Present     int
	AMFUENGAPID *AMFUENGAPID `aper:"referenceFieldValue:10"`
	RANUENGAPID *RANUENGAPID `aper:"referenceFieldValue:85"`
	RoutingID   *RoutingID   `aper:"referenceFieldValue:89"`
	NRPPaPDU    *NRPPaPDU    `aper:"referenceFieldValue:46"`
}

type UplinkUEAssociatedNRPPaTransportIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UplinkUEAssociatedNRPPaTransportIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UplinkUEAssociatedNRPPaTransportIEsPresentNothing int = iota /* No components present */
	UplinkUEAssociatedNRPPaTransportIEsPresentAMFUENGAPID
	UplinkUEAssociatedNRPPaTransportIEsPresentRANUENGAPID
	UplinkUEAssociatedNRPPaTransportIEsPresentRoutingID
	UplinkUEAssociatedNRPPaTransportIEsPresentNRPPaPDU
)

type UplinkUEAssociatedNRPPaTransportIEsValue struct {
	Present     int
	AMFUENGAPID *AMFUENGAPID `aper:"referenceFieldValue:10"`
	RANUENGAPID *RANUENGAPID `aper:"referenceFieldValue:85"`
	RoutingID   *RoutingID   `aper:"referenceFieldValue:89"`
	NRPPaPDU    *NRPPaPDU    `aper:"referenceFieldValue:46"`
}

type DownlinkNonUEAssociatedNRPPaTransportIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       DownlinkNonUEAssociatedNRPPaTransportIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DownlinkNonUEAssociatedNRPPaTransportIEsPresentNothing int = iota /* No components present */
	DownlinkNonUEAssociatedNRPPaTransportIEsPresentRoutingID
	DownlinkNonUEAssociatedNRPPaTransportIEsPresentNRPPaPDU
)

type DownlinkNonUEAssociatedNRPPaTransportIEsValue struct {
	Present   int
	RoutingID *RoutingID `aper:"referenceFieldValue:89"`
	NRPPaPDU  *NRPPaPDU  `aper:"referenceFieldValue:46"`
}

type UplinkNonUEAssociatedNRPPaTransportIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UplinkNonUEAssociatedNRPPaTransportIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UplinkNonUEAssociatedNRPPaTransportIEsPresentNothing int = iota /* No components present */
	UplinkNonUEAssociatedNRPPaTransportIEsPresentRoutingID
	UplinkNonUEAssociatedNRPPaTransportIEsPresentNRPPaPDU
)

type UplinkNonUEAssociatedNRPPaTransportIEsValue struct {
	Present   int
	RoutingID *RoutingID `aper:"referenceFieldValue:89"`
	NRPPaPDU  *NRPPaPDU  `aper:"referenceFieldValue:46"`
}

type TraceStartIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       TraceStartIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TraceStartIEsPresentNothing int = iota /* No components present */
	TraceStartIEsPresentAMFUENGAPID
	TraceStartIEsPresentRANUENGAPID
	TraceStartIEsPresentTraceActivation
)

type TraceStartIEsValue struct {
	Present         int
	AMFUENGAPID     *AMFUENGAPID     `aper:"referenceFieldValue:10"`
	RANUENGAPID     *RANUENGAPID     `aper:"referenceFieldValue:85"`
	TraceActivation *TraceActivation `aper:"valueExt,referenceFieldValue:108"`
}

type TraceFailureIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       TraceFailureIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	TraceFailureIndicationIEsPresentNothing int = iota /* No components present */
	TraceFailureIndicationIEsPresentAMFUENGAPID
	TraceFailureIndicationIEsPresentRANUENGAPID
	TraceFailureIndicationIEsPresentNGRANTraceID
	TraceFailureIndicationIEsPresentCause
)

type TraceFailureIndicationIEsValue struct {
	Present      int
	AMFUENGAPID  *AMFUENGAPID  `aper:"referenceFieldValue:10"`
	RANUENGAPID  *RANUENGAPID  `aper:"referenceFieldValue:85"`
	NGRANTraceID *NGRANTraceID `aper:"referenceFieldValue:44"`
	Cause        *Cause        `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
}

type DeactivateTraceIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       DeactivateTraceIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	DeactivateTraceIEsPresentNothing int = iota /* No components present */
	DeactivateTraceIEsPresentAMFUENGAPID
	DeactivateTraceIEsPresentRANUENGAPID
	DeactivateTraceIEsPresentNGRANTraceID
)

type DeactivateTraceIEsValue struct {
	Present      int
	AMFUENGAPID  *AMFUENGAPID  `aper:"referenceFieldValue:10"`
	RANUENGAPID  *RANUENGAPID  `aper:"referenceFieldValue:85"`
	NGRANTraceID *NGRANTraceID `aper:"referenceFieldValue:44"`
}

type CellTrafficTraceIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       CellTrafficTraceIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	CellTrafficTraceIEsPresentNothing int = iota /* No components present */
	CellTrafficTraceIEsPresentAMFUENGAPID
	CellTrafficTraceIEsPresentRANUENGAPID
	CellTrafficTraceIEsPresentNGRANTraceID
	CellTrafficTraceIEsPresentNGRANCGI
	CellTrafficTraceIEsPresentTraceCollectionEntityIPAddress
)

type CellTrafficTraceIEsValue struct {
	Present                        int
	AMFUENGAPID                    *AMFUENGAPID           `aper:"referenceFieldValue:10"`
	RANUENGAPID                    *RANUENGAPID           `aper:"referenceFieldValue:85"`
	NGRANTraceID                   *NGRANTraceID          `aper:"referenceFieldValue:44"`
	NGRANCGI                       *NGRANCGI              `aper:"referenceFieldValue:43,valueLB:0,valueUB:2"`
	TraceCollectionEntityIPAddress *TransportLayerAddress `aper:"referenceFieldValue:109"`
}

type LocationReportingControlIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       LocationReportingControlIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	LocationReportingControlIEsPresentNothing int = iota /* No components present */
	LocationReportingControlIEsPresentAMFUENGAPID
	LocationReportingControlIEsPresentRANUENGAPID
	LocationReportingControlIEsPresentLocationReportingRequestType
)

type LocationReportingControlIEsValue struct {
	Present                      int
	AMFUENGAPID                  *AMFUENGAPID                  `aper:"referenceFieldValue:10"`
	RANUENGAPID                  *RANUENGAPID                  `aper:"referenceFieldValue:85"`
	LocationReportingRequestType *LocationReportingRequestType `aper:"valueExt,referenceFieldValue:33"`
}

type LocationReportingFailureIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       LocationReportingFailureIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	LocationReportingFailureIndicationIEsPresentNothing int = iota /* No components present */
	LocationReportingFailureIndicationIEsPresentAMFUENGAPID
	LocationReportingFailureIndicationIEsPresentRANUENGAPID
	LocationReportingFailureIndicationIEsPresentCause
)

type LocationReportingFailureIndicationIEsValue struct {
	Present     int
	AMFUENGAPID *AMFUENGAPID `aper:"referenceFieldValue:10"`
	RANUENGAPID *RANUENGAPID `aper:"referenceFieldValue:85"`
	Cause       *Cause       `aper:"referenceFieldValue:15,valueLB:0,valueUB:5"`
}

type LocationReportIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       LocationReportIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	LocationReportIEsPresentNothing int = iota /* No components present */
	LocationReportIEsPresentAMFUENGAPID
	LocationReportIEsPresentRANUENGAPID
	LocationReportIEsPresentUserLocationInformation
	LocationReportIEsPresentUEPresenceInAreaOfInterestList
	LocationReportIEsPresentLocationReportingRequestType
)

type LocationReportIEsValue struct {
	Present                        int
	AMFUENGAPID                    *AMFUENGAPID                    `aper:"referenceFieldValue:10"`
	RANUENGAPID                    *RANUENGAPID                    `aper:"referenceFieldValue:85"`
	UserLocationInformation        *UserLocationInformation        `aper:"referenceFieldValue:121,valueLB:0,valueUB:3"`
	UEPresenceInAreaOfInterestList *UEPresenceInAreaOfInterestList `aper:"referenceFieldValue:116"`
	LocationReportingRequestType   *LocationReportingRequestType   `aper:"valueExt,referenceFieldValue:33"`
}

type UETNLABindingReleaseRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UETNLABindingReleaseRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UETNLABindingReleaseRequestIEsPresentNothing int = iota /* No components present */
	UETNLABindingReleaseRequestIEsPresentAMFUENGAPID
	UETNLABindingReleaseRequestIEsPresentRANUENGAPID
)

type UETNLABindingReleaseRequestIEsValue struct {
	Present     int
	AMFUENGAPID *AMFUENGAPID `aper:"referenceFieldValue:10"`
	RANUENGAPID *RANUENGAPID `aper:"referenceFieldValue:85"`
}

type UERadioCapabilityInfoIndicationIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UERadioCapabilityInfoIndicationIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UERadioCapabilityInfoIndicationIEsPresentNothing int = iota /* No components present */
	UERadioCapabilityInfoIndicationIEsPresentAMFUENGAPID
	UERadioCapabilityInfoIndicationIEsPresentRANUENGAPID
	UERadioCapabilityInfoIndicationIEsPresentUERadioCapability
	UERadioCapabilityInfoIndicationIEsPresentUERadioCapabilityForPaging
)

type UERadioCapabilityInfoIndicationIEsValue struct {
	Present                    int
	AMFUENGAPID                *AMFUENGAPID                `aper:"referenceFieldValue:10"`
	RANUENGAPID                *RANUENGAPID                `aper:"referenceFieldValue:85"`
	UERadioCapability          *UERadioCapability          `aper:"referenceFieldValue:117"`
	UERadioCapabilityForPaging *UERadioCapabilityForPaging `aper:"valueExt,referenceFieldValue:118"`
}

type UERadioCapabilityCheckRequestIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UERadioCapabilityCheckRequestIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UERadioCapabilityCheckRequestIEsPresentNothing int = iota /* No components present */
	UERadioCapabilityCheckRequestIEsPresentAMFUENGAPID
	UERadioCapabilityCheckRequestIEsPresentRANUENGAPID
	UERadioCapabilityCheckRequestIEsPresentUERadioCapability
)

type UERadioCapabilityCheckRequestIEsValue struct {
	Present           int
	AMFUENGAPID       *AMFUENGAPID       `aper:"referenceFieldValue:10"`
	RANUENGAPID       *RANUENGAPID       `aper:"referenceFieldValue:85"`
	UERadioCapability *UERadioCapability `aper:"referenceFieldValue:117"`
}

type UERadioCapabilityCheckResponseIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       UERadioCapabilityCheckResponseIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	UERadioCapabilityCheckResponseIEsPresentNothing int = iota /* No components present */
	UERadioCapabilityCheckResponseIEsPresentAMFUENGAPID
	UERadioCapabilityCheckResponseIEsPresentRANUENGAPID
	UERadioCapabilityCheckResponseIEsPresentIMSVoiceSupportIndicator
	UERadioCapabilityCheckResponseIEsPresentCriticalityDiagnostics
)

type UERadioCapabilityCheckResponseIEsValue struct {
	Present                  int
	AMFUENGAPID              *AMFUENGAPID              `aper:"referenceFieldValue:10"`
	RANUENGAPID              *RANUENGAPID              `aper:"referenceFieldValue:85"`
	IMSVoiceSupportIndicator *IMSVoiceSupportIndicator `aper:"referenceFieldValue:30"`
	CriticalityDiagnostics   *CriticalityDiagnostics   `aper:"valueExt,referenceFieldValue:19"`
}

type PDUSessionResourceModifyRequestTransferIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceModifyRequestTransferIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceModifyRequestTransferIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceModifyRequestTransferIEsPresentPDUSessionAggregateMaximumBitRate
	PDUSessionResourceModifyRequestTransferIEsPresentULNGUUPTNLModifyList
	PDUSessionResourceModifyRequestTransferIEsPresentNetworkInstance
	PDUSessionResourceModifyRequestTransferIEsPresentQosFlowAddOrModifyRequestList
	PDUSessionResourceModifyRequestTransferIEsPresentQosFlowToReleaseList
	PDUSessionResourceModifyRequestTransferIEsPresentAdditionalULNGUUPTNLInformation
)

type PDUSessionResourceModifyRequestTransferIEsValue struct {
	Present                           int
	PDUSessionAggregateMaximumBitRate *PDUSessionAggregateMaximumBitRate `aper:"valueExt,referenceFieldValue:130"`
	ULNGUUPTNLModifyList              *ULNGUUPTNLModifyList              `aper:"referenceFieldValue:140"`
	NetworkInstance                   *NetworkInstance                   `aper:"referenceFieldValue:129"`
	QosFlowAddOrModifyRequestList     *QosFlowAddOrModifyRequestList     `aper:"referenceFieldValue:135"`
	QosFlowToReleaseList              *QosFlowList                       `aper:"referenceFieldValue:137"`
	AdditionalULNGUUPTNLInformation   *UPTransportLayerInformation       `aper:"referenceFieldValue:126,valueLB:0,valueUB:1"`
}

type PDUSessionResourceSetupRequestTransferIEs struct {
	Id          ProtocolIEID
	Criticality Criticality
	Value       PDUSessionResourceSetupRequestTransferIEsValue `aper:"openType,referenceFieldName:Id"`
}

const (
	PDUSessionResourceSetupRequestTransferIEsPresentNothing int = iota /* No components present */
	PDUSessionResourceSetupRequestTransferIEsPresentPDUSessionAggregateMaximumBitRate
	PDUSessionResourceSetupRequestTransferIEsPresentULNGUUPTNLInformation
	PDUSessionResourceSetupRequestTransferIEsPresentAdditionalULNGUUPTNLInformation
	PDUSessionResourceSetupRequestTransferIEsPresentDataForwardingNotPossible
	PDUSessionResourceSetupRequestTransferIEsPresentPDUSessionType
	PDUSessionResourceSetupRequestTransferIEsPresentSecurityIndication
	PDUSessionResourceSetupRequestTransferIEsPresentNetworkInstance
	PDUSessionResourceSetupRequestTransferIEsPresentQosFlowSetupRequestList
)

type PDUSessionResourceSetupRequestTransferIEsValue struct {
	Present                           int
	PDUSessionAggregateMaximumBitRate *PDUSessionAggregateMaximumBitRate `aper:"valueExt,referenceFieldValue:130"`
	ULNGUUPTNLInformation             *UPTransportLayerInformation       `aper:"referenceFieldValue:139,valueLB:0,valueUB:1"`
	AdditionalULNGUUPTNLInformation   *UPTransportLayerInformation       `aper:"referenceFieldValue:126,valueLB:0,valueUB:1"`
	DataForwardingNotPossible         *DataForwardingNotPossible         `aper:"referenceFieldValue:127"`
	PDUSessionType                    *PDUSessionType                    `aper:"referenceFieldValue:134"`
	SecurityIndication                *SecurityIndication                `aper:"valueExt,referenceFieldValue:138"`
	NetworkInstance                   *NetworkInstance                   `aper:"referenceFieldValue:129"`
	QosFlowSetupRequestList           *QosFlowSetupRequestList           `aper:"referenceFieldValue:136"`
}
