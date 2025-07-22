package processor

import (
	"github.com/free5gc/smf/internal/sbi/consumer"
	"github.com/free5gc/smf/pkg/app"
)

const (
	CONTEXT_NOT_FOUND = "CONTEXT_NOT_FOUND"
)

type ProcessorSmf interface {
	app.App

	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorSmf
}

func NewProcessor(smf ProcessorSmf) (*Processor, error) {
	p := &Processor{
		ProcessorSmf: smf,
	}
	return p, nil
}
