package pcf_consumer

import (
	"context"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/pcf_context"
	"net/http"
	"strings"
)

func BuildNFInstance(context *pcf_context.PCFContext) (profile models.NfProfile, err error) {
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_PCF
	profile.NfStatus = models.NfStatus_REGISTERED
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, context.HttpIPv4Address)
	service := []models.NfService{}
	for _, nfService := range context.NfService {
		service = append(service, nfService)
	}
	profile.NfServices = &service
	profile.PcfInfo = &models.PcfInfo{
		DnnList: []string{
			"free5gc",
			"internet",
		},
		// SupiRanges: &[]models.SupiRange{
		// 	{
		// 		//from TS 29.510 6.1.6.2.9 example2
		//		//no need to set supirange in this moment 2019/10/4
		// 		Start:   "123456789040000",
		// 		End:     "123456789059999",
		// 		Pattern: "^imsi-12345678904[0-9]{4}$",
		// 	},
		// },
	}
	return
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (resouceNrfUri string, err error) {

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response
	_, res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
	if err != nil {
		return
	}
	if res != nil {
		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			resouceNrfUri = nrfUri
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
		}
	}
	return
}
