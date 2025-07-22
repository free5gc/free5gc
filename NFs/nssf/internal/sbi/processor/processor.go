package processor

import (
	"github.com/free5gc/nssf/pkg/app"
)

type Processor struct {
	app.NssfApp
}

func NewProcessor(nssf app.NssfApp) *Processor {
	return &Processor{
		NssfApp: nssf,
	}
}
