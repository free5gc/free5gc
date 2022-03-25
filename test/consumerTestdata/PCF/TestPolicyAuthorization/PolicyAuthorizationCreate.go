package TestPolicyAuthorization

import (
	"github.com/free5gc/openapi/models"
)

func GetPostAppSessionsData_Normal() models.AppSessionContext {
	PostAppSessionsData := models.AppSessionContext{
		AscReqData: &models.AppSessionContextReqData{
			AfRoutReq: &models.AfRoutingRequirement{},
			Dnn:       "internet",
			SliceInfo: &models.Snssai{
				Sst: 1,
				Sd:  "010203",
			},
			MedComponents: map[string]models.MediaComponent{
				"1": {
					MedCompN: 1,
					MarBwDl:  "400 Mbps",
					MarBwUl:  "400 Mbps",
					MirBwDl:  "20 Mbps",
					MirBwUl:  "20 Mbps",
					MedType:  models.MediaType_AUDIO,
					FStatus:  models.FlowStatus_ENABLED,
					MedSubComps: map[string]models.MediaSubComponent{
						"1": {
							FNum:    1,
							FDescs:  []string{"permit out ip from 127.0.0.1 to 45.45.0.2"},
							FStatus: models.FlowStatus_ENABLED,
						},
					},
				},
			},
			EvSubsc: &models.EventsSubscReqData{
				Events: []models.AfEventSubscription{
					{
						Event:       models.AfEvent_ACCESS_TYPE_CHANGE,
						NotifMethod: models.AfNotifMethod_EVENT_DETECTION,
					},
					{
						Event: models.AfEvent_QOS_NOTIF,
					},
					{
						Event: models.AfEvent_PLMN_CHG,
					},
					{
						Event: models.AfEvent_FAILED_RESOURCES_ALLOCATION,
					},
					{
						Event: models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION,
					},
					{
						Event: models.AfEvent_USAGE_REPORT,
					},
				},
				NotifUri: "https://127.0.0.1:12345",
				UsgThres: &models.UsageThreshold{
					Duration:    100,
					TotalVolume: 30000,
				},
			},
			NotifUri: "https://127.0.0.1:12345",
			SuppFeat: "5", //b'0111'
			Supi:     "imsi-2089300007487",
			UeIpv4:   "45.45.0.2",
		},
	}
	return PostAppSessionsData
}

func GetPostAppSessionsData_AFInfluenceOnTrafficRouting() models.AppSessionContext {
	PostAppSessionsData := models.AppSessionContext{
		AscReqData: &models.AppSessionContextReqData{
			AfAppId:  "edge",
			Dnn:      "internet",
			SuppFeat: "03",
			Supi:     "imsi-2089300007487",
			UeIpv4:   "10.60.0.1",
			SliceInfo: &models.Snssai{
				Sst: 1, Sd: "010203",
			},
			AfRoutReq: &models.AfRoutingRequirement{
				AppReloc: false,
				UpPathChgSub: &models.UpPathChgEvent{
					DnaiChgType:     models.DnaiChangeType_LATE,
					NotificationUri: "http://127.0.0.100:8000/nnef-callback/v1/traffic-influence/edge",
					NotifCorreId:    "1234",
				},
				RouteToLocs: []models.RouteToLocation{
					{
						Dnai:        "edge",
						RouteProfId: "MEC1",
					},
				},
			},
			NotifUri: "http://127.0.0.100:8000/nnef-callback/v1/applications/edge",
			IpDomain: "edgeIPDomain",
		},
	}
	return PostAppSessionsData
}

func GetPostAppSessionsData_Flow3() models.AppSessionContext {
	PostAppSessionsData := GetPostAppSessionsData_Normal()
	medComp := PostAppSessionsData.AscReqData.MedComponents["1"]
	medComp.MedSubComps["2"] = models.MediaSubComponent{
		FNum:    2,
		FDescs:  []string{"permit in ip from 127.0.0.2 to 45.45.0.2"},
		MarBwDl: "200 Mbps",
		FStatus: models.FlowStatus_ENABLED,
	}
	medComp.MedSubComps["3"] = models.MediaSubComponent{
		FNum:    3,
		FDescs:  []string{"permit inout ip from 127.0.0.3 to 45.45.0.2"},
		MarBwDl: "500 Mbps",
		FStatus: models.FlowStatus_ENABLED,
	}
	PostAppSessionsData.AscReqData.MedComponents["1"] = medComp
	return PostAppSessionsData
}

