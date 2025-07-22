package processor

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi/models"
)

func (p *Processor) HandleCreateAMFEventSubscription(c *gin.Context,
	createEventSubscription models.AmfCreateEventSubscription,
) {
	createdEventSubscription, problemDetails := p.CreateAMFEventSubscriptionProcedure(createEventSubscription)
	if createdEventSubscription != nil {
		c.JSON(http.StatusCreated, createdEventSubscription)
	} else if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED_NF_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	}
}

// TODO: handle event filter
func (p *Processor) CreateAMFEventSubscriptionProcedure(createEventSubscription models.AmfCreateEventSubscription) (
	*models.AmfCreatedEventSubscription, *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	createdEventSubscription := &models.AmfCreatedEventSubscription{}
	subscription := createEventSubscription.Subscription
	if subscription == nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "SUBSCRIPTION_EMPTY",
		}
		return nil, problemDetails
	}
	contextEventSubscription := &context.AMFContextEventSubscription{}
	contextEventSubscription.EventSubscription = *subscription
	var isImmediate bool
	var immediateFlags []bool
	var reportlist []models.AmfEventReport

	id, err := amfSelf.EventSubscriptionIDGenerator.Allocate()
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED_NF_FAILURE",
		}
		return nil, problemDetails
	}
	newSubscriptionID := strconv.Itoa(int(id))

	// store subscription in context
	ueEventSubscription := context.AmfUeEventSubscription{}
	extCtxEventSub := models.ExtAmfEventSubscription{
		EventList:                     contextEventSubscription.EventSubscription.EventList,
		EventNotifyUri:                contextEventSubscription.EventSubscription.EventNotifyUri,
		NotifyCorrelationId:           contextEventSubscription.EventSubscription.NotifyCorrelationId,
		NfId:                          contextEventSubscription.EventSubscription.NfId,
		SubsChangeNotifyUri:           contextEventSubscription.EventSubscription.SubsChangeNotifyUri,
		SubsChangeNotifyCorrelationId: contextEventSubscription.EventSubscription.SubsChangeNotifyCorrelationId,
		Supi:                          contextEventSubscription.EventSubscription.Supi,
		GroupId:                       contextEventSubscription.EventSubscription.GroupId,
		ExcludeSupiList:               contextEventSubscription.EventSubscription.ExcludeSupiList,
		ExcludeGpsiList:               contextEventSubscription.EventSubscription.ExcludeGpsiList,
		IncludeSupiList:               contextEventSubscription.EventSubscription.IncludeSupiList,
		IncludeGpsiList:               contextEventSubscription.EventSubscription.IncludeGpsiList,
		Gpsi:                          contextEventSubscription.EventSubscription.Gpsi,
		Pei:                           contextEventSubscription.EventSubscription.Pei,
		AnyUE:                         contextEventSubscription.EventSubscription.AnyUE,
		Options:                       contextEventSubscription.EventSubscription.Options,
		SourceNfType:                  contextEventSubscription.EventSubscription.SourceNfType,
	}
	ueEventSubscription.EventSubscription = &extCtxEventSub
	ueEventSubscription.Timestamp = time.Now().UTC()

	if subscription.Options != nil && subscription.Options.Trigger == models.AmfEventTrigger_CONTINUOUS {
		ueEventSubscription.RemainReports = new(int32)
		*ueEventSubscription.RemainReports = subscription.Options.MaxReports
	}

	if subscription.EventList == nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "SUBSCRIPTION_EMPTY",
		}
		return nil, problemDetails
	}

	for _, events := range subscription.EventList {
		immediateFlags = append(immediateFlags, events.ImmediateFlag)
		if events.ImmediateFlag {
			isImmediate = true
		}
	}

	if subscription.AnyUE {
		contextEventSubscription.IsAnyUe = true
		ueEventSubscription.AnyUe = true
		amfSelf.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*context.AmfUe)
			ue.Lock.Lock()
			ue.EventSubscriptionsInfo[newSubscriptionID] = new(context.AmfUeEventSubscription)
			*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
			contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)
			ue.Lock.Unlock()
			return true
		})
	} else if subscription.GroupId != "" {
		contextEventSubscription.IsGroupUe = true
		ueEventSubscription.AnyUe = true
		amfSelf.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*context.AmfUe)
			ue.Lock.Lock()
			if ue.GroupID == subscription.GroupId {
				ue.EventSubscriptionsInfo[newSubscriptionID] = new(context.AmfUeEventSubscription)
				*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
				contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)
			}
			ue.Lock.Unlock()
			return true
		})
	} else {
		if ue, ok := amfSelf.AmfUeFindBySupi(subscription.Supi); !ok {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "UE_NOT_SERVED_BY_AMF",
			}
			return nil, problemDetails
		} else {
			ue.Lock.Lock()
			ue.EventSubscriptionsInfo[newSubscriptionID] = new(context.AmfUeEventSubscription)
			*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
			contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)
			ue.Lock.Unlock()
		}
	}

	// delete subscription
	if subscription.Options != nil {
		contextEventSubscription.Expiry = subscription.Options.Expiry
	}
	amfSelf.NewEventSubscription(newSubscriptionID, contextEventSubscription)

	// build response

	createdEventSubscription.Subscription = subscription
	createdEventSubscription.SubscriptionId = newSubscriptionID

	// for immediate use
	if subscription.AnyUE {
		amfSelf.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*context.AmfUe)
			ue.Lock.Lock()
			defer ue.Lock.Unlock()

			if isImmediate {
				p.subReports(ue, newSubscriptionID)
			}
			for i, flag := range immediateFlags {
				if flag {
					report, ok := p.newAmfEventReport(ue, subscription.EventList[i].Type, newSubscriptionID)
					if ok {
						reportlist = append(reportlist, report)
					}
				}
			}
			// delete subscription
			if reportlistLen := len(reportlist); reportlistLen > 0 && (!reportlist[reportlistLen-1].State.Active) {
				delete(ue.EventSubscriptionsInfo, newSubscriptionID)
			}
			return true
		})
	} else if subscription.GroupId != "" {
		amfSelf.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*context.AmfUe)
			ue.Lock.Lock()
			defer ue.Lock.Unlock()

			if isImmediate {
				p.subReports(ue, newSubscriptionID)
			}
			if ue.GroupID == subscription.GroupId {
				for i, flag := range immediateFlags {
					if flag {
						report, ok := p.newAmfEventReport(ue, subscription.EventList[i].Type, newSubscriptionID)
						if ok {
							reportlist = append(reportlist, report)
						}
					}
				}
				// delete subscription
				if reportlistLen := len(reportlist); reportlistLen > 0 && (!reportlist[reportlistLen-1].State.Active) {
					delete(ue.EventSubscriptionsInfo, newSubscriptionID)
				}
			}
			return true
		})
	} else {
		ue, _ := amfSelf.AmfUeFindBySupi(subscription.Supi)
		ue.Lock.Lock()
		defer ue.Lock.Unlock()

		if isImmediate {
			p.subReports(ue, newSubscriptionID)
		}
		for i, flag := range immediateFlags {
			if flag {
				report, ok := p.newAmfEventReport(ue, subscription.EventList[i].Type, newSubscriptionID)
				if ok {
					reportlist = append(reportlist, report)
				}
			}
		}
		// delete subscription
		if reportlistLen := len(reportlist); reportlistLen > 0 && (!reportlist[reportlistLen-1].State.Active) {
			delete(ue.EventSubscriptionsInfo, newSubscriptionID)
		}
	}
	if len(reportlist) > 0 {
		createdEventSubscription.ReportList = reportlist
		// delete subscription
		if !reportlist[0].State.Active {
			amfSelf.DeleteEventSubscription(newSubscriptionID)
		}
	}

	return createdEventSubscription, nil
}

