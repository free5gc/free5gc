package udr_service

import (
	"bufio"
	"fmt"
	"free5gc/lib/http2_util"
	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/udr/DataRepository"
	"free5gc/src/udr/factory"
	"free5gc/src/udr/logger"
	"free5gc/src/udr/udr_consumer"
	"free5gc/src/udr/udr_handler"
	"free5gc/src/udr/udr_util"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type UDR struct{}

type (
	// Config information.
	Config struct {
		udrcfg string
	}
)

var config Config

var udrCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "udrcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*UDR) GetCliCmd() (flags []cli.Flag) {
	return udrCLi
}

func (*UDR) Initialize(c *cli.Context) {

	config = Config{
		udrcfg: c.String("udrcfg"),
	}

	if config.udrcfg != "" {
		factory.InitConfigFactory(config.udrcfg)
	} else {
		DefaultUdrConfigPath := path_util.Gofree5gcPath("free5gc/config/udrcfg.conf")
		factory.InitConfigFactory(DefaultUdrConfigPath)
	}

	initLog.Traceln("UDR debug level(string):", app.ContextSelf().Logger.UDR.DebugLevel)
	if app.ContextSelf().Logger.UDR.DebugLevel != "" {
		initLog.Infoln("UDR debug level(string):", app.ContextSelf().Logger.UDR.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.UDR.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.UDR.ReportCaller)

}

func (udr *UDR) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range udr.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (udr *UDR) Start() {
	// get config file info
	config := factory.UdrConfig
	sbi := config.Configuration.Sbi
	mongodb := config.Configuration.Mongodb
	nrfUri := config.Configuration.NrfUri

	initLog.Infof("UDR Config Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	// Connect to MongoDB
	DataRepository.SetMongoDB(mongodb.Name, mongodb.Url)

	initLog.Infoln("Server started")

	router := gin.Default()

	DataRepository.AddService(router)

	udrLogPath := udr_util.UdrLogPath
	udrPemPath := udr_util.UdrPemPath
	udrKeyPath := udr_util.UdrKeyPath
	if sbi.Tls != nil {
		udrLogPath = path_util.Gofree5gcPath(sbi.Tls.Log)
		udrPemPath = path_util.Gofree5gcPath(sbi.Tls.Pem)
		udrKeyPath = path_util.Gofree5gcPath(sbi.Tls.Key)
	}

	addr := fmt.Sprintf("%s:%d", sbi.IPv4Addr, sbi.Port)
	profile := udr_consumer.BuildNFInstance()
	var newNrfUri string
	var err error
	newNrfUri, profile.NfInstanceId, err = udr_consumer.SendRegisterNFInstance(nrfUri, profile.NfInstanceId, profile)
	if err == nil {
		config.Configuration.NrfUri = newNrfUri
	} else {
		initLog.Errorf("Send Register NFInstance Error[%s]", err.Error())
	}

	go udr_handler.Handle()
	server, err := http2_util.NewServer(addr, udrLogPath, router)
	if err == nil && server != nil {
		initLog.Infoln(server.ListenAndServeTLS(udrPemPath, udrKeyPath))
	}
}

func (udr *UDR) Exec(c *cli.Context) error {

	//UDR.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("udrcfg"))
	args := udr.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./udr", args...)

	udr.Initialize(c)

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
			fmt.Println("command.Start Fails!")
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
