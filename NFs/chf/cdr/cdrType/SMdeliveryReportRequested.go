package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SMdeliveryReportRequestedPresentYes asn.Enumerated = 0
	SMdeliveryReportRequestedPresentNo  asn.Enumerated = 1
)

type SMdeliveryReportRequested struct {
	Value asn.Enumerated
}
