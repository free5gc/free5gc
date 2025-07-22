package consumer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/util"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	Nnrf_NFManagement "github.com/free5gc/openapi/nrf/NFManagement"
)

type nnrfService struct {
	consumer *Consumer

	nfMngmntMu sync.RWMutex
	nfDiscMu   sync.RWMutex

	nfMngmntClients map[string]*Nnrf_NFManagement.APIClient
	nfDiscClients   map[string]*Nnrf_NFDiscovery.APIClient
}

func (s *nnrfService) getNFManagementClient(uri string) *Nnrf_NFManagement.APIClient {
	if uri == "" {
		return nil
	}
	s.nfMngmntMu.RLock()
	client, ok := s.nfMngmntClients[uri]
	if ok {
		s.nfMngmntMu.RUnlock()
		return client
	}

	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nnrf_NFManagement.NewAPIClient(configuration)

	s.nfMngmntMu.RUnlock()
	s.nfMngmntMu.Lock()
	defer s.nfMngmntMu.Unlock()
	s.nfMngmntClients[uri] = client
	return client
}

func (s *nnrfService) getNFDiscClient(uri string) *Nnrf_NFDiscovery.APIClient {
	if uri == "" {
		return nil
	}
	s.nfDiscMu.RLock()
	client, ok := s.nfDiscClients[uri]
	if ok {
		s.nfDiscMu.RUnlock()
		return client
	}

	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nnrf_NFDiscovery.NewAPIClient(configuration)

	s.nfDiscMu.RUnlock()
	s.nfDiscMu.Lock()
	defer s.nfDiscMu.Unlock()
	s.nfDiscClients[uri] = client
	return client
}

func (s *nnrfService) SendSearchNFInstances(nrfUri string, targetNfType, requestNfType models.NrfNfManagementNfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesRequest,
) (*models.SearchResult, error) {
	// Set client and set url
	param.TargetNfType = &targetNfType
	param.RequesterNfType = &requestNfType
	client := s.getNFDiscClient(nrfUri)
	if client == nil {
		return nil, openapi.ReportError("nrf not found")
	}

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return nil, err
	}
	res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, param)
	var result *models.SearchResult
	if err != nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
	}
	if res != nil {
		result = &res.SearchResult
	}
	return result, err
}

func (s *nnrfService) SearchUdmSdmInstance(
	ue *amf_context.AmfUe, nrfUri string, targetNfType, requestNfType models.NrfNfManagementNfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesRequest,
) error {
	resp, localErr := s.SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)
	if localErr != nil {
		return localErr
	}

	// select the first UDM_SDM, TODO: select base on other info
	var sdmUri string
	for index := range resp.NfInstances {
		ue.UdmId = resp.NfInstances[index].NfInstanceId
		sdmUri = util.SearchNFServiceUri(&resp.NfInstances[index], models.ServiceName_NUDM_SDM,
			models.NfServiceStatus_REGISTERED)
		if sdmUri != "" {
			break
		}
	}
	ue.NudmSDMUri = sdmUri
	if ue.NudmSDMUri == "" {
		err := fmt.Errorf("AMF can not select an UDM by NRF")
		logger.ConsumerLog.Error(err)
		return err
	}
	return nil
}

func (s *nnrfService) SearchNssfNSSelectionInstance(
	ue *amf_context.AmfUe, nrfUri string, targetNfType, requestNfType models.NrfNfManagementNfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesRequest,
) error {
	resp, localErr := s.SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)
	if localErr != nil {
		return localErr
	}

	// select the first NSSF, TODO: select base on other info
	var nssfUri string
	for index := range resp.NfInstances {
		ue.NssfId = resp.NfInstances[index].NfInstanceId
		nssfUri = util.SearchNFServiceUri(&resp.NfInstances[index], models.ServiceName_NNSSF_NSSELECTION,
			models.NfServiceStatus_REGISTERED)
		if nssfUri != "" {
			break
		}
	}
	ue.NssfUri = nssfUri
	if ue.NssfUri == "" {
		return fmt.Errorf("AMF can not select an NSSF by NRF")
	}
	return nil
}

func (s *nnrfService) SearchAmfCommunicationInstance(ue *amf_context.AmfUe, nrfUri string, targetNfType,
	requestNfType models.NrfNfManagementNfType, param *Nnrf_NFDiscovery.SearchNFInstancesRequest,
) (err error) {
	resp, localErr := s.SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)
	if localErr != nil {
		err = localErr
		return
	}

	// select the first AMF, TODO: select base on other info
	var amfUri string
	for index := range resp.NfInstances {
		if resp.NfInstances[index].NfInstanceId == amf_context.GetSelf().NfId {
			continue
		}
		ue.TargetAmfProfile = &resp.NfInstances[index]
		amfUri = util.SearchNFServiceUri(&resp.NfInstances[index], models.ServiceName_NAMF_COMM,
			models.NfServiceStatus_REGISTERED)
		if amfUri != "" {
			break
		}
	}
	ue.TargetAmfUri = amfUri
	if ue.TargetAmfUri == "" {
		err = fmt.Errorf("AMF can not select an target AMF by NRF")
	}
	return
}

