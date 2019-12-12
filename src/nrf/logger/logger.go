package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var AppLog *logrus.Entry
var InitLog *logrus.Entry
var HandlerLog *logrus.Entry
var nrf_handler *logrus.Entry
var ManagementLog *logrus.Entry

func init() {
	log = logrus.New()
	log.SetReportCaller(true)

	log.Formatter = &logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			orgFilename, _ := os.Getwd()
			log.Traceln("orgFilename", orgFilename)
			repopath := orgFilename
			repopath = strings.Replace(repopath, "/bin", "", 1)
			filename := strings.Replace(f.File, repopath, "", -1)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	AppLog = log.WithFields(logrus.Fields{"NRF": "app"})
	InitLog = log.WithFields(logrus.Fields{"NRF": "init"})
	HandlerLog = log.WithFields(logrus.Fields{"NRF": "HandlerLog"})
	ManagementLog = log.WithFields(logrus.Fields{"NRF": "ManagementLog"})
	nrf_handler = log.WithFields(logrus.Fields{"NRF": "nrf_handler"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
