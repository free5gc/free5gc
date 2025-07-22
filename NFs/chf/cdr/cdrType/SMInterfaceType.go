package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SMInterfaceTypePresentUnkown                 asn.Enumerated = 0
	SMInterfaceTypePresentMobileOriginating      asn.Enumerated = 1
	SMInterfaceTypePresentMobileTerminating      asn.Enumerated = 2
	SMInterfaceTypePresentApplicationOriginating asn.Enumerated = 3
	SMInterfaceTypePresentApplicationTerminating asn.Enumerated = 4
	SMInterfaceTypePresentDeviceTrigger          asn.Enumerated = 5
)

type SMInterfaceType struct {
	Value asn.Enumerated
}
