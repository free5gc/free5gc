package amf_service

import (
	"bufio"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/amf/Communication"
	"free5gc/src/amf/EventExposure"
	"free5gc/src/amf/HttpCallback"
	"free5gc/src/amf/Location"
	"free5gc/src/amf/MT"
	"free5gc/src/amf/OAM"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/amf_ngap/ngap_message"
	"free5gc/src/amf/amf_ngap/ngap_sctp"
	"free5gc/src/amf/amf_producer/amf_producer_callback"
	"free5gc/src/amf/amf_util"
	"free5gc/src/amf/factory"
	"free5gc/src/amf/logger"
	"free5gc/src/app"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

type AMF struct{}

type (
	// Config information.
	Config struct {
		amfcfg string
	}
)

var config Config

var amfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "amfcfg",
		Usage: "amf config file",
	},
}

var initLog *logrus.Entry
var sctpListener *amf_ngap_sctp.SCTPListener

func init() {
	initLog = logger.InitLog
}

func (*AMF) GetCliCmd() (flags []cli.Flag) {
	return amfCLi
}

func (*AMF) Initialize(c *cli.Context) {

	config = Config{
		amfcfg: c.String("amfcfg"),
	}

	if config.amfcfg != "" {
		factory.InitConfigFactory(path_util.Gofree5gcPath(config.amfcfg))
	} else {
		factory.InitConfigFactory(amf_util.DefaultAmfConfigPath)
	}

	initLog.Traceln("AMF debug level(string):", app.ContextSelf().Logger.AMF.DebugLevel)
	if app.ContextSelf().Logger.AMF.DebugLevel != "" {
		initLog.Infoln("AMF debug level(string):", app.ContextSelf().Logger.AMF.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.AMF.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.AMF.ReportCaller)

}

func (amf *AMF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range amf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (amf *AMF) Start() {
	initLog.Infoln("Server started")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	Namf_Callback.AddService(router)
	Namf_OAM.AddService(router)
	for _, serviceName := range factory.AmfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serviceName) {
		case models.ServiceName_NAMF_COMM:
			Communication.AddService(router)
		case models.ServiceName_NAMF_EVTS:
			EventExposure.AddService(router)
		case models.ServiceName_NAMF_MT:
			Namf_MT.AddService(router)
		case models.ServiceName_NAMF_LOC:
			Namf_Location.AddService(router)
		}
	}

	self := amf_context.AMF_Self()
	amf_util.InitAmfContext(self)

	addr := fmt.Sprintf("%s:%d", self.HttpIPv4Address, self.HttpIpv4Port)

	for _, ngapAddr := range self.NgapIpList {
		sctpListener = amf_ngap_sctp.Server(ngapAddr)
	}
	go amf_handler.Handle()

	// Register to NRF
	profile, err := amf_consumer.BuildNFInstance(self)
	if err != nil {
		initLog.Error("Build AMF Profile Error")
	}

	_, self.NfId, _ = amf_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		amf.Terminate()
		os.Exit(0)
	}()

	server, err := http2_util.NewServer(addr, amf_util.AmfLogPath, router)
	if err == nil && server != nil {
		initLog.Infoln(server.ListenAndServeTLS(amf_util.AmfPemPath, amf_util.AmfKeyPath))
	} else {
		initLog.Errorf("Initialize http2 server failed: %+v", err)
	}
}

func (amf *AMF) Exec(c *cli.Context) error {

	//AMF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("amfcfg"))
	args := amf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./amf", args...)

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
			initLog.Errorf("AMF Start error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}

// Used in AMF planned removal procedure
func (amf *AMF) Terminate() {
	logger.InitLog.Infof("Terminating AMF...")
	amfSelf := amf_context.AMF_Self()

	// TODO: forward registered UE contexts to target AMF in the same AMF set if there is one

	// deregister with NRF
	problemDetails, err := amf_consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("[AMF] Deregister from NRF successfully")
	}

	// send AMF status indication to ran to notify ran that this AMF will be unavailable
	logger.InitLog.Infof("Send AMF Status Indication to Notify RANs due to AMF terminating")
	unavailableGuamiList := ngap_message.BuildUnavailableGUAMIList(amfSelf.ServedGuamiList)
	for _, ran := range amfSelf.AmfRanPool {
		ngap_message.SendAMFStatusIndication(ran, unavailableGuamiList)
	}

	logger.InitLog.Infof("Close SCTP server...")
	sctpListener.Close()
	logger.InitLog.Infof("SCTP server closed")

	amf_producer_callback.SendAmfStatusChangeNotify((string)(models.StatusChange_UNAVAILABLE), amfSelf.ServedGuamiList)
	logger.InitLog.Infof("AMF terminated")
}
