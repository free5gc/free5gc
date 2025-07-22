package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type LocationReportingChargingInformation struct { /* Set Type */
	LocationReportingMessagetype LocationReportingMessageType       `ber:"tagNum:0"`
	UserIdentifier               *InvolvedParty                     `ber:"tagNum:1,optional"`
	UserEquipmentInfo            *SubscriberEquipmentNumber         `ber:"tagNum:2,optional"`
	SUPIunauthenticatedFlag      *asn.NULL                          `ber:"tagNum:3,optional"`
	UserRoamerInOut              *RoamerInOut                       `ber:"tagNum:4,optional"`
	UserLocationInformation      *UserLocationInformation           `ber:"tagNum:5,optional"`
	UserLocationInfoTime         *TimeStamp                         `ber:"tagNum:6,optional"`
	UETimeZone                   *MSTimeZone                        `ber:"tagNum:7,optional"`
	PresenceReportingAreaInfo    *PresenceReportingAreaInfo         `ber:"tagNum:8,optional"`
	RATType                      *RATType                           `ber:"tagNum:9,optional"`
	PSCellInformation            *PSCellInformation                 `ber:"tagNum:10,optional"`
	UserLocationInformationASN1  *UserLocationInformationStructured `ber:"tagNum:11,optional"`
}
