package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type RegistrationChargingInformation struct { /* Set Type */
	RegistrationMessagetype RegistrationMessageType    `ber:"tagNum:0"`
	UserIdentifier          *InvolvedParty             `ber:"tagNum:1,optional"`
	UserEquipmentInfo       *SubscriberEquipmentNumber `ber:"tagNum:2,optional"`
	SUPIunauthenticatedFlag *asn.NULL                  `ber:"tagNum:3,optional"`
	UserRoamerInOut         *RoamerInOut               `ber:"tagNum:4,optional"`
	UserLocationInformation *UserLocationInformation   `ber:"tagNum:5,optional"`
	UserLocationInfoTime    *TimeStamp                 `ber:"tagNum:6,optional"`
	UETimeZone              *MSTimeZone                `ber:"tagNum:7,optional"`
	RATType                 *RATType                   `ber:"tagNum:8,optional"`
	MICOModeIndication      *MICOModeIndication        `ber:"tagNum:9,optional"`
	SmsIndication           *SmsIndication             `ber:"tagNum:10,optional"`
	/* Sequence of = 35, FULL Name = struct RegistrationChargingInformation__taiList */
	/* TAI */
	TaiList                []TAI                   `ber:"tagNum:11,optional"`
	ServiceAreaRestriction *ServiceAreaRestriction `ber:"tagNum:12,optional"`
	/* Sequence of = 35, FULL Name = struct RegistrationChargingInformation__requestedNSSAI */
	/* SingleNSSAI */
	RequestedNSSAI []SingleNSSAI `ber:"tagNum:13,optional"`
	/* Sequence of = 35, FULL Name = struct RegistrationChargingInformation__allowedNSSAI */
	/* SingleNSSAI */
	AllowedNSSAI []SingleNSSAI `ber:"tagNum:14,optional"`
	/* Sequence of = 35, FULL Name = struct RegistrationChargingInformation__rejectedNSSAI */
	/* SingleNSSAI */
	RejectedNSSAI     []SingleNSSAI      `ber:"tagNum:15,optional"`
	PSCellInformation *PSCellInformation `ber:"tagNum:16,optional"`
	FiveGMMCapability *FiveGMMCapability `ber:"tagNum:17,optional"`
	/* Sequence of = 35, FULL Name = struct RegistrationChargingInformation__nSSAIMapList */
	/* NSSAIMap */
	NSSAIMapList                []NSSAIMap                         `ber:"tagNum:18,optional"`
	AmfUeNgapId                 *AmfUeNgapId                       `ber:"tagNum:19,optional"`
	RanUeNgapId                 *RanUeNgapId                       `ber:"tagNum:20,optional"`
	RanNodeId                   *GlobalRanNodeId                   `ber:"tagNum:21,optional"`
	UserLocationInformationASN1 *UserLocationInformationStructured `ber:"tagNum:22,optional"`
}