func (p *Processor) HandleDeleteAMFEventSubscription(c *gin.Context) {
	logger.EeLog.Infoln("Handle Delete AMF Event Subscription")

	subscriptionID := c.Param("subscriptionId")

	problemDetails := p.DeleteAMFEventSubscriptionProcedure(subscriptionID)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		c.JSON(http.StatusOK, nil)
	}
}

func (p *Processor) DeleteAMFEventSubscriptionProcedure(subscriptionID string) *models.ProblemDetails {
	amfSelf := context.GetSelf()

	subscription, ok := amfSelf.FindEventSubscription(subscriptionID)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "SUBSCRIPTION_NOT_FOUND",
		}
		return problemDetails
	}

	for _, supi := range subscription.UeSupiList {
		if ue, okAmfUeFindBySupi := amfSelf.AmfUeFindBySupi(supi); okAmfUeFindBySupi {
			ue.Lock.Lock()
			delete(ue.EventSubscriptionsInfo, subscriptionID)
			ue.Lock.Unlock()
		}
	}
	amfSelf.DeleteEventSubscription(subscriptionID)
	return nil
}

func (p *Processor) HandleModifyAMFEventSubscription(c *gin.Context,
	modifySubscriptionRequest models.ModifySubscriptionRequest,
) {
	logger.EeLog.Infoln("Handle Modify AMF Event Subscription")

	subscriptionID := c.Param("subscriptionId")

	updatedEventSubscription, problemDetails := p.
		ModifyAMFEventSubscriptionProcedure(subscriptionID, modifySubscriptionRequest)
	if updatedEventSubscription != nil {
		c.JSON(http.StatusOK, updatedEventSubscription)
	} else if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED_NF_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	}
}

