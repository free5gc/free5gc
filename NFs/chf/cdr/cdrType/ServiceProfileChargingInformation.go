package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type ServiceProfileChargingInformation struct { /* Set Type */
	ServiceProfileIdentifier *asn.OctetString `ber:"tagNum:0,optional"`
	/* Sequence of = 35, FULL Name = struct ServiceProfileChargingInformation__sNSSAIList */
	/* SingleNSSAI */
	SNSSAIList                    []SingleNSSAI                  `ber:"tagNum:1,optional"`
	SST                           *SliceServiceType              `ber:"tagNum:2,optional"`
	Latency                       *int64                         `ber:"tagNum:3,optional"`
	Availability                  *int64                         `ber:"tagNum:4,optional"`
	ResourceSharingLevel          *SharingLevel                  `ber:"tagNum:5,optional"`
	Jitter                        *int64                         `ber:"tagNum:6,optional"`
	Reliability                   *asn.OctetString               `ber:"tagNum:7,optional"`
	MaxNumberofUEs                *int64                         `ber:"tagNum:8,optional"`
	CoverageArea                  *asn.OctetString               `ber:"tagNum:9,optional"`
	UEMobilityLevel               *MobilityLevel                 `ber:"tagNum:10,optional"`
	DelayToleranceIndicator       *DelayToleranceIndicator       `ber:"tagNum:11,optional"`
	DLThroughtputPerSlice         *Throughput                    `ber:"tagNum:12,optional"`
	DLThroughtputPerUE            *Throughput                    `ber:"tagNum:13,optional"`
	ULThroughtputPerSlice         *Throughput                    `ber:"tagNum:14,optional"`
	ULThroughtputPerUE            *Throughput                    `ber:"tagNum:15,optional"`
	MaxNumberofPDUsessions        *int64                         `ber:"tagNum:16,optional"`
	KPIsMonitoringList            *asn.OctetString               `ber:"tagNum:17,optional"`
	SupportedAccessTechnology     *int64                         `ber:"tagNum:18,optional"`
	V2XCommunicationMode          *V2XCommunicationModeIndicator `ber:"tagNum:19,optional"`
	AddServiceProfileChargingInfo *asn.OctetString               `ber:"tagNum:100,optional"`
}
