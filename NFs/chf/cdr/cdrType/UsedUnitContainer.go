package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type UsedUnitContainer struct { /* Sequence Type */
	ServiceIdentifier *ServiceIdentifier `ber:"tagNum:0,optional"`
	Time              *CallDuration      `ber:"tagNum:1,optional"`
	/* Sequence of = 35, FULL Name = struct UsedUnitContainer__triggers */
	/* Trigger */
	Triggers                    []Trigger                 `ber:"tagNum:2,optional"`
	TriggerTimeStamp            *TimeStamp                `ber:"tagNum:3,optional"`
	DataTotalVolume             *DataVolumeOctets         `ber:"tagNum:4,optional"`
	DataVolumeUplink            *DataVolumeOctets         `ber:"tagNum:5,optional"`
	DataVolumeDownlink          *DataVolumeOctets         `ber:"tagNum:6,optional"`
	ServiceSpecificUnits        *int64                    `ber:"tagNum:7,optional"`
	EventTimeStamp              *TimeStamp                `ber:"tagNum:8,optional"`
	LocalSequenceNumber         *LocalSequenceNumber      `ber:"tagNum:9,optional"`
	RatingIndicator             *RatingIndicator          `ber:"tagNum:10,optional"`
	PDUContainerInformation     *PDUContainerInformation  `ber:"tagNum:11,optional"`
	QuotaManagementIndicator    *bool                     `ber:"tagNum:12,optional"`
	QuotaManagementIndicatorExt *QuotaManagementIndicator `ber:"tagNum:13,optional"`
	NSPAContainerInformation    *NSPAContainerInformation `ber:"tagNum:14,optional"`
	/* Sequence of = 35, FULL Name = struct UsedUnitContainer__eventTimeStampExt */
	/* TimeStamp */
	EventTimeStampExt []TimeStamp `ber:"tagNum:15,optional"`
}
