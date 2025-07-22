package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/nrf/pkg/service"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/version"
)

var NRF *service.NrfApp

func main() {
	app := cli.NewApp()
	app.Name = "nrf"
	app.Usage = "5G Network Repository Function (NRF)"
	app.Action = action
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
		},
		&cli.StringSliceFlag{
			Name:    "log",
			Aliases: []string{"l"},
			Usage:   "Output NF log to `FILE`",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.MainLog.Errorf("NRF Run Error: %v\n", err)
	}
}

func action(cliCtx *cli.Context) error {
	tlsKeyLogPath, err := initLogFile(cliCtx.StringSlice("log"))
	if err != nil {
		return err
	}
	logger.MainLog.Infoln("NRF version: ", version.GetVersion())

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh  // Wait for interrupt signal to gracefully shutdown
		cancel() // Notify each goroutine and wait them stopped
	}()

	cfg, err := factory.ReadConfig(cliCtx.String("config"))
	if err != nil {
		return err
	}
	factory.NrfConfig = cfg

	nrf, err := service.NewApp(ctx, cfg, tlsKeyLogPath)
	if err != nil {
		return err
	}
	NRF = nrf

	nrf.Start()
	return nil
}

func initLogFile(logNfPath []string) (string, error) {
	logTlsKeyPath := ""

	for _, path := range logNfPath {
		if err := logger_util.LogFileHook(logger.Log, path); err != nil {
			return "", err
		}

		if logTlsKeyPath != "" {
			continue
		}

		nfDir, _ := filepath.Split(path)
		tmpDir := filepath.Join(nfDir, "key")
		if err := os.MkdirAll(tmpDir, 0o775); err != nil {
			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
			return "", err
		}
		_, name := filepath.Split(factory.NrfDefaultTLSKeyLogPath)
		logTlsKeyPath = filepath.Join(tmpDir, name)
	}
	return logTlsKeyPath, nil
}
