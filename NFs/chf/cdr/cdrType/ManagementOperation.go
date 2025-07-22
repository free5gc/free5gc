package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	ManagementOperationPresentCreateMOI           asn.Enumerated = 0
	ManagementOperationPresentModifyMOIAttributes asn.Enumerated = 1
	ManagementOperationPresentDeleteMOI           asn.Enumerated = 2
)

type ManagementOperation struct {
	Value asn.Enumerated
}
