package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	DIRECT_DEBITING RequestedAction = 0
	REFUND_ACCOUNT  RequestedAction = 1
	CHECK_BALANCE   RequestedAction = 2
	PRICE_ENQUIRY   RequestedAction = 3
)

type RequestedAction diam_datatype.Enumerated
