package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type ServiceRating struct {
	ServiceIdentifier              diam_datatype.Unsigned32       `avp:"Service-Identifier"`
	DestinationID                  diam_datatype.Grouped          `avp:"DestinationID"`
	ServiceInformation             diam_datatype.Grouped          `avp:"ServiceInformation"`
	Extension                      diam_datatype.Grouped          `avp:"Extension"`
	RequestSubType                 RequestSubType                 `avp:"RequestSubType"`
	Price                          diam_datatype.Unsigned32       `avp:"Price"`
	BillingInfo                    diam_datatype.UTF8String       `avp:"BillingInfo"`
	ImpactOnCounter                *diam_datatype.Grouped         `avp:"ImpactonCounter"`
	RequestedUnits                 diam_datatype.Unsigned32       `avp:"RequestedUnits"`
	ConsumedUnits                  diam_datatype.Unsigned32       `avp:"ConsumedUnits"`
	ConsumedUnitsAfterTariffSwitch diam_datatype.Unsigned32       `avp:"ConsumedUnitsAfterTariffSwitch"`
	TariffSwitchTime               diam_datatype.Unsigned32       `avp:"TariffSwitchTime"`
	MonetaryTariff                 *MonetaryTariff                `avp:"MonetaryTariff"`
	NextMonetaryTariff             *NextMonetaryTariff            `avp:"NextMonetaryTariff"`
	ExpiryTime                     diam_datatype.Time             `avp:"ExpiryTime"`
	ValidUnits                     diam_datatype.Unsigned32       `avp:"ValidUnits"`
	MonetaryTariffAfterValidUnits  *MonetaryTariffAfterValidUnits `avp:"MonetaryTariffAfterValidUnits"`
	MonetaryQuota                  diam_datatype.Unsigned32       `avp:"MonetaryQuota"`
	MinimalRequestedUnits          diam_datatype.Unsigned32       `avp:"MinimalRequestedUnits"`
	AllowedUnits                   diam_datatype.Unsigned32       `avp:"AllowedUnits"`
}
