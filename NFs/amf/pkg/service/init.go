package service

import (
	"context"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/ngap"
	ngap_message "github.com/free5gc/amf/internal/ngap/message"
	ngap_service "github.com/free5gc/amf/internal/ngap/service"
	"github.com/free5gc/amf/internal/sbi"
	"github.com/free5gc/amf/internal/sbi/consumer"
	"github.com/free5gc/amf/internal/sbi/processor"
	callback "github.com/free5gc/amf/internal/sbi/processor/notifier"
	"github.com/free5gc/amf/pkg/app"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/openapi/models"
)

type AmfAppInterface interface {
	app.App
	consumer.ConsumerAmf
	Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

var AMF AmfAppInterface

type AmfApp struct {
	AmfAppInterface

	cfg    *factory.Config
	amfCtx *amf_context.AMFContext
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	processor *processor.Processor
	consumer  *consumer.Consumer
	sbiServer *sbi.Server
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*AmfApp, error) {
	amf := &AmfApp{
		cfg: cfg,
	}
	amf.SetLogEnable(cfg.GetLogEnable())
	amf.SetLogLevel(cfg.GetLogLevel())
	amf.SetReportCaller(cfg.GetLogReportCaller())

	consumer, err := consumer.NewConsumer(amf)
	if err != nil {
		return amf, err
	}
	amf.consumer = consumer

	processor, err_p := processor.NewProcessor(amf)
	if err_p != nil {
		return amf, err_p
	}
	amf.processor = processor

	amf.ctx, amf.cancel = context.WithCancel(ctx)
	amf.amfCtx = amf_context.GetSelf()

	if amf.sbiServer, err = sbi.NewServer(amf, tlsKeyLogPath); err != nil {
		return nil, err
	}

	AMF = amf

	return amf, nil
}

func (a *AmfApp) SetLogEnable(enable bool) {
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

func (a *AmfApp) SetLogLevel(level string) {
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

func (a *AmfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *AmfApp) Start() {
	self := a.Context()
	amf_context.InitAmfContext(self)

	ngapHandler := ngap_service.NGAPHandler{
		HandleMessage:         ngap.Dispatch,
		HandleNotification:    ngap.HandleSCTPNotification,
		HandleConnectionError: ngap.HandleSCTPConnError,
	}

	sctpConfig := ngap_service.NewSctpConfig(factory.AmfConfig.GetSctpConfig())
	ngap_service.Run(a.Context().NgapIpList, a.Context().NgapPort, ngapHandler, sctpConfig)
	logger.InitLog.Infoln("Server started")

	a.wg.Add(1)
	go a.listenShutdownEvent()

	var profile models.NrfNfManagementNfProfile
	if profileTmp, err1 := a.Consumer().BuildNFInstance(a.Context()); err1 != nil {
		logger.InitLog.Error("Build AMF Profile Error")
	} else {
		profile = profileTmp
	}
	_, nfId, err_reg := a.Consumer().SendRegisterNFInstance(a.ctx, a.Context().NrfUri, a.Context().NfId, &profile)
	if err_reg != nil {
		logger.InitLog.Warnf("Send Register NF Instance failed: %+v", err_reg)
	} else {
		a.Context().NfId = nfId
	}

	if err := a.sbiServer.Run(context.Background(), &a.wg); err != nil {
		logger.MainLog.Fatalf("Run SBI server failed: %+v", err)
	}
	a.WaitRoutineStopped()
}

// Used in AMF planned removal procedure
func (a *AmfApp) Terminate() {
	a.cancel()
}

func (a *AmfApp) Config() *factory.Config {
	return a.cfg
}

func (a *AmfApp) Context() *amf_context.AMFContext {
	return a.amfCtx
}

func (a *AmfApp) CancelContext() context.Context {
	return a.ctx
}

func (a *AmfApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *AmfApp) Processor() *processor.Processor {
	return a.processor
}

func (a *AmfApp) listenShutdownEvent() {
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

func (a *AmfApp) CallServerStop() {
	if a.sbiServer != nil {
		a.sbiServer.Stop()
	}
}

func (a *AmfApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.MainLog.Infof("AMF App is terminated")
}

func (a *AmfApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating AMF...")
	a.CallServerStop()
	// deregister with NRF
	problemDetails, err_deg := a.Consumer().SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.MainLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err_deg != nil {
		logger.MainLog.Errorf("Deregister NF instance Error[%+v]", err_deg)
	} else {
		logger.MainLog.Infof("[AMF] Deregister from NRF successfully")
	}

	// TODO: forward registered UE contexts to target AMF in the same AMF set if there is one

	// ngap
	// send AMF status indication to ran to notify ran that this AMF will be unavailable
	logger.MainLog.Infof("Send AMF Status Indication to Notify RANs due to AMF terminating")
	amfSelf := a.Context()
	unavailableGuamiList := ngap_message.BuildUnavailableGUAMIList(amfSelf.ServedGuamiList)
	amfSelf.AmfRanPool.Range(func(key, value interface{}) bool {
		ran := value.(*amf_context.AmfRan)
		ngap_message.SendAMFStatusIndication(ran, unavailableGuamiList)
		return true
	})
	ngap_service.Stop()
	callback.SendAmfStatusChangeNotify((string)(models.StatusChange_UNAVAILABLE), amfSelf.ServedGuamiList)
}
