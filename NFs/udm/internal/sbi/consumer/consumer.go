package consumer

import (
	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	Nnrf_NFManagement "github.com/free5gc/openapi/nrf/NFManagement"
	Nudm_SubscriberDataManagement "github.com/free5gc/openapi/udm/SubscriberDataManagement"
	Nudm_UEContextManagement "github.com/free5gc/openapi/udm/UEContextManagement"
	Nudr_DataRepository "github.com/free5gc/openapi/udr/DataRepository"
	"github.com/free5gc/udm/pkg/app"
)

type ConsumerUdm interface {
	app.App
}

type Consumer struct {
	ConsumerUdm

	// consumer services
	*nnrfService
	*nudrService
	*nudmService
}

func NewConsumer(udm ConsumerUdm) (*Consumer, error) {
	c := &Consumer{
		ConsumerUdm: udm,
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfMngmntClients: make(map[string]*Nnrf_NFManagement.APIClient),
		nfDiscClients:   make(map[string]*Nnrf_NFDiscovery.APIClient),
	}

	c.nudrService = &nudrService{
		consumer:    c,
		nfDRClients: make(map[string]*Nudr_DataRepository.APIClient),
	}

	c.nudmService = &nudmService{
		consumer:      c,
		nfSDMClients:  make(map[string]*Nudm_SubscriberDataManagement.APIClient),
		nfUECMClients: make(map[string]*Nudm_UEContextManagement.APIClient),
	}
	return c, nil
}
