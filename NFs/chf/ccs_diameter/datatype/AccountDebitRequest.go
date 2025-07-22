package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type AccountDebitRequest struct {
	SessionId                     diam_datatype.UTF8String       `avp:"Session-Id"`
	OriginHost                    diam_datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm                   diam_datatype.DiameterIdentity `avp:"Origin-Realm"`
	DestinationRealm              diam_datatype.DiameterIdentity `avp:"Destination-Realm"`
	DestinationHost               diam_datatype.DiameterIdentity `avp:"Destination-Host"`
	UserName                      diam_datatype.OctetString      `avp:"User-Name"`
	OriginStateId                 diam_datatype.Unsigned32       `avp:"Origin-State-Id"`
	RequestedAction               RequestedAction                `avp:"Requested-Action"`
	AuthApplicationId             diam_datatype.Unsigned32       `avp:"Auth-Application-Id"`
	ServiceContextId              diam_datatype.UTF8String       `avp:"Service-Context-Id"`
	CcRequestType                 CcRequestType                  `avp:"CC-Request-Type"`
	CcRequestNumber               diam_datatype.Unsigned32       `avp:"CC-Request-Number"`
	EventTimestamp                diam_datatype.Time             `avp:"Event-Timestamp"`
	TerminationCause              TerminationCause               `avp:"Termination-Cause"`
	SubscriptionId                *SubscriptionId                `avp:"Subscription-Id"`
	ServiceIdentifier             diam_datatype.Unsigned32       `avp:"Service-Identifier"`
	ServiceInformation            diam_datatype.Grouped          `avp:"Service-Information"`
	MultipleServicesIndicator     MultipleServicesIndicator      `avp:"Multiple-Services-Indicator"`
	ProxyInfo                     diam_datatype.Grouped          `avp:"Proxy-Info"`
	MultipleServicesCreditControl *MultipleServicesCreditControl `avp:"Multiple-Services-Credit-Control"`
}
