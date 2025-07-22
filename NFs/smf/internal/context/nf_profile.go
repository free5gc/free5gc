package context

import (
	"fmt"
	"time"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/pkg/factory"
)

type NFProfile struct {
	NFServices       *[]models.NrfNfManagementNfService
	NFServiceVersion *[]models.NfServiceVersion
	SMFInfo          *models.SmfInfo
	PLMNList         *[]models.PlmnId
}

func (c *SMFContext) SetupNFProfile(nfProfileconfig *factory.Config) {
	// Set time
	nfSetupTime := time.Now()

	// set NfServiceVersion
	c.NfProfile.NFServiceVersion = &[]models.NfServiceVersion{
		{
			ApiVersionInUri: "v1",
			ApiFullVersion: fmt.
				Sprintf("https://%s:%d"+factory.SmfPdusessionResUriPrefix, GetSelf().RegisterIPv4, GetSelf().SBIPort),
			Expiry: &nfSetupTime,
		},
	}

	// set NFServices
	c.NfProfile.NFServices = new([]models.NrfNfManagementNfService)
	for _, serviceName := range nfProfileconfig.Configuration.ServiceNameList {
		*c.NfProfile.NFServices = append(*c.NfProfile.NFServices, models.NrfNfManagementNfService{
			ServiceInstanceId: GetSelf().NfInstanceID + serviceName,
			ServiceName:       models.ServiceName(serviceName),
			Versions:          *c.NfProfile.NFServiceVersion,
			Scheme:            models.UriScheme_HTTPS,
			NfServiceStatus:   models.NfServiceStatus_REGISTERED,
			ApiPrefix:         fmt.Sprintf("%s://%s:%d", GetSelf().URIScheme, GetSelf().RegisterIPv4, GetSelf().SBIPort),
			IpEndPoints: []models.IpEndPoint{
				{
					Ipv4Address: GetSelf().RegisterIPv4,
					Port:        int32(GetSelf().SBIPort),
				},
			},
		})
	}

	// set smfInfo
	c.NfProfile.SMFInfo = &models.SmfInfo{
		SNssaiSmfInfoList: SNssaiSmfInfo(),
	}

	// set PlmnList if exists
	if plmnList := nfProfileconfig.Configuration.PLMNList; plmnList != nil {
		c.NfProfile.PLMNList = new([]models.PlmnId)
		for _, plmn := range plmnList {
			*c.NfProfile.PLMNList = append(*c.NfProfile.PLMNList, models.PlmnId{
				Mcc: plmn.Mcc,
				Mnc: plmn.Mnc,
			})
		}
	}
}

func SNssaiSmfInfo() []models.SnssaiSmfInfoItem {
	snssaiInfo := make([]models.SnssaiSmfInfoItem, 0)
	for _, snssai := range smfContext.SnssaiInfos {
		var snssaiInfoModel models.SnssaiSmfInfoItem
		snssaiInfoModel.SNssai = &models.ExtSnssai{
			Sst: snssai.Snssai.Sst,
			Sd:  snssai.Snssai.Sd,
		}
		dnnModelList := make([]models.DnnSmfInfoItem, 0)

		for dnn := range snssai.DnnInfos {
			dnnModelList = append(dnnModelList, models.DnnSmfInfoItem{
				Dnn: dnn,
			})
		}

		snssaiInfoModel.DnnSmfInfoList = dnnModelList

		snssaiInfo = append(snssaiInfo, snssaiInfoModel)
	}
	return snssaiInfo
}
