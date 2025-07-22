package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	PartialRecordMethodPresentDefault    asn.Enumerated = 0
	PartialRecordMethodPresentIndividual asn.Enumerated = 1
)

type PartialRecordMethod struct {
	Value asn.Enumerated
}
