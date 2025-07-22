package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type ServiceUsageRequest struct {
	SessionId           diam_datatype.UTF8String       `avp:"Session-Id"`
	OriginHost          diam_datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm         diam_datatype.DiameterIdentity `avp:"Origin-Realm"`
	DestinationRealm    diam_datatype.DiameterIdentity `avp:"Destination-Realm"`
	DestinationHost     diam_datatype.DiameterIdentity `avp:"Destination-Host"`
	VendorSpecificAppId diam_datatype.Grouped          `avp:"Vendor-Specific-Application-Id"`
	UserName            diam_datatype.OctetString      `avp:"User-Name"`
	EventTimestamp      diam_datatype.Time             `avp:"Event-Timestamp"`
	BeginTime           diam_datatype.Time             `avp:"BeginTime"`
	ActualTime          diam_datatype.Time             `avp:"ActualTime"`
	SubscriptionId      *SubscriptionId                `avp:"Subscription-Id"`
	ServiceRating       *ServiceRating                 `avp:"Service-Rating"`
}
