package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type QosFlowsUsageReport struct { /* Sequence Type */
	QosFlowId          *QoSFlowId       `ber:"tagNum:0,optional"`
	StartTime          TimeStamp        `ber:"tagNum:1"`
	EndTime            TimeStamp        `ber:"tagNum:2"`
	DataVolumeDownlink DataVolumeOctets `ber:"tagNum:3"`
	DataVolumeUplink   DataVolumeOctets `ber:"tagNum:4"`
}
