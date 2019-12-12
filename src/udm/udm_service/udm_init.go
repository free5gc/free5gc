package udm_service

import (
	"bufio"
	"fmt"
	"free5gc/lib/http2_util"
	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/udm/EventExposure"
	"free5gc/src/udm/ParameterProvision"
	"free5gc/src/udm/SubscriberDataManagement"
	"free5gc/src/udm/UEAuthentication"
	"free5gc/src/udm/UEContextManagement"
	"free5gc/src/udm/factory"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_consumer"
	"free5gc/src/udm/udm_context"
	"free5gc/src/udm/udm_handler"
	"free5gc/src/udm/udm_util"

	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type UDM struct{}

type (
	// Config information.
	Config struct {
		udmcfg string
	}
)

var config Config

var udmCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "udmcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*UDM) GetCliCmd() (flags []cli.Flag) {
	return udmCLi
}

func (*UDM) Initialize(c *cli.Context) {

	config = Config{
		udmcfg: c.String("udmcfg"),
	}

	if config.udmcfg != "" {
		factory.InitConfigFactory(path_util.Gofree5gcPath(config.udmcfg))
	} else {
		DefaultUdmConfigPath := path_util.Gofree5gcPath("free5gc/config/udmcfg.conf")
		factory.InitConfigFactory(DefaultUdmConfigPath)
	}

	initLog.Traceln("UDM debug level(string):", app.ContextSelf().Logger.UDM.DebugLevel)
	if app.ContextSelf().Logger.UDM.DebugLevel != "" {
		initLog.Infoln("UDM debug level(string):", app.ContextSelf().Logger.UDM.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.UDM.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.UDM.ReportCaller)

}

func (udm *UDM) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range udm.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (udm *UDM) Start() {
	config := factory.UdmConfig
	configuration := config.Configuration
	sbi := configuration.Sbi
	serviceName := configuration.ServiceNameList

	initLog.Infof("UDM Config Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	initLog.Infoln("Server started")

	router := gin.Default()

	EventExposure.AddService(router)
	ParameterProvision.AddService(router)
	SubscriberDataManagement.AddService(router)
	UEAuthentication.AddService(router)
	UEContextManagement.AddService(router)

	udmLogPath := path_util.Gofree5gcPath("free5gc/udmsslkey.log")
	udmPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.pem")
	udmKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.key")
	if sbi.Tls != nil {
		udmLogPath = path_util.Gofree5gcPath(sbi.Tls.Log)
		udmPemPath = path_util.Gofree5gcPath(sbi.Tls.Pem)
		udmKeyPath = path_util.Gofree5gcPath(sbi.Tls.Key)
	}
	addr := fmt.Sprintf("%s:%d", sbi.IPv4Addr, sbi.Port)

	self := udm_context.UDM_Self()
	udm_util.InitUDMContext(self)
	udm_context.UDM_Self().InitNFService(serviceName, config.Info.Version)

	proflie, err := udm_consumer.BuildNFInstance(self)
	if err != nil {
		logger.InitLog.Errorln(err.Error())
	} else {
		var newNrfUri string
		var err1 error
		newNrfUri, self.NfId, err1 = udm_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, proflie)
		if err1 != nil {
			logger.InitLog.Errorln(err1.Error())
		} else {
			self.NrfUri = newNrfUri
		}
	}

	go udm_handler.Handle()
	server, err := http2_util.NewServer(addr, udmLogPath, router)
	if err == nil && server != nil {
		initLog.Infoln(server.ListenAndServeTLS(udmPemPath, udmKeyPath))
	}
}

func (udm *UDM) Exec(c *cli.Context) error {

	//UDM.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("udmcfg"))
	args := udm.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./udm", args...)

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
			fmt.Printf("UDM Start error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
