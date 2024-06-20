package consumer

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/models"
)

type nnrfService struct {
	consumer *Consumer

	nfMngmntMu sync.RWMutex

	nfMngmntClients map[string]*Nnrf_NFManagement.APIClient
}

func (s *nnrfService) getNFManagementClient(uri string) *Nnrf_NFManagement.APIClient {
	if uri == "" {
		return nil
	}
	s.nfMngmntMu.RLock()
	client, ok := s.nfMngmntClients[uri]
	if ok {
		defer s.nfMngmntMu.RUnlock()
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

func (s *nnrfService) SendNFStatusNotify(
	notification_event models.NotificationEventType,
	nfInstanceUri string,
	url string,
	nfProfile *models.NfProfile,
) *models.ProblemDetails {
	logger.ConsumerLog.Infoln("SendNFStatusNotify")

	client := s.getNFManagementClient(url)
	if client == nil {
		return &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "NEW_CLIENT_ERROR",
			Detail: fmt.Sprintf("Can't Get/New Client for url for [%+v]", url),
		}
	}
	s.nfMngmntMu.RLock()
	defer s.nfMngmntMu.RUnlock()

	notifcationData := models.NotificationData{
		Event:         notification_event,
		NfInstanceUri: nfInstanceUri,
	}
	if nfProfile != nil {
		buildNotificationDataFromNfProfile(notifcationData.NfProfile, nfProfile)
	}

	res, err := client.NotificationApi.NotificationPost(context.TODO(), notifcationData)
	if err != nil {
		logger.NfmLog.Infof("Notify fail: %v", err)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "NOTIFICATION_ERROR",
			Detail: err.Error(),
		}
		return problemDetails
	}
	if res != nil {
		defer func() {
			if resCloseErr := res.Body.Close(); resCloseErr != nil {
				logger.NfmLog.Errorf("NotificationApi response body cannot close: %+v", resCloseErr)
			}
		}()
		if status := res.StatusCode; status != http.StatusNoContent {
			logger.NfmLog.Warnln("Error status in NotificationPost: ", status)
			problemDetails := &models.ProblemDetails{
				Status: int32(status),
				Cause:  "NOTIFICATION_ERROR",
			}
			return problemDetails
		}
	}
	return nil
}

func buildNotificationDataFromNfProfile(notifProfile *models.NfProfileNotificationData, nfProfile *models.NfProfile) {
	notifProfile.NfInstanceId = nfProfile.NfInstanceId
	notifProfile.NfType = nfProfile.NfType
	notifProfile.NfStatus = nfProfile.NfStatus
	notifProfile.HeartBeatTimer = nfProfile.HeartBeatTimer
	notifProfile.PlmnList = *nfProfile.PlmnList
	notifProfile.SNssais = *nfProfile.SNssais
	notifProfile.PerPlmnSnssaiList = nfProfile.PerPlmnSnssaiList
	notifProfile.NsiList = nfProfile.NsiList
	notifProfile.Fqdn = nfProfile.Fqdn
	notifProfile.InterPlmnFqdn = nfProfile.InterPlmnFqdn
	notifProfile.Ipv4Addresses = nfProfile.Ipv4Addresses
	notifProfile.Ipv6Addresses = nfProfile.Ipv6Addresses
	notifProfile.AllowedPlmns = *nfProfile.AllowedPlmns
	notifProfile.AllowedNfTypes = nfProfile.AllowedNfTypes
	notifProfile.AllowedNfDomains = nfProfile.AllowedNfDomains
	notifProfile.AllowedNssais = *nfProfile.AllowedNssais
	notifProfile.Priority = nfProfile.Priority
	notifProfile.Capacity = nfProfile.Capacity
	notifProfile.Load = nfProfile.Load
	notifProfile.Locality = nfProfile.Locality
	notifProfile.UdrInfo = nfProfile.UdrInfo
	notifProfile.UdmInfo = nfProfile.UdmInfo
	notifProfile.AusfInfo = nfProfile.AusfInfo
	notifProfile.AmfInfo = nfProfile.AmfInfo
	notifProfile.SmfInfo = nfProfile.SmfInfo
	notifProfile.UpfInfo = nfProfile.UpfInfo
	notifProfile.PcfInfo = nfProfile.PcfInfo
	notifProfile.BsfInfo = nfProfile.BsfInfo
	notifProfile.ChfInfo = nfProfile.ChfInfo
	notifProfile.NrfInfo = nfProfile.NrfInfo
	notifProfile.CustomInfo = nfProfile.CustomInfo
	notifProfile.RecoveryTime = nfProfile.RecoveryTime
	notifProfile.NfServicePersistence = nfProfile.NfServicePersistence
	notifProfile.NfServices = *nfProfile.NfServices
	notifProfile.NfProfileChangesSupportInd = nfProfile.NfProfileChangesSupportInd
	notifProfile.NfProfileChangesInd = nfProfile.NfProfileChangesInd
	notifProfile.DefaultNotificationSubscriptions = nfProfile.DefaultNotificationSubscriptions
}
