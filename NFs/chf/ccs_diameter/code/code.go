package code

// AVP that are specifies in 32.296 with AVP code xxxx

const (
	Re_interface        = 16777218
	ServiceUsageMessage = 111
	ABMF_CreditControl  = 272
)

const (
	BeginTime = iota + 7000
	ActualTime
	ServiceRating
	DestinationID
	Extension
	Price
	BillingInfo
	TariffSwitchTime
	MonetaryTariff
	NextMonetaryTariff
	ExpiryTime
	ValidUnits
	MonetaryTariffAfterValidUnits
	RequestSubType
	ConsumedUnits
	ConsumedUnitsAfterTariffSwitch
	MonetaryQuota
	RequestedUnits
	MinimalRequestedUnits
	ServiceInformation
	ImpactonCounter
	AllowedUnits
	RequestedCounters
	CounterTariff
	CounterPrice
	BasicPriceTimeStamp
	Counter
	VendorSpecificAppId
	ABResponse
	AcctBalanceId
)
