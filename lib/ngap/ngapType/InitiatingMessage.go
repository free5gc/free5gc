//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type InitiatingMessage struct {
	ProcedureCode ProcedureCode
	Criticality   Criticality
	Value         InitiatingMessageValue `aper:"openType,referenceFieldName:ProcedureCode"`
}

const (
	InitiatingMessagePresentNothing int = iota /* No components present */
	InitiatingMessagePresentAMFConfigurationUpdate
	InitiatingMessagePresentHandoverCancel
	InitiatingMessagePresentHandoverRequired
	InitiatingMessagePresentHandoverRequest
	InitiatingMessagePresentInitialContextSetupRequest
	InitiatingMessagePresentNGReset
	InitiatingMessagePresentNGSetupRequest
	InitiatingMessagePresentPathSwitchRequest
	InitiatingMessagePresentPDUSessionResourceModifyRequest
	InitiatingMessagePresentPDUSessionResourceModifyIndication
	InitiatingMessagePresentPDUSessionResourceReleaseCommand
	InitiatingMessagePresentPDUSessionResourceSetupRequest
	InitiatingMessagePresentPWSCancelRequest
	InitiatingMessagePresentRANConfigurationUpdate
	InitiatingMessagePresentUEContextModificationRequest
	InitiatingMessagePresentUEContextReleaseCommand
	InitiatingMessagePresentUERadioCapabilityCheckRequest
	InitiatingMessagePresentWriteReplaceWarningRequest
	InitiatingMessagePresentAMFStatusIndication
	InitiatingMessagePresentCellTrafficTrace
	InitiatingMessagePresentDeactivateTrace
	InitiatingMessagePresentDownlinkNASTransport
	InitiatingMessagePresentDownlinkNonUEAssociatedNRPPaTransport
	InitiatingMessagePresentDownlinkRANConfigurationTransfer
	InitiatingMessagePresentDownlinkRANStatusTransfer
	InitiatingMessagePresentDownlinkUEAssociatedNRPPaTransport
	InitiatingMessagePresentErrorIndication
	InitiatingMessagePresentHandoverNotify
	InitiatingMessagePresentInitialUEMessage
	InitiatingMessagePresentLocationReport
	InitiatingMessagePresentLocationReportingControl
	InitiatingMessagePresentLocationReportingFailureIndication
	InitiatingMessagePresentNASNonDeliveryIndication
	InitiatingMessagePresentOverloadStart
	InitiatingMessagePresentOverloadStop
	InitiatingMessagePresentPaging
	InitiatingMessagePresentPDUSessionResourceNotify
	InitiatingMessagePresentPrivateMessage
	InitiatingMessagePresentPWSFailureIndication
	InitiatingMessagePresentPWSRestartIndication
	InitiatingMessagePresentRerouteNASRequest
	InitiatingMessagePresentRRCInactiveTransitionReport
	InitiatingMessagePresentTraceFailureIndication
	InitiatingMessagePresentTraceStart
	InitiatingMessagePresentUEContextReleaseRequest
	InitiatingMessagePresentUERadioCapabilityInfoIndication
	InitiatingMessagePresentUETNLABindingReleaseRequest
	InitiatingMessagePresentUplinkNASTransport
	InitiatingMessagePresentUplinkNonUEAssociatedNRPPaTransport
	InitiatingMessagePresentUplinkRANConfigurationTransfer
	InitiatingMessagePresentUplinkRANStatusTransfer
	InitiatingMessagePresentUplinkUEAssociatedNRPPaTransport
)

