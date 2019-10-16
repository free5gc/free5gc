package TestBDTPolicy

import (
	"free5gc/lib/openapi/models"
	"time"
)

func GetCreateTestData() (models.BdtReqData, models.BdtReqData, models.BdtReqData, models.BdtReqData) {
	startTime := time.Now()
	mm, _ := time.ParseDuration("10m")
	stopTime := time.Now().Add(mm)
	bdtReqData1 := models.BdtReqData{
		AspId: "123456",
		DesTimeInt: &models.TimeWindow{
			StartTime: &startTime,
			StopTime:  &stopTime,
		},
		NumOfUes: 1,
		VolPerUe: &models.UsageThreshold{
			Duration:       1,
			TotalVolume:    1,
			DownlinkVolume: 1,
			UplinkVolume:   1,
		},
		SuppFeat: "123456",
	}

	bdtReqData2 := models.BdtReqData{
		AspId: "456789",
		DesTimeInt: &models.TimeWindow{
			StartTime: &startTime,
			StopTime:  &stopTime,
		},
		NumOfUes: 1,
		VolPerUe: &models.UsageThreshold{
			Duration:       2,
			TotalVolume:    2,
			DownlinkVolume: 2,
			UplinkVolume:   2,
		},
		SuppFeat: "9876545",
	}

	bdtReqData3 := models.BdtReqData{
		AspId: "1012",
		DesTimeInt: &models.TimeWindow{
			StartTime: &startTime,
			StopTime:  &stopTime,
		},
		NumOfUes: 2,
		VolPerUe: &models.UsageThreshold{
			Duration:       3,
			TotalVolume:    3,
			DownlinkVolume: 2,
			UplinkVolume:   2,
		},
		SuppFeat: "000011",
	}

	bdtReqDataNil := models.BdtReqData{
		AspId: "",
		DesTimeInt: &models.TimeWindow{
			StartTime: nil,
			StopTime:  nil,
		},
		NumOfUes: 0,
		VolPerUe: &models.UsageThreshold{
			Duration:       0,
			TotalVolume:    0,
			DownlinkVolume: 0,
			UplinkVolume:   0,
		},
		SuppFeat: "",
	}

	return bdtReqData1, bdtReqData2, bdtReqData3, bdtReqDataNil
}

func GetUpdateTestData() models.BdtPolicyDataPatch {
	bdtPolicyDataPatch := models.BdtPolicyDataPatch{
		SelTransPolicyId: 1,
	}
	return bdtPolicyDataPatch
}
