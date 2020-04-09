#ifndef __PFCP_MESSAGE_H__
#define __PFCP_MESSAGE_H__

#include <stdint.h>

#include "utlt_debug.h"
#include "utlt_lib.h"
#include "utlt_buff.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct _TlvOctet {
    unsigned long presence;
    uint16_t type;
    uint16_t len;
    void *value;
} __attribute__((packed)) TlvOctet;

typedef struct _IeDescription {
    uint16_t msgType;
    uint16_t msgLen; // msg struct size
    _Bool isTlvObj;
    int numToParse;
    int next[35];
} __attribute__((packed)) IeDescription;

/* 5.1 General format */
#define PFCP_HEADER_LEN     16
#define PFCP_SEID_LEN       8
typedef struct _PfcpHeader {
    union {
        struct {
            ENDIAN4(uint8_t version:3;,
                    uint8_t spare0:3;,
                    uint8_t mp:1;,
                    uint8_t seidP:1;)
        };
        uint8_t flags;
    };
    uint8_t type;
    uint16_t length;
    union {
        struct {
            uint64_t seid;
#define PfcpTransactionId2Sqn(__transactionId) htonl(((__transactionId) << 8))
#define PfcpSqn2TransactionId(__sqn) (ntohl(__sqn) >> 8)
            uint32_t sqn;
        };
        uint32_t sqn_only;
    };
} __attribute__ ((packed)) PfcpHeader;
    
#define RESERVED 0
#define PFCP_HEARTBEAT_REQUEST 1
#define PFCP_HEARTBEAT_RESPONSE 2
#define PFCPPFD_MANAGEMENT_REQUEST 3
#define PFCPPFD_MANAGEMENT_RESPONSE 4
#define PFCP_ASSOCIATION_SETUP_REQUEST 5
#define PFCP_ASSOCIATION_SETUP_RESPONSE 6
#define PFCP_ASSOCIATION_UPDATE_REQUEST 7
#define PFCP_ASSOCIATION_UPDATE_RESPONSE 8
#define PFCP_ASSOCIATION_RELEASE_REQUEST 9
#define PFCP_ASSOCIATION_RELEASE_RESPONSE 10
#define PFCP_VERSION_NOT_SUPPORTED_RESPONSE 11
#define PFCP_NODE_REPORT_REQUEST 12
#define PFCP_NODE_REPORT_RESPONSE 13
#define PFCP_SESSION_SET_DELETION_REQUEST 14
#define PFCP_SESSION_SET_DELETION_RESPONSE 15
#define PFCP_SESSION_ESTABLISHMENT_REQUEST 50
#define PFCP_SESSION_ESTABLISHMENT_RESPONSE 51
#define PFCP_SESSION_MODIFICATION_REQUEST 52
#define PFCP_SESSION_MODIFICATION_RESPONSE 53
#define PFCP_SESSION_DELETION_REQUEST 54
#define PFCP_SESSION_DELETION_RESPONSE 55
#define PFCP_SESSION_REPORT_REQUEST 56
#define PFCP_SESSION_REPORT_RESPONSE 57

