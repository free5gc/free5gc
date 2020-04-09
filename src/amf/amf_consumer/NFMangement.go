package amf_consumer

import (
	"context"
	"fmt"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_util"
	"free5gc/src/amf/logger"
	"net/http"
	"strings"
	"time"
)

func BuildNFInstance(context *amf_context.AMFContext) (profile models.NfProfile, err error) {
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_AMF
	profile.NfStatus = models.NfStatus_REGISTERED
	var plmns []models.PlmnId
	for _, plmnItem := range context.PlmnSupportList {
		plmns = append(plmns, plmnItem.PlmnId)
	}
	if len(plmns) > 0 {
		profile.PlmnList = &plmns
		// TODO: change to Per Plmn Support Snssai List
		profile.SNssais = &context.PlmnSupportList[0].SNssaiList
	}
	amfInfo := models.AmfInfo{}
	if len(context.ServedGuamiList) == 0 {
		err = fmt.Errorf("Gumai List is Empty in AMF")
		return
	}
	regionId, setId, _, err1 := amf_util.SeperateAmfId(context.ServedGuamiList[0].AmfId)
	if err1 != nil {
		err = err1
		return
	}
	amfInfo.AmfRegionId = regionId
	amfInfo.AmfSetId = setId
	amfInfo.GuamiList = &context.ServedGuamiList
	if len(context.SupportTaiLists) == 0 {
		err = fmt.Errorf("SupportTaiList is Empty in AMF")
		return
	}
	amfInfo.TaiList = &context.SupportTaiLists
	profile.AmfInfo = &amfInfo
	if context.HttpIPv4Address == "" {
		err = fmt.Errorf("AMF Address is empty")
		return
	}
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, context.HttpIPv4Address)
	service := []models.NfService{}
	for _, nfService := range context.NfService {
		service = append(service, nfService)
	}
	if len(service) > 0 {
		profile.NfServices = &service
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
			fmt.Println(fmt.Errorf("AMF register to NRF Error[%s]", err.Error()))
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
			fmt.Println(fmt.Errorf("handler returned wrong status code %d", status))
			fmt.Println(fmt.Errorf("NRF return wrong status code %d", status))
		}
	}
	return
}

func SendDeregisterNFInstance() (problemDetails *models.ProblemDetails, err error) {

	logger.ConsumerLog.Infof("[AMF] Send Deregister NFInstance")

	amfSelf := amf_context.AMF_Self()
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(amfSelf.NrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response

	res, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(context.Background(), amfSelf.NfId)
	if err == nil {
		return
	} else if res != nil {
		if res.Status != err.Error() {
			return
		}
		problem := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = common.ReportError("server no response")
	}
	return
}
