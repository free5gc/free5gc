package consumer

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/nrf/NFManagement"
)

type nnrfService struct {
	consumer *Consumer

	nfMngmntMu sync.RWMutex

	nfMngmntClients map[string]*NFManagement.APIClient
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

func (s *nnrfService) SendNFStatusNotify(
	ctx context.Context,
	notification_event models.NotificationEventType,
	nfInstanceUri string,
	url string,
	nfProfile *models.NrfNfManagementNfProfile,
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

	notifcationData := models.NrfNfManagementNotificationData{
		Event:         notification_event,
		NfInstanceUri: nfInstanceUri,
	}
	if nfProfile != nil {
		notifcationData.NfProfile = nfProfile
	}

	request := &NFManagement.CreateSubscriptionOnNFStatusEventPostRequest{
		NrfNfManagementNotificationData: &notifcationData,
	}

	res, err := client.SubscriptionsCollectionApi.CreateSubscriptionOnNFStatusEventPost(
		ctx, nfInstanceUri, request)
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
		if res.AcceptEncoding != "" && res.AcceptEncoding != "application/json" {
			logger.NfmLog.Warnln("Unexpected Accept-Encoding in response: ", res.AcceptEncoding)
			problemDetails := &models.ProblemDetails{
				Status: http.StatusUnsupportedMediaType,
				Cause:  "NOTIFICATION_ERROR",
			}
			return problemDetails
		}
	}
	return nil
}