#define PFCP_Reserved_TYPE 0
#define PFCP_CreatePDR_TYPE 1
#define PFCP_PDI_TYPE 2
#define PFCP_CreateFAR_TYPE 3
#define PFCP_ForwardingParameters_TYPE 4
#define PFCP_DuplicatingParameters_TYPE 5
#define PFCP_CreateURR_TYPE 6
#define PFCP_CreateQER_TYPE 7
#define PFCP_CreatedPDR_TYPE 8
#define PFCP_UpdatePDR_TYPE 9
#define PFCP_UpdateFAR_TYPE 10
#define PFCP_UpdateForwardingParameters_TYPE 11
#define PFCP_UpdateBARPFCPSessionReportResponse_TYPE 12
#define PFCP_UpdateURR_TYPE 13
#define PFCP_UpdateQER_TYPE 14
#define PFCP_RemovePDR_TYPE 15
#define PFCP_RemoveFAR_TYPE 16
#define PFCP_RemoveURR_TYPE 17
#define PFCP_RemoveQER_TYPE 18
#define PFCP_Cause_TYPE 19
#define PFCP_SourceInterface_TYPE 20
#define PFCP_FTEID_TYPE 21
#define PFCP_NetworkInstance_TYPE 22
#define PFCP_SDFFilter_TYPE 23
#define PFCP_ApplicationID_TYPE 24
#define PFCP_GateStatus_TYPE 25
#define PFCP_MBR_TYPE 26
#define PFCP_GBR_TYPE 27
#define PFCP_QERCorrelationID_TYPE 28
#define PFCP_Precedence_TYPE 29
#define PFCP_TransportLevelMarking_TYPE 30
#define PFCP_VolumeThreshold_TYPE 31
#define PFCP_TimeThreshold_TYPE 32
#define PFCP_MonitoringTime_TYPE 33
#define PFCP_SubsequentVolumeThreshold_TYPE 34
#define PFCP_SubsequentTimeThreshold_TYPE 35
#define PFCP_InactivityDetectionTime_TYPE 36
#define PFCP_ReportingTriggers_TYPE 37
#define PFCP_RedirectInformation_TYPE 38
#define PFCP_ReportType_TYPE 39
#define PFCP_OffendingIE_TYPE 40
#define PFCP_ForwardingPolicy_TYPE 41
#define PFCP_DestinationInterface_TYPE 42
#define PFCP_UPFunctionFeatures_TYPE 43
#define PFCP_ApplyAction_TYPE 44
#define PFCP_DownlinkDataServiceInformation_TYPE 45
#define PFCP_DownlinkDataNotificationDelay_TYPE 46
#define PFCP_DLBufferingDuration_TYPE 47
#define PFCP_DLBufferingSuggestedPacketCount_TYPE 48
#define PFCP_PFCPSMReqFlags_TYPE 49
#define PFCP_PFCPSRRspFlags_TYPE 50
#define PFCP_LoadControlInformation_TYPE 51
#define PFCP_SequenceNumber_TYPE 52
#define PFCP_Metric_TYPE 53
#define PFCP_OverloadControlInformation_TYPE 54
#define PFCP_Timer_TYPE 55
#define PFCP_PacketDetectionRuleID_TYPE 56
#define PFCP_FSEID_TYPE 57
#define PFCP_ApplicationIDsPFDs_TYPE 58
#define PFCP_PFDContext_TYPE 59
#define PFCP_NodeID_TYPE 60
#define PFCP_PFDContents_TYPE 61
#define PFCP_MeasurementMethod_TYPE 62
#define PFCP_UsageReportTrigger_TYPE 63
#define PFCP_MeasurementPeriod_TYPE 64
#define PFCP_FQCSID_TYPE 65
#define PFCP_VolumeMeasurement_TYPE 66
#define PFCP_DurationMeasurement_TYPE 67
#define PFCP_ApplicationDetectionInformation_TYPE 68
#define PFCP_TimeOfFirstPacket_TYPE 69
#define PFCP_TimeOfLastPacket_TYPE 70
#define PFCP_QuotaHoldingTime_TYPE 71
#define PFCP_DroppedDLTrafficThreshold_TYPE 72
#define PFCP_VolumeQuota_TYPE 73
#define PFCP_TimeQuota_TYPE 74
#define PFCP_StartTime_TYPE 75
#define PFCP_EndTime_TYPE 76
#define PFCP_QueryURR_TYPE 77
#define PFCP_UsageReportPFCPSessionModificationResponse_TYPE 78
#define PFCP_UsageReportPFCPSessionDeletionResponse_TYPE 79
#define PFCP_UsageReportPFCPSessionReportRequest_TYPE 80
#define PFCP_URRID_TYPE 81
#define PFCP_LinkedURRID_TYPE 82
#define PFCP_DownlinkDataReport_TYPE 83
#define PFCP_OuterHeaderCreation_TYPE 84
#define PFCP_CreateBAR_TYPE 85
#define PFCP_UpdateBARPFCPSessionModificationRequest_TYPE 86
#define PFCP_RemoveBAR_TYPE 87
#define PFCP_BARID_TYPE 88
#define PFCP_CPFunctionFeatures_TYPE 89
#define PFCP_UsageInformation_TYPE 90
#define PFCP_ApplicationInstanceID_TYPE 91
#define PFCP_FlowInformation_TYPE 92
#define PFCP_UEIPAddress_TYPE 93
#define PFCP_PacketRate_TYPE 94
#define PFCP_OuterHeaderRemoval_TYPE 95
#define PFCP_RecoveryTimeStamp_TYPE 96
#define PFCP_DLFlowLevelMarking_TYPE 97
#define PFCP_HeaderEnrichment_TYPE 98
#define PFCP_ErrorIndicationReport_TYPE 99
#define PFCP_MeasurementInformation_TYPE 100
#define PFCP_NodeReportType_TYPE 101
#define PFCP_UserPlanePathFailureReport_TYPE 102
#define PFCP_RemoteGTPUPeer_TYPE 103
#define PFCP_URSEQN_TYPE 104
#define PFCP_UpdateDuplicatingParameters_TYPE 105
#define PFCP_ActivatePredefinedRules_TYPE 106
#define PFCP_DeactivatePredefinedRules_TYPE 107
#define PFCP_FARID_TYPE 108
#define PFCP_QERID_TYPE 109
#define PFCP_OCIFlags_TYPE 110
#define PFCP_PFCPAssociationReleaseRequest_TYPE 111
#define PFCP_GracefulReleasePeriod_TYPE 112
#define PFCP_PDNType_TYPE 113
#define PFCP_FailedRuleID_TYPE 114
#define PFCP_TimeQuotaMechanism_TYPE 115
#define PFCP_UserPlaneIPResourceInformation_TYPE 116
#define PFCP_UserPlaneInactivityTimer_TYPE 117
#define PFCP_AggregatedURRs_TYPE 118
#define PFCP_Multiplier_TYPE 119
#define PFCP_AggregatedURRID_TYPE 120
#define PFCP_SubsequentVolumeQuota_TYPE 121
#define PFCP_SubsequentTimeQuota_TYPE 122
#define PFCP_RQI_TYPE 123
#define PFCP_QFI_TYPE 124
#define PFCP_QueryURRReference_TYPE 125
#define PFCP_AdditionalUsageReportsInformation_TYPE 126
#define PFCP_CreateTrafficEndpoint_TYPE 127
#define PFCP_CreatedTrafficEndpoint_TYPE 128
#define PFCP_UpdateTrafficEndpoint_TYPE 129
#define PFCP_RemoveTrafficEndpoint_TYPE 130
#define PFCP_TrafficEndpointID_TYPE 131
#define PFCP_EthernetPacketFilter_TYPE 132
#define PFCP_MACAddress_TYPE 133
#define PFCP_CTAG_TYPE 134
#define PFCP_STAG_TYPE 135
#define PFCP_Ethertype_TYPE 136
#define PFCP_Proxying_TYPE 137
#define PFCP_EthernetFilterID_TYPE 138
#define PFCP_EthernetFilterProperties_TYPE 139
#define PFCP_SuggestedBufferingPacketsCount_TYPE 140
#define PFCP_UserID_TYPE 141
#define PFCP_EthernetPDUSessionInformation_TYPE 142
#define PFCP_EthernetTrafficInformation_TYPE 143
#define PFCP_MACAddressesDetected_TYPE 144
#define PFCP_MACAddressesRemoved_TYPE 145
#define PFCP_EthernetInactivityTimer_TYPE 146
#define PFCP_AdditionalMonitoringTime_TYPE 147
#define PFCP_EventInformation_TYPE 148
#define PFCP_EventReporting_TYPE 149
#define PFCP_EventID_TYPE 150
#define PFCP_EventThreshold_TYPE 151
#define PFCP_TraceInformation_TYPE 152
#define PFCP_FramedRoute_TYPE 153
#define PFCP_FramedRouting_TYPE 154
#define PFCP_FramedIPv6Route_TYPE 155

