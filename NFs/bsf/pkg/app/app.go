/*
 * BSF App
 */

package app

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/logger"
	businessMetrics "github.com/free5gc/bsf/internal/metrics/business"
	"github.com/free5gc/bsf/internal/sbi"
	"github.com/free5gc/bsf/internal/sbi/consumer"
	"github.com/free5gc/bsf/pkg/factory"
	"github.com/free5gc/util/metrics"
	sbiMetrics "github.com/free5gc/util/metrics/sbi"
	"github.com/free5gc/util/metrics/utils"
)

type App struct {
	cfg           *factory.Config
	ctx           context.Context
	tlsKeyPath    string
	bsfCtx        *bsfContext.BSFContext
	metricsServer *metrics.Server
	wg            sync.WaitGroup
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) *App {
	a := &App{
		cfg:        cfg,
		ctx:        ctx,
		tlsKeyPath: tlsKeyLogPath,
		bsfCtx:     bsfContext.BsfSelf,
	}

	// Initialize metrics if enabled
	if a.cfg.AreMetricsEnabled() {
		sbiMetrics.EnableSbiMetrics()

		features := map[utils.MetricTypeEnabled]bool{utils.SBI: true}
		customMetrics := make(map[utils.MetricTypeEnabled][]prometheus.Collector)

		var err error
		if a.metricsServer, err = metrics.NewServer(
			getInitMetrics(cfg, features, customMetrics), tlsKeyLogPath, logger.MainLog); err != nil {
			logger.MainLog.Warnf("Failed to create metrics server: %+v", err)
		}
	}

	return a
}

func getInitMetrics(
	cfg *factory.Config,
	features map[utils.MetricTypeEnabled]bool,
	customMetrics map[utils.MetricTypeEnabled][]prometheus.Collector,
) metrics.InitMetrics {
	metricsInfo := metrics.Metrics{
		BindingIPv4: cfg.GetMetricsBindingAddr(),
		Scheme:      cfg.GetMetricsScheme(),
		Namespace:   cfg.GetMetricsNamespace(),
		Port:        cfg.GetMetricsPort(),
		Tls: metrics.Tls{
			Key: cfg.GetMetricsCertKeyPath(),
			Pem: cfg.GetMetricsCertPemPath(),
		},
	}

	// Enable business metrics if configured
	if cfg.AreMetricsEnabled() {
		businessMetrics.EnableBindingMetrics()
		businessMetrics.EnableDiscoveryMetrics()

		// Add BSF business metrics
		if customMetrics == nil {
			customMetrics = make(map[utils.MetricTypeEnabled][]prometheus.Collector)
		}

		// Add binding metrics
		customMetrics[utils.SBI] = append(
			customMetrics[utils.SBI],
			businessMetrics.GetBindingHandlerMetrics(cfg.GetMetricsNamespace())...)

		// Add discovery metrics
		customMetrics[utils.SBI] = append(
			customMetrics[utils.SBI],
			businessMetrics.GetDiscoveryHandlerMetrics(cfg.GetMetricsNamespace())...)
	}

	return metrics.NewInitMetrics(metricsInfo, "bsf", features, customMetrics)
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

	// Start metrics server if enabled
	if a.cfg.AreMetricsEnabled() && a.metricsServer != nil {
		go func() {
			a.metricsServer.Run(&a.wg)
		}()
		logger.MainLog.Infof("BSF metrics server enabled on %s://%s",
			a.cfg.GetMetricsScheme(), a.cfg.GetMetricsBindingAddr())
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
	sbi.AddService(router) // Add CORS
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
