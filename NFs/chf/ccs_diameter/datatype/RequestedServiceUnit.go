package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type RequestedServiceUnit struct {
	TariffChangeUsage      diam_datatype.Enumerated `avp:"Tariff-Change-Usage"`
	CCTime                 diam_datatype.Unsigned32 `avp:"CC-Time"`
	CCMoney                *CCMoney                 `avp:"CC-Money"`
	CCTotalOctets          diam_datatype.Unsigned64 `avp:"CC-Total-Octets"`
	CCInputOctets          diam_datatype.Unsigned64 `avp:"CC-Input-Octets"`
	CCOutputOctets         diam_datatype.Unsigned64 `avp:"CC-Output-Octets"`
	CCServiceSpecificUnits diam_datatype.Unsigned64 `avp:"CC-Service-Specific-Units"`
}
