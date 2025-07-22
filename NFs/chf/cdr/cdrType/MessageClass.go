package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	MessageClassPresentPersonal           asn.Enumerated = 0
	MessageClassPresentAdvertisement      asn.Enumerated = 1
	MessageClassPresentInformationService asn.Enumerated = 2
	MessageClassPresentAuto               asn.Enumerated = 3
)

type MessageClass struct {
	Value asn.Enumerated
}
