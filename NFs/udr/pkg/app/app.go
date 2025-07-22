package app

import (
	udr_context "github.com/free5gc/udr/internal/context"
	"github.com/free5gc/udr/pkg/factory"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *udr_context.UDRContext
	Config() *factory.Config
}
