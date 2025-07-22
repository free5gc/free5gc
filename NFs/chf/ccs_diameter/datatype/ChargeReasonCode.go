package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	UNKNOWN                      ChargeReasonCode = 0
	USAGE                        ChargeReasonCode = 1
	COMMUNICATION_ATTEMPT_CHARGE ChargeReasonCode = 2
	SETUP_CHARGE                 ChargeReasonCode = 3
	ADD_ON_CHARGE                ChargeReasonCode = 4
)

type ChargeReasonCode diam_datatype.Enumerated
