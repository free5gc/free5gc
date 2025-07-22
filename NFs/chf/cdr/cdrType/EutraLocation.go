package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type EutraLocation struct { /* Sequence Type */
	Tai                      *TAI                      `ber:"tagNum:0,optional"`
	Ecgi                     *Ecgi                     `ber:"tagNum:1,optional"`
	AgeOfLocationInformation *AgeOfLocationInformation `ber:"tagNum:3,optional"`
	UeLocationTimestamp      *TimeStamp                `ber:"tagNum:4,optional"`
	GeographicalInformation  *GeographicalInformation  `ber:"tagNum:5,optional"`
	GeodeticInformation      *GeodeticInformation      `ber:"tagNum:6,optional"`
	GlobalNgenbId            *GlobalRanNodeId          `ber:"tagNum:7,optional"`
	GlobalENbId              *GlobalRanNodeId          `ber:"tagNum:8,optional"`
}
