package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	AccessTypePresentThreeGPPAccess    asn.Enumerated = 0
	AccessTypePresentNonThreeGPPAccess asn.Enumerated = 1
)

type AccessType struct {
	Value asn.Enumerated
}
