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
	"time"
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

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (resourceNrfUri string, retrieveNfInstanceId string, err error) {
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	apiClient := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response
	for {
		_, res, err = apiClient.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
		if err != nil || res == nil {
			//TODO : add log
			fmt.Println(fmt.Errorf("NSSF register to NRF Error[%s]", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}
		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			break
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resourceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			break
		} else {
			fmt.Println("NRF return wrong status code", status)
		}
	}
	return
}
