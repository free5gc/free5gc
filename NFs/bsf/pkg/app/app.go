/*
 * BSF App
 */

package app

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/logger"
	"github.com/free5gc/bsf/internal/sbi"
	"github.com/free5gc/bsf/internal/sbi/consumer"
	"github.com/free5gc/bsf/pkg/factory"
)

type App struct {
	cfg        *factory.Config
	ctx        context.Context
	tlsKeyPath string
	bsfCtx     *bsfContext.BSFContext
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) *App {
	return &App{
		cfg:        cfg,
		ctx:        ctx,
		tlsKeyPath: tlsKeyLogPath,
		bsfCtx:     bsfContext.BsfSelf,
	}
}

func (a *App) Start() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	// Initialize MongoDB connection
	if err := a.bsfCtx.ConnectMongoDB(); err != nil {
		logger.MainLog.Warnf("MongoDB connection failed: %+v", err)
	}

	// Register with NRF
	go func() {
		if _, err := consumer.SendRegisterNFInstance(); err != nil {
			logger.MainLog.Errorf("BSF register to NRF Error[%+v]", err)
		} else {
			logger.MainLog.Infof("BSF successfully registered with NRF")
		}
	}()

	// Start SBI server
	router := gin.Default()
	sbi.AddService(router)

	// Add CORS
	router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "User-Agent",
			"Referrer", "Host", "Token", "X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))

	bindAddr := fmt.Sprintf("%s:%d", a.cfg.Configuration.Sbi.BindingIPv4, a.cfg.Configuration.Sbi.Port)
	logger.MainLog.Infof("BSF SBI Server started on %s", bindAddr)

	server := &http.Server{
		Addr:           bindAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if a.cfg.Configuration.Sbi.Scheme == "http" {
		err := server.ListenAndServe()
		if err != nil {
			logger.MainLog.Fatalf("HTTP server setup failed: %+v", err)
		}
	} else if a.cfg.Configuration.Sbi.Scheme == "https" {
		err := server.ListenAndServeTLS(
			a.cfg.Configuration.Sbi.Tls.Pem,
			a.cfg.Configuration.Sbi.Tls.Key,
		)
		if err != nil {
			logger.MainLog.Fatalf("HTTPS server setup failed: %+v", err)
		}
	}
}
