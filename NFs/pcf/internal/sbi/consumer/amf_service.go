package consumer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/amf/Communication"
	"github.com/free5gc/openapi/models"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/pkg/factory"
)

type namfService struct {
	consumer *Consumer

	nfComMu sync.RWMutex

	nfComClients map[string]*Communication.APIClient
}

func (s *namfService) getNFCommunicationClient(uri string) *Communication.APIClient {
	if uri == "" {
		return nil
	}
	s.nfComMu.RLock()
	client, ok := s.nfComClients[uri]
	if ok {
		defer s.nfComMu.RUnlock()
		return client
	}

	configuration := Communication.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Communication.NewAPIClient(configuration)

	s.nfComMu.RUnlock()
	s.nfComMu.Lock()
	defer s.nfComMu.Unlock()
	s.nfComClients[uri] = client
	return client
}

func (s *namfService) AmfStatusChangeSubscribe(amfUri string, guamiList []models.Guami) (
	problemDetails *models.ProblemDetails, err error,
) {
	logger.ConsumerLog.Debugf("PCF Subscribe to AMF status[%+v]", amfUri)
	pcfContext := s.consumer.pcf.Context()

	// Set client and set url
	client := s.getNFCommunicationClient(amfUri)

	subscriptionData := models.AmfCommunicationSubscriptionData{
		AmfStatusUri: fmt.Sprintf("%s"+factory.PcfCallbackResUriPrefix+"/amfstatus", pcfContext.GetIPv4Uri()),
		GuamiList:    guamiList,
	}
	amfStausChangeRequest := &Communication.AMFStatusChangeSubscribeRequest{}
	amfStausChangeRequest.SetAmfCommunicationSubscriptionData(subscriptionData)
	ctx, pd, err := pcfContext.GetTokenCtx(models.ServiceName_NAMF_COMM, models.NrfNfManagementNfType_AMF)
	if err != nil {
		return pd, err
	}
	res, localErr := client.SubscriptionsCollectionCollectionApi.AMFStatusChangeSubscribe(
		ctx, amfStausChangeRequest)

	if localErr == nil {
		locationHeader := res.Location
		logger.ConsumerLog.Debugf("location header: %+v", locationHeader)

		subscriptionID := locationHeader[strings.LastIndex(locationHeader, "/")+1:]
		amfStatusSubsData := pcf_context.AMFStatusSubscriptionData{
			AmfUri:       amfUri,
			AmfStatusUri: res.AmfCommunicationSubscriptionData.AmfStatusUri,
			GuamiList:    res.AmfCommunicationSubscriptionData.GuamiList,
		}
		pcfContext.NewAmfStatusSubscription(subscriptionID, amfStatusSubsData)
	}

	if genericErr, ok := localErr.(openapi.GenericOpenAPIError); ok {
		if problemDetails, ok := genericErr.Model().(models.ProblemDetails); ok {
			return &problemDetails, nil
		}

		logger.ConsumerLog.Errorf("openapi error: %+v", localErr)
		return nil, localErr
	}
	return nil, localErr
}
