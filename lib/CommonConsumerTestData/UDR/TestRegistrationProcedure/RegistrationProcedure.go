package TestRegistrationProcedure

import (
	"free5gc/lib/openapi/models"
)

const (
	FREE5GC_CASE = "free5gc"
)

var TestAmDataTable = make(map[string]models.AccessAndMobilitySubscriptionData)
var TestSmfSelDataTable = make(map[string]models.SmfSelectionSubscriptionData)

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
						Dnn: "internet",
					},
				},
			},
		},
	}
}
