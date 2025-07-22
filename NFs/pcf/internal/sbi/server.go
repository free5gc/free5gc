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
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/sbi/consumer"
	"github.com/free5gc/pcf/internal/sbi/processor"
	"github.com/free5gc/pcf/internal/util"
	"github.com/free5gc/pcf/pkg/app"
	"github.com/free5gc/pcf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	APIFunc gin.HandlerFunc
}

func applyRoutes(group *gin.RouterGroup, routes []Route) {
	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.APIFunc)
		case "POST":
			group.POST(route.Pattern, route.APIFunc)
		case "PUT":
			group.PUT(route.Pattern, route.APIFunc)
		case "PATCH":
			group.PATCH(route.Pattern, route.APIFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.APIFunc)
		}
	}
}

type pcf interface {
	app.App
	Processor() *processor.Processor
	Consumer() *consumer.Consumer
	CancelContext() context.Context
}

type Server struct {
	pcf

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(pcf pcf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		pcf:    pcf,
		router: logger_util.NewGinWithLogrus(logger.GinLog),
	}

	smPolicyRoutes := s.getSmPolicyRoutes()
	smPolicyGroup := s.router.Group(factory.PcfSMpolicyCtlResUriPrefix)
	applyRoutes(smPolicyGroup, smPolicyRoutes)

	amPolicyRoutes := s.getAmPolicyRoutes()
	amPolicyGroup := s.router.Group(factory.PcfAMpolicyCtlResUriPrefix)
	amRouterAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NPCF_AM_POLICY_CONTROL)
	amPolicyGroup.Use(func(c *gin.Context) {
		amRouterAuthorizationCheck.Check(c, s.Context())
	})
	applyRoutes(amPolicyGroup, amPolicyRoutes)

	bdtPolicyRoutes := s.getBdtPolicyRoutes()
	bdtPolicyGroup := s.router.Group(factory.PcfBdtPolicyCtlResUriPrefix)
	bdtRouterAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NPCF_BDTPOLICYCONTROL)
	bdtPolicyGroup.Use(func(c *gin.Context) {
		bdtRouterAuthorizationCheck.Check(c, s.Context())
	})
	applyRoutes(bdtPolicyGroup, bdtPolicyRoutes)

	httpcallbackRoutes := s.getHttpCallBackRoutes()
	httpcallbackGroup := s.router.Group(factory.PcfCallbackResUriPrefix)
	applyRoutes(httpcallbackGroup, httpcallbackRoutes)

	oamRoutes := s.getOamRoutes()
	oamGroup := s.router.Group(factory.PcfOamResUriPrefix)
	oamRouterAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NPCF_OAM)
	oamGroup.Use(func(c *gin.Context) {
		oamRouterAuthorizationCheck.Check(c, s.Context())
	})
	applyRoutes(oamGroup, oamRoutes)

	policyAuthorizationRoutes := s.getPolicyAuthorizationRoutes()
	policyAuthorizationGroup := s.router.Group(factory.PcfPolicyAuthResUriPrefix)
	policyAuthorizationRouterAuthorizationCheck := util.
		NewRouterAuthorizationCheck(models.ServiceName_NPCF_POLICYAUTHORIZATION)
	policyAuthorizationGroup.Use(func(c *gin.Context) {
		policyAuthorizationRouterAuthorizationCheck.Check(c, s.Context())
	})
	applyRoutes(policyAuthorizationGroup, policyAuthorizationRoutes)

	uePolicyRoutes := s.getUePolicyRoutes()
	uePolicyGroup := s.router.Group(factory.PcfUePolicyCtlResUriPrefix)
	applyRoutes(uePolicyGroup, uePolicyRoutes)

	cfg := s.Config()
	bindAddr := cfg.GetSbiBindingAddr()
	logger.SBILog.Infof("Binding addr: [%s]", bindAddr)
	var err error
	if s.httpServer, err = httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, s.router); err != nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}
	s.httpServer.ErrorLog = log.New(logger.SBILog.WriterLevel(logrus.ErrorLevel), "HTTP2: ", 0)

	return s, nil
}

func (s *Server) Run(traceCtx context.Context, wg *sync.WaitGroup) error {
	var err error
	_, s.Context().NfId, err = s.Consumer().SendRegisterNFInstance(s.CancelContext())
	if err != nil {
		logger.InitLog.Errorf("PCF register to NRF Error[%s]", err.Error())
	}

	wg.Add(1)
	go s.startServer(wg)

	return nil
}

func (s *Server) Shutdown(traceCtx context.Context) {
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

	var err error
	cfg := s.Config()
	scheme := cfg.GetSbiScheme()
	switch scheme {
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
