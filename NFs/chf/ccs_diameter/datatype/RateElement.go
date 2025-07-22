package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type RateElement struct {
	CCUnitType         CCUnitType               `avp:"CC-Unit-Type"`
	ChargeReasonCode   ChargeReasonCode         `avp:"Charge-Reason-Code"`
	UnitValue          *UnitValue               `avp:"Unit-Value"`
	UnitCost           *UnitCost                `avp:"Unit-Cost"`
	UnitQuotaThreshold diam_datatype.Unsigned32 `avp:"Unit-Quota-Threshold"`
}
