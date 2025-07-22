package consumer

import (
	"github.com/free5gc/openapi/nrf/NFManagement"
	"github.com/free5gc/udr/pkg/app"
)

type Consumer struct {
	app.App

	*NrfService
}

func NewConsumer(udr app.App) *Consumer {
	configuration := NFManagement.NewConfiguration()
	configuration.SetBasePath(udr.Context().NrfUri)
	nrfService := &NrfService{
		nfMngmntClients: make(map[string]*NFManagement.APIClient),
	}

	return &Consumer{
		App:        udr,
		NrfService: nrfService,
	}
}
