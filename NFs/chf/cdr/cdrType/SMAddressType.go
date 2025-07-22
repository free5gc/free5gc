package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SMAddressTypePresentEmailAddress          asn.Enumerated = 0
	SMAddressTypePresentMSISDN                asn.Enumerated = 1
	SMAddressTypePresentIPv4Address           asn.Enumerated = 2
	SMAddressTypePresentIPv6Address           asn.Enumerated = 3
	SMAddressTypePresentNumericShortCode      asn.Enumerated = 4
	SMAddressTypePresentAlphanumericShortCode asn.Enumerated = 5
	SMAddressTypePresentOther                 asn.Enumerated = 6
	SMAddressTypePresentIMSI                  asn.Enumerated = 7
	SMAddressTypePresentNAI                   asn.Enumerated = 8
	SMAddressTypePresentExternalId            asn.Enumerated = 9
)

type SMAddressType struct {
	Value asn.Enumerated
}
