package webui_service

import (
	"bufio"
	"fmt"
	"github.com/gin-contrib/cors"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/udr/factory"
	"free5gc/webconsole/backend/WebUI"
	"free5gc/webconsole/backend/logger"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type WEBUI struct{}

type (
	// Config information.
	Config struct {
		webuicfg string
	}
)

var config Config

var webuiCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "webuicfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*WEBUI) GetCliCmd() (flags []cli.Flag) {
	return webuiCLi
}

func (*WEBUI) Initialize(c *cli.Context) {

	config = Config{
		webuicfg: c.String("webuicfg"),
	}

	if config.webuicfg != "" {
		factory.InitConfigFactory(config.webuicfg)
	} else {
		DefaultUdrConfigPath := path_util.Gofree5gcPath("free5gc/config/udrcfg.conf")
		factory.InitConfigFactory(DefaultUdrConfigPath)
	}

	initLog.Traceln("WEBUI debug level(string):", app.ContextSelf().Logger.WEBUI.DebugLevel)
	if app.ContextSelf().Logger.WEBUI.DebugLevel != "" {
		initLog.Infoln("WEBUI debug level(string):", app.ContextSelf().Logger.WEBUI.DebugLevel)
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.WEBUI.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(app.ContextSelf().Logger.WEBUI.ReportCaller)

}

func (webui *WEBUI) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range webui.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (webui *WEBUI) Start() {
	// get config file info from UdrConfig
	mongodb := factory.UdrConfig.Configuration.Mongodb

	// Connect to MongoDB
	MongoDBLibrary.SetMongoDB(mongodb.Name, mongodb.Url)

	initLog.Infoln("Server started")

	router := WebUI.NewRouter()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	// router.Use(cors.Default())

	router.NoRoute(ReturnPublic())

	initLog.Infoln(router.Run(":5000"))
}

func (webui *WEBUI) Exec(c *cli.Context) error {

	//WEBUI.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("webuicfg"))
	args := webui.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./webui", args...)

	webui.Initialize(c)

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
