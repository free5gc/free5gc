package app

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/internal/pfcp"
	"github.com/free5gc/go-upf/pkg/factory"
)

type UpfApp struct {
	ctx        context.Context
	wg         sync.WaitGroup
	cfg        *factory.Config
	driver     forwarder.Driver
	pfcpServer *pfcp.PfcpServer
}

func NewApp(cfg *factory.Config) (*UpfApp, error) {
	upf := &UpfApp{
		cfg: cfg,
	}
	upf.SetLogLevel(cfg.Logger.Level)
	upf.SetReportCaller(cfg.Logger.ReportCaller)
	return upf, nil
}

func (u *UpfApp) Config() *factory.Config {
	return u.cfg
}

func (a *UpfApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}

	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}

	logger.Log.SetLevel(lvl)
}

func (a *UpfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	logger.Log.SetReportCaller(reportCaller)
}

func (a *UpfApp) SetLogEnable(enable bool) {

}

func (u *UpfApp) Run() error {
	var cancel context.CancelFunc
	u.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	u.wg.Add(1)
	/* Go Routine is spawned here for listening for cancellation event on
	 * context */
	go u.listenShutdownEvent()

	var err error
	u.driver, err = forwarder.NewDriver(&u.wg, u.cfg)
	if err != nil {
		return err
	}

	u.pfcpServer = pfcp.NewPfcpServer(u.cfg, u.driver)
	u.driver.HandleReport(u.pfcpServer)
	u.pfcpServer.Start(&u.wg)

	logger.MainLog.Infoln("UPF started")

	// Wait for interrupt signal to gracefully shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Receive the interrupt signal
	logger.MainLog.Infof("Shutdown UPF ...")
	// Notify each goroutine and wait them stopped
	cancel()
	u.WaitRoutineStopped()
	logger.MainLog.Infof("UPF exited")
	return nil
}

func (u *UpfApp) listenShutdownEvent() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}

		u.wg.Done()
	}()

	<-u.ctx.Done()
	if u.pfcpServer != nil {
		u.pfcpServer.Stop()
	}
	if u.driver != nil {
		u.driver.Close()
	}
}

func (u *UpfApp) WaitRoutineStopped() {
	u.wg.Wait()
	u.Terminate()
}

func (u *UpfApp) Start() {
	if err := u.Run(); err != nil {
		logger.MainLog.Errorf("UPF Run err: %v", err)
	}
}

func (u *UpfApp) Terminate() {
	logger.MainLog.Infof("Terminating UPF...")
	logger.MainLog.Infof("UPF terminated")
}
