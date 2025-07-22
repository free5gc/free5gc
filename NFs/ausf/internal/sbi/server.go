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

	ausf_context "github.com/free5gc/ausf/internal/context"
	"github.com/free5gc/ausf/internal/logger"
	"github.com/free5gc/ausf/internal/sbi/consumer"
	"github.com/free5gc/ausf/internal/sbi/processor"
	"github.com/free5gc/ausf/internal/util"
	"github.com/free5gc/ausf/pkg/app"
	"github.com/free5gc/ausf/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

type ServerAusf interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

type Server struct {
	ServerAusf

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(ausf ServerAusf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerAusf: ausf,
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

	for _, serviceName := range factory.AusfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serviceName) {
		case models.ServiceName_NAUSF_AUTH:
			ausfUeAuthenticationGroup := router.Group(factory.AusfAuthResUriPrefix)
			ausfUeAuthenticationRoutes := s.getUeAuthenticationRoutes()
			routerAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NAUSF_AUTH)
			ausfUeAuthenticationGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, ausf_context.GetSelf())
			})
			applyRoutes(ausfUeAuthenticationGroup, ausfUeAuthenticationRoutes)
		case models.ServiceName_NAUSF_SORPROTECTION:
			ausfSorprotectionGroup := router.Group(factory.AusfSorprotectionResUriPrefix)
			ausfSorprotectionRoutes := s.getSorprotectionRoutes()
			routerAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NAUSF_SORPROTECTION)
			ausfSorprotectionGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, ausf_context.GetSelf())
			})
			applyRoutes(ausfSorprotectionGroup, ausfSorprotectionRoutes)
		case models.ServiceName_NAUSF_UPUPROTECTION:
			ausfUpuprotectionGroup := router.Group(factory.AusfUpuprotectionResUriPrefix)
			ausfUpuprotectionRoutes := s.getUpuprotectionRoutes()
			routerAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NAUSF_UPUPROTECTION)
			ausfUpuprotectionGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, ausf_context.GetSelf())
			})
			applyRoutes(ausfUpuprotectionGroup, ausfUpuprotectionRoutes)

		default:
			logger.SBILog.Warnf("Unsupported service name: %s", serviceName)
		}
	}

	return router
}

func (s *Server) Run(traceCtx context.Context, wg *sync.WaitGroup) error {
	var err error
	_, s.Context().NfId, err = s.Consumer().RegisterNFInstance(context.Background())
	if err != nil {
		logger.InitLog.Errorf("AUSF register to NRF Error[%s]", err.Error())
	}

	wg.Add(1)
	go s.startServer(wg)

	return nil
}

func (s *Server) Shutdown() {
	s.shutdownHttpServer()
}

func (s *Server) shutdownHttpServer() {
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
