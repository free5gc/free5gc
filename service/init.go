package service

import (
	"bufio"
	"fmt"
	"free5gc/lib/logger_util"
	"free5gc/lib/path_util"
	nrf_context "free5gc/src/nrf/context"
	"free5gc/src/nrf/util"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/http2_util"
	"free5gc/src/app"
	"free5gc/src/nrf/accesstoken"
	"free5gc/src/nrf/discovery"
	"free5gc/src/nrf/factory"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/management"
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
		DefaultNrfConfigPath := path_util.Gofree5gcPath("free5gc/config/nrfcfg.conf")
		factory.InitConfigFactory(DefaultNrfConfigPath)
	}

	if app.ContextSelf().Logger.NRF.DebugLevel != "" {
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.NRF.DebugLevel)
		if err != nil {
			initLog.Warnf("Log level [%s] is not valid, set to [info] level", app.ContextSelf().Logger.NRF.DebugLevel)
			logger.SetLogLevel(logrus.InfoLevel)
		} else {
			logger.SetLogLevel(level)
			initLog.Infof("Log level is set to [%s] level", level)
		}
	} else {
		initLog.Infoln("Log level is default set to [info] level")
		logger.SetLogLevel(logrus.InfoLevel)
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

	router := logger_util.NewGinWithLogrus(logger.GinLog)

	accesstoken.AddService(router)
	discovery.AddService(router)
	management.AddService(router)

	nrf_context.InitNrfContext()

	uri := fmt.Sprintf("%s:%d", factory.NrfConfig.Configuration.Sbi.IPv4Addr, factory.NrfConfig.Configuration.Sbi.Port)
	initLog.Infoln(uri)
	server, err := http2_util.NewServer(uri, util.NrfLogPath, router)

	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.NrfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.NrfPemPath, util.NrfKeyPath)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
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
		if err = command.Start(); err != nil {
			fmt.Printf("NRF Start error: %v", err)
		}
		fmt.Println("NRF  end")
		wg.Done()
	}()

	wg.Wait()

	return err
}
