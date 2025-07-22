package consumer

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antihax/optional"
	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/models"
)

const (
	RetryRegisterNrfDuration = 2 * time.Second
)

var serviceNfType map[models.ServiceName]models.NfType

func init() {
	serviceNfType = make(map[models.ServiceName]models.NfType)
	serviceNfType[models.ServiceName_NNRF_NFM] = models.NfType_NRF
	serviceNfType[models.ServiceName_NNRF_DISC] = models.NfType_NRF
	serviceNfType[models.ServiceName_NUDM_SDM] = models.NfType_UDM
	serviceNfType[models.ServiceName_NUDM_UECM] = models.NfType_UDM
	serviceNfType[models.ServiceName_NUDM_UEAU] = models.NfType_UDM
	serviceNfType[models.ServiceName_NUDM_EE] = models.NfType_UDM
	serviceNfType[models.ServiceName_NUDM_PP] = models.NfType_UDM
	serviceNfType[models.ServiceName_NAMF_COMM] = models.NfType_AMF
	serviceNfType[models.ServiceName_NAMF_EVTS] = models.NfType_AMF
	serviceNfType[models.ServiceName_NAMF_MT] = models.NfType_AMF
	serviceNfType[models.ServiceName_NAMF_LOC] = models.NfType_AMF
	serviceNfType[models.ServiceName_NSMF_PDUSESSION] = models.NfType_SMF
	serviceNfType[models.ServiceName_NSMF_EVENT_EXPOSURE] = models.NfType_SMF
	serviceNfType[models.ServiceName_NAUSF_AUTH] = models.NfType_AUSF
	serviceNfType[models.ServiceName_NAUSF_SORPROTECTION] = models.NfType_AUSF
	serviceNfType[models.ServiceName_NAUSF_UPUPROTECTION] = models.NfType_AUSF
	serviceNfType[models.ServiceName_NNEF_PFDMANAGEMENT] = models.NfType_NEF
	serviceNfType[models.ServiceName_NPCF_AM_POLICY_CONTROL] = models.NfType_PCF
	serviceNfType[models.ServiceName_NPCF_SMPOLICYCONTROL] = models.NfType_PCF
	serviceNfType[models.ServiceName_NPCF_POLICYAUTHORIZATION] = models.NfType_PCF
	serviceNfType[models.ServiceName_NPCF_BDTPOLICYCONTROL] = models.NfType_PCF
	serviceNfType[models.ServiceName_NPCF_EVENTEXPOSURE] = models.NfType_PCF
	serviceNfType[models.ServiceName_NPCF_UE_POLICY_CONTROL] = models.NfType_PCF
	serviceNfType[models.ServiceName_NSMSF_SMS] = models.NfType_SMSF
	serviceNfType[models.ServiceName_NNSSF_NSSELECTION] = models.NfType_NSSF
	serviceNfType[models.ServiceName_NNSSF_NSSAIAVAILABILITY] = models.NfType_NSSF
	serviceNfType[models.ServiceName_NUDR_DR] = models.NfType_UDR
	serviceNfType[models.ServiceName_NLMF_LOC] = models.NfType_LMF
	serviceNfType[models.ServiceName_N5G_EIR_EIC] = models.NfType__5_G_EIR
	serviceNfType[models.ServiceName_NBSF_MANAGEMENT] = models.NfType_BSF
	serviceNfType[models.ServiceName_NCHF_SPENDINGLIMITCONTROL] = models.NfType_CHF
	serviceNfType[models.ServiceName_NCHF_CONVERGEDCHARGING] = models.NfType_CHF
	serviceNfType[models.ServiceName_NNWDAF_EVENTSSUBSCRIPTION] = models.NfType_NWDAF
	serviceNfType[models.ServiceName_NNWDAF_ANALYTICSINFO] = models.NfType_NWDAF
}

type nnrfService struct {
	consumer *Consumer

	nfDiscMu      sync.RWMutex
	nfDiscClients map[string]*Nnrf_NFDiscovery.APIClient

	nfMngmntMu      sync.RWMutex
	nfMngmntClients map[string]*Nnrf_NFManagement.APIClient
}

