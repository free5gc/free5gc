package TestAmf

import (
	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
)

func BuildSmfNfProfile() (uuId string, profile models.NfProfile) {
	uuId = uuid.New().String()
	profile = models.NfProfile{
		NfInstanceId: uuId,
		NfType:       models.NfType_SMF,
		NfStatus:     models.NfStatus_REGISTERED,
		SNssais: &[]models.Snssai{
			{
				Sst: 1,
				Sd:  "010203",
			},
		},
		PlmnList: &[]models.PlmnId{
			{
				Mcc: "208",
				Mnc: "93",
			},
		},
		NfServices: &[]models.NfService{
			{

				ServiceInstanceId: "1",
				ServiceName:       models.ServiceName_NSMF_PDUSESSION,
				Scheme:            models.UriScheme_HTTPS,
				NfServiceStatus:   models.NfServiceStatus_REGISTERED,
				Versions: &[]models.NfServiceVersion{
					{
						ApiVersionInUri: "v1",
						ApiFullVersion:  "1.0.0",
					},
				},
				ApiPrefix: "https://localhost:29502",
				IpEndPoints: &[]models.IpEndPoint{
					{
						Ipv4Address: "127.0.0.1",
						Port:        29502,
					},
				},
			},
		},
		SmfInfo: &models.SmfInfo{
			SNssaiSmfInfoList: &[]models.SnssaiSmfInfoItem{
				{
					SNssai: &models.Snssai{
						Sst: 1,
						Sd:  "010203",
					},
					DnnSmfInfoList: &[]models.DnnSmfInfoItem{
						{
							Dnn: "internet",
						},
					},
				},
			},
		},
	}
	return

}
