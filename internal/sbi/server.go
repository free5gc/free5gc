package sbi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/internal/sbi/processor"
	"github.com/free5gc/nrf/internal/util"
	"github.com/free5gc/nrf/pkg/app"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/metrics"
)

type ServerNrf interface {
	app.App

	// Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

type Server struct {
	ServerNrf

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(nrf ServerNrf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerNrf: nrf,
		router:    logger_util.NewGinWithLogrus(logger.GinLog),
	}
	s.router.Use(metrics.InboundMetrics())
	cfg := s.Config()
	bindAddr := cfg.GetSbiBindingAddr()
	logger.SBILog.Infof("Binding addr: [%s]", bindAddr)

	s.applyService()

	var err error
	if s.httpServer, err = httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, s.router); err != nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}
	return s, nil
}

func (s *Server) GetLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.NfmLog.Error(err)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (s *Server) applyService() {
	accesstokenRoutes := s.getAccesstokenRoutes()
	accesstokenGroup := s.router.Group("") // accesstoken service didn't have api prefix
	applyRoutes(accesstokenGroup, accesstokenRoutes)

	bootstrappingRoutes := s.getBootstrappingRoutes()
	bootstrappingGroup := s.router.Group(factory.NrfBootstrappingPrefix)
	applyRoutes(bootstrappingGroup, bootstrappingRoutes)

	discoveryRoutes := s.getNfDiscoveryRoutes()
	discoveryGroup := s.router.Group(factory.NrfDiscResUriPrefix)
	discAuthCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NNRF_DISC)
	discoveryGroup.Use(func(c *gin.Context) {
		discAuthCheck.Check(c, s.Context())
	})
	applyRoutes(discoveryGroup, discoveryRoutes)

	// OAuth2 must exclude NfRegister
	nfRegisterRoute := s.getNfRegisterRoute()
	nfRegisterGroup := s.router.Group(factory.NrfNfmResUriPrefix)
	applyRoutes(nfRegisterGroup, nfRegisterRoute)

	managementRoutes := s.getNfManagementRoute()
	managementGroup := s.router.Group(factory.NrfNfmResUriPrefix)
	managementAuthCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NNRF_NFM)
	managementGroup.Use(func(c *gin.Context) {
		managementAuthCheck.Check(c, s.Context())
	})
	applyRoutes(managementGroup, managementRoutes)
}

func (s *Server) Run(wg *sync.WaitGroup) error {
	wg.Add(1)
	go s.startServer(wg)

	logger.SBILog.Infoln("SBI server started")
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

	cfg := s.Config()
	serverScheme := cfg.GetSbiScheme()

	var err error
	switch serverScheme {
	case "http":
		err = s.httpServer.ListenAndServe()
	case "https":
		// TODO: support TLS mutual authentication for OAuth
		err = s.httpServer.ListenAndServeTLS(
			cfg.GetNrfCertPemPath(),
			cfg.GetNrfPrivKeyPath())
	default:
		err = fmt.Errorf("not support this scheme[%s]", serverScheme)
	}

	if err != nil && err != http.ErrServerClosed {
		logger.SBILog.Errorf("SBI server error: %v", err)
	}
	logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
}

func (s *Server) Stop() {
	// server stop
	const defaultShutdownTimeout time.Duration = 2 * time.Second

	toCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(toCtx); err != nil {
		logger.SBILog.Errorf("Could not close SBI server: %#v", err)
	}
}
