package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SteerModeValuePresentActiveStandby asn.Enumerated = 0
	SteerModeValuePresentLoadBalancing asn.Enumerated = 1
	SteerModeValuePresentSmallestDelay asn.Enumerated = 2
	SteerModeValuePresentPriorityBased asn.Enumerated = 3
)

type SteerModeValue struct {
	Value asn.Enumerated
}
