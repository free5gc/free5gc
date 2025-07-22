package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	NOT_APPLICABLE LowBalanceIndication = 0
	YES            LowBalanceIndication = 1
)

type LowBalanceIndication diam_datatype.Enumerated
