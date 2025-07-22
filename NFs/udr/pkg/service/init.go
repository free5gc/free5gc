package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/nrf/NFManagement"
	udr_context "github.com/free5gc/udr/internal/context"
	"github.com/free5gc/udr/internal/logger"
	"github.com/free5gc/udr/internal/sbi"
	"github.com/free5gc/udr/internal/sbi/consumer"
	"github.com/free5gc/udr/internal/sbi/processor"
	"github.com/free5gc/udr/pkg/app"
	"github.com/free5gc/udr/pkg/factory"
	"github.com/free5gc/util/mongoapi"
)

type UdrApp struct {
	cfg    *factory.Config
	udrCtx *udr_context.UDRContext

	ctx    context.Context
	cancel context.CancelFunc

	wg        sync.WaitGroup
	sbiServer *sbi.Server
	processor *processor.Processor
	consumer  *consumer.Consumer
}

var _ app.App = &UdrApp{}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*UdrApp, error) {
	udr_context.Init()

	udr := &UdrApp{
		cfg:    cfg,
		udrCtx: udr_context.GetSelf(),
		wg:     sync.WaitGroup{},
	}
	udr.ctx, udr.cancel = context.WithCancel(ctx)

	udr.SetLogEnable(cfg.GetLogEnable())
	udr.SetLogLevel(cfg.GetLogLevel())
	udr.SetReportCaller(cfg.GetLogReportCaller())

	processor := processor.NewProcessor(udr)
	udr.processor = processor

	consumer := consumer.NewConsumer(udr)
	udr.consumer = consumer

	udr.sbiServer = sbi.NewServer(udr, tlsKeyLogPath)

	return udr, nil
}

func (a *UdrApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *UdrApp) Processor() *processor.Processor {
	return a.processor
}

func (a *UdrApp) Config() *factory.Config {
	return a.cfg
}

func (a *UdrApp) Context() *udr_context.UDRContext {
	return a.udrCtx
}

func (a *UdrApp) SetLogEnable(enable bool) {
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

func (a *UdrApp) SetLogLevel(level string) {
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

func (a *UdrApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}
	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (u *UdrApp) registerToNrf(ctx context.Context) error {
	udrContext := u.udrCtx

	nrfUri, nfId, err := u.consumer.SendRegisterNFInstance(ctx, udrContext.NrfUri)
	if err != nil {
		return fmt.Errorf("send register NFInstance error[%s]", err.Error())
	}
	udrContext.NrfUri = nrfUri
	udrContext.NfId = nfId

	return nil
}

func (a *UdrApp) deregisterFromNrf() {
	err := a.consumer.SendDeregisterNFInstance()
	if err != nil {
		switch apiErr := err.(type) {
		case openapi.GenericOpenAPIError:
			switch errModel := apiErr.Model().(type) {
			case NFManagement.DeregisterNFInstanceError:
				pd := &errModel.ProblemDetails
				logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", pd)
			case error:
				logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
			}
		case error:
			logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
		}
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	}

	logger.InitLog.Infof("Deregister from NRF successfully")
}

func (a *UdrApp) Start() {
	err := a.registerToNrf(a.ctx)
	if err != nil {
		logger.InitLog.Errorf("register to NRF failed: %v", err)
	} else {
		logger.InitLog.Infof("register to NRF successfully")
	}

	// get config file info
	logger.InitLog.Infoln("Server started")
	config := factory.UdrConfig
	mongodb := config.Configuration.Mongodb

	logger.InitLog.Infof("UDR Config Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	// Connect to MongoDB
	if err := mongoapi.SetMongoDB(mongodb.Name, mongodb.Url); err != nil {
		logger.InitLog.Errorf("UDR start set MongoDB error: %+v", err)
		return
	}

	// Graceful deregister when panic
	defer func() {
		if p := recover(); p != nil {
			logger.InitLog.Errorf("panic: %v\n%s", p, string(debug.Stack()))
			a.deregisterFromNrf()
		}
	}()

	a.wg.Add(1)
	go a.listenShutdown(a.ctx)

	a.sbiServer.Run(&a.wg)
	a.WaitRoutineStopped()
}

func (a *UdrApp) listenShutdown(ctx context.Context) {
	defer a.wg.Done()

	<-ctx.Done()
	a.terminateProcedure()
}

func (a *UdrApp) Terminate() {
	a.cancel()
}

func (a *UdrApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating UDR...")
	a.CallServerStop()
	a.deregisterFromNrf()
}

func (a *UdrApp) CallServerStop() {
	if a.sbiServer != nil {
		a.sbiServer.Shutdown()
	}
}

func (a *UdrApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.MainLog.Infof("UDR terminated")
}
