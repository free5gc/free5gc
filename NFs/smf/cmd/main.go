package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/pkg/factory"
	"github.com/free5gc/smf/pkg/service"
	"github.com/free5gc/smf/pkg/utils"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/version"
)

var SMF *service.SmfApp

func main() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	app := cli.NewApp()
	app.Name = "smf"
	app.Usage = "5G Session Management Function (SMF)"
	app.Action = action
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
		},
		&cli.StringFlag{
			Name:    "uerouting",
			Aliases: []string{"u"},
			Usage:   "Load uerouting configuration from `FILE`",
		},
		&cli.StringSliceFlag{
			Name:    "log",
			Aliases: []string{"l"},
			Usage:   "Output NF log to `FILE`",
		},
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))

	if err := app.Run(os.Args); err != nil {
		logger.MainLog.Errorf("SMF Run error: %v\n", err)
	}
}

func action(cliCtx *cli.Context) error {
	tlsKeyLogPath, err := initLogFile(cliCtx.StringSlice("log"))
	if err != nil {
		return err
	}

	logger.MainLog.Infoln("SMF version: ", version.GetVersion())

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh  // Wait for interrupt signal to gracefully shutdown
		cancel() // Notify each goroutine and wait them stopped
	}()

	cfg, err := factory.ReadConfig(cliCtx.String("config"))
	if err != nil {
		sigCh <- nil
		return err
	}
	factory.SmfConfig = cfg

	ueRoutingCfg, err := factory.ReadUERoutingConfig(cliCtx.String("uerouting"))
	if err != nil {
		sigCh <- nil
		return err
	}
	factory.UERoutingConfig = ueRoutingCfg

	pfcpStart, pfcpTerminate := utils.InitPFCPFunc(ctx)
	smf, err := service.NewApp(ctx, cfg, tlsKeyLogPath, pfcpStart, pfcpTerminate)
	if err != nil {
		sigCh <- nil
		return err
	}
	SMF = smf

	smf.Start()

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
		_, name := filepath.Split(factory.SmfDefaultTLSKeyLogPath)
		logTlsKeyPath = filepath.Join(tmpDir, name)
	}

	return logTlsKeyPath, nil
}
