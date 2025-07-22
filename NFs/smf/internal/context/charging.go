package context

import (
	"github.com/free5gc/openapi/models"
)

type ChargingLevel uint8

// For a rating group that is pdu session charging level, all volume in a pdu session will be charged
// For a rating group that is flow charging level (or Rating group level (32.255)),
// only volume in a flow will be charged
const (
	PduSessionCharging ChargingLevel = iota
	FlowCharging
)

type RequestType uint8

// For each charging event, it will have a corresponding charging request type, see 32.255 Table 5.2.1.4.1
const (
	CHARGING_INIT RequestType = iota
	CHARGING_UPDATE
	CHARGING_RELEASE
)

type ChargingInfo struct {
	ChargingMethod         models.QuotaManagementIndicator
	VolumeLimitExpiryTimer *Timer
	EventLimitExpiryTimer  *Timer
	ChargingLevel          ChargingLevel
	RatingGroup            int32
	UpfId                  string
}
