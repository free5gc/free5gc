package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type N2ConnectionChargingInformation struct { /* Set Type */
	N2ConnectionMessageType N2ConnectionMessageType    `ber:"tagNum:0"`
	UserIdentifier          *InvolvedParty             `ber:"tagNum:1,optional"`
	UserEquipmentInfo       *SubscriberEquipmentNumber `ber:"tagNum:2,optional"`
	SUPIunauthenticatedFlag *asn.NULL                  `ber:"tagNum:3,optional"`
	UserRoamerInOut         *RoamerInOut               `ber:"tagNum:4,optional"`
	UserLocationInformation *UserLocationInformation   `ber:"tagNum:5,optional"`
	UserLocationInfoTime    *TimeStamp                 `ber:"tagNum:6,optional"`
	UETimeZone              *MSTimeZone                `ber:"tagNum:7,optional"`
	RATType                 *RATType                   `ber:"tagNum:8,optional"`
	RanUeNgapId             *RanUeNgapId               `ber:"tagNum:9,optional"`
	RanNodeId               *GlobalRanNodeId           `ber:"tagNum:10,optional"`
	/* Sequence of = 35, FULL Name = struct N2ConnectionChargingInformation__restrictedRatList */
	/* RATType */
	RestrictedRatList []RATType `ber:"tagNum:11,optional"`
	/* Sequence of = 35, FULL Name = struct N2ConnectionChargingInformation__forbiddenAreaList */
	/* Area */
	ForbiddenAreaList      []Area                  `ber:"tagNum:12,optional"`
	ServiceAreaRestriction *ServiceAreaRestriction `ber:"tagNum:13,optional"`
	/* Sequence of = 35, FULL Name = struct N2ConnectionChargingInformation__restrictedCnList */
	/* CoreNetworkType */
	RestrictedCnList []CoreNetworkType `ber:"tagNum:14,optional"`
	/* Sequence of = 35, FULL Name = struct N2ConnectionChargingInformation__allowedNSSAI */
	/* SingleNSSAI */
	AllowedNSSAI                []SingleNSSAI                      `ber:"tagNum:15,optional"`
	RrcEstablishmentCause       *RrcEstablishmentCause             `ber:"tagNum:16,optional"`
	PSCellInformation           *PSCellInformation                 `ber:"tagNum:17,optional"`
	AmfUeNgapId                 *AmfUeNgapId                       `ber:"tagNum:18,optional"`
	UserLocationInformationASN1 *UserLocationInformationStructured `ber:"tagNum:19,optional"`
}
