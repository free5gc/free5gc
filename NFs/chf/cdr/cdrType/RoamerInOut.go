package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	RoamerInOutPresentRoamerInBound  asn.Enumerated = 0
	RoamerInOutPresentRoamerOutBound asn.Enumerated = 1
)

type RoamerInOut struct {
	Value asn.Enumerated
}
