package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type CostInformation struct {
	CurrencyCode diam_datatype.Unsigned32 `avp:"Currency-Code"`
	UnitValue    *UnitValue               `avp:"Unit-Value"`
	UnitCost     *UnitCost                `avp:"Unit-Cost"`
}
