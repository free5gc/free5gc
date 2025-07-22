package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type NrLocation struct { /* Sequence Type */
	Tai                      *TAI                      `ber:"tagNum:0,optional"`
	Ncgi                     *Ncgi                     `ber:"tagNum:1,optional"`
	AgeOfLocationInformation *AgeOfLocationInformation `ber:"tagNum:2,optional"`
	UeLocationTimestamp      *TimeStamp                `ber:"tagNum:3,optional"`
	GeographicalInformation  *GeographicalInformation  `ber:"tagNum:4,optional"`
	GeodeticInformation      *GeodeticInformation      `ber:"tagNum:5,optional"`
	GlobalGnbId              *GlobalRanNodeId          `ber:"tagNum:6,optional"`
}
