//go:binary-only-package

package pfcp

import (
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/util_3gpp"
)

type Message struct {
	Header Header
	Body   interface{}
}

type HeartbeatRequest struct {
	RecoveryTimeStamp *pfcpType.RecoveryTimeStamp `tlv:"96"`
}

type HeartbeatResponse struct {
	RecoveryTimeStamp *pfcpType.RecoveryTimeStamp `tlv:"96"`
}

type PFCPPFDManagementRequest struct {
	ApplicationIDsPFDs *ApplicationIDsPFDs `tlv:"58"`
}

type ApplicationIDsPFDs struct {
	ApplicationID *pfcpType.ApplicationID `tlv:"24"`
	PFD           *pfcpType.PFDContext    `tlv:"59"`
}

type PFD struct {
	PFDContents *pfcpType.PFDContents `tlv:"61"`
}

type PFCPPFDManagementResponse struct {
	Cause       *pfcpType.Cause       `tlv:"19"`
	OffendingIE *pfcpType.OffendingIE `tlv:"40"`
}

type PFCPAssociationSetupRequest struct {
	NodeID                         *pfcpType.NodeID                         `tlv:"60"`
	RecoveryTimeStamp              *pfcpType.RecoveryTimeStamp              `tlv:"96"`
	UPFunctionFeatures             *pfcpType.UPFunctionFeatures             `tlv:"43"`
	CPFunctionFeatures             *pfcpType.CPFunctionFeatures             `tlv:"89"`
	UserPlaneIPResourceInformation *pfcpType.UserPlaneIPResourceInformation `tlv:"116"`
}

type PFCPAssociationSetupResponse struct {
	NodeID                         *pfcpType.NodeID                         `tlv:"60"`
	Cause                          *pfcpType.Cause                          `tlv:"19"`
	RecoveryTimeStamp              *pfcpType.RecoveryTimeStamp              `tlv:"96"`
	UPFunctionFeatures             *pfcpType.UPFunctionFeatures             `tlv:"43"`
	CPFunctionFeatures             *pfcpType.CPFunctionFeatures             `tlv:"89"`
	UserPlaneIPResourceInformation *pfcpType.UserPlaneIPResourceInformation `tlv:"116"`
}

type PFCPAssociationUpdateRequest struct {
	NodeID                         *pfcpType.NodeID                         `tlv:"60"`
	UPFunctionFeatures             *pfcpType.UPFunctionFeatures             `tlv:"43"`
	CPFunctionFeatures             *pfcpType.CPFunctionFeatures             `tlv:"89"`
	PFCPAssociationReleaseRequest  *PFCPAssociationReleaseRequest           `tlv:"111"`
	GracefulReleasePeriod          *pfcpType.GracefulReleasePeriod          `tlv:"112"`
	UserPlaneIPResourceInformation *pfcpType.UserPlaneIPResourceInformation `tlv:"116"`
}

type PFCPAssociationUpdateResponse struct {
	NodeID             *pfcpType.NodeID             `tlv:"60"`
	Cause              *pfcpType.Cause              `tlv:"19"`
	UPFunctionFeatures *pfcpType.UPFunctionFeatures `tlv:"43"`
	CPFunctionFeatures *pfcpType.CPFunctionFeatures `tlv:"89"`
}

type PFCPAssociationReleaseRequest struct {
	NodeID *pfcpType.NodeID `tlv:"60"`
}

type PFCPAssociationReleaseResponse struct {
	NodeID *pfcpType.NodeID `tlv:"60"`
	Cause  *pfcpType.Cause  `tlv:"19"`
}

type PFCPNodeReportRequest struct {
	NodeID                     *pfcpType.NodeID                     `tlv:"60"`
	NodeReportType             *pfcpType.NodeReportType             `tlv:"101"`
	UserPlanePathFailureReport *pfcpType.UserPlanePathFailureReport `tlv:"102"`
}

type UserPlanePathFailure struct {
	RemoteGTPUPeer *pfcpType.RemoteGTPUPeer `tlv:"103"`
}

type PFCPNodeReportResponse struct {
	NodeID      *pfcpType.NodeID      `tlv:"60"`
	Cause       *pfcpType.Cause       `tlv:"19"`
	OffendingIE *pfcpType.OffendingIE `tlv:"40"`
}

