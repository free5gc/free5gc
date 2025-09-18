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
	"github.com/free5gc/bsf/internal/sbi/processor"
	"github.com/free5gc/bsf/pkg/factory"
	"github.com/free5gc/util/metrics"
	sbiMetrics "github.com/free5gc/util/metrics/sbi"
	"github.com/free5gc/util/metrics/utils"
)

type App struct {
	ctx           context.Context
	config        *factory.Config
	bsfCtx        *bsfContext.BSFContext
	metricsServer *metrics.Server
	consumer      *consumer.Consumer
	wg            sync.WaitGroup
}

func NewApp(cfg *factory.Config) (*App, error) {
	bsf := &App{
		config: cfg,
		bsfCtx: bsfContext.BsfSelf,
	}

	// Initialize consumer
	var err error
	if bsf.consumer, err = consumer.NewConsumer(bsf); err != nil {
		return nil, fmt.Errorf("failed to initialize consumer: %w", err)
	}

	// Initialize processor singleton
	if _, err := processor.NewProcessor(bsf); err != nil {
		return nil, fmt.Errorf("failed to initialize processor: %w", err)
	}

	// Set BSF context configuration
	bsf.bsfCtx.NrfUri = cfg.Configuration.NrfUri

	// Initialize metrics if enabled - need to check proper method name
	var tlsKeyLogPath string
	if cfg.AreMetricsEnabled() {
		sbiMetrics.EnableSbiMetrics()

		features := map[utils.MetricTypeEnabled]bool{utils.SBI: true}
		customMetrics := make(map[utils.MetricTypeEnabled][]prometheus.Collector)

		var err error
		if bsf.metricsServer, err = metrics.NewServer(
			getInitMetrics(cfg, features, customMetrics), tlsKeyLogPath, logger.MainLog); err != nil {
			logger.MainLog.Warnf("Failed to create metrics server: %+v", err)
		}
	}

	return bsf, nil
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

	// Enable business metrics if configured - preserve your existing approach
	if cfg.AreMetricsEnabled() {
		businessMetrics.EnableBindingMetrics()
		businessMetrics.EnableDiscoveryMetrics()

		// Add BSF business metrics using your existing functions
		if customMetrics == nil {
			customMetrics = make(map[utils.MetricTypeEnabled][]prometheus.Collector)
		}

		// Add binding metrics using your existing function
		customMetrics[utils.SBI] = append(
			customMetrics[utils.SBI],
			businessMetrics.GetBindingHandlerMetrics(cfg.GetMetricsNamespace())...)

		// Add discovery metrics using your existing function
		customMetrics[utils.SBI] = append(
			customMetrics[utils.SBI],
			businessMetrics.GetDiscoveryHandlerMetrics(cfg.GetMetricsNamespace())...)
	}

	return metrics.NewInitMetrics(metricsInfo, "bsf", features, customMetrics)
}

func (a *App) Config() *factory.Config {
	return a.config
}

func (a *App) Context() *bsfContext.BSFContext {
	return a.bsfCtx
}

func (a *App) CancelContext() context.Context {
	return a.ctx
}

func (a *App) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *App) Start() error {
	defer func() {
		if p := recover(); p != nil {
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	// Initialize MongoDB connection
	if err := a.bsfCtx.ConnectMongoDB(); err != nil {
		logger.MainLog.Warnf("MongoDB connection failed: %+v", err)
	} else {
		// Load existing bindings from MongoDB
		if err := a.bsfCtx.LoadPcfBindingsFromMongoDB(); err != nil {
			logger.MainLog.Warnf("Failed to load PCF bindings from MongoDB: %+v", err)
		}
	}

	// Start cleanup routine for expired and inactive bindings
	a.bsfCtx.StartCleanupRoutine()

	// Start metrics server if enabled
	if a.config.AreMetricsEnabled() && a.metricsServer != nil {
		go func() {
			a.metricsServer.Run(&a.wg)
		}()
		logger.MainLog.Infof("BSF metrics server enabled on %s://%s",
			a.config.GetMetricsScheme(), a.config.GetMetricsBindingAddr())
	}

	// Register with NRF - moved to consumer
	if err := a.consumer.RegisterWithNRF(); err != nil {
		logger.MainLog.Errorf("BSF register to NRF Error: %+v", err)
		return fmt.Errorf("failed to register with NRF: %w", err)
	}
	logger.MainLog.Infof("BSF successfully registered with NRF")

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

	bindAddr := fmt.Sprintf("%s:%d", a.config.Configuration.Sbi.BindingIPv4, a.config.Configuration.Sbi.Port)
	logger.MainLog.Infof("BSF SBI Server started on %s", bindAddr)

	server := &http.Server{
		Addr:           bindAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if a.config.Configuration.Sbi.Scheme == "http" {
		return server.ListenAndServe()
	} else if a.config.Configuration.Sbi.Scheme == "https" {
		return server.ListenAndServeTLS(
			a.config.Configuration.Sbi.Tls.Pem,
			a.config.Configuration.Sbi.Tls.Key,
		)
	}

	return fmt.Errorf("unsupported scheme: %s", a.config.Configuration.Sbi.Scheme)
}

func (a *App) Terminate() error {
	logger.MainLog.Infof("Terminating BSF...")

	// Deregister from NRF using consumer
	if err := a.consumer.DeregisterWithNRF(); err != nil {
		logger.MainLog.Errorf("BSF deregister from NRF Error: %+v", err)
		// Don't return error here as termination should continue
	}

	// Stop cleanup routine
	a.bsfCtx.StopCleanupRoutine()

	// Disconnect from MongoDB
	if err := a.bsfCtx.DisconnectMongoDB(); err != nil {
		logger.MainLog.Errorf("Error disconnecting from MongoDB: %+v", err)
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	return nil
}
