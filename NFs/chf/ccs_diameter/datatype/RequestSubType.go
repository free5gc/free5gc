package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	REQ_SUBTYPE_AOC     RequestSubType = 0
	REQ_SUBTYPE_RESERVE RequestSubType = 1
	REQ_SUBTYPE_DEBIT   RequestSubType = 2
	REQ_SUBTYPE_RELEASE RequestSubType = 3
)

type RequestSubType diam_datatype.Enumerated
