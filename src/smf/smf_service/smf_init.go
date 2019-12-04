package smf_service

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/app"
	"free5gc/src/smf/EventExposure"
	"free5gc/src/smf/PDUSession"
	"free5gc/src/smf/factory"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/smf_consumer"
	"free5gc/src/smf/smf_context"
	"free5gc/src/smf/smf_handler"
	"free5gc/src/smf/smf_pfcp/pfcp_message"
	"free5gc/src/smf/smf_pfcp/pfcp_udp"
	"free5gc/src/smf/smf_util"
	"net"
	"os/exec"
	"sync"
	"time"
)

type SMF struct{}

type (
	// Config information.
	Config struct {
		smfcfg string
	}
)

var config Config

var smfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "smfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*SMF) GetCliCmd() (flags []cli.Flag) {
	return smfCLi
}

func (*SMF) Initialize(c *cli.Context) {

	config = Config{
		smfcfg: c.String("smfcfg"),
	}

	if config.smfcfg != "" {
		factory.InitConfigFactory(config.smfcfg)
	} else {
		DefaultSmfConfigPath := path_util.Gofree5gcPath("free5gc/config/smfcfg.conf")
		factory.InitConfigFactory(DefaultSmfConfigPath)
	}

	initLog.Traceln("SMF debug level(string):", app.ContextSelf().Logger.SMF.DebugLevel)
	if app.ContextSelf().Logger.SMF.DebugLevel != "" {
		initLog.Infoln("SMF debug level(string):", app.ContextSelf().Logger.SMF.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.SMF.DebugLevel)
		if err != nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.SMF.ReportCaller)
}

func (smf *SMF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range smf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (smf *SMF) Start() {
	smf_context.InitSmfContext(&factory.SmfConfig)

	initLog.Infoln("Server started")
	router := gin.Default()

	err := smf_consumer.SendNFRegistration()

	if err != nil {
		retry_err := smf_consumer.RetrySendNFRegistration(10)
		if retry_err != nil {
			logger.InitLog.Errorln(retry_err)
			return
		}
	}

	for _, serviceName := range factory.SmfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serviceName) {
		case models.ServiceName_NSMF_PDUSESSION:
			PDUSession.AddService(router)
		case models.ServiceName_NSMF_EVENT_EXPOSURE:
			EventExposure.AddService(router)
		}
	}
	pfcp_udp.Run()

	for _, upf := range factory.SmfConfig.Configuration.UPF {
		if upf.Port == 0 {
			upf.Port = pfcpUdp.PFCP_PORT
		}
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", upf.Addr, upf.Port))
		if err != nil {
			logger.InitLog.Warnln("UPF addr error")
		}
		pfcp_message.SendPfcpAssociationSetupRequest(addr)
	}

	time.Sleep(1000 * time.Millisecond)

	go smf_handler.Handle()
	HTTPAddr := fmt.Sprintf("%s:%d", smf_context.SMF_Self().HTTPAddress, smf_context.SMF_Self().HTTPPort)
	server, _ := http2_util.NewServer(HTTPAddr, smf_util.SmfLogPath, router)

	initLog.Infoln(server.ListenAndServeTLS(smf_util.SmfPemPath, smf_util.SmfKeyPath))
}

func (smf *SMF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("smfcfg"))
	args := smf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./smf", args...)

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
			initLog.Errorf("SMF Start error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
