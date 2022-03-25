package TestBDTPolicy

import (
	"time"

	"github.com/free5gc/openapi/models"
)

func GetCreateTestData() models.BdtReqData {
	startTime := time.Now().Format(time.RFC3339)
	mm, _ := time.ParseDuration("10m")
	stopTime := time.Now().Add(mm).Format(time.RFC3339)
	bdtReqData := models.BdtReqData{
		AspId: "123456",
		DesTimeInt: &models.TimeWindow{
			StartTime: startTime,
			StopTime:  stopTime,
		},
		NumOfUes: 1,
		VolPerUe: &models.UsageThreshold{
			Duration:       1,
			TotalVolume:    1,
			DownlinkVolume: 1,
			UplinkVolume:   1,
		},
		NwAreaInfo: &models.NetworkAreaInfo{
			Tais: []models.Tai{
				{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					Tac: "000001",
				},
			},
			Ncgis: []models.Ncgi{
				{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					NrCellId: "000000001",
				},
			},
			GRanNodeIds: []models.GlobalRanNodeId{
				{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					GNbId: &models.GNbId{
						BitLength: 24,
						GNBValue:  "000001",
					},
				},
			},
		},
		SuppFeat: "",
	}

	return bdtReqData
}

func GetUpdateTestData() models.BdtPolicyDataPatch {
	bdtPolicyDataPatch := models.BdtPolicyDataPatch{
		SelTransPolicyId: 1,
	}
	return bdtPolicyDataPatch
}
