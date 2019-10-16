package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"free5gc/src/amf/amf_service"
	"free5gc/src/amf/logger"
	"free5gc/src/app"
	"os"
)

var AMF = &amf_service.AMF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "amf"
	appLog.Infoln(app.Name)
	app.Usage = "-free5gccfg common configuration file -amfcfg amf configuration file"
	app.Action = action
	app.Flags = AMF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorf("AMF Run error: %v", err)
	}
}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gccfg"))
	AMF.Initialize(c)
	AMF.Start()
}
