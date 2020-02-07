package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/logger_conf"
	"github.com/free5gc/logger_util"
)

var (
	log            *logrus.Logger
	AppLog         *logrus.Entry
	InitLog        *logrus.Entry
	CfgLog         *logrus.Entry
	HandlerLog     *logrus.Entry
	ManagementLog  *logrus.Entry
	AccessTokenLog *logrus.Entry
	DiscoveryLog   *logrus.Entry
	GinLog         *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	free5gcLogHook, err := logger_util.NewFileHook(logger_conf.Free5gcLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(free5gcLogHook)
	}

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"nrf.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "NRF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "NRF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "NRF", "category": "CFG"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "NRF", "category": "HDLR"})
	ManagementLog = log.WithFields(logrus.Fields{"component": "NRF", "category": "MGMT"})
	AccessTokenLog = log.WithFields(logrus.Fields{"component": "NRF", "category": "Token"})
	DiscoveryLog = log.WithFields(logrus.Fields{"component": "NRF", "category": "DSCV"})
	GinLog = log.WithFields(logrus.Fields{"component": "NRF", "category": "GIN"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
