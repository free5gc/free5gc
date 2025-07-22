package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SharingLevelPresentSHARED    asn.Enumerated = 0
	SharingLevelPresentNONSHARED asn.Enumerated = 1
)

type SharingLevel struct {
	Value asn.Enumerated
}
