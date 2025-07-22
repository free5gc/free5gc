package consumer

import (
	"context"

	"github.com/free5gc/openapi/amf/Communication"
	"github.com/free5gc/openapi/nrf/NFDiscovery"
	"github.com/free5gc/openapi/nrf/NFManagement"
	"github.com/free5gc/openapi/pcf/AMPolicyControl"
	"github.com/free5gc/openapi/udr/DataRepository"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/pkg/factory"
)

type pcf interface {
	Config() *factory.Config
	Context() *pcf_context.PCFContext
	CancelContext() context.Context
}

type Consumer struct {
	pcf

	// consumer services
	*nnrfService
	*namfService
	*nudrService
	*npcfService
}

func NewConsumer(pcf pcf) (*Consumer, error) {
	c := &Consumer{
		pcf: pcf,
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfMngmntClients: make(map[string]*NFManagement.APIClient),
		nfDiscClients:   make(map[string]*NFDiscovery.APIClient),
	}

	c.namfService = &namfService{
		consumer:     c,
		nfComClients: make(map[string]*Communication.APIClient),
	}

	c.nudrService = &nudrService{
		consumer:         c,
		nfDataSubClients: make(map[string]*DataRepository.APIClient),
	}

	c.npcfService = &npcfService{
		consumer:                c,
		nfAMPolicyControlClient: make(map[string]*AMPolicyControl.APIClient),
	}

	return c, nil
}
