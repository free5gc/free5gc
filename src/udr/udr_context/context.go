package udr_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
)

var udrContext = UDRContext{}

type ueId = string
type ueGroupId = string
type subsId = string

type UDRServiceType int

const (
	NUDR_DR UDRServiceType = iota
)

func init() {
	UDR_Self().Name = "udr"
	UDR_Self().EeSubscriptionIDGenerator = 1
	UDR_Self().SdmSubscriptionIDGenerator = 1
	UDR_Self().SubscriptionDataSubscriptionIDGenerator = 1
	UDR_Self().PolicyDataSubscriptionIDGenerator = 1
	UDR_Self().UESubsCollection = make(map[ueId]*UESubsData)
	UDR_Self().UEGroupCollection = make(map[ueGroupId]*UEGroupSubsData)
	UDR_Self().SubscriptionDataSubscriptions = make(map[subsId]*models.SubscriptionDataSubscriptions)
	UDR_Self().PolicyDataSubscriptions = make(map[subsId]*models.PolicyDataSubscription)
}

type UDRContext struct {
	Name                                    string
	UriScheme                               models.UriScheme
	HttpIpv4Port                            int
	HttpIPv4Address                         string
	HttpIPv6Address                         string
	NfId                                    string
	NrfUri                                  string
	EeSubscriptionIDGenerator               int
	SdmSubscriptionIDGenerator              int
	PolicyDataSubscriptionIDGenerator       int
	UESubsCollection                        map[ueId]*UESubsData
	UEGroupCollection                       map[ueGroupId]*UEGroupSubsData
	SubscriptionDataSubscriptionIDGenerator int
	SubscriptionDataSubscriptions           map[subsId]*models.SubscriptionDataSubscriptions
	PolicyDataSubscriptions                 map[subsId]*models.PolicyDataSubscription
}

type UESubsData struct {
	EeSubscriptionCollection map[subsId]*EeSubscriptionCollection
	SdmSubscriptions         map[subsId]*models.SdmSubscription
}

type UEGroupSubsData struct {
	EeSubscriptions map[subsId]*models.EeSubscription
}

type EeSubscriptionCollection struct {
	EeSubscriptions      *models.EeSubscription
	AmfSubscriptionInfos []models.AmfSubscriptionInfo
}

// Reset UDR Context
func (context *UDRContext) Reset() {
	for key := range context.UESubsCollection {
		delete(context.UESubsCollection, key)
	}
	for key := range context.UEGroupCollection {
		delete(context.UEGroupCollection, key)
	}
	for key := range context.SubscriptionDataSubscriptions {
		delete(context.SubscriptionDataSubscriptions, key)
	}
	for key := range context.PolicyDataSubscriptions {
		delete(context.PolicyDataSubscriptions, key)
	}
	context.EeSubscriptionIDGenerator = 1
	context.SdmSubscriptionIDGenerator = 1
	context.SubscriptionDataSubscriptionIDGenerator = 1
	context.PolicyDataSubscriptionIDGenerator = 1
	context.UriScheme = models.UriScheme_HTTPS
	context.Name = "udr"
}

func (context *UDRContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, context.HttpIpv4Port)
}

func (context *UDRContext) GetIPv4GroupUri(udrServiceType UDRServiceType) string {
	var serviceUri string

	switch udrServiceType {
	case NUDR_DR:
		serviceUri = "/nudr-dr/v1"
	default:
		serviceUri = ""
	}

	return fmt.Sprintf("%s://%s:%d%s", context.UriScheme, context.HttpIPv4Address, context.HttpIpv4Port, serviceUri)
}

// Create new UDR context
func UDR_Self() *UDRContext {
	return &udrContext
}
