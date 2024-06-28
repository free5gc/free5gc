package app

import (
	"context"
)

type NetworkFunction interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)
	Start()
	Terminate()
}

type NFstruct struct {
	Nf     NetworkFunction
	Ctx    *context.Context
	Cancel *context.CancelFunc
}
