package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var log *logrus.Logger
var AppLog *logrus.Entry
var InitLog *logrus.Entry
var GsmLog *logrus.Entry
var PfcpLog *logrus.Entry
var PduSessLog *logrus.Entry
var CtxLog *logrus.Entry

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

	AppLog = log.WithFields(logrus.Fields{"SMF": "app"})
	InitLog = log.WithFields(logrus.Fields{"SMF": "init"})
	PfcpLog = log.WithFields(logrus.Fields{"SMF": "pfcp"})
	PduSessLog = log.WithFields(logrus.Fields{"SMF": "pdu_session"})
	GsmLog = log.WithFields(logrus.Fields{"SMF": "GSM"})
	CtxLog = log.WithFields(logrus.Fields{"SMF": "Context"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