func (s *nnrfService) getNFDiscoveryClient(uri string) *Nnrf_NFDiscovery.APIClient {
	s.nfDiscMu.RLock()
	if client, ok := s.nfDiscClients[uri]; ok {
		defer s.nfDiscMu.RUnlock()
		return client
	} else {
		configuration := Nnrf_NFDiscovery.NewConfiguration()
		configuration.SetBasePath(uri)
		cli := Nnrf_NFDiscovery.NewAPIClient(configuration)

		s.nfDiscMu.RUnlock()
		s.nfDiscMu.Lock()
		defer s.nfDiscMu.Unlock()
		s.nfDiscClients[uri] = cli
		return cli
	}
}

func (s *nnrfService) getNFManagementClient(uri string) *Nnrf_NFManagement.APIClient {
	s.nfMngmntMu.RLock()
	if client, ok := s.nfMngmntClients[uri]; ok {
		defer s.nfMngmntMu.RUnlock()
		return client
	} else {
		configuration := Nnrf_NFManagement.NewConfiguration()
		configuration.SetBasePath(uri)
		cli := Nnrf_NFManagement.NewAPIClient(configuration)

		s.nfMngmntMu.RUnlock()
		s.nfMngmntMu.Lock()
		defer s.nfMngmntMu.Unlock()
		s.nfMngmntClients[uri] = cli
		return cli
	}
}

func (s *nnrfService) RegisterNFInstance() error {
	var rsp *http.Response
	var nf models.NfProfile
	var err error

	client := s.getNFManagementClient(s.consumer.Config().NrfUri())
	nfProfile, err := s.buildNfProfile()
	if err != nil {
		return fmt.Errorf("RegisterNFInstance err: %+v", err)
	}

	for {
		nf, rsp, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(
			context.TODO(), s.consumer.Context().NfInstID(), *nfProfile)
		if rsp != nil && rsp.Body != nil {
			if bodyCloseErr := rsp.Body.Close(); bodyCloseErr != nil {
				logger.ConsumerLog.Errorf("response body cannot close: %+v", bodyCloseErr)
			}
		}

		if err != nil || rsp == nil {
			logger.ConsumerLog.Infof("NEF register to NRF Error[%v], sleep 2s and retry", err)
			time.Sleep(RetryRegisterNrfDuration)
			continue
		}

		status := rsp.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			logger.ConsumerLog.Infof("NFRegister Update")
			break
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := rsp.Header.Get("Location")
			// resouceNrfUri := resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			s.consumer.Context().SetNfInstID(resourceUri[strings.LastIndex(resourceUri, "/")+1:])

			oauth2 := false
			if nf.CustomInfo != nil {
				v, ok := nf.CustomInfo["oauth2"].(bool)
				if ok {
					oauth2 = v
					logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
				}
			}
			s.consumer.Context().OAuth2Required = oauth2
			if oauth2 && s.consumer.Context().Config().NrfCertPem() == "" {
				logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
			}

			logger.ConsumerLog.Infof("NFRegister Created")
			break
		} else {
			logger.ConsumerLog.Infof("NRF return wrong status: %d", status)
		}
	}
	return nil
}

func (s *nnrfService) buildNfProfile() (*models.NfProfile, error) {
	profile := &models.NfProfile{
		NfInstanceId: s.consumer.Context().NfInstID(),
		NfType:       models.NfType_NEF,
		NfStatus:     models.NfStatus_REGISTERED,
	}

	cfg := s.consumer.Config()
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, cfg.SbiRegisterIP())
	nfServices := cfg.NFServices()
	if len(nfServices) == 0 {
		return nil, fmt.Errorf("buildNfProfile err: NFServices is Empty")
	}
	profile.NfServices = &nfServices
	return profile, nil
}

func (s *nnrfService) DeregisterNFInstance() error {
	logger.ConsumerLog.Infof("DeregisterNFInstance")

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NfType_NRF)
	if err != nil {
		return nil
	}

	client := s.getNFManagementClient(s.consumer.Config().NrfUri())

	rsp, err := client.NFInstanceIDDocumentApi.DeregisterNFInstance(
		ctx, s.consumer.Context().NfInstID())
	if rsp != nil && rsp.Body != nil {
		if bodyCloseErr := rsp.Body.Close(); bodyCloseErr != nil {
			logger.ConsumerLog.Errorf("response body cannot close: %+v", bodyCloseErr)
		}
	}
	if err != nil {
		if rsp == nil {
			return fmt.Errorf("DeregisterNFInstance Error: server no response")
		} else if rsp.Status != err.Error() {
			return fmt.Errorf("DeregisterNFInstance Error[%+v]", err)
		}
		pd := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		return fmt.Errorf("DeregisterNFInstance Failed: Problem[%+v]", pd)
	}
	return nil
}

