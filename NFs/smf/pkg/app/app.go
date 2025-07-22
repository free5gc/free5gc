package app

import (
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/pkg/factory"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *smf_context.SMFContext
	Config() *factory.Config
}
