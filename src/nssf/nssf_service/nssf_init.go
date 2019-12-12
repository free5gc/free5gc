/*
 * NSSF Service
 */

package nssf_service

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"free5gc/lib/http2_util"
	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/nssf/NSSAIAvailability"
	"free5gc/src/nssf/NSSelection"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/logger"
	"free5gc/src/nssf/nssf_consumer"
	"free5gc/src/nssf/nssf_context"
	"free5gc/src/nssf/nssf_handler"
	"free5gc/src/nssf/util"
)

type NSSF struct{}

type (
	// Config information.
	Config struct {
		nssfcfg string
	}
)

var config Config

var nssfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "nssfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*NSSF) GetCliCmd() (flags []cli.Flag) {
	return nssfCLi
}

func (*NSSF) Initialize(c *cli.Context) {

	config = Config{
		nssfcfg: c.String("nssfcfg"),
	}

	if config.nssfcfg != "" {
		initLog.Infof("Read configuration file %s", config.nssfcfg)
		factory.InitConfigFactory(path_util.Gofree5gcPath(config.nssfcfg))
	} else {
		initLog.Warnln("No configuration file is provided")
		factory.InitConfigFactory(path_util.Gofree5gcPath(util.DEFAULT_CONFIG))
	}

	nssf_context.InitNssfContext()

	initLog.Traceln("NSSF debug level(string):", app.ContextSelf().Logger.NSSF.DebugLevel)
	if app.ContextSelf().Logger.NSSF.DebugLevel != "" {
		initLog.Infoln("NSSF debug level(string):", app.ContextSelf().Logger.NSSF.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.NSSF.DebugLevel)
		if err != nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.NSSF.ReportCaller)
}

func (nssf *NSSF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range nssf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (nssf *NSSF) Start() {
	initLog.Infoln("Server started")

	router := gin.Default()

	NSSAIAvailability.AddService(router)
	NSSelection.AddService(router)

	go nssf_handler.Handle()

	self := nssf_context.NSSF_Self()
	addr := fmt.Sprintf("%s:%d", self.HttpIpv4Address, self.Port)

	// Register to NRF
	profile, err := nssf_consumer.BuildNFProfile(self)
	if err != nil {
		initLog.Error("Failed to build NSSF profile")
	}
	_, self.NfId, err = nssf_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	if err != nil {
		initLog.Errorf("Failed to register NSSF to NRF: %s", err.Error())
	}

	server, err := http2_util.NewServer(addr, util.NSSF_LOG_PATH, router)
	if err == nil && server != nil {
		initLog.Infoln(server.ListenAndServeTLS(util.NSSF_PEM_PATH, util.NSSF_KEY_PATH))
	}
}

func (nssf *NSSF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("nssfcfg"))
	args := nssf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./nssf", args...)

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
		if err := command.Start(); err != nil {
			fmt.Printf("NSSF Start error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
