package amf_producer

import (
	"github.com/sirupsen/logrus"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/gmm/gmm_state"
	"free5gc/src/amf/logger"
	"strconv"
	"time"
)

var HttpLog *logrus.Entry

func init() {
	// init Pool
	HttpLog = logger.HttpLog
}

// TODO: handle event filter
func CreateAMFEventSubscription(context *amf_context.AMFContext, request models.AmfCreateEventSubscription, recieveTime time.Time) (response *models.AmfCreatedEventSubscription, err models.ProblemDetails) {
	response = &models.AmfCreatedEventSubscription{}
	subscription := request.Subscription
	contextEventSubscription := &amf_context.AMFContextEventSubscription{}
	contextEventSubscription.EventSubscription = *subscription
	newSubscriptionID := strconv.Itoa(context.EventSubscriptionIDGenerator)

	var isImmediate bool
	var immediateFlags []bool
	var reportlist []models.AmfEventReport

	// store subscription in amf_context
	ueEventSubscription := amf_context.AmfUeEventSubscription{}
	ueEventSubscription.EventSubscription = &contextEventSubscription.EventSubscription
	ueEventSubscription.Timestamp = recieveTime

	if subscription.Options != nil && subscription.Options.Trigger == models.AmfEventTrigger_CONTINUOUS {
		ueEventSubscription.RemainReports = new(int32)
		*ueEventSubscription.RemainReports = subscription.Options.MaxReports
	}
	for _, events := range *subscription.EventList {
		immediateFlags = append(immediateFlags, events.ImmediateFlag)
		if events.ImmediateFlag {
			isImmediate = true
		}
	}

	if subscription.AnyUE {
		contextEventSubscription.IsAnyUe = true
		ueEventSubscription.AnyUe = true
		for _, ue := range context.UePool {
			ue.EventSubscriptionsInfo[newSubscriptionID] = new(amf_context.AmfUeEventSubscription)
			*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
			contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)
		}
	} else if subscription.GroupId != "" {
		contextEventSubscription.IsGroupUe = true
		ueEventSubscription.AnyUe = true
		for _, ue := range context.UePool {
			if ue.GroupID == subscription.GroupId {
				ue.EventSubscriptionsInfo[newSubscriptionID] = new(amf_context.AmfUeEventSubscription)
				*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
				contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)
			}
		}
	} else {
		if ue, ok := context.UePool[subscription.Supi]; !ok {
			err.Status = 403
			err.Cause = "UE_NOT_SERVED_BY_AMF"
			return nil, err
		} else {
			ue.EventSubscriptionsInfo[newSubscriptionID] = new(amf_context.AmfUeEventSubscription)
			*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
			contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)

		}
	}

	// delete subscription
	if subscription.Options != nil {
		contextEventSubscription.Expiry = subscription.Options.Expiry
	}
	context.EventSubscriptionIDGenerator++
	context.EventSubscriptions[newSubscriptionID] = contextEventSubscription

	// build response

	response.Subscription = subscription
	response.SubscriptionId = newSubscriptionID

	// for immediate use
	if subscription.AnyUE {
		for _, ue := range context.UePool {
			if isImmediate {
				subReports(ue, newSubscriptionID)
			}
			for i, flag := range immediateFlags {
				if flag {
					report, ok := NewAmfEventReport(ue, (*subscription.EventList)[i].Type, newSubscriptionID)
					if ok {
						reportlist = append(reportlist, report)
					}
				}
			}
			// delete subscription
			if len := len(reportlist); len > 0 && (!reportlist[len-1].State.Active) {
				delete(ue.EventSubscriptionsInfo, newSubscriptionID)
			}
		}
	} else if subscription.GroupId != "" {
		for _, ue := range context.UePool {
			if isImmediate {
				subReports(ue, newSubscriptionID)
			}
			if ue.GroupID == subscription.GroupId {
				for i, flag := range immediateFlags {
					if flag {
						report, ok := NewAmfEventReport(ue, (*subscription.EventList)[i].Type, newSubscriptionID)
						if ok {
							reportlist = append(reportlist, report)
						}
					}
				}
				// delete subscription
				if len := len(reportlist); len > 0 && (!reportlist[len-1].State.Active) {
					delete(ue.EventSubscriptionsInfo, newSubscriptionID)
				}
			}
		}
	} else {
		ue := context.UePool[subscription.Supi]
		if isImmediate {
			subReports(ue, newSubscriptionID)
		}
		for i, flag := range immediateFlags {
			if flag {
				report, ok := NewAmfEventReport(ue, (*subscription.EventList)[i].Type, newSubscriptionID)
				if ok {
					reportlist = append(reportlist, report)
				}
			}
		}
		// delete subscription
		if len := len(reportlist); len > 0 && (!reportlist[len-1].State.Active) {
			delete(ue.EventSubscriptionsInfo, newSubscriptionID)
		}
	}
	if len(reportlist) > 0 {
		response.ReportList = reportlist
		// delete subscription
		if !reportlist[0].State.Active {
			delete(context.EventSubscriptions, newSubscriptionID)
		}
	}

	return
}

