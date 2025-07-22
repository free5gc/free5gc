package consumer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/nrf/NFDiscovery"
	"github.com/free5gc/openapi/nrf/NFManagement"
	"github.com/free5gc/openapi/udm/SubscriberDataManagement"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/internal/logger"
)

type nnrfService struct {
	consumer *Consumer

	NFManagementgMu sync.RWMutex
	NFDiscoveryMu   sync.RWMutex

	NFManagementClients map[string]*NFManagement.APIClient
	NFDiscoveryClients  map[string]*NFDiscovery.APIClient
}

func (s *nnrfService) getNFManagementClient(uri string) *NFManagement.APIClient {
	if uri == "" {
		return nil
	}
	s.NFManagementgMu.RLock()
	client, ok := s.NFManagementClients[uri]
	if ok {
		s.NFManagementgMu.RUnlock()
		return client
	}

	configuration := NFManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client = NFManagement.NewAPIClient(configuration)

	s.NFManagementgMu.RUnlock()
	s.NFManagementgMu.Lock()
	defer s.NFManagementgMu.Unlock()
	s.NFManagementClients[uri] = client
	return client
}

func (s *nnrfService) getNFDiscoveryClient(uri string) *NFDiscovery.APIClient {
	if uri == "" {
		return nil
	}
	s.NFDiscoveryMu.RLock()
	client, ok := s.NFDiscoveryClients[uri]
	if ok {
		s.NFDiscoveryMu.RUnlock()
		return client
	}

	configuration := NFDiscovery.NewConfiguration()
	configuration.SetBasePath(uri)
	client = NFDiscovery.NewAPIClient(configuration)

	s.NFDiscoveryMu.RUnlock()
	s.NFDiscoveryMu.Lock()
	defer s.NFDiscoveryMu.Unlock()
	s.NFDiscoveryClients[uri] = client
	return client
}

func (s *nnrfService) RegisterNFInstance(ctx context.Context) error {
	smfContext := s.consumer.Context()
	client := s.getNFManagementClient(smfContext.NrfUri)
	nfProfile, err := s.buildNfProfile(smfContext)
	if err != nil {
		return errors.Wrap(err, "RegisterNFInstance buildNfProfile()")
	}

	var nf models.NrfNfManagementNfProfile
	var res *NFManagement.RegisterNFInstanceResponse
	registerNFInstanceRequest := &NFManagement.RegisterNFInstanceRequest{
		NfInstanceID:             &smfContext.NfInstanceID,
		NrfNfManagementNfProfile: &nfProfile,
	}

	// Check data (Use RESTful PUT)
	finish := false
	for !finish {
		select {
		case <-ctx.Done():
			return fmt.Errorf("RegisterNFInstance context done")
		default:
			res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(ctx, registerNFInstanceRequest)
			if err != nil || res == nil {
				logger.ConsumerLog.Errorf("SMF register to NRF Error[%s]", err.Error())
				time.Sleep(2 * time.Second)
				continue
			}
			nf = res.NrfNfManagementNfProfile

			// http.StatusOK
			if res.Location == "" {
				// NFUpdate
				finish = true
			} else { // http.StatusCreated
				// NFRegister
				resourceUri := res.Location
				smfContext.NfInstanceID = resourceUri[strings.LastIndex(resourceUri, "/")+1:]

				oauth2 := false
				if nf.CustomInfo != nil {
					v, ok := nf.CustomInfo["oauth2"].(bool)
					if ok {
						oauth2 = v
						logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
					}
				}
				smfContext.OAuth2Required = oauth2
				if oauth2 && smfContext.NrfCertPem == "" {
					logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
				}
				finish = true
			}
		}
	}

	logger.InitLog.Infof("SMF Registration to NRF %v", nf)
	return nil
}

func (s *nnrfService) buildNfProfile(smfContext *smf_context.SMFContext) (
	profile models.NrfNfManagementNfProfile, err error,
) {
	smfProfile := smfContext.NfProfile

	sNssais := []models.ExtSnssai{}
	for _, snssaiSmfInfo := range smfProfile.SMFInfo.SNssaiSmfInfoList {
		sNssais = append(sNssais, *snssaiSmfInfo.SNssai)
	}

	// set nfProfile
	profile = models.NrfNfManagementNfProfile{
		NfInstanceId:  smfContext.NfInstanceID,
		NfType:        models.NrfNfManagementNfType_SMF,
		NfStatus:      models.NrfNfManagementNfStatus_REGISTERED,
		Ipv4Addresses: []string{smfContext.RegisterIPv4},
		NfServices:    *smfProfile.NFServices,
		SmfInfo:       smfProfile.SMFInfo,
		SNssais:       sNssais,
		PlmnList:      *smfProfile.PLMNList,
	}
	if smfContext.Locality != "" {
		profile.Locality = smfContext.Locality
	}
	return profile, err
}

