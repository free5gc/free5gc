package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	INITIAL_REQUEST     CcRequestType = 1
	UPDATE_REQUEST      CcRequestType = 2
	TERMINATION_REQUEST CcRequestType = 3
	EVENT_REQUEST       CcRequestType = 4
)

type CcRequestType diam_datatype.Enumerated
