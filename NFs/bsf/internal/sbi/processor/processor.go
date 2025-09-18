package processor

import (
	"context"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/sbi/consumer"
	"github.com/free5gc/bsf/pkg/factory"
)

var processor *Processor

type ProcessorBsf interface {
	Config() *factory.Config
	Context() *bsfContext.BSFContext
	CancelContext() context.Context
	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorBsf
}

func GetProcessor() *Processor {
	return processor
}

func NewProcessor(bsf ProcessorBsf) (*Processor, error) {
	p := &Processor{
		ProcessorBsf: bsf,
	}
	processor = p
	return p, nil
}
