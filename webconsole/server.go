package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"free5gc/src/app"
	"free5gc/webconsole/backend/logger"
	"free5gc/webconsole/backend/webui_service"
	"os"
)

var WEBUI = &webui_service.WEBUI{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "webui"
	appLog.Infoln(app.Name)
	app.Usage = "-free5gccfg common configuration file -webuicfg webui configuration file"
	app.Action = action
	app.Flags = WEBUI.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Warnf("Error args: %v", err)
	}
}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gccfg"))
	WEBUI.Initialize(c)
	WEBUI.Start()
}
