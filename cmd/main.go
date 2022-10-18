package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"
	"github.com/urfave/cli"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/internal/util"
	nrf_service "github.com/free5gc/nrf/pkg/service"
	"github.com/free5gc/util/version"
)

var NRF = &nrf_service.NRF{}

func main() {
	app := cli.NewApp()
	app.Name = "nrf"
	app.Usage = "5G Network Repository Function (NRF)"
	app.Action = action
	app.Flags = NRF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorf("NRF Run Error: %v\n", err)
	}
}

func action(c *cli.Context) error {
	if err := initLogFile(c.String("log"), c.String("log5gc")); err != nil {
		logger.AppLog.Errorf("%+v", err)
		return err
	}

	if err := NRF.Initialize(c); err != nil {
		switch errType := err.(type) {
		case govalidator.Errors:
			validErrs := err.(govalidator.Errors).Errors()
			for _, validErr := range validErrs {
				logger.CfgLog.Errorf("%+v", validErr)
			}
		default:
			logger.CfgLog.Errorf("%+v", errType)
		}
		logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
		return fmt.Errorf("Failed to initialize !!")
	}

	logger.AppLog.Infoln(c.App.Name)
	logger.AppLog.Infoln("NRF version: ", version.GetVersion())

	NRF.Start()

	return nil
}

func initLogFile(logNfPath, log5gcPath string) error {
	NRF.KeyLogPath = util.NrfDefaultKeyLogPath

	if err := logger.LogFileHook(logNfPath, log5gcPath); err != nil {
		return err
	}

	if logNfPath != "" {
		nfDir, _ := filepath.Split(logNfPath)
		tmpDir := filepath.Join(nfDir, "key")
		if err := os.MkdirAll(tmpDir, 0o775); err != nil {
			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
			return err
		}
		_, name := filepath.Split(util.NrfDefaultKeyLogPath)
		NRF.KeyLogPath = filepath.Join(tmpDir, name)
	}

	return nil
}
