package nrf_service

import (
	"bufio"
	"fmt"
	"free5gc/lib/path_util"
	"free5gc/src/nrf/nrf_context"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/http2_util"
	"free5gc/src/app"
	"free5gc/src/nrf/AccessToken"
	"free5gc/src/nrf/Discovery"
	"free5gc/src/nrf/Management"
	"free5gc/src/nrf/factory"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/nrf_handler"
	"free5gc/src/nrf/nrf_util"
)

type NRF struct{}

type (
	// Config information.
	Config struct {
		nrfcfg string
	}
)

var config Config

var nrfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "nrfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*NRF) GetCliCmd() (flags []cli.Flag) {
	return nrfCLi
}

func (*NRF) Initialize(c *cli.Context) {

	config = Config{
		nrfcfg: c.String("nrfcfg"),
	}

	if config.nrfcfg != "" {
		factory.InitConfigFactory(config.nrfcfg)
	} else {
		factory.InitConfigFactory(path_util.Gofree5gcPath("free5gc/config/nrfcfg.conf"))
	}

	initLog.Traceln("NRF debug level(string):", app.ContextSelf().Logger.NRF.DebugLevel)
	if app.ContextSelf().Logger.NRF.DebugLevel != "" {
		initLog.Infoln("NRF debug level(string):", app.ContextSelf().Logger.NRF.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.NRF.DebugLevel)
		if err != nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.NRF.ReportCaller)
}

func (nrf *NRF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range nrf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (nrf *NRF) Start() {
	MongoDBLibrary.SetMongoDB(factory.NrfConfig.Configuration.MongoDBName, factory.NrfConfig.Configuration.MongoDBUrl)
	initLog.Infoln("Server started")

	router := gin.Default()

	AccessToken.AddService(router)
	Discovery.AddService(router)
	Management.AddService(router)

	nrf_context.InitNrfContext()

	go nrf_handler.Handle()

	uri := fmt.Sprintf("%s:%d", factory.NrfConfig.Configuration.Sbi.IPv4Addr, factory.NrfConfig.Configuration.Sbi.Port)
	initLog.Infoln(uri)
	server, err := http2_util.NewServer(uri, nrf_util.NrfLogPath, router)
	if err == nil && server != nil {
		initLog.Infoln(server.ListenAndServeTLS(nrf_util.NrfPemPath, nrf_util.NrfKeyPath))
	}
}

func (nrf *NRF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("nrfcfg"))
	args := nrf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./nrf", args...)

	nrf.Initialize(c)

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
		fmt.Println("NRF log start")
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		fmt.Println("NRF  start")
		if err := command.Start(); err != nil {
			fmt.Printf("NRF Start error: %v", err)
		}
		fmt.Println("NRF  end")
		wg.Done()
	}()

	wg.Wait()

	return err
}