typedef TlvOctet Reserved;
typedef TlvOctet Cause;
typedef TlvOctet SourceInterface;
typedef TlvOctet FTEID;
typedef TlvOctet NetworkInstance;
typedef TlvOctet SDFFilter;
typedef TlvOctet ApplicationID;
typedef TlvOctet GateStatus;
typedef TlvOctet MBR;
typedef TlvOctet GBR;
typedef TlvOctet QERCorrelationID;
typedef TlvOctet Precedence;
typedef TlvOctet TransportLevelMarking;
typedef TlvOctet VolumeThreshold;
typedef TlvOctet TimeThreshold;
typedef TlvOctet MonitoringTime;
typedef TlvOctet SubsequentVolumeThreshold;
typedef TlvOctet SubsequentTimeThreshold;
typedef TlvOctet InactivityDetectionTime;
typedef TlvOctet ReportingTriggers;
typedef TlvOctet RedirectInformation;
typedef TlvOctet ReportType;
typedef TlvOctet OffendingIE;
typedef TlvOctet ForwardingPolicy;
typedef TlvOctet DestinationInterface;
typedef TlvOctet UPFunctionFeatures;
typedef TlvOctet ApplyAction;
typedef TlvOctet DownlinkDataServiceInformation;
typedef TlvOctet DownlinkDataNotificationDelay;
typedef TlvOctet DLBufferingDuration;
typedef TlvOctet DLBufferingSuggestedPacketCount;
typedef TlvOctet PFCPSMReqFlags;
typedef TlvOctet PFCPSRRspFlags;
typedef TlvOctet SequenceNumber;
typedef TlvOctet Metric;
typedef TlvOctet Timer;
typedef TlvOctet PacketDetectionRuleID;
typedef TlvOctet FSEID;
//typedef TlvOctet PFDContext;
typedef TlvOctet NodeID;
typedef TlvOctet PFDContents;
typedef TlvOctet MeasurementMethod;
typedef TlvOctet UsageReportTrigger;
typedef TlvOctet MeasurementPeriod;
typedef TlvOctet FQCSID;
typedef TlvOctet VolumeMeasurement;
typedef TlvOctet DurationMeasurement;
typedef TlvOctet TimeOfFirstPacket;
typedef TlvOctet TimeOfLastPacket;
typedef TlvOctet QuotaHoldingTime;
typedef TlvOctet DroppedDLTrafficThreshold;
typedef TlvOctet VolumeQuota;
typedef TlvOctet TimeQuota;
typedef TlvOctet StartTime;
typedef TlvOctet EndTime;
typedef TlvOctet URRID;
typedef TlvOctet LinkedURRID;
typedef TlvOctet OuterHeaderCreation;
typedef TlvOctet BARID;
typedef TlvOctet CPFunctionFeatures;
typedef TlvOctet UsageInformation;
typedef TlvOctet ApplicationInstanceID;
typedef TlvOctet FlowInformation;
typedef TlvOctet UEIPAddress;
typedef TlvOctet PacketRate;
typedef TlvOctet OuterHeaderRemoval;
typedef TlvOctet RecoveryTimeStamp;
typedef TlvOctet DLFlowLevelMarking;
typedef TlvOctet HeaderEnrichment;
typedef TlvOctet MeasurementInformation;
typedef TlvOctet NodeReportType;
//typedef TlvOctet UserPlanePathFailureReport;
typedef TlvOctet RemoteGTPUPeer;
typedef TlvOctet URSEQN;
typedef TlvOctet ActivatePredefinedRules;
typedef TlvOctet DeactivatePredefinedRules;
typedef TlvOctet FARID;
typedef TlvOctet QERID;
typedef TlvOctet OCIFlags;
typedef TlvOctet GracefulReleasePeriod;
typedef TlvOctet PDNType;
typedef TlvOctet FailedRuleID;
typedef TlvOctet TimeQuotaMechanism;
typedef TlvOctet UserPlaneIPResourceInformation;
typedef TlvOctet UserPlaneInactivityTimer;
typedef TlvOctet Multiplier;
typedef TlvOctet AggregatedURRID;
typedef TlvOctet SubsequentVolumeQuota;
typedef TlvOctet SubsequentTimeQuota;
typedef TlvOctet RQI;
typedef TlvOctet QFI;
typedef TlvOctet QueryURRReference;
typedef TlvOctet AdditionalUsageReportsInformation;
typedef TlvOctet TrafficEndpointID;
typedef TlvOctet MACAddress;
typedef TlvOctet CTAG;
typedef TlvOctet STAG;
typedef TlvOctet Ethertype;
typedef TlvOctet Proxying;
typedef TlvOctet EthernetFilterID;
typedef TlvOctet EthernetFilterProperties;
typedef TlvOctet SuggestedBufferingPacketsCount;
typedef TlvOctet UserID;
typedef TlvOctet EthernetPDUSessionInformation;
typedef TlvOctet MACAddressesDetected;
typedef TlvOctet MACAddressesRemoved;
typedef TlvOctet EthernetInactivityTimer;
typedef TlvOctet EventID;
typedef TlvOctet EventThreshold;
typedef TlvOctet TraceInformation;
typedef TlvOctet FramedRoute;
typedef TlvOctet FramedRouting;
typedef TlvOctet FramedIPv6Route;

