package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type AccountDebitResponse struct {
	SessionId                     diam_datatype.UTF8String       `avp:"Session-Id"`
	ResultCode                    diam_datatype.Unsigned32       `avp:"Result-Code"`
	OriginHost                    diam_datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm                   diam_datatype.DiameterIdentity `avp:"Origin-Realm"`
	AuthApplicationId             diam_datatype.Unsigned32       `avp:"Auth-Application-Id"`
	ServiceContextId              diam_datatype.UTF8String       `avp:"Service-Context-Id"`
	CcRequestType                 CcRequestType                  `avp:"CC-Request-Type"`
	CcRequestNumber               diam_datatype.Unsigned32       `avp:"CC-Request-Number"`
	CCSessionFailover             CcSessionFailover              `avp:"CC-Session-Failover"`
	CostInformation               *CostInformation               `avp:"Cost-Information"`
	LowBalanceIndication          LowBalanceIndication           `avp:"Low-Balance-Indication"`
	EventTimestamp                diam_datatype.Time             `avp:"Event-Timestamp"`
	RemainingBalance              *RemainingBalance              `avp:"Remaining-Balance"`
	ABResponse                    *ABResponse                    `avp:"AB-Response"`
	ProxyInfo                     diam_datatype.Grouped          `avp:"Proxy-Info"`
	MultipleServicesCreditControl *MultipleServicesCreditControl `avp:"Multiple-Services-Credit-Control"`
}
