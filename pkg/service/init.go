package service

import (
	"context"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	nrf_context "github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/internal/sbi/accesstoken"
	"github.com/free5gc/nrf/internal/sbi/discovery"
	"github.com/free5gc/nrf/internal/sbi/management"
	"github.com/free5gc/nrf/pkg/app"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/mongoapi"
)

var NRF *NrfApp

var _ app.App = &NrfApp{}

type NrfApp struct {
	app.App

	cfg    *factory.Config
	nrfCtx *nrf_context.NRFContext

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	server *http.Server
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*NrfApp, error) {
	nrf := &NrfApp{
		cfg: cfg,
		wg:  sync.WaitGroup{},
	}
	nrf.SetLogEnable(cfg.GetLogEnable())
	nrf.SetLogLevel(cfg.GetLogLevel())
	nrf.SetReportCaller(cfg.GetLogReportCaller())

	err := nrf_context.InitNrfContext()
	if err != nil {
		logger.InitLog.Errorln(err)
		return nrf, err
	}

	nrf.nrfCtx = nrf_context.GetSelf()
	nrf.ctx, nrf.cancel = context.WithCancel(ctx)

	return nrf, nil
}

func (a *NrfApp) SetLogEnable(enable bool) {
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

func (a *NrfApp) SetLogLevel(level string) {
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

func (a *NrfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *NrfApp) Start() {
	if err := mongoapi.SetMongoDB(factory.NrfConfig.Configuration.MongoDBName,
		factory.NrfConfig.Configuration.MongoDBUrl); err != nil {
		logger.InitLog.Errorf("SetMongoDB failed: %+v", err)
		return
	}
	logger.InitLog.Infoln("Server starting")

	router := logger_util.NewGinWithLogrus(logger.GinLog)

	accesstoken.AddService(router)
	discovery.AddService(router)
	management.AddService(router)

	tlsKeyLogPath := ""
	bindAddr := factory.NrfConfig.GetSbiBindingAddr()
	logger.InitLog.Infof("Binding addr: [%s]", bindAddr)
	server, err := httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, router)
	if err != nil {
		logger.InitLog.Warnf("Initialize HTTP server: +%v", err)
		return
	}
	a.server = server

	a.wg.Add(1)
	go a.listenShutdownEvent()

	serverScheme := factory.NrfConfig.GetSbiScheme()
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		// TODO: support TLS mutual authentication for OAuth
		err = server.ListenAndServeTLS(
			factory.NrfConfig.GetNrfCertPemPath(),
			factory.NrfConfig.GetNrfPrivKeyPath())
	}

	if err != nil && err != http.ErrServerClosed {
		logger.MainLog.Errorf("SBI server error: %v", err)
	}
	logger.MainLog.Warnf("SBI server (listen on %s) stopped", server.Addr)
}

func (a *NrfApp) listenShutdownEvent() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		a.wg.Done()
	}()

	<-a.ctx.Done()
	a.Terminate()
}

func (a *NrfApp) Terminate() {
	logger.InitLog.Infof("Terminating NRF...")

	logger.InitLog.Infof("Waiting for 2s for other NFs to deregister")
	time.Sleep(2 * time.Second)

	a.cancel()

	logger.InitLog.Infof("Remove NF Profile...")
	err := mongoapi.Drop("NfProfile")
	if err != nil {
		logger.InitLog.Errorf("Drop NfProfile collection failed: %+v", err)
	}

	// server stop
	const defaultShutdownTimeout time.Duration = 2 * time.Second

	toCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()
	if err := a.server.Shutdown(toCtx); err != nil {
		logger.MainLog.Errorf("Could not close SBI server: %#v", err)
	}
}

func (a *NrfApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.InitLog.Infof("NRF App terminated")
}
