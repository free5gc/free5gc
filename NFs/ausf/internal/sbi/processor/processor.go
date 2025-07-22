package processor

import (
	"github.com/free5gc/ausf/internal/sbi/consumer"
	"github.com/free5gc/ausf/pkg/app"
)

type ProcessorAusf interface {
	app.App

	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorAusf
}

func NewProcessor(ausf ProcessorAusf) (*Processor, error) {
	p := &Processor{
		ProcessorAusf: ausf,
	}
	return p, nil
}
