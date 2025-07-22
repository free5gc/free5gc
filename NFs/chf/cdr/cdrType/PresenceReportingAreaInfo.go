package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type PresenceReportingAreaInfo struct { /* Sequence Type */
	PresenceReportingAreaIdentifier   asn.OctetString                    `ber:"tagNum:0"`
	PresenceReportingAreaStatus       *PresenceReportingAreaStatus       `ber:"tagNum:1,optional"`
	PresenceReportingAreaElementsList *PresenceReportingAreaElementsList `ber:"tagNum:2,optional"`
	PresenceReportingAreaNode         *PresenceReportingAreaNode         `ber:"tagNum:3,optional"`
}