func (s *nnrfService) BuildNFInstance(context *amf_context.AMFContext) (
	profile models.NrfNfManagementNfProfile, err error,
) {
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NrfNfManagementNfType_AMF
	profile.NfStatus = models.NrfNfManagementNfStatus_REGISTERED
	var plmns []models.PlmnId
	for _, plmnItem := range context.PlmnSupportList {
		plmns = append(plmns, *plmnItem.PlmnId)
	}
	if len(plmns) > 0 {
		profile.PlmnList = plmns
		// TODO: change to Per Plmn Support Snssai List
		var SnssaiList []models.ExtSnssai
		for _, snssaiItem := range context.PlmnSupportList[0].SNssaiList {
			SnssaiList = append(SnssaiList, util.SnssaiModelsToExtSnssai(snssaiItem))
		}
		profile.SNssais = SnssaiList
	}
	amfInfo := models.NrfNfManagementAmfInfo{}
	if len(context.ServedGuamiList) == 0 {
		err = fmt.Errorf("gumai List is Empty in AMF")
		return profile, err
	}
	regionId, setId, _, err1 := util.SeperateAmfId(context.ServedGuamiList[0].AmfId)
	if err1 != nil {
		err = err1
		return profile, err
	}
	amfInfo.AmfRegionId = regionId
	amfInfo.AmfSetId = setId
	amfInfo.GuamiList = context.ServedGuamiList
	if len(context.SupportTaiLists) == 0 {
		err = fmt.Errorf("SupportTaiList is Empty in AMF")
		return profile, err
	}
	amfInfo.TaiList = context.SupportTaiLists
	profile.AmfInfo = &amfInfo
	if context.RegisterIPv4 == "" {
		err = fmt.Errorf("AMF Address is empty")
		return profile, err
	}
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, context.RegisterIPv4)
	service := []models.NrfNfManagementNfService{}
	for _, nfService := range context.NfService {
		service = append(service, nfService)
	}
	if len(service) > 0 {
		profile.NfServices = service
	}

	defaultNotificationSubscription := models.DefaultNotificationSubscription{
		CallbackUri:      fmt.Sprintf("%s"+factory.AmfCallbackResUriPrefix+"/n1-message-notify", context.GetIPv4Uri()),
		NotificationType: models.NrfNfManagementNotificationType_N1_MESSAGES,
		N1MessageClass:   models.N1MessageClass__5_GMM,
	}
	profile.DefaultNotificationSubscriptions = append(profile.DefaultNotificationSubscriptions,
		defaultNotificationSubscription)
	return profile, err
}

func (s *nnrfService) SendRegisterNFInstance(ctx context.Context, nrfUri, nfInstanceId string,
	profile *models.NrfNfManagementNfProfile) (
	resouceNrfUri string, retrieveNfInstanceId string, err error,
) {
	// Set client and set url
	client := s.getNFManagementClient(nrfUri)
	if client == nil {
		return "", "", openapi.ReportError("nrf not found")
	}

	var res *Nnrf_NFManagement.RegisterNFInstanceResponse
	var nf models.NrfNfManagementNfProfile
	registerNFInstanceRequest := &Nnrf_NFManagement.RegisterNFInstanceRequest{
		NfInstanceID:             &nfInstanceId,
		NrfNfManagementNfProfile: profile,
	}
	finish := false
	for !finish {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("context done")
		default:
			res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(ctx, registerNFInstanceRequest)
			if err != nil || res == nil {
				// TODO : add log
				logger.ConsumerLog.Errorf("AMF register to NRF Error[%s]", err.Error())
				time.Sleep(2 * time.Second)
				continue
			}
			if res.Location == "" {
				// NFUpdate
				finish = true
			} else {
				// NFRegister
				resourceUri := res.Location
				nf = res.NrfNfManagementNfProfile
				index := strings.Index(resourceUri, "/nnrf-nfm/")
				if index >= 0 {
					resouceNrfUri = resourceUri[:index]
				}
				// resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
				retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]

				oauth2 := false
				if nf.CustomInfo != nil {
					v, ok := nf.CustomInfo["oauth2"].(bool)
					if ok {
						oauth2 = v
						logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
					}
				}
				amf_context.GetSelf().OAuth2Required = oauth2
				if oauth2 && amf_context.GetSelf().NrfCertPem == "" {
					logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
				}
				finish = true
			}
		}
	}
	return resouceNrfUri, retrieveNfInstanceId, err
}

func (s *nnrfService) SendDeregisterNFInstance() (problemDetails *models.ProblemDetails, err error) {
	logger.ConsumerLog.Infof("[AMF] Send Deregister NFInstance")
	amfContext := s.consumer.Context()

	client := s.getNFManagementClient(amfContext.NrfUri)
	if client == nil {
		return nil, openapi.ReportError("nrf not found")
	}

	ctx, pd, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return pd, err
	}

	request := &Nnrf_NFManagement.DeregisterNFInstanceRequest{
		NfInstanceID: &amfContext.NfId,
	}

	_, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, request)
	if err != nil {
		switch apiErr := err.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errModel := apiErr.Model().(type) {
			case Nnrf_NFManagement.DeregisterNFInstanceError:
				problemDetails = &errModel.ProblemDetails
			case error:
				problemDetails = openapi.ProblemDetailsSystemFailure(errModel.Error())
			default:
				err = openapi.ReportError("openapi error")
			}
		case error:
			problemDetails = openapi.ProblemDetailsSystemFailure(apiErr.Error())
		default:
			err = openapi.ReportError("server no response")
		}
	}

	return problemDetails, err
}
