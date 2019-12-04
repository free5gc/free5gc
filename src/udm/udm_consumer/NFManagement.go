package udm_consumer

import (
	"context"
	"fmt"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/udm_context"
	"net/http"
	"strings"
	"time"
)

func BuildNFInstance(udmContext *udm_context.UDMContext) (profile models.NfProfile, err error) {
	profile.NfInstanceId = udmContext.NfId
	profile.NfStatus = models.NfStatus_REGISTERED
	profile.NfType = models.NfType_UDM
	services := []models.NfService{}
	for _, nfservice := range udmContext.NfService {
		services = append(services, nfservice)
	}
	if len(services) > 0 {
		profile.NfServices = &services
	}

	var udmInfo models.UdmInfo
	profile.UdmInfo = &udmInfo
	profile.UdmInfo.GroupId = udmContext.GroupId
	if udmContext.HttpIPv4Address == "" {
		err = fmt.Errorf("UDM Address is empty")
		return
	}
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, udmContext.HttpIPv4Address)

	return
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (resouceNrfUri string, retrieveNfInstanceId string, err error) {

	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response
	for {
		_, res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
		if err != nil || res == nil {
			//TODO : add log
			fmt.Println(fmt.Errorf("UDM register to NRF Error[%v]", err.Error()))
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
			resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			break
		} else {
			fmt.Println("handler returned wrong status code", status)
			fmt.Println("NRF return wrong status code", status)
		}
	}
	return
}
