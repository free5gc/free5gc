package n3iwf_service

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	//"free5gc/src/amf/Communication"
	//"free5gc/src/amf/EventExposure"
	//"free5gc/src/amf/amf_context"
	"free5gc/src/app"
	"free5gc/src/n3iwf/logger"
	"os/exec"
	"sync"
)

type N3IWF struct{}

type (
	// Config information.
	Config struct {
		n3iwfcfg string
	}
)

var config Config

var n3iwfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "cfg",
		Usage: "n3iwf configuration file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*N3IWF) GetCliCmd() (flags []cli.Flag) {
	return n3iwfCLi
}

func (*N3IWF) Initialize(cfgPath string, c *cli.Context) {

	config = Config{
		n3iwfcfg: c.String("n3iwfcfg"),
	}

	app.AppInitializeWillInitialize(cfgPath)

	initLog.Traceln("N3IWF debug level(string):", app.ContextSelf().Logger.N3IWF.DebugLevel)
	if app.ContextSelf().Logger.N3IWF.DebugLevel != "" {
		initLog.Infoln("W3IWF debug level(string):", app.ContextSelf().Logger.N3IWF.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.N3IWF.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.N3IWF.ReportCaller)

}

func (n3iwf *N3IWF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range n3iwf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (n3iwf *N3IWF) Start() {
	initLog.Infoln("Server started")

	//self := amf_context.AMF_Self()
	//supi := "imsi-0010202"
	//ue := self.NewAmfUe(supi)
	//ue.GroupID = "12121212-208-93-01010101"
	//ue.TimeZone = "UTC"

}

func (n3iwf *N3IWF) Exec(c *cli.Context) error {

	//N3IWF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("n3iwfcfg"))
	args := n3iwf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./n3iwf", args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	stderr, err := command.StderrPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		command.Start()
		wg.Done()
	}()

	wg.Wait()

	return err
}
