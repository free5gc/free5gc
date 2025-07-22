package app

import (
	udm_context "github.com/free5gc/udm/internal/context"
	"github.com/free5gc/udm/pkg/factory"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *udm_context.UDMContext
	Config() *factory.Config
}
