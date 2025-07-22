package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/free5gc/nef/internal/logger"
	nefapp "github.com/free5gc/nef/pkg/app"
	"github.com/free5gc/nef/pkg/factory"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/version"
	"github.com/urfave/cli/v2"
)

func main() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	app := cli.NewApp()
	app.Name = "nef"
	app.Usage = "5G Network Exposure Function (NEF)"
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
		logger.MainLog.Errorf("NEF Cli Run err: %v\n", err)
	}
}

func action(cliCtx *cli.Context) error {
	tlsKeyLogPath, err := initLogFile(cliCtx.StringSlice("log"))
	if err != nil {
		return err
	}

	logger.MainLog.Infoln("NEF version: ", version.GetVersion())

	cfg, err := factory.ReadConfig(cliCtx.String("config"))
	if err != nil {
		return err
	}

	nef, err := nefapp.NewApp(cfg, tlsKeyLogPath)
	if err != nil {
		return fmt.Errorf("new NEF err: %+v", err)
	}

	if err := nef.Run(); err != nil {
		return nil
	}

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
		_, name := filepath.Split(factory.NefDefaultTLSKeyLogPath)
		logTlsKeyPath = filepath.Join(tmpDir, name)
	}

	return logTlsKeyPath, nil
}