type PFCPSessionSetDeletionRequest struct {
	NodeID     *pfcpType.NodeID `tlv:"60"`
	SGWCFQCSID *pfcpType.FQCSID `tlv:"65"`
	PGWCFQCSID *pfcpType.FQCSID `tlv:"65"`
	SGWUFQCSID *pfcpType.FQCSID `tlv:"65"`
	PGWUFQCSID *pfcpType.FQCSID `tlv:"65"`
	TWANFQCSID *pfcpType.FQCSID `tlv:"65"`
	EPDGFQCSID *pfcpType.FQCSID `tlv:"65"`
	MMEFQCSID  *pfcpType.FQCSID `tlv:"65"`
}

type PFCPSessionSetDeletionResponse struct {
	NodeID      *pfcpType.NodeID      `tlv:"60"`
	Cause       *pfcpType.Cause       `tlv:"19"`
	OffendingIE *pfcpType.OffendingIE `tlv:"40"`
}

type PFCPSessionEstablishmentRequest struct {
	NodeID                   *pfcpType.NodeID                   `tlv:"60"`
	CPFSEID                  *pfcpType.FSEID                    `tlv:"57"`
	CreatePDR                []*CreatePDR                       `tlv:"1"`
	CreateFAR                []*CreateFAR                       `tlv:"3"`
	CreateURR                []*CreateURR                       `tlv:"6"`
	CreateQER                []*CreateQER                       `tlv:"7"`
	CreateBAR                []*CreateBAR                       `tlv:"85"`
	CreateTrafficEndpoint    *CreateTrafficEndpoint             `tlv:"127"`
	PDNType                  *pfcpType.PDNType                  `tlv:"113"`
	SGWCFQCSID               *pfcpType.FQCSID                   `tlv:"65"`
	MMEFQCSID                *pfcpType.FQCSID                   `tlv:"65"`
	PGWCFQCSID               *pfcpType.FQCSID                   `tlv:"65"`
	EPDGFQCSID               *pfcpType.FQCSID                   `tlv:"65"`
	TWANFQCSID               *pfcpType.FQCSID                   `tlv:"65"`
	UserPlaneInactivityTimer *pfcpType.UserPlaneInactivityTimer `tlv:"117"`
	UserID                   *pfcpType.UserID                   `tlv:"141"`
	TraceInformation         *pfcpType.TraceInformation         `tlv:"152"`
}

type CreatePDR struct {
	PDRID                   *pfcpType.PacketDetectionRuleID   `tlv:"56"`
	Precedence              *pfcpType.Precedence              `tlv:"29"`
	PDI                     *PDI                              `tlv:"2"`
	OuterHeaderRemoval      *pfcpType.OuterHeaderRemoval      `tlv:"95"`
	FARID                   *pfcpType.FARID                   `tlv:"108"`
	URRID                   *pfcpType.URRID                   `tlv:"81"`
	QERID                   *pfcpType.QERID                   `tlv:"109"`
	ActivatePredefinedRules *pfcpType.ActivatePredefinedRules `tlv:"106"`
}

type PDI struct {
	SourceInterface               *pfcpType.SourceInterface               `tlv:"20"`
	LocalFTEID                    *pfcpType.FTEID                         `tlv:"21"`
	NetworkInstance               *util_3gpp.Dnn                          `tlv:"22"`
	UEIPAddress                   *pfcpType.UEIPAddress                   `tlv:"93"`
	TrafficEndpointID             *pfcpType.TrafficEndpointID             `tlv:"131"`
	SDFFilter                     *pfcpType.SDFFilter                     `tlv:"23"`
	ApplicationID                 *pfcpType.ApplicationID                 `tlv:"24"`
	EthernetPDUSessionInformation *pfcpType.EthernetPDUSessionInformation `tlv:"142"`
	EthernetPacketFilter          *EthernetPacketFilter                   `tlv:"132"`
	QFI                           *pfcpType.QFI                           `tlv:"124"`
	FramedRoute                   *pfcpType.FramedRoute                   `tlv:"153"`
	FramedRouting                 *pfcpType.FramedRouting                 `tlv:"154"`
	FramedIPv6Route               *pfcpType.FramedIPv6Route               `tlv:"155"`
}

type EthernetPacketFilter struct {
	EthernetFilterID         *pfcpType.EthernetFilterID         `tlv:"138"`
	EthernetFilterProperties *pfcpType.EthernetFilterProperties `tlv:"139"`
	MACAddress               *pfcpType.MACAddress               `tlv:"133"`
	Ethertype                *pfcpType.Ethertype                `tlv:"136"`
	CTAG                     *pfcpType.CTAG                     `tlv:"134"`
	STAG                     *pfcpType.STAG                     `tlv:"135"`
	SDFFilter                *pfcpType.SDFFilter                `tlv:"23"`
}

