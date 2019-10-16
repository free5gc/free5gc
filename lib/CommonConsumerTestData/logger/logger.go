package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var log *logrus.Logger
var CommTestCommLog *logrus.Entry
var CommTestAmfLog *logrus.Entry
var CommTestSmfLog *logrus.Entry

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
	CommTestCommLog = log.WithFields(logrus.Fields{"CommonTest": "comm"})
	CommTestAmfLog = log.WithFields(logrus.Fields{"CommonTest": "amf"})
	CommTestSmfLog = log.WithFields(logrus.Fields{"CommonTest": "smf"})
}

func SetLogLevel(level logrus.Level) {
	CommTestCommLog.Infoln("set log level :", level)
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	CommTestCommLog.Infoln("set report call :", bool)
	log.SetReportCaller(bool)
}
