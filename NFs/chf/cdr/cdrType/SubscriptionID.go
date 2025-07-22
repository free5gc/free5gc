package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type SubscriptionID struct { /* Set Type */
	SubscriptionIDType SubscriptionIDType `ber:"tagNum:0"`
	SubscriptionIDData asn.UTF8String     `ber:"tagNum:1"`
}
