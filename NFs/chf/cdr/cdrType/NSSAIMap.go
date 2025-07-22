package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type NSSAIMap struct { /* Sequence Type */
	ServingSnssai SingleNSSAI `ber:"tagNum:0"`
	HomeSnssai    SingleNSSAI `ber:"tagNum:1"`
}
