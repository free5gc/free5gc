package main

import (
	"math/rand"
	"os"
	"runtime/debug"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/free5gc/go-upf/internal/logger"
	upfapp "github.com/free5gc/go-upf/pkg/app"
	"github.com/free5gc/go-upf/pkg/factory"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/version"
)

func main() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	app := cli.NewApp()
	app.Name = "upf"
	app.Usage = "5G User Plane Function (UPF)"
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

	// rand.Seed(time.Now().UnixNano()) // rand.Seed has been deprecated
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	randSeed.Uint64()

	if err := app.Run(os.Args); err != nil {
		logger.MainLog.Errorf("UPF Cli Run Error: %v", err)
	}
}

func action(cliCtx *cli.Context) error {
	err := initLogFile(cliCtx.StringSlice("log"))
	if err != nil {
		return err
	}

	logger.MainLog.Infoln("UPF version: ", version.GetVersion())

	cfg, err := factory.ReadConfig(cliCtx.String("config"))
	if err != nil {
		return err
	}

	upf, err := upfapp.NewApp(cfg)
	if err != nil {
		return err
	}

	if err := upf.Run(); err != nil {
		return err
	}

	return nil
}

func initLogFile(logNfPath []string) error {
	for _, path := range logNfPath {
		if err := logger_util.LogFileHook(logger.Log, path); err != nil {
			return err
		}
	}
	return nil
}