typedef struct _HeartbeatRequest {
    unsigned long presence;
    RecoveryTimeStamp recoveryTimeStamp;
} __attribute__((packed)) HeartbeatRequest;

typedef struct _HeartbeatResponse {
    unsigned long presence;
    RecoveryTimeStamp recoveryTimeStamp;
} __attribute__((packed)) HeartbeatResponse;

typedef struct _PFDContext {
    unsigned long presence;
    PFDContents pFDContents;
} __attribute__((packed)) PFDContext;

typedef struct _PFCPPFDManagementResponse {
    unsigned long presence;
    Cause cause;
    OffendingIE offendingIE;
} __attribute__((packed)) PFCPPFDManagementResponse;

typedef struct _PFCPAssociationSetupRequest {
    unsigned long presence;
    NodeID nodeID;
    RecoveryTimeStamp recoveryTimeStamp;
    UPFunctionFeatures uPFunctionFeatures;
    CPFunctionFeatures cPFunctionFeatures;
    UserPlaneIPResourceInformation userPlaneIPResourceInformation;
} __attribute__((packed)) PFCPAssociationSetupRequest;

typedef struct _PFCPAssociationSetupResponse {
    unsigned long presence;
    NodeID nodeID;
    Cause cause;
    RecoveryTimeStamp recoveryTimeStamp;
    UPFunctionFeatures uPFunctionFeatures;
    CPFunctionFeatures cPFunctionFeatures;
    UserPlaneIPResourceInformation userPlaneIPResourceInformation;
} __attribute__((packed)) PFCPAssociationSetupResponse;

typedef struct _PFCPAssociationUpdateResponse {
    unsigned long presence;
    NodeID nodeID;
    Cause cause;
    UPFunctionFeatures uPFunctionFeatures;
    CPFunctionFeatures cPFunctionFeatures;
} __attribute__((packed)) PFCPAssociationUpdateResponse;

typedef struct _PFCPAssociationReleaseRequest {
    unsigned long presence;
    NodeID nodeID;
} __attribute__((packed)) PFCPAssociationReleaseRequest;

typedef struct _PFCPAssociationReleaseResponse {
    unsigned long presence;
    NodeID nodeID;
    Cause cause;
} __attribute__((packed)) PFCPAssociationReleaseResponse;

typedef struct _UserPlanePathFailureReport {
    unsigned long presence;
    RemoteGTPUPeer remoteGTPUPeer;
} __attribute__((packed)) UserPlanePathFailureReport;

typedef struct _PFCPNodeReportRequest {
    unsigned long presence;
    NodeID nodeID;
    NodeReportType nodeReportType;
    UserPlanePathFailureReport userPlanePathFailureReport;
} __attribute__((packed)) PFCPNodeReportRequest;

typedef struct _PFCPNodeReportResponse {
    unsigned long presence;
    NodeID nodeID;
    Cause cause;
    OffendingIE offendingIE;
} __attribute__((packed)) PFCPNodeReportResponse;

typedef struct _PFCPSessionSetDeletionRequest {
    unsigned long presence;
    NodeID nodeID;
    FQCSID sGWCFQCSID;
    FQCSID pGWCFQCSID;
    FQCSID sGWUFQCSID;
    FQCSID pGWUFQCSID;
    FQCSID tWANFQCSID;
    FQCSID ePDGFQCSID;
    FQCSID mMEFQCSID;
} __attribute__((packed)) PFCPSessionSetDeletionRequest;

typedef struct _PFCPSessionSetDeletionResponse {
    unsigned long presence;
    NodeID nodeID;
    Cause cause;
    OffendingIE offendingIE;
} __attribute__((packed)) PFCPSessionSetDeletionResponse;

typedef struct _EthernetPacketFilter {
    unsigned long presence;
    EthernetFilterID ethernetFilterID;
    EthernetFilterProperties ethernetFilterProperties;
    MACAddress mACAddress;
    Ethertype ethertype;
    CTAG cTAG;
    STAG sTAG;
    SDFFilter sDFFilter;
} __attribute__((packed)) EthernetPacketFilter;

typedef struct _ForwardingParameters {
    unsigned long presence;
    DestinationInterface destinationInterface;
    NetworkInstance networkInstance;
    RedirectInformation redirectInformation;
    OuterHeaderCreation outerHeaderCreation;
    TransportLevelMarking transportLevelMarking;
    ForwardingPolicy forwardingPolicy;
    HeaderEnrichment headerEnrichment;
    TrafficEndpointID linkedTrafficEndpointID;
    Proxying proxying;
} __attribute__((packed)) ForwardingParameters;

typedef struct _DuplicatingParameters {
    unsigned long presence;
    DestinationInterface destinationInterface;
    OuterHeaderCreation outerHeaderCreation;
    TransportLevelMarking transportLevelMarking;
    ForwardingPolicy forwardingPolicy;
} __attribute__((packed)) DuplicatingParameters;

