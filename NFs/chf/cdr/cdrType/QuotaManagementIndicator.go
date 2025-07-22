package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	QuotaManagementIndicatorPresentOnlineCharging           asn.Enumerated = 0
	QuotaManagementIndicatorPresentOfflineCharging          asn.Enumerated = 1
	QuotaManagementIndicatorPresentQuotaManagementSuspended asn.Enumerated = 2
)

type QuotaManagementIndicator struct {
	Value asn.Enumerated
}