func (s *nnrfService) SearchNFInstances(
	nrfUri string,
	srvName models.ServiceName,
	param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts,
) (*models.NfProfile, string, error) {
	if param == nil {
		param = &Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}
	}
	param.ServiceNames = optional.NewInterface([]models.ServiceName{srvName})

	client := s.getNFDiscoveryClient(nrfUri)

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NfType_NRF)
	if err != nil {
		return nil, "", err
	}

	res, rsp, err := client.NFInstancesStoreApi.SearchNFInstances(ctx,
		serviceNfType[srvName], models.NfType_NEF, param)
	if rsp != nil && rsp.Body != nil {
		if bodyCloseErr := rsp.Body.Close(); bodyCloseErr != nil {
			logger.ConsumerLog.Errorf("SearchNFInstances err: response body cannot close: %+v", bodyCloseErr)
		}
	}
	if rsp != nil && rsp.StatusCode == http.StatusTemporaryRedirect {
		err = fmt.Errorf("SearchNFInstances err: Temporary Redirect")
	}
	if err != nil {
		return nil, "", err
	}

	nfProf, uri, err := getProfileAndUri(res.NfInstances, srvName)
	if err != nil {
		logger.ConsumerLog.Errorf("%s", err.Error())
		return nil, "", err
	}
	return nfProf, uri, nil
}

func getProfileAndUri(nfInstances []models.NfProfile, srvName models.ServiceName) (*models.NfProfile, string, error) {
	// select the first ServiceName
	// TODO: select base on other info
	var profile *models.NfProfile
	var uri string
	for _, nfProfile := range nfInstances {
		profile = &nfProfile
		uri = searchNFServiceUri(nfProfile, srvName, models.NfServiceStatus_REGISTERED)
		if uri != "" {
			break
		}
	}
	if uri == "" {
		return nil, "", fmt.Errorf("no uri for %s found", srvName)
	}
	return profile, uri, nil
}

// searchNFServiceUri returns NF Uri derived from NfProfile with corresponding service
func searchNFServiceUri(nfProfile models.NfProfile, serviceName models.ServiceName,
	nfServiceStatus models.NfServiceStatus,
) string {
	if nfProfile.NfServices == nil {
		return ""
	}

	nfUri := ""
	for _, service := range *nfProfile.NfServices {
		if service.ServiceName == serviceName && service.NfServiceStatus == nfServiceStatus {
			if service.Fqdn != "" {
				nfUri = string(service.Scheme) + "://" + service.Fqdn
			} else if nfProfile.Fqdn != "" {
				nfUri = string(service.Scheme) + "://" + nfProfile.Fqdn
			} else if service.ApiPrefix != "" {
				u, err := url.Parse(service.ApiPrefix)
				if err != nil {
					return nfUri
				}
				nfUri = u.Scheme + "://" + u.Host
			} else if len(*service.IpEndPoints) != 0 {
				// Select the first IpEndPoint
				// TODO: select others when failure
				point := (*service.IpEndPoints)[0]
				if point.Ipv4Address != "" {
					nfUri = getUriFromIpEndPoint(service.Scheme, point.Ipv4Address, point.Port)
				} else if len(nfProfile.Ipv4Addresses) != 0 {
					nfUri = getUriFromIpEndPoint(service.Scheme, nfProfile.Ipv4Addresses[0], point.Port)
				}
			}
		}
		if nfUri != "" {
			break
		}
	}

	return nfUri
}

func getUriFromIpEndPoint(scheme models.UriScheme, ipv4Address string, port int32) string {
	uri := ""
	if port != 0 {
		uri = string(scheme) + "://" + ipv4Address + ":" + strconv.Itoa(int(port))
	} else {
		switch scheme {
		case models.UriScheme_HTTP:
			uri = string(scheme) + "://" + ipv4Address + ":80"
		case models.UriScheme_HTTPS:
			uri = string(scheme) + "://" + ipv4Address + ":443"
		}
	}
	return uri
}
