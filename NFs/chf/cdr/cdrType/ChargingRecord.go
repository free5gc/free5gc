package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type ChargingRecord struct { /* Set Type */
	RecordType                   RecordType                 `ber:"tagNum:0"`
	RecordingNetworkFunctionID   NetworkFunctionName        `ber:"tagNum:1"`
	SubscriberIdentifier         *SubscriptionID            `ber:"tagNum:2,optional"`
	NFunctionConsumerInformation NetworkFunctionInformation `ber:"tagNum:3"`
	/* Sequence of = 35, FULL Name = struct ChargingRecord__triggers */
	/* Trigger */
	Triggers []Trigger `ber:"tagNum:4,optional"`
	/* Sequence of = 35, FULL Name = struct ChargingRecord__listOfMultipleUnitUsage */
	/* MultipleUnitUsage */
	ListOfMultipleUnitUsage              []MultipleUnitUsage                   `ber:"tagNum:5,optional"`
	RecordOpeningTime                    TimeStamp                             `ber:"tagNum:6"`
	Duration                             CallDuration                          `ber:"tagNum:7"`
	RecordSequenceNumber                 *int64                                `ber:"tagNum:8,optional"`
	CauseForRecClosing                   CauseForRecClosing                    `ber:"tagNum:9"`
	Diagnostics                          *Diagnostics                          `ber:"tagNum:10,optional"`
	LocalRecordSequenceNumber            *LocalSequenceNumber                  `ber:"tagNum:11,optional"`
	RecordExtensions                     *ManagementExtensions                 `ber:"tagNum:12,optional"`
	PDUSessionChargingInformation        *PDUSessionChargingInformation        `ber:"tagNum:13,optional"`
	RoamingQBCInformation                *RoamingQBCInformation                `ber:"tagNum:14,optional"`
	SMSChargingInformation               *SMSChargingInformation               `ber:"tagNum:15,optional"`
	ChargingSessionIdentifier            *ChargingSessionIdentifier            `ber:"tagNum:16,optional"`
	ServiceSpecificationInformation      *asn.OctetString                      `ber:"tagNum:17,optional"`
	ExposureFunctionAPIInformation       *ExposureFunctionAPIInformation       `ber:"tagNum:18,optional"`
	RegistrationChargingInformation      *RegistrationChargingInformation      `ber:"tagNum:19,optional"`
	N2ConnectionChargingInformation      *N2ConnectionChargingInformation      `ber:"tagNum:20,optional"`
	LocationReportingChargingInformation *LocationReportingChargingInformation `ber:"tagNum:21,optional"`
	IncompleteCDRIndication              *IncompleteCDRIndication              `ber:"tagNum:22,optional"`
	TenantIdentifier                     *TenantIdentifier                     `ber:"tagNum:23,optional"`
	MnSConsumerIdentifier                *MnSConsumerIdentifier                `ber:"tagNum:24,optional"`
	NSMChargingInformation               *NSMChargingInformation               `ber:"tagNum:25,optional"`
	NSPAChargingInformation              *NSPAChargingInformation              `ber:"tagNum:26,optional"`
	ChargingID                           *ChargingID                           `ber:"tagNum:27,optional"`
}