type CreateFAR struct {
	FARID                 *pfcpType.FARID                 `tlv:"108"`
	ApplyAction           *pfcpType.ApplyAction           `tlv:"44"`
	ForwardingParameters  *ForwardingParametersIEInFAR    `tlv:"4"`
	DuplicatingParameters *pfcpType.DuplicatingParameters `tlv:"5"`
	BARID                 *pfcpType.BARID                 `tlv:"88"`
}

type ForwardingParametersIEInFAR struct {
	DestinationInterface    *pfcpType.DestinationInterface  `tlv:"42"`
	NetworkInstance         *util_3gpp.Dnn                  `tlv:"22"`
	RedirectInformation     *pfcpType.RedirectInformation   `tlv:"38"`
	OuterHeaderCreation     *pfcpType.OuterHeaderCreation   `tlv:"84"`
	TransportLevelMarking   *pfcpType.TransportLevelMarking `tlv:"30"`
	ForwardingPolicy        *pfcpType.ForwardingPolicy      `tlv:"41"`
	HeaderEnrichment        *pfcpType.HeaderEnrichment      `tlv:"98"`
	LinkedTrafficEndpointID *pfcpType.TrafficEndpointID     `tlv:"131"`
	Proxying                *pfcpType.Proxying              `tlv:"137"`
}

type DuplicatingParametersIEInFAR struct {
	DestinationInterface  *pfcpType.DestinationInterface  `tlv:"42"`
	OuterHeaderCreation   *pfcpType.OuterHeaderCreation   `tlv:"84"`
	TransportLevelMarking *pfcpType.TransportLevelMarking `tlv:"30"`
	ForwardingPolicy      *pfcpType.ForwardingPolicy      `tlv:"41"`
}

type CreateURR struct {
	URRID                     *pfcpType.URRID                     `tlv:"81"`
	MeasurementMethod         *pfcpType.MeasurementMethod         `tlv:"62"`
	ReportingTriggers         *pfcpType.ReportingTriggers         `tlv:"37"`
	MeasurementPeriod         *pfcpType.MeasurementPeriod         `tlv:"64"`
	VolumeThreshold           *pfcpType.VolumeThreshold           `tlv:"31"`
	VolumeQuota               *pfcpType.VolumeQuota               `tlv:"73"`
	TimeThreshold             *pfcpType.TimeThreshold             `tlv:"32"`
	TimeQuota                 *pfcpType.TimeQuota                 `tlv:"74"`
	QuotaHoldingTime          *pfcpType.QuotaHoldingTime          `tlv:"71"`
	DroppedDLTrafficThreshold *pfcpType.DroppedDLTrafficThreshold `tlv:"72"`
	MonitoringTime            *pfcpType.MonitoringTime            `tlv:"33"`
	EventInformation          *EventInformation                   `tlv:"148"`
	SubsequentVolumeThreshold *pfcpType.SubsequentVolumeThreshold `tlv:"34"`
	SubsequentTimeThreshold   *pfcpType.SubsequentTimeThreshold   `tlv:"35"`
	SubsequentVolumeQuota     *pfcpType.SubsequentVolumeQuota     `tlv:"121"`
	SubsequentTimeQuota       *pfcpType.SubsequentTimeQuota       `tlv:"122"`
	InactivityDetectionTime   *pfcpType.InactivityDetectionTime   `tlv:"36"`
	LinkedURRID               *pfcpType.LinkedURRID               `tlv:"82"`
	MeasurementInformation    *pfcpType.MeasurementInformation    `tlv:"100"`
	TimeQuotaMechanism        *pfcpType.TimeQuotaMechanism        `tlv:"115"`
	AggregatedURRs            *AggregatedURRs                     `tlv:"118"`
	FARIDForQuotaAction       *pfcpType.FARID                     `tlv:"108"`
	EthernetInactivityTimer   *pfcpType.EthernetInactivityTimer   `tlv:"146"`
	AdditionalMonitoringTime  *AdditionalMonitoringTime           `tlv:"147"`
}

type AggregatedURRs struct {
	AggregatedURRID *pfcpType.AggregatedURRID `tlv:"120"`
	Multiplier      *pfcpType.Multiplier      `tlv:"119"`
}

