package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type ABResponse struct {
	AcctBalance *AcctBalance          `avp:"Acct-Balance"`
	Counter     diam_datatype.Grouped `avp:"Counter"`
}
