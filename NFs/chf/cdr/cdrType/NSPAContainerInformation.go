package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type NSPAContainerInformation struct { /* Sequence Type */
	Latency                         *int64                 `ber:"tagNum:0,optional"`
	Throughput                      *Throughput            `ber:"tagNum:1,optional"`
	MaximumPacketLossRate           *asn.UTF8String        `ber:"tagNum:3,optional"`
	ServiceExperienceStatisticsData *ServiceExperienceInfo `ber:"tagNum:4,optional"`
	NumberOfPDUSessions             *int64                 `ber:"tagNum:5,optional"`
	NumberOfRegisteredSubscribers   *int64                 `ber:"tagNum:6,optional"`
	LoadLevel                       *NsiLoadLevelInfo      `ber:"tagNum:7,optional"`
}