typedef struct _AggregatedURRs {
    unsigned long presence;
    AggregatedURRID aggregatedURRID;
    Multiplier multiplier;
} __attribute__((packed)) AggregatedURRs;

typedef struct _AdditionalMonitoringTime {
    unsigned long presence;
    MonitoringTime monitoringTime;
    SubsequentVolumeThreshold subsequentVolumeThreshold;
    SubsequentTimeThreshold subsequentTimeThreshold;
    SubsequentVolumeQuota subsequentVolumeQuota;
    SubsequentTimeQuota subsequentTimeQuota;
} __attribute__((packed)) AdditionalMonitoringTime;

typedef struct _EventInformation {
    unsigned long presence;
    EventID eventID;
    EventThreshold eventThreshold;
} __attribute__((packed)) EventInformation;

typedef struct _CreateQER {
    unsigned long presence;
    QERID qERID;
    QERCorrelationID qERCorrelationID;
    GateStatus gateStatus;
    MBR maximumBitrate;
    GBR guaranteedBitrate;
    PacketRate packetRate;
    DLFlowLevelMarking dLFlowLevelMarking;
    QFI qoSFlowIdentifier;
    RQI reflectiveQoS;
} __attribute__((packed)) CreateQER;

typedef struct _CreateBAR {
    unsigned long presence;
    BARID bARID;
    DownlinkDataNotificationDelay downlinkDataNotificationDelay;
    SuggestedBufferingPacketsCount suggestedBufferingPacketsCount;
} __attribute__((packed)) CreateBAR;

typedef struct _CreateTrafficEndpoint {
    unsigned long presence;
    TrafficEndpointID trafficEndpointID;
    FTEID localFTEID;
    NetworkInstance networkInstance;
    UEIPAddress uEIPAddress;
    EthernetPDUSessionInformation ethernetPDUSessionInformation;
    FramedRoute framedRoute;
    FramedRouting framedRouting;
    FramedIPv6Route framedIPv6Route;
} __attribute__((packed)) CreateTrafficEndpoint;

typedef struct _LoadControlInformation {
    unsigned long presence;
    SequenceNumber loadControlSequenceNumber;
    Metric loadMetric;
} __attribute__((packed)) LoadControlInformation;

typedef struct _OverloadControlInformation {
    unsigned long presence;
    SequenceNumber overloadControlSequenceNumber;
    Metric overloadReductionMetric;
    Timer periodOfValidity;
    OCIFlags overloadControlInformationFlags;
} __attribute__((packed)) OverloadControlInformation;

typedef struct _CreatedTrafficEndpoint {
    unsigned long presence;
    TrafficEndpointID trafficEndpointID;
    FTEID localFTEID;
} __attribute__((packed)) CreatedTrafficEndpoint;

typedef struct _UpdateForwardingParameters {
    unsigned long presence;
    DestinationInterface destinationInterface;
    NetworkInstance networkInstance;
    RedirectInformation redirectInformation;
    OuterHeaderCreation outerHeaderCreation;
    TransportLevelMarking transportLevelMarking;
    ForwardingPolicy forwardingPolicy;
    HeaderEnrichment headerEnrichment;
    PFCPSMReqFlags pFCPSMReqFlags;
    TrafficEndpointID linkedTrafficEndpointID;
} __attribute__((packed)) UpdateForwardingParameters;

typedef struct _UpdateDuplicatingParameters {
    unsigned long presence;
    DestinationInterface destinationInterface;
    OuterHeaderCreation outerHeaderCreation;
    TransportLevelMarking transportLevelMarking;
    ForwardingPolicy forwardingPolicy;
} __attribute__((packed)) UpdateDuplicatingParameters;

typedef struct _UpdateURR {
    unsigned long presence;
    URRID uRRID;
    MeasurementMethod measurementMethod;
    ReportingTriggers reportingTriggers;
    MeasurementPeriod measurementPeriod;
    VolumeThreshold volumeThreshold;
    VolumeQuota volumeQuota;
    TimeThreshold timeThreshold;
    TimeQuota timeQuota;
    QuotaHoldingTime quotaHoldingTime;
    DroppedDLTrafficThreshold droppedDLTrafficThreshold;
    MonitoringTime monitoringTime;
    EventInformation eventInformation;
    SubsequentVolumeThreshold subsequentVolumeThreshold;
    SubsequentTimeThreshold subsequentTimeThreshold;
    SubsequentVolumeQuota subsequentVolumeQuota;
    SubsequentTimeQuota subsequentTimeQuota;
    InactivityDetectionTime inactivityDetectionTime;
    LinkedURRID linkedURRID;
    MeasurementInformation measurementInformation;
    TimeQuotaMechanism timeQuotaMechanism;
    AggregatedURRs aggregatedURRs;
    FARID fARIDForQuotaAction;
    EthernetInactivityTimer ethernetInactivityTimer;
    AdditionalMonitoringTime additionalMonitoringTime;
} __attribute__((packed)) UpdateURR;

typedef struct _UpdateQER {
    unsigned long presence;
    QERID qERID;
    QERCorrelationID qERCorrelationID;
    GateStatus gateStatus;
    MBR maximumBitrate;
    GBR guaranteedBitrate;
    PacketRate packetRate;
    DLFlowLevelMarking dLFlowLevelMarking;
    QFI qoSFlowIdentifier;
    RQI reflectiveQoS;
} __attribute__((packed)) UpdateQER;

