package smf_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/smf/factory"
	"time"
)

var NFServices *[]models.NfService

var NfServiceVersion *[]models.NfServiceVersion

var SmfInfo *models.SmfInfo

func SetupNFProfile(config *factory.Config) {
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
	NFServices = new([]models.NfService)
	for _, serviceName := range config.Configuration.ServiceNameList {
		*NFServices = append(*NFServices, models.NfService{
			ServiceInstanceId: SMF_Self().NfInstanceID + serviceName,
			ServiceName:       models.ServiceName(serviceName),
			Versions:          NfServiceVersion,
			Scheme:            models.UriScheme_HTTPS,
			NfServiceStatus:   models.NfServiceStatus_REGISTERED,
			ApiPrefix:         fmt.Sprintf("%s://%s:%d", SMF_Self().URIScheme, SMF_Self().HTTPAddress, SMF_Self().HTTPPort),
		})
	}

	//set smfInfo
	SmfInfo = &models.SmfInfo{
		SNssaiSmfInfoList: &smfContext.SnssaiInfos,
	}
}