type AdditionalMonitoringTime struct {
	MonitoringTime            *pfcpType.MonitoringTime            `tlv:"33"`
	SubsequentVolumeThreshold *pfcpType.SubsequentVolumeThreshold `tlv:"34"`
	SubsequentTimeThreshold   *pfcpType.SubsequentTimeThreshold   `tlv:"35"`
	SubsequentVolumeQuota     *pfcpType.SubsequentVolumeQuota     `tlv:"121"`
	SubsequentTimeQuota       *pfcpType.SubsequentTimeQuota       `tlv:"122"`
}

type EventInformation struct {
	EventID        *pfcpType.EventID        `tlv:"150"`
	EventThreshold *pfcpType.EventThreshold `tlv:"151"`
}

type CreateQER struct {
	QERID              *pfcpType.QERID              `tlv:"109"`
	QERCorrelationID   *pfcpType.QERCorrelationID   `tlv:"28"`
	GateStatus         *pfcpType.GateStatus         `tlv:"25"`
	MaximumBitrate     *pfcpType.MBR                `tlv:"26"`
	GuaranteedBitrate  *pfcpType.GBR                `tlv:"27"`
	PacketRate         *pfcpType.PacketRate         `tlv:"94"`
	DLFlowLevelMarking *pfcpType.DLFlowLevelMarking `tlv:"97"`
	QoSFlowIdentifier  *pfcpType.QFI                `tlv:"124"`
	ReflectiveQoS      *pfcpType.RQI                `tlv:"123"`
}

type CreateBAR struct {
	BARID                          *pfcpType.BARID                          `tlv:"88"`
	DownlinkDataNotificationDelay  *pfcpType.DownlinkDataNotificationDelay  `tlv:"46"`
	SuggestedBufferingPacketsCount *pfcpType.SuggestedBufferingPacketsCount `tlv:"140"`
}

type CreateTrafficEndpoint struct {
	TrafficEndpointID             *pfcpType.TrafficEndpointID             `tlv:"131"`
	LocalFTEID                    *pfcpType.FTEID                         `tlv:"21"`
	NetworkInstance               *util_3gpp.Dnn                          `tlv:"22"`
	UEIPAddress                   *pfcpType.UEIPAddress                   `tlv:"93"`
	EthernetPDUSessionInformation *pfcpType.EthernetPDUSessionInformation `tlv:"142"`
	FramedRoute                   *pfcpType.FramedRoute                   `tlv:"153"`
	FramedRouting                 *pfcpType.FramedRouting                 `tlv:"154"`
	FramedIPv6Route               *pfcpType.FramedIPv6Route               `tlv:"155"`
}

type PFCPSessionEstablishmentResponse struct {
	NodeID                     *pfcpType.NodeID            `tlv:"60"`
	Cause                      *pfcpType.Cause             `tlv:"19"`
	OffendingIE                *pfcpType.OffendingIE       `tlv:"40"`
	UPFSEID                    *pfcpType.FSEID             `tlv:"57"`
	CreatedPDR                 *CreatedPDR                 `tlv:"8"`
	LoadControlInformation     *LoadControlInformation     `tlv:"51"`
	OverloadControlInformation *OverloadControlInformation `tlv:"54"`
	SGWUFQCSID                 *pfcpType.FQCSID            `tlv:"65"`
	PGWUFQCSID                 *pfcpType.FQCSID            `tlv:"65"`
	FailedRuleID               *pfcpType.FailedRuleID      `tlv:"114"`
	CreatedTrafficEndpoint     *CreatedTrafficEndpoint     `tlv:"128"`
}

type CreatedPDR struct {
	PDRID      *pfcpType.PacketDetectionRuleID `tlv:"56"`
	LocalFTEID *pfcpType.FTEID                 `tlv:"21"`
}

type LoadControlInformation struct {
	LoadControlSequenceNumber *pfcpType.SequenceNumber `tlv:"52"`
	LoadMetric                *pfcpType.Metric         `tlv:"53"`
}

type OverloadControlInformation struct {
	OverloadControlSequenceNumber   *pfcpType.SequenceNumber `tlv:"52"`
	OverloadReductionMetric         *pfcpType.Metric         `tlv:"53"`
	PeriodOfValidity                *pfcpType.Timer          `tlv:"55"`
	OverloadControlInformationFlags *pfcpType.OCIFlags       `tlv:"110"`
}

