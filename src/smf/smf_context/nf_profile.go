package smf_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"time"
)

var NFServices *[]models.NfService

var NfServiceVersion *[]models.NfServiceVersion

var SmfInfo *models.SmfInfo

func SetupNFProfile() {
	//Set time
	date := time.Now()
	dateFormat, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))

	//set NfServiceVersion
	NfServiceVersion = &[]models.NfServiceVersion{
		{
			ApiVersionInUri: "v1",
			ApiFullVersion:  fmt.Sprintf("https://%s:%d/nsmf-pdusession/v1", SMF_Self().HTTPAddress, SMF_Self().HTTPPort),
			Expiry:          &dateFormat,
		},
	}

	//set NFServices
	NFServices = &[]models.NfService{
		{
			ServiceInstanceId: SMF_Self().NfInstanceID + "nsmf-pdusession",
			ServiceName:       models.ServiceName_NSMF_PDUSESSION,
			Versions:          NfServiceVersion,
			Scheme:            models.UriScheme_HTTPS,
			NfServiceStatus:   models.NfServiceStatus_REGISTERED,
			ApiPrefix:         fmt.Sprintf("%s://%s:%d", SMF_Self().URIScheme, SMF_Self().HTTPAddress, SMF_Self().HTTPPort),
		},
	}

	//set smfInfo
	SmfInfo = &models.SmfInfo{
		SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
			{
				SNssai: &models.Snssai{},
				DnnSmfInfoList: &[]models.DnnSmfInfoItem{
					{
						Dnn: "internet",
					},
				},
			},
		},
	}
}
