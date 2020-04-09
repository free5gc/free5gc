package udr_consumer

import (
	"context"
	"fmt"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/udr/factory"
	"free5gc/src/udr/udr_context"
	"net/http"
	"strings"
	"time"
)

func BuildNFInstance(context *udr_context.UDRContext) (profile models.NfProfile) {
	config := factory.UdrConfig
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_UDR
	profile.NfStatus = models.NfStatus_REGISTERED
	version := config.Info.Version
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	apiPrefix := fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, context.HttpIpv4Port)
	services := []models.NfService{
		{
			ServiceInstanceId: "DataRepository",
			ServiceName:       models.ServiceName_NUDR_DR,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          models.UriScheme(context.UriScheme),
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       apiPrefix,
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: context.HttpIPv4Address,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(context.HttpIpv4Port),
				},
			},
		},
	}
	profile.NfServices = &services
	// TODO: finish the Udr Info
	profile.UdrInfo = &models.UdrInfo{
		SupportedDataSets: []models.DataSetId{
			// models.DataSetId_APPLICATION,
			// models.DataSetId_EXPOSURE,
			// models.DataSetId_POLICY,
			models.DataSetId_SUBSCRIPTION,
		},
	}
	return
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (resouceNrfUri string, retrieveNfInstanceId string, err error) {

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response
	for {
		_, res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
		if err != nil || res == nil {
			//TODO : add log
			fmt.Println(fmt.Errorf("UDR register to NRF Error[%s]", err.Error()))
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
