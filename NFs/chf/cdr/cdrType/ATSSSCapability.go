package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	ATSSSCapabilityPresentATSSSLL               asn.Enumerated = 0
	ATSSSCapabilityPresentMPTCPATSSLL           asn.Enumerated = 1
	ATSSSCapabilityPresentMPTCPATSSLLASModeUL   asn.Enumerated = 2
	ATSSSCapabilityPresentMPTCPATSSLLExSDModeUL asn.Enumerated = 3
	ATSSSCapabilityPresentMPTCPATSSLLASModeDLUL asn.Enumerated = 4
)

type ATSSSCapability struct {
	Value asn.Enumerated
}