func (s *nnrfService) SendDeregisterNFInstance() (err error) {
	logger.ConsumerLog.Infof("Send Deregister NFInstance")

	smfContext := s.consumer.Context()
	ctx, pd, err := smfContext.GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
	if err != nil {
		logger.ConsumerLog.Errorf("Get token context failed, problem details: %+v", pd)
		return err
	}

	client := s.getNFManagementClient(smfContext.NrfUri)
	request := &NFManagement.DeregisterNFInstanceRequest{
		NfInstanceID: &smfContext.NfInstanceID,
	}

	_, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, request)

	return err
}

func (s *nnrfService) SendSearchNFInstances(
	nrfUri string,
	targetNfType, requestNfType models.NrfNfManagementNfType,
	param *NFDiscovery.SearchNFInstancesRequest,
) (*models.SearchResult, error) {
	// Set client and set url
	smfContext := s.consumer.Context()
	client := s.getNFDiscoveryClient(smfContext.NrfUri)

	if client == nil {
		return nil, openapi.ReportError("nrf not found")
	}

	ctx, _, err := smfContext.GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return nil, err
	}

	res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, param)
	if err != nil || res == nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
		return nil, err
	}
	result := res.SearchResult
	return &result, err
}

func (s *nnrfService) NFDiscoveryUDM(ctx context.Context) (
	result models.SearchResult, localErr error,
) {
	targetNfType := models.NrfNfManagementNfType_UDM
	requesterNfType := models.NrfNfManagementNfType_SMF
	request := &NFDiscovery.SearchNFInstancesRequest{
		TargetNfType:    &targetNfType,
		RequesterNfType: &requesterNfType,
	}

	smfContext := s.consumer.Context()

	client := s.getNFDiscoveryClient(smfContext.NrfUri)
	// Check data
	res, localErr := client.NFInstancesStoreApi.SearchNFInstances(ctx, request)
	if res != nil {
		result = res.SearchResult
	}
	return result, localErr
}

func (s *nnrfService) NFDiscoveryPCF(ctx context.Context) (
	result models.SearchResult, localErr error,
) {
	targetNfType := models.NrfNfManagementNfType_PCF
	requesterNfType := models.NrfNfManagementNfType_SMF
	request := &NFDiscovery.SearchNFInstancesRequest{
		TargetNfType:    &targetNfType,
		RequesterNfType: &requesterNfType,
	}

	smfContext := s.consumer.Context()

	client := s.getNFDiscoveryClient(smfContext.NrfUri)
	// Check data
	res, localErr := client.NFInstancesStoreApi.SearchNFInstances(ctx, request)
	if res != nil {
		result = res.SearchResult
	}
	return result, localErr
}

func (s *nnrfService) NFDiscoveryAMF(smContext *smf_context.SMContext, ctx context.Context) (
	result models.SearchResult, localErr error,
) {
	targetNfType := models.NrfNfManagementNfType_AMF
	requesterNfType := models.NrfNfManagementNfType_SMF
	request := &NFDiscovery.SearchNFInstancesRequest{
		TargetNfType:       &targetNfType,
		RequesterNfType:    &requesterNfType,
		TargetNfInstanceId: &smContext.ServingNfId,
	}

	smfContext := s.consumer.Context()

	client := s.getNFDiscoveryClient(smfContext.NrfUri)
	// Check data
	res, localErr := client.NFInstancesStoreApi.SearchNFInstances(ctx, request)
	if res != nil {
		result = res.SearchResult
	}
	return result, localErr
}

func (s *nnrfService) SendNFDiscoveryUDM() (*models.ProblemDetails, error) {
	smfContext := s.consumer.Context()
	ctx, pd, err := smfContext.GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return pd, err
	}

	// Check data
	result, localErr := s.NFDiscoveryUDM(ctx)

	switch err := localErr.(type) {
	case openapi.GenericOpenAPIError:
		switch errModel := err.Model().(type) {
		case NFDiscovery.SearchNFInstancesError:
			return &errModel.ProblemDetails, nil
		case error:
			return openapi.ProblemDetailsSystemFailure(errModel.Error()), nil
		default:
			return nil, openapi.ReportError("openapi error")
		}
	case error:
		return openapi.ProblemDetailsSystemFailure(err.Error()), nil
	case nil:
		smfContext.UDMProfile = result.NfInstances[0]

		var client *SubscriberDataManagement.APIClient
		for _, service := range smfContext.UDMProfile.NfServices {
			if service.ServiceName == models.ServiceName_NUDM_SDM {
				client = s.consumer.nudmService.getSubscribeDataManagementClient(service.ApiPrefix)
			}
		}
		if client == nil {
			logger.ConsumerLog.Traceln("Get Subscribe Data Management Client Failed")
			return nil, fmt.Errorf("get Subscribe Data Management Client Failed")
		}
	default:
		return nil, openapi.ReportError("server no response")
	}

	return nil, nil
}

