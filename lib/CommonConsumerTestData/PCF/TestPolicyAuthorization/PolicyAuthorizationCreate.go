package TestPolicyAuthorization

import (
	"free5gc/lib/openapi/models"
)

func GetPostAppSessions201Data() models.AppSessionContext {
	PostAppSessions201Data := models.AppSessionContext{
		AscReqData: &models.AppSessionContextReqData{
			AfRoutReq: &models.AfRoutingRequirement{},
			MedComponents: map[string]models.MediaComponent{
				"mediacomponent1": {
					FStatus: models.FlowStatus_ENABLED,
				},
			},
			EvSubsc: &models.EventsSubscReqData{
				Events: []models.AfEventSubscription{
					{
						Event:       "ACCESS_TYPE_CHANGE",
						NotifMethod: "EVENT_DETECTION",
					},
				},
			},
			NotifUri: "NotifUri",
			SuppFeat: "0",
			Supi:     "string1",
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
		AfAppId:  "AfAppId",
		AspId:    "AspId",
		BdtRefId: "BdtRefId",
		AfRoutReq: &models.AfRoutingRequirementRm{
			AppReloc: true,
			RouteToLocs: []models.RouteToLocation{
				{
					Dnai: "Dnai",
					RouteInfo: &models.RouteInformation{
						Ipv4Addr:   "111.11.11.1",
						Ipv6Addr:   "222.22.22.2",
						PortNumber: 9999,
					},
					RouteProfId: "RouteProfId",
				},
			},
			SpVal: &models.SpatialValidityRm{
				PresenceInfoList: map[string]models.PresenceInfo{
					"additionalProp1": {
						PraId: "11",
					},
				},
			},
			TempVals: []models.TemporalValidity{
				{},
			},
			UpPathChgSub: &models.UpPathChgEvent{},
		},
		EvSubsc: &models.EventsSubscReqDataRm{
			NotifUri: "EvSubsc_NotifUri",
			Events: []models.AfEventSubscription{
				{
					Event:       "ACCESS_TYPE_CHANGE",
					NotifMethod: "EVENT_DETECTION",
				},
			},
			UsgThres: &models.UsageThresholdRm{
				Duration:       1,
				TotalVolume:    1,
				DownlinkVolume: 1,
				UplinkVolume:   1},
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
