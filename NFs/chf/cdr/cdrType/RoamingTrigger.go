package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type RoamingTrigger struct { /* Sequence Type */
	Trigger                 *SMFTrigger       `ber:"tagNum:0,optional"`
	TriggerCategory         *TriggerCategory  `ber:"tagNum:1,optional"`
	TimeLimit               *CallDuration     `ber:"tagNum:2,optional"`
	VolumeLimit             *DataVolumeOctets `ber:"tagNum:3,optional"`
	MaxNbChargingConditions *int64            `ber:"tagNum:4,optional"`
}