typedef struct _RemoveFAR {
    unsigned long presence;
    FARID fARID;
} __attribute__((packed)) RemoveFAR;

typedef struct _RemoveURR {
    unsigned long presence;
    URRID uRRID;
} __attribute__((packed)) RemoveURR;

typedef struct _RemoveQER {
    unsigned long presence;
    QERID qERID;
} __attribute__((packed)) RemoveQER;

typedef struct _QueryURR {
    unsigned long presence;
    URRID uRRID;
} __attribute__((packed)) QueryURR;

typedef struct _UpdateBARPFCPSessionModificationRequest {
    unsigned long presence;
    BARID bARID;
    DownlinkDataNotificationDelay downlinkDataNotificationDelay;
    SuggestedBufferingPacketsCount suggestedBufferingPacketsCount;
} __attribute__((packed)) UpdateBARPFCPSessionModificationRequest;

typedef struct _RemoveBAR {
    unsigned long presence;
    BARID bARID;
} __attribute__((packed)) RemoveBAR;

typedef struct _UpdateTrafficEndpoint {
    unsigned long presence;
    TrafficEndpointID trafficEndpointID;
    FTEID localFTEID;
    NetworkInstance networkInstance;
    UEIPAddress uEIPAddress;
    FramedRoute framedRoute;
    FramedRouting framedRouting;
    FramedIPv6Route framedIPv6Route;
} __attribute__((packed)) UpdateTrafficEndpoint;

typedef struct _RemoveTrafficEndpoint {
    unsigned long presence;
    TrafficEndpointID trafficEndpointID;
} __attribute__((packed)) RemoveTrafficEndpoint;

typedef struct _ApplicationDetectionInformation {
    unsigned long presence;
    ApplicationID applicationID;
    ApplicationInstanceID applicationInstanceID;
    FlowInformation flowInformation;
} __attribute__((packed)) ApplicationDetectionInformation;

typedef struct _EventReporting {
    unsigned long presence;
    EventID eventID;
} __attribute__((packed)) EventReporting;

typedef struct _EthernetTrafficInformation {
    unsigned long presence;
    MACAddressesDetected mACAddressesDetected;
    MACAddressesRemoved mACAddressesRemoved;
} __attribute__((packed)) EthernetTrafficInformation;

typedef struct _ErrorIndicationReport {
    unsigned long presence;
    FTEID remoteFTEID;
} __attribute__((packed)) ErrorIndicationReport;

typedef struct _UpdateBARPFCPSessionReportResponse {
    unsigned long presence;
    BARID bARID;
    DownlinkDataNotificationDelay downlinkDataNotificationDelay;
    DLBufferingDuration dLBufferingDuration;
    DLBufferingSuggestedPacketCount dLBufferingSuggestedPacketCount;
    SuggestedBufferingPacketsCount suggestedBufferingPacketsCount;
} __attribute__((packed)) UpdateBARPFCPSessionReportResponse;

typedef struct _ApplicationIDsPFDs {
    unsigned long presence;
    ApplicationID applicationID;
    PFDContext pFD;
} __attribute__((packed)) ApplicationIDsPFDs;

typedef struct _PFCPAssociationUpdateRequest {
    unsigned long presence;
    NodeID nodeID;
    UPFunctionFeatures uPFunctionFeatures;
    CPFunctionFeatures cPFunctionFeatures;
    PFCPAssociationReleaseRequest pFCPAssociationReleaseRequest;
    GracefulReleasePeriod gracefulReleasePeriod;
    UserPlaneIPResourceInformation userPlaneIPResourceInformation;
} __attribute__((packed)) PFCPAssociationUpdateRequest;

typedef struct _PDI {
    unsigned long presence;
    SourceInterface sourceInterface;
    FTEID localFTEID;
    NetworkInstance networkInstance;
    UEIPAddress uEIPAddress;
    TrafficEndpointID trafficEndpointID;
    SDFFilter sDFFilter;
    ApplicationID applicationID;
    EthernetPDUSessionInformation ethernetPDUSessionInformation;
    EthernetPacketFilter ethernetPacketFilter;
    QFI qFI;
    FramedRoute framedRoute;
    FramedRouting framedRouting;
    FramedIPv6Route framedIPv6Route;
} __attribute__((packed)) PDI;

typedef struct _CreateFAR {
    unsigned long presence;
    FARID fARID;
    ApplyAction applyAction;
    ForwardingParameters forwardingParameters;
    DuplicatingParameters duplicatingParameters;
    BARID bARID;
} __attribute__((packed)) CreateFAR;

typedef struct _CreateURR {
    unsigned long presence;
    URRID uRRID;
    MeasurementMethod measurementMethod;
    ReportingTriggers reportingTriggers;
    MeasurementPeriod measurementPeriod;
    VolumeThreshold volumeThreshold;
    VolumeQuota volumeQuota;
    TimeThreshold timeThreshold;
    TimeQuota timeQuota;
    QuotaHoldingTime quotaHoldingTime;
    DroppedDLTrafficThreshold droppedDLTrafficThreshold;
    MonitoringTime monitoringTime;
    EventInformation eventInformation;
    SubsequentVolumeThreshold subsequentVolumeThreshold;
    SubsequentTimeThreshold subsequentTimeThreshold;
    SubsequentVolumeQuota subsequentVolumeQuota;
    SubsequentTimeQuota subsequentTimeQuota;
    InactivityDetectionTime inactivityDetectionTime;
    LinkedURRID linkedURRID;
    MeasurementInformation measurementInformation;
    TimeQuotaMechanism timeQuotaMechanism;
    AggregatedURRs aggregatedURRs;
    FARID fARIDForQuotaAction;
    EthernetInactivityTimer ethernetInactivityTimer;
    AdditionalMonitoringTime additionalMonitoringTime;
} __attribute__((packed)) CreateURR;

