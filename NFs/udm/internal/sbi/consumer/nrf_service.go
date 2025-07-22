package consumer

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/free5gc/openapi/models"
	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	Nnrf_NFManagement "github.com/free5gc/openapi/nrf/NFManagement"
	udm_context "github.com/free5gc/udm/internal/context"
	"github.com/free5gc/udm/internal/logger"
	"github.com/free5gc/udm/internal/util"
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
		defer s.nfDiscMu.RUnlock()
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

func (s *nnrfService) SendSearchNFInstances(
	nrfUri string, param Nnrf_NFDiscovery.SearchNFInstancesRequest) (
	*models.SearchResult, error,
) {
	// Set client and set url
	udmContext := s.consumer.Context()

	client := s.getNFDiscClient(udmContext.NrfUri)

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return nil, err
	}

	searchNfInstancesRsp, err1 := client.NFInstancesStoreApi.SearchNFInstances(ctx, &param)
	result := searchNfInstancesRsp.SearchResult
	if err1 != nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
	}

	return &result, nil
}

func (s *nnrfService) SendNFInstancesUDR(id string, types int) string {
	self := udm_context.GetSelf()
	targetNfType := models.NrfNfManagementNfType_UDR
	requestNfType := models.NrfNfManagementNfType_UDM
	searchNFinstanceRequest := Nnrf_NFDiscovery.SearchNFInstancesRequest{
		// 	DataSet: optional.NewInterface(models.DataSetId_SUBSCRIPTION),
	}
	searchNFinstanceRequest.RequesterNfType = &requestNfType
	searchNFinstanceRequest.TargetNfType = &targetNfType

	result, err := s.SendSearchNFInstances(self.NrfUri, searchNFinstanceRequest)
	if err != nil {
		logger.ConsumerLog.Error(err.Error())
		return ""
	}
	for _, profile := range result.NfInstances {
		return util.SearchNFServiceUri(profile, models.ServiceName_NUDR_DR, models.NfServiceStatus_REGISTERED)
	}
	return ""
}

func (s *nnrfService) SendDeregisterNFInstance() (err error) {
	logger.ConsumerLog.Infof("Send Deregister NFInstance")

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return err
	}

	udmContext := s.consumer.Context()
	client := s.getNFManagementClient(udmContext.NrfUri)

	var derigisterNfInstanceRequest Nnrf_NFManagement.DeregisterNFInstanceRequest
	derigisterNfInstanceRequest.NfInstanceID = &udmContext.NfId
	_, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, &derigisterNfInstanceRequest)

	return err
}

func (s *nnrfService) RegisterNFInstance(ctx context.Context) (
	resouceNrfUri string, retrieveNfInstanceID string, err error,
) {
	udmContext := s.consumer.Context()
	client := s.getNFManagementClient(udmContext.NrfUri)
	nfProfile, err := s.buildNfProfile(udmContext)
	if err != nil {
		return "", "", errors.Wrap(err, "RegisterNFInstance buildNfProfile()")
	}
	var registerNfInstanceRequest Nnrf_NFManagement.RegisterNFInstanceRequest
	registerNfInstanceRequest.NfInstanceID = &udmContext.NfId
	registerNfInstanceRequest.NrfNfManagementNfProfile = &nfProfile
	var res *Nnrf_NFManagement.RegisterNFInstanceResponse
	for {
		select {
		case <-ctx.Done():
			return "", "", errors.Errorf("Context Cancel before RegisterNFInstance")
		default:
		}

		res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(ctx, &registerNfInstanceRequest)

		if err != nil || res == nil {
			logger.ConsumerLog.Errorf("UDM register to NRF Error[%v]", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}

		if res.Location == "" {
			// NFUpdate
			break
		} else { // http.statusCreated
			// NFRegister
			resourceUri := res.Location
			resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceID = resourceUri[strings.LastIndex(resourceUri, "/")+1:]

			oauth2 := false
			if res.NrfNfManagementNfProfile.CustomInfo != nil {
				v, ok := res.NrfNfManagementNfProfile.CustomInfo["oauth2"].(bool)
				if ok {
					oauth2 = v
					logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
				}
			}
			udm_context.GetSelf().OAuth2Required = oauth2
			if oauth2 && udm_context.GetSelf().NrfCertPem == "" {
				logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
			}

			break
		}
	}
	return resouceNrfUri, retrieveNfInstanceID, err
}

func (s *nnrfService) buildNfProfile(udmContext *udm_context.UDMContext) (
	profile models.NrfNfManagementNfProfile, err error,
) {
	profile.NfInstanceId = udmContext.NfId
	profile.NfType = models.NrfNfManagementNfType_UDM
	profile.NfStatus = models.NrfNfManagementNfStatus_REGISTERED
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, udmContext.RegisterIPv4)
	for _, nfService := range udmContext.NfService {
		profile.NfServices = append(profile.NfServices, nfService)
	}
	profile.UdmInfo = &models.UdmInfo{
		// Todo
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
