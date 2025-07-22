package app

import (
	"context"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	nef_context "github.com/free5gc/nef/internal/context"
	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/internal/sbi"
	"github.com/free5gc/nef/internal/sbi/consumer"
	"github.com/free5gc/nef/internal/sbi/notifier"
	"github.com/free5gc/nef/internal/sbi/processor"
	"github.com/free5gc/nef/pkg/factory"
	"github.com/sirupsen/logrus"
)

type NefApp struct {
	ctx       context.Context
	wg        sync.WaitGroup
	cfg       *factory.Config
	nefCtx    *nef_context.NefContext
	consumer  *consumer.Consumer
	notifier  *notifier.Notifier
	proc      *processor.Processor
	sbiServer *sbi.Server
}

func NewApp(cfg *factory.Config, tlsKeyLogPath string) (*NefApp, error) {
	var err error
	nef := &NefApp{cfg: cfg}
	nef.SetLogEnable(cfg.GetLogEnable())
	nef.SetLogLevel(cfg.GetLogLevel())
	nef.SetReportCaller(cfg.GetLogReportCaller())

	if nef.nefCtx, err = nef_context.NewContext(nef); err != nil {
		return nil, err
	}
	if nef.consumer, err = consumer.NewConsumer(nef); err != nil {
		return nil, err
	}
	if nef.notifier, err = notifier.NewNotifier(); err != nil {
		return nil, err
	}
	if nef.proc, err = processor.NewProcessor(nef); err != nil {
		return nil, err
	}
	if nef.sbiServer, err = sbi.NewServer(nef, tlsKeyLogPath); err != nil {
		return nil, err
	}
	return nef, nil
}

func (a *NefApp) Config() *factory.Config {
	return a.cfg
}

func (a *NefApp) Context() *nef_context.NefContext {
	return a.nefCtx
}

func (a *NefApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *NefApp) Notifier() *notifier.Notifier {
	return a.notifier
}

func (a *NefApp) Processor() *processor.Processor {
	return a.proc
}

func (a *NefApp) SbiServer() *sbi.Server {
	return a.sbiServer
}

func (a *NefApp) SetLogEnable(enable bool) {
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

func (a *NefApp) SetLogLevel(level string) {
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

func (a *NefApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *NefApp) Run() error {
	var cancel context.CancelFunc
	a.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	a.wg.Add(1)
	/* Go Routine is spawned here for listening for cancellation event on
	 * context */
	go a.listenShutdownEvent()

	if err := a.sbiServer.Run(&a.wg); err != nil {
		return err
	}

	if err := a.consumer.RegisterNFInstance(); err != nil {
		return err
	}

	// Wait for interrupt signal to gracefully shutdown UPF
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Receive the interrupt signal
	logger.MainLog.Infof("Shutdown NEF ...")
	// Notify each goroutine and wait them stopped
	cancel()
	a.WaitRoutineStopped()
	logger.MainLog.Infof("NEF exited")
	return nil
}

func (a *NefApp) listenShutdownEvent() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}

		a.wg.Done()
	}()

	<-a.ctx.Done()
	a.sbiServer.Stop()
}

func (a *NefApp) WaitRoutineStopped() {
	a.wg.Wait()
	a.Terminate()
}

func (a *NefApp) Start() {
	if err := a.Run(); err != nil {
		logger.MainLog.Errorf("NEF Run err: %v", err)
	}
}

func (a *NefApp) Terminate() {
	logger.MainLog.Infof("Terminating NEF...")

	// deregister with NRF
	if err := a.consumer.DeregisterNFInstance(); err != nil {
		logger.MainLog.Error(err)
	} else {
		logger.MainLog.Infof("Deregister from NRF successfully")
	}
	logger.MainLog.Infof("NEF terminated")
}
