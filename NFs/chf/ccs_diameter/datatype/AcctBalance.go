package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type AcctBalance struct {
	AcctBalanceId diam_datatype.Unsigned64 `avp:"Acct-Balance-Id"`
	UnitValue     *UnitValue               `avp:"Unit-Value"`
}
