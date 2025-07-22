package processor

import (
	"github.com/free5gc/udm/internal/sbi/consumer"
	"github.com/free5gc/udm/pkg/app"
)

type ProcessorUdm interface {
	app.App

	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorUdm
}

func NewProcessor(udm ProcessorUdm) (*Processor, error) {
	p := &Processor{
		ProcessorUdm: udm,
	}
	return p, nil
}
