package sbi

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/pkg/app"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"

	"github.com/free5gc/nrf/internal/sbi/accesstoken"
	"github.com/free5gc/nrf/internal/sbi/discovery"
	"github.com/free5gc/nrf/internal/sbi/management"
)

type ServerNrf interface {
	app.App

	// Consumer() *consumer.Consumer
	// Processor() *processor.Processor
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
	cfg := s.Config()
	bindAddr := cfg.GetSbiBindingAddr()
	logger.SBILog.Infof("Binding addr: [%s]", bindAddr)

	accesstoken.AddService(s.router)
	discovery.AddService(s.router)
	management.AddService(s.router)

	var err error
	if s.httpServer, err = httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, s.router); err != nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}
	return s, nil
}

func (s *Server) Run(wg *sync.WaitGroup) error {
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

	cfg := s.Config()
	serverScheme := cfg.GetSbiScheme()

	var err error
	if serverScheme == "http" {
		err = s.httpServer.ListenAndServe()
	} else if serverScheme == "https" {
		// TODO: support TLS mutual authentication for OAuth
		err = s.httpServer.ListenAndServeTLS(
			cfg.GetNrfCertPemPath(),
			cfg.GetNrfPrivKeyPath())
	} else {
		err = fmt.Errorf("No support this scheme[%s]", serverScheme)
	}

	if err != nil && err != http.ErrServerClosed {
		logger.SBILog.Errorf("SBI server error: %v", err)
	}
	logger.SBILog.Warnf("SBI server (listen on %s) stopped", s.httpServer.Addr)
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