type CreatedTrafficEndpoint struct {
	TrafficEndpointID *pfcpType.TrafficEndpointID `tlv:"131"`
	LocalFTEID        *pfcpType.FTEID             `tlv:"21"`
}

type PFCPSessionModificationRequest struct {
	CPFSEID                  *pfcpType.FSEID                          `tlv:"57"`
	RemovePDR                []*RemovePDR                             `tlv:"15"`
	RemoveFAR                []*RemoveFAR                             `tlv:"16"`
	RemoveURR                []*RemoveURR                             `tlv:"17"`
	RemoveQER                []*pfcpType.RemoveQER                    `tlv:"18"`
	RemoveBAR                []*RemoveBAR                             `tlv:"87"`
	RemoveTrafficEndpoint    *RemoveTrafficEndpoint                   `tlv:"130"`
	CreatePDR                []*CreatePDR                             `tlv:"1"`
	CreateFAR                []*CreateFAR                             `tlv:"3"`
	CreateURR                []*CreateURR                             `tlv:"6"`
	CreateQER                []*CreateQER                             `tlv:"7"`
	CreateBAR                []*CreateBAR                             `tlv:"85"`
	CreateTrafficEndpoint    *CreateTrafficEndpoint                   `tlv:"127"`
	UpdatePDR                []*UpdatePDR                             `tlv:"9"`
	UpdateFAR                []*UpdateFAR                             `tlv:"10"`
	UpdateURR                []*UpdateURR                             `tlv:"13"`
	UpdateQER                []*UpdateQER                             `tlv:"14"`
	UpdateBAR                *UpdateBARPFCPSessionModificationRequest `tlv:"86"`
	UpdateTrafficEndpoint    *UpdateTrafficEndpoint                   `tlv:"129"`
	PFCPSMReqFlags           *pfcpType.PFCPSMReqFlags                 `tlv:"49"`
	QueryURR                 *QueryURR                                `tlv:"77"`
	PGWCFQCSID               *pfcpType.FQCSID                         `tlv:"65"`
	SGWCFQCSID               *pfcpType.FQCSID                         `tlv:"65"`
	MMEFQCSID                *pfcpType.FQCSID                         `tlv:"65"`
	EPDGFQCSID               *pfcpType.FQCSID                         `tlv:"65"`
	TWANFQCSID               *pfcpType.FQCSID                         `tlv:"65"`
	UserPlaneInactivityTimer *pfcpType.UserPlaneInactivityTimer       `tlv:"117"`
	QueryURRReference        *pfcpType.QueryURRReference              `tlv:"125"`
	TraceInformation         *pfcpType.TraceInformation               `tlv:"152"`
}

type UpdatePDR struct {
	PDRID                     *pfcpType.PacketDetectionRuleID     `tlv:"56"`
	OuterHeaderRemoval        *pfcpType.OuterHeaderRemoval        `tlv:"95"`
	Precedence                *pfcpType.Precedence                `tlv:"29"`
	PDI                       *PDI                                `tlv:"2"`
	FARID                     *pfcpType.FARID                     `tlv:"108"`
	URRID                     *pfcpType.URRID                     `tlv:"81"`
	QERID                     *pfcpType.QERID                     `tlv:"109"`
	ActivatePredefinedRules   *pfcpType.ActivatePredefinedRules   `tlv:"106"`
	DeactivatePredefinedRules *pfcpType.DeactivatePredefinedRules `tlv:"107"`
}

type UpdateFAR struct {
	FARID                       *pfcpType.FARID                       `tlv:"108"`
	ApplyAction                 *pfcpType.ApplyAction                 `tlv:"44"`
	UpdateForwardingParameters  *UpdateForwardingParametersIEInFAR    `tlv:"11"`
	UpdateDuplicatingParameters *pfcpType.UpdateDuplicatingParameters `tlv:"105"`
	BARID                       *pfcpType.BARID                       `tlv:"88"`
}

type UpdateForwardingParametersIEInFAR struct {
	DestinationInterface    *pfcpType.DestinationInterface  `tlv:"42"`
	NetworkInstance         *util_3gpp.Dnn                  `tlv:"22"`
	RedirectInformation     *pfcpType.RedirectInformation   `tlv:"38"`
	OuterHeaderCreation     *pfcpType.OuterHeaderCreation   `tlv:"84"`
	TransportLevelMarking   *pfcpType.TransportLevelMarking `tlv:"30"`
	ForwardingPolicy        *pfcpType.ForwardingPolicy      `tlv:"41"`
	HeaderEnrichment        *pfcpType.HeaderEnrichment      `tlv:"98"`
	PFCPSMReqFlags          *pfcpType.PFCPSMReqFlags        `tlv:"49"`
	LinkedTrafficEndpointID *pfcpType.TrafficEndpointID     `tlv:"131"`
}

