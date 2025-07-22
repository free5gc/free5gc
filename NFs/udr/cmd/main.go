package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/free5gc/udr/internal/logger"
	"github.com/free5gc/udr/pkg/factory"
	"github.com/free5gc/udr/pkg/service"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/version"
)

var UDR *service.UdrApp

func main() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	app := cli.NewApp()
	app.Name = "udr"
	app.Usage = "5G Unified Data Repository (UDR)"
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
		logger.MainLog.Errorf("UDR Run error: %v\n", err)
	}
}

func action(cliCtx *cli.Context) error {
	tlsKeyLogPath, err := initLogFile(cliCtx.StringSlice("log"))
	if err != nil {
		return err
	}

	logger.MainLog.Infoln("UDR version: ", version.GetVersion())
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	cfg, err := factory.ReadConfig(cliCtx.String("config"))
	if err != nil {
		return err
	}
	factory.UdrConfig = cfg
	udr, err := service.NewApp(ctx, cfg, tlsKeyLogPath)
	if err != nil {
		return err
	}
	UDR = udr

	udr.Start()

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
		_, name := filepath.Split(factory.UdrDefaultTLSKeyLogPath)
		logTlsKeyPath = filepath.Join(tmpDir, name)
	}

	return logTlsKeyPath, nil
}
