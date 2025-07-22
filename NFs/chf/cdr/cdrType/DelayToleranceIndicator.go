package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	DelayToleranceIndicatorPresentDTSupported    asn.Enumerated = 0
	DelayToleranceIndicatorPresentDTNotSupported asn.Enumerated = 1
)

type DelayToleranceIndicator struct {
	Value asn.Enumerated
}
