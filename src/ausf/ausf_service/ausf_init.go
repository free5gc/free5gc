package ausf_service

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"free5gc/lib/http2_util"
	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/ausf/UEAuthentication"
	"free5gc/src/ausf/ausf_consumer"
	"free5gc/src/ausf/ausf_context"
	"free5gc/src/ausf/ausf_handler"
	"free5gc/src/ausf/ausf_util"
	"free5gc/src/ausf/factory"
	"free5gc/src/ausf/logger"
	"os/exec"
	"sync"
)

type AUSF struct{}

type (
	// Config information.
	Config struct {
		ausfcfg string
	}
)

var config Config

var ausfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "ausfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*AUSF) GetCliCmd() (flags []cli.Flag) {
	return ausfCLi
}

func (*AUSF) Initialize(c *cli.Context) {

	config = Config{
		ausfcfg: c.String("ausfcfg"),
	}

	if config.ausfcfg != "" {
		factory.InitConfigFactory(config.ausfcfg)
	} else {
		DefaultAusfConfigPath := path_util.Gofree5gcPath("free5gc/config/ausfcfg.conf")
		factory.InitConfigFactory(DefaultAusfConfigPath)
	}

	initLog.Traceln("AUSF debug level(string):", app.ContextSelf().Logger.AUSF.DebugLevel)
	if app.ContextSelf().Logger.AUSF.DebugLevel != "" {
		initLog.Infoln("AUSF debug level(string):", app.ContextSelf().Logger.AUSF.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.AUSF.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.AUSF.ReportCaller)

}

func (ausf *AUSF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range ausf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (ausf *AUSF) Start() {
	initLog.Infoln("Server started")

	router := gin.Default()
	UEAuthentication.AddService(router)

	ausf_context.Init()
	self := ausf_context.GetSelf()
	// Register to NRF
	profile, err := ausf_consumer.BuildNFInstance(self)
	if err != nil {
		initLog.Error("Build AUSF Profile Error")
	}
	_, self.NfId, err = ausf_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	if err != nil {
		initLog.Errorf("AUSF register to NRF Error[%s]", err.Error())
	}

	ausfLogPath := ausf_util.AusfLogPath
	ausfPemPath := ausf_util.AusfPemPath
	ausfKeyPath := ausf_util.AusfKeyPath

	go ausf_handler.Handle()
	server, err := http2_util.NewServer(":29509", ausfLogPath, router)
	if err == nil && server != nil {
		initLog.Infoln(server.ListenAndServeTLS(ausfPemPath, ausfKeyPath))
	}
}

func (ausf *AUSF) Exec(c *cli.Context) error {

	//AUSF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("ausfcfg"))
	args := ausf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./ausf", args...)

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
		startErr := command.Start()
		if startErr != nil {
			initLog.Fatalln(startErr)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