type UpdateDuplicatingParametersIEInFAR struct {
	DestinationInterface  *pfcpType.DestinationInterface  `tlv:"42"`
	OuterHeaderCreation   *pfcpType.OuterHeaderCreation   `tlv:"84"`
	TransportLevelMarking *pfcpType.TransportLevelMarking `tlv:"30"`
	ForwardingPolicy      *pfcpType.ForwardingPolicy      `tlv:"41"`
}

type UpdateURR struct {
	URRID                     *pfcpType.URRID                     `tlv:"81"`
	MeasurementMethod         *pfcpType.MeasurementMethod         `tlv:"62"`
	ReportingTriggers         *pfcpType.ReportingTriggers         `tlv:"37"`
	MeasurementPeriod         *pfcpType.MeasurementPeriod         `tlv:"64"`
	VolumeThreshold           *pfcpType.VolumeThreshold           `tlv:"31"`
	VolumeQuota               *pfcpType.VolumeQuota               `tlv:"73"`
	TimeThreshold             *pfcpType.TimeThreshold             `tlv:"32"`
	TimeQuota                 *pfcpType.TimeQuota                 `tlv:"74"`
	QuotaHoldingTime          *pfcpType.QuotaHoldingTime          `tlv:"71"`
	DroppedDLTrafficThreshold *pfcpType.DroppedDLTrafficThreshold `tlv:"72"`
	MonitoringTime            *pfcpType.MonitoringTime            `tlv:"33"`
	EventInformation          *EventInformation                   `tlv:"148"`
	SubsequentVolumeThreshold *pfcpType.SubsequentVolumeThreshold `tlv:"34"`
	SubsequentTimeThreshold   *pfcpType.SubsequentTimeThreshold   `tlv:"35"`
	SubsequentVolumeQuota     *pfcpType.SubsequentVolumeQuota     `tlv:"121"`
	SubsequentTimeQuota       *pfcpType.SubsequentTimeQuota       `tlv:"122"`
	InactivityDetectionTime   *pfcpType.InactivityDetectionTime   `tlv:"36"`
	LinkedURRID               *pfcpType.LinkedURRID               `tlv:"82"`
	MeasurementInformation    *pfcpType.MeasurementInformation    `tlv:"100"`
	TimeQuotaMechanism        *pfcpType.TimeQuotaMechanism        `tlv:"115"`
	AggregatedURRs            *AggregatedURRs                     `tlv:"118"`
	FARIDForQuotaAction       *pfcpType.FARID                     `tlv:"108"`
	EthernetInactivityTimer   *pfcpType.EthernetInactivityTimer   `tlv:"146"`
	AdditionalMonitoringTime  *AdditionalMonitoringTime           `tlv:"147"`
}

type UpdateQER struct {
	QERID              *pfcpType.QERID              `tlv:"109"`
	QERCorrelationID   *pfcpType.QERCorrelationID   `tlv:"28"`
	GateStatus         *pfcpType.GateStatus         `tlv:"25"`
	MaximumBitrate     *pfcpType.MBR                `tlv:"26"`
	GuaranteedBitrate  *pfcpType.GBR                `tlv:"27"`
	PacketRate         *pfcpType.PacketRate         `tlv:"94"`
	DLFlowLevelMarking *pfcpType.DLFlowLevelMarking `tlv:"97"`
	QoSFlowIdentifier  *pfcpType.QFI                `tlv:"124"`
	ReflectiveQoS      *pfcpType.RQI                `tlv:"123"`
}

type RemovePDR struct {
	PDRID *pfcpType.PacketDetectionRuleID `tlv:"56"`
}

type RemoveFAR struct {
	FARID *pfcpType.FARID `tlv:"108"`
}

type RemoveURR struct {
	URRID *pfcpType.URRID `tlv:"81"`
}

type RemoveQERIEPFCPSessionModificationRequest struct {
	QERID *pfcpType.QERID `tlv:"109"`
}

type QueryURR struct {
	URRID *pfcpType.URRID `tlv:"81"`
}

