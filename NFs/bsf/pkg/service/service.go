package service

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

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
	wg         sync.WaitGroup
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*BsfApp, error) {
	bsf := &BsfApp{
		cfg:        cfg,
		tlsKeyPath: tlsKeyLogPath,
	}
	bsf.SetLogEnable(cfg.Logger.BSF)
	bsf.ctx, bsf.cancel = context.WithCancel(ctx)

	bsfContext.InitBsfContext()
	var err error
	if bsf.app, err = app.NewApp(bsf.cfg); err != nil {
		return nil, fmt.Errorf("failed to create BSF app: %w", err)
	}
	return bsf, nil
}

func (bsf *BsfApp) Start() error {
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
	return bsf.app.Start()
}

func (bsf *BsfApp) Terminate() error {
	logger.MainLog.Infof("Terminating BSF...")

	// Terminate app first
	if err := bsf.app.Terminate(); err != nil {
		logger.MainLog.Errorf("Error terminating app: %+v", err)
		// Continue with cleanup even if app termination fails
	}

	bsf.cancel()
	logger.MainLog.Infof("BSF App is terminated")
	return nil
}

func (bsf *BsfApp) SetLogEnable(cfg *factory.LogSetting) {
	logger.MainLog.Infof("BSF Log enable")
	if cfg.DebugLevel != "" {
		if level, err := logrus.ParseLevel(cfg.DebugLevel); err != nil {
			logger.MainLog.Warnf("BSF Log level [%s] is invalid, set to [info] level",
				cfg.DebugLevel)
			logger.Log.SetLevel(logrus.InfoLevel)
		} else {
			logger.MainLog.Infof("BSF Log level is set to [%s] level", level)
			logger.Log.SetLevel(level)
		}
	} else {
		logger.MainLog.Infoln("BSF Log level not set. Default set to [info] level")
		logger.Log.SetLevel(logrus.InfoLevel)
	}
	logger.Log.SetReportCaller(cfg.ReportCaller)
}

func (bsf *BsfApp) SetLogLevel(level logrus.Level) {
	logger.MainLog.Infof("set log level : %+v", level)
	logger.Log.SetLevel(level)
}

func (bsf *BsfApp) GetContext() context.Context {
	return bsf.ctx
}
