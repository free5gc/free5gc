package consumer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/nrf/NFDiscovery"
	"github.com/free5gc/openapi/nrf/NFManagement"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

type nnrfService struct {
	consumer *Consumer

	nfMngmntMu sync.RWMutex
	nfDiscMu   sync.RWMutex

	nfMngmntClients map[string]*NFManagement.APIClient
	nfDiscClients   map[string]*NFDiscovery.APIClient
}

func (s *nnrfService) getNFManagementClient(uri string) *NFManagement.APIClient {
	if uri == "" {
		return nil
	}
	s.nfMngmntMu.RLock()
	client, ok := s.nfMngmntClients[uri]
	if ok {
		defer s.nfMngmntMu.RUnlock()
		return client
	}

	configuration := NFManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client = NFManagement.NewAPIClient(configuration)

	s.nfMngmntMu.RUnlock()
	s.nfMngmntMu.Lock()
	defer s.nfMngmntMu.Unlock()
	s.nfMngmntClients[uri] = client
	return client
}

func (s *nnrfService) getNFDiscClient(uri string) *NFDiscovery.APIClient {
	if uri == "" {
		return nil
	}
	s.nfDiscMu.RLock()
	client, ok := s.nfDiscClients[uri]
	if ok {
		defer s.nfDiscMu.RUnlock()
		return client
	}

	configuration := NFDiscovery.NewConfiguration()
	configuration.SetBasePath(uri)
	client = NFDiscovery.NewAPIClient(configuration)

	s.nfDiscMu.RUnlock()
	s.nfDiscMu.Lock()
	defer s.nfDiscMu.Unlock()
	s.nfDiscClients[uri] = client
	return client
}

func (s *nnrfService) SendSearchNFInstances(
	nrfUri string, targetNfType, requestNfType models.NrfNfManagementNfType, param NFDiscovery.SearchNFInstancesRequest) (
	*models.SearchResult, error,
) {
	// Set client and set url
	client := s.getNFDiscClient(nrfUri)

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return nil, err
	}
	param.TargetNfType = &targetNfType
	param.RequesterNfType = &requestNfType
	res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, &param)
	if err != nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
		return nil, err
	}

	result := res.SearchResult

	return &result, nil
}

func (s *nnrfService) SendNFInstancesUDR(nrfUri, id string) string {
	targetNfType := models.NrfNfManagementNfType_UDR
	requestNfType := models.NrfNfManagementNfType_PCF
	localVarOptionals := NFDiscovery.SearchNFInstancesRequest{
		// 	DataSet: optional.NewInterface(models.DataSetId_SUBSCRIPTION),
	}

	result, err := s.SendSearchNFInstances(nrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.ConsumerLog.Error(err.Error())
		return ""
	}
	for _, profile := range result.NfInstances {
		if uri := util.SearchNFServiceUri(profile, models.ServiceName_NUDR_DR, models.NfServiceStatus_REGISTERED); uri != "" {
			return uri
		}
	}
	return ""
}

func (s *nnrfService) SendNFInstancesBSF(nrfUri string) string {
	targetNfType := models.NrfNfManagementNfType_BSF
	requestNfType := models.NrfNfManagementNfType_PCF
	localVarOptionals := NFDiscovery.SearchNFInstancesRequest{}

	result, err := s.SendSearchNFInstances(nrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.ConsumerLog.Error(err.Error())
		return ""
	}
	for _, profile := range result.NfInstances {
		if uri := util.SearchNFServiceUri(profile, models.ServiceName_NBSF_MANAGEMENT,
			models.NfServiceStatus_REGISTERED); uri != "" {
			return uri
		}
	}
	return ""
}

func (s *nnrfService) SendNFInstancesAMF(nrfUri string, guami models.Guami, serviceName models.ServiceName) string {
	targetNfType := models.NrfNfManagementNfType_AMF
	requestNfType := models.NrfNfManagementNfType_PCF

	localVarOptionals := NFDiscovery.SearchNFInstancesRequest{
		Guami: &guami,
	}

	result, err := s.SendSearchNFInstances(nrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.ConsumerLog.Error(err.Error())
		return ""
	}
	for _, profile := range result.NfInstances {
		return util.SearchNFServiceUri(profile, serviceName, models.NfServiceStatus_REGISTERED)
	}
	return ""
}

// management
func (s *nnrfService) BuildNFInstance(
	context *pcf_context.PCFContext,
) (profile models.NrfNfManagementNfProfile, err error) {
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NrfNfManagementNfType_PCF
	profile.NfStatus = models.NrfNfManagementNfStatus_REGISTERED
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, context.RegisterIPv4)
	services := []models.NrfNfManagementNfService{}
	for _, nfService := range context.NfService {
		services = append(services, nfService)
	}
	if len(services) > 0 {
		profile.NfServices = services
	}
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
	if context.Locality != "" {
		profile.Locality = context.Locality
	}
	return profile, nil
}

func (s *nnrfService) SendRegisterNFInstance(ctx context.Context) (
	resouceNrfUri string, retrieveNfInstanceID string, err error,
) {
	// Set client and set url
	pcfContext := s.consumer.Context()

	client := s.getNFManagementClient(pcfContext.NrfUri)
	nfProfile, err := s.BuildNFInstance(pcfContext)
	if err != nil {
		return "", "",
			errors.Wrap(err, "RegisterNFInstance buildNfProfile()")
	}

	var nf models.NrfNfManagementNfProfile
	var res *NFManagement.RegisterNFInstanceResponse

	finish := false
	for !finish {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("RegisterNFInstance context done")
		default:
			req := &NFManagement.RegisterNFInstanceRequest{
				NfInstanceID:             &pcfContext.NfId,
				NrfNfManagementNfProfile: &nfProfile,
			}
			res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(ctx, req)
			if err != nil || res == nil {
				logger.ConsumerLog.Errorf("PCF register to NRF Error[%v]", err)
				time.Sleep(2 * time.Second)
				continue
			}
			nf = res.NrfNfManagementNfProfile

			if res.Location == "" {
				// NFUpdate
				finish = true
			} else {
				// NFRegister
				resourceUri := res.Location
				resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
				retrieveNfInstanceID = resourceUri[strings.LastIndex(resourceUri, "/")+1:]

				oauth2 := false
				if nf.CustomInfo != nil {
					v, ok := nf.CustomInfo["oauth2"].(bool)
					if ok {
						oauth2 = v
						logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
					}
				}
				pcf_context.GetSelf().OAuth2Required = oauth2
				if oauth2 && pcf_context.GetSelf().NrfCertPem == "" {
					logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
				}

				finish = true
			}
		}
	}

	return resouceNrfUri, retrieveNfInstanceID, err
}

func (s *nnrfService) SendDeregisterNFInstance() (problemDetails *models.ProblemDetails, err error) {
	logger.ConsumerLog.Infof("Send Deregister NFInstance")

	ctx, pd, err := pcf_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return pd, err
	}

	pcfContext := s.consumer.Context()
	client := s.getNFManagementClient(pcfContext.NrfUri)
	request := &NFManagement.DeregisterNFInstanceRequest{
		NfInstanceID: &pcfContext.NfId,
	}

	_, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, request)

	return problemDetails, err
}
