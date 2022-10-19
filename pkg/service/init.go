package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	nrf_context "github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/internal/sbi/accesstoken"
	"github.com/free5gc/nrf/internal/sbi/discovery"
	"github.com/free5gc/nrf/internal/sbi/management"
	"github.com/free5gc/nrf/internal/util"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/mongoapi"
)

type NRF struct {
	KeyLogPath string
}

type (
	// Commands information.
	Commands struct {
		config string
	}
)

var commands Commands

var cliCmd = []cli.Flag{
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Load configuration from `FILE`",
	},
	cli.StringFlag{
		Name:  "log, l",
		Usage: "Output NF log to `FILE`",
	},
	cli.StringFlag{
		Name:  "log5gc, lc",
		Usage: "Output free5gc log to `FILE`",
	},
}

func (*NRF) GetCliCmd() (flags []cli.Flag) {
	return cliCmd
}

func (nrf *NRF) Initialize(c *cli.Context) error {
	commands = Commands{
		config: c.String("config"),
	}

	if commands.config != "" {
		if err := factory.InitConfigFactory(commands.config); err != nil {
			return err
		}
	} else {
		if err := factory.InitConfigFactory(util.NrfDefaultConfigPath); err != nil {
			return err
		}
	}

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	if _, err := factory.NrfConfig.Validate(); err != nil {
		return err
	}

	nrf.SetLogLevel()

	return nil
}

func (nrf *NRF) SetLogLevel() {
	if factory.NrfConfig.Logger == nil {
		logger.InitLog.Warnln("NRF config without log level setting!!!")
		return
	}

	if factory.NrfConfig.Logger.NRF != nil {
		if factory.NrfConfig.Logger.NRF.DebugLevel != "" {
			level, err := logrus.ParseLevel(factory.NrfConfig.Logger.NRF.DebugLevel)
			if err != nil {
				logger.InitLog.Warnf("NRF Log level [%s] is invalid, set to [info] level",
					factory.NrfConfig.Logger.NRF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				logger.InitLog.Infof("NRF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			logger.InitLog.Infoln("NRF Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.NrfConfig.Logger.NRF.ReportCaller)
	}
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
	if err := mongoapi.SetMongoDB(factory.NrfConfig.Configuration.MongoDBName,
		factory.NrfConfig.Configuration.MongoDBUrl); err != nil {
		logger.InitLog.Errorf("SetMongoDB failed: %+v", err)
		return
	}
	logger.InitLog.Infoln("Server starting")

	router := logger_util.NewGinWithLogrus(logger.GinLog)

	accesstoken.AddService(router)
	discovery.AddService(router)
	management.AddService(router)

	err := nrf_context.InitNrfContext()
	if err != nil {
		logger.InitLog.Errorln(err)
		return
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		<-signalChannel
		// Waiting for other NFs to deregister
		time.Sleep(2 * time.Second)
		nrf.Terminate()
		os.Exit(0)
	}()

	bindAddr := factory.NrfConfig.GetSbiBindingAddr()
	logger.InitLog.Infof("Binding addr: [%s]", bindAddr)
	server, err := httpwrapper.NewHttp2Server(bindAddr, nrf.KeyLogPath, router)
	if err != nil {
		logger.InitLog.Warnf("Initialize HTTP server: +%v", err)
		return
	}

	serverScheme := factory.NrfConfig.GetSbiScheme()
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		// TODO: support TLS mutual authentication for OAuth
		err = server.ListenAndServeTLS(
			factory.NrfConfig.GetNrfCertPemPath(),
			factory.NrfConfig.GetNrfPrivKeyPath())
	}

	if err != nil {
		logger.InitLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (nrf *NRF) Exec(c *cli.Context) error {
	logger.InitLog.Traceln("args:", c.String("nrfcfg"))
	args := nrf.FilterCli(c)
	logger.InitLog.Traceln("filter: ", args)
	command := exec.Command("./nrf", args...)

	if err := nrf.Initialize(c); err != nil {
		return err
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		logger.InitLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	stderr, err := command.StderrPipe()
	if err != nil {
		logger.InitLog.Fatalln(err)
	}
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		in := bufio.NewScanner(stderr)
		fmt.Println("NRF log start")
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

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

func (nrf *NRF) Terminate() {
	logger.InitLog.Infof("Terminating NRF...")

	logger.InitLog.Infof("Remove NF Profile...")
	err := mongoapi.Drop("NfProfile")
	if err != nil {
		logger.InitLog.Errorf("Drop NfProfile collection failed: %+v", err)
	}

	logger.InitLog.Infof("NRF terminated")
}
