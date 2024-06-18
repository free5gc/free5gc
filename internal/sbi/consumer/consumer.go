package consumer

import (
	"github.com/free5gc/nrf/pkg/app"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
)

type ConsumerNrf interface {
	app.App
}

type Consumer struct {
	ConsumerNrf

	*nnrfService
}

func NewConsumer(nrf ConsumerNrf) (*Consumer, error) {
	c := &Consumer{
		ConsumerNrf: nrf,
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfMngmntClients: make(map[string]*Nnrf_NFManagement.APIClient),
	}
	return c, nil
}
