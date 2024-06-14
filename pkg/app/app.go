package app

import (
	nrf_context "github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/pkg/factory"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *nrf_context.NRFContext
	Config() *factory.Config
}
