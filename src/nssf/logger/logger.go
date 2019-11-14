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
var ContextLog *logrus.Entry
var FactoryLog *logrus.Entry
var HandlerLog *logrus.Entry
var InitLog *logrus.Entry
var Nsselection *logrus.Entry
var Nssaiavailability *logrus.Entry
var Util *logrus.Entry

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

	AppLog = log.WithFields(logrus.Fields{"NSSF": "app"})
	ContextLog = log.WithFields(logrus.Fields{"NSSF": "context"})
	FactoryLog = log.WithFields(logrus.Fields{"NSSF": "factory"})
	HandlerLog = log.WithFields(logrus.Fields{"NSSF": "handler"})
	InitLog = log.WithFields(logrus.Fields{"NSSF": "init"})
	Nsselection = log.WithFields(logrus.Fields{"NSSF": "nsselection"})
	Nssaiavailability = log.WithFields(logrus.Fields{"NSSF": "nssaiavailability"})
	Util = log.WithFields(logrus.Fields{"NSSF": "util"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
