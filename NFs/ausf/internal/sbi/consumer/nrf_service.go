package consumer

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	ausf_context "github.com/free5gc/ausf/internal/context"
	"github.com/free5gc/ausf/internal/logger"
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

func (s *nnrfService) SendSearchNFInstances(
	nrfUri string, targetNfType,
	requestNfType models.NrfNfManagementNfType,
	param Nnrf_NFDiscovery.SearchNFInstancesRequest,
) (
	*models.SearchResult, error,
) {
	// Set client and set url
	client := s.getNFDiscClient(nrfUri)
	if client == nil {
		return nil, openapi.ReportError("nrf not found")
	}

	ctx, _, err := ausf_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return nil, err
	}

	res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, &param)

	if err != nil || res == nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
		return nil, err
	}

	result := res.SearchResult
	return &result, err
}

func (s *nnrfService) SendDeregisterNFInstance() (*models.ProblemDetails, error) {
	logger.ConsumerLog.Infof("[AUSF] Send Deregister NFInstance")

	ctx, pd, err := ausf_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return pd, err
	}

	ausfContext := s.consumer.Context()
	client := s.getNFManagementClient(ausfContext.NrfUri)
	request := &Nnrf_NFManagement.DeregisterNFInstanceRequest{
		NfInstanceID: &ausfContext.NfId,
	}

	_, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, request)
	if apiErr, ok := err.(openapi.GenericOpenAPIError); ok {
		// API error
		if deregNfError, okDeg := apiErr.Model().(Nnrf_NFManagement.DeregisterNFInstanceError); okDeg {
			return &deregNfError.ProblemDetails, err
		}
		return nil, err
	}
	return nil, err
}

func (s *nnrfService) RegisterNFInstance(ctx context.Context) (
	resouceNrfUri string, retrieveNfInstanceID string, err error,
) {
	ausfContext := s.consumer.Context()
	client := s.getNFManagementClient(ausfContext.NrfUri)
	nfProfile, err := s.buildNfProfile(ausfContext)
	if err != nil {
		return "", "", errors.Wrap(err, "RegisterNFInstance buildNfProfile()")
	}

	var nf models.NrfNfManagementNfProfile
	var res *Nnrf_NFManagement.RegisterNFInstanceResponse
	registerNFInstanceRequest := &Nnrf_NFManagement.RegisterNFInstanceRequest{
		NfInstanceID:             &ausfContext.NfId,
		NrfNfManagementNfProfile: &nfProfile,
	}
	for {
		select {
		case <-ctx.Done():
			return "", "", errors.Errorf("Context Cancel before RegisterNFInstance")
		default:
		}
		res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(ctx, registerNFInstanceRequest)
		if err != nil || res == nil {
			logger.ConsumerLog.Errorf("AUSF register to NRF Error[%v]", err)
			time.Sleep(2 * time.Second)
			continue
		}
		nf = res.NrfNfManagementNfProfile

		// http.StatusOK
		if res.Location == "" {
			// NFUpdate
			break
		} else { // http.StatusCreated
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
			ausf_context.GetSelf().OAuth2Required = oauth2
			if oauth2 && ausf_context.GetSelf().NrfCertPem == "" {
				logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
			}

			break
		}
	}
	return resouceNrfUri, retrieveNfInstanceID, err
}

func (s *nnrfService) buildNfProfile(ausfContext *ausf_context.AUSFContext) (
	profile models.NrfNfManagementNfProfile, err error,
) {
	profile.NfInstanceId = ausfContext.NfId
	profile.NfType = models.NrfNfManagementNfType_AUSF
	profile.NfStatus = models.NrfNfManagementNfStatus_REGISTERED
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, ausfContext.RegisterIPv4)
	services := []models.NrfNfManagementNfService{}
	for _, nfService := range ausfContext.NfService {
		services = append(services, nfService)
	}
	if len(services) > 0 {
		profile.NfServices = services
	}
	profile.AusfInfo = &models.AusfInfo{
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

func (s *nnrfService) GetUdmUrl(nrfUri string) string {
	udmUrl := "https://localhost:29503" // default
	targetNfType := models.NrfNfManagementNfType_UDM
	requestNfType := models.NrfNfManagementNfType_AUSF
	nfDiscoverParam := Nnrf_NFDiscovery.SearchNFInstancesRequest{
		RequesterNfType: &requestNfType,
		TargetNfType:    &targetNfType,
		ServiceNames:    []models.ServiceName{models.ServiceName_NUDM_UEAU},
	}
	res, err := s.SendSearchNFInstances(
		nrfUri,
		models.NrfNfManagementNfType_UDM,
		models.NrfNfManagementNfType_AUSF,
		nfDiscoverParam,
	)
	if err != nil {
		logger.ConsumerLog.Errorln("[Search UDM UEAU] ", err.Error(), "use defalt udmUrl", udmUrl)
	} else if len(res.NfInstances) > 0 {
		udmInstance := res.NfInstances[0]
		if len(udmInstance.Ipv4Addresses) > 0 && udmInstance.NfServices != nil {
			ueauService := udmInstance.NfServices[0]
			ueauEndPoint := ueauService.IpEndPoints[0]
			udmUrl = string(ueauService.Scheme) + "://" + ueauEndPoint.Ipv4Address + ":" + strconv.Itoa(int(ueauEndPoint.Port))
		}
	} else {
		logger.ConsumerLog.Errorln("[Search UDM UEAU] len(NfInstances) = 0")
	}
	return udmUrl
}
