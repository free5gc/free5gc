package processor

import (
	nef_context "github.com/free5gc/nef/internal/context"
	"github.com/free5gc/nef/internal/sbi/consumer"
	"github.com/free5gc/nef/internal/sbi/notifier"
	"github.com/free5gc/nef/pkg/factory"
)

type nef interface {
	Context() *nef_context.NefContext
	Config() *factory.Config
	Consumer() *consumer.Consumer
	Notifier() *notifier.Notifier
}

type Processor struct {
	nef
}

type HandlerResponse struct {
	Status  int
	Headers map[string][]string
	Body    interface{}
}

func NewProcessor(nef nef) (*Processor, error) {
	handler := &Processor{
		nef: nef,
	}

	return handler, nil
}

func addLocationheader(header map[string][]string, location string) {
	locations := header["Location"]
	if locations == nil {
		header["Location"] = []string{location}
	} else {
		header["Location"] = append(locations, location)
	}
}