func GetPostAppSessionsData_403Forbidden() models.AppSessionContext {
	PostAppSessionsData := GetPostAppSessionsData_Normal()
	medComp := PostAppSessionsData.AscReqData.MedComponents["1"]
	medComp.MedSubComps["1"] = models.MediaSubComponent{
		FNum:    1,
		FDescs:  []string{"permit in ip from 127.0.0.4 to 45.45.0.2"},
		FStatus: models.FlowStatus_ENABLED,
	}
	medComp.MirBwUl = "500 Mbps"
	PostAppSessionsData.AscReqData.MedComponents["1"] = medComp
	return PostAppSessionsData
}

func GetPostAppSessionsData_400() models.AppSessionContext {
	PostAppSessionsData := GetPostAppSessionsData_Normal()
	PostAppSessionsData.AscReqData.MedComponents = nil
	return PostAppSessionsData
}

func GetPostAppSessionsData_NoEvent() models.AppSessionContext {
	PostAppSessionsData := GetPostAppSessionsData_Normal()
	PostAppSessionsData.AscReqData.EvSubsc = nil
	return PostAppSessionsData
}

func GetDeleteAppSession204Data() models.AppSessionContext {
	DeleteAppSession204Data := models.AppSessionContext{
		AscReqData: &models.AppSessionContextReqData{
			Supi:     "123",
			NotifUri: "https://127.0.0.1:12345",
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
			{
				Event:       models.AfEvent_ACCESS_TYPE_CHANGE,
				NotifMethod: models.AfNotifMethod_EVENT_DETECTION,
			},
			{
				Event:       models.AfEvent_PLMN_CHG,
				NotifMethod: models.AfNotifMethod_EVENT_DETECTION,
			},
		},
		NotifUri: "https://127.0.0.1:12345",
	}
	return UpdateEventsSubsc201Data
}

func GetUpdateEventsSubsc200Data() models.EventsSubscReqData {
	UpdateEventsSubsc200Data := models.EventsSubscReqData{
		Events: []models.AfEventSubscription{
			{
				Event:       models.AfEvent_PLMN_CHG,
				NotifMethod: models.AfNotifMethod_EVENT_DETECTION,
			},
		},
		NotifUri: "https://127.0.0.1:12345",
	}
	return UpdateEventsSubsc200Data
}

func GetUpdateEventsSubsc204Data() models.EventsSubscReqData {
	UpdateEventsSubsc204Data := models.EventsSubscReqData{
		Events: []models.AfEventSubscription{
			{
				Event:       models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION,
				NotifMethod: models.AfNotifMethod_EVENT_DETECTION,
			},
		},
		NotifUri: "https://127.0.0.1:12345",
	}
	return UpdateEventsSubsc204Data
}

func GetUpdateEventsSubsc400Data() models.EventsSubscReqData {
	UpdateEventsSubsc400Data := models.EventsSubscReqData{
		UsgThres: &models.UsageThreshold{
			Duration:       0,
			TotalVolume:    0,
			DownlinkVolume: 0,
			UplinkVolume:   0},
	}
	return UpdateEventsSubsc400Data
}

func GetModAppSession200Data() models.AppSessionContextUpdateData {
	ModAppSession200Data := models.AppSessionContextUpdateData{
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
			UpPathChgSub: &models.UpPathChgEvent{},
		},
		EvSubsc: &models.EventsSubscReqDataRm{
			NotifUri: "EvSubsc_NotifUri",
			Events: []models.AfEventSubscription{
				{
					Event:       models.AfEvent_ACCESS_TYPE_CHANGE,
					NotifMethod: models.AfNotifMethod_EVENT_DETECTION,
				},
			},
			UsgThres: &models.UsageThresholdRm{
				Duration:    10,
				TotalVolume: 10,
			},
		},
		MedComponents: map[string]models.MediaComponentRm{
			"1": {
				MedCompN: 1,
				MarBwDl:  "40 Mbps",
				MarBwUl:  "40 Mbps",
				MirBwDl:  "20 Mbps",
				MirBwUl:  "20 Mbps",
				MedType:  models.MediaType_AUDIO,
				FStatus:  models.FlowStatus_ENABLED,
				MedSubComps: map[string]models.MediaSubComponentRm{
					"1": {
						FNum:    1,
						FDescs:  []string{"permit out ip from 127.0.0.9 to 45.45.0.2"},
						FStatus: models.FlowStatus_ENABLED,
					},
				},
			},
		},
	}
	return ModAppSession200Data
}

func GetModAppSession403Data() models.AppSessionContextUpdateData {
	ModAppSession403Data := models.AppSessionContextUpdateData{}
	return ModAppSession403Data
}
