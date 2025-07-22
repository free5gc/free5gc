package consumer

import (
	"sync"

	Nudm_SubscriberDataManagement "github.com/free5gc/openapi/udm/SubscriberDataManagement"
	Nudm_UEContextManagement "github.com/free5gc/openapi/udm/UEContextManagement"
)

type nudmService struct {
	consumer *Consumer

	nfSDMMu  sync.RWMutex
	nfUECMMu sync.RWMutex

	nfSDMClients  map[string]*Nudm_SubscriberDataManagement.APIClient
	nfUECMClients map[string]*Nudm_UEContextManagement.APIClient
}

func (s *nudmService) GetSDMClient(uri string) *Nudm_SubscriberDataManagement.APIClient {
	if uri == "" {
		return nil
	}
	s.nfSDMMu.RLock()
	client, ok := s.nfSDMClients[uri]
	if ok {
		s.nfSDMMu.RUnlock()
		return client
	}

	configuration := Nudm_SubscriberDataManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nudm_SubscriberDataManagement.NewAPIClient(configuration)

	s.nfSDMMu.RUnlock()
	s.nfSDMMu.Lock()
	defer s.nfSDMMu.Unlock()
	s.nfSDMClients[uri] = client
	return client
}

func (s *nudmService) GetUECMClient(uri string) *Nudm_UEContextManagement.APIClient {
	if uri == "" {
		return nil
	}
	s.nfUECMMu.RLock()
	client, ok := s.nfUECMClients[uri]
	if ok {
		defer s.nfUECMMu.RUnlock()
		return client
	}

	configuration := Nudm_UEContextManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nudm_UEContextManagement.NewAPIClient(configuration)

	s.nfUECMMu.RUnlock()
	s.nfUECMMu.Lock()
	defer s.nfUECMMu.Unlock()
	s.nfUECMClients[uri] = client
	return client
}
