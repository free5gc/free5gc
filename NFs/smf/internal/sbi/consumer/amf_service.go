package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/free5gc/openapi/amf/Communication"
	"github.com/free5gc/openapi/models"
)

type namfService struct {
	consumer *Consumer

	CommunicationMu sync.RWMutex

	CommunicationClients map[string]*Communication.APIClient
}

func (s *namfService) getCommunicationClient(uri string) *Communication.APIClient {
	if uri == "" {
		return nil
	}
	s.CommunicationMu.RLock()
	client, ok := s.CommunicationClients[uri]
	if ok {
		s.CommunicationMu.RUnlock()
		return client
	}

	configuration := Communication.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Communication.NewAPIClient(configuration)

	s.CommunicationMu.RUnlock()
	s.CommunicationMu.Lock()
	defer s.CommunicationMu.Unlock()
	s.CommunicationClients[uri] = client
	return client
}

func (s *namfService) N1N2MessageTransfer(
	ctx context.Context, supi string, n1n2Request models.N1N2MessageTransferRequest, apiPrefix string,
) (*models.N1N2MessageTransferRspData, error) {
	client := s.getCommunicationClient(apiPrefix)
	if client == nil {
		return nil, fmt.Errorf("N1N2MessageTransfer client is nil: (%v)", apiPrefix)
	}

	n1n2MessageTransferRequest := &Communication.N1N2MessageTransferRequest{
		UeContextId:                &supi,
		N1N2MessageTransferRequest: &n1n2Request,
	}

	rsp, err := client.N1N2MessageCollectionCollectionApi.N1N2MessageTransfer(ctx, n1n2MessageTransferRequest)
	if err != nil || rsp == nil {
		return nil, err
	}

	return &rsp.N1N2MessageTransferRspData, err
}
