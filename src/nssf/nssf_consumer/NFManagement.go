/*
 * NSSF Consumer
 *
 * Network Function Management
 */

package nssf_consumer

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/nssf/nssf_context"
)

func BuildNFProfile(context *nssf_context.NSSFContext) (profile models.NfProfile, err error) {
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_NSSF
	profile.NfStatus = models.NfStatus_REGISTERED
	profile.PlmnList = &context.SupportedPlmnList
	profile.Ipv4Addresses = []string{context.HttpIpv4Address}
	var services []models.NfService
	for _, nfService := range context.NfService {
		services = append(services, nfService)
	}
	if len(services) > 0 {
		profile.NfServices = &services
	}
	return
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (resourceNrfUri string, err error) {
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	apiClient := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response
	_, res, err = apiClient.NFInstanceIDDocumentApi.RegisterNFInstance(context.Background(), nfInstanceId, profile)
	if err != nil {
		return
	}
	if res != nil {
		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			resourceNrfUri = nrfUri
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resourceUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
		} else {
			err = fmt.Errorf("Unrecognized status code %d", status)
		}
	}
	return
}
