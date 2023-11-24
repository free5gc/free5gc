package TestRegistrationProcedure

import (
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/webconsole/backend/WebUI"
)

const (
	FREE5GC_CASE = "free5gc"
)

var TestAmDataTable = make(map[string]models.AccessAndMobilitySubscriptionData)
var TestSmfSelDataTable = make(map[string]models.SmfSelectionSubscriptionData)
var TestSmSelDataTable = make(map[string][]models.SessionManagementSubscriptionData)
var TestAmPolicyDataTable = make(map[string]models.AmPolicyData)
var TestSmPolicyDataTable = make(map[string]models.SmPolicyData)
var TestChargingDataTable = make(map[string][]WebUI.ChargingData)
var TestFlowRuleTable = make(map[string][]WebUI.FlowRule)
var TestQoSFlowTable = make(map[string][]WebUI.QosFlow)

func init() {
	TestAmDataTable[FREE5GC_CASE] = models.AccessAndMobilitySubscriptionData{
		Gpsis: []string{
			"msisdn-0900000000",
		},
		Nssai: &models.Nssai{
			DefaultSingleNssais: []models.Snssai{
				{
					Sst: 1,
					Sd:  "FEDCBA",
				},
				{
					Sst: 1,
					Sd:  "112233",
				},
			},
			SingleNssais: []models.Snssai{
				{
					Sst: 1,
					Sd:  "FEDCBA",
				},
				{
					Sst: 1,
					Sd:  "112233",
				},
			},
		},
		SubscribedUeAmbr: &models.AmbrRm{
			Uplink:   "1000 Kbps",
			Downlink: "1000 Kbps",
		},
	}

	TestSmfSelDataTable[FREE5GC_CASE] = models.SmfSelectionSubscriptionData{
		SubscribedSnssaiInfos: map[string]models.SnssaiInfo{
			"01FEDCBA": { // sst:1, sd:FEDCBA
				DnnInfos: []models.DnnInfo{
					{
						Dnn: "internet",
					},
				},
			},
			"01112233": { // sst:1, sd:112233
				DnnInfos: []models.DnnInfo{
					{
						Dnn: "internet2",
					},
				},
			},
		},
	}

	TestAmPolicyDataTable[FREE5GC_CASE] = models.AmPolicyData{
		SubscCats: []string{
			"free5gc",
		},
	}

	TestSmPolicyDataTable[FREE5GC_CASE] = models.SmPolicyData{
		SmPolicySnssaiData: map[string]models.SmPolicySnssaiData{
			"01FEDCBA": {
				Snssai: &models.Snssai{
					Sd:  "FEDCBA",
					Sst: 1,
				},
				SmPolicyDnnData: map[string]models.SmPolicyDnnData{
					"internet": {
						Dnn: "internet",
					},
				},
			},
			"01112233": {
				Snssai: &models.Snssai{
					Sd:  "112233",
					Sst: 1,
				},
				SmPolicyDnnData: map[string]models.SmPolicyDnnData{
					"internet2": {
						Dnn: "internet2",
					},
				},
			},
		},
	}

	TestSmSelDataTable[FREE5GC_CASE] = []models.SessionManagementSubscriptionData{
		{
			SingleNssai: &models.Snssai{
				Sst: 1,
				Sd:  "FEDCBA",
			},
			DnnConfigurations: map[string]models.DnnConfiguration{
				"internet": {
					SscModes: &models.SscModes{
						DefaultSscMode:  models.SscMode__1,
						AllowedSscModes: []models.SscMode{models.SscMode__1, models.SscMode__2, models.SscMode__3},
					},
					PduSessionTypes: &models.PduSessionTypes{DefaultSessionType: models.PduSessionType_IPV4,
						AllowedSessionTypes: []models.PduSessionType{models.PduSessionType_IPV4},
					},
					SessionAmbr: &models.Ambr{
						Uplink:   "1000 Kbps",
						Downlink: "1000 Kbps",
					},
					Var5gQosProfile: &models.SubscribedDefaultQos{
						Var5qi: 9,
						Arp: &models.Arp{
							PriorityLevel: 8,
						},
						PriorityLevel: 8,
					},
				},
			},
		},
		{
			SingleNssai: &models.Snssai{
				Sst: 1,
				Sd:  "112233",
			},
			DnnConfigurations: map[string]models.DnnConfiguration{
				"internet2": {
					SscModes: &models.SscModes{
						DefaultSscMode:  models.SscMode__1,
						AllowedSscModes: []models.SscMode{models.SscMode__1, models.SscMode__2, models.SscMode__3},
					},
					PduSessionTypes: &models.PduSessionTypes{DefaultSessionType: models.PduSessionType_IPV4,
						AllowedSessionTypes: []models.PduSessionType{models.PduSessionType_IPV4},
					},
					SessionAmbr: &models.Ambr{
						Uplink:   "1000 Kbps",
						Downlink: "1000 Kbps",
					},
					Var5gQosProfile: &models.SubscribedDefaultQos{
						Var5qi: 9,
						Arp: &models.Arp{
							PriorityLevel: 8,
						},
						PriorityLevel: 8,
					},
				},
			},
		},
	}

	TestChargingDataTable[FREE5GC_CASE] = []WebUI.ChargingData{
		{
			Snssai:         "01FEDCBA",
			Dnn:            "",
			Filter:         "",
			QosRef:         3,
			ChargingMethod: "Online",
			Quota:          "100000",
			UnitCost:       "1",
		},
		{
			Snssai:         "01FEDCBA",
			Dnn:            "internet",
			Filter:         "10.10.0.84/32",
			QosRef:         1,
			ChargingMethod: "Online",
			Quota:          "100000",
			UnitCost:       "2",
		},
		{
			Snssai:         "01112233",
			Dnn:            "",
			Filter:         "",
			QosRef:         4,
			ChargingMethod: "Online",
			Quota:          "100000",
			UnitCost:       "2",
		},
		{
			Snssai:         "01112233",
			Dnn:            "internet",
			Filter:         "10.10.0.84/32",
			QosRef:         2,
			ChargingMethod: "Online",
			Quota:          "100000",
			UnitCost:       "4",
		},
	}

	TestFlowRuleTable[FREE5GC_CASE] = []WebUI.FlowRule{
		{
			Filter:     "10.10.0.84/32",
			Precedence: 128,
			Snssai:     "01FEDCBA",
			Dnn:        "internet",
			QosRef:     1,
		},
		{
			Filter:     "10.10.0.84/32",
			Precedence: 128,
			Snssai:     "01112233",
			Dnn:        "internet",
			QosRef:     2,
		},
	}

	TestQoSFlowTable[FREE5GC_CASE] = []WebUI.QosFlow{
		{
			Snssai: "01FEDCBA",
			Dnn:    "internet",
			MBRUL:  "208 Mbps",
			MBRDL:  "208 Mbps",
			GBRUL:  "108 Mbps",
			GBRDL:  "108 Mbps",
			QosRef: 1,
			Var5QI: 8,
		},
		{
			Snssai: "01112233",
			Dnn:    "internet",
			MBRUL:  "208 Mbps",
			MBRDL:  "208 Mbps",
			GBRUL:  "108 Mbps",
			GBRDL:  "108 Mbps",
			QosRef: 2,
			Var5QI: 7,
		},
	}
}
