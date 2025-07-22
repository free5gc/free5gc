package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	FAILOVER_NOT_SUPPORTED CcSessionFailover = 0
	FAILOVER_SUPPORTED     CcSessionFailover = 1
)

type CcSessionFailover diam_datatype.Enumerated
