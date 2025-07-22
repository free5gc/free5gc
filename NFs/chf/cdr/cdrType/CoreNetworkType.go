package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	CoreNetworkTypePresentFiveGC asn.Enumerated = 0
	CoreNetworkTypePresentEPC    asn.Enumerated = 1
)

type CoreNetworkType struct {
	Value asn.Enumerated
}
