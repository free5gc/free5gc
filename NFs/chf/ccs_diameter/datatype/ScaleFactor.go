package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type ScaleFactor struct {
	ValueDigits diam_datatype.Integer64 `avp:"Value-Digits"`
	Exponent    diam_datatype.Integer32 `avp:"Exponent"`
}
