package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type ServiceUsageResponse struct {
	SessionId           diam_datatype.UTF8String       `avp:"Session-Id"`
	OriginHost          diam_datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm         diam_datatype.DiameterIdentity `avp:"Origin-Realm"`
	VendorSpecificAppId diam_datatype.Grouped          `avp:"Vendor-Specific-Application-Id"`
	EventTimestamp      diam_datatype.Time             `avp:"Event-Timestamp"`
	ServiceRating       *ServiceRating                 `avp:"Service-Rating"`
}
