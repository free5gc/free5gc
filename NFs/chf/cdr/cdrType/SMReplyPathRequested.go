package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SMReplyPathRequestedPresentNoReplyPathSet asn.Enumerated = 0
	SMReplyPathRequestedPresentReplyPathSet   asn.Enumerated = 1
)

type SMReplyPathRequested struct {
	Value asn.Enumerated
}