type UpdateBARPFCPSessionModificationRequest struct {
	BARID                          *pfcpType.BARID                          `tlv:"88"`
	DownlinkDataNotificationDelay  *pfcpType.DownlinkDataNotificationDelay  `tlv:"46"`
	SuggestedBufferingPacketsCount *pfcpType.SuggestedBufferingPacketsCount `tlv:"140"`
}

type RemoveBAR struct {
	BARID *pfcpType.BARID `tlv:"88"`
}

type UpdateTrafficEndpoint struct {
	TrafficEndpointID *pfcpType.TrafficEndpointID `tlv:"131"`
	LocalFTEID        *pfcpType.FTEID             `tlv:"21"`
	NetworkInstance   *util_3gpp.Dnn              `tlv:"22"`
	UEIPAddress       *pfcpType.UEIPAddress       `tlv:"93"`
	FramedRoute       *pfcpType.FramedRoute       `tlv:"153"`
	FramedRouting     *pfcpType.FramedRouting     `tlv:"154"`
	FramedIPv6Route   *pfcpType.FramedIPv6Route   `tlv:"155"`
}

type RemoveTrafficEndpoint struct {
	TrafficEndpointID *pfcpType.TrafficEndpointID `tlv:"131"`
}

type PFCPSessionModificationResponse struct {
	Cause                             *pfcpType.Cause                             `tlv:"19"`
	OffendingIE                       *pfcpType.OffendingIE                       `tlv:"40"`
	CreatedPDR                        *CreatedPDR                                 `tlv:"8"`
	LoadControlInformation            *LoadControlInformation                     `tlv:"51"`
	OverloadControlInformation        *OverloadControlInformation                 `tlv:"54"`
	UsageReport                       *UsageReportPFCPSessionModificationResponse `tlv:"78"`
	FailedRuleID                      *pfcpType.FailedRuleID                      `tlv:"114"`
	AdditionalUsageReportsInformation *pfcpType.AdditionalUsageReportsInformation `tlv:"126"`
	CreatedUpdatedTrafficEndpoint     *CreatedTrafficEndpoint                     `tlv:"128"`
}

type UsageReportPFCPSessionModificationResponse struct {
	URRID                      *pfcpType.URRID               `tlv:"81"`
	URSEQN                     *pfcpType.URSEQN              `tlv:"104"`
	UsageReportTrigger         *pfcpType.UsageReportTrigger  `tlv:"63"`
	StartTime                  *pfcpType.StartTime           `tlv:"75"`
	EndTime                    *pfcpType.EndTime             `tlv:"76"`
	VolumeMeasurement          *pfcpType.VolumeMeasurement   `tlv:"66"`
	DurationMeasurement        *pfcpType.DurationMeasurement `tlv:"67"`
	TimeOfFirstPacket          *pfcpType.TimeOfFirstPacket   `tlv:"69"`
	TimeOfLastPacket           *pfcpType.TimeOfLastPacket    `tlv:"70"`
	UsageInformation           *pfcpType.UsageInformation    `tlv:"90"`
	QueryURRReference          *pfcpType.QueryURRReference   `tlv:"125"`
	EthernetTrafficInformation *EthernetTrafficInformation   `tlv:"143"`
}

type PFCPSessionDeletionRequest struct {
}

type PFCPSessionDeletionResponse struct {
	Cause                      *pfcpType.Cause                         `tlv:"19"`
	OffendingIE                *pfcpType.OffendingIE                   `tlv:"40"`
	LoadControlInformation     *LoadControlInformation                 `tlv:"51"`
	OverloadControlInformation *OverloadControlInformation             `tlv:"54"`
	UsageReport                *UsageReportPFCPSessionDeletionResponse `tlv:"79"`
}

type UsageReportPFCPSessionDeletionResponse struct {
	URRID                      *pfcpType.URRID               `tlv:"81"`
	URSEQN                     *pfcpType.URSEQN              `tlv:"104"`
	UsageReportTrigger         *pfcpType.UsageReportTrigger  `tlv:"63"`
	StartTime                  *pfcpType.StartTime           `tlv:"75"`
	EndTime                    *pfcpType.EndTime             `tlv:"76"`
	VolumeMeasurement          *pfcpType.VolumeMeasurement   `tlv:"66"`
	DurationMeasurement        *pfcpType.DurationMeasurement `tlv:"67"`
	TimeOfFirstPacket          *pfcpType.TimeOfFirstPacket   `tlv:"69"`
	TimeOfLastPacket           *pfcpType.TimeOfLastPacket    `tlv:"70"`
	UsageInformation           *pfcpType.UsageInformation    `tlv:"90"`
	EthernetTrafficInformation *EthernetTrafficInformation   `tlv:"143"`
}

