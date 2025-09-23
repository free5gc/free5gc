package consumer

import (
	"context"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/pkg/factory"
)

var consumer *Consumer

type ConsumerBsf interface {
	Config() *factory.Config
	Context() *bsfContext.BSFContext
	CancelContext() context.Context
}

type Consumer struct {
	ConsumerBsf

	// consumer services
	*nnrfService
}

type nnrfService struct {
	consumer *Consumer
}

func GetConsumer() *Consumer {
	return consumer
}

func NewConsumer(bsf ConsumerBsf) (*Consumer, error) {
	c := &Consumer{
		ConsumerBsf: bsf,
	}

	c.nnrfService = &nnrfService{
		consumer: c,
	}

	consumer = c
	return c, nil
}

// RegisterWithNRF calls the existing NRF service registration function
func (c *Consumer) RegisterWithNRF(ctx context.Context) error {
	_, err := SendRegisterNFInstance(ctx)
	return err
}

// DeregisterWithNRF calls the existing NRF service deregistration function
func (c *Consumer) DeregisterWithNRF() error {
	_, err := SendDeregisterNFInstance()
	return err
}
