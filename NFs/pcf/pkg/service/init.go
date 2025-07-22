package service

import (
	"context"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"

	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/sbi"
	"github.com/free5gc/pcf/internal/sbi/consumer"
	"github.com/free5gc/pcf/internal/sbi/processor"
	"github.com/free5gc/pcf/pkg/app"
	"github.com/free5gc/pcf/pkg/factory"
)

var PCF *PcfApp

var _ app.App = &PcfApp{}

type PcfApp struct {
	app.App
	cfg    *factory.Config
	pcfCtx *pcf_context.PCFContext
	ctx    context.Context
	cancel context.CancelFunc

	consumer  *consumer.Consumer
	processor *processor.Processor
	sbiServer *sbi.Server
	wg        sync.WaitGroup
}

func NewApp(
	ctx context.Context,
	cfg *factory.Config,
	tlsKeyLogPath string,
) (*PcfApp, error) {
	pcf := &PcfApp{
		cfg: cfg,
		wg:  sync.WaitGroup{},
	}
	pcf.SetLogEnable(cfg.GetLogEnable())
	pcf.SetLogLevel(cfg.GetLogLevel())
	pcf.SetReportCaller(cfg.GetLogReportCaller())

	pcf.ctx, pcf.cancel = context.WithCancel(ctx)
	pcf_context.Init()
	pcf.pcfCtx = pcf_context.GetSelf()

	// consumer
	consumer, err := consumer.NewConsumer(pcf)
	if err != nil {
		return pcf, err
	}
	pcf.consumer = consumer

	// processor
	p, err := processor.NewProcessor(pcf)
	if err != nil {
		return pcf, err
	}
	pcf.processor = p

	if pcf.sbiServer, err = sbi.NewServer(pcf, tlsKeyLogPath); err != nil {
		return nil, err
	}
	PCF = pcf

	return pcf, nil
}

func (a *PcfApp) Config() *factory.Config {
	return a.cfg
}

func (a *PcfApp) Context() *pcf_context.PCFContext {
	return a.pcfCtx
}

func (a *PcfApp) CancelContext() context.Context {
	return a.ctx
}

func (a *PcfApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *PcfApp) Processor() *processor.Processor {
	return a.processor
}

func (a *PcfApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == io.Discard {
		return
	}

	a.cfg.SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(io.Discard)
	}
}

func (a *PcfApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}

	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}

	a.cfg.SetLogLevel(level)
	logger.Log.SetLevel(lvl)
}

func (a *PcfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *PcfApp) Start() {
	logger.InitLog.Infoln("Server started")
	a.wg.Add(1)
	go a.listenShutdownEvent()
	if err := a.sbiServer.Run(context.Background(), &a.wg); err != nil {
		logger.InitLog.Fatalf("Run SBI server failed: %+v", err)
	}
	a.WaitRoutineStopped()
}

func (a *PcfApp) listenShutdownEvent() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		a.wg.Done()
	}()

	<-a.ctx.Done()
	a.terminateProcedure()
}

func (a *PcfApp) CallServerStop() {
	if a.sbiServer != nil {
		a.sbiServer.Shutdown(context.Background())
	}
}

func (a *PcfApp) Terminate() {
	a.cancel()
}

func (a *PcfApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating PCF...")
	a.CallServerStop()
	// deregister with NRF
	problemDetails, err := a.Consumer().SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("Deregister from NRF successfully")
	}
	logger.InitLog.Infof("PCF terminated")
}

func (a *PcfApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.MainLog.Infof("PCF App is terminated")
}