type PFCPSessionReportRequest struct {
	ReportType                        *pfcpType.ReportType                        `tlv:"39"`
	DownlinkDataReport                *DownlinkDataReport                         `tlv:"83"`
	UsageReport                       *UsageReportPFCPSessionReportRequest        `tlv:"80"`
	ErrorIndicationReport             *ErrorIndicationReport                      `tlv:"99"`
	LoadControlInformation            *LoadControlInformation                     `tlv:"51"`
	OverloadControlInformation        *OverloadControlInformation                 `tlv:"54"`
	AdditionalUsageReportsInformation *pfcpType.AdditionalUsageReportsInformation `tlv:"126"`
}

type DownlinkDataReport struct {
	PDRID                          *pfcpType.PacketDetectionRuleID          `tlv:"56"`
	DownlinkDataServiceInformation *pfcpType.DownlinkDataServiceInformation `tlv:"45"`
}

type UsageReportPFCPSessionReportRequest struct {
	URRID                           *pfcpType.URRID                  `tlv:"81"`
	URSEQN                          *pfcpType.URSEQN                 `tlv:"104"`
	UsageReportTrigger              *pfcpType.UsageReportTrigger     `tlv:"63"`
	StartTime                       *pfcpType.StartTime              `tlv:"75"`
	EndTime                         *pfcpType.EndTime                `tlv:"76"`
	VolumeMeasurement               *pfcpType.VolumeMeasurement      `tlv:"66"`
	DurationMeasurement             *pfcpType.DurationMeasurement    `tlv:"67"`
	ApplicationDetectionInformation *ApplicationDetectionInformation `tlv:"68"`
	UEIPAddress                     *pfcpType.UEIPAddress            `tlv:"93"`
	NetworkInstance                 *util_3gpp.Dnn                   `tlv:"22"`
	TimeOfFirstPacket               *pfcpType.TimeOfFirstPacket      `tlv:"69"`
	TimeOfLastPacket                *pfcpType.TimeOfLastPacket       `tlv:"70"`
	UsageInformation                *pfcpType.UsageInformation       `tlv:"90"`
	QueryURRReference               *pfcpType.QueryURRReference      `tlv:"125"`
	EventReporting                  *EventReporting                  `tlv:"149"`
	EthernetTrafficInformation      *EthernetTrafficInformation      `tlv:"143"`
}

type ApplicationDetectionInformation struct {
	ApplicationID         *pfcpType.ApplicationID         `tlv:"24"`
	ApplicationInstanceID *pfcpType.ApplicationInstanceID `tlv:"91"`
	FlowInformation       *pfcpType.FlowInformation       `tlv:"92"`
}

type EventReporting struct {
	EventID *pfcpType.EventID `tlv:"150"`
}

type EthernetTrafficInformation struct {
	MACAddressesDetected *pfcpType.MACAddressesDetected `tlv:"144"`
	MACAddressesRemoved  *pfcpType.MACAddressesRemoved  `tlv:"145"`
}

type ErrorIndicationReport struct {
	RemoteFTEID *pfcpType.FTEID `tlv:"21"`
}

type PFCPSessionReportResponse struct {
	Cause        *pfcpType.Cause                              `tlv:"19"`
	OffendingIE  *pfcpType.OffendingIE                        `tlv:"40"`
	UpdateBAR    *pfcpType.UpdateBARPFCPSessionReportResponse `tlv:"12"`
	SxSRRspFlags *pfcpType.PFCPSRRspFlags                     `tlv:"50"`
}

type UpdateBARIEInPFCPSessionReportResponse struct {
	BARID                           *pfcpType.BARID                           `tlv:"88"`
	DownlinkDataNotificationDelay   *pfcpType.DownlinkDataNotificationDelay   `tlv:"46"`
	DLBufferingDuration             *pfcpType.DLBufferingDuration             `tlv:"47"`
	DLBufferingSuggestedPacketCount *pfcpType.DLBufferingSuggestedPacketCount `tlv:"48"`
	SuggestedBufferingPacketsCount  *pfcpType.SuggestedBufferingPacketsCount  `tlv:"140"`
}
