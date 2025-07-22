package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type SMSChargingInformation struct { /* Set Type */
	OriginatorInfo *OriginatorInfo `ber:"tagNum:1,optional"`
	/* Sequence of = 35, FULL Name = struct SMSChargingInformation__recipientInfos */
	/* RecipientInfo */
	RecipientInfos              []RecipientInfo                    `ber:"tagNum:2,optional"`
	UserEquipmentInfo           *SubscriberEquipmentNumber         `ber:"tagNum:3,optional"`
	UserLocationInformation     *UserLocationInformation           `ber:"tagNum:4,optional"`
	UETimeZone                  *MSTimeZone                        `ber:"tagNum:5,optional"`
	RATType                     *RATType                           `ber:"tagNum:6,optional"`
	SMSCAddress                 *AddressString                     `ber:"tagNum:7,optional"`
	Eventtimestamp              TimeStamp                          `ber:"tagNum:8"`
	SMDataCodingScheme          *int64                             `ber:"tagNum:20,optional"`
	SMMessageType               *SMMessageType                     `ber:"tagNum:21,optional"`
	SMReplyPathRequested        *SMReplyPathRequested              `ber:"tagNum:22,optional"`
	SMUserDataHeader            *asn.OctetString                   `ber:"tagNum:23,optional"`
	SMSStatus                   *SMSStatus                         `ber:"tagNum:24,optional"`
	SMDischargeTime             *TimeStamp                         `ber:"tagNum:25,optional"`
	SMTotalNumber               *int64                             `ber:"tagNum:26,optional"`
	SMServiceType               *SMServiceType                     `ber:"tagNum:27,optional"`
	SMSequenceNumber            *int64                             `ber:"tagNum:28,optional"`
	SMSResult                   *SMSResult                         `ber:"tagNum:29,optional"`
	SubmissionTime              *TimeStamp                         `ber:"tagNum:30,optional"`
	SMPriority                  *PriorityType                      `ber:"tagNum:31,optional"`
	MessageReference            *MessageReference                  `ber:"tagNum:32,optional"`
	MessageSize                 *int64                             `ber:"tagNum:33,optional"`
	MessageClass                *MessageClass                      `ber:"tagNum:34,optional"`
	SMdeliveryReportRequested   *SMdeliveryReportRequested         `ber:"tagNum:35,optional"`
	MessageClassTokenText       *asn.UTF8String                    `ber:"tagNum:36,optional"`
	UserRoamerInOut             *RoamerInOut                       `ber:"tagNum:37,optional"`
	UserLocationInformationASN1 *UserLocationInformationStructured `ber:"tagNum:38,optional"`
}
