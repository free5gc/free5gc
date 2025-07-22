package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type RemainingBalance struct {
	UnitValue    *UnitValue               `avp:"Unit-Value"`
	CurrencyCode diam_datatype.Unsigned32 `avp:"Currency-Code"`
}
