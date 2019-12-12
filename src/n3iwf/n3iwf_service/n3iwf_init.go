package n3iwf_service

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/n3iwf/factory"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_handler"
	"free5gc/src/n3iwf/n3iwf_ngap/n3iwf_sctp"
	"free5gc/src/n3iwf/n3iwf_util"
	//"free5gc/src/n3iwf/n3iwf_context"
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
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "n3iwfcfg",
		Usage: "n3iwf config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*N3IWF) GetCliCmd() (flags []cli.Flag) {
	return n3iwfCLi
}

func (*N3IWF) Initialize(c *cli.Context) {

	config = Config{
		n3iwfcfg: c.String("n3iwfcfg"),
	}

	if config.n3iwfcfg != "" {
		factory.InitConfigFactory(path_util.Gofree5gcPath(config.n3iwfcfg))
	} else {
		DefaultSmfConfigPath := path_util.Gofree5gcPath("free5gc/config/n3iwfcfg.conf")
		factory.InitConfigFactory(DefaultSmfConfigPath)
	}

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

	n3iwf_util.InitN3IWFContext()

	go n3iwf_handler.Handle()

	wg := sync.WaitGroup{}

	n3iwf_sctp.InitiateSCTP(&wg)

	wg.Wait()

	//self := n3iwf_context.N3IWFSelf()
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

	wg := sync.WaitGroup{}
	wg.Add(3)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
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
		if err := command.Start(); err != nil {
			initLog.Errorf("N3IWF start error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
