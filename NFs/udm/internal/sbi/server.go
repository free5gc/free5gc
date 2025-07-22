package sbi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udm/internal/logger"
	"github.com/free5gc/udm/internal/sbi/consumer"
	"github.com/free5gc/udm/internal/sbi/processor"
	"github.com/free5gc/udm/internal/util"
	"github.com/free5gc/udm/pkg/app"
	"github.com/free5gc/udm/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

type ServerUdm interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor
	CancelContext() context.Context
}

type Server struct {
	ServerUdm

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(udm ServerUdm, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerUdm: udm,
		router:    logger_util.NewGinWithLogrus(logger.GinLog),
	}

	cfg := s.Config()
	bindAddr := cfg.GetSbiBindingAddr()
	logger.SBILog.Infof("Binding addr: [%s]", bindAddr)
	var err error
	if s.httpServer, err = httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, s.router); err != nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}
	s.httpServer.ErrorLog = log.New(logger.SBILog.WriterLevel(logrus.ErrorLevel), "HTTP2: ", 0)

	return s, err
}

func (s *Server) Run(traceCtx context.Context, wg *sync.WaitGroup) error {
	logger.SBILog.Info("Starting server...")

	var err error
	_, s.Context().NfId, err = s.Consumer().RegisterNFInstance(s.CancelContext())
	if err != nil {
		logger.InitLog.Errorf("UDM register to NRF Error[%s]", err.Error())
	}

	wg.Add(1)
	go s.startServer(wg)

	return nil
}

func (s *Server) startServer(wg *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.SBILog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			s.Terminate()
		}
		wg.Done()
	}()

	logger.SBILog.Infof("Start SBI server (listen on %s)", s.httpServer.Addr)

	s.router = newRouter(s)

	var err error
	cfg := s.Config()
	scheme := cfg.GetSbiScheme()
	switch s.Config().GetSbiScheme() {
	case "http":
		err = s.httpServer.ListenAndServe()
	case "https":
		err = s.httpServer.ListenAndServeTLS(
			cfg.GetCertPemPath(),
			cfg.GetCertKeyPath())
	default:
		err = fmt.Errorf("no support this scheme[%s]", scheme)
	}

	if err != nil && err != http.ErrServerClosed {
		logger.SBILog.Errorf("SBI server error: %v", err)
	}
	logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
}

func (s *Server) Shutdown() {
	s.shutdownHttpServer()
}

func (s *Server) Stop() {
	const defaultShutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer != nil {
		logger.SBILog.Infof("Stop SBI server (listen on %s)", s.httpServer.Addr)
		toCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(toCtx); err != nil {
			logger.SBILog.Errorf("Could not close SBI server: %#v", err)
		}
	}
}

func (s *Server) shutdownHttpServer() {
	const shutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer == nil {
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := s.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.SBILog.Errorf("HTTP server shutdown failed: %+v", err)
	}
}

func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	// EE
	udmEERoutes := s.getEventExposureRoutes()
	udmEEGroup := s.router.Group(factory.UdmEeResUriPrefix)
	routerAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_EE)
	udmEEGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmEEGroup, udmEERoutes)

	// Callback
	udmCallBackRoutes := s.getHttpCallBackRoutes()
	udmCallNackGroup := s.router.Group("")
	AddService(udmCallNackGroup, udmCallBackRoutes)

	// UEAU
	udmUEAURoutes := s.getUEAuthenticationRoutes()
	udmUEAUGroup := s.router.Group(factory.UdmUeauResUriPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_UEAU)
	udmUEAUGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmUEAUGroup, udmUEAURoutes)

	ueauTwoLayerPath := "/:supi/:twoLayer"
	udmUEAUGroup.Any(ueauTwoLayerPath, s.UEAUTwoLayerPathHandlerFunc)

	ueauThreeLayerPath := "/:supi/:twoLayer/:thirdLayer"
	udmUEAUGroup.Any(ueauThreeLayerPath, s.UEAUThreeLayerPathHandlerFunc)

	generateAvPath := "/:supi/hss-security-information/:hssAuthType/generate-av"
	udmUEAUGroup.Any(generateAvPath, s.HandleGenerateAv)

	// UECM
	udmUECMRoutes := s.getUEContextManagementRoutes()
	udmUECMGroup := s.router.Group(factory.UdmUecmResUriPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_UECM)
	udmUECMGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmUECMGroup, udmUECMRoutes)

	// SDM
	udmSDMRoutes := s.getSubscriberDataManagementRoutes()
	udmSDMGroup := s.router.Group(factory.UdmSdmResUriPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_SDM)
	udmSDMGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmSDMGroup, udmSDMRoutes)

	oneLayerPath := "/:supi"
	udmSDMGroup.Any(oneLayerPath, s.OneLayerPathHandlerFunc)

	twoLayerPath := "/:supi/:subscriptionId"
	udmSDMGroup.Any(twoLayerPath, s.TwoLayerPathHandlerFunc)

	threeLayerPath := "/:supi/:subscriptionId/:thirdLayer"
	udmSDMGroup.Any(threeLayerPath, s.ThreeLayerPathHandlerFunc)

	// PP
	udmPPRoutes := s.getParameterProvisionRoutes()
	udmPPGroup := s.router.Group(factory.UdmPpResUriPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_PP)
	udmPPGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmPPGroup, udmPPRoutes)

	// MT
	udmMTRoutes := s.getMTRoutes()
	udmMTGroup := s.router.Group(factory.UdmMtResUrdPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_MT)
	udmMTGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmMTGroup, udmMTRoutes)

	// NIDDAU
	udmNIDDAURoutes := s.getNIDDAuthenticationRoutes()
	udmNIDDAUGroup := s.router.Group(factory.UdmNiddauResUriPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_NIDDAU)
	udmNIDDAUGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmNIDDAUGroup, udmNIDDAURoutes)

	// RSDS
	udmRSDSRoutes := s.getReportSMDeliveryStatusRoutes()
	udmRSDSGroup := s.router.Group(factory.UdmRsdsResUriPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_RSDS)
	udmRSDSGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmRSDSGroup, udmRSDSRoutes)

	// SSAU
	udmSSAURoutes := s.getServiceSpecificAuthorizationRoutes()
	udmSSAUGroup := s.router.Group(factory.UdmSsauResUriPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_SSAU)
	udmSSAUGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmSSAUGroup, udmSSAURoutes)

	// UEID
	udmUEIDRoutes := s.getUEIDRoutes()
	udmUEIDGroup := s.router.Group(factory.UdmUeidResUriPrefix)
	routerAuthorizationCheck = util.NewRouterAuthorizationCheck(models.ServiceName_NUDM_UEID)
	udmUEIDGroup.Use(func(c *gin.Context) {
		routerAuthorizationCheck.Check(c, s.Context())
	})
	AddService(udmUEIDGroup, udmUEIDRoutes)

	return router
}
