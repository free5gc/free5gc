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

	"github.com/free5gc/chf/internal/logger"
	"github.com/free5gc/chf/internal/sbi/consumer"
	"github.com/free5gc/chf/internal/sbi/processor"
	"github.com/free5gc/chf/internal/util"
	"github.com/free5gc/chf/pkg/app"
	"github.com/free5gc/chf/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

type ServerChf interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor
	CancelContext() context.Context
}

type Server struct {
	ServerChf

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(chf ServerChf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerChf: chf,
		router:    logger_util.NewGinWithLogrus(logger.GinLog),
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

	for _, serviceName := range s.Config().Configuration.ServiceNameList {
		switch models.ServiceName(serviceName) {
		case models.ServiceName_NCHF_CONVERGEDCHARGING:
			chfConvergedChargingGroup := router.Group(factory.ConvergedChargingResUriPrefix)
			chfConvergedChargingGroup.Use(func(c *gin.Context) {
				// oauth middleware
				util.NewRouterAuthorizationCheck(models.ServiceName(serviceName)).Check(c, s.Context())
			})
			chfConvergedChargingRoutes := s.getConvergenChargingRoutes()
			applyRoutes(chfConvergedChargingGroup, chfConvergedChargingRoutes)

		case models.ServiceName_NCHF_OFFLINEONLYCHARGING:
			chfOfflineOnlyChargingGroup := router.Group(factory.OfflineOnlyChargingResUriPrefix)
			chfOfflineOnlyChargingGroup.Use(func(c *gin.Context) {
				// oauth middleware
				util.NewRouterAuthorizationCheck(models.ServiceName(serviceName)).Check(c, s.Context())
			})

			chfOfflineOnlyChargingGroupRoutes := s.getOfflineOnlyChargingRoutes()
			applyRoutes(chfOfflineOnlyChargingGroup, chfOfflineOnlyChargingGroupRoutes)

		case models.ServiceName_NCHF_SPENDINGLIMITCONTROL:
			chfSpendingLimitControlGroup := router.Group(factory.SpendingLimitControlResUriPrefix)
			chfSpendingLimitControlGroup.Use(func(c *gin.Context) {
				// oauth middleware
				util.NewRouterAuthorizationCheck(models.ServiceName(serviceName)).Check(c, s.Context())
			})
			chfSpendingLimitControlRoutes := s.getSpendingLimitControlRoutes()
			applyRoutes(chfSpendingLimitControlGroup, chfSpendingLimitControlRoutes)

		default:
			logger.SBILog.Warnf("Unsupported service name: %s", serviceName)
		}
	}

	return router
}

func (s *Server) Run(traceCtx context.Context, wg *sync.WaitGroup) error {
	var err error
	_, s.Context().NfId, err = s.Consumer().RegisterNFInstance(s.CancelContext())
	if err != nil {
		logger.InitLog.Errorf("CHF register to NRF Error[%s]", err.Error())
	}

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
	logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
}
