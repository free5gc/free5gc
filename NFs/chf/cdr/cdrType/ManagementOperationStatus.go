package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	ManagementOperationStatusPresentOPERATIONSUCCEEDED asn.Enumerated = 0
	ManagementOperationStatusPresentOPERATIONFAILED    asn.Enumerated = 1
)

type ManagementOperationStatus struct {
	Value asn.Enumerated
}
