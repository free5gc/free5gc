package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_service"
	"os"
)

type (
	// Config information.
	Config struct {
		n3iwfcfg string
	}
)

var config Config

var N3IWF = &n3iwf_service.N3IWF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "n3iwf"
	appLog.Infoln(app.Name)
	app.Usage = "This is a network function. The abrreviation of Non-3GPP Interworking Function."
	app.Action = action
	app.Flags = N3IWF.GetCliCmd()
	app.Run(os.Args)

}

func action(c *cli.Context) {
	config = Config{
		n3iwfcfg: c.String("n3iwfcfg"),
	}
	appLog.Debugln("n3iwfcfgFile:", config.n3iwfcfg)

	N3IWF.Initialize(config.n3iwfcfg, c)
	N3IWF.Start()
}
