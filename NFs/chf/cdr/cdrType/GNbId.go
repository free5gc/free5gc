package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type GNbId struct { /* Sequence Type */
	BitLength int64         `ber:"tagNum:0"`
	GNbValue  asn.IA5String `ber:"tagNum:1"`
}
