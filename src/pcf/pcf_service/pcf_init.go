package pcf_service

import (
	"bufio"
	"fmt"
	"free5gc/lib/http2_util"
	"free5gc/src/app"
	"free5gc/src/pcf/AMPolicy"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler"
	"free5gc/src/pcf/pcf_util"

	"free5gc/src/pcf/BDTPolicy"
	"free5gc/src/pcf/PolicyAuthorization"
	"free5gc/src/pcf/SMPolicy"
	"free5gc/src/pcf/UEPolicy"
	"free5gc/src/pcf/logger"
	"os/exec"
	"sync"

	"free5gc/src/pcf/factory"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type PCF struct{}

type (
	// Config information.
	Config struct {
		pcfcfg string
	}
)

var config Config

var pcfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "pcfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*PCF) GetCliCmd() (flags []cli.Flag) {
	return pcfCLi
}

func (*PCF) Initialize(c *cli.Context) {

	config = Config{
		pcfcfg: c.String("pcfcfg"),
	}
	if config.pcfcfg != "" {
		factory.InitConfigFactory(config.pcfcfg)
	} else {
		DefaultPcfConfigPath := pcf_util.PCF_CONFIG_PATH
		factory.InitConfigFactory(DefaultPcfConfigPath)
	}

	initLog.Traceln("PCF debug level(string):", app.ContextSelf().Logger.PCF.DebugLevel)
	if app.ContextSelf().Logger.PCF.DebugLevel != "" {
		initLog.Infoln("PCF debug level(string):", app.ContextSelf().Logger.PCF.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.PCF.DebugLevel)
		if err != nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.PCF.ReportCaller)
}

func (pcf *PCF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range pcf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (pcf *PCF) Start() {
	initLog.Infoln("Server started")
	router := gin.Default()

	BDTPolicy.AddService(router)
	SMPolicy.AddService(router)
	AMPolicy.AddService(router)
	UEPolicy.AddService(router)
	PolicyAuthorization.AddService(router)

	go pcf_handler.Handle()
	self := pcf_context.PCF_Self()
	pcf_util.InitpcfContext(self)
	addr := fmt.Sprintf("%s:%d", self.HttpIPv4Address, self.HttpIpv4Port)
	server, err := http2_util.NewServer(addr, pcf_util.PCF_LOG_PATH, router)
	if err == nil && server != nil {
		initLog.Infoln(server.ListenAndServeTLS(pcf_util.PCF_PEM_PATH, pcf_util.PCF_KEY_PATH))
	}
}

func (pcf *PCF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("pcfcfg"))
	args := pcf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./pcf", args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(4)
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
		fmt.Println("PCF log start")
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		fmt.Println("PCF start")
		if err := command.Start(); err != nil {
			fmt.Printf("command.Start() error: %v", err)
		}
		fmt.Println("PCF end")
		wg.Done()
	}()

	wg.Wait()

	return err
}
