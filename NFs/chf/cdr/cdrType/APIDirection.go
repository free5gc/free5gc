package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	APIDirectionPresentInvocation   asn.Enumerated = 0
	APIDirectionPresentNotification asn.Enumerated = 1
)

type APIDirection struct {
	Value asn.Enumerated
}
