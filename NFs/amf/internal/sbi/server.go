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

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/sbi/consumer"
	"github.com/free5gc/amf/internal/sbi/processor"
	util_oauth "github.com/free5gc/amf/internal/util"
	"github.com/free5gc/amf/pkg/app"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

var (
	reqbody         = "[Request Body] "
	applicationjson = "application/json"
	multipartrelate = "multipart/related"
)

type ServerAmf interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

type Server struct {
	ServerAmf

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(amf ServerAmf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerAmf: amf,
	}

	s.router = newRouter(s)

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

func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	amfHttpCallBackGroup := router.Group(factory.AmfCallbackResUriPrefix)
	amfHttpCallBackRoutes := s.getHttpCallBackRoutes()
	applyRoutes(amfHttpCallBackGroup, amfHttpCallBackRoutes)

	for _, serverName := range factory.AmfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serverName) {
		case models.ServiceName_NAMF_COMM:
			amfCommunicationGroup := router.Group(factory.AmfCommResUriPrefix)
			amfCommunicationRoutes := s.getCommunicationRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_COMM)
			amfCommunicationGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfCommunicationGroup, amfCommunicationRoutes)
		case models.ServiceName_NAMF_EVTS:
			amfEventExposureGroup := router.Group(factory.AmfEvtsResUriPrefix)
			amfEventExposureRoutes := s.getEventexposureRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_EVTS)
			amfEventExposureGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfEventExposureGroup, amfEventExposureRoutes)
		case models.ServiceName_NAMF_MT:
			amfMTGroup := router.Group(factory.AmfMtResUriPrefix)
			amfMTRoutes := s.getMTRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_MT)
			amfMTGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfMTGroup, amfMTRoutes)
		case models.ServiceName_NAMF_LOC:
			amfLocationGroup := router.Group(factory.AmfLocResUriPrefix)
			amfLocationRoutes := s.getLocationRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_LOC)
			amfLocationGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfLocationGroup, amfLocationRoutes)
		case models.ServiceName_NAMF_OAM:
			amfOAMGroup := router.Group(factory.AmfOamResUriPrefix)
			amfOAMRoutes := s.getOAMRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_OAM)
			amfOAMGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfOAMGroup, amfOAMRoutes)
		case models.ServiceName_NAMF_MBS_COMM:
			amfMbsComGroup := router.Group(factory.AmfMbsComResUriPrefix)
			amfMbsComRoutes := s.getMbsCommunicationRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_MBS_COMM)
			amfMbsComGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfMbsComGroup, amfMbsComRoutes)
		case models.ServiceName_NAMF_MBS_BC:
			amfMbsBCGroup := router.Group(factory.AmfMbsBCResUriPrefix)
			amfMbsBCRoutes := s.getMbsBroadcastRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_MBS_BC)
			amfMbsBCGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfMbsBCGroup, amfMbsBCRoutes)
		}
	}

	return router
}

func (s *Server) Run(traceCtx context.Context, wg *sync.WaitGroup) error {
	wg.Add(1)
	go s.startServer(wg)

	return nil
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
	logger.SBILog.Warnf("SBI server (listen on %s) stopped", s.httpServer.Addr)
}