typedef struct _CreatedPDR {
    unsigned long presence;
    PacketDetectionRuleID pDRID;
    FTEID localFTEID;
} __attribute__((packed)) CreatedPDR;

typedef struct _UpdatePDR {
    unsigned long presence;
    PacketDetectionRuleID pDRID;
    OuterHeaderRemoval outerHeaderRemoval;
    Precedence precedence;
    PDI pDI;
    FARID fARID;
    URRID uRRID;
    QERID qERID;
    ActivatePredefinedRules activatePredefinedRules;
    DeactivatePredefinedRules deactivatePredefinedRules;
} __attribute__((packed)) UpdatePDR;

typedef struct _UpdateFAR {
    unsigned long presence;
    FARID fARID;
    ApplyAction applyAction;
    UpdateForwardingParameters updateForwardingParameters;
    UpdateDuplicatingParameters updateDuplicatingParameters;
    BARID bARID;
} __attribute__((packed)) UpdateFAR;

typedef struct _RemovePDR {
    unsigned long presence;
    PacketDetectionRuleID pDRID;
} __attribute__((packed)) RemovePDR;

typedef struct _UsageReportPFCPSessionModificationResponse {
    unsigned long presence;
    URRID uRRID;
    URSEQN uRSEQN;
    UsageReportTrigger usageReportTrigger;
    StartTime startTime;
    EndTime endTime;
    VolumeMeasurement volumeMeasurement;
    DurationMeasurement durationMeasurement;
    TimeOfFirstPacket timeOfFirstPacket;
    TimeOfLastPacket timeOfLastPacket;
    UsageInformation usageInformation;
    QueryURRReference queryURRReference;
    EthernetTrafficInformation ethernetTrafficInformation;
} __attribute__((packed)) UsageReportPFCPSessionModificationResponse;

typedef struct _PFCPSessionDeletionRequest {
    unsigned long presence;
} __attribute__((packed)) PFCPSessionDeletionRequest;

typedef struct _UsageReportPFCPSessionDeletionResponse {
    unsigned long presence;
    URRID uRRID;
    URSEQN uRSEQN;
    UsageReportTrigger usageReportTrigger;
    StartTime startTime;
    EndTime endTime;
    VolumeMeasurement volumeMeasurement;
    DurationMeasurement durationMeasurement;
    TimeOfFirstPacket timeOfFirstPacket;
    TimeOfLastPacket timeOfLastPacket;
    UsageInformation usageInformation;
    EthernetTrafficInformation ethernetTrafficInformation;
} __attribute__((packed)) UsageReportPFCPSessionDeletionResponse;

typedef struct _DownlinkDataReport {
    unsigned long presence;
    PacketDetectionRuleID pDRID;
    DownlinkDataServiceInformation downlinkDataServiceInformation;
} __attribute__((packed)) DownlinkDataReport;

typedef struct _UsageReportPFCPSessionReportRequest {
    unsigned long presence;
    URRID uRRID;
    URSEQN uRSEQN;
    UsageReportTrigger usageReportTrigger;
    StartTime startTime;
    EndTime endTime;
    VolumeMeasurement volumeMeasurement;
    DurationMeasurement durationMeasurement;
    ApplicationDetectionInformation applicationDetectionInformation;
    UEIPAddress uEIPAddress;
    NetworkInstance networkInstance;
    TimeOfFirstPacket timeOfFirstPacket;
    TimeOfLastPacket timeOfLastPacket;
    UsageInformation usageInformation;
    QueryURRReference queryURRReference;
    EventReporting eventReporting;
    EthernetTrafficInformation ethernetTrafficInformation;
} __attribute__((packed)) UsageReportPFCPSessionReportRequest;

typedef struct _PFCPSessionReportResponse {
    unsigned long presence;
    Cause cause;
    OffendingIE offendingIE;
    UpdateBARPFCPSessionReportResponse updateBAR;
    PFCPSRRspFlags sxSRRspFlags;
} __attribute__((packed)) PFCPSessionReportResponse;

typedef struct _PFCPPFDManagementRequest {
    unsigned long presence;
    ApplicationIDsPFDs applicationIDsPFDs;
} __attribute__((packed)) PFCPPFDManagementRequest;

typedef struct _CreatePDR {
    unsigned long presence;
    PacketDetectionRuleID pDRID;
    Precedence precedence;
    PDI pDI;
    OuterHeaderRemoval outerHeaderRemoval;
    FARID fARID;
    URRID uRRID;
    QERID qERID;
    ActivatePredefinedRules activatePredefinedRules;
} __attribute__((packed)) CreatePDR;

typedef struct _PFCPSessionEstablishmentResponse {
    unsigned long presence;
    NodeID nodeID;
    Cause cause;
    OffendingIE offendingIE;
    FSEID uPFSEID;
    CreatedPDR createdPDR;
    LoadControlInformation loadControlInformation;
    OverloadControlInformation overloadControlInformation;
    FQCSID sGWUFQCSID;
    FQCSID pGWUFQCSID;
    FailedRuleID failedRuleID;
    CreatedTrafficEndpoint createdTrafficEndpoint;
} __attribute__((packed)) PFCPSessionEstablishmentResponse;

