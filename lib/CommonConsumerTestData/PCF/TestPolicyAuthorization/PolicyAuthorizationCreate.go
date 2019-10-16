package TestPolicyAuthorization

import (
	"free5gc/lib/openapi/models"
)

func GetPostAppSessions201Data() models.AppSessionContext {
	PostAppSessions201Data := models.AppSessionContext{
		AscReqData: &models.AppSessionContextReqData{
			Supi:     "123",
			NotifUri: "NotifUri",
			SuppFeat: "0",
		},
		AscRespData: &models.AppSessionContextRespData{},
		EvsNotif: &models.EventsNotification{
			EvSubsUri: "EvSubsUri",
			EvNotifs: []models.AfEventNotification{
				{Event: "USAGE_REPORT"},
			},
		},
	}
	return PostAppSessions201Data
}

func GetPostAppSessions403Data() models.AppSessionContext {
	PostAppSessions403Data := models.AppSessionContext{
		AscReqData:  &models.AppSessionContextReqData{},
		AscRespData: &models.AppSessionContextRespData{},
		EvsNotif:    &models.EventsNotification{},
	}
	return PostAppSessions403Data
}

func GetDeleteAppSession204Data() models.AppSessionContext {
	DeleteAppSession204Data := models.AppSessionContext{
		AscReqData: &models.AppSessionContextReqData{
			Supi:     "123",
			NotifUri: "NotifUri",
			SuppFeat: "0",
		},
		AscRespData: &models.AppSessionContextRespData{},
		EvsNotif:    &models.EventsNotification{},
	}
	return DeleteAppSession204Data
}

func GetUpdateEventsSubsc201Data() models.EventsSubscReqData {
	UpdateEventsSubsc201Data := models.EventsSubscReqData{
		Events: []models.AfEventSubscription{
			{Event: "ACCESS_TYPE_CHANGE",
				NotifMethod: "EVENT_DETECTION"},
		},
		NotifUri: "Put_NotifUri",
		UsgThres: &models.UsageThreshold{
			Duration:       0,
			TotalVolume:    0,
			DownlinkVolume: 0,
			UplinkVolume:   0},
	}
	return UpdateEventsSubsc201Data
}

func GetUpdateEventsSubsc200Data() models.EventsSubscReqData {
	UpdateEventsSubsc200Data := models.EventsSubscReqData{
		Events: []models.AfEventSubscription{
			{Event: "FAILED_RESOURCES_ALLOCATION",
				NotifMethod: "ONE_TIME"},
		},
		NotifUri: "Put_NotifUri",
		UsgThres: &models.UsageThreshold{
			Duration:       0,
			TotalVolume:    0,
			DownlinkVolume: 0,
			UplinkVolume:   0},
	}
	return UpdateEventsSubsc200Data
}

func GetUpdateEventsSubsc403Data() models.EventsSubscReqData {
	UpdateEventsSubsc403Data := models.EventsSubscReqData{
		UsgThres: &models.UsageThreshold{
			Duration:       0,
			TotalVolume:    0,
			DownlinkVolume: 0,
			UplinkVolume:   0},
	}
	return UpdateEventsSubsc403Data
}

func GetModAppSession200Data() models.AppSessionContextUpdateData {
	ModAppSession200Data := models.AppSessionContextUpdateData{
		AfAppId:  "string",
		AspId:    "AspId",
		BdtRefId: "BdtRefId",
		EvSubsc: &models.EventsSubscReqDataRm{
			NotifUri: "evSubsc_notifUri",
			Events: []models.AfEventSubscription{
				{
					Event:       "ACCESS_TYPE_CHANGE",
					NotifMethod: "EVENT_DETECTION",
				},
			},
		},
		MedComponents: map[string]models.MediaComponentRm{
			"additionalProp1": {
				AfAppId: "med_afAppId",
				AfRoutReq: &models.AfRoutingRequirementRm{
					SpVal: &models.SpatialValidityRm{
						PresenceInfoList: map[string]models.PresenceInfo{
							"additionalProp1": {
								PraId: "11",
							},
						},
					},
				},
				MedCompN: 1,
			},
		},
		MpsId:  "MpsId",
		SponId: "SponId",
	}
	return ModAppSession200Data
}

func GetModAppSession403Data() models.AppSessionContextUpdateData {
	ModAppSession403Data := models.AppSessionContextUpdateData{}
	return ModAppSession403Data
}