func DeleteAMFEventSubscription(context *amf_context.AMFContext, subscriptionId string) (err models.ProblemDetails) {
	contextSubscription, ok := context.EventSubscriptions[subscriptionId]
	if !ok {
		err.Status = 404
		err.Cause = "SUBSCRIPTION_NOT_FOUND"
		return
	}
	for _, supi := range contextSubscription.UeSupiList {
		ue, ok := context.UePool[supi]
		if ok {
			delete(ue.EventSubscriptionsInfo, subscriptionId)
		}
	}
	delete(context.EventSubscriptions, subscriptionId)
	return
}

func ModifyAMFEventSubscription(context *amf_context.AMFContext, subscriptionId string, request models.ModifySubscriptionRequest) (err models.ProblemDetails) {
	contextSubscription, ok := context.EventSubscriptions[subscriptionId]
	if !ok {
		err.Status = 404
		err.Cause = "SUBSCRIPTION_NOT_FOUND"
		return
	}
	if request.OptionItem != nil {
		contextSubscription.Expiry = request.OptionItem.Value
	} else if request.SubscriptionItemInner != nil {
		subscription := &contextSubscription.EventSubscription
		if !contextSubscription.IsAnyUe && !contextSubscription.IsGroupUe {
			if _, ok := context.UePool[subscription.Supi]; !ok {
				err.Status = 403
				err.Cause = "UE_NOT_SERVED_BY_AMF"
				return err
			}
		}
		op := request.SubscriptionItemInner.Op
		index, _ := strconv.Atoi(request.SubscriptionItemInner.Path[11:])
		lists := (*subscription.EventList)
		len := len(*subscription.EventList)
		switch op {
		case "replace":
			event := *request.SubscriptionItemInner.Value
			if index < len {
				(*subscription.EventList)[index] = event
			}
		case "remove":
			if index < len {
				*subscription.EventList = append(lists[:index], lists[index+1:]...)
			}
		case "add":
			event := *request.SubscriptionItemInner.Value
			*subscription.EventList = append(lists, event)
		}
	}
	return
}

func subReports(ue *amf_context.AmfUe, subscriptionId string) {
	remainReport := ue.EventSubscriptionsInfo[subscriptionId].RemainReports
	if remainReport == nil {
		return
	}
	*remainReport--
}

// DO NOT handle AmfEventType_PRESENCE_IN_AOI_REPORT and AmfEventType_UES_IN_AREA_REPORT(about area)
func NewAmfEventReport(ue *amf_context.AmfUe, Type models.AmfEventType, subscriptionId string) (report models.AmfEventReport, ok bool) {
	ueSubscription, ok := ue.EventSubscriptionsInfo[subscriptionId]
	if !ok {
		return
	}

	report.AnyUe = ueSubscription.AnyUe
	report.Supi = ue.Supi
	report.Type = Type
	report.TimeStamp = &ueSubscription.Timestamp
	report.State = new(models.AmfEventState)
	mode := ueSubscription.EventSubscription.Options
	if mode == nil {
		report.State.Active = true
	} else if mode.Trigger == models.AmfEventTrigger_ONE_TIME {
		report.State.Active = false
	} else if *ueSubscription.RemainReports <= 0 {
		report.State.Active = false
	} else {
		report.State.Active = getDuration(mode.Expiry, &report.State.RemainDuration)
		if report.State.Active {
			report.State.RemainReports = *ueSubscription.RemainReports
		}
	}

	switch Type {
	case models.AmfEventType_LOCATION_REPORT:
		report.Location = &ue.Location
	// case models.AmfEventType_PRESENCE_IN_AOI_REPORT:
	// report.AreaList = (*subscription.EventList)[eventIndex].AreaList
	case models.AmfEventType_TIMEZONE_REPORT:
		report.Timezone = ue.TimeZone
	case models.AmfEventType_ACCESS_TYPE_REPORT:
		for accessType, sm := range ue.Sm {
			if sm.Check(gmm_state.REGISTERED) {
				report.AccessTypeList = append(report.AccessTypeList, accessType)
			}
		}
	case models.AmfEventType_REGISTRATION_STATE_REPORT:
		var rmInfos []models.RmInfo
		for accessType, sm := range ue.Sm {
			rmInfo := models.RmInfo{
				RmState:    models.RmState_DEREGISTERED,
				AccessType: accessType,
			}
			if sm.Check(gmm_state.REGISTERED) {
				rmInfo.RmState = models.RmState_REGISTERED
			}
			rmInfos = append(rmInfos, rmInfo)
		}
		report.RmInfoList = rmInfos
	case models.AmfEventType_CONNECTIVITY_STATE_REPORT:
		report.CmInfoList = ue.GetCmInfo()
	case models.AmfEventType_REACHABILITY_REPORT:
		report.Reachability = ue.Reachability
	case models.AmfEventType_SUBSCRIBED_DATA_REPORT:
		report.SubscribedData = &ue.SubscribedData
	case models.AmfEventType_COMMUNICATION_FAILURE_REPORT:
		// TODO : report.CommFailure
	case models.AmfEventType_SUBSCRIPTION_ID_CHANGE:
		report.SubscriptionId = subscriptionId
	case models.AmfEventType_SUBSCRIPTION_ID_ADDITION:
		report.SubscriptionId = subscriptionId
	}
	return

}

func getDuration(expiry *time.Time, remainDuration *int32) bool {

	if expiry != nil {
		if time.Now().After(*expiry) {
			return false
		} else {
			duration := time.Until(*expiry)
			*remainDuration = int32(duration.Seconds())
		}
	}
	return true

}