typedef struct _PFCPSessionModificationRequest {
    unsigned long presence;
    FSEID cPFSEID;
    RemovePDR removePDR;
    RemoveFAR removeFAR;
    RemoveURR removeURR;
    RemoveQER removeQER;
    RemoveBAR removeBAR;
    RemoveTrafficEndpoint removeTrafficEndpoint;
    CreatePDR createPDR[2];
    CreateFAR createFAR[2];
    CreateURR createURR;
    CreateQER createQER;
    CreateBAR createBAR;
    CreateTrafficEndpoint createTrafficEndpoint;
    UpdatePDR updatePDR;
    UpdateFAR updateFAR;
    UpdateURR updateURR;
    UpdateQER updateQER;
    UpdateBARPFCPSessionModificationRequest updateBAR;
    UpdateTrafficEndpoint updateTrafficEndpoint;
    PFCPSMReqFlags pFCPSMReqFlags;
    QueryURR queryURR;
    FQCSID pGWCFQCSID;
    FQCSID sGWCFQCSID;
    FQCSID mMEFQCSID;
    FQCSID ePDGFQCSID;
    FQCSID tWANFQCSID;
    UserPlaneInactivityTimer userPlaneInactivityTimer;
    QueryURRReference queryURRReference;
    TraceInformation traceInformation;
} __attribute__((packed)) PFCPSessionModificationRequest;

typedef struct _PFCPSessionModificationResponse {
    unsigned long presence;
    Cause cause;
    OffendingIE offendingIE;
    CreatedPDR createdPDR;
    LoadControlInformation loadControlInformation;
    OverloadControlInformation overloadControlInformation;
    UsageReportPFCPSessionModificationResponse usageReport;
    FailedRuleID failedRuleID;
    AdditionalUsageReportsInformation additionalUsageReportsInformation;
    CreatedTrafficEndpoint createdUpdatedTrafficEndpoint;
} __attribute__((packed)) PFCPSessionModificationResponse;

typedef struct _PFCPSessionDeletionResponse {
    unsigned long presence;
    Cause cause;
    OffendingIE offendingIE;
    LoadControlInformation loadControlInformation;
    OverloadControlInformation overloadControlInformation;
    UsageReportPFCPSessionDeletionResponse usageReport;
} __attribute__((packed)) PFCPSessionDeletionResponse;

typedef struct _PFCPSessionReportRequest {
    unsigned long presence;
    ReportType reportType;
    DownlinkDataReport downlinkDataReport;
    UsageReportPFCPSessionReportRequest usageReport;
    ErrorIndicationReport errorIndicationReport;
    LoadControlInformation loadControlInformation;
    OverloadControlInformation overloadControlInformation;
    AdditionalUsageReportsInformation additionalUsageReportsInformation;
} __attribute__((packed)) PFCPSessionReportRequest;

typedef struct _PFCPSessionEstablishmentRequest {
    unsigned long presence;
    NodeID nodeID;
    FSEID cPFSEID;
    CreatePDR createPDR[2];
    CreateFAR createFAR[2];
    CreateURR createURR;
    CreateQER createQER;
    CreateBAR createBAR;
    CreateTrafficEndpoint createTrafficEndpoint;
    PDNType pDNType;
    FQCSID sGWCFQCSID;
    FQCSID mMEFQCSID;
    FQCSID pGWCFQCSID;
    FQCSID ePDGFQCSID;
    FQCSID tWANFQCSID;
    UserPlaneInactivityTimer userPlaneInactivityTimer;
    UserID userID;
    TraceInformation traceInformation;
} __attribute__((packed)) PFCPSessionEstablishmentRequest;


typedef struct _PfcpMessage {
    PfcpHeader header;
    union {
        HeartbeatRequest heartbeatRequest;
        HeartbeatResponse heartbeatResponse;
        PFCPPFDManagementRequest pFCPPFDManagementRequest;
        PFCPPFDManagementResponse pFCPPFDManagementResponse;
        PFCPAssociationSetupRequest pFCPAssociationSetupRequest;
        PFCPAssociationSetupResponse pFCPAssociationSetupResponse;
        PFCPAssociationUpdateRequest pFCPAssociationUpdateRequest;
        PFCPAssociationUpdateResponse pFCPAssociationUpdateResponse;
        PFCPAssociationReleaseRequest pFCPAssociationReleaseRequest;
        PFCPAssociationReleaseResponse pFCPAssociationReleaseResponse;
        PFCPNodeReportRequest pFCPNodeReportRequest;
        PFCPNodeReportResponse pFCPNodeReportResponse;
        PFCPSessionSetDeletionRequest pFCPSessionSetDeletionRequest;
        PFCPSessionSetDeletionResponse pFCPSessionSetDeletionResponse;
        PFCPSessionEstablishmentRequest pFCPSessionEstablishmentRequest;
        PFCPSessionEstablishmentResponse pFCPSessionEstablishmentResponse;
        PFCPSessionModificationRequest pFCPSessionModificationRequest;
        PFCPSessionModificationResponse pFCPSessionModificationResponse;
        PFCPSessionDeletionRequest pFCPSessionDeletionRequest;
        PFCPSessionDeletionResponse pFCPSessionDeletionResponse;
        PFCPSessionReportRequest pFCPSessionReportRequest;
        PFCPSessionReportResponse pFCPSessionReportResponse;
    };
} PfcpMessage;

Status PfcpParseMessage(PfcpMessage *pfcpMessage, Bufblk *buf);

Status PfcpBuildMessage(Bufblk **bufBlkPtr, PfcpMessage *pfcpMessage);

  Status PfcpStructFree(PfcpMessage *pfcpMessage);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __PFCP_MESSAGE_H__ */

