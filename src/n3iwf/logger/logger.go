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
var NgapLog *logrus.Entry
var HandlerLog *logrus.Entry
var ContextLog *logrus.Entry
var IKELog *logrus.Entry
var RelayLog *logrus.Entry

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
			repopath := orgFilename
			repopath = strings.Replace(repopath, "/bin", "", 1)
			filename := strings.Replace(f.File, repopath, "", -1)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	AppLog = log.WithFields(logrus.Fields{"N3IWF": "app"})
	InitLog = log.WithFields(logrus.Fields{"N3IWF": "init"})
	NgapLog = log.WithFields(logrus.Fields{"N3IWF": "NGAP"})
	HandlerLog = log.WithFields(logrus.Fields{"N3IWF": "handler"})
	ContextLog = log.WithFields(logrus.Fields{"N3IWF": "context"})
	IKELog = log.WithFields(logrus.Fields{"N3IWF": "IKE"})
	RelayLog = log.WithFields(logrus.Fields{"N3IWF": "relay"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
