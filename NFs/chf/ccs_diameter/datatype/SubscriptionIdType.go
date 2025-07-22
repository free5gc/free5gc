package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	END_USER_E164    SubscriptionIdType = 0
	END_USER_IMSI    SubscriptionIdType = 1
	END_USER_SIP_URI SubscriptionIdType = 2
	END_USER_NAI     SubscriptionIdType = 3
)

type SubscriptionIdType diam_datatype.Enumerated
