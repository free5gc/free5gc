package processor

import (
	"github.com/free5gc/udr/internal/database"
	"github.com/free5gc/udr/pkg/app"
)

type Processor struct {
	app.App
	database.DbConnector
}

func NewProcessor(udr app.App) *Processor {
	return &Processor{
		App:         udr,
		DbConnector: database.NewDbConnector(udr.Config().Configuration.DbConnectorType),
	}
}
