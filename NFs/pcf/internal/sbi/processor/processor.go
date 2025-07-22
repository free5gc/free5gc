package processor

import (
	"github.com/free5gc/pcf/internal/sbi/consumer"
	"github.com/free5gc/pcf/pkg/app"
)

type PCF interface {
	app.App
	Consumer() *consumer.Consumer
}

type Processor struct {
	PCF
}

func NewProcessor(pcf PCF) (*Processor, error) {
	p := &Processor{
		PCF: pcf,
	}

	return p, nil
}
