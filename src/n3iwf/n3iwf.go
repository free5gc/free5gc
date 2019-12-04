package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"free5gc/src/app"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_service"
	"os"
)

var N3IWF = &n3iwf_service.N3IWF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "n3iwf"
	appLog.Infoln(app.Name)
	app.Usage = "-free5gccfg common configuration file -n3iwfcfg n3iwf configuration file"
	app.Action = action
	app.Flags = N3IWF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorf("N3IWF Run Error: %v", err)
	}
}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gccfg"))
	N3IWF.Initialize(c)
	N3IWF.Start()
}
