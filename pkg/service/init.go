package service

import (
	"io/ioutil"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	nrf_context "github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/internal/sbi/accesstoken"
	"github.com/free5gc/nrf/internal/sbi/discovery"
	"github.com/free5gc/nrf/internal/sbi/management"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/mongoapi"
)

type NrfApp struct {
	cfg    *factory.Config
	nrfCtx *nrf_context.NRFContext
}

func NewApp(cfg *factory.Config) (*NrfApp, error) {
	nrf := &NrfApp{cfg: cfg}
	nrf.SetLogEnable(cfg.GetLogEnable())
	nrf.SetLogLevel(cfg.GetLogLevel())
	nrf.SetReportCaller(cfg.GetLogReportCaller())

	err := nrf_context.InitNrfContext()
	if err != nil {
		logger.InitLog.Errorln(err)
		return nrf, err
	}
	nrf.nrfCtx = nrf_context.GetSelf()
	return nrf, nil
}

func (a *NrfApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == ioutil.Discard {
		return
	}

	a.cfg.SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(ioutil.Discard)
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

func (a *NrfApp) Start(tlsKeyLogPath string) {
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

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		<-signalChannel
		// Waiting for other NFs to deregister
		time.Sleep(2 * time.Second)
		a.Terminate()
		os.Exit(0)
	}()

	bindAddr := factory.NrfConfig.GetSbiBindingAddr()
	logger.InitLog.Infof("Binding addr: [%s]", bindAddr)
	server, err := httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, router)
	if err != nil {
		logger.InitLog.Warnf("Initialize HTTP server: +%v", err)
		return
	}

	serverScheme := factory.NrfConfig.GetSbiScheme()
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		// TODO: support TLS mutual authentication for OAuth
		err = server.ListenAndServeTLS(
			factory.NrfConfig.GetNrfCertPemPath(),
			factory.NrfConfig.GetNrfPrivKeyPath())
	}

	if err != nil {
		logger.InitLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (a *NrfApp) Terminate() {
	logger.InitLog.Infof("Terminating NRF...")

	logger.InitLog.Infof("Remove NF Profile...")
	err := mongoapi.Drop("NfProfile")
	if err != nil {
		logger.InitLog.Errorf("Drop NfProfile collection failed: %+v", err)
	}

	logger.InitLog.Infof("NRF terminated")
}