func (s *nnrfService) SendNFDiscoveryPCF() (*models.ProblemDetails, error) {
	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return pd, err
	}

	// Check data
	result, localErr := s.NFDiscoveryPCF(ctx)

	switch err := localErr.(type) {
	case openapi.GenericOpenAPIError:
		switch errModel := err.Model().(type) {
		case NFDiscovery.SearchNFInstancesError:
			return &errModel.ProblemDetails, nil
		case error:
			return openapi.ProblemDetailsSystemFailure(errModel.Error()), nil
		default:
			return nil, openapi.ReportError("openapi error")
		}
	case error:
		return openapi.ProblemDetailsSystemFailure(err.Error()), nil
	case nil:
		logger.ConsumerLog.Traceln(result.NfInstances)
	default:
		return nil, openapi.ReportError("server no response")
	}

	return nil, nil
}

func (s *nnrfService) SendNFDiscoveryServingAMF(smContext *smf_context.SMContext) (*models.ProblemDetails, error) {
	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return pd, err
	}

	// Check data
	result, localErr := s.NFDiscoveryAMF(smContext, ctx)

	switch err := localErr.(type) {
	case openapi.GenericOpenAPIError:
		switch errModel := err.Model().(type) {
		case NFDiscovery.SearchNFInstancesError:
			return &errModel.ProblemDetails, nil
		case error:
			return openapi.ProblemDetailsSystemFailure(errModel.Error()), nil
		default:
			return nil, openapi.ReportError("openapi error")
		}
	case error:
		return openapi.ProblemDetailsSystemFailure(err.Error()), nil
	case nil:
		if result.NfInstances == nil {
			logger.ConsumerLog.Warnln("NfInstances is nil")
			return nil, openapi.ReportError("NfInstances is nil")
		}
		logger.ConsumerLog.Info("SendNFDiscoveryServingAMF ok")
		smContext.AMFProfile = result.NfInstances[0]
	default:
		return nil, openapi.ReportError("server no response")
	}

	return nil, nil
}

// CHFSelection will select CHF for this SM Context
func (s *nnrfService) CHFSelection(smContext *smf_context.SMContext) error {
	// Send NFDiscovery for find CHF
	targetNfType := models.NrfNfManagementNfType_CHF
	requesterNfType := models.NrfNfManagementNfType_SMF
	request := &NFDiscovery.SearchNFInstancesRequest{
		TargetNfType:    &targetNfType,
		RequesterNfType: &requesterNfType,
		// Supi:            &smContext.Supi,
	}

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return err
	}

	client := s.getNFDiscoveryClient(s.consumer.Context().NrfUri)
	// Check data
	res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, request)
	if err != nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
		return err
	}

	// Select CHF from available CHF
	if res != nil && len(res.SearchResult.NfInstances) > 0 {
		smContext.SelectedCHFProfile = res.SearchResult.NfInstances[0]
		return nil
	}
	return fmt.Errorf("no CHF found in CHFSelection")
}

// PCFSelection will select PCF for this SM Context
func (s *nnrfService) PCFSelection(smContext *smf_context.SMContext) error {
	ctx, _, errToken := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_DISC, "NRF")
	if errToken != nil {
		return errToken
	}
	// Send NFDiscovery for find PCF
	targetNfType := models.NrfNfManagementNfType_PCF
	requesterNfType := models.NrfNfManagementNfType_SMF
	request := &NFDiscovery.SearchNFInstancesRequest{
		TargetNfType:    &targetNfType,
		RequesterNfType: &requesterNfType,
	}

	if s.consumer.Context().Locality != "" {
		request.PreferredLocality = &s.consumer.Context().Locality
	}

	client := s.getNFDiscoveryClient(s.consumer.Context().NrfUri)
	// Check data
	res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, request)
	if err != nil {
		return err
	}

	// Select PCF from available PCF

	smContext.SelectedPCFProfile = res.SearchResult.NfInstances[0]

	return nil
}

func (s *nnrfService) SearchNFInstances(
	ctx context.Context,
	targetNfType, requesterNfType models.NrfNfManagementNfType,
	localVarOptionals *NFDiscovery.SearchNFInstancesRequest,
) (*models.SearchResult, error) {
	client := s.getNFDiscoveryClient(s.consumer.Context().NrfUri)

	request := &NFDiscovery.SearchNFInstancesRequest{
		TargetNfType:    &targetNfType,
		RequesterNfType: &requesterNfType,
	}

	res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, request)
	if err != nil {
		return nil, err
	}

	return &res.SearchResult, err
}
