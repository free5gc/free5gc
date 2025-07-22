package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type FiveGQoSInformation struct { /* Sequence Type */
	FiveQi                 *int64                       `ber:"tagNum:1,optional"`
	ARP                    *AllocationRetentionPriority `ber:"tagNum:2,optional"`
	QoSNotificationControl *bool                        `ber:"tagNum:3,optional"`
	ReflectiveQos          *bool                        `ber:"tagNum:4,optional"`
	MaxbitrateUL           *Bitrate                     `ber:"tagNum:5,optional"`
	MaxbitrateDL           *Bitrate                     `ber:"tagNum:6,optional"`
	GuaranteedbitrateUL    *Bitrate                     `ber:"tagNum:7,optional"`
	GuaranteedbitrateDL    *Bitrate                     `ber:"tagNum:8,optional"`
	PriorityLevel          *int64                       `ber:"tagNum:9,optional"`
	AverWindow             *int64                       `ber:"tagNum:10,optional"`
	MaxDataBurstVol        *int64                       `ber:"tagNum:11,optional"`
	MaxPacketLossRateDL    *int64                       `ber:"tagNum:12,optional"`
	MaxPacketLossRateUL    *int64                       `ber:"tagNum:13,optional"`
}
