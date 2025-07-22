package service

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"runtime/debug"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/ike"
	"github.com/free5gc/n3iwf/internal/ike/xfrm"
	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/n3iwf/internal/ngap"
	"github.com/free5gc/n3iwf/internal/nwucp"
	"github.com/free5gc/n3iwf/internal/nwuup"
	"github.com/free5gc/n3iwf/pkg/app"
	"github.com/free5gc/n3iwf/pkg/factory"
)

var N3IWF *N3iwfApp

var _ app.App = &N3iwfApp{}

type N3iwfApp struct {
	n3iwfCtx    *n3iwf_context.N3IWFContext
	cfg         *factory.Config
	ngapServer  *ngap.Server
	nwucpServer *nwucp.Server
	nwuupServer *nwuup.Server
	ikeServer   *ike.Server

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewApp(
	ctx context.Context,
	cfg *factory.Config,
	tlsKeyLogPath string,
) (*N3iwfApp, error) {
	var err error
	n3iwf := &N3iwfApp{
		cfg: cfg,
		wg:  sync.WaitGroup{},
	}
	n3iwf.ctx, n3iwf.cancel = context.WithCancel(ctx)

	n3iwf.SetLogEnable(cfg.GetLogEnable())
	n3iwf.SetLogLevel(cfg.GetLogLevel())
	n3iwf.SetReportCaller(cfg.GetLogReportCaller())

	if n3iwf.n3iwfCtx, err = n3iwf_context.NewContext(n3iwf); err != nil {
		return nil, errors.Wrap(err, "NewApp()")
	}
	if n3iwf.ngapServer, err = ngap.NewServer(n3iwf); err != nil {
		return nil, errors.Wrap(err, "NewApp()")
	}
	if n3iwf.nwucpServer, err = nwucp.NewServer(n3iwf); err != nil {
		return nil, errors.Wrap(err, "NewApp()")
	}
	if n3iwf.nwuupServer, err = nwuup.NewServer(n3iwf); err != nil {
		return nil, errors.Wrap(err, "NewApp()")
	}
	if n3iwf.ikeServer, err = ike.NewServer(n3iwf); err != nil {
		return nil, errors.Wrap(err, "NewApp()")
	}
	N3IWF = n3iwf
	return n3iwf, nil
}

func (a *N3iwfApp) CancelContext() context.Context {
	return a.ctx
}

func (a *N3iwfApp) Context() *n3iwf_context.N3IWFContext {
	return a.n3iwfCtx
}

func (a *N3iwfApp) Config() *factory.Config {
	return a.cfg
}

func (a *N3iwfApp) SetLogEnable(enable bool) {
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

func (a *N3iwfApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	mainLog := logger.MainLog
	if err != nil {
		mainLog.Warnf("Log level [%s] is invalid", level)
		return
	}

	mainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}

	a.cfg.SetLogLevel(level)
	logger.Log.SetLevel(lvl)
}

func (a *N3iwfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *N3iwfApp) Run() error {
	if err := a.initDefaultXfrmInterface(); err != nil {
		return err
	}
	mainLog := logger.MainLog

	a.wg.Add(1)
	go a.listenShutdownEvent()

	// NGAP
	if err := a.ngapServer.Run(&a.wg); err != nil {
		return errors.Wrapf(err, "Run()")
	}
	mainLog.Infof("NGAP service running.")

	// Relay listeners
	// Control plane
	if err := a.nwucpServer.Run(&a.wg); err != nil {
		return errors.Wrapf(err, "Listen NWu control plane traffic failed")
	}
	mainLog.Infof("NAS TCP server successfully started.")

	// User plane of N3IWF
	if err := a.nwuupServer.Run(&a.wg); err != nil {
		return errors.Wrapf(err, "Listen NWu user plane traffic failed")
	}
	mainLog.Infof("Listening NWu user plane traffic")

	// IKE
	if err := a.ikeServer.Run(&a.wg); err != nil {
		return errors.Wrapf(err, "Start IKE service failed")
	}
	mainLog.Infof("IKE service running")

	mainLog.Infof("N3IWF started")

	a.WaitRoutineStopped()
	return nil
}

func (a *N3iwfApp) Start() {
	if err := a.Run(); err != nil {
		logger.MainLog.Errorf("N3IWF Run err: %v", err)
	}
}

func (a *N3iwfApp) listenShutdownEvent() {
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

func (a *N3iwfApp) WaitRoutineStopped() {
	a.wg.Wait()
	// Waiting for negotiatioon with netlink for deleting interfaces
	a.removeIPsecInterfaces()
	logger.MainLog.Infof("N3IWF App is terminated")
}

func (a *N3iwfApp) initDefaultXfrmInterface() error {
	// Setup default IPsec interface for Control Plane
	var linkIPSec netlink.Link
	var err error
	n3iwfCtx := a.n3iwfCtx
	cfg := a.Config()
	mainLog := logger.MainLog
	n3iwfIPAddr := net.ParseIP(cfg.GetIPSecGatewayAddr()).To4()
	n3iwfIPAddrAndSubnet := net.IPNet{IP: n3iwfIPAddr, Mask: n3iwfCtx.IPSecInnerIPPool.IPSubnet.Mask}
	newXfrmiName := fmt.Sprintf("%s-default", cfg.GetXfrmIfaceName())

	if linkIPSec, err = xfrm.SetupIPsecXfrmi(newXfrmiName, n3iwfCtx.XfrmParentIfaceName,
		cfg.GetXfrmIfaceId(), n3iwfIPAddrAndSubnet); err != nil {
		mainLog.Errorf("Setup XFRM interface %s fail: %+v", newXfrmiName, err)
		return err
	}

	route := &netlink.Route{
		LinkIndex: linkIPSec.Attrs().Index,
		Dst:       n3iwfCtx.IPSecInnerIPPool.IPSubnet,
	}

	if err := netlink.RouteAdd(route); err != nil {
		mainLog.Warnf("netlink.RouteAdd: %+v", err)
	}

	mainLog.Infof("Setup XFRM interface %s ", newXfrmiName)

	n3iwfCtx.XfrmIfaces.LoadOrStore(cfg.GetXfrmIfaceId(), linkIPSec)
	n3iwfCtx.XfrmIfaceIdOffsetForUP = 1

	return nil
}

func (a *N3iwfApp) removeIPsecInterfaces() {
	mainLog := logger.MainLog
	a.n3iwfCtx.XfrmIfaces.Range(
		func(key, value interface{}) bool {
			iface := value.(netlink.Link)
			if err := netlink.LinkDel(iface); err != nil {
				mainLog.Errorf("Delete interface %s fail: %+v", iface.Attrs().Name, err)
			} else {
				mainLog.Infof("Delete interface: %s", iface.Attrs().Name)
			}
			return true
		},
	)
}

func (a *N3iwfApp) Terminate() {
	a.cancel()
}

func (a *N3iwfApp) terminateProcedure() {
	logger.MainLog.Info("Stopping service created by N3IWF")

	a.ngapServer.Stop()
	a.nwucpServer.Stop()
	a.nwuupServer.Stop()
	a.ikeServer.Stop()
}

func (a *N3iwfApp) SendNgapEvt(evt n3iwf_context.NgapEvt) {
	a.ngapServer.SendNgapEvt(evt)
}

func (a *N3iwfApp) SendIkeEvt(evt n3iwf_context.IkeEvt) {
	a.ikeServer.SendIkeEvt(evt)
}