func (p *Processor) ModifyAMFEventSubscriptionProcedure(
	subscriptionID string,
	modifySubscriptionRequest models.ModifySubscriptionRequest) (
	*models.AmfUpdatedEventSubscription, *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	contextSubscription, ok := amfSelf.FindEventSubscription(subscriptionID)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "SUBSCRIPTION_NOT_FOUND",
		}
		return nil, problemDetails
	}

	if len(modifySubscriptionRequest.OptionItem) != 0 {
		contextSubscription.Expiry = modifySubscriptionRequest.OptionItem[0].Value
	} else if len(modifySubscriptionRequest.SubscriptionItem) != 0 {
		subscription := &contextSubscription.EventSubscription
		if !contextSubscription.IsAnyUe && !contextSubscription.IsGroupUe {
			if _, okAmfUeFindBySupi := amfSelf.AmfUeFindBySupi(subscription.Supi); !okAmfUeFindBySupi {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusForbidden,
					Cause:  "UE_NOT_SERVED_BY_AMF",
				}
				return nil, problemDetails
			}
		}
		op := modifySubscriptionRequest.SubscriptionItem[0].Op
		index, err := strconv.Atoi(modifySubscriptionRequest.SubscriptionItem[0].Path[11:])
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "UNSPECIFIED_NF_FAILURE",
			}
			return nil, problemDetails
		}
		lists := (subscription.EventList)
		eventlistLen := len(subscription.EventList)
		switch op {
		case "replace":
			event := *modifySubscriptionRequest.SubscriptionItem[0].Value
			if index < eventlistLen {
				(subscription.EventList)[index] = event
			}
		case "remove":
			if index < eventlistLen {
				eventlist := []models.AmfEvent{}
				eventlist = append(eventlist, lists[:index]...)
				eventlist = append(eventlist, lists[index+1:]...)
				subscription.EventList = eventlist
			}
		case "add":
			event := *modifySubscriptionRequest.SubscriptionItem[0].Value
			eventlist := []models.AmfEvent{}
			eventlist = append(eventlist, lists...)
			eventlist = append(eventlist, event)
			subscription.EventList = eventlist
		}
	}

	updatedEventSubscription := &models.AmfUpdatedEventSubscription{
		Subscription: &contextSubscription.EventSubscription,
	}
	return updatedEventSubscription, nil
}

func (p *Processor) subReports(ue *context.AmfUe, subscriptionId string) {
	remainReport := ue.EventSubscriptionsInfo[subscriptionId].RemainReports
	if remainReport == nil {
		return
	}
	*remainReport--
}

// DO NOT handle AmfEventType_PRESENCE_IN_AOI_REPORT and AmfEventType_UES_IN_AREA_REPORT(about area)
func (p *Processor) newAmfEventReport(ue *context.AmfUe, amfEventType models.AmfEventType, subscriptionId string) (
	report models.AmfEventReport, ok bool,
) {
	ueSubscription, ok := ue.EventSubscriptionsInfo[subscriptionId]
	if !ok {
		return report, ok
	}

	report.AnyUe = ueSubscription.AnyUe
	report.Supi = ue.Supi
	report.Type = amfEventType
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
		report.State.Active = p.getDuration(mode.Expiry, &report.State.RemainDuration)
		if report.State.Active {
			report.State.RemainReports = *ueSubscription.RemainReports
		}
	}

	switch amfEventType {
	case models.AmfEventType_LOCATION_REPORT:
		report.Location = &ue.Location
	// case models.AmfEventType_PRESENCE_IN_AOI_REPORT:
	// report.AreaList = (*subscription.EventList)[eventIndex].AreaList
	case models.AmfEventType_TIMEZONE_REPORT:
		report.Timezone = ue.TimeZone
	case models.AmfEventType_ACCESS_TYPE_REPORT:
		for accessType, state := range ue.State {
			if state.Is(context.Registered) {
				report.AccessTypeList = append(report.AccessTypeList, accessType)
			}
		}
	case models.AmfEventType_REGISTRATION_STATE_REPORT:
		var rmInfos []models.RmInfo
		for accessType, state := range ue.State {
			rmInfo := models.RmInfo{
				RmState:    models.RmState_DEREGISTERED,
				AccessType: accessType,
			}
			if state.Is(context.Registered) {
				rmInfo.RmState = models.RmState_REGISTERED
			}
			rmInfos = append(rmInfos, rmInfo)
		}
		report.RmInfoList = rmInfos
	case models.AmfEventType_CONNECTIVITY_STATE_REPORT:
		report.CmInfoList = ue.GetCmInfo()
	case models.AmfEventType_REACHABILITY_REPORT:
		report.Reachability = ue.Reachability
	case models.AmfEventType_COMMUNICATION_FAILURE_REPORT:
		// TODO : report.CommFailure
	case models.AmfEventType_SUBSCRIPTION_ID_CHANGE:
		report.SubscriptionId = subscriptionId
	case models.AmfEventType_SUBSCRIPTION_ID_ADDITION:
		report.SubscriptionId = subscriptionId
	}
	return report, ok
}

func (p *Processor) getDuration(expiry *time.Time, remainDuration *int32) bool {
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