type InitiatingMessageValue struct {
	Present                               int
	AMFConfigurationUpdate                *AMFConfigurationUpdate                `aper:"valueExt,referenceFieldValue:0"`
	HandoverCancel                        *HandoverCancel                        `aper:"valueExt,referenceFieldValue:10"`
	HandoverRequired                      *HandoverRequired                      `aper:"valueExt,referenceFieldValue:12"`
	HandoverRequest                       *HandoverRequest                       `aper:"valueExt,referenceFieldValue:13"`
	InitialContextSetupRequest            *InitialContextSetupRequest            `aper:"valueExt,referenceFieldValue:14"`
	NGReset                               *NGReset                               `aper:"valueExt,referenceFieldValue:20"`
	NGSetupRequest                        *NGSetupRequest                        `aper:"valueExt,referenceFieldValue:21"`
	PathSwitchRequest                     *PathSwitchRequest                     `aper:"valueExt,referenceFieldValue:25"`
	PDUSessionResourceModifyRequest       *PDUSessionResourceModifyRequest       `aper:"valueExt,referenceFieldValue:26"`
	PDUSessionResourceModifyIndication    *PDUSessionResourceModifyIndication    `aper:"valueExt,referenceFieldValue:27"`
	PDUSessionResourceReleaseCommand      *PDUSessionResourceReleaseCommand      `aper:"valueExt,referenceFieldValue:28"`
	PDUSessionResourceSetupRequest        *PDUSessionResourceSetupRequest        `aper:"valueExt,referenceFieldValue:29"`
	PWSCancelRequest                      *PWSCancelRequest                      `aper:"valueExt,referenceFieldValue:32"`
	RANConfigurationUpdate                *RANConfigurationUpdate                `aper:"valueExt,referenceFieldValue:35"`
	UEContextModificationRequest          *UEContextModificationRequest          `aper:"valueExt,referenceFieldValue:40"`
	UEContextReleaseCommand               *UEContextReleaseCommand               `aper:"valueExt,referenceFieldValue:41"`
	UERadioCapabilityCheckRequest         *UERadioCapabilityCheckRequest         `aper:"valueExt,referenceFieldValue:43"`
	WriteReplaceWarningRequest            *WriteReplaceWarningRequest            `aper:"valueExt,referenceFieldValue:51"`
	AMFStatusIndication                   *AMFStatusIndication                   `aper:"valueExt,referenceFieldValue:1"`
	CellTrafficTrace                      *CellTrafficTrace                      `aper:"valueExt,referenceFieldValue:2"`
	DeactivateTrace                       *DeactivateTrace                       `aper:"valueExt,referenceFieldValue:3"`
	DownlinkNASTransport                  *DownlinkNASTransport                  `aper:"valueExt,referenceFieldValue:4"`
	DownlinkNonUEAssociatedNRPPaTransport *DownlinkNonUEAssociatedNRPPaTransport `aper:"valueExt,referenceFieldValue:5"`
	DownlinkRANConfigurationTransfer      *DownlinkRANConfigurationTransfer      `aper:"valueExt,referenceFieldValue:6"`
	DownlinkRANStatusTransfer             *DownlinkRANStatusTransfer             `aper:"valueExt,referenceFieldValue:7"`
	DownlinkUEAssociatedNRPPaTransport    *DownlinkUEAssociatedNRPPaTransport    `aper:"valueExt,referenceFieldValue:8"`
	ErrorIndication                       *ErrorIndication                       `aper:"valueExt,referenceFieldValue:9"`
	HandoverNotify                        *HandoverNotify                        `aper:"valueExt,referenceFieldValue:11"`
	InitialUEMessage                      *InitialUEMessage                      `aper:"valueExt,referenceFieldValue:15"`
	LocationReport                        *LocationReport                        `aper:"valueExt,referenceFieldValue:18"`
	LocationReportingControl              *LocationReportingControl              `aper:"valueExt,referenceFieldValue:16"`
	LocationReportingFailureIndication    *LocationReportingFailureIndication    `aper:"valueExt,referenceFieldValue:17"`
	NASNonDeliveryIndication              *NASNonDeliveryIndication              `aper:"valueExt,referenceFieldValue:19"`
	OverloadStart                         *OverloadStart                         `aper:"valueExt,referenceFieldValue:22"`
	OverloadStop                          *OverloadStop                          `aper:"valueExt,referenceFieldValue:23"`
	Paging                                *Paging                                `aper:"valueExt,referenceFieldValue:24"`
	PDUSessionResourceNotify              *PDUSessionResourceNotify              `aper:"valueExt,referenceFieldValue:30"`
	PrivateMessage                        *PrivateMessage                        `aper:"valueExt,referenceFieldValue:31"`
	PWSFailureIndication                  *PWSFailureIndication                  `aper:"valueExt,referenceFieldValue:33"`
	PWSRestartIndication                  *PWSRestartIndication                  `aper:"valueExt,referenceFieldValue:34"`
	RerouteNASRequest                     *RerouteNASRequest                     `aper:"valueExt,referenceFieldValue:36"`
	RRCInactiveTransitionReport           *RRCInactiveTransitionReport           `aper:"valueExt,referenceFieldValue:37"`
	TraceFailureIndication                *TraceFailureIndication                `aper:"valueExt,referenceFieldValue:38"`
	TraceStart                            *TraceStart                            `aper:"valueExt,referenceFieldValue:39"`
	UEContextReleaseRequest               *UEContextReleaseRequest               `aper:"valueExt,referenceFieldValue:42"`
	UERadioCapabilityInfoIndication       *UERadioCapabilityInfoIndication       `aper:"valueExt,referenceFieldValue:44"`
	UETNLABindingReleaseRequest           *UETNLABindingReleaseRequest           `aper:"valueExt,referenceFieldValue:45"`
	UplinkNASTransport                    *UplinkNASTransport                    `aper:"valueExt,referenceFieldValue:46"`
	UplinkNonUEAssociatedNRPPaTransport   *UplinkNonUEAssociatedNRPPaTransport   `aper:"valueExt,referenceFieldValue:47"`
	UplinkRANConfigurationTransfer        *UplinkRANConfigurationTransfer        `aper:"valueExt,referenceFieldValue:48"`
	UplinkRANStatusTransfer               *UplinkRANStatusTransfer               `aper:"valueExt,referenceFieldValue:49"`
	UplinkUEAssociatedNRPPaTransport      *UplinkUEAssociatedNRPPaTransport      `aper:"valueExt,referenceFieldValue:50"`
}
