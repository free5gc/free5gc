package service

import (
	"context"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"

	ausf_context "github.com/free5gc/ausf/internal/context"
	"github.com/free5gc/ausf/internal/logger"
	"github.com/free5gc/ausf/internal/sbi"
	"github.com/free5gc/ausf/internal/sbi/consumer"
	"github.com/free5gc/ausf/internal/sbi/processor"
	"github.com/free5gc/ausf/pkg/app"
	"github.com/free5gc/ausf/pkg/factory"
)

var AUSF *AusfApp

var _ app.App = &AusfApp{}

type AusfApp struct {
	ausfCtx *ausf_context.AUSFContext
	cfg     *factory.Config

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	sbiServer *sbi.Server
	consumer  *consumer.Consumer
	processor *processor.Processor
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*AusfApp, error) {
	ausf := &AusfApp{
		cfg: cfg,
		wg:  sync.WaitGroup{},
	}
	ausf.SetLogEnable(cfg.GetLogEnable())
	ausf.SetLogLevel(cfg.GetLogLevel())
	ausf.SetReportCaller(cfg.GetLogReportCaller())
	ausf_context.Init()

	processor, err_p := processor.NewProcessor(ausf)
	if err_p != nil {
		return ausf, err_p
	}
	ausf.processor = processor

	consumer, err := consumer.NewConsumer(ausf)
	if err != nil {
		return ausf, err
	}
	ausf.consumer = consumer

	ausf.ctx, ausf.cancel = context.WithCancel(ctx)
	ausf.ausfCtx = ausf_context.GetSelf()

	if ausf.sbiServer, err = sbi.NewServer(ausf, tlsKeyLogPath); err != nil {
		return nil, err
	}
	AUSF = ausf

	return ausf, nil
}

func (a *AusfApp) CancelContext() context.Context {
	return a.ctx
}

func (a *AusfApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *AusfApp) Processor() *processor.Processor {
	return a.processor
}

func (a *AusfApp) Context() *ausf_context.AUSFContext {
	return a.ausfCtx
}

func (a *AusfApp) Config() *factory.Config {
	return a.cfg
}

func (a *AusfApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == io.Discard {
		return
	}

	a.Config().SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(io.Discard)
	}
}

func (a *AusfApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}

	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}

	a.Config().SetLogLevel(level)
	logger.Log.SetLevel(lvl)
}

func (a *AusfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	a.Config().SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *AusfApp) Start() {
	logger.InitLog.Infoln("Server started")

	a.wg.Add(1)
	go a.listenShutdownEvent()

	if err := a.sbiServer.Run(context.Background(), &a.wg); err != nil {
		logger.MainLog.Fatalf("Run SBI server failed: %+v", err)
	}
	a.WaitRoutineStopped()
}

func (a *AusfApp) listenShutdownEvent() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		a.wg.Done()
	}()

	<-a.ctx.Done()
	a.terminateProcedure()
}

func (a *AusfApp) Terminate() {
	a.cancel()
}

func (a *AusfApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating AUSF...")
	a.CallServerStop()

	// deregister with NRF
	problemDetails, err := a.Consumer().SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.MainLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.MainLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.MainLog.Infof("Deregister from NRF successfully")
	}
	logger.MainLog.Infof("CHF SBI Server terminated")
}

func (a *AusfApp) CallServerStop() {
	if a.sbiServer != nil {
		a.sbiServer.Shutdown()
	}
}

func (a *AusfApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.MainLog.Infof("AUSF App is terminated")
}
