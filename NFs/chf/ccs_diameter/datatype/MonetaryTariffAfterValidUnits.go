package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type MonetaryTariffAfterValidUnits struct {
	CurrencyCode diam_datatype.Unsigned32 `avp:"Currency-Code"`
	ScaleFactor  *ScaleFactor             `avp:"Scale-Factor"`
	RateElement  *RateElement             `avp:"Rate-Element"`
}
