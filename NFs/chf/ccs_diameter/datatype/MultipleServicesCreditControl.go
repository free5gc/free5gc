package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type FinalUnitIndication struct {
	FinalUnitAction       FinalUnitAction            `avp:"Final-Unit-Action"`
	RestrictionFilterRule diam_datatype.IPFilterRule `avp:"Restriction-Filter-Rule"`
	FilterId              diam_datatype.UTF8String   `avp:"Filter-Id"`
	RedirectServer        diam_datatype.Grouped      `avp:"Redirect-Server"`
}
