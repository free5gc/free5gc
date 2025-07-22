package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type AllocationRetentionPriority struct { /* Sequence Type */
	PriorityLevel           int64                   `ber:"tagNum:1"`
	PreemptionCapability    PreemptionCapability    `ber:"tagNum:2"`
	PreemptionVulnerability PreemptionVulnerability `ber:"tagNum:3"`
}
