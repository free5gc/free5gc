package smf_consumer

import (
	"context"
	"fmt"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/Nudm_SubscriberDataManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/smf_context"
	"net/http"

	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/mohae/deepcopy"
)

func SendNFRegistration() error {

	//set nfProfile
	profile := models.NfProfile{
		NfInstanceId:  smf_context.SMF_Self().NfInstanceID,
		NfType:        models.NfType_SMF,
		NfStatus:      models.NfStatus_REGISTERED,
		Ipv4Addresses: []string{smf_context.SMF_Self().HTTPAddress},
		NfServices:    smf_context.NFServices,
		SmfInfo:       smf_context.SmfInfo,
	}
	var rep models.NfProfile
	var res *http.Response
	var err error

	// Check data (Use RESTful PUT)
	for {
		rep, res, err = smf_context.SMF_Self().NFManagementClient.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), smf_context.SMF_Self().NfInstanceID, profile)
		if err != nil || res == nil {
			logger.AppLog.Infof("SMF register to NRF Error[%s]", err.Error())
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
			// resouceNrfUri := resourceUri[strings.LastIndex(resourceUri, "/"):]
			smf_context.SMF_Self().NfInstanceID = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			break
		} else {
			logger.AppLog.Infof("handler returned wrong status code %d", status)
			// fmt.Errorf("NRF return wrong status code %d", status)
		}
	}

	logger.InitLog.Infof("SMF Registration to NRF %v", rep)
	return nil
}

func RetrySendNFRegistration(MaxRetry int) error {

	retryCount := 0
	for retryCount < MaxRetry {
		err := SendNFRegistration()
		if err == nil {
			return nil
		}
		logger.AppLog.Warnf("Send NFRegistration Failed by %v", err)
		retryCount++
	}

	return fmt.Errorf("[SMF] Retry NF Registration has meet maximum")
}

func SendNFDeregistration() {

	// Check data (Use RESTful DELETE)
	res, err := smf_context.SMF_Self().NFManagementClient.NFInstanceIDDocumentApi.DeregisterNFInstance(context.TODO(), smf_context.SMF_Self().NfInstanceID)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusNoContent {
			logger.AppLog.Info("handler returned wrong status code", status)
		}
	}

}

func SendNFDiscoveryUDM() {
	if smf_context.SMF_Self().SubscriberDataManagementClient != nil {
		return
	}

	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}

	// Check data
	rep, res, err := smf_context.SMF_Self().NFDiscoveryClient.NFInstancesStoreApi.SearchNFInstances(context.TODO(), models.NfType_UDM, models.NfType_SMF, &localVarOptionals)
	if err != nil {
		return
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			logger.AppLog.Info("handler returned wrong status code", status)
		}
	}

	smf_context.SMF_Self().UDMProfile = rep.NfInstances[0]

	for _, service := range *smf_context.SMF_Self().UDMProfile.NfServices {
		if service.ServiceName == models.ServiceName_NUDM_SDM {
			SDMConf := Nudm_SubscriberDataManagement.NewConfiguration()
			SDMConf.SetBasePath(service.ApiPrefix)
			smf_context.SMF_Self().SubscriberDataManagementClient = Nudm_SubscriberDataManagement.NewAPIClient(SDMConf)
		}
	}

	if smf_context.SMF_Self().SubscriberDataManagementClient == nil {
		logger.AppLog.Warnln("sdm client failed")
	}
}

func SendNFDiscoveryPCF() {

	// Set targetNfType
	targetNfType := models.NfType_PCF
	// Set requestNfType
	requesterNfType := models.NfType_SMF
	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}

	// Check data
	rep, res, err := smf_context.SMF_Self().NFDiscoveryClient.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
	if err != nil {
		return
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			logger.AppLog.Info("handler returned wrong status code", status)
		}
	}

	logger.AppLog.Traceln(rep.NfInstances)
}

func SendNFDiscoveryServingAMF(smContext *smf_context.SMContext) {
	targetNfType := models.NfType_AMF
	requesterNfType := models.NfType_SMF

	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}

	localVarOptionals.TargetNfInstanceId = optional.NewInterface(smContext.ServingNfId)

	// Check data
	rep, res, err := smf_context.SMF_Self().NFDiscoveryClient.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requesterNfType, &localVarOptionals)
	if err != nil {
		logger.AppLog.Panic(err)
		return
	}
	if rep.NfInstances == nil {
		if status := res.StatusCode; status != http.StatusOK {
			logger.AppLog.Info("handler returned wrong status code", status)
		}
		logger.AppLog.Info("rep.NfInstances == nil")
		return
	}
	logger.AppLog.Info("SendNFDiscoveryServingAMF ok")
	smContext.AMFProfile = deepcopy.Copy(rep.NfInstances[0]).(models.NfProfile)

}
