package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	MULTIPLE_SERVICES_NOT_SUPPORTED MultipleServicesIndicator = 0
	MULTIPLE_SERVICES_SUPPORTED     MultipleServicesIndicator = 1
)

type MultipleServicesIndicator diam_datatype.Enumerated
