package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type SubscriptionId struct {
	SubscriptionIdType SubscriptionIdType       `avp:"Subscription-Id-Type"`
	SubscriptionIdData diam_datatype.UTF8String `avp:"Subscription-Id-Data"`
}
