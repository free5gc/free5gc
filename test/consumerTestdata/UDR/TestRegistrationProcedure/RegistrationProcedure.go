package TestRegistrationProcedure

import (
	"github.com/free5gc/openapi/models"
)

const (
	FREE5GC_CASE = "free5gc"
)

var TestAmDataTable = make(map[string]models.AccessAndMobilitySubscriptionData)
var TestSmfSelDataTable = make(map[string]models.SmfSelectionSubscriptionData)
var TestSmSelDataTable = make(map[string][]models.SessionManagementSubscriptionData)
var TestAmPolicyDataTable = make(map[string]models.AmPolicyData)
var TestSmPolicyDataTable = make(map[string]models.SmPolicyData)

func init() {
	TestAmDataTable[FREE5GC_CASE] = models.AccessAndMobilitySubscriptionData{
		Gpsis: []string{
			"msisdn-0900000000",
		},
		Nssai: &models.Nssai{
			DefaultSingleNssais: []models.Snssai{
				{
					Sst: 1,
					Sd:  "010203",
				},
				{
					Sst: 1,
					Sd:  "112233",
				},
			},
			SingleNssais: []models.Snssai{
				{
					Sst: 1,
					Sd:  "010203",
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
			"01010203": { // sst:1, sd:010203
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
			"01010203": {
				Snssai: &models.Snssai{
					Sd:  "010203",
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
				Sd:  "010203",
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
}
