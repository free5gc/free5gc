package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type NsiLoadLevelInfo struct { /* Sequence Type */
	LoadLevelInformation *int64           `ber:"tagNum:0,optional"`
	Snssai               *SingleNSSAI     `ber:"tagNum:1,optional"`
	NsiId                *asn.OctetString `ber:"tagNum:2,optional"`
}
