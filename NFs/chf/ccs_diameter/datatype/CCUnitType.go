package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	TIME                 CCUnitType = 0
	MONEY                CCUnitType = 1
	TOTALOCTETS          CCUnitType = 2
	INPUTOCTETS          CCUnitType = 3
	OUTPUTOCTETS         CCUnitType = 4
	SERVICESPECIFICUNITS CCUnitType = 5
)

type CCUnitType diam_datatype.Enumerated
