package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/logger"
	"github.com/free5gc/bsf/pkg/app"
	"github.com/free5gc/bsf/pkg/factory"
)

type BsfApp struct {
	cfg        *factory.Config
	ctx        context.Context
	cancel     context.CancelFunc
	tlsKeyPath string
	app        *app.App
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*BsfApp, error) {
	bsf := &BsfApp{
		cfg:        cfg,
		tlsKeyPath: tlsKeyLogPath,
	}
	bsf.SetLogEnable(cfg.Logger.Enable)
	bsf.SetLogLevel(cfg.Logger.Level)
	bsf.SetReportCaller(cfg.Logger.ReportCaller)
	bsf.ctx, bsf.cancel = context.WithCancel(ctx)

	bsfContext.InitBsfContext()
	var err error
	if bsf.app, err = app.NewApp(bsf.ctx, bsf.cfg); err != nil {
		return nil, fmt.Errorf("failed to create BSF app: %w", err)
	}
	return bsf, nil
}

func (bsf *BsfApp) Start() {
	defer func() {
		if p := recover(); p != nil {
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	logger.MainLog.Infoln("BSF started")

	if err := factory.CheckConfigVersion(); err != nil {
		logger.MainLog.Warnf("Config version error: %v", err)
	}

	// Start the app and return error directly instead of using goroutine
	bsf.app.Start()
}

func (bsf *BsfApp) Terminate() {
	logger.MainLog.Infof("Terminating BSF...")

	// Terminate app first
	if err := bsf.app.Terminate(); err != nil {
		logger.MainLog.Errorf("Error terminating app: %+v", err)
		// Continue with cleanup even if app termination fails
	}

	bsf.cancel()
	logger.MainLog.Infof("BSF App is terminated")
}

func (bsf *BsfApp) Config() *factory.Config {
	return bsf.cfg
}

func (bsf *BsfApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == io.Discard {
		return
	}

	bsf.Config().SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(io.Discard)
	}
}

func (bsf *BsfApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}

	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}

	bsf.Config().SetLogLevel(level)
	logger.Log.SetLevel(lvl)
}

func (bsf *BsfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	bsf.Config().SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (bsf *BsfApp) GetContext() context.Context {
	return bsf.ctx
}
