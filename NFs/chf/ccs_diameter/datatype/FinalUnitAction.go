package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	TERMINATE       FinalUnitAction = 0
	REDIRECT        FinalUnitAction = 1
	RESTRICT_ACCESS FinalUnitAction = 2
)

type FinalUnitAction diam_datatype.Enumerated
