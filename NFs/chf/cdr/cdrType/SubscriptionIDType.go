package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SubscriptionIDTypePresentENDUSERE164    asn.Enumerated = 0
	SubscriptionIDTypePresentENDUSERIMSI    asn.Enumerated = 1
	SubscriptionIDTypePresentENDUSERSIPURI  asn.Enumerated = 2
	SubscriptionIDTypePresentENDUSERNAI     asn.Enumerated = 3
	SubscriptionIDTypePresentENDUSERPRIVATE asn.Enumerated = 4
)

type SubscriptionIDType struct {
	Value asn.Enumerated
}
