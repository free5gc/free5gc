package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	TriggerCategoryPresentImmediateReport asn.Enumerated = 0
	TriggerCategoryPresentDeferredReport  asn.Enumerated = 1
)

type TriggerCategory struct {
	Value asn.Enumerated
}
